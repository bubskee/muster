package engine

import "github.com/bubskee/muster/state"

func recordControlResult(
	entry *TraceEntry,
	result ControlResult,
) {
	for i := range entry.Controls {
		if entry.Controls[i].ControlID == result.ControlID {
			entry.Controls[i] = result
			return
		}
	}

	entry.Controls = append(entry.Controls, result)
}

func evaluateAction(
	event Event,
	current *state.State,
	controls []Control,
	action ControlAction,
) []ControlResult {
	var results []ControlResult

	for _, control := range controls {
		result := control.Evaluate(event, current)

		if !result.Matched || result.Action != action {
			continue
		}

		result.Acted = true

		for _, effect := range result.Effects {
			effect.Apply(current)
		}

		results = append(results, result)
	}

	return results
}

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
		// Record every control considered for this reachable event.
		//
		// This is the initial snapshot. Later phases may re-evaluate and
		// overwrite the result as state changes.
		for _, control := range controls {
			recordControlResult(
				&entry,
				control.Evaluate(event, current),
			)
		}

		// Phase 1: Vedettes observe the attempted event.
		//
		// Observation may produce state, e.g.
		// alert:credential-access.
		observations := evaluateAction(
			event,
			current,
			controls,
			ActionObserve,
		)

		for _, control := range observations {
			recordControlResult(&entry, control)

			entry.ObservedBy = append(
				entry.ObservedBy,
				control.ControlID,
			)
		}

		// Phase 2: Pickets may stop the event before its effects occur.
		blocks := evaluateAction(
			event,
			current,
			controls,
			ActionBlock,
		)

		for _, control := range blocks {
			recordControlResult(&entry, control)

			entry.BlockedBy = append(
				entry.BlockedBy,
				control.ControlID,
			)
		}

		if len(entry.BlockedBy) > 0 {
			entry.Status = EventBlocked
			entry.After = current.Facts()

			result.Trace = append(result.Trace, entry)
			continue
		}

		// Phase 3: the incident event succeeds.
		event.Apply(current)

		entry.Status = EventApplied
		entry.Effects = append([]Effect(nil), event.Effects...)

		// Phase 4: reviewers consume surfaced alerts.
		reviews := evaluateAction(
			event,
			current,
			controls,
			ActionReview,
		)

		for _, control := range reviews {
			recordControlResult(&entry, control)
		}

		// Phase 5: reviewed signals are assessed for criticality.
		criticality := evaluateAction(
			event,
			current,
			controls,
			ActionCritical,
		)

		for _, control := range criticality {
			recordControlResult(&entry, control)
		}

		// Phase 6: critical signals may be escalated.
		escalations := evaluateAction(
			event,
			current,
			controls,
			ActionEscalate,
		)

		for _, control := range escalations {
			recordControlResult(&entry, control)
		}

		// Phase 7: Reserves may respond to an escalated incident.
		responses := evaluateAction(
			event,
			current,
			controls,
			ActionRespond,
		)

		for _, control := range responses {
			recordControlResult(&entry, control)

			entry.RespondedBy = append(
				entry.RespondedBy,
				control.ControlID,
			)

			entry.ResponseEffects = append(
				entry.ResponseEffects,
				control.Effects...,
			)
		}

		entry.After = current.Facts()
		result.Trace = append(result.Trace, entry)
	}

	return result
}
