package experiment

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteFileResolvesRelativePathsAndReportsMetrics(t *testing.T) {
	dir := t.TempDir()

	write := func(name, contents string) {
		t.Helper()
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("scenario.yaml", `
id: demo
initial_facts:
  - start
events:
  - id: breach
    requires:
      present:
        - start
    effects:
      add:
        - attacker:inside
`)
	write("controls.yaml", `
id: none
controls: []
`)
	write("experiment.yaml", `
id: demo-experiment
metrics:
  - id: inside
    fact: attacker:inside
runs:
  - id: baseline
    label: Baseline
    scenario: scenario.yaml
    controls: controls.yaml
`)

	result, err := ExecuteFile(filepath.Join(dir, "experiment.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Runs) != 1 {
		t.Fatalf("runs = %d, want 1", len(result.Runs))
	}
	if !result.Runs[0].Metrics["inside"] {
		t.Fatal("inside metric = false, want true")
	}

	var markdown bytes.Buffer
	if err := WriteMarkdown(&markdown, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(markdown.String(), "| Baseline | none | breach") {
		t.Fatalf("markdown missing run row:\n%s", markdown.String())
	}

	var csv bytes.Buffer
	if err := WriteCSV(&csv, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(csv.String(), "baseline,Baseline,demo,none,breach") {
		t.Fatalf("csv missing run row:\n%s", csv.String())
	}
}

func TestLoadRejectsDuplicateMetricIDs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "experiment.yaml")
	data := `
id: demo
metrics:
  - id: outcome
    fact: a
  - id: outcome
    fact: b
runs:
  - id: run
    scenario: scenario.yaml
    controls: controls.yaml
`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil || !strings.Contains(err.Error(), "duplicate metric id") {
		t.Fatalf("err = %v, want duplicate metric id", err)
	}
}
