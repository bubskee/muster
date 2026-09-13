# Muster MVP Development Plan

> **Goal:** Build the smallest useful counterfactual replay engine for Watchline.
>
> Hold an incident trace fixed, vary the defensive composition, and observe how the trajectory changes.

## Scope

Muster is **not** a cyber range, RL environment, network simulator, SIEM emulator, or attacker model.

For the MVP, Muster should answer one question well:

> Given the same incident sequence, how do different combinations of **Vedettes**, **Pickets**, and **Reserves** change the outcome?

The first fixture will be a small, semantically faithful micro-replay derived from the July 2026 Hugging Face intrusion.

---

## MVP Invariants

Keep these true until the first re-evaluation gate:

- Incident events are fixed and deterministic.
- Defensive composition is the independent variable.
- Observation does **not** imply prevention.
- Prevention does **not** imply recovery.
- Vedette, Picket, and Reserve roles are explicit.
- Events describe semantic security transitions, not raw telemetry.
- Controls act on events/state; they do not simulate real infrastructure.
- No attacker agent chooses new actions.
- No stochasticity before the deterministic model proves useful.

---

## Initial Execution Model

For each event:

```text
check preconditions
      ↓
eligible Vedettes observe
      ↓
eligible Pickets intercept
      ↓
if not blocked:
    apply event effects
      ↓
eligible Reserves react
      ↓
apply reserve effects
      ↓
record trace
```

Detection must be able to occur without changing the incident trajectory.

That separation is a core Watchline property, not an implementation detail.

---

## Core Types

The exact Go API can evolve, but the MVP needs roughly these concepts.

```go
type Event struct {
    ID            string
    Time          int
    Kind          EventKind
    Preconditions []Condition
    Effects       []Effect
    Tags          []string
}
```

```go
type Role string

const (
    Vedette Role = "vedette"
    Picket  Role = "picket"
    Reserve Role = "reserve"
)
```

```go
type Control interface {
    Evaluate(Event, State) Outcome
}
```

```go
type Outcome struct {
    Observed bool
    Blocked  bool
    Actions  []Effect
}
```

The role should initially be metadata on a common control abstraction rather than three unrelated APIs. Some real controls may eventually span roles.

---

# Commit Plan

## Commit 1 — Event, Effect, State

Implement the minimum semantic model:

- `Event`
- `Effect`
- `Condition`
- `State`
- application of effects to state

No controls yet.

Success condition:

> A sequence of events can mutate state deterministically.

---

## Commit 2 — Deterministic Replay

Implement the replay loop:

```text
Scenario → ordered Events → State transitions → RunResult
```

Support:

- precondition checks
- skipped/unreachable events
- applied effects
- event-by-event trace

No Watchline controls yet.

Success condition:

> A fixed scenario produces a reproducible terminal state and trace.

---

## Commit 3 — Explicit V/P/R Controls

Add the defensive abstraction:

- `Control`
- explicit `Role`
- event matching
- evaluation result

Roles:

```text
Vedette  → observe
Picket   → block/intercept
Reserve  → respond/recover
```

Success condition:

> The engine knows *which role* acted and records it in the trace.

---

## Commit 4 — Observation vs Blocking

Make the central semantic distinction executable.

A Vedette may observe an event while the event still succeeds.

A Picket may block an event.

A Reserve may alter state after an event.

Required trace behavior:

```text
E3 credential-access
  observed by: token-read-detector
  blocked by:  none
  effects:     credential:k8s-service-account
```

Success condition:

> "Excellent watchmen, open gates" is representable without hacks.

---

## Commit 5 — `hf-worker-to-node` Fixture

Add the first incident micro-replay.

Initial six-event path:

```text
E1  malicious dataset reaches worker
E2  worker executes attacker-controlled code
E3  service-account credential accessed
E4  Kubernetes permissions/API discovered
E5  privileged/hostPath workload created
E6  node access obtained
```

The fixture should remain deliberately modest.

Call it:

```text
hf-worker-to-node
```

Do **not** call it a full Hugging Face incident replay.

Success condition:

> With no controls, the trace reaches node access.

---

## Commit 6 — Vedettes-Only Watchline

Add the first meaningful defensive configuration.

Example composition:

```text
Vedette: worker execution detector
Vedette: credential-access detector
Vedette: Kubernetes discovery detector
```

No prevention.

No recovery.

Expected result:

- early detection
- no blocked events
- incident still reaches node access
- non-zero `events_after_detection`

Success condition:

> Muster demonstrates that strong observation alone does not necessarily change the trajectory.

---

## Commit 7 — Picket / Interleaved Watchline

Add at least one preventive composition.

Candidate controls:

```text
Picket: scoped workload identity
Picket: deny privileged hostPath workloads
```

Optionally combine with existing Vedettes:

```text
V worker execution
P workload identity
V credential reuse
P admission control
```

Expected result:

- same incident fixture
- same initial state
- different terminal trajectory
- a clear first blocked transition

