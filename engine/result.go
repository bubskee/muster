package engine

type EventStatus string

const (
	EventApplied EventStatus = "applied"
	EventSkipped EventStatus = "skipped"
	EventBlocked EventStatus = "blocked"
)

type TraceEntry struct {
	Index int

	EventID string
	Status  EventStatus

	Unsatisfied []Condition

	// Event effects that actually occurred.
	Effects []Effect

	// All controls evaluated against this event.
	Controls []ControlResult

	// Convenience summaries of controls that actually acted.
	ObservedBy  []string
	BlockedBy   []string
	RespondedBy []string

	// State changes caused by Reserve responses rather than the event itself.
	ResponseEffects []Effect

	Before []string
	After  []string
}

type RunResult struct {
	ScenarioID string

	InitialState  []string
	TerminalState []string

	Trace []TraceEntry
}
