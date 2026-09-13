package engine_test

import (
	"testing"

	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/scenarios"
	"github.com/bubskee/muster/watchlines"
)

func TestSummaryDistinguishesDefensiveCompositions(t *testing.T) {
	scenario := scenarios.HFWorkerToNode()

	noControls := engine.Summarize(
		engine.Replay(scenario),
	)

	vedettes := engine.Summarize(
		engine.Replay(
			scenario,
			watchlines.HFVedettesOnly()...,
		),
	)

	interleaved := engine.Summarize(
		engine.Replay(
			scenario,
			watchlines.HFInterleaved()...,
		),
	)

	// No controls: attack succeeds, but it was never detected.
	if noControls.EventsAfterDetection != nil {
		t.Fatalf(
			"no-controls events after detection = %v, want nil",
			*noControls.EventsAfterDetection,
		)
	}

	if noControls.FurthestApplied != scenarios.EventNodeAccess {
		t.Fatalf(
			"no-controls furthest applied = %q, want %q",
			noControls.FurthestApplied,
			scenarios.EventNodeAccess,
		)
	}

	// Vedettes-only: detected at E2, then four more attacker events succeed.
	if vedettes.EventsAfterDetection == nil {
		t.Fatal("vedettes events after detection = nil, want 4")
	}

	if *vedettes.EventsAfterDetection != 4 {
		t.Fatalf(
			"vedettes events after detection = %d, want 4",
			*vedettes.EventsAfterDetection,
		)
	}

	if vedettes.FurthestApplied != scenarios.EventNodeAccess {
		t.Fatalf(
			"vedettes furthest applied = %q, want %q",
			vedettes.FurthestApplied,
			scenarios.EventNodeAccess,
		)
	}

	// Interleaved: E3 is attempted and blocked; nothing after detection succeeds.
	if interleaved.EventsAfterDetection == nil {
		t.Fatal("interleaved events after detection = nil, want 0")
	}

	if *interleaved.EventsAfterDetection != 0 {
		t.Fatalf(
			"interleaved events after detection = %d, want 0",
			*interleaved.EventsAfterDetection,
		)
	}

	if interleaved.FurthestApplied != scenarios.EventWorkerCodeExecution {
		t.Fatalf(
			"interleaved furthest applied = %q, want %q",
			interleaved.FurthestApplied,
			scenarios.EventWorkerCodeExecution,
		)
	}

	if interleaved.FurthestAttempted != scenarios.EventServiceAccountAccess {
		t.Fatalf(
			"interleaved furthest attempted = %q, want %q",
			interleaved.FurthestAttempted,
			scenarios.EventServiceAccountAccess,
		)
	}

	if interleaved.FirstBlock != scenarios.EventServiceAccountAccess {
		t.Fatalf(
			"interleaved first block = %q, want %q",
			interleaved.FirstBlock,
			scenarios.EventServiceAccountAccess,
		)
	}
}
