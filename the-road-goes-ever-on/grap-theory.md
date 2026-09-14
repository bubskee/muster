# Grap Theory

> *Under mild conditions, graps exhibit substantial graphness.*

An append-only notebook for comparing **Muster / Watchlines** with adjacent formalisms.

The rule for this file is simple: **do not rewrite history**. Add corrections, refinements, and reversals as later entries. If an earlier claim turns out to be wrong, leave it in place and append what changed.

---

## 2026-09-14 — Starting point

### Why this exists

Muster began as a lightweight replay engine for a practical question:

> Given a concrete incident trace and a proposed defense, what attacker capability remains after replay?

Its representation was intentionally conventional. The interesting part was not supposed to be a new graph formalism. The interesting part was making defensive counterfactuals **executable, inspectable, and easy to run against a known incident**.

After the sprint, the obvious question became:

> How much of Muster is merely a convenient presentation of ideas already native to existing state-transition formalisms?

This notebook records the answer as we explore it.

The name **Grap Theory** is, regrettably, sticking.

---

## 2026-09-14 — Minimal Muster semantics

For present purposes, a Muster incident can be approximated as:

```text
initial state
+ ordered event trace
+ event preconditions
+ event effects
+ controls / defensive state
--------------------------------
counterfactual replay
```

A simplified event looks like:

```text
Event:
    requires: facts
    adds: facts
    removes: facts
```

Replay proceeds left-to-right through a fixed incident trace.

For each event:

```text
if preconditions are satisfied:
    apply its effects
else:
    skip it

continue to the next historical event
```

The resulting state tells us what attacker capability remains.

Two details already appear important:

1. **The incident order is supplied.** Muster is not primarily searching over every possible execution.
2. **Skipping is part of the execution semantics.** A failed historical event does not halt replay; replay continues to later events.

This is likely to matter when comparing Muster with formalisms whose default question is reachability over all enabled transitions.

---

## 2026-09-14 — Watchlines as state effects

Watchlines currently use three defensive roles:

```text
Vedette  — see: something is happening
Picket   — interrupt: not through here
Reserve  — contain: we were wrong; limit the damage
```

A rough state-transition interpretation is:

```text
Attacker state = A
Defender state = D
```

### Vedette

Primarily changes defender state:

```text
(A, D) -> (A, D')
```

The attacker may be in the same visible position, but the defender now knows or remembers something that can matter later.

This maps naturally onto the "defender history matters" result from Muster.

### Picket

Prevents or alters an attacker transition:

```text
A --t--> A'
```

becomes blocked, redirected, or otherwise unavailable.

### Reserve

Acts after some compromise or uncertainty has already occurred, changing later reachable state:

```text
(A, D) -> (A', D')
```

The distinctions are not yet a formal taxonomy. One research question is whether they correspond to useful, non-arbitrary classes of transition in a more explicit formal model.

---

## 2026-09-14 — Petri nets: first mapping

A basic Petri-net mapping is straightforward:

```text
Muster fact              <-> Petri place
fact currently true      <-> token in that place
Muster event             <-> transition
event precondition       <-> input/read arc
event effect             <-> output/consuming arc
Muster state             <-> marking
execute event            <-> fire transition
```

This immediately shows substantial overlap.

However, ordinary Petri nets naturally ask something closer to:

> From this marking, which transitions are enabled, and what markings are reachable?

Muster naturally asks:

> Here is the next event in the historical trace. Does it still happen under this defense? Then continue to the next historical event.

That difference in default problem shape became the main thing to investigate.

---

## 2026-09-14 — First Petri encoding: too much plumbing

Our first JACKPOT encoding used ordinary places and transitions.

We initially represented boolean negation with explicit complement places such as:

```text
no_evidence
no_worker_isolated
no_credential_revoked
no_node_access
```

This kept the state space finite, but it made the encoding look much clumsier than Muster.

That comparison was unfair.

SNAKES supports richer arc annotations, including:

```text
Test(...)       — require a token without consuming it
Inhibitor(...)  — require that a matching token be absent
```

Once we used those idiomatically, most of the `no_*` scaffolding disappeared.

**Correction:** absence-testing and read-only preconditions are not a meaningful strike against the Petri-net formulation when using SNAKES.

---

## 2026-09-14 — Fixed-order replay with phase places

To reproduce JACKPOT exactly, we next introduced explicit incident phases:

```text
phase_observe
phase_response
phase_replay
done
```

The baseline path became:

