package engine

import "github.com/bubskee/muster/state"

func Replay(scenario Scenario, controls ...Control) RunResult {
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

		// Evaluate every control against the reachable event.
		for _, control := range controls {
			entry.Controls = append(
				entry.Controls,
				control.Evaluate(event, current),
			)
		}

		// Vedettes observe. Pickets may stop the event before its effects occur.
		for i := range entry.Controls {
			control := &entry.Controls[i]

			if !control.Matched {
				continue
			}

			switch control.Action {
			case ActionObserve:
				control.Acted = true
				entry.ObservedBy = append(
					entry.ObservedBy,
					control.ControlID,
				)

			case ActionBlock:
				control.Acted = true
				entry.BlockedBy = append(
					entry.BlockedBy,
					control.ControlID,
				)
			}
		}

		// Any successful Picket interception prevents the incident event.
		if len(entry.BlockedBy) > 0 {
			entry.Status = EventBlocked
			entry.After = current.Facts()

			result.Trace = append(result.Trace, entry)
			continue
		}

		// The incident event succeeds.
		event.Apply(current)

		entry.Status = EventApplied
		entry.Effects = append([]Effect(nil), event.Effects...)

		// Reserves respond after a successful event.
		for i := range entry.Controls {
			control := &entry.Controls[i]

			if !control.Matched || control.Action != ActionRespond {
				continue
			}

			control.Acted = true
			entry.RespondedBy = append(
				entry.RespondedBy,
				control.ControlID,
			)

			for _, effect := range control.Effects {
				effect.Apply(current)

				entry.ResponseEffects = append(
					entry.ResponseEffects,
					effect,
				)
			}
		}

		entry.After = current.Facts()
		result.Trace = append(result.Trace, entry)
	}

	result.TerminalState = current.Facts()

	return result
}
