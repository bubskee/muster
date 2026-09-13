package experiment

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bubskee/muster/config"
	"github.com/bubskee/muster/engine"
	"gopkg.in/yaml.v3"
)

type MetricSpec struct {
	ID   string `yaml:"id"`
	Fact string `yaml:"fact"`
}

type RunSpec struct {
	ID       string `yaml:"id"`
	Label    string `yaml:"label"`
	Scenario string `yaml:"scenario"`
	Controls string `yaml:"controls"`
}

type Spec struct {
	ID      string       `yaml:"id"`
	Metrics []MetricSpec `yaml:"metrics"`
	Runs    []RunSpec    `yaml:"runs"`
}

type RunResult struct {
	ID           string            `json:"id"`
	Label        string            `json:"label,omitempty"`
	ScenarioPath string            `json:"scenario_path"`
	ControlsPath string            `json:"controls_path"`
	ControlSetID string            `json:"control_set_id"`
	Metrics      map[string]bool   `json:"metrics"`
	Summary      engine.RunSummary `json:"summary"`
	Result       engine.RunResult  `json:"result"`
}

type Result struct {
	ExperimentID string       `json:"experiment_id"`
	Metrics      []MetricSpec `json:"metrics"`
	Runs         []RunResult  `json:"runs"`
}

func Load(path string) (Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Spec{}, fmt.Errorf("read experiment %q: %w", path, err)
	}

	var spec Spec
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&spec); err != nil {
		return Spec{}, fmt.Errorf("decode experiment YAML: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return Spec{}, fmt.Errorf("multiple YAML documents are not supported")
		}
		return Spec{}, fmt.Errorf("decode trailing experiment YAML: %w", err)
	}

	if err := validate(spec); err != nil {
		return Spec{}, err
	}

	return spec, nil
}

func ExecuteFile(path string) (Result, error) {
	spec, err := Load(path)
	if err != nil {
		return Result{}, err
	}

	baseDir := filepath.Dir(path)
	out := Result{
		ExperimentID: spec.ID,
		Metrics:      append([]MetricSpec(nil), spec.Metrics...),
		Runs:         make([]RunResult, 0, len(spec.Runs)),
	}

	for _, run := range spec.Runs {
		scenarioPath := resolve(baseDir, run.Scenario)
		controlsPath := resolve(baseDir, run.Controls)

		scenario, err := config.LoadScenario(scenarioPath)
		if err != nil {
			return Result{}, fmt.Errorf("run %q: %w", run.ID, err)
		}
		controls, err := config.LoadControlSet(controlsPath)
		if err != nil {
			return Result{}, fmt.Errorf("run %q: %w", run.ID, err)
		}

		replay := engine.Replay(scenario, controls.Controls...)
		metrics := make(map[string]bool, len(spec.Metrics))
		for _, metric := range spec.Metrics {
			metrics[metric.ID] = contains(replay.TerminalState, metric.Fact)
		}

		out.Runs = append(out.Runs, RunResult{
			ID:           run.ID,
			Label:        run.Label,
			ScenarioPath: run.Scenario,
			ControlsPath: run.Controls,
			ControlSetID: controls.ID,
			Metrics:      metrics,
			Summary:      engine.Summarize(replay),
			Result:       replay,
		})
	}

	return out, nil
}

func WriteMarkdown(w io.Writer, result Result) error {
	if _, err := fmt.Fprintf(w, "# Experiment: %s\n\n", markdownCell(result.ExperimentID)); err != nil {
		return err
	}

	headers := []string{"run", "controls", "furthest applied", "first observation", "first block"}
	for _, metric := range result.Metrics {
		headers = append(headers, metric.ID)
	}

	if err := writeMarkdownRow(w, headers); err != nil {
		return err
	}
	separator := make([]string, len(headers))
	for i := range separator {
		separator[i] = "---"
	}
	if err := writeMarkdownRow(w, separator); err != nil {
		return err
	}

	for _, run := range result.Runs {
		name := run.ID
		if run.Label != "" {
			name = run.Label
		}
		row := []string{
			name,
			run.ControlSetID,
			run.Summary.FurthestApplied,
			run.Summary.FirstObservation,
			run.Summary.FirstBlock,
		}
		for _, metric := range result.Metrics {
			row = append(row, yesNo(run.Metrics[metric.ID]))
		}
		if err := writeMarkdownRow(w, row); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w, "\nMetrics:"); err != nil {
		return err
	}
	for _, metric := range result.Metrics {
		if _, err := fmt.Fprintf(w, "- `%s` → `%s`\n", metric.ID, metric.Fact); err != nil {
			return err
		}
	}
	return nil
}

func WriteCSV(w io.Writer, result Result) error {
	writer := csv.NewWriter(w)

	header := []string{"run_id", "label", "scenario", "control_set", "furthest_applied", "first_observation", "first_block"}
	for _, metric := range result.Metrics {
		header = append(header, metric.ID)
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, run := range result.Runs {
		row := []string{
			run.ID,
			run.Label,
			run.Result.ScenarioID,
			run.ControlSetID,
			run.Summary.FurthestApplied,
			run.Summary.FirstObservation,
			run.Summary.FirstBlock,
		}
		for _, metric := range result.Metrics {
			row = append(row, fmt.Sprintf("%t", run.Metrics[metric.ID]))
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}

func validate(spec Spec) error {
	if strings.TrimSpace(spec.ID) == "" {
		return fmt.Errorf("experiment id is required")
	}
	if len(spec.Metrics) == 0 {
		return fmt.Errorf("experiment must define at least one metric")
	}
	if len(spec.Runs) == 0 {
		return fmt.Errorf("experiment must define at least one run")
	}

	metricIDs := make(map[string]struct{}, len(spec.Metrics))
	for i, metric := range spec.Metrics {
		where := fmt.Sprintf("metrics[%d]", i)
		if strings.TrimSpace(metric.ID) == "" {
			return fmt.Errorf("%s.id is required", where)
		}
		if strings.TrimSpace(metric.Fact) == "" {
			return fmt.Errorf("%s.fact is required", where)
		}
		if _, ok := metricIDs[metric.ID]; ok {
			return fmt.Errorf("duplicate metric id %q", metric.ID)
		}
		metricIDs[metric.ID] = struct{}{}
	}

	runIDs := make(map[string]struct{}, len(spec.Runs))
	for i, run := range spec.Runs {
		where := fmt.Sprintf("runs[%d]", i)
		if strings.TrimSpace(run.ID) == "" {
			return fmt.Errorf("%s.id is required", where)
		}
		if _, ok := runIDs[run.ID]; ok {
			return fmt.Errorf("duplicate run id %q", run.ID)
		}
		runIDs[run.ID] = struct{}{}
		if strings.TrimSpace(run.Scenario) == "" {
			return fmt.Errorf("%s.scenario is required", where)
		}
		if strings.TrimSpace(run.Controls) == "" {
			return fmt.Errorf("%s.controls is required", where)
		}
	}
	return nil
}

func resolve(baseDir, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(baseDir, path))
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func writeMarkdownRow(w io.Writer, values []string) error {
	for i := range values {
		values[i] = markdownCell(values[i])
	}
	_, err := fmt.Fprintf(w, "| %s |\n", strings.Join(values, " | "))
	return err
}

func markdownCell(value string) string {
	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.ReplaceAll(value, "\n", " ")
	return value
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "—"
}
