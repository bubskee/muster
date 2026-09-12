package scenarios

import "github.com/bubskee/muster/engine"

const (
	HFWorkerToNodeID = "hf-worker-to-node"

	EventMaliciousDataset       = "malicious-dataset-reaches-worker"
	EventWorkerCodeExecution    = "worker-code-execution"
	EventServiceAccountAccess   = "service-account-credential-access"
	EventKubernetesAPIDiscovery = "kubernetes-api-discovery"
	EventPrivilegedHostPath     = "privileged-hostpath-workload"
	EventNodeAccess             = "node-access"
)

// HFWorkerToNode returns a deliberately small semantic replay derived from the
// July 2026 Hugging Face incident.
//
// It is not intended to reproduce the full incident or its telemetry. It
// isolates one useful path:
//
// malicious dataset
//
//	→ worker code execution
//	→ service-account credential access
//	→ Kubernetes discovery
//	→ privileged hostPath workload
//	→ node access
func HFWorkerToNode() engine.Scenario {
	return engine.Scenario{
		ID: HFWorkerToNodeID,
		Events: []engine.Event{
			{
				ID: EventMaliciousDataset,
				Effects: []engine.Effect{
					{
						Fact: "dataset:malicious-on-worker",
						Op:   engine.EffectAdd,
					},
				},
			},
			{
				ID: EventWorkerCodeExecution,
				Preconditions: []engine.Condition{
					{
						Fact: "dataset:malicious-on-worker",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "access:worker",
						Op:   engine.EffectAdd,
					},
				},
			},
			{
				ID: EventServiceAccountAccess,
				Preconditions: []engine.Condition{
					{
						Fact: "access:worker",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "credential:k8s-service-account",
						Op:   engine.EffectAdd,
					},
				},
			},
			{
				ID: EventKubernetesAPIDiscovery,
				Preconditions: []engine.Condition{
					{
						Fact: "credential:k8s-service-account",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "discovery:k8s-api",
						Op:   engine.EffectAdd,
					},
				},
			},
			{
				ID: EventPrivilegedHostPath,
				Preconditions: []engine.Condition{
					{
						Fact: "credential:k8s-service-account",
						Op:   engine.ConditionPresent,
					},
					{
						Fact: "discovery:k8s-api",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "workload:privileged-hostpath",
						Op:   engine.EffectAdd,
					},
				},
			},
			{
				ID: EventNodeAccess,
				Preconditions: []engine.Condition{
					{
						Fact: "workload:privileged-hostpath",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "access:node",
						Op:   engine.EffectAdd,
					},
				},
			},
		},
	}
}
