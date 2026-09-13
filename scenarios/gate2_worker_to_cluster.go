package scenarios

import "github.com/bubskee/muster/engine"

const (
	Gate2NoFailureID      = "gate2-worker-cluster-none"
	Gate2EarlyFailureID   = "gate2-worker-cluster-early"
	Gate2BetweenFailureID = "gate2-worker-cluster-between"

	EventGate2WorkerExecution     = "worker-code-execution"
	EventGate2CredentialTheft     = "cluster-credential-theft"
	EventGate2K8sDiscovery        = "kubernetes-api-discovery-from-worker"
	EventGate2CredentialReplay    = "stolen-credential-replay"
	EventGate2PrivilegedWorkload  = "privileged-workload-creation"
	EventGate2NodeAccess          = "node-access"
	EventGate2TelemetryImpairment = "worker-telemetry-impairment"
)

const (
	FactGate2MaliciousDataset           = "dataset:malicious-on-worker"
	FactGate2WorkerAccess               = "attacker:worker-access"
	FactGate2ClusterCredential          = "attacker:cluster-credential"
	FactGate2ClusterKnowledge           = "attacker:cluster-knowledge"
	FactGate2ClusterSession             = "attacker:cluster-session"
	FactGate2PrivilegedWorkload         = "workload:privileged-hostpath"
	FactGate2NodeAccess                 = "attacker:node-access"
	FactGate2WorkerTelemetryCompromised = "compromised:worker-telemetry"
)

func Gate2NoFailure() engine.Scenario {
	return gate2Scenario(Gate2NoFailureID, []engine.Event{
		gate2WorkerExecution(),
		gate2CredentialTheft(),
		gate2K8sDiscovery(),
		gate2CredentialReplay(),
		gate2PrivilegedWorkload(),
		gate2NodeAccess(),
	})
}

func Gate2EarlyFailure() engine.Scenario {
	return gate2Scenario(Gate2EarlyFailureID, []engine.Event{
		gate2WorkerExecution(),
		gate2TelemetryImpairment(),
		gate2CredentialTheft(),
		gate2K8sDiscovery(),
		gate2CredentialReplay(),
		gate2PrivilegedWorkload(),
		gate2NodeAccess(),
	})
}

func Gate2BetweenFailure() engine.Scenario {
	return gate2Scenario(Gate2BetweenFailureID, []engine.Event{
		gate2WorkerExecution(),
		gate2CredentialTheft(),
		gate2TelemetryImpairment(),
		gate2K8sDiscovery(),
		gate2CredentialReplay(),
		gate2PrivilegedWorkload(),
		gate2NodeAccess(),
	})
}

func gate2Scenario(id string, events []engine.Event) engine.Scenario {
	return engine.Scenario{
		ID:           id,
		InitialFacts: []string{FactGate2MaliciousDataset},
		Events:       events,
	}
}

func gate2WorkerExecution() engine.Event {
	return engine.Event{
		ID:            EventGate2WorkerExecution,
		Preconditions: []engine.Condition{{Fact: FactGate2MaliciousDataset, Op: engine.ConditionPresent}},
		Effects:       []engine.Effect{{Fact: FactGate2WorkerAccess, Op: engine.EffectAdd}},
	}
}

func gate2CredentialTheft() engine.Event {
	return engine.Event{
		ID:            EventGate2CredentialTheft,
		Preconditions: []engine.Condition{{Fact: FactGate2WorkerAccess, Op: engine.ConditionPresent}},
		Effects:       []engine.Effect{{Fact: FactGate2ClusterCredential, Op: engine.EffectAdd}},
	}
}

func gate2K8sDiscovery() engine.Event {
	return engine.Event{
		ID: EventGate2K8sDiscovery,
		Preconditions: []engine.Condition{
			{Fact: FactGate2WorkerAccess, Op: engine.ConditionPresent},
			{Fact: FactGate2ClusterCredential, Op: engine.ConditionPresent},
		},
		Effects: []engine.Effect{{Fact: FactGate2ClusterKnowledge, Op: engine.EffectAdd}},
	}
}

func gate2CredentialReplay() engine.Event {
	return engine.Event{
		ID:            EventGate2CredentialReplay,
		Preconditions: []engine.Condition{{Fact: FactGate2ClusterCredential, Op: engine.ConditionPresent}},
		Effects:       []engine.Effect{{Fact: FactGate2ClusterSession, Op: engine.EffectAdd}},
	}
}

func gate2PrivilegedWorkload() engine.Event {
	return engine.Event{
		ID:            EventGate2PrivilegedWorkload,
		Preconditions: []engine.Condition{{Fact: FactGate2ClusterSession, Op: engine.ConditionPresent}},
		Effects:       []engine.Effect{{Fact: FactGate2PrivilegedWorkload, Op: engine.EffectAdd}},
	}
}

func gate2NodeAccess() engine.Event {
	return engine.Event{
		ID:            EventGate2NodeAccess,
		Preconditions: []engine.Condition{{Fact: FactGate2PrivilegedWorkload, Op: engine.ConditionPresent}},
		Effects:       []engine.Effect{{Fact: FactGate2NodeAccess, Op: engine.EffectAdd}},
	}
}

func gate2TelemetryImpairment() engine.Event {
	return engine.Event{
		ID:            EventGate2TelemetryImpairment,
		Preconditions: []engine.Condition{{Fact: FactGate2WorkerAccess, Op: engine.ConditionPresent}},
		Effects:       []engine.Effect{{Fact: FactGate2WorkerTelemetryCompromised, Op: engine.EffectAdd}},
	}
}
