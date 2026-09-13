package watchlines

import (
	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/scenarios"
)

const (
	Gate2AlertCredentialTheft       = "alert:credential-theft"
	Gate2AlertK8sDiscovery          = "alert:k8s-discovery"
	Gate2ReviewedClusterIntegrity   = "reviewed:cluster-integrity"
	Gate2CriticalClusterIntegrity   = "critical:cluster-integrity"
	Gate2EscalatedClusterIntegrity  = "escalated:cluster-integrity"
	Gate2EscalatedWorkerContainment = "escalated:worker-containment"
	Gate2OverloadedSOCQueue         = scenarios.FactGate2SOCQueueOverloaded
	Gate2AgentReviewerUnavailable   = scenarios.FactGate2AgentReviewerUnavailable
)

func Gate2V1CredentialTheft() engine.EventControl {
	return engine.EventControl{
		ID: "worker-credential-theft", Role: engine.Vedette,
		EventIDs:     []string{scenarios.EventGate2CredentialTheft},
		Substrate:    "worker-telemetry",
		SuppressedBy: []engine.Condition{{Fact: scenarios.FactGate2WorkerTelemetryCompromised, Op: engine.ConditionPresent}},
		Action:       engine.ActionObserve,
		Effects:      []engine.Effect{{Fact: Gate2AlertCredentialTheft, Op: engine.EffectAdd}},
	}
}

func Gate2V2WorkerK8sDiscovery() engine.EventControl {
	return engine.EventControl{
		ID: "worker-k8s-discovery", Role: engine.Vedette,
		EventIDs:     []string{scenarios.EventGate2K8sDiscovery},
		Substrate:    "worker-telemetry",
		SuppressedBy: []engine.Condition{{Fact: scenarios.FactGate2WorkerTelemetryCompromised, Op: engine.ConditionPresent}},
		Action:       engine.ActionObserve,
		Effects:      []engine.Effect{{Fact: Gate2AlertK8sDiscovery, Op: engine.EffectAdd}},
	}
}

func Gate2V3K8sAuditDiscovery() engine.EventControl {
	return engine.EventControl{
		ID: "k8s-audit-discovery", Role: engine.Vedette,
		EventIDs:  []string{scenarios.EventGate2K8sDiscovery},
		Substrate: "k8s-audit",
		Action:    engine.ActionObserve,
		Effects:   []engine.Effect{{Fact: Gate2AlertK8sDiscovery, Op: engine.EffectAdd}},
	}
}

func Gate2H1CredentialHistoryReview() engine.EventControl {
	return engine.EventControl{
		ID:           "review-credential-history",
		EventIDs:     []string{scenarios.EventGate2K8sDiscovery},
		Substrate:    "soc-queue",
		Requires:     []engine.Condition{{Fact: Gate2AlertCredentialTheft, Op: engine.ConditionPresent}},
		SuppressedBy: []engine.Condition{{Fact: Gate2OverloadedSOCQueue, Op: engine.ConditionPresent}},
		Action:       engine.ActionReview,
		Effects:      []engine.Effect{{Fact: Gate2ReviewedClusterIntegrity, Op: engine.EffectAdd}},
	}
}

func Gate2H2K8sDiscoveryReview() engine.EventControl {
	return engine.EventControl{
		ID:           "review-k8s-discovery",
		EventIDs:     []string{scenarios.EventGate2K8sDiscovery},
		Substrate:    "soc-queue",
		Requires:     []engine.Condition{{Fact: Gate2AlertK8sDiscovery, Op: engine.ConditionPresent}},
		SuppressedBy: []engine.Condition{{Fact: Gate2OverloadedSOCQueue, Op: engine.ConditionPresent}},
		Action:       engine.ActionReview,
		Effects:      []engine.Effect{{Fact: Gate2ReviewedClusterIntegrity, Op: engine.EffectAdd}},
	}
}

func Gate2C1ClusterIntegrityCriticality() engine.EventControl {
	return engine.EventControl{
		ID:       "assess-cluster-integrity",
		EventIDs: []string{scenarios.EventGate2K8sDiscovery},
		Requires: []engine.Condition{
			{Fact: Gate2ReviewedClusterIntegrity, Op: engine.ConditionPresent},
			{Fact: scenarios.FactGate2ClusterCredential, Op: engine.ConditionPresent},
		},
		Action:  engine.ActionCritical,
		Effects: []engine.Effect{{Fact: Gate2CriticalClusterIntegrity, Op: engine.EffectAdd}},
	}
}

