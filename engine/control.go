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

	ActionReview   ControlAction = "review"
	ActionCritical ControlAction = "critical"
	ActionEscalate ControlAction = "escalate"

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

	// Conditions that must hold for this control to act.
	Requires []Condition

	Action  ControlAction
	Effects []Effect
}

func (c EventControl) Evaluate(event Event, s *state.State) ControlResult {
	matched := c.matches(event)

	result := ControlResult{
		ControlID: c.ID,
		Role:      c.Role,
		Matched:   matched,
		Action:    ActionNone,
	}

	if !matched {
		return result
	}

	for _, condition := range c.Requires {
		if !condition.SatisfiedBy(s) {
			return result
		}
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

func applyControls(
	event Event,
	current *state.State,
	controls []Control,
	action ControlAction,
	entry *TraceEntry,
) {
	for _, control := range controls {
		result := control.Evaluate(event, current)

		if !result.Matched || result.Action != action {
			continue
		}

		result.Acted = true

		for _, effect := range result.Effects {
			effect.Apply(current)
		}

		entry.Controls = append(entry.Controls, result)
	}
}
