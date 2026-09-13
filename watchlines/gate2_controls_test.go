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

func TestGate2MemorylessAblationCollapsesEarlyBetweenDistinction(t *testing.T) {
	controls := gate2BroadCorrelatedMemoryless()

	early := engine.Replay(
		scenarios.Gate2EarlyFailure(),
		controls...,
	)

	between := engine.Replay(
		scenarios.Gate2BetweenFailure(),
		controls...,
	)

	earlyNode := gate2ContainsFact(
		early.TerminalState,
		scenarios.FactGate2NodeAccess,
	)

	betweenNode := gate2ContainsFact(
		between.TerminalState,
		scenarios.FactGate2NodeAccess,
	)

	// Removing the ability to consume persisted evidence should collapse
	// the original EARLY/BETWEEN outcome distinction.
	if !earlyNode || !betweenNode {
		t.Fatalf(
			"memoryless ablation should make both runs reach node: early=%v between=%v",
			early.TerminalState,
			between.TerminalState,
		)
	}

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
		"review-credential-history-memoryless",
	)

	betweenH1 := gate2ControlResult(
		t,
		betweenE3,
		"review-credential-history-memoryless",
	)

	// EARLY has no historical alert because telemetry was already dead
	// when E2 occurred.
	if earlyH1.Disposition != engine.DispositionWaiting || earlyH1.Acted {
		t.Fatalf(
			"EARLY memoryless H1 = %+v, want waiting and not acted",
			earlyH1,
		)
	}

	if earlyH1.Reason != "missing:"+Gate2AlertCredentialTheft {
		t.Fatalf(
			"EARLY memoryless H1 reason = %q, want missing historical alert",
			earlyH1.Reason,
		)
	}

	// BETWEEN still has the historical alert, but the ablation prevents
	// that alert from being consumed after its originating telemetry
	// substrate has disappeared.
	if betweenH1.Disposition != engine.DispositionWaiting || betweenH1.Acted {
		t.Fatalf(
			"BETWEEN memoryless H1 = %+v, want waiting and not acted",
			betweenH1,
		)
	}

	wantBetweenReason := "present:" + scenarios.FactGate2WorkerTelemetryCompromised
	if betweenH1.Reason != wantBetweenReason {
		t.Fatalf(
			"BETWEEN memoryless H1 reason = %q, want %q",
			betweenH1.Reason,
			wantBetweenReason,
		)
	}

	// Neither run should revoke the credential, so replay remains possible.
	for name, result := range map[string]engine.RunResult{
		"early":   early,
		"between": between,
	} {
		if !gate2ContainsFact(
			result.TerminalState,
			scenarios.FactGate2ClusterCredential,
		) {
			t.Fatalf(
				"%s unexpectedly revoked cluster credential: %v",
				name,
				result.TerminalState,
			)
		}

		e4 := gate2TraceEntry(
			t,
			result,
			scenarios.EventGate2CredentialReplay,
		)

		if e4.Status != engine.EventApplied {
			t.Fatalf(
				"%s E4 status = %q, want applied",
				name,
				e4.Status,
			)
		}
	}
}

