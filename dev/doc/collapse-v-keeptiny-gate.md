# Muster Gate 2 — Shared Observation and Review Failure Domains

## Why this gate exists

Gate 1 established that Muster's replay engine works, but did **not** establish that the current experiment is interesting.

The current `hf-worker-to-node` fixture is a linear chain. In that fixture:

- the Vedette-only result is tautological because observation is causally inert;
- defensive outcomes collapse largely to the location of the earliest blocking Picket;
- the terminal summary contains fields that are currently misleading;
- the Hugging Face case is a poor fit for a simple "detection worked, prevention failed" story.

Before expanding Muster, test a narrower and more falsifiable hypothesis:

> **Can a Watchline fail despite strong sensing because its observation and review paths share failure domains?**

This is intended to test something that ordinary earliest-cut reasoning does not automatically answer.

---

# Pre-gate hygiene

Before running the experiment, fix the existing summary layer.

## 1. Fix `events_after_detection`

Current behavior is misleading because skipped trace entries are counted and an undefended run can score `0`.

Change the semantics to:

- undefined / `nil` if no event was ever observed;
- otherwise count only `EventApplied` events after the first successful observation;
- do not count skipped or blocked trace entries.

The absence of detection must never look like the best possible detection score.

## 2. Split `furthest_event`

Replace the ambiguous field with:

```text
furthest_applied
furthest_attempted
```

A blocked event may be the furthest attempted transition without ever having occurred.

## 3. Do not add additional headline metrics

These are correctness fixes, not an invitation to build a scoring system.

`TerminalState` remains the most trustworthy outcome representation.

---

# Updated causal model

The current engine effectively models:

```text
incident event
    ↓
Vedette matches
    ↓
nothing necessarily happens
```

That is too weak for the hypothesis under test.

The smallest acceptable causal path for Gate 2 is:

```text
world event
    ↓
Vedette detects
    ↓
signal is surfaced
    ↓
human or agent reviewer notices / reads / interprets
    ↓
criticality is assessed from the threatened security outcome
    ↓
signal is escalated
    ↓
Reserve may respond
    ↓
response changes incident state
```

The exact implementation may remain very small.

The purpose is **not** to build a SOC simulator.

The purpose is to make this distinction executable:

```text
no detection
≠
detected but unread
≠
read but under-prioritized
≠
recognized as critical but not escalated
≠
escalated and acted upon
```

---

# Gate hypothesis

> Defensive composition can matter because nominally redundant controls may share sensing or review failure domains.

A Watchline with more detectors should not necessarily outperform one with fewer detectors if:

- the detectors depend on the same compromised telemetry substrate;
- their alerts flow through the same failed or overloaded review path;
- their reviewers fail to recognize that the observed condition threatens a critical security outcome;
- or some combination of these.

The variable under test is **independence of observation-to-action paths**, not simply blocker position.

---

# Minimal new semantics

Only add what is necessary to test this hypothesis.

## Evidence substrate

A Vedette may declare where its evidence comes from.

Examples:

```text
worker-telemetry
host-audit
k8s-audit
cloud-control-plane
```

The incident state may contain facts such as:

```text
compromised:worker-telemetry
```

A Vedette whose evidence substrate is compromised cannot produce a trustworthy signal.

The trace must distinguish:

```text
did not match
```

from:

```text
would have matched, but evidence substrate failed
```

Do not collapse both cases into `Matched: false`.

---

## Review substrate

A surfaced signal may also depend on a reviewer or review channel.

Examples:

```text
human-oncall
soc-queue
agent-reviewer
independent-security-agent
```

The incident or environment may contain facts such as:

```text
unavailable:human-oncall
overloaded:soc-queue
compromised:agent-reviewer
```

A signal can therefore exist without becoming reviewed or escalated.

This is the executable version of:

> "Didn't look, didn't read."

---

## Alert / review / criticality / escalation state

