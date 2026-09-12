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

	Unsatisfied []Condition
	Effects     []Effect

	// Controls records defensive controls evaluated for this reachable event.
	Controls []ControlResult

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