func TestGate2PrematureContainmentRequiresBothL1AndR2(t *testing.T) {
	baseline := Gate2FullResponseWithoutPrematureContainment()

	tests := []struct {
		name     string
		controls []engine.Control
		wantNode bool
	}{
		{
			name:     "baseline",
			controls: baseline,
			wantNode: false,
		},
		{
			name: "plus-l1-only",
			controls: append(
				append([]engine.Control(nil), baseline...),
				Gate2L1WorkerContainmentEscalation(),
			),
			wantNode: false,
		},
		{
			name: "plus-r2-only",
			controls: append(
				append([]engine.Control(nil), baseline...),
				Gate2R2IsolateWorker(),
			),
			wantNode: false,
		},
		{
			name: "plus-l1-and-r2",
			controls: append(
				append([]engine.Control(nil), baseline...),
				Gate2L1WorkerContainmentEscalation(),
				Gate2R2IsolateWorker(),
			),
			wantNode: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Replay(
				scenarios.Gate2NoFailure(),
				tt.controls...,
			)

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

func TestGate2L1AloneEscalatesWithoutContainment(t *testing.T) {
	controls := append(
		Gate2FullResponseWithoutPrematureContainment(),
		Gate2L1WorkerContainmentEscalation(),
	)

	result := engine.Replay(
		scenarios.Gate2NoFailure(),
		controls...,
	)

	e2 := gate2TraceEntry(
		t,
		result,
		scenarios.EventGate2CredentialTheft,
	)

	l1 := gate2ControlResult(
		t,
		e2,
		"escalate-worker-containment",
	)

	if l1.Disposition != engine.DispositionReady || !l1.Acted {
		t.Fatalf(
			"L1 = %+v, want ready and acted",
			l1,
		)
	}

	if !gate2ContainsFact(
		e2.After,
		scenarios.FactGate2WorkerAccess,
	) {
		t.Fatalf(
			"L1 alone unexpectedly removed worker access: %v",
			e2.After,
		)
	}

	e3 := gate2TraceEntry(
		t,
		result,
		scenarios.EventGate2K8sDiscovery,
	)

	if e3.Status != engine.EventApplied {
		t.Fatalf(
			"E3 status = %q, want applied",
			e3.Status,
		)
	}

	if gate2ContainsFact(
		result.TerminalState,
		scenarios.FactGate2NodeAccess,
	) {
		t.Fatalf(
			"L1 alone worsened outcome: %v",
			result.TerminalState,
		)
	}
}

func TestGate2PrematureContainmentInvariantToControlDeclarationOrder(t *testing.T) {
	tests := []struct {
		name     string
		controls []engine.Control
		wantNode bool
	}{
		{
			name:     "baseline",
			controls: Gate2FullResponseWithoutPrematureContainment(),
			wantNode: false,
		},
		{
			name:     "jackpot",
			controls: Gate2FullResponseWithPrematureContainment(),
			wantNode: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forward := engine.Replay(
				scenarios.Gate2NoFailure(),
				tt.controls...,
			)

			reversed := engine.Replay(
				scenarios.Gate2NoFailure(),
				gate2ReverseControls(tt.controls)...,
			)

			forwardNode := gate2ContainsFact(
				forward.TerminalState,
				scenarios.FactGate2NodeAccess,
			)

			reversedNode := gate2ContainsFact(
				reversed.TerminalState,
				scenarios.FactGate2NodeAccess,
			)

			if forwardNode != tt.wantNode {
				t.Fatalf(
					"forward node-access = %v, want %v; terminal=%v",
					forwardNode,
					tt.wantNode,
					forward.TerminalState,
				)
			}

			if reversedNode != tt.wantNode {
				t.Fatalf(
					"reversed node-access = %v, want %v; terminal=%v",
					reversedNode,
					tt.wantNode,
					reversed.TerminalState,
				)
			}

			if !gate2SameFacts(
				forward.TerminalState,
				reversed.TerminalState,
			) {
				t.Fatalf(
					"terminal state depends on control declaration order:\nforward=%v\nreversed=%v",
					forward.TerminalState,
					reversed.TerminalState,
				)
			}
		})
	}
}

func gate2SameFacts(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	seen := make(map[string]int, len(a))

	for _, fact := range a {
		seen[fact]++
	}

	for _, fact := range b {
		seen[fact]--
		if seen[fact] < 0 {
			return false
		}
	}

	for _, count := range seen {
		if count != 0 {
			return false
		}
	}

	return true
}

func gate2CollapsedCredentialAlertResponse() engine.EventControl {
	return engine.EventControl{
		ID:       "collapsed-response-credential-alert",
		Role:     engine.Reserve,
		EventIDs: []string{scenarios.EventGate2K8sDiscovery},
		Requires: []engine.Condition{
			{
				Fact: Gate2AlertCredentialTheft,
				Op:   engine.ConditionPresent,
			},
			{
				Fact: scenarios.FactGate2ClusterCredential,
				Op:   engine.ConditionPresent,
			},
		},
		Action: engine.ActionRespond,
		Effects: []engine.Effect{
			{
				Fact: scenarios.FactGate2ClusterCredential,
				Op:   engine.EffectRemove,
			},
		},
	}
}

func gate2CollapsedK8sAlertResponse() engine.EventControl {
	return engine.EventControl{
		ID:       "collapsed-response-k8s-alert",
		Role:     engine.Reserve,
		EventIDs: []string{scenarios.EventGate2K8sDiscovery},
		Requires: []engine.Condition{
			{
				Fact: Gate2AlertK8sDiscovery,
				Op:   engine.ConditionPresent,
			},
			{
				Fact: scenarios.FactGate2ClusterCredential,
				Op:   engine.ConditionPresent,
			},
		},
		Action: engine.ActionRespond,
		Effects: []engine.Effect{
			{
				Fact: scenarios.FactGate2ClusterCredential,
				Op:   engine.EffectRemove,
			},
		},
	}
}

func gate2BroadCorrelatedCollapsedPipeline() []engine.Control {
	return []engine.Control{
		Gate2V1CredentialTheft(),
		Gate2V2WorkerK8sDiscovery(),

		gate2CollapsedCredentialAlertResponse(),
		gate2CollapsedK8sAlertResponse(),
	}
}

func gate2NarrowDiverseCollapsedPipeline() []engine.Control {
	return []engine.Control{
		Gate2V2WorkerK8sDiscovery(),
		Gate2V3K8sAuditDiscovery(),

		gate2CollapsedCredentialAlertResponse(),
		gate2CollapsedK8sAlertResponse(),
	}
}

func TestGate2PipelineCollapsePreservesObservationBudgetOutcomes(t *testing.T) {
	tests := []struct {
		name      string
		scenario  engine.Scenario
		full      []engine.Control
		collapsed []engine.Control
	}{
		{
			name:      "broad-none",
			scenario:  scenarios.Gate2NoFailure(),
			full:      Gate2BroadCorrelated(),
			collapsed: gate2BroadCorrelatedCollapsedPipeline(),
		},
		{
			name:      "broad-early",
			scenario:  scenarios.Gate2EarlyFailure(),
			full:      Gate2BroadCorrelated(),
			collapsed: gate2BroadCorrelatedCollapsedPipeline(),
		},
		{
			name:      "broad-between",
			scenario:  scenarios.Gate2BetweenFailure(),
			full:      Gate2BroadCorrelated(),
			collapsed: gate2BroadCorrelatedCollapsedPipeline(),
		},
		{
			name:      "diverse-none",
			scenario:  scenarios.Gate2NoFailure(),
			full:      Gate2NarrowDiverse(),
			collapsed: gate2NarrowDiverseCollapsedPipeline(),
		},
		{
			name:      "diverse-early",
			scenario:  scenarios.Gate2EarlyFailure(),
			full:      Gate2NarrowDiverse(),
			collapsed: gate2NarrowDiverseCollapsedPipeline(),
		},
		{
			name:      "diverse-between",
			scenario:  scenarios.Gate2BetweenFailure(),
			full:      Gate2NarrowDiverse(),
			collapsed: gate2NarrowDiverseCollapsedPipeline(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			full := engine.Replay(
				tt.scenario,
				tt.full...,
			)

			collapsed := engine.Replay(
				tt.scenario,
				tt.collapsed...,
			)

			fullNode := gate2ContainsFact(
				full.TerminalState,
				scenarios.FactGate2NodeAccess,
			)

			collapsedNode := gate2ContainsFact(
				collapsed.TerminalState,
				scenarios.FactGate2NodeAccess,
			)

			if fullNode != collapsedNode {
				t.Fatalf(
					"pipeline collapse changed terminal outcome: full=%v collapsed=%v\nfull-state=%v\ncollapsed-state=%v",
					fullNode,
					collapsedNode,
					full.TerminalState,
					collapsed.TerminalState,
				)
			}
		})
	}
}

func gate2CollapsedFullResponseWithoutPrematureContainment() []engine.Control {
	return []engine.Control{
		Gate2V1CredentialTheft(),
		Gate2V2WorkerK8sDiscovery(),
		gate2CollapsedK8sAlertResponse(),
	}
}

func gate2CollapsedFullResponseWithPrematureContainment() []engine.Control {
	controls := gate2CollapsedFullResponseWithoutPrematureContainment()

	return append(
		controls,
		Gate2L1WorkerContainmentEscalation(),
		Gate2R2IsolateWorker(),
	)
}

func TestGate2PipelineCollapsePreservesPrematureContainmentNonMonotonicity(t *testing.T) {
	fullWithout := engine.Replay(
		scenarios.Gate2NoFailure(),
		Gate2FullResponseWithoutPrematureContainment()...,
	)

	fullWith := engine.Replay(
		scenarios.Gate2NoFailure(),
		Gate2FullResponseWithPrematureContainment()...,
	)

	collapsedWithout := engine.Replay(
		scenarios.Gate2NoFailure(),
		gate2CollapsedFullResponseWithoutPrematureContainment()...,
	)

	collapsedWith := engine.Replay(
		scenarios.Gate2NoFailure(),
		gate2CollapsedFullResponseWithPrematureContainment()...,
	)

	fullWithoutNode := gate2ContainsFact(
		fullWithout.TerminalState,
		scenarios.FactGate2NodeAccess,
	)

	fullWithNode := gate2ContainsFact(
		fullWith.TerminalState,
		scenarios.FactGate2NodeAccess,
	)

	collapsedWithoutNode := gate2ContainsFact(
		collapsedWithout.TerminalState,
		scenarios.FactGate2NodeAccess,
	)

	collapsedWithNode := gate2ContainsFact(
		collapsedWith.TerminalState,
		scenarios.FactGate2NodeAccess,
	)

	if fullWithoutNode || collapsedWithoutNode {
		t.Fatalf(
			"baseline should remain secure: full=%v collapsed=%v",
			fullWithout.TerminalState,
			collapsedWithout.TerminalState,
		)
	}

	if !fullWithNode || !collapsedWithNode {
		t.Fatalf(
			"premature-containment arm should remain adverse: full=%v collapsed=%v",
			fullWith.TerminalState,
			collapsedWith.TerminalState,
		)
	}
}

func TestGate2BReviewBudgetMatrix(t *testing.T) {
	tests := []struct {
		name     string
		scenario engine.Scenario
		controls []engine.Control
		wantNode bool
	}{
		{
			name:     "correlated-none",
			scenario: scenarios.Gate2NoFailure(),
			controls: Gate2CorrelatedReview(),
			wantNode: false,
		},
		{
			name:     "correlated-soc-failure",
			scenario: scenarios.Gate2SOCReviewFailure(),
			controls: Gate2CorrelatedReview(),
			wantNode: true,
		},
		{
			name:     "correlated-agent-failure",
			scenario: scenarios.Gate2AgentReviewFailure(),
			controls: Gate2CorrelatedReview(),
			wantNode: false,
		},
		{
			name:     "diverse-none",
			scenario: scenarios.Gate2NoFailure(),
			controls: Gate2DiverseReview(),
			wantNode: false,
		},
		{
			name:     "diverse-soc-failure",
			scenario: scenarios.Gate2SOCReviewFailure(),
			controls: Gate2DiverseReview(),
			wantNode: false,
		},
		{
			name:     "diverse-agent-failure",
			scenario: scenarios.Gate2AgentReviewFailure(),
			controls: Gate2DiverseReview(),
			wantNode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Replay(
				tt.scenario,
				tt.controls...,
			)

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

func TestGate2BSOCFailureRewardsReviewIndependence(t *testing.T) {
	correlated := engine.Replay(
		scenarios.Gate2SOCReviewFailure(),
		Gate2CorrelatedReview()...,
	)

	diverse := engine.Replay(
		scenarios.Gate2SOCReviewFailure(),
		Gate2DiverseReview()...,
	)

	correlatedE3 := gate2TraceEntry(
		t,
		correlated,
		scenarios.EventGate2K8sDiscovery,
	)

	h1 := gate2ControlResult(
		t,
		correlatedE3,
		"review-credential-history",
	)

	h2 := gate2ControlResult(
		t,
		correlatedE3,
		"review-k8s-discovery",
	)

	if h1.Disposition != engine.DispositionSuppressed || h1.Acted {
		t.Fatalf(
			"correlated H1 = %+v, want suppressed and not acted",
			h1,
		)
	}

	if h2.Disposition != engine.DispositionSuppressed || h2.Acted {
		t.Fatalf(
			"correlated H2 = %+v, want suppressed and not acted",
			h2,
		)
	}

	if h1.Reason != scenarios.FactGate2SOCQueueOverloaded {
		t.Fatalf(
			"H1 suppression reason = %q, want %q",
			h1.Reason,
			scenarios.FactGate2SOCQueueOverloaded,
		)
	}

	if h2.Reason != scenarios.FactGate2SOCQueueOverloaded {
		t.Fatalf(
			"H2 suppression reason = %q, want %q",
			h2.Reason,
			scenarios.FactGate2SOCQueueOverloaded,
		)
	}

	if !gate2ContainsFact(
		correlated.TerminalState,
		scenarios.FactGate2NodeAccess,
	) {
		t.Fatalf(
			"correlated SOC failure should reach node: %v",
			correlated.TerminalState,
		)
	}

	diverseE3 := gate2TraceEntry(
		t,
		diverse,
		scenarios.EventGate2K8sDiscovery,
	)

	diverseH1 := gate2ControlResult(
		t,
		diverseE3,
		"review-credential-history",
	)

	a1 := gate2ControlResult(
		t,
		diverseE3,
		"review-k8s-discovery-agent",
	)

	if diverseH1.Disposition != engine.DispositionSuppressed || diverseH1.Acted {
		t.Fatalf(
			"diverse H1 = %+v, want suppressed and not acted",
			diverseH1,
		)
	}

	if a1.Disposition != engine.DispositionReady || !a1.Acted {
		t.Fatalf(
			"diverse A1 = %+v, want ready and acted",
			a1,
		)
	}

	r1 := gate2ControlResult(
		t,
		diverseE3,
		"revoke-cluster-credential",
	)

	if r1.Disposition != engine.DispositionReady || !r1.Acted {
		t.Fatalf(
			"diverse R1 = %+v, want ready and acted",
			r1,
		)
	}

	if gate2ContainsFact(
		diverse.TerminalState,
		scenarios.FactGate2NodeAccess,
	) {
		t.Fatalf(
			"diverse SOC failure reached node: %v",
			diverse.TerminalState,
		)
	}
}

func TestGate2BAgentFailureDoesNotMakeDiverseReviewMagicallyBetter(t *testing.T) {
	result := engine.Replay(
		scenarios.Gate2AgentReviewFailure(),
		Gate2DiverseReview()...,
	)

	e3 := gate2TraceEntry(
		t,
		result,
		scenarios.EventGate2K8sDiscovery,
	)

	a1 := gate2ControlResult(
		t,
		e3,
		"review-k8s-discovery-agent",
	)

	h1 := gate2ControlResult(
		t,
		e3,
		"review-credential-history",
	)

	if a1.Disposition != engine.DispositionSuppressed || a1.Acted {
		t.Fatalf(
			"A1 = %+v, want suppressed and not acted",
			a1,
		)
	}

	if a1.Reason != scenarios.FactGate2AgentReviewerUnavailable {
		t.Fatalf(
			"A1 suppression reason = %q, want %q",
			a1.Reason,
			scenarios.FactGate2AgentReviewerUnavailable,
		)
	}

	if h1.Disposition != engine.DispositionReady || !h1.Acted {
		t.Fatalf(
			"H1 = %+v, want ready and acted",
			h1,
		)
	}

	if gate2ContainsFact(
		result.TerminalState,
		scenarios.FactGate2NodeAccess,
	) {
		t.Fatalf(
			"surviving H1 should preserve secure outcome: %v",
			result.TerminalState,
		)
	}
}

func TestGate2CIntegratedCompositionMatrix(t *testing.T) {
	type architecture struct {
		name     string
		controls func() []engine.Control
	}

	type treatment struct {
		name     string
		scenario func() engine.Scenario
	}

	architectures := []architecture{
		{"cc", Gate2CC},
		{"cd", Gate2CD},
		{"dc", Gate2DC},
		{"dd", Gate2DD},
	}

	treatments := []treatment{
		{"none", scenarios.Gate2BetweenFailure},
		{"soc-failure", scenarios.Gate2BetweenSOCReviewFailure},
		{"agent-failure", scenarios.Gate2BetweenAgentReviewFailure},
	}

	wantNode := map[string]map[string]bool{
		"cc": {
			"none":          false,
			"soc-failure":   true,
			"agent-failure": false,
		},
		"cd": {
			"none":          false,
			"soc-failure":   true,
			"agent-failure": false,
		},
		"dc": {
			"none":          false,
			"soc-failure":   true,
			"agent-failure": false,
		},
		"dd": {
			"none":          false,
			"soc-failure":   false,
			"agent-failure": true,
		},
	}

	for _, arch := range architectures {
		for _, treatment := range treatments {
			t.Run(arch.name+"-"+treatment.name, func(t *testing.T) {
				result := engine.Replay(
					treatment.scenario(),
					arch.controls()...,
				)

				gotNode := gate2ContainsFact(
					result.TerminalState,
					scenarios.FactGate2NodeAccess,
				)

				want := wantNode[arch.name][treatment.name]

				if gotNode != want {
					t.Fatalf(
						"node-access = %v, want %v; terminal=%v",
						gotNode,
						want,
						result.TerminalState,
					)
				}
			})
		}
	}
}
