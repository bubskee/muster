package engine

import "fmt"

type RunSummary struct {
	ScenarioID string

	FurthestEvent    string
	FirstObservation string
	FirstBlock       string

	EventsAfterDetection int

	TerminalState []string
}

func Summarize(result RunResult) RunSummary {
	summary := RunSummary{
		ScenarioID:    result.ScenarioID,
		TerminalState: append([]string(nil), result.TerminalState...),
	}

	firstObservedIndex := -1

	for i, entry := range result.Trace {
		if entry.Status == EventApplied || entry.Status == EventBlocked {
			summary.FurthestEvent = entry.EventID
		}

		if summary.FirstObservation == "" && len(entry.ObservedBy) > 0 {
			summary.FirstObservation = entry.EventID
			firstObservedIndex = i
		}

		if summary.FirstBlock == "" && len(entry.BlockedBy) > 0 {
			summary.FirstBlock = entry.EventID
		}
	}

	if firstObservedIndex >= 0 {
		summary.EventsAfterDetection =
			len(result.Trace) - firstObservedIndex - 1
	}

	return summary
}

func (s RunSummary) String() string {
	return fmt.Sprintf(
		"scenario: %s\n"+
			"furthest_event: %s\n"+
			"first_observation: %s\n"+
			"first_block: %s\n"+
			"events_after_detection: %d\n"+
			"terminal_state: %v\n",
		s.ScenarioID,
		s.FurthestEvent,
		s.FirstObservation,
		s.FirstBlock,
		s.EventsAfterDetection,
		s.TerminalState,
	)
}
