package engine_test

import (
	"testing"

	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/state"
)

func TestEventControlRequiresState(t *testing.T) {
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

	t.Run("requirements satisfied", func(t *testing.T) {
		s := state.New("alert:credential-access")

		result := control.Evaluate(event, s)

		if !result.Matched {
			t.Fatal("control did not match event")
		}

		if result.Action != engine.ActionReview {
			t.Fatalf(
				"action = %q, want %q",
				result.Action,
				engine.ActionReview,
			)
		}
	})

	t.Run("requirements missing", func(t *testing.T) {
		s := state.New()

		result := control.Evaluate(event, s)

		if !result.Matched {
			t.Fatal("control should still match the event")
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
