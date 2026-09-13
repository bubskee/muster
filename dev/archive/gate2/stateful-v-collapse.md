# Muster Gate 2 — Stateful Replay or Collapse

## Why this gate exists

Gate 1 established that Muster's replay engine behaves coherently.

It did **not** establish that a replay engine is necessary.

The Hugging Face incident remains useful as a motivating and calibration case because many of its intervention points are straightforward in hindsight. That is desirable: an obvious incident gives Muster somewhere to demonstrate that its representation is sane.

But:

> **HF being avoidable is not evidence that Muster deserves to exist.**

Gate 2 asks the narrower and more hostile question:

> **Does stateful replay expose a defensive composition effect that a simpler causal/dependency model cannot express just as clearly?**

The default result is:

```text
COLLAPSE
```

Muster must earn:

```text
KEEP TINY
```

`EXPAND` is not an available outcome.

---

# 1. The competitor Muster must beat

Do not compare Muster against a naive attack graph that knows nothing about defensive dependencies.

Give the static alternative every reasonable advantage.

Let the incident be a fixed event DAG:

$$
G=(V,E)
$$

For an attack/world event \(v\):

$$
reachable(v)=
\left(
\bigwedge_{u \in Pred(v)} succeeds(u)
\right)
\land
\neg blocked(v,C)
$$

where \(C\) is the defensive composition.

Controls may have static dependencies:

$$
blocked(v,C)=
\bigvee_{c \in C}
[
c \text{ applies to } v
\land substrateAvailable(c,v)
\land c \text{ blocks } v
]
$$

The baseline may also know that an earlier reached incident event disabled a substrate:

$$
substrateAvailable(c,v)
=
\neg \exists u \in Ancestors(v):
succeeds(u)
\land
u \text{ disables substrate}(c)
$$

The baseline therefore gets:

```text
event reachability
causal prerequisites
control coverage
control substrate dependencies
common-mode substrate failure
ordinary prevention
downstream causal pruning
```

It may conclude things such as:

```text
credential acquisition blocked
→ later credential use unreachable
```

without requiring replay.

This is not considered stateful merely because events are evaluated in causal order.

---

# 2. The boundary: defender history

The static baseline does **not** receive mutable defensive history beyond what can be reconstructed from reached incident nodes.

The decisive question is:

> Can the future at event \(v\) be determined solely from the incident graph, the defensive composition, the incident nodes reached so far, current substrate availability, and attacker capabilities?

If yes:

```text
STATIC MODEL SUFFICIENT
```

If another persistent variable is required because of **how the execution arrived there**, replay may be justified.

Examples include:

```text
alert already emitted
alert persisted after detector loss
signal already reviewed
partial response already executed
credential already revoked
resource already quarantined
one-shot response already consumed
earlier observation changed later interpretation
```

The strongest diagnostic is **history equivalence**.

Suppose two executions reach the same incident point with:

```text
the same reached attack/world nodes
the same currently available substrates
the same defensive composition
the same attacker capabilities
```

If their futures must be identical, the static representation is sufficient.

If their futures can differ because one execution contains a persistent defensive fact created earlier, Muster has crossed into genuinely stateful replay.

Formally, the static model aims for:

$$
Outcome = F(G,C)
$$

or, with an explicitly varied ordering:

$$
Outcome = F(G_{\tau},C)
$$

Muster earns its state machine only if the useful analysis requires something like:

$$
S_{t+1}=T(S_t,E_t,C)
$$

and:

$$
Outcome=H(S_n)
$$

where some relevant component of \(S_t\) cannot be reconstructed merely from the set of incident nodes reached so far.

Otherwise:

> `State` is graph reachability wearing a struct.

---

# 3. Shared failure domains are calibration, not a result

Gate 2 must not be passed by:

```text
detector depends on worker telemetry
worker telemetry fails
detector fails
```

Nor by:

```text
three detectors share one substrate
shared substrate fails
all three detectors fail
```

Those outcomes are expected from the declarations themselves.

They establish that Muster can represent common-mode failure.

They do not establish analytical value.

Therefore Cases A–C remain useful, but are reclassified as **calibration fixtures**.

