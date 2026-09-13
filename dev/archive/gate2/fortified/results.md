# Gate 2 Results

Gate 2 tested whether Muster's stateful replay representation earns its complexity relative to simpler dependency reasoning.

## Summary

```text
2A        REPLAY EARNED
2B        calibration
2C        static composition; no additional replay evidence
JACKPOT   good synthetic non-monotone interaction
```

Overall verdict:

```text
KEEP TINY
```

Muster does not establish that stateful replay is formally necessary or uniquely expressive. It does demonstrate a compact class of cases where persistent defender history is causally load-bearing and where executable replay makes counterfactual defensive interactions explicit.

## Gate 2A — Observation Budget

The Broad/Correlated Watchline produced a history-dependent divergence between the `EARLY` and `BETWEEN` telemetry-failure treatments.

At an aligned checkpoint, both executions had the same:

* reached incident/world events;
* attacker capabilities;
* current telemetry availability;
* defensive composition.

The only relevant difference was persistent defensive state:

```text
BETWEEN: alert:credential-theft
EARLY:   no persisted alert
```

That persisted alert later enabled review and credential revocation in `BETWEEN`, while `EARLY` reached node access.

**Verdict: REPLAY EARNED.**

This result is specific to the tested Watchline. The Narrow/Diverse sensing composition remained robust through ordinary evidence-source independence and does not require defender-history reasoning.

## Gate 2B — Review Budget

Correlated and diverse review paths behaved according to ordinary dependency propagation.

```text
                     NONE      SOC-FAILURE   AGENT-FAILURE

Correlated H1+H2      secure      NODE          secure
Diverse    H1+A1      secure      secure        secure
```

No additional history-dependent effect appeared.

**Verdict: calibration only.**

Running 2B exposed a trace-fidelity bug: a control that was initially `WAITING` could remain recorded that way even after an earlier replay phase changed its prerequisites and it should subsequently have been recorded as `SUPPRESSED`. Execution outcomes were unaffected. Replay phase tracing was corrected and regression-tested.

## Gate 2C — Integrated Composition

The integrated observation × review matrix matched the frozen static-null prediction.

The main architectural lesson was that independence does not compose automatically: the surviving evidence source must connect to a surviving review path.

In particular, the fully diverse architecture could survive one failure domain while becoming the uniquely failing architecture under another.

**Verdict: useful composition result, but no additional evidence for stateful replay.**

## Non-Monotonicity Fixture — Synthetic JACKPOT

The dedicated premature-containment fixture produced a non-monotone defensive interaction.

For nested functioning-control sets:

$$
W_1 \subset W_2
$$

the baseline remained secure, while adding the functioning `L1 → R2` early-containment path caused node access.

Mechanism:

```text
credential stolen
→ early containment removes worker access
→ later discovery event becomes unreachable
→ later evidence is never generated
→ stronger response never triggers
→ stolen portable credential survives
→ credential replay succeeds
→ node access
```

Component ablations showed that neither `L1` nor `R2` alone produced the adverse outcome; the interaction required both.

Control declaration order did not affect the result.

**Verdict: good synthetic result / synthetic JACKPOT.**

This was a preregistered constructed mechanism, not a discovery about the real Hugging Face incident.

## Ablations

### Memoryless defender state

Replacing historical-alert consumption with a memoryless reviewer caused the `EARLY` and `BETWEEN` 2A outcomes to converge to the same adverse result.

This supports the claim that persistent defender history is causally load-bearing.

### Pipeline collapse

Collapsing:

```text
review
→ criticality
→ escalation
```

into the smallest equivalent response gate preserved both:

* the 2A history-dependent outcome;
* the synthetic non-monotonicity result.

Therefore these intermediate stages are useful explanatory / trace vocabulary in the current experiments, but are not themselves load-bearing dynamics.

The smaller causal core is approximately:

```text
persistent defensive state
+ incident reachability
+ attacker capability state
+ response effects
→ terminal outcome
```

## Interpretation

Gate 2 supports keeping Muster small.

The strongest evidence is not that Muster expresses something impossible to encode in another formalism. Rather:

1. persistent defender history materially changes a future outcome in 2A;
2. a deliberately constructed defensive interaction becomes non-monotone under replay;
3. simpler cases correctly collapse to ordinary dependency reasoning;
4. ablations identify which parts of the model are genuinely causal and which are primarily explanatory structure.

That is sufficient for Muster to remain a compact executable counterfactual harness rather than collapsing entirely into a static table or dependency diagram.
