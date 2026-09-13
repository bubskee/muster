# Gate 2 Frozen Incident — Worker-to-Cluster Credential Replay

## Status

**FROZEN BEFORE WATCHLINE DESIGN**

This is a synthetic micro derived from the shape of the July 2026 Hugging Face intrusion.

It is **not** claimed to reproduce the exact Hugging Face timeline.

The incident/world events are frozen before defining competing Watchlines.

Controls may not change these event semantics.

---

## Purpose

The incident must support testing whether persistent defensive history and ordering matter beyond ordinary causal pruning.

It therefore contains:

* a worker foothold;
* theft of a portable cluster credential;
* a worker-local Kubernetes discovery opportunity;
* later reuse of the stolen credential that does not require continued worker access;
* a path to privileged workload creation and node access;
* one explicitly synthetic telemetry-impairment event whose timing can vary without changing attacker capability.

The incident itself contains **no defensive controls, alerts, reviews, criticality assessments, escalations, or responses**.

Those belong to the later Watchline definition.

---

## Initial state

```text id="2mq6ms"
dataset:malicious-on-worker
```

No defensive-substrate failure is present initially.

---

## Core incident events

| ID | Event                                  | Preconditions                                           | Effects                          |
| -- | -------------------------------------- | ------------------------------------------------------- | -------------------------------- |
| E1 | `worker-code-execution`                | `dataset:malicious-on-worker`                           | `+ attacker:worker-access`       |
| E2 | `cluster-credential-theft`             | `attacker:worker-access`                                | `+ attacker:cluster-credential`  |
| E3 | `kubernetes-api-discovery-from-worker` | `attacker:worker-access`, `attacker:cluster-credential` | `+ attacker:cluster-knowledge`   |
| E4 | `stolen-credential-replay`             | `attacker:cluster-credential`                           | `+ attacker:cluster-session`     |
| E5 | `privileged-workload-creation`         | `attacker:cluster-session`                              | `+ workload:privileged-hostpath` |
| E6 | `node-access`                          | `workload:privileged-hostpath`                          | `+ attacker:node-access`         |

---

## Important causal property

E3 is **not** a prerequisite for E4.

That is intentional.

The attacker may perform worker-local discovery while the original foothold remains available, but possession of the copied credential can outlive that foothold.

Therefore:

```text id="ig0wm5"
loss of attacker:worker-access
```

may prevent:

```text id="zzt5rj"
E3 kubernetes-api-discovery-from-worker
```

without necessarily preventing:

```text id="d10212"
E4 stolen-credential-replay
```

provided:

```text id="yf9mkw"
attacker:cluster-credential
```

was already acquired.

This is a property of the frozen incident, not of either Watchline.

---

## Synthetic perturbation event

Define one additional world event:

```text id="m4a49d"
F1 worker-telemetry-impairment

requires:
    attacker:worker-access

effects:
    + compromised:worker-telemetry
```

`F1` changes no attacker capability.

It exists only to alter the availability of defenses that later declare a dependency on `worker-telemetry`.

For Gate 2, **no control may directly observe, block, review, or respond to F1**.

This avoids making detection of the perturbation itself part of the experiment.

---

## Frozen timing arms

The same F1 event has exactly three permitted treatments.

```text id="8csmcp"
NO-FAILURE

E1
E2
E3
E4
E5
E6
```

```text id="zx1mcy"
EARLY-FAILURE

E1
F1
E2
E3
E4
E5
E6
```

```text id="bvy1fi"
BETWEEN-OPPORTUNITIES

E1
E2
F1
E3
E4
E5
E6
```

No other F1 position is part of Gate 2.

The event definitions are identical across all three arms.

Only ordering changes.

---

## History-equivalence checkpoint

The EARLY-FAILURE and BETWEEN-OPPORTUNITIES runs create a deliberate checkpoint immediately after both E2 and F1 have occurred.

At that point, both incident histories have:

```text id="btn36p"
E1 reached
E2 reached
F1 reached

attacker:worker-access
attacker:cluster-credential
compromised:worker-telemetry
```

They therefore have:

```text id="r3g38s"
the same reached incident nodes
the same attacker capabilities
the same current substrate availability
```

The incident itself supplies no reason for their futures to differ.

If a later Watchline produces different futures because E2 generated a persistent defensive artifact before telemetry impairment in one ordering but not the other, that difference is defender history rather than attack-DAG reachability.

This is the intended history-equivalence challenge.

---

## Branch relevant to premature containment

Because E3 requires:

```text id="golzn4"
attacker:worker-access
```

a defensive response after E2 may make E3 unreachable by removing that fact.

However, because E4 requires only:

```text id="5o4e2u"
attacker:cluster-credential
```

such a response does not automatically prevent later credential replay.

This permits — but does not guarantee — the following class of interaction:

```text id="np456b"
early partial response
→ original worker lost
→ E3 does not occur
→ some later defensive evidence may never exist
→ stolen credential may nevertheless remain usable
```

Whether any Watchline actually produces that behavior is not part of the incident definition.

---

## Adverse terminal predicate

The primary Gate 2 adverse outcome is:

```text id="wk025j"
attacker:node-access
```

present in terminal state.

Therefore:

```text id="liq38c"
SECURITY FAILURE
    iff terminal state contains attacker:node-access
```

Other facts may be reported diagnostically, especially:

```text id="a95et8"
attacker:worker-access
attacker:cluster-credential
attacker:cluster-session
attacker:cluster-knowledge
workload:privileged-hostpath
```

but they do not replace the preregistered primary outcome.

---

## Baseline undefended trajectory

With no controls and no F1:

```text id="1ar391"
E1 applied
E2 applied
E3 applied
E4 applied
E5 applied
E6 applied
```

Terminal state includes:

```text id="00re16"
attacker:node-access
```

With F1 inserted at either permitted failure position, F1 does not itself alter attacker capability, so the undefended attack trajectory remains otherwise unchanged.

---

## Provenance boundary

The following incident shape is grounded in the Hugging Face case:

```text id="ywwjlv"
worker execution
→ credential access
→ Kubernetes / cluster interaction
→ stolen-credential reuse
→ privileged workload / deeper cluster access
→ node-level compromise
```

The exact simplification into E1–E6 is ours.

The following element is explicitly synthetic:

```text id="jgsuxu"
F1 worker-telemetry-impairment
```

No claim is made that Hugging Face's relevant detection telemetry was disabled at this point in the real intrusion.

The micro uses this perturbation solely to test Gate 2's state/history hypothesis.

---

## Containment rule

From this point forward, Gate 2 may vary:

```text id="lpbwm7"
Watchline composition
control substrate dependencies
F1 treatment among the three frozen arms
```

It may not vary:

```text id="o6f2pa"
E1–E6 identities
their prerequisites
their effects
the meaning of F1
the three permitted F1 treatments
the adverse terminal predicate
```

If the incident must be changed to obtain an interesting result, record that as evidence against the current Gate 2 experiment rather than silently repairing the micro.
