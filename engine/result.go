package engine

type EventStatus string

const (
	EventApplied EventStatus = "applied"
	EventSkipped EventStatus = "skipped"
)

// TraceEntry records what happened when replay reached one Event.
type TraceEntry struct {
	Index int

	EventID string
	Status  EventStatus

	// Unsatisfied contains the preconditions that prevented a skipped
	// event from occurring.
	Unsatisfied []Condition

	// Effects contains the effects applied by a successful event.
	Effects []Effect

	// State snapshots make the trace self-contained and easy to inspect.
	Before []string
	After  []string
}

// RunResult is the complete deterministic result of replaying one Scenario.
type RunResult struct {
	ScenarioID string

	InitialState  []string
	TerminalState []string

	Trace []TraceEntry
}
