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
	DispositionUnmatched  ControlDisposition = "unmatched"
	DispositionWaiting    ControlDisposition = "waiting"
	DispositionReady      ControlDisposition = "ready"
	DispositionSuppressed ControlDisposition = "suppressed"
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

	Requires []Condition

	// Substrate names the failure domain this control depends on.
	// It is descriptive; SuppressedBy defines the executable failure condition.
	Substrate string

	// If any of these conditions are satisfied, the control cannot act.
	SuppressedBy []Condition

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

	result.Disposition = DispositionWaiting

	for _, condition := range c.Requires {
		if !condition.SatisfiedBy(s) {
			result.Reason = unsatisfiedConditionReason(condition)
			return result
		}
	}

	for _, condition := range c.SuppressedBy {
		if condition.SatisfiedBy(s) {
			result.Disposition = DispositionSuppressed
			result.Reason = suppressionReason(condition)
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

// Phase reports the replay phase in which this control is evaluated.
func (c EventControl) Phase() ControlAction {
	return c.Action
}

func suppressionReason(condition Condition) string {
	switch condition.Op {
	case ConditionPresent:
		return condition.Fact
	case ConditionAbsent:
		return "missing:" + condition.Fact
	default:
		return "suppressed:" + condition.Fact
	}
}
