package state

import "sort"

// State is the accumulated semantic state of a replay.
//
// Facts are deliberately opaque strings at this stage. The engine does not
// need to understand what "access:node" means; it only needs to know whether
// the fact is present.
//
// Examples:
//
//	access:worker
//	credential:k8s-service-account
//	discovery:k8s-api
//	access:node
type State struct {
	facts map[string]struct{}
}

func New(facts ...string) *State {
	s := &State{
		facts: make(map[string]struct{}, len(facts)),
	}

	for _, fact := range facts {
		s.Add(fact)
	}

	return s
}

// Has reports whether fact is currently present.
func (s *State) Has(fact string) bool {
	if s == nil {
		return false
	}

	_, ok := s.facts[fact]
	return ok
}

// Add makes fact present.
func (s *State) Add(fact string) {
	if s.facts == nil {
		s.facts = make(map[string]struct{})
	}

	s.facts[fact] = struct{}{}
}

// Remove makes fact absent.
func (s *State) Remove(fact string) {
	if s == nil {
		return
	}

	delete(s.facts, fact)
}

// Facts returns the current facts in deterministic order.
//
// Returning a copy rather than exposing the underlying map keeps callers from
// mutating State behind its back.
func (s *State) Facts() []string {
	if s == nil {
		return nil
	}

	facts := make([]string, 0, len(s.facts))
	for fact := range s.facts {
		facts = append(facts, fact)
	}

	sort.Strings(facts)
	return facts
}

// Clone returns an independent copy of State.
func (s *State) Clone() *State {
	clone := New()

	if s == nil {
		return clone
	}

	for fact := range s.facts {
		clone.Add(fact)
	}

	return clone
}