Use ordinary Muster state facts where possible.

Illustrative facts:

```text
alert:worker-code-execution
reviewed:worker-code-execution
critical:worker-code-execution
escalated:worker-code-execution
```

The exact naming is flexible.

The important rule is causal:

> A Reserve must not respond merely because a Vedette matched.

A Reserve response must depend on an appropriate escalated fact.

### Criticality rule

Do **not** make criticality a property of the Vedette itself.

The Vedette answers:

> What happened?

Criticality answers:

> Does what happened threaten an outcome important enough to require action?

For Gate 2, keep this deliberately coarse. A Boolean-like distinction is enough:

```text
critical:<signal>
```

versus the absence of that fact.

Do not build a general severity scale.

Criticality should be derived from the threatened security objective or journey, not from the loudness of the individual alert. This allows the gate to represent an important failure mode:

```text
alert exists
    ↓
alert is reviewed correctly
    ↓
criticality is assigned too low
    ↓
no escalation
```

That is distinct from both sensor failure and "didn't look, didn't read."

---

# What not to model yet

Do **not** add:

- stochastic reviewer attention;
- realistic alert queues;
- human reaction-time distributions;
- multi-level severity scoring;
- numerical risk scores;
- automatic blast-radius-to-severity formulas;
- dynamic priority queues;
- false-positive simulation;
- false-negative probabilities;
- adaptive attacker agents;
- reinforcement learning;
- generalized DAG infrastructure;
- network topology;
- SIEM emulation;
- real-time scheduling;
- a general trust-domain framework.

Keep the implementation embarrassingly small.

---

# Experimental compositions

Hold control count and overall defensive "budget" approximately fixed.

The exact controls may change to fit the final sourced or synthetic micro.

## Watchline A — correlated strength

A superficially strong configuration:

```text
V worker-exec detector      → evidence: worker-telemetry
                              review:   soc-queue

V credential detector       → evidence: worker-telemetry
                              review:   soc-queue

V k8s discovery detector    → evidence: worker-telemetry
                              review:   soc-queue

R credential revocation     → requires critical escalation
```

This Watchline has many detectors but shares both sensing and review failure domains.

---

## Watchline B — independent coverage

A superficially smaller configuration:

```text
V worker-exec detector      → evidence: worker-telemetry
                              review:   agent-reviewer

V credential detector       → evidence: host-audit
                              review:   human-oncall

R credential revocation     → requires critical escalation
```

This Watchline has fewer detectors but more independent paths from observation to action.

---

# Incident micro

Keep the incident approximately fixed.

Add the **smallest** event or condition necessary to compromise or disable one observation or review substrate.

Prefer a historically sourced transition when the public incident supports it.

If the Hugging Face record does not cleanly support the required failure mode:

> Use a clearly labeled synthetic micro.

Do not bend the incident history to fit Watchline.

The purpose of Gate 2 is to test the abstraction, not to force Hugging Face to prove it.

---

# Required cases

At minimum, run these cases against both Watchlines.

## Case A — evidence failure

Example:

```text
compromised:worker-telemetry
```

Expected effect:

- controls sharing worker telemetry lose observation;
- independently sourced observation may survive.

---

## Case B — review / prioritization failure

Gate 2 must be able to distinguish two post-detection failures.

### B1 — unread

Example:

```text
overloaded:soc-queue
```

Expected effect:

- alerts may still be produced;
- signals routed only through that review path are not meaningfully consumed;
- no criticality assessment or escalation occurs for unread signals;
- independently reviewed signals may still proceed.

### B2 — under-prioritized

Example:

```text
underprioritized:cluster-integrity
```

Expected effect:

- the alert exists;
- the alert is reviewed;
- the observation is not classified as threatening a sufficiently critical security outcome;
- no escalation occurs;
- the Reserve therefore does not respond.

The implementation only needs to make B1 and B2 representably distinct. They do not both need elaborate standalone scenarios.

