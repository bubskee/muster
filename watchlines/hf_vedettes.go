package watchlines

import (
	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/scenarios"
)

// HFVedettesOnly models strong observation along the worker-to-node path.
//
// It deliberately contains no prevention or recovery controls.
func HFVedettesOnly() []engine.Control {
	return []engine.Control{
		engine.EventControl{
			ID:       "worker-execution-detector",
			Role:     engine.Vedette,
			EventIDs: []string{scenarios.EventWorkerCodeExecution},
			Action:   engine.ActionObserve,
		},
		engine.EventControl{
			ID:       "credential-access-detector",
			Role:     engine.Vedette,
			EventIDs: []string{scenarios.EventServiceAccountAccess},
			Action:   engine.ActionObserve,
		},
		engine.EventControl{
			ID:       "kubernetes-discovery-detector",
			Role:     engine.Vedette,
			EventIDs: []string{scenarios.EventKubernetesAPIDiscovery},
			Action:   engine.ActionObserve,
		},
	}
}
