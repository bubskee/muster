package watchlines_test

import (
	"testing"

	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/scenarios"
	"github.com/bubskee/muster/watchlines"
)

func TestHFVedettesObserveButDoNotStopIncident(t *testing.T) {
	scenario := scenarios.HFWorkerToNode()
	controls := watchlines.HFVedettesOnly()

	result := engine.Replay(scenario, controls...)

	var firstObserved int = -1
	var blocked int

	for i, entry := range result.Trace {
		if firstObserved == -1 && len(entry.ObservedBy) > 0 {
			firstObserved = i
		}

		blocked += len(entry.BlockedBy)
	}

	if firstObserved == -1 {
		t.Fatal("incident was never observed")
	}

	if result.Trace[firstObserved].EventID != scenarios.EventWorkerCodeExecution {
		t.Fatalf(
			"first observed event = %q, want %q",
			result.Trace[firstObserved].EventID,
			scenarios.EventWorkerCodeExecution,
		)
	}

	if blocked != 0 {
		t.Fatalf("blocked controls = %d, want 0", blocked)
	}

	eventsAfterDetection := len(result.Trace) - firstObserved - 1

	if eventsAfterDetection == 0 {
		t.Fatal("expected incident progression after first detection")
	}

	last := result.Trace[len(result.Trace)-1]

	if last.EventID != scenarios.EventNodeAccess {
		t.Fatalf(
			"final event = %q, want %q",
			last.EventID,
			scenarios.EventNodeAccess,
		)
	}

	if last.Status != engine.EventApplied {
		t.Fatalf(
			"node access status = %q, want applied",
			last.Status,
		)
	}
}