func Gate2E1ClusterIntegrityEscalation() engine.EventControl {
	return engine.EventControl{
		ID:       "escalate-cluster-integrity",
		EventIDs: []string{scenarios.EventGate2K8sDiscovery},
		Requires: []engine.Condition{{Fact: Gate2CriticalClusterIntegrity, Op: engine.ConditionPresent}},
		Action:   engine.ActionEscalate,
		Effects:  []engine.Effect{{Fact: Gate2EscalatedClusterIntegrity, Op: engine.EffectAdd}},
	}
}

func Gate2R1RevokeClusterCredential() engine.EventControl {
	return engine.EventControl{
		ID: "revoke-cluster-credential", Role: engine.Reserve,
		EventIDs: []string{scenarios.EventGate2K8sDiscovery},
		Requires: []engine.Condition{{Fact: Gate2EscalatedClusterIntegrity, Op: engine.ConditionPresent}},
		Action:   engine.ActionRespond,
		Effects:  []engine.Effect{{Fact: scenarios.FactGate2ClusterCredential, Op: engine.EffectRemove}},
	}
}

func Gate2L1WorkerContainmentEscalation() engine.EventControl {
	return engine.EventControl{
		ID:       "escalate-worker-containment",
		EventIDs: []string{scenarios.EventGate2CredentialTheft},
		Requires: []engine.Condition{{Fact: Gate2AlertCredentialTheft, Op: engine.ConditionPresent}},
		Action:   engine.ActionEscalate,
		Effects:  []engine.Effect{{Fact: Gate2EscalatedWorkerContainment, Op: engine.EffectAdd}},
	}
}

func Gate2R2IsolateWorker() engine.EventControl {
	return engine.EventControl{
		ID: "isolate-worker", Role: engine.Reserve,
		EventIDs: []string{scenarios.EventGate2CredentialTheft},
		Requires: []engine.Condition{{Fact: Gate2EscalatedWorkerContainment, Op: engine.ConditionPresent}},
		Action:   engine.ActionRespond,
		Effects:  []engine.Effect{{Fact: scenarios.FactGate2WorkerAccess, Op: engine.EffectRemove}},
	}
}

func Gate2BroadCorrelated() []engine.Control {
	return []engine.Control{
		Gate2V1CredentialTheft(), Gate2V2WorkerK8sDiscovery(),
		Gate2H1CredentialHistoryReview(), Gate2H2K8sDiscoveryReview(),
		Gate2C1ClusterIntegrityCriticality(), Gate2E1ClusterIntegrityEscalation(),
		Gate2R1RevokeClusterCredential(),
	}
}

func Gate2NarrowDiverse() []engine.Control {
	return []engine.Control{
		Gate2V2WorkerK8sDiscovery(), Gate2V3K8sAuditDiscovery(),
		Gate2H1CredentialHistoryReview(), Gate2H2K8sDiscoveryReview(),
		Gate2C1ClusterIntegrityCriticality(), Gate2E1ClusterIntegrityEscalation(),
		Gate2R1RevokeClusterCredential(),
	}
}

func Gate2FullResponseWithoutPrematureContainment() []engine.Control {
	return []engine.Control{
		Gate2V1CredentialTheft(),

		// Full response depends on later E3 evidence.
		Gate2V2WorkerK8sDiscovery(),
		Gate2H2K8sDiscoveryReview(),

		Gate2C1ClusterIntegrityCriticality(),
		Gate2E1ClusterIntegrityEscalation(),
		Gate2R1RevokeClusterCredential(),
	}
}

func Gate2FullResponseWithPrematureContainment() []engine.Control {
	controls := Gate2FullResponseWithoutPrematureContainment()

	return append(
		controls,
		Gate2L1WorkerContainmentEscalation(),
		Gate2R2IsolateWorker(),
	)
}

func gate2MemorylessCredentialHistoryReview() engine.EventControl {
	return engine.EventControl{
		ID:       "review-credential-history-memoryless",
		EventIDs: []string{scenarios.EventGate2K8sDiscovery},
		Requires: []engine.Condition{
			{
				Fact: Gate2AlertCredentialTheft,
				Op:   engine.ConditionPresent,
			},
			{
				// Ablation: historical evidence is usable only while its
				// originating telemetry substrate is still healthy.
				Fact: scenarios.FactGate2WorkerTelemetryCompromised,
				Op:   engine.ConditionAbsent,
			},
		},
		Substrate: "soc-queue",
		SuppressedBy: []engine.Condition{
			{
				Fact: Gate2OverloadedSOCQueue,
				Op:   engine.ConditionPresent,
			},
		},
		Action: engine.ActionReview,
		Effects: []engine.Effect{
			{
				Fact: Gate2ReviewedClusterIntegrity,
				Op:   engine.EffectAdd,
			},
		},
	}
}

