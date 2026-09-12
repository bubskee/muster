package engine

// Scenario is an ordered incident trace with a known initial state.
//
// InitialFacts rather than a *state.State keeps scenarios declarative and
// makes them straightforward to serialize later.
type Scenario struct {
	ID           string
	InitialFacts []string
	Events       []Event
}
