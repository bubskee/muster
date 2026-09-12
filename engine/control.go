package engine

import "github.com/bubskee/muster/state"

type Role string

const (
	Vedette Role = "vedette"
	Picket  Role = "picket"
	Reserve Role = "reserve"
)

type ControlAction string

const (
	ActionNone    ControlAction = "none"
	ActionObserve ControlAction = "observe"
	ActionBlock   ControlAction = "block"
	ActionRespond ControlAction = "respond"
)

type ControlResult struct {
	ControlID string
	Role      Role
	Matched   bool

	Action ControlAction
	Acted  bool

	Effects []Effect
}

type Control interface {
	Evaluate(Event, *state.State) ControlResult
}

type EventControl struct {
	ID       string
	Role     Role
	EventIDs []string

	Action  ControlAction
	Effects []Effect
}

func (c EventControl) Evaluate(event Event, _ *state.State) ControlResult {
	matched := c.matches(event)

	result := ControlResult{
		ControlID: c.ID,
		Role:      c.Role,
		Matched:   matched,
	}

	if !matched {
		return result
	}

	result.Action = c.Action
	result.Effects = append([]Effect(nil), c.Effects...)

	return result
}

func (c EventControl) matches(event Event) bool {
	for _, id := range c.EventIDs {
		if id == event.ID {
			return true
		}
	}

	return false
}