---

## Case C — no shared failure

No relevant sensing or review substrate is compromised.

Expected effect:

- the larger correlated Watchline may perform at least as well as the smaller independent one;
- relevant observations are reviewed;
- critical threats are correctly classified and escalated;
- independence should not be treated as automatically superior.

This case guards against baking the desired result into the comparison.

---

# Pre-registration

Before running the experiment, write down:

1. which controls should detect in each case;
2. which signals should be surfaced;
3. which signals should be reviewed;
4. which reviewed signals should be classified as critical;
5. which signals should be escalated;
6. whether the Reserve should respond;
7. the expected terminal security state;
8. which Watchline should perform better.

Do this **before examining the output**.

The prediction should be capable of being wrong.

---

# Falsification criteria

## COLLAPSE Muster if:

The result remains fully explainable as:

```text
the earliest surviving blocker wins
```

or:

```text
the configuration with one surviving detector wins
```

without any meaningful interaction between:

- evidence independence;
- review independence;
- criticality assignment;
- escalation;
- response.

Also collapse if the same analysis is materially clearer as ordinary attack-graph / cut-set analysis plus a small script.

In that case, Muster should become a very thin evaluator, potentially over an existing incident representation such as MITRE Attack Flow.

---

## KEEP TINY if:

The experiment produces a clear composition effect such as:

- a Watchline with more detectors loses because they share a compromised evidence substrate;
- independent evidence sources preserve observation;
- signals exist but shared review failure prevents criticality assessment or escalation;
- signals are reviewed but under-prioritized and therefore never escalate;
- independently reviewed and correctly prioritized signals still trigger response;
- changing which sensing or review substrate fails changes the ranking between Watchlines;
- the result cannot be summarized solely by earliest blocker position.

This earns another gate.

It does **not** earn unrestricted feature development.

---

# EXPAND is not an outcome of this gate

Even a successful result is only one micro.

Expansion requires later evidence that the abstraction:

- survives another independently encoded incident;
- remains useful when controls are defined without reference to one scenario;
- exposes interactions not readily captured by ordinary attack graphs / attack-defense trees;
- helps a reader or analyst reason better than an equivalent prose or notebook treatment.

---

# Explicit non-claims

Gate 2 does **not** attempt to show:

- that V/P/R is a novel security taxonomy;
- that `Role` currently has executable semantics;
- that the original Hugging Face incident was primarily a detection-vs-prevention failure;
- that Muster predicts adaptive attacker behavior;
- that exact event-ID controls model realistic detection precision;
- that Muster is a general security simulator;
- that human reviewers behave deterministically in reality;
- that criticality can be reduced to a universal numeric score;
- that this gate provides a general incident-severity framework.

For this gate, V/P/R remains an analytical vocabulary.

The narrower hypothesis is:

> **Does modeling observation-to-action paths — including sensing, review, criticality assignment, escalation, and response — as stateful, failure-prone, and potentially independent reveal a meaningful defensive composition effect that simple blocker placement does not?**

---

# Role of V/P/R in this gate

Do not force `Role` to become load-bearing merely to satisfy prior criticism.

Instead:

- `Vedette` names the observation role;
- `Picket` names interception/prevention;
- `Reserve` names response/recovery;
- executable behavior continues to live in explicit actions and state transitions.

If Gate 2 succeeds, revisit whether `Role` deserves first-class executable semantics.

If it does not, `Role` may remain documentation or be removed.

---

# Gate stop condition

Stop after:

```text
Case A: evidence failure
Case B: review / prioritization failure
Case C: no shared failure
```

against the two pre-registered Watchlines.

Do not add a feature because one case exposes something Muster cannot represent.

At the end of the gate, choose only:

```text
COLLAPSE
```

or:

```text
KEEP TINY
```

The purpose of Gate 2 is to discover whether Muster has a non-trivial experimental reason to exist.

Not to prove that it does.
