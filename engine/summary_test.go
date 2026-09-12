package engine_test

import (
	"testing"

	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/scenarios"
	"github.com/bubskee/muster/watchlines"
)

func TestSummaryDistinguishesDefensiveCompositions(t *testing.T) {
	scenario := scenarios.HFWorkerToNode()

	vedettes := engine.Replay(
		scenario,
		watchlines.HFVedettesOnly()...,
	)

	interleaved := engine.Replay(
		scenario,
		watchlines.HFInterleaved()...,
	)

	vedetteSummary := engine.Summarize(vedettes)
	interleavedSummary := engine.Summarize(interleaved)

	if vedetteSummary.FirstObservation != scenarios.EventWorkerCodeExecution {
		t.Fatalf(
			"vedettes first observation = %q",
			vedetteSummary.FirstObservation,
		)
	}

	if vedetteSummary.FirstBlock != "" {
		t.Fatalf(
			"vedettes first block = %q, want none",
			vedetteSummary.FirstBlock,
		)
	}

	if vedetteSummary.FurthestEvent != scenarios.EventNodeAccess {
		t.Fatalf(
			"vedettes furthest event = %q, want %q",
			vedetteSummary.FurthestEvent,
			scenarios.EventNodeAccess,
		)
	}

	if interleavedSummary.FirstBlock != scenarios.EventServiceAccountAccess {
		t.Fatalf(
			"interleaved first block = %q, want %q",
			interleavedSummary.FirstBlock,
			scenarios.EventServiceAccountAccess,
		)
	}

	if interleavedSummary.FurthestEvent != scenarios.EventServiceAccountAccess {
		t.Fatalf(
			"interleaved furthest event = %q, want %q",
			interleavedSummary.FurthestEvent,
			scenarios.EventServiceAccountAccess,
		)
	}
}
