# Gate 2 Frozen Control Pool

## Status

**FROZEN AFTER INCIDENT, BEFORE WATCHLINE COMPOSITION**

Controls are defined independently of the Watchlines in which they may later appear.

The incident remains the previously frozen E1–E6 + F1 micro.

No control may observe or respond directly to F1.

---

# Observation controls

| ID | Control                   | Event                                     | Substrate          | Effect                     |
| -- | ------------------------- | ----------------------------------------- | ------------------ | -------------------------- |
| V1 | `worker-credential-theft` | E2 `cluster-credential-theft`             | `worker-telemetry` | `+ alert:credential-theft` |
| V2 | `worker-k8s-discovery`    | E3 `kubernetes-api-discovery-from-worker` | `worker-telemetry` | `+ alert:k8s-discovery`    |
| V3 | `k8s-audit-discovery`     | E3 `kubernetes-api-discovery-from-worker` | `k8s-audit`        | `+ alert:k8s-discovery`    |

All three have:

```text
Role:   Vedette
Action: Observe
```

V1 and V2 are suppressed by:

```text
compromised:worker-telemetry
```

V3 is not.

### V1 — worker credential-theft detector

```text
matches:
    E2 cluster-credential-theft

substrate:
    worker-telemetry

suppressed by:
    compromised:worker-telemetry

effect:
    + alert:credential-theft
```

This represents worker-local observation of access to credential material.

The exact worker-telemetry mapping is synthetic but plausible and must not be presented as a specific published Hugging Face control.

### V2 — worker-side Kubernetes discovery detector

```text
matches:
    E3 kubernetes-api-discovery-from-worker

substrate:
    worker-telemetry

suppressed by:
    compromised:worker-telemetry

effect:
    + alert:k8s-discovery
```

This represents host/EDR/process-side observation of Kubernetes discovery activity originating from the worker.

### V3 — Kubernetes audit discovery detector

```text
matches:
    E3 kubernetes-api-discovery-from-worker

substrate:
    k8s-audit

effect:
    + alert:k8s-discovery
```

This represents control-plane observation of Kubernetes resource enumeration.

V2 and V3 therefore observe **the same incident event through different evidence domains**.

That pairing exists specifically so dependency structure can be varied without changing event coverage.

---

# Shared history-aware review controls

## H1 — review persisted credential alert

```text
ID:
    review-credential-history

matches:
    E3 kubernetes-api-discovery-from-worker

requires:
    alert:credential-theft

substrate:
    soc-queue

suppressed by:
    overloaded:soc-queue

action:
    Review

effect:
    + reviewed:cluster-integrity
```

H1 deliberately consumes an alert produced by an **earlier incident event**.

This is the control that makes the history-equivalence fixture possible.

If E2 produced an alert before worker telemetry failed, H1 may still consume that persisted alert at E3.

If telemetry failed before E2, the alert never existed and H1 remains waiting.

---

## H2 — review current Kubernetes discovery alert

```text
ID:
    review-k8s-discovery

matches:
    E3 kubernetes-api-discovery-from-worker

requires:
    alert:k8s-discovery

substrate:
    soc-queue

suppressed by:
    overloaded:soc-queue

action:
    Review

effect:
    + reviewed:cluster-integrity
```

H1 and H2 provide two routes to the same reviewed security concern:

```text
credential-theft history
             \
              → reviewed:cluster-integrity
             /
current k8s discovery
```

This gives us OR semantics without adding OR logic to the engine.

---

# Shared critical-response chain

## C1 — cluster-integrity criticality

```text
ID:
    assess-cluster-integrity

matches:
    E3 kubernetes-api-discovery-from-worker

requires:
    reviewed:cluster-integrity
    attacker:cluster-credential

action:
    Critical

effect:
    + critical:cluster-integrity
```

The credential must still exist for this condition to represent an active cluster-integrity threat.

---

## E1 — critical escalation

```text
ID:
    escalate-cluster-integrity

matches:
    E3 kubernetes-api-discovery-from-worker

requires:
    critical:cluster-integrity

action:
    Escalate

effect:
    + escalated:cluster-integrity
```

---

## R1 — full credential revocation

