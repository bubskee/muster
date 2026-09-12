package scenarios

import "github.com/bubskee/muster/engine"

const HFWorkerToNodeID = "hf-worker-to-node"

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
				ID: "malicious-dataset-reaches-worker",
				Effects: []engine.Effect{
					{
						Fact: "dataset:malicious-on-worker",
						Op:   engine.EffectAdd,
					},
				},
			},
			{
				ID: "worker-code-execution",
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
				ID: "service-account-credential-access",
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
				ID: "kubernetes-api-discovery",
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
				ID: "privileged-hostpath-workload",
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
				ID: "node-access",
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
