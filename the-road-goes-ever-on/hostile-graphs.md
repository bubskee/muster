# Hostile Graphs

> Companion notebook to [`grap-theory.md`](./grap-theory.md). Grap Theory began by asking how Muster maps onto general state-transition formalisms, using Petri nets as the first comparison. This notebook moves to **security-native attack and defense graph formalisms** and asks a sharper question: do they already make Muster's useful semantics natural, redundant, or awkward?

Append-only. Corrections and reversals get added later rather than silently rewriting earlier conclusions.

---

## 2026-09-14 — Starting question

Petri nets established that Muster's stateful behavior is readily expressible in a general state-transition formalism. The remaining question is less about expressive power and more about fit.

Attack graphs, attack trees, attack-defense trees, and related security formalisms are closer to Muster's home territory. The comparison therefore asks:

> Given a concrete incident, proposed controls, and defender state, how naturally can a security-native graph formalism answer what attacker capability remains after the counterfactual?

Things to watch:

- representation of a concrete historical incident rather than an abstract attack space
- ordering and state carried across events
- residual attacker capability rather than only goal reachability
- defensive controls that observe, interrupt, or respond
- persistent defender history
- non-monotone interactions such as JACKPOT
- how much machinery or annotation is required to reproduce Muster-style replay
- whether existing formalisms already provide a cleaner or more general account of the same ideas

The goal is not to defend Muster. The goal is to find the strongest nearby formalism and see what survives the comparison.

## 2026-09-14 — Attack graphs: two close ancestors

The first two papers in this branch already look much closer to Muster than the Petri-net comparison did.

### Sheyner et al. (2002): scenario attack graphs

Sheyner et al. model the network as a finite-state machine. Network configuration and attacker privileges are encoded in state; attacker actions are transitions; a security property defines the bad states. Model checking is then used to enumerate attack scenarios that can reach those states.

The resulting attack graph is therefore fundamentally a **state-space object**:

```text
system model
    ->
possible attacker transitions
    ->
reachable bad states / attack scenarios
```

This is a clear ancestor of Muster's event/precondition/effect machinery.

The important difference is direction of analysis.

Sheyner asks roughly:

> Given this system model, what attack sequences can reach the security goal?

Muster asks:

> Given this concrete incident sequence, what changes if we alter the defense and replay it?

So far, the best distinction is not expressive power. It is the choice of **unit of analysis**: attack-space exploration versus fixed-trace defensive counterfactual replay.

This also suggests a useful formalization route for Watchlines. We can likely reuse a state-transition model rather than inventing a bespoke mathematical foundation.

A natural extension is to split state into attacker/system state and defender state:

```text
x = (A, D)
```

where `A` contains attacker-visible capability and system facts, while `D` contains defender observations, memory, and response state.

That makes the defender-history result ordinary state dependence: two states can agree on `A` while differing on `D`, producing different later reachability.

---

### Ou, Boyer & McQueen (2006): logical attack graphs

Ou et al. explicitly respond to the state explosion in scenario attack graphs.

Instead of making each graph node a complete network state, they make nodes **logical statements** about configuration or attacker privilege. Edges encode causal dependencies between those facts. Their own distinction is useful:

```text
scenario attack graph:  how the attack can happen
logical attack graph:   why the attack can happen
```

The logical graph is bipartite:

```text
fact
  ->
derivation
  ->
fact
```

A derivation node acts like an AND node: all prerequisite facts are required. A derived fact acts like an OR node: several derivations may establish the same fact.

The mapping to Muster is immediate:

```text
Muster requires      ~ prerequisite fact nodes
Muster event         ~ derivation node / interaction rule
Muster adds          ~ derived fact
```

MulVAL expresses these rules in Datalog. Primitive facts come from configuration input; derived facts represent capabilities that follow by repeatedly applying interaction rules.

This is much closer to Muster's dependency structure than a raw state-space graph.

But the two systems still appear to privilege different questions.

MulVAL derives:

> What attacker capabilities follow from this configuration and rule set?

Muster replays:

> Given that these events happened in this order, which still happen under this defensive counterfactual, and what remains true afterward?

---

### A terminology trap: "attack simulation trace"

Ou et al. also use the phrase **attack simulation trace**, but it does not mean an incident trace in Muster's sense.

Their trace records successful logical derivations:

```text
because(rule, fact, conjunct)
```

or approximately:

```text
these prerequisite facts
    ->
this rule
    ->
therefore this fact is derivable
```

It is a proof / justification trace produced while evaluating the Datalog program.

Muster's trace is chronological:

```text
event_1
event_2
event_3
...
```

The distinction may matter a great deal:

```text
MulVAL trace: why a capability is derivable
Muster trace: what happened, in this order
```

---

### The monotonicity seam

The most interesting point so far is the attack-graph literature's treatment of **monotonicity**.

Ou et al. discuss prior work that assumes attacker privilege only increases: launching an attack does not require giving up capability already obtained. They then argue that even many apparently non-monotonic attacks can be treated as monotonic at a sufficiently abstract level, because their causes can be represented propositionally as configuration facts.

That abstraction is excellent for scalable vulnerability analysis.

It may also erase exactly the kind of interaction JACKPOT exposes.

JACKPOT depends on:

```text
later attacker action
    ->
evidence produced
    ->
defensive response enabled
    ->
credential revoked
```

An earlier defensive intervention blocks the attacker action, but also suppresses the evidence:

```text
better local prevention
    ->
less defender state
    ->
later response disabled
    ->
more attacker capability survives
```

This suggests a concrete comparison question:

> What extra structure does a logical attack graph need to distinguish current derivability from historical production and observation of a fact?

If the answer is "none; this is already native," that weakens Watchlines.

If the answer is "explicit temporal / defender-state machinery," that may identify a real seam between dependency analysis and incident replay.

---

### Current map of the lineage

A tentative genealogy:

```text
state-space / scenario attack graphs
        |
        v
logical / exploit-dependency attack graphs
        |
        v
abstraction, grouping, hardening, configuration analysis
        |
        ?
        v
incident-oriented defensive counterfactuals
```

We should follow this lineage before claiming the final step is missing.

In particular, we should look for later work by Ou, McQueen, Jajodia, Noel, and adjacent authors on:

- graph abstraction and grouping
- automated hardening / configuration management
- explicit defensive actions
- dynamic or temporal attack graphs
- defender observations and memory
- attack-defense taxonomies
- incident response and SOC workflows

The research question is narrowing:

> Is Muster merely a lightweight workflow over familiar attack-graph semantics, or does fixed-trace replay plus explicit defender state expose a class of defensive counterfactuals that these formalisms do not make convenient?

Too early to answer. Close enough to keep digging.
