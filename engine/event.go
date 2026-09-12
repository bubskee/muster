package engine

import "github.com/bubskee/muster/state"

type ConditionOp string

const (
	ConditionPresent ConditionOp = "present"
	ConditionAbsent  ConditionOp = "absent"
)

// Condition describes something that must be true before an Event can occur.
type Condition struct {
	Fact string
	Op   ConditionOp
}

func (c Condition) SatisfiedBy(s *state.State) bool {
	switch c.Op {
	case ConditionPresent:
		return s.Has(c.Fact)
	case ConditionAbsent:
		return !s.Has(c.Fact)
	default:
		return false
	}
}

type EffectOp string

const (
	EffectAdd    EffectOp = "add"
	EffectRemove EffectOp = "remove"
)

// Effect describes a semantic state change caused by an Event.
type Effect struct {
	Fact string
	Op   EffectOp
}

func (e Effect) Apply(s *state.State) {
	switch e.Op {
	case EffectAdd:
		s.Add(e.Fact)
	case EffectRemove:
		s.Remove(e.Fact)
	}
}

// Event is one semantic transition in an incident replay.
//
// At this stage an Event contains only what is necessary to determine whether
// it can occur and what state changes if it does.
type Event struct {
	ID            string
	Preconditions []Condition
	Effects       []Effect
}

// UnsatisfiedPreconditions returns the conditions that do not currently hold.
func (e Event) UnsatisfiedPreconditions(s *state.State) []Condition {
	var unsatisfied []Condition

	for _, condition := range e.Preconditions {
		if !condition.SatisfiedBy(s) {
			unsatisfied = append(unsatisfied, condition)
		}
	}

	return unsatisfied
}

// CanApply reports whether all of the Event's preconditions currently hold.
func (e Event) CanApply(s *state.State) bool {
	return len(e.UnsatisfiedPreconditions(s)) == 0
}

// Apply applies the Event's effects if all preconditions hold.
//
// The bool reports whether the Event occurred.
func (e Event) Apply(s *state.State) bool {
	if !e.CanApply(s) {
		return false
	}

	for _, effect := range e.Effects {
		effect.Apply(s)
	}

	return true
}
