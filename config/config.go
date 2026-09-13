package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bubskee/muster/engine"
	"gopkg.in/yaml.v3"
)

type conditionSet struct {
	Present []string `yaml:"present"`
	Absent  []string `yaml:"absent"`
}

type effectSet struct {
	Add    []string `yaml:"add"`
	Remove []string `yaml:"remove"`
}

type eventSpec struct {
	ID       string       `yaml:"id"`
	Requires conditionSet `yaml:"requires"`
	Effects  effectSet    `yaml:"effects"`
}

type scenarioSpec struct {
	ID           string      `yaml:"id"`
	InitialFacts []string    `yaml:"initial_facts"`
	Events       []eventSpec `yaml:"events"`
}

type controlSpec struct {
	ID           string       `yaml:"id"`
	Role         string       `yaml:"role"`
	Events       []string     `yaml:"events"`
	Requires     conditionSet `yaml:"requires"`
	Substrate    string       `yaml:"substrate"`
	SuppressedBy conditionSet `yaml:"suppressed_by"`
	Action       string       `yaml:"action"`
	Effects      effectSet    `yaml:"effects"`
}

type controlSetSpec struct {
	ID       string        `yaml:"id"`
	Controls []controlSpec `yaml:"controls"`
}

// ControlSet is a named collection of controls loaded from one YAML file.
// The ID is presentation metadata; Replay receives only Controls.
type ControlSet struct {
	ID       string
	Controls []engine.Control
}

func LoadScenario(path string) (engine.Scenario, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return engine.Scenario{}, fmt.Errorf("read scenario %q: %w", path, err)
	}

	scenario, err := ParseScenario(data)
	if err != nil {
		return engine.Scenario{}, fmt.Errorf("parse scenario %q: %w", path, err)
	}

	return scenario, nil
}

func ParseScenario(data []byte) (engine.Scenario, error) {
	var spec scenarioSpec
	if err := decodeStrict(data, &spec); err != nil {
		return engine.Scenario{}, err
	}

	if strings.TrimSpace(spec.ID) == "" {
		return engine.Scenario{}, fmt.Errorf("scenario id is required")
	}

	if err := validateFacts("initial_facts", spec.InitialFacts); err != nil {
		return engine.Scenario{}, err
	}

	seen := make(map[string]struct{}, len(spec.Events))
	events := make([]engine.Event, 0, len(spec.Events))

	for i, event := range spec.Events {
		where := fmt.Sprintf("events[%d]", i)
		if strings.TrimSpace(event.ID) == "" {
			return engine.Scenario{}, fmt.Errorf("%s.id is required", where)
		}
		if _, ok := seen[event.ID]; ok {
			return engine.Scenario{}, fmt.Errorf("duplicate event id %q", event.ID)
		}
		seen[event.ID] = struct{}{}

		requires, err := compileConditions(where+".requires", event.Requires)
		if err != nil {
			return engine.Scenario{}, err
		}
		effects, err := compileEffects(where+".effects", event.Effects)
		if err != nil {
			return engine.Scenario{}, err
		}

		events = append(events, engine.Event{
			ID:            event.ID,
			Preconditions: requires,
			Effects:       effects,
		})
	}

	return engine.Scenario{
		ID:           spec.ID,
		InitialFacts: append([]string(nil), spec.InitialFacts...),
		Events:       events,
	}, nil
}

func LoadControlSet(path string) (ControlSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ControlSet{}, fmt.Errorf("read controls %q: %w", path, err)
	}

	set, err := ParseControlSet(data)
	if err != nil {
		return ControlSet{}, fmt.Errorf("parse controls %q: %w", path, err)
	}

	return set, nil
}

