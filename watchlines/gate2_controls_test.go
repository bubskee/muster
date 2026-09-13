package watchlines

import (
	"testing"

	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/scenarios"
)

func TestGate2AObservationBudgetMatrix(t *testing.T) {
	tests := []struct {
		name     string
		scenario engine.Scenario
		controls []engine.Control
		wantNode bool
	}{
		{"broad-none", scenarios.Gate2NoFailure(), Gate2BroadCorrelated(), false},
		{"broad-early", scenarios.Gate2EarlyFailure(), Gate2BroadCorrelated(), true},
		{"broad-between", scenarios.Gate2BetweenFailure(), Gate2BroadCorrelated(), false},
		{"diverse-none", scenarios.Gate2NoFailure(), Gate2NarrowDiverse(), false},
		{"diverse-early", scenarios.Gate2EarlyFailure(), Gate2NarrowDiverse(), false},
		{"diverse-between", scenarios.Gate2BetweenFailure(), Gate2NarrowDiverse(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Replay(tt.scenario, tt.controls...)

			gotNode := gate2ContainsFact(
				result.TerminalState,
				scenarios.FactGate2NodeAccess,
			)

			if gotNode != tt.wantNode {
				t.Fatalf(
					"terminal state = %v, node-access = %v, want %v",
					result.TerminalState,
					gotNode,
					tt.wantNode,
				)
			}
		})
	}
}

func TestGate2ABetweenCarriesPersistedAlertAcrossTelemetryFailure(t *testing.T) {
	early := engine.Replay(
		scenarios.Gate2EarlyFailure(),
		Gate2BroadCorrelated()...,
	)

	between := engine.Replay(
		scenarios.Gate2BetweenFailure(),
		Gate2BroadCorrelated()...,
	)

	earlyE3 := gate2TraceEntry(
		t,
		early,
		scenarios.EventGate2K8sDiscovery,
	)

	betweenE3 := gate2TraceEntry(
		t,
		between,
		scenarios.EventGate2K8sDiscovery,
	)

	earlyH1 := gate2ControlResult(
		t,
		earlyE3,
		"review-credential-history",
	)

	betweenH1 := gate2ControlResult(
		t,
		betweenE3,
		"review-credential-history",
	)

	if earlyH1.Disposition != engine.DispositionWaiting || earlyH1.Acted {
		t.Fatalf(
			"EARLY H1 = %+v, want waiting and not acted",
			earlyH1,
		)
	}

	if betweenH1.Disposition != engine.DispositionReady || !betweenH1.Acted {
		t.Fatalf(
			"BETWEEN H1 = %+v, want ready and acted",
			betweenH1,
		)
	}

	earlyNode := gate2ContainsFact(
		early.TerminalState,
		scenarios.FactGate2NodeAccess,
	)

	betweenNode := gate2ContainsFact(
		between.TerminalState,
		scenarios.FactGate2NodeAccess,
	)

	if earlyNode == betweenNode {
		t.Fatalf(
			"EARLY and BETWEEN should diverge: early=%v between=%v",
			early.TerminalState,
			between.TerminalState,
		)
	}
}

