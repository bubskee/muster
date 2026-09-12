package engine_test

import (
	"testing"

	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/state"
)

func TestEventSequenceMutatesStateDeterministically(t *testing.T) {
	s := state.New("dataset:untrusted")

	events := []engine.Event{
		{
			ID: "worker-rce",
			Preconditions: []engine.Condition{
				{Fact: "dataset:untrusted", Op: engine.ConditionPresent},
			},
			Effects: []engine.Effect{
				{Fact: "access:worker", Op: engine.EffectAdd},
			},
		},
		{
			ID: "read-service-account-token",
			Preconditions: []engine.Condition{
				{Fact: "access:worker", Op: engine.ConditionPresent},
			},
			Effects: []engine.Effect{
				{Fact: "credential:k8s-service-account", Op: engine.EffectAdd},
			},
		},
	}

	for _, event := range events {
		if !event.Apply(s) {
			t.Fatalf("event %q unexpectedly failed to apply", event.ID)
		}
	}

	want := []string{
		"access:worker",
		"credential:k8s-service-account",
		"dataset:untrusted",
	}

	got := s.Facts()

	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestEventDoesNotApplyWhenPreconditionFails(t *testing.T) {
	s := state.New()

	event := engine.Event{
		ID: "node-access",
		Preconditions: []engine.Condition{
			{Fact: "access:worker", Op: engine.ConditionPresent},
		},
		Effects: []engine.Effect{
			{Fact: "access:node", Op: engine.EffectAdd},
		},
	}

	if event.Apply(s) {
		t.Fatal("event applied despite unsatisfied precondition")
	}

	if s.Has("access:node") {
		t.Fatal("failed event mutated state")
	}
}
