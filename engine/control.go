package engine

import "github.com/bubskee/muster/state"

type Role string

const (
	Vedette Role = "vedette"
	Picket  Role = "picket"
	Reserve Role = "reserve"
)

// ControlResult records the result of evaluating one control against one event.
//
// For Commit 3, controls only report whether they match. They do not yet
// observe, block, or recover. Those semantics arrive in Commit 4.
type ControlResult struct {
	ControlID string
	Role      Role
	Matched   bool
}

// Control is one defensive mechanism in a Watchline.
//
// Evaluate receives the pre-event state. This matters later: a Picket should
// be able to judge an event before its effects are applied.
type Control interface {
	Evaluate(Event, *state.State) ControlResult
}

// EventControl is the first deliberately simple Control implementation.
//
// It matches explicit event IDs. Later we can grow matching to semantic kinds,
// tags, state conditions, etc. without changing the replay contract.
type EventControl struct {
	ID       string
	Role     Role
	EventIDs []string
}

func (c EventControl) Evaluate(event Event, _ *state.State) ControlResult {
	return ControlResult{
		ControlID: c.ID,
		Role:      c.Role,
		Matched:   c.matches(event),
	}
}

func (c EventControl) matches(event Event) bool {
	for _, id := range c.EventIDs {
		if id == event.ID {
			return true
		}
	}

	return false
}
