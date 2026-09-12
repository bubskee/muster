package watchlines

import (
	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/scenarios"
)

// HFInterleaved models a small defensive composition with both observation
// and prevention along the worker-to-node path.
//
// The incident fixture is unchanged. Only the defensive composition differs
// from HFVedettesOnly.
func HFInterleaved() []engine.Control {
	return []engine.Control{
		engine.EventControl{
			ID:       "worker-execution-detector",
			Role:     engine.Vedette,
			EventIDs: []string{scenarios.EventWorkerCodeExecution},
			Action:   engine.ActionObserve,
		},
		engine.EventControl{
			ID:       "service-account-credential-guard",
			Role:     engine.Picket,
			EventIDs: []string{scenarios.EventServiceAccountAccess},
			Action:   engine.ActionBlock,
		},
		engine.EventControl{
			ID:       "kubernetes-discovery-detector",
			Role:     engine.Vedette,
			EventIDs: []string{scenarios.EventKubernetesAPIDiscovery},
			Action:   engine.ActionObserve,
		},
		engine.EventControl{
			ID:       "privileged-hostpath-admission",
			Role:     engine.Picket,
			EventIDs: []string{scenarios.EventPrivilegedHostPath},
			Action:   engine.ActionBlock,
		},
	}
}