```text
ID:
    revoke-cluster-credential

Role:
    Reserve

matches:
    E3 kubernetes-api-discovery-from-worker

requires:
    escalated:cluster-integrity

action:
    Respond

effects:
    - attacker:cluster-credential
```

Because E4 requires:

```text
attacker:cluster-credential
```

successful R1 makes E4 and therefore E5/E6 unreachable.

This is ordinary causal pruning and is **not by itself evidence for keeping Muster**.

---

# Optional premature-containment path

These controls are frozen now because they are required for the preregistered non-monotonicity challenge.

They are not automatically present in every Watchline.

## L1 — low-confidence escalation

```text
ID:
    escalate-worker-containment

matches:
    E2 cluster-credential-theft

requires:
    alert:credential-theft

action:
    Escalate

effect:
    + escalated:worker-containment
```

This represents a rapid but limited reaction to the first signal.

---

## R2 — partial worker containment

```text
ID:
    isolate-worker

Role:
    Reserve

matches:
    E2 cluster-credential-theft

requires:
    escalated:worker-containment

action:
    Respond

effects:
    - attacker:worker-access
```

R2 deliberately does **not** remove:

```text
attacker:cluster-credential
```

because E2 has already succeeded before the response phase.

Consequently:

```text
R2
→ E3 becomes unreachable
```

while:

```text
E4 may remain reachable
```

This creates the possibility that early partial containment suppresses a later observation opportunity without removing the already-stolen portable credential.

Whether that actually produces a worse terminal outcome is left to the frozen experiment.

---

# Control-pool invariants

The following semantics may not change after observing Gate 2 results.

| Control | Event coverage | Substrate        | Result                                  |
| ------- | -------------- | ---------------- | --------------------------------------- |
| V1      | E2             | worker telemetry | credential-theft alert                  |
| V2      | E3             | worker telemetry | Kubernetes-discovery alert              |
| V3      | E3             | Kubernetes audit | Kubernetes-discovery alert              |
| H1      | E3             | SOC queue        | review from historical credential alert |
| H2      | E3             | SOC queue        | review from current discovery alert     |
| C1      | E3             | —                | critical cluster-integrity state        |
| E1      | E3             | —                | critical escalation                     |
| R1      | E3             | —                | revoke cluster credential               |
| L1      | E2             | —                | low-confidence escalation               |
| R2      | E2             | —                | remove worker access only               |

The experiment may select subsets of this pool according to a separately frozen composition rule.

It may not alter individual control semantics per composition.

---

# Provenance boundary

The distinction between worker-side observation and Kubernetes control-plane observation is realistic.

MITRE ATT&CK describes Kubernetes resource-discovery detection through Kubernetes API-server audit data, and also notes host/EDR observation of commands such as `kubectl get pods`.

ATT&CK also describes detection of stolen container service-account tokens through Kubernetes API/authentication telemetry.

The exact V1 and V2 controls above are synthetic instances of those broad evidence classes.

The experiment tests the abstraction, not the efficacy of any named commercial detection product.

---

# Explicit exclusions

Gate 2 does not add controls for:

```text
F1 telemetry impairment itself
E5 privileged workload creation
E6 node access
```

No late “save” is added after the decisive response opportunity at E3.

No control receives probabilistic detection.

No alert expires.

No review latency distribution exists.

No queue depth exists.

No adaptive attacker exists.

No new substrate is introduced after results are observed.

---

# Why this pool is sufficient

V2 versus V3 gives an equal-coverage, different-substrate mechanism fixture.

V1 + V2 gives broader temporal/event coverage concentrated in worker telemetry.

V2 + V3 gives narrower unique-event coverage but evidence-source diversity.

H1 makes an alert created before substrate failure persist into a later defensive opportunity.

H2 permits independent later evidence to reach the same downstream response path.

L1 + R2 make it possible for an early functioning detector to cause partial containment that removes the opportunity for E3 while leaving the stolen credential intact.

Therefore the same frozen pool can test:

```text
common-mode failure
history-equivalence
coverage versus independence
memoryless ablation
pipeline collapse
non-monotonic defensive interaction
```

without adding another engine feature.
