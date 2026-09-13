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

type ControlDisposition string

const (
	DispositionUnmatched ControlDisposition = "unmatched"
	DispositionWaiting   ControlDisposition = "waiting"
	DispositionReady     ControlDisposition = "ready"
)

type ControlResult struct {
	ControlID string
	Role      Role
	Matched   bool

	Disposition ControlDisposition
	Reason      string

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
		ControlID:   c.ID,
		Role:        c.Role,
		Matched:     matched,
		Disposition: DispositionUnmatched,
		Action:      ActionNone,
	}

	if !matched {
		return result
	}

	// The control is relevant to this event, but may still be waiting
	// on causal inputs produced by an earlier phase.
	result.Disposition = DispositionWaiting

	for _, condition := range c.Requires {
		if !condition.SatisfiedBy(s) {
			result.Reason = unsatisfiedConditionReason(condition)
			return result
		}
	}

	result.Disposition = DispositionReady
	result.Action = c.Action
	result.Effects = append([]Effect(nil), c.Effects...)

	return result
}

func unsatisfiedConditionReason(condition Condition) string {
	switch condition.Op {
	case ConditionPresent:
		return "missing:" + condition.Fact
	case ConditionAbsent:
		return "present:" + condition.Fact
	default:
		return "unsatisfied:" + condition.Fact
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