---

# 4. Calibration fixtures

## A — evidence failure

Represent:

```text
matched
but
evidence substrate unavailable
```

as distinct from:

```text
unmatched
```

Expected dispositions:

```text
UNMATCHED
WAITING
SUPPRESSED
READY
```

must remain semantically distinct.

This tests the engine.

It does not help Muster pass Gate 2.

---

## B1 — review failure

Represent:

```text
alert exists
review path unavailable
alert remains unread
no escalation
```

This demonstrates that sensing and consumption are distinct.

Again, representational value only.

---

## B2 — criticality failure

Represent:

```text
alert exists
alert reviewed
criticality insufficient
no escalation
```

This is valuable taxonomy.

It becomes analytically load-bearing only if criticality depends on accumulated state rather than merely acting as another Boolean gate attached to the same event.

A preferred Gate 2 design makes criticality depend on evidence accumulated across more than one event, for example:

```text
alert:worker-exec
AND
alert:credential-access
→ critical:cluster-integrity
```

If criticality can be deleted or collapsed into escalation without changing any important outcome, treat it as presentation rather than engine semantics.

---

## C — healthy path

No relevant substrate fails.

All controls behave sensibly.

This prevents the experiment from defining independence as automatically superior.

It still cannot earn `KEEP TINY`.

---

# 5. Experimental construction order

To prevent result-by-construction, freeze things in this order.

## Step 1 — freeze the incident

Define:

```text
incident events
event prerequisites
event effects
initial state
designated adverse terminal outcome
```

before defining competing Watchlines.

If the micro is synthetic, label it synthetic.

Do not construct the incident after deciding which Watchline should win.

---

## Step 2 — define controls independently

Define each control without reference to the composition in which it will later appear.

For each control freeze:

```text
event coverage
role
action
Requires
Effects
substrate dependency
SuppressedBy
```

Where practical for real controls, inherit evidence/log-source dependencies from published detector mappings rather than choosing substrates specifically to manufacture correlation.

Synthetic dependencies must be fixed before composition comparison.

---

## Step 3 — define the allowed compositions

Use a small fixed control pool and a fixed defensive budget.

If the control pool is small enough, enumerate all sensible compositions rather than selecting only a preferred hero and villain.

Named Watchlines may still be used as explanatory examples.

---

## Step 4 — write the null prediction first

Before executing Muster, evaluate the same cases using the static causal-pruning model.

Write down:

```text
reachable incident nodes
disabled controls
surviving paths
expected terminal outcome
predicted composition ranking
```

The static analysis is the null.

Agreement with it is **not validation**.

Agreement everywhere is evidence toward collapse.

---

# 6. Equal-coverage mechanism fixture

Run one deliberately controlled fixture in which competing compositions have:

```text
equal control count
equal event coverage
equal Reserve capability
equal review / criticality / escalation semantics
```

and differ primarily in dependency wiring.

This isolates the mechanism cleanly.

Its purpose is:

> Show that the engine represents the intended dependency interaction without coverage confounds.

This is a **mechanism fixture**, not the decisive experiment.

A successful equal-coverage fixture does not earn `KEEP TINY`.

---

# 7. The decisive experiment: fixed-budget crossover

The real comparison should permit the realistic trade-off between:

```text
coverage
and
independence
```

Hold defensive cost/budget fixed.

One composition may obtain broader sensing by reusing a shared substrate.

Another may obtain less coverage but greater independence.

Neither should weakly dominate the other by construction.

The important output is not:

```text
independence wins
```

but a conditional decision boundary such as:

```text
under condition X: A > B
under condition Y: B > A
```

The same controls must retain the same semantics across all runs.

---

# 8. Failure timing must matter for more than prefix length

At least one Gate 2 case must vary **when** a defensive substrate becomes unavailable.

Minimum timing set:

```text
substrate never fails
substrate fails before first relevant opportunity
substrate fails between two relevant opportunities
```

However:

> A ranking reversal caused merely by “which detector occurs before position k” does not pass the gate.

In a linear chain, failure position can collapse into prefix length.

That is still static reasoning.

A valid temporal result must involve persistent state or another interaction that is not reducible to:

