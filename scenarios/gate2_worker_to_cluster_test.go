package scenarios

import (
	"reflect"
	"testing"

	"github.com/bubskee/muster/engine"
)

func TestGate2UndefendedTrajectoriesReachNode(t *testing.T) {
	for _, scenario := range []engine.Scenario{
		Gate2NoFailure(),
		Gate2EarlyFailure(),
		Gate2BetweenFailure(),
	} {
		t.Run(scenario.ID, func(t *testing.T) {
			result := engine.Replay(scenario)
			if !gate2HasFact(result.TerminalState, FactGate2NodeAccess) {
				t.Fatalf("terminal state = %v, want %q", result.TerminalState, FactGate2NodeAccess)
			}
		})
	}
}

func TestGate2FailureTimingsReachEquivalentWorldState(t *testing.T) {
	early := engine.Replay(Gate2EarlyFailure())
	between := engine.Replay(Gate2BetweenFailure())

	earlyCheckpoint := early.Trace[2].After
	betweenCheckpoint := between.Trace[2].After

	if !reflect.DeepEqual(earlyCheckpoint, betweenCheckpoint) {
		t.Fatalf("checkpoint states differ:\nearly:   %v\nbetween: %v", earlyCheckpoint, betweenCheckpoint)
	}
}

func gate2HasFact(facts []string, want string) bool {
	for _, fact := range facts {
		if fact == want {
			return true
		}
	}
	return false
}
