package engine

import (
	"fmt"
	"strings"
)

type RunSummary struct {
	ScenarioID string

	FurthestApplied   string
	FurthestAttempted string

	FirstObservation string
	FirstBlock       string

	// Nil means the incident was never observed.
	// A pointer to 0 means it was observed, but no later events succeeded.
	EventsAfterDetection *int

	TerminalState []string
}

func Summarize(result RunResult) RunSummary {
	summary := RunSummary{
		ScenarioID:    result.ScenarioID,
		TerminalState: append([]string(nil), result.TerminalState...),
	}

	firstObservedIndex := -1

	for i, entry := range result.Trace {
		switch entry.Status {
		case EventApplied:
			summary.FurthestApplied = entry.EventID
			summary.FurthestAttempted = entry.EventID

		case EventBlocked:
			// The event was attempted but did not occur.
			summary.FurthestAttempted = entry.EventID
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
		count := 0

		for _, entry := range result.Trace[firstObservedIndex+1:] {
			if entry.Status == EventApplied {
				count++
			}
		}

		summary.EventsAfterDetection = &count
	}

	return summary
}

func (s RunSummary) String() string {
	eventsAfterDetection := "n/a"
	if s.EventsAfterDetection != nil {
		eventsAfterDetection = fmt.Sprintf("%d", *s.EventsAfterDetection)
	}

	return fmt.Sprintf(
		"scenario: %s\n"+
			"furthest_applied: %s\n"+
			"furthest_attempted: %s\n"+
			"first_observation: %s\n"+
			"first_block: %s\n"+
			"events_after_detection: %s\n"+
			"terminal_state: [%s]\n",
		s.ScenarioID,
		s.FurthestApplied,
		s.FurthestAttempted,
		s.FirstObservation,
		s.FirstBlock,
		eventsAfterDetection,
		strings.Join(s.TerminalState, ", "),
	)
}