```text
discover
-> revoke
-> replay blocked
```

The early-isolation counterfactual became:

```text
isolate early
-> revocation skipped
-> replay succeeds
```

Both produced four reachable states and exactly reproduced the reversal:

```text
BASELINE
credential revoked
no node access
```

versus:

```text
EARLY ISOLATION
credential survives
node access obtained
```

This worked, but the phase places were clearly encoding something Muster already has natively: **the supplied order of the incident trace**.

---

## 2026-09-14 — Colored token as program counter

We then replaced four phase places with one colored token:

```text
pc = 0
pc = 1
pc = 2
pc = 3
```

Transitions explicitly advanced the token:

```text
discover : 0 -> 1
revoke   : 1 -> 2
replay   : 2 -> 3
```

This cleaned up the representation, but did not remove the underlying requirement.

The historical order was still being imposed on the Petri net.

The encoding had become cleaner, not more native.

---

## 2026-09-14 — The trace itself as a colored token

The next step was more revealing.

Instead of a numeric program counter, we placed the remaining incident trace directly in a colored token:

```text
("discover", "revoke", "replay")
        ->
("revoke", "replay")
        ->
("replay",)
        ->
()
```

For the counterfactual:

```text
("isolate_early", "revoke", "replay")
```

The security facts remained Petri places. The historical incident became data carried by the net.

This representation is much closer to the conceptual structure of Muster:

```text
security state + remaining incident trace
```

The result still reproduced JACKPOT exactly.

At this point the Petri model no longer looked like "a Petri net that happens to resemble Muster."

It looked increasingly like **an interpreter for a Muster trace implemented using Petri-net machinery**.

---

## 2026-09-14 — Compiling Muster-shaped events into a Petri net

We then stopped handwriting individual trace transitions and described events declaratively:

```text
Event:
    name
    requires
    adds
    removes
```

A compiler converts each event occurrence into Petri transitions.

For an event whose requirements are met:

```text
apply:event
```

fires and advances the trace token.

For unmet preconditions, the compiler emits skip transitions such as:

```text
skip:event:missing:fact
```

which advance the historical trace without applying the event effects.

This directly implements Muster's core rule:

```text
if prerequisites met:
    execute event
else:
    skip event
continue
```

The compiled Petri net again reproduced the same baseline and early-isolation results.

So we now have:

```text
Muster-shaped event list
        |
        v
Petri-net compiler
        |
        v
colored Petri net
        |
        v
reachability graph
        |
        v
same JACKPOT reversal
```

This is the strongest result of the exploration so far.

---

## 2026-09-14 — Current interpretation

Petri nets can plainly express the relevant Muster semantics.

The interesting difference is therefore **not expressive possibility**.

The emerging distinction is about what each representation treats as primitive.

### Petri nets make primitive

```text
state / marking
enabled transitions
firing
reachability
concurrency / alternative execution order
```

### Muster makes primitive

```text
a concrete historical trace
ordered replay
event preconditions and effects
skip-and-continue behavior
residual attacker capability
defender history carried through replay
```

To reproduce Muster's workflow in a Petri net, we added or compiled:

```text
ordered trace as data
apply-vs-skip semantics
event preconditions/effects
state carried across the trace
```

None of these are beyond Petri nets.

But they are not the default question Petri nets appear designed to make convenient.

A tentative characterization:

> **Muster is a deliberately restricted state-transition system whose defaults are shaped around counterfactual replay of a known incident rather than exploration of all possible executions.**

This is a hypothesis, not yet a conclusion.

---

## 2026-09-14 — JACKPOT and non-monotonic defense

The synthetic JACKPOT example has survived translation into Petri-net semantics.

The mechanism is:

```text
later attacker action
    ->
evidence generated
    ->
credential revoked
    ->
later replay blocked
```

Adding an earlier Picket removes the attacker action:

```text
early isolation
    -X->
discovery
```

but therefore also removes:

```text
discovery
-> evidence
-> revocation
```

so the stolen credential survives and later replay succeeds.

Abstractly:

```text
Picket removes transition t
t produces defender evidence e
Reserve r requires e
therefore removing t disables r
and may enlarge later attacker reachability
```

The important claim is only an **existence result**:

> A locally effective defensive intervention can worsen later modeled attacker capability by suppressing state that another defense depends on.

The Petri-net translation strengthens confidence that this is not an artifact of Muster's implementation.

It does **not** establish prevalence or practical frequency.

---

## 2026-09-14 — What we have not established