```text
earliest surviving detector
```

Prefer a micro in which:

```text
an earlier defensive artifact persists
the source that produced it later disappears
a later control can still act on the artifact
```

or another equivalent form of path-dependent state.

---

# 9. History-equivalence challenge

Construct at least one paired case in which, at some checkpoint, two executions have:

```text
same reached incident nodes
same current substrate availability
same control composition
same attacker capability
```

but different defensive history.

Example shape:

```text
History A

suspicious operation
→ alert emitted and persisted
→ telemetry destroyed
→ checkpoint
→ reviewer later consumes preserved alert
→ response
```

versus:

```text
History B

telemetry destroyed
→ suspicious operation
→ no alert exists
→ checkpoint
→ reviewer has nothing to consume
→ no response
```

If the future differs because of the latent alert state, replay is doing something the defined static baseline does not represent.

If all such cases can be reduced cleanly to additional incident nodes without effectively encoding defender-history state into the graph, the static representation wins.

---

# 10. Non-monotonicity challenge

Gate 2 should also test whether the defensive system is **coherent** in the ordinary sense that more functioning controls cannot worsen the designated security outcome.

Let \(W\) be the set of functioning/non-suppressed controls.

For two configurations:

$$
W_1 \subseteq W_2
$$

a monotone defensive system should never produce:

```text
W1: acceptable terminal outcome
W2: adverse terminal outcome
```

simply because the additional control in \(W_2\) functioned.

Define the designated adverse terminal predicate before running the experiment.

Then exhaustively check the small run set for monotonicity violations.

A particularly strong Gate 2 result is:

> A functioning defensive control causes an earlier or weaker response that prevents later evidence accumulation, producing a worse terminal security outcome than when that control is suppressed.

Illustrative synthetic shape:

```text
early alert
→ weak/partial escalation
→ partial containment
→ later evidence never materializes
→ stronger response never triggers
→ important attacker capability survives
```

while:

```text
early detector suppressed
→ attack proceeds far enough to expose second signal
→ combined evidence becomes critical
→ full response
→ important attacker capability removed
```

Thus:

```text
working control
→ worse terminal state
```

This is not required because “perverse defenses are interesting.”

It is useful because it breaks the monotonicity assumed by the simplest cut-set/reliability treatment.

### Harsh pass rule

For this gate:

> If no case in the frozen run set exhibits either a history-equivalence violation or a defensible non-monotone defensive interaction, assume the replay engine has not earned its complexity.

If all substantive conclusions remain monotone causal pruning:

```text
COLLAPSE
```

---

# 11. Memoryless ablation

Run the decisive cases through an intentionally weaker evaluator with no persistent defender history.

It may use:

```text
current event
fixed composition
current substrate availability
incident-DAG reachability
```

It may not retain:

```text
earlier alerts
earlier reviews
earlier escalations
pending signals
response-created latent defensive state
one-shot resource consumption
other facts whose meaning depends on history
```

Compare:

```text
terminal outcome
composition ranking
substantive explanation
```

If the memoryless evaluator reproduces what matters:

```text
COLLAPSE
```

Do not respond by adding richer state.

---

# 12. Pipeline-collapse ablation

Run the same cases with:

```text
review
→ criticality
→ escalation
```

collapsed into the smallest equivalent gate.

If this produces the same meaningful terminal outcomes and composition rankings:

```text
review / criticality / escalation
```

are explanatory labels, not load-bearing dynamics.

Keep them only if useful for trace readability.

Do not cite them as evidence that replay is necessary.

Likewise, if the descriptive `Substrate` field changes nothing beyond what `SuppressedBy` already encodes, treat it as documentation rather than executable semantics.

---

# 13. Freeze rule

Before examining decisive outputs, freeze:

```text
incident micro
initial state
control pool
control coverage
control effects
substrate dependencies
SuppressedBy conditions
criticality logic
Reserve effects
allowed compositions
failure positions
adverse terminal predicate
static-null predictions
```

After seeing the results:

> Weak results are not permission to add features.

Bug fixes are allowed.

They must be recorded.

If the experiment becomes interesting only after adding a new capability, that is evidence against the current abstraction.