func ParseControlSet(data []byte) (ControlSet, error) {
	var spec controlSetSpec
	if err := decodeStrict(data, &spec); err != nil {
		return ControlSet{}, err
	}

	if strings.TrimSpace(spec.ID) == "" {
		return ControlSet{}, fmt.Errorf("control-set id is required")
	}

	seen := make(map[string]struct{}, len(spec.Controls))
	controls := make([]engine.Control, 0, len(spec.Controls))

	for i, control := range spec.Controls {
		where := fmt.Sprintf("controls[%d]", i)
		if strings.TrimSpace(control.ID) == "" {
			return ControlSet{}, fmt.Errorf("%s.id is required", where)
		}
		if _, ok := seen[control.ID]; ok {
			return ControlSet{}, fmt.Errorf("duplicate control id %q", control.ID)
		}
		seen[control.ID] = struct{}{}

		if len(control.Events) == 0 {
			return ControlSet{}, fmt.Errorf("%s.events must not be empty", where)
		}
		if err := validateFacts(where+".events", control.Events); err != nil {
			return ControlSet{}, err
		}

		role, err := parseRole(control.Role)
		if err != nil {
			return ControlSet{}, fmt.Errorf("%s.role: %w", where, err)
		}
		action, err := parseAction(control.Action)
		if err != nil {
			return ControlSet{}, fmt.Errorf("%s.action: %w", where, err)
		}
		requires, err := compileConditions(where+".requires", control.Requires)
		if err != nil {
			return ControlSet{}, err
		}
		suppressedBy, err := compileConditions(where+".suppressed_by", control.SuppressedBy)
		if err != nil {
			return ControlSet{}, err
		}
		effects, err := compileEffects(where+".effects", control.Effects)
		if err != nil {
			return ControlSet{}, err
		}

		controls = append(controls, engine.EventControl{
			ID:           control.ID,
			Role:         role,
			EventIDs:     append([]string(nil), control.Events...),
			Requires:     requires,
			Substrate:    control.Substrate,
			SuppressedBy: suppressedBy,
			Action:       action,
			Effects:      effects,
		})
	}

	return ControlSet{ID: spec.ID, Controls: controls}, nil
}

func decodeStrict(data []byte, out any) error {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("decode YAML: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("multiple YAML documents are not supported")
		}
		return fmt.Errorf("decode trailing YAML: %w", err)
	}

	return nil
}

func compileConditions(where string, spec conditionSet) ([]engine.Condition, error) {
	if err := validateFacts(where+".present", spec.Present); err != nil {
		return nil, err
	}
	if err := validateFacts(where+".absent", spec.Absent); err != nil {
		return nil, err
	}

	conditions := make([]engine.Condition, 0, len(spec.Present)+len(spec.Absent))
	for _, fact := range spec.Present {
		conditions = append(conditions, engine.Condition{Fact: fact, Op: engine.ConditionPresent})
	}
	for _, fact := range spec.Absent {
		conditions = append(conditions, engine.Condition{Fact: fact, Op: engine.ConditionAbsent})
	}
	return conditions, nil
}

func compileEffects(where string, spec effectSet) ([]engine.Effect, error) {
	if err := validateFacts(where+".add", spec.Add); err != nil {
		return nil, err
	}
	if err := validateFacts(where+".remove", spec.Remove); err != nil {
		return nil, err
	}

	effects := make([]engine.Effect, 0, len(spec.Add)+len(spec.Remove))
	for _, fact := range spec.Add {
		effects = append(effects, engine.Effect{Fact: fact, Op: engine.EffectAdd})
	}
	for _, fact := range spec.Remove {
		effects = append(effects, engine.Effect{Fact: fact, Op: engine.EffectRemove})
	}
	return effects, nil
}

func validateFacts(where string, facts []string) error {
	for i, fact := range facts {
		if strings.TrimSpace(fact) == "" {
			return fmt.Errorf("%s[%d] must not be empty", where, i)
		}
	}
	return nil
}

func parseRole(value string) (engine.Role, error) {
	role := engine.Role(value)
	switch role {
	case "", engine.Vedette, engine.Picket, engine.Reserve:
		return role, nil
	default:
		return "", fmt.Errorf("unknown role %q", value)
	}
}

func parseAction(value string) (engine.ControlAction, error) {
	action := engine.ControlAction(value)
	switch action {
	case engine.ActionObserve,
		engine.ActionBlock,
		engine.ActionReview,
		engine.ActionCritical,
		engine.ActionEscalate,
		engine.ActionRespond:
		return action, nil
	default:
		return engine.ActionNone, fmt.Errorf("unknown action %q", value)
	}
}