func gate2BroadCorrelatedMemoryless() []engine.Control {
	return []engine.Control{
		Gate2V1CredentialTheft(),
		Gate2V2WorkerK8sDiscovery(),

		gate2MemorylessCredentialHistoryReview(),
		Gate2H2K8sDiscoveryReview(),

		Gate2C1ClusterIntegrityCriticality(),
		Gate2E1ClusterIntegrityEscalation(),
		Gate2R1RevokeClusterCredential(),
	}
}

func gate2ReverseControls(controls []engine.Control) []engine.Control {
	reversed := append([]engine.Control(nil), controls...)

	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	return reversed
}

func Gate2A1K8sDiscoveryReview() engine.EventControl {
	return engine.EventControl{
		ID:        "review-k8s-discovery-agent",
		EventIDs:  []string{scenarios.EventGate2K8sDiscovery},
		Substrate: "independent-agent-reviewer",
		Requires: []engine.Condition{
			{
				Fact: Gate2AlertK8sDiscovery,
				Op:   engine.ConditionPresent,
			},
		},
		SuppressedBy: []engine.Condition{
			{
				Fact: scenarios.FactGate2AgentReviewerUnavailable,
				Op:   engine.ConditionPresent,
			},
		},
		Action: engine.ActionReview,
		Effects: []engine.Effect{
			{
				Fact: Gate2ReviewedClusterIntegrity,
				Op:   engine.EffectAdd,
			},
		},
	}
}

func Gate2CorrelatedReview() []engine.Control {
	return []engine.Control{
		Gate2V1CredentialTheft(),
		Gate2V3K8sAuditDiscovery(),

		Gate2H1CredentialHistoryReview(),
		Gate2H2K8sDiscoveryReview(),

		Gate2C1ClusterIntegrityCriticality(),
		Gate2E1ClusterIntegrityEscalation(),
		Gate2R1RevokeClusterCredential(),
	}
}

func Gate2DiverseReview() []engine.Control {
	return []engine.Control{
		Gate2V1CredentialTheft(),
		Gate2V3K8sAuditDiscovery(),

		Gate2H1CredentialHistoryReview(),
		Gate2A1K8sDiscoveryReview(),

		Gate2C1ClusterIntegrityCriticality(),
		Gate2E1ClusterIntegrityEscalation(),
		Gate2R1RevokeClusterCredential(),
	}
}

func Gate2CC() []engine.Control {
	return []engine.Control{
		// Correlated observation.
		Gate2V1CredentialTheft(),
		Gate2V2WorkerK8sDiscovery(),

		// Correlated review.
		Gate2H1CredentialHistoryReview(),
		Gate2H2K8sDiscoveryReview(),

		Gate2C1ClusterIntegrityCriticality(),
		Gate2E1ClusterIntegrityEscalation(),
		Gate2R1RevokeClusterCredential(),
	}
}

func Gate2CD() []engine.Control {
	return []engine.Control{
		// Correlated observation.
		Gate2V1CredentialTheft(),
		Gate2V2WorkerK8sDiscovery(),

		// Diverse review.
		Gate2H1CredentialHistoryReview(),
		Gate2A1K8sDiscoveryReview(),

		Gate2C1ClusterIntegrityCriticality(),
		Gate2E1ClusterIntegrityEscalation(),
		Gate2R1RevokeClusterCredential(),
	}
}

func Gate2DC() []engine.Control {
	return []engine.Control{
		// Diverse observation.
		Gate2V2WorkerK8sDiscovery(),
		Gate2V3K8sAuditDiscovery(),

		// Correlated review.
		Gate2H1CredentialHistoryReview(),
		Gate2H2K8sDiscoveryReview(),

		Gate2C1ClusterIntegrityCriticality(),
		Gate2E1ClusterIntegrityEscalation(),
		Gate2R1RevokeClusterCredential(),
	}
}

func Gate2DD() []engine.Control {
	return []engine.Control{
		// Diverse observation.
		Gate2V2WorkerK8sDiscovery(),
		Gate2V3K8sAuditDiscovery(),

		// Diverse review.
		Gate2H1CredentialHistoryReview(),
		Gate2A1K8sDiscoveryReview(),

		Gate2C1ClusterIntegrityCriticality(),
		Gate2E1ClusterIntegrityEscalation(),
		Gate2R1RevokeClusterCredential(),
	}
}
