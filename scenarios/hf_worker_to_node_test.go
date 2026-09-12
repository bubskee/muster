package scenarios_test

import (
	"testing"

	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/scenarios"
)

func TestHFWorkerToNodeReachesNodeAccessWithoutControls(t *testing.T) {
	scenario := scenarios.HFWorkerToNode()

	result := engine.Replay(scenario)

	if len(result.Trace) != 6 {
		t.Fatalf("trace length = %d, want 6", len(result.Trace))
	}

	for _, entry := range result.Trace {
		if entry.Status != engine.EventApplied {
			t.Fatalf(
				"event %q status = %q, want applied",
				entry.EventID,
				entry.Status,
			)
		}
	}

	foundNodeAccess := false

	for _, fact := range result.TerminalState {
		if fact == "access:node" {
			foundNodeAccess = true
			break
		}
	}

	if !foundNodeAccess {
		t.Fatalf(
			"terminal state does not contain access:node: %v",
			result.TerminalState,
		)
	}
}