Success condition:

> Changing only defensive composition materially changes the replay outcome.

---

## Commit 8 — Terminal Summary

Produce the first useful comparison output.

Minimum per-run summary:

```text
furthest_event
first_observation
first_block
events_after_detection
terminal_state
```

Useful additional metrics if they fall out naturally:

```text
blast_radius
time_to_detection
time_to_containment
```

Example:

```text
HF-MICRO / vedettes-only
  furthest_event:         E6 node-access
  first_observation:      E2 worker-rce
  first_block:            none
  events_after_detection: 4
  terminal_state:         node-compromised

HF-MICRO / admission-control
  furthest_event:         E5 privileged-pod
  first_observation:      E2 worker-rce
  first_block:            E5 privileged-pod
  events_after_detection: 2
  terminal_state:         contained
```

---

# RE-EVALUATION GATE

**Stop after Commit 8.**

Do not automatically proceed to more features.

At this point, inspect the experiment and ask:

### 1. Does the abstraction reveal anything?

Does comparing Watchlines make the Watchline argument clearer than prose alone?

In particular, can Muster make failures such as this obvious?

```text
excellent detection
+ weak prevention
+ slow/no recovery
= large post-detection incident progression
```

### 2. Does the output distinguish architectures meaningfully?

If two superficially capable defensive stacks produce materially different trajectories because their V/P/R composition differs, Muster is doing useful work.

### 3. Is the engine simpler than the explanation?

If demonstrating the concept requires a large amount of simulator machinery, the abstraction may be wrong.

### 4. Are we accidentally rebuilding an existing tool?

Re-check against:

- CyberBattleSim
- CybORG
- Atomic Red Team
- CALDERA
- existing Hugging Face incident emulators

Muster should remain focused on **counterfactual defensive composition**, not generic cyber simulation or telemetry replay.

### Gate decision

Proceed only if the answer is approximately:

> Holding the incident fixed and varying Watchline composition produces a compact, legible, and non-trivial comparison that supports the Watchline thesis.

If instead the output amounts to:

> "Controls that block attacks block attacks,"

stop and redesign before adding features.

---

# Explicitly Out of Scope Before the Gate

Do **not** add:

- network topology simulation
- graph traversal machinery
- adaptive attacker agents
- reinforcement learning
- shell command execution
- Kubernetes emulation
- ECS / Elastic ingestion
- real SIEM integrations
- MITRE ATT&CK framework machinery
- generalized plugin systems
- GUIs
- notebooks as the primary engine
- stochastic control behavior
- probabilistic attacker behavior
- automatic incident extraction
- CyberBattleSim integration
- real-time scheduling
- human-in-the-loop timing models

These may become useful later. None is required to validate Muster.

---

# Likely Post-Gate Extensions

Only consider these after Commit 8 survives review.

## Latency

Especially important for Reserves:

```yaml
latency: 3
```

This enables questions such as:

> How much damage occurs between detection and containment?

## Reliability

```yaml
reliability: 0.85
```

This enables redundancy and failure-rate experiments.

Adding stochasticity implies:

- seeded runs
- repeated trials
- distributions
- confidence intervals

Do not incur that complexity prematurely.

## Expanded Hugging Face Micro

Possible additional events:

```text
E7  cloud/cluster credentials harvested
E8  lateral cluster access
E9  unusual public-service C2
```

## Presentation Layer

Once the engine is stable:

```text
scenario YAML ─┐
               ├─> Muster engine ─> RunResult.json ─> notebook / figures
watchline YAML ┘
```

Useful CLI target:

```bash
muster run scenarios/hf-worker-to-node.yaml watchlines/*.yaml
```

Possible outputs:

- human-readable trace
- JSON
- CSV

The notebook should consume Muster output; **Muster is not the notebook**.

---

# Candidate Repo Layout

```text
muster/
├── cmd/
│   └── muster/
├── engine/
│   ├── replay.go
│   ├── state.go
│   └── result.go
├── model/
│   ├── event.go
│   ├── control.go
│   └── effect.go
├── scenarios/
│   └── hf-worker-to-node.yaml
├── watchlines/
│   ├── baseline.yaml
│   ├── vedettes-only.yaml
│   ├── pickets.yaml
│   └── interleaved.yaml
├── notebooks/
│   └── hf_micro.ipynb
└── README.md
```

Do not optimize this layout before the core model works.

---

# Definition of MVP Success

Muster passes its first muster if:

1. One fixed Hugging Face micro-scenario can be replayed deterministically.
2. Vedette, Picket, and Reserve behavior remain semantically distinct.
3. A detection-heavy Watchline can observe an incident without stopping it.
4. A differently composed Watchline changes the terminal trajectory.
5. The terminal summary makes that difference immediately legible.
6. The implementation remains small enough that the research idea, not simulator engineering, is the center of gravity.

**Next decision point: Commit 8, terminal summary.**
