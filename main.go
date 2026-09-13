package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/bubskee/muster/config"
	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/experiment"
)

type replayOutput struct {
	ControlSetID string           `json:"control_set_id"`
	Result       engine.RunResult `json:"result"`
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "replay":
		if err := replay(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "muster:", err)
			os.Exit(1)
		}
	case "experiment":
		if err := runExperiment(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "muster:", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "muster: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func replay(args []string) error {
	flags := flag.NewFlagSet("replay", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	scenarioPath := flags.String("scenario", "", "scenario YAML file")
	controlsPath := flags.String("controls", "", "control-set YAML file")
	asJSON := flags.Bool("json", false, "emit replay result as JSON")

	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}
	if *scenarioPath == "" {
		return fmt.Errorf("--scenario is required")
	}
	if *controlsPath == "" {
		return fmt.Errorf("--controls is required")
	}

	scenario, err := config.LoadScenario(*scenarioPath)
	if err != nil {
		return err
	}
	controls, err := config.LoadControlSet(*controlsPath)
	if err != nil {
		return err
	}

	result := engine.Replay(scenario, controls.Controls...)

	if *asJSON {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(replayOutput{
			ControlSetID: controls.ID,
			Result:       result,
		})
	}

	fmt.Printf("control_set: %s\n", controls.ID)
	fmt.Print(engine.Summarize(result).String())
	return nil
}

func runExperiment(args []string) error {
	flags := flag.NewFlagSet("experiment", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	path := flags.String("file", "", "experiment YAML file")
	format := flags.String("format", "markdown", "output format: markdown, json, or csv")

	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}
	if *path == "" {
		return fmt.Errorf("--file is required")
	}

	result, err := experiment.ExecuteFile(*path)
	if err != nil {
		return err
	}

	switch *format {
	case "markdown", "md":
		return experiment.WriteMarkdown(os.Stdout, result)
	case "csv":
		return experiment.WriteCSV(os.Stdout, result)
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	default:
		return fmt.Errorf("unknown --format %q; want markdown, json, or csv", *format)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  muster replay --scenario SCENARIO.yaml --controls CONTROLS.yaml [--json]
  muster experiment --file EXPERIMENT.yaml [--format markdown|json|csv]

examples:
  muster replay --scenario examples/gate2-between.yaml --controls examples/gate2-broad-correlated.yaml
  muster replay --scenario examples/gate2-between.yaml --controls examples/gate2-broad-correlated.yaml --json
  muster experiment --file experiments/hf-july-2026.yaml
  muster experiment --file experiments/hf-july-2026.yaml --format json`)
}
