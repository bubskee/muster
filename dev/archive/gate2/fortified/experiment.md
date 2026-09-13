# Gate 2 Experiment Design

## Status

**FROZEN BEFORE IMPLEMENTATION / EXECUTION**

This document defines the permitted Gate 2 compositions and comparisons over the already-frozen incident and control semantics.

The experiments use **structural slot budgets**, not claims about real monetary, engineering, or operational cost.

Observation, review, and response capacity are budgeted separately.

$$
B=(O,H,R)
$$

where:

* \(O\) = observation slots
* \(H\) = review / interpretation slots
* \(R\) = response capability

Comparisons vary one budget dimension at a time.

---

# Common rules

All experiments use the same frozen incident:

```text
E1 worker-code-execution
E2 cluster-credential-theft
E3 kubernetes-api-discovery-from-worker
E4 stolen-credential-replay
E5 privileged-workload-creation
E6 node-access

F1 worker-telemetry-impairment
```

with the three frozen F1 treatments:

```text
NONE
EARLY
BETWEEN
```

The primary adverse terminal predicate is:

```text
attacker:node-access
```

All control semantics are fixed before execution.

No composition may alter:

```text
event semantics
control event coverage
control Effects
control Requires
substrate dependencies
Reserve behavior
terminal outcome definition
```

---

# Gate 2A — Observation Budget

## Question

Given a fixed downstream review/response system and two observation slots:

> Should the second observation investment buy additional incident coverage on a shared evidence substrate, or redundant coverage on an independent evidence substrate?

## Budget

$$
B_{2A}=(O=2,\;H=\text{fixed},\;R=\text{fixed})
$$

Only the two Vedette slots vary.

---

## Composition A — Broad / Correlated

```text
V1 credential-theft detector
    event:     E2
    substrate: worker-telemetry

V2 worker-side k8s discovery detector
    event:     E3
    substrate: worker-telemetry
```

Properties:

```text
unique event coverage: E2 + E3
evidence domains:       1
```

---

## Composition B — Narrow / Diverse

```text
V2 worker-side k8s discovery detector
    event:     E3
    substrate: worker-telemetry

V3 k8s-audit discovery detector
    event:     E3
    substrate: k8s-audit
```

Properties:

```text
unique event coverage: E3
evidence domains:       2
```

V2 is common to both compositions.

The second slot therefore encodes exactly this trade:

```text
V1 = additional temporal/event coverage
V3 = evidence-source independence
```

---

## Downstream controls

The same downstream machinery is available to both compositions.

Controls that have no input in a given composition remain present but inert.

No downstream control is swapped to compensate for observation choice.

---

## Run matrix

```text
Broad/Correlated × NONE
Broad/Correlated × EARLY
Broad/Correlated × BETWEEN

Narrow/Diverse   × NONE
Narrow/Diverse   × EARLY
Narrow/Diverse   × BETWEEN
```

Six runs.

---

## Static null

Before executing Muster, predict all six using the frozen causal-pruning baseline.

Record:

```text
reachable events
available Vedettes
expected alerts
expected response
terminal attacker:node-access
A/B ranking
```

Shared-substrate failure alone does not count as a Muster result.

---

## Interesting result

The useful result is not:

```text
diverse survives telemetry failure
```

That is expected.

The interesting result requires ordering/history to change an outcome or ranking in a way not captured by static causal pruning.

---

# Gate 2B — Review Budget

## Question

Given fixed observation and response machinery and two review slots:

> Should review redundancy be concentrated in one operational failure domain, or split across independent human/agent review paths?

## Budget

$$
B_{2B}=(O=\text{fixed},\;H=2,\;R=\text{fixed})
$$

Observation is identical across compositions.

Response capability is identical across compositions.

Only review-path dependency structure varies.

---

## Fixed observation

Use the same frozen observation set in both arms.

The preferred set is:

```text
V1 credential-theft detector
V3 k8s-audit discovery detector
```

This supplies observations from two evidence domains without making observation structure itself the independent variable.

---

## Composition A — Correlated Review

Two review opportunities share one operational domain:

```text
H1 credential-history reviewer
    substrate: soc-queue / human-oncall

H2 k8s-discovery reviewer
    substrate: soc-queue / human-oncall
```

Both ultimately contribute to the same frozen criticality / escalation / response machinery.

---

## Composition B — Diverse Review

One review opportunity remains on the human/SOC path.

The other uses an independently failing agent-review path:

```text
H1 credential-history reviewer
    substrate: soc-queue / human-oncall

A1 k8s-discovery reviewer
    substrate: independent-agent-reviewer
```

H2 and A1 must:

```text
consume the same alert
produce the same reviewed fact
feed the same downstream machinery
```

They differ only in review substrate.

No semantics such as:

```text
agent is faster
human is smarter
agent is more accurate
```

