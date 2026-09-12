package watchlines_test

import (
	"testing"

	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/scenarios"
	"github.com/bubskee/muster/watchlines"
)

func TestHFInterleavedChangesIncidentTrajectory(t *testing.T) {
	scenario := scenarios.HFWorkerToNode()

	result := engine.Replay(
		scenario,
		watchlines.HFInterleaved()...,
	)

	if len(result.Trace) != 6 {
		t.Fatalf("trace length = %d, want 6", len(result.Trace))
	}

	if result.Trace[0].Status != engine.EventApplied {
		t.Fatalf(
			"%s status = %q, want applied",
			result.Trace[0].EventID,
			result.Trace[0].Status,
		)
	}

	if result.Trace[1].Status != engine.EventApplied {
		t.Fatalf(
			"%s status = %q, want applied",
			result.Trace[1].EventID,
			result.Trace[1].Status,
		)
	}

	credentialAccess := result.Trace[2]

	if credentialAccess.EventID != scenarios.EventServiceAccountAccess {
		t.Fatalf(
			"trace[2] = %q, want %q",
			credentialAccess.EventID,
			scenarios.EventServiceAccountAccess,
		)
	}

	if credentialAccess.Status != engine.EventBlocked {
		t.Fatalf(
			"credential access status = %q, want blocked",
			credentialAccess.Status,
		)
	}

	if len(credentialAccess.BlockedBy) != 1 {
		t.Fatalf(
			"blocked by = %v, want one control",
			credentialAccess.BlockedBy,
		)
	}

	if credentialAccess.BlockedBy[0] != "service-account-credential-guard" {
		t.Fatalf(
			"blocked by = %q, want service-account-credential-guard",
			credentialAccess.BlockedBy[0],
		)
	}

	for i := 3; i < len(result.Trace); i++ {
		if result.Trace[i].Status != engine.EventSkipped {
			t.Fatalf(
				"%s status = %q, want skipped",
				result.Trace[i].EventID,
				result.Trace[i].Status,
			)
		}
	}

	for _, fact := range result.TerminalState {
		if fact == "access:node" {
			t.Fatalf(
				"terminal state unexpectedly contains access:node: %v",
				result.TerminalState,
			)
		}
	}
}

func TestDefensiveCompositionChangesOutcome(t *testing.T) {
	scenario := scenarios.HFWorkerToNode()

	vedettesOnly := engine.Replay(
		scenario,
		watchlines.HFVedettesOnly()...,
	)

	interleaved := engine.Replay(
		scenario,
		watchlines.HFInterleaved()...,
	)

	if vedettesOnly.Trace[5].Status != engine.EventApplied {
		t.Fatalf(
			"vedettes-only node access = %q, want applied",
			vedettesOnly.Trace[5].Status,
		)
	}

	if interleaved.Trace[5].Status != engine.EventSkipped {
		t.Fatalf(
			"interleaved node access = %q, want skipped",
			interleaved.Trace[5].Status,
		)
	}
}
