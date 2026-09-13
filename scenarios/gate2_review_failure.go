package scenarios

import "github.com/bubskee/muster/engine"

const (
	Gate2SOCReviewFailureID   = "gate2-review-soc-failure"
	Gate2AgentReviewFailureID = "gate2-review-agent-failure"

	FactGate2SOCQueueOverloaded       = "overloaded:soc-queue"
	FactGate2AgentReviewerUnavailable = "unavailable:independent-agent-reviewer"
)

func Gate2SOCReviewFailure() engine.Scenario {
	scenario := Gate2NoFailure()
	scenario.ID = Gate2SOCReviewFailureID
	scenario.InitialFacts = append(
		scenario.InitialFacts,
		FactGate2SOCQueueOverloaded,
	)
	return scenario
}

func Gate2AgentReviewFailure() engine.Scenario {
	scenario := Gate2NoFailure()
	scenario.ID = Gate2AgentReviewFailureID
	scenario.InitialFacts = append(
		scenario.InitialFacts,
		FactGate2AgentReviewerUnavailable,
	)
	return scenario
}

func Gate2BetweenSOCReviewFailure() engine.Scenario {
	scenario := Gate2BetweenFailure()
	scenario.ID = "gate2-between-soc-review-failure"
	scenario.InitialFacts = append(
		scenario.InitialFacts,
		FactGate2SOCQueueOverloaded,
	)
	return scenario
}

func Gate2BetweenAgentReviewFailure() engine.Scenario {
	scenario := Gate2BetweenFailure()
	scenario.ID = "gate2-between-agent-review-failure"
	scenario.InitialFacts = append(
		scenario.InitialFacts,
		FactGate2AgentReviewerUnavailable,
	)
	return scenario
}