func TestGate2PrematureContainmentCanWorsenOutcome(t *testing.T) {
	without := engine.Replay(
		scenarios.Gate2NoFailure(),
		Gate2FullResponseWithoutPrematureContainment()...,
	)

	with := engine.Replay(
		scenarios.Gate2NoFailure(),
		Gate2FullResponseWithPrematureContainment()...,
	)

	// First establish the terminal non-monotonicity:
	//
	// adding the functioning L1/R2 containment path makes
	// the security outcome strictly worse.
	if gate2ContainsFact(
		without.TerminalState,
		scenarios.FactGate2NodeAccess,
	) {
		t.Fatalf(
			"without premature containment reached node: %v",
			without.TerminalState,
		)
	}

	if !gate2ContainsFact(
		with.TerminalState,
		scenarios.FactGate2NodeAccess,
	) {
		t.Fatalf(
			"with premature containment did not reach node: %v",
			with.TerminalState,
		)
	}

	// Now verify the mechanism rather than accepting the
	// terminal-state difference alone.
	//
	// Without premature containment, E3 occurs and generates
	// the later evidence required for the full response path.
	withoutE3 := gate2TraceEntry(
		t,
		without,
		scenarios.EventGate2K8sDiscovery,
	)

	if withoutE3.Status != engine.EventApplied {
		t.Fatalf(
			"without containment E3 status = %q, want applied",
			withoutE3.Status,
		)
	}

	v2 := gate2ControlResult(
		t,
		withoutE3,
		"worker-k8s-discovery",
	)

	if v2.Disposition != engine.DispositionReady || !v2.Acted {
		t.Fatalf(
			"without containment V2 = %+v, want ready and acted",
			v2,
		)
	}

	h2 := gate2ControlResult(
		t,
		withoutE3,
		"review-k8s-discovery",
	)

	if h2.Disposition != engine.DispositionReady || !h2.Acted {
		t.Fatalf(
			"without containment H2 = %+v, want ready and acted",
			h2,
		)
	}

	r1 := gate2ControlResult(
		t,
		withoutE3,
		"revoke-cluster-credential",
	)

	if r1.Disposition != engine.DispositionReady || !r1.Acted {
		t.Fatalf(
			"without containment R1 = %+v, want ready and acted",
			r1,
		)
	}

	// The full response should revoke the portable cluster
	// credential before E4 can replay it.
	if gate2ContainsFact(
		without.TerminalState,
		scenarios.FactGate2ClusterCredential,
	) {
		t.Fatalf(
			"without containment retained cluster credential: %v",
			without.TerminalState,
		)
	}

	// With the extra functioning containment path, worker access
	// is removed after E2. That destroys the later E3 evidence
	// opportunity but leaves the already-stolen credential intact.
	withE3 := gate2TraceEntry(
		t,
		with,
		scenarios.EventGate2K8sDiscovery,
	)

	if withE3.Status != engine.EventSkipped {
		t.Fatalf(
			"with containment E3 status = %q, want skipped",
			withE3.Status,
		)
	}

	withE2 := gate2TraceEntry(
		t,
		with,
		scenarios.EventGate2CredentialTheft,
	)

	l1 := gate2ControlResult(
		t,
		withE2,
		"escalate-worker-containment",
	)

	if l1.Disposition != engine.DispositionReady || !l1.Acted {
		t.Fatalf(
			"with containment L1 = %+v, want ready and acted",
			l1,
		)
	}

	r2 := gate2ControlResult(
		t,
		withE2,
		"isolate-worker",
	)

	if r2.Disposition != engine.DispositionReady || !r2.Acted {
		t.Fatalf(
			"with containment R2 = %+v, want ready and acted",
			r2,
		)
	}

	if gate2ContainsFact(
		withE2.After,
		scenarios.FactGate2WorkerAccess,
	) {
		t.Fatalf(
			"R2 acted but worker access survived E2: %v",
			withE2.After,
		)
	}

	if !gate2ContainsFact(
		with.TerminalState,
		scenarios.FactGate2ClusterCredential,
	) {
		t.Fatalf(
			"with containment unexpectedly removed cluster credential: %v",
			with.TerminalState,
		)
	}

	// And the credential replay should consequently succeed.
	withE4 := gate2TraceEntry(
		t,
		with,
		scenarios.EventGate2CredentialReplay,
	)

	if withE4.Status != engine.EventApplied {
		t.Fatalf(
			"with containment E4 status = %q, want applied",
			withE4.Status,
		)
	}
}

func gate2TraceEntry(
	t *testing.T,
	result engine.RunResult,
	eventID string,
) engine.TraceEntry {
	t.Helper()

	for _, entry := range result.Trace {
		if entry.EventID == eventID {
			return entry
		}
	}

	t.Fatalf("event %q not found in trace", eventID)
	return engine.TraceEntry{}
}

func gate2ControlResult(
	t *testing.T,
	entry engine.TraceEntry,
	controlID string,
) engine.ControlResult {
	t.Helper()

	for _, result := range entry.Controls {
		if result.ControlID == controlID {
			return result
		}
	}

	t.Fatalf(
		"control %q not found for event %q",
		controlID,
		entry.EventID,
	)

	return engine.ControlResult{}
}

func gate2ContainsFact(facts []string, want string) bool {
	for _, fact := range facts {
		if fact == want {
			return true
		}
	}

	return false
}