---

# 14. Anti-cherry-picking rule

Prefer enumeration over hand selection.

For a sufficiently small control pool:

```text
enumerate every allowed composition
×
every preregistered substrate condition / failure position
```

Then inspect:

```text
ranking
terminal outcome
history dependence
monotonicity
```

Do not search through scenarios until a favorite Watchline wins.

The incident remains fixed.

---

# 15. COLLAPSE criteria

Choose:

```text
COLLAPSE
```

if any of the following describes the substantive result.

### Static dependency analysis is sufficient

Everything reduces to:

```text
control depends on S
S fails
control fails
```

or:

```text
control blocks E
descendants of E become unreachable
```

---

### Timing is merely prefix length

Changing failure position only changes which detector happens to occur first.

---

### Defender history is reconstructible from attack reachability

No latent defensive fact independently affects the future.

---

### Memoryless ablation preserves the result

The same composition ranking and substantive explanation survive without defender history.

---

### All useful defensive interactions remain monotone

Additional functioning controls never make the designated security outcome worse, and the resulting behavior is adequately represented by causal pruning / cut-set analysis.

---

### Pipeline stages are decorative

Review, criticality, and escalation can be collapsed without changing anything important.

---

### The result depends on authored failure assignments

The experiment succeeds only because substrates or failures were selected specifically to disable one composition.

---

### A simpler representation communicates the result better

Attack Flow, an attack-defense tree, a fault-tree/cut-set representation, or a small dependency script answers the same question more clearly.

In that case Muster should collapse into the thinner representation.

That is a successful research result.

---

# 16. KEEP TINY criteria

Choose:

```text
KEEP TINY
```

only if the frozen experiment demonstrates a compact class of cases in which:

1. The static causal-pruning baseline is written down first.
2. The incident and controls were frozen before results were known.
3. The result depends on ordering plus persistent defensive state.
4. At least one important defensive fact cannot be reconstructed merely from reached incident nodes and current substrate availability.
5. The memoryless ablation loses an important terminal distinction or gets a composition ranking wrong.
6. The phenomenon survives enumeration or modest nearby perturbations rather than appearing only in one hand-selected pairing.
7. The engine requires no major feature expansion to produce the result.
8. The resulting explanation is clearer in replay than in the equivalent static encoding.

A non-monotone control interaction is especially strong evidence.

If no such interaction appears, the history-equivalence test must independently demonstrate why persistent replay state is load-bearing.

Even then:

```text
KEEP TINY
```

means only:

> Muster survived one serious falsification attempt.

It earns another gate.

Nothing more.

---

# 17. What Gate 2 does not establish

Gate 2 does not establish:

```text
that V/P/R is a novel security taxonomy
that Muster predicts attacker behavior
that deterministic review is realistic
that Hugging Face required a replay engine to understand
that shared failure domains are novel
that Muster is better than dynamic fault trees
that Muster is better than Petri nets
that one synthetic micro demonstrates reuse
```

More expressive formalisms can represent order and state.

Muster's possible advantage over them is currently only:

```text
smallness
incident-native vocabulary
counterfactual readability
```

Those are ergonomic claims.

They require later evidence.

---

# 18. Gate 2 stop condition

Stop after:

```text
calibration fixtures
equal-coverage mechanism fixture
static-null analysis
fixed-budget crossover
timing/history challenge
memoryless ablation
pipeline-collapse ablation
monotonicity check
```

Do not add YAML.

Do not add stochastic timing.

Do not add queues.

Do not add adaptive attackers.

Do not add a generic graph engine.

Do not add a SOC simulator.

Do not add another incident to rescue this one.

At the end, answer:

```text
1. What did the static model predict?

2. What did replay reveal that static causal pruning did not?

3. Which persistent fact made history load-bearing?

4. Did any additional functioning control make security worse?

5. Did the memoryless ablation preserve the result?

6. Would we still build Muster if this were the only result we ever obtained?
```

Then choose exactly one:

```text
COLLAPSE
```

or:

```text
KEEP TINY
```

The purpose of Gate 2 is not to prove that Muster works.

The purpose is to determine whether Muster is doing work.
