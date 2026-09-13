package config_test

import (
	"testing"

	"github.com/bubskee/muster/config"
	"github.com/bubskee/muster/engine"
)

func TestParseScenario(t *testing.T) {
	scenario, err := config.ParseScenario([]byte(`
id: demo
initial_facts:
  - seed

events:
  - id: first
    requires:
      present: [seed]
      absent: [stopped]
    effects:
      add: [access]
      remove: [seed]
`))
	if err != nil {
		t.Fatal(err)
	}

	if scenario.ID != "demo" {
		t.Fatalf("id = %q, want demo", scenario.ID)
	}
	if len(scenario.Events) != 1 {
		t.Fatalf("events = %d, want 1", len(scenario.Events))
	}

	event := scenario.Events[0]
	if len(event.Preconditions) != 2 {
		t.Fatalf("preconditions = %d, want 2", len(event.Preconditions))
	}
	if event.Preconditions[0].Op != engine.ConditionPresent || event.Preconditions[1].Op != engine.ConditionAbsent {
		t.Fatalf("unexpected preconditions: %+v", event.Preconditions)
	}
	if len(event.Effects) != 2 || event.Effects[0].Op != engine.EffectAdd || event.Effects[1].Op != engine.EffectRemove {
		t.Fatalf("unexpected effects: %+v", event.Effects)
	}
}

func TestParseControlSet(t *testing.T) {
	set, err := config.ParseControlSet([]byte(`
id: demo-controls
controls:
  - id: detector
    role: vedette
    events: [credential-theft]
    substrate: worker-telemetry
    suppressed_by:
      present: [telemetry:dead]
    action: observe
    effects:
      add: [alert:credential-theft]
`))
	if err != nil {
		t.Fatal(err)
	}

	if set.ID != "demo-controls" || len(set.Controls) != 1 {
		t.Fatalf("set = %+v", set)
	}

	control, ok := set.Controls[0].(engine.EventControl)
	if !ok {
		t.Fatalf("control type = %T, want engine.EventControl", set.Controls[0])
	}
	if control.Role != engine.Vedette || control.Action != engine.ActionObserve {
		t.Fatalf("control = %+v", control)
	}
	if len(control.SuppressedBy) != 1 || control.SuppressedBy[0].Fact != "telemetry:dead" {
		t.Fatalf("suppressed_by = %+v", control.SuppressedBy)
	}
}

func TestParseRejectsUnknownField(t *testing.T) {
	_, err := config.ParseScenario([]byte(`
id: demo
initial_facts: []
banana: true
events: []
`))
	if err == nil {
		t.Fatal("expected unknown YAML field to fail")
	}
}

func TestParseRejectsUnknownAction(t *testing.T) {
	_, err := config.ParseControlSet([]byte(`
id: demo
controls:
  - id: nope
    events: [event]
    action: teleport
`))
	if err == nil {
		t.Fatal("expected unknown action to fail")
	}
}