are permitted.

Human/agent identity is used only to represent distinct operational failure domains.

---

## Review-failure treatments

Use a small frozen set:

```text
NONE

SOC-FAILURE
    correlated human/SOC review unavailable or suppressed

AGENT-FAILURE
    independent agent review unavailable or suppressed
```

If the existing Gate 2 machinery represents only one review-failure treatment initially, implement the minimum required facts without adding queue simulation, timing distributions, or reviewer quality parameters.

---

## Run matrix

```text
Correlated Review × NONE
Correlated Review × SOC-FAILURE
Correlated Review × AGENT-FAILURE

Diverse Review    × NONE
Diverse Review    × SOC-FAILURE
Diverse Review    × AGENT-FAILURE
```

The AGENT-FAILURE rows serve as a control against encoding:

```text
agent independence = automatically better
```

---

## Static null

Predict all rows before execution using ordinary dependency propagation.

If results reduce entirely to:

```text
reviewer depends on failed substrate
→ reviewer unavailable
```

Gate 2B is calibration only.

---

## HF relevance

Gate 2B is motivated by the Hugging Face observation that sensing/correlation and appropriate prioritization/escalation are distinct defensive functions.

It is not intended to claim that an independent agent reviewer would have prevented the real incident.

The human/agent split is an architectural counterfactual, not an HF historical reconstruction.

---

# Gate 2C — Integrated Composition Check

## Status

**PREREGISTERED BUT CONDITIONAL**

Do not run 2C merely because 2A and 2B exist.

Run 2C only if at least one of 2A or 2B produces evidence that stateful replay remains worth evaluating.

If both 2A and 2B collapse cleanly to the static baseline:

```text
STOP
COLLAPSE
```

Do not use 2C to rescue Muster.

---

## Question

When observation and review independence are both allowed to vary:

> Do their effects compose in the obvious static way, or does persistent defensive history create an interaction that neither factor reveals alone?

---

## Budget

$$
B_{2C}=(O=2,\;H=2,\;R=\text{fixed})
$$

This yields a small \(2\times2\) architecture matrix:

```text
                    REVIEW
               correlated   diverse

OBS correlated      CC         CD

    diverse         DC         DD
```

where:

```text
CC = correlated sensing + correlated review
CD = correlated sensing + diverse review
DC = diverse sensing    + correlated review
DD = diverse sensing    + diverse review
```

---

## Purpose

2C is **not** intended to show:

```text
DD is strongest
```

That result would be nearly tautological under independent failures.

The only reason to run 2C is to test whether:

```text
observation architecture
×
review architecture
×
persistent defensive state
```

produces an interaction not predicted by independently analyzing the two factors.

If the combined result equals the obvious composition of 2A and 2B:

```text
no additional Gate 2 evidence
```

---

# Separate non-monotonicity fixture

The premature-containment controls:

```text
L1 low-confidence escalation
R2 partial worker containment
```

are not part of the 2A/2B slot budget.

They belong to the separately preregistered **non-monotonicity challenge**.

That challenge asks:

> Can adding a functioning defensive path make the designated security outcome worse?

The test compares otherwise nested functioning-control sets:

$$
W_1 \subset W_2
$$

and checks whether:

$$
Outcome(W_2) < Outcome(W_1)
$$

under the frozen security ordering.

If the additional working control causes:

```text
early alert
→ partial containment
→ later evidence opportunity disappears
→ full response never triggers
→ node access succeeds
```

while its absence permits stronger later containment, the defensive system is non-monotone.

This fixture must not be mixed into 2A or 2B merely to manufacture a crossover.

---

# Required pre-execution artifacts

Before implementing/running the decisive matrix, produce:

```text
1. static-null prediction table for 2A
2. static-null prediction table for 2B
3. predicted monotonicity of the dedicated fixture
```

For each row record:

```text
expected reachable events
expected observations
expected review state
expected escalation
expected response
expected terminal node-access
expected winning composition, if any
reason
```

These predictions are frozen before observing Muster output.

---

# Stop rules

## COLLAPSE immediately if

2A and 2B are completely explained by the static dependency model **and**

the non-monotonicity/history fixture does not produce load-bearing persistent state.

Do not run 2C.

---

## Continue to 2C only if

2A or 2B exposes a replay-relevant history effect,

or the dedicated history/non-monotonicity fixture demonstrates that the state machine captures a phenomenon the static baseline does not.

---

# Interpretation discipline

A slot means:

> one normalized architectural opportunity within this experiment.

It does **not** mean equal:

```text
dollars
engineering hours
operational burden
latency
accuracy
maintenance cost
staffing requirement
```

No claim about real-world cost equivalence is made.

The purpose of the budget is to force explicit architectural trade-offs without inventing a quantitative security-economics model.