We have not shown that Muster introduces a new formalism.

We have not shown that Petri nets are unable to represent any Muster behavior.

We have not shown that Muster is more expressive.

We have not shown that a Petri-net implementation is inherently harder to use.

We have not yet compared against the strongest security-native alternatives.

We have not yet established that Muster produces materially important real-world insight unavailable from existing attack-defense or state-transition analysis.

The current result is narrower and more useful:

> Muster's semantics can be compiled into a colored Petri-net model, but doing so makes explicit several assumptions that Muster treats as first-class execution semantics.

---

## 2026-09-14 — Questions exposed by the Petri experiment

The useful next questions are:

- Can the compiler be made generic enough to consume actual Muster scenario definitions?
- How much Petri-net structure is generated per Muster event?
- Does the Petri version remain inspectable as incident size increases?
- Are skip transitions merely implementation plumbing, or is there a more idiomatic Petri representation?
- Can the same defender-history experiments be encoded without adding special machinery?
- Can Watchline roles be characterized directly in terms of transition effects?
- Does a conventional Petri reachability query answer the same practical questions as Muster's residual-capability output?
- Is the chief difference representational convenience, workflow, or something deeper?

---

## 2026-09-14 — Formalisms still to visit

### Attack / attack-defense graphs and trees

Likely the most important security-native comparison.

Questions:

```text
How naturally do they encode a concrete incident trace?
How naturally do they encode residual attacker capability?
How naturally do they encode defender history?
How naturally do they encode the JACKPOT reversal?
```

### Alloy

Potentially useful for small-model exploration and counterexample search.

Especially attractive for questions like:

```text
Does adding a Picket always reduce attacker reachability?
```

If false, Alloy should be able to find a small counterexample.

### TLA+

Potentially useful for explicit state-transition semantics, temporal ordering, and invariants.

Could be a good fit if the question becomes:

```text
What properties of Watchline execution should hold over all traces?
```

### Lean / Coq / Agda

Only worth reaching for if we discover properties worth machine-checking.

The language itself is not the contribution; the target would need to be explicit theorems such as determinism, equivalence, or monotonicity under restricted classes of controls.

### Haskell

Potentially attractive as a small reference implementation because pure functions and algebraic data types map cleanly onto the semantics.

But:

> Haskell is not proof merely by being Haskell.

Its value would be clarity and property testing, not academic gravitas.

---

## 2026-09-14 — Design lesson for Muster

If redesigning from scratch, the likely architecture would not be "rewrite Muster in a more formal language."

A more compelling separation is:

```text
Watchline / replay semantics
        |
        | explicit mathematical or formal definition
        v
reference model / model checker
        |
        +----------------+
        |                |
        v                v
formal exploration     Muster
                      (Go)
                        |
                        v
                 executable tool
```

Go remains attractive for Muster-the-tool:

```text
simple state
explicit code
easy CLI
easy distribution
good tests
low ceremony
```

The formal semantics can live elsewhere.

That keeps Muster executable without forcing the implementation language to carry the burden of proof.

---

## 2026-09-14 — Grap Theory, theorem zero

**Theorem 0 (informal).**

> Under mild conditions, graps exhibit substantial graphness.

**Evidence.**

A Muster-shaped event trace was compiled into a colored Petri net and reproduced the same JACKPOT counterfactual reversal.

**Caveat.**

The compiler carried enough Muster semantics into the Petri representation that the result should not be mistaken for a claim of novelty.

**Corollary.**

> Every sufficiently disciplined grap can be embedded in a graph, provided one is willing to carry around enough snakes.

Peer review pending.

## 2026-09-14 — Defender history in Petri nets

The defender-history result required no special Petri-net machinery.

Two runs were constructed with identical attacker-visible state and identical remaining incident trace after an observation event:

```text
attacker: credential_held
trace: ("review", "replay")
```

They differed only in defender state:

```text
persistent history:    alert_seen
current-state-only:    -
```

The resulting reachability diverged. With `alert_seen`, review revoked the credential and later replay failed. Without persistent defender state, review was skipped and replay reached `node_access`.

This suggests that the substantive phenomenon is ordinary state dependence:

> Two executions can agree on attacker-visible state while differing in defender-history state, producing different later reachability.

Petri nets represent this directly as different markings.

The distinction therefore does not support a claim that Muster introduces a new state concept. Muster's contribution, if any, is in making defender history explicit and convenient within a fixed incident-replay workflow.
