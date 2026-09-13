package engine_test

import (
	"testing"

	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/state"
)

func TestEventControlDisposition(t *testing.T) {
	event := engine.Event{
		ID: "credential-access",
	}

	control := engine.EventControl{
		ID:       "soc-review",
		EventIDs: []string{"credential-access"},
		Requires: []engine.Condition{
			{
				Fact: "alert:credential-access",
				Op:   engine.ConditionPresent,
			},
		},
		Action: engine.ActionReview,
	}

	t.Run("unmatched", func(t *testing.T) {
		otherEvent := engine.Event{
			ID: "unrelated-event",
		}

		result := control.Evaluate(otherEvent, state.New())

		if result.Matched {
			t.Fatal("control unexpectedly matched event")
		}

		if result.Disposition != engine.DispositionUnmatched {
			t.Fatalf(
				"disposition = %q, want %q",
				result.Disposition,
				engine.DispositionUnmatched,
			)
		}

		if result.Action != engine.ActionNone {
			t.Fatalf(
				"action = %q, want %q",
				result.Action,
				engine.ActionNone,
			)
		}
	})

	t.Run("waiting for required state", func(t *testing.T) {
		result := control.Evaluate(event, state.New())

		if !result.Matched {
			t.Fatal("control should match event")
		}

		if result.Disposition != engine.DispositionWaiting {
			t.Fatalf(
				"disposition = %q, want %q",
				result.Disposition,
				engine.DispositionWaiting,
			)
		}

		if result.Reason != "missing:alert:credential-access" {
			t.Fatalf(
				"reason = %q, want %q",
				result.Reason,
				"missing:alert:credential-access",
			)
		}

		if result.Action != engine.ActionNone {
			t.Fatalf(
				"action = %q, want %q",
				result.Action,
				engine.ActionNone,
			)
		}
	})

	t.Run("ready", func(t *testing.T) {
		result := control.Evaluate(
			event,
			state.New("alert:credential-access"),
		)

		if !result.Matched {
			t.Fatal("control did not match event")
		}

		if result.Disposition != engine.DispositionReady {
			t.Fatalf(
				"disposition = %q, want %q",
				result.Disposition,
				engine.DispositionReady,
			)
		}

		if result.Reason != "" {
			t.Fatalf(
				"reason = %q, want empty",
				result.Reason,
			)
		}

		if result.Action != engine.ActionReview {
			t.Fatalf(
				"action = %q, want %q",
				result.Action,
				engine.ActionReview,
			)
		}
	})
}

func TestEventControlSuppressedBySubstrate(t *testing.T) {
	event := engine.Event{
		ID: "credential-access",
	}

	control := engine.EventControl{
		ID:        "credential-detector",
		Role:      engine.Vedette,
		EventIDs:  []string{"credential-access"},
		Substrate: "worker-telemetry",
		SuppressedBy: []engine.Condition{
			{
				Fact: "compromised:worker-telemetry",
				Op:   engine.ConditionPresent,
			},
		},
		Action: engine.ActionObserve,
		Effects: []engine.Effect{
			{
				Fact: "alert:credential-access",
				Op:   engine.EffectAdd,
			},
		},
	}

	t.Run("healthy substrate", func(t *testing.T) {
		result := control.Evaluate(event, state.New())

		if result.Disposition != engine.DispositionReady {
			t.Fatalf(
				"disposition = %q, want %q",
				result.Disposition,
				engine.DispositionReady,
			)
		}

		if result.Action != engine.ActionObserve {
			t.Fatalf(
				"action = %q, want %q",
				result.Action,
				engine.ActionObserve,
			)
		}
	})

	t.Run("compromised substrate", func(t *testing.T) {
		result := control.Evaluate(
			event,
			state.New("compromised:worker-telemetry"),
		)

		if result.Disposition != engine.DispositionSuppressed {
			t.Fatalf(
				"disposition = %q, want %q",
				result.Disposition,
				engine.DispositionSuppressed,
			)
		}

		if result.Reason != "compromised:worker-telemetry" {
			t.Fatalf(
				"reason = %q, want compromised substrate",
				result.Reason,
			)
		}

		if result.Action != engine.ActionNone {
			t.Fatalf(
				"action = %q, want %q",
				result.Action,
				engine.ActionNone,
			)
		}
	})
}

func TestWaitingPrecedesSuppression(t *testing.T) {
	event := engine.Event{
		ID: "credential-access",
	}

	control := engine.EventControl{
		ID:       "soc-review",
		EventIDs: []string{"credential-access"},
		Requires: []engine.Condition{
			{
				Fact: "alert:credential-access",
				Op:   engine.ConditionPresent,
			},
		},
		Substrate: "soc-queue",
		SuppressedBy: []engine.Condition{
			{
				Fact: "overloaded:soc-queue",
				Op:   engine.ConditionPresent,
			},
		},
		Action: engine.ActionReview,
	}

	result := control.Evaluate(
		event,
		state.New("overloaded:soc-queue"),
	)

	if result.Disposition != engine.DispositionWaiting {
		t.Fatalf(
			"disposition = %q, want %q",
			result.Disposition,
			engine.DispositionWaiting,
		)
	}

	if result.Reason != "missing:alert:credential-access" {
		t.Fatalf(
			"reason = %q, want missing alert",
			result.Reason,
		)
	}
}
