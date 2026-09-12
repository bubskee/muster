package engine

import "github.com/bubskee/muster/state"

// Replay executes a Scenario from its initial state.
//
// Events are considered in order. An Event whose preconditions are not met is
// recorded as skipped, but replay continues. This allows later independent
// Events to remain reachable without introducing graph semantics.
func Replay(scenario Scenario) RunResult {
	current := state.New(scenario.InitialFacts...)

	result := RunResult{
		ScenarioID:   scenario.ID,
		InitialState: current.Facts(),
		Trace:        make([]TraceEntry, 0, len(scenario.Events)),
	}

	for i, event := range scenario.Events {
		entry := TraceEntry{
			Index:   i,
			EventID: event.ID,
			Before:  current.Facts(),
		}

		unsatisfied := event.UnsatisfiedPreconditions(current)
		if len(unsatisfied) > 0 {
			entry.Status = EventSkipped
			entry.Unsatisfied = append([]Condition(nil), unsatisfied...)
			entry.After = current.Facts()

			result.Trace = append(result.Trace, entry)
			continue
		}

		event.Apply(current)

		entry.Status = EventApplied
		entry.Effects = append([]Effect(nil), event.Effects...)
		entry.After = current.Facts()

		result.Trace = append(result.Trace, entry)
	}

	result.TerminalState = current.Facts()

	return result
}
