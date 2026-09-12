package engine_test

import (
	"reflect"
	"testing"

	"github.com/bubskee/muster/engine"
)

func TestReplayProducesDeterministicTraceAndTerminalState(t *testing.T) {
	scenario := engine.Scenario{
		ID: "test-replay",
		InitialFacts: []string{
			"dataset:untrusted",
		},
		Events: []engine.Event{
			{
				ID: "worker-rce",
				Preconditions: []engine.Condition{
					{
						Fact: "dataset:untrusted",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "access:worker",
						Op:   engine.EffectAdd,
					},
				},
			},
			{
				ID: "node-access",
				Preconditions: []engine.Condition{
					{
						Fact: "credential:node",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "access:node",
						Op:   engine.EffectAdd,
					},
				},
			},
			{
				ID: "service-account-token-read",
				Preconditions: []engine.Condition{
					{
						Fact: "access:worker",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "credential:k8s-service-account",
						Op:   engine.EffectAdd,
					},
				},
			},
		},
	}

	first := engine.Replay(scenario)
	second := engine.Replay(scenario)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf(
			"replay was not deterministic:\nfirst:  %#v\nsecond: %#v",
			first,
			second,
		)
	}

	wantTerminal := []string{
		"access:worker",
		"credential:k8s-service-account",
		"dataset:untrusted",
	}

	if !reflect.DeepEqual(first.TerminalState, wantTerminal) {
		t.Fatalf(
			"terminal state = %v, want %v",
			first.TerminalState,
			wantTerminal,
		)
	}

	if len(first.Trace) != 3 {
		t.Fatalf("trace length = %d, want 3", len(first.Trace))
	}

	if first.Trace[0].Status != engine.EventApplied {
		t.Errorf("worker-rce status = %q, want applied", first.Trace[0].Status)
	}

	if first.Trace[1].Status != engine.EventSkipped {
		t.Errorf("node-access status = %q, want skipped", first.Trace[1].Status)
	}

	if len(first.Trace[1].Unsatisfied) != 1 {
		t.Fatalf(
			"node-access unsatisfied conditions = %d, want 1",
			len(first.Trace[1].Unsatisfied),
		)
	}

	if first.Trace[2].Status != engine.EventApplied {
		t.Errorf(
			"service-account-token-read status = %q, want applied",
			first.Trace[2].Status,
		)
	}
}
