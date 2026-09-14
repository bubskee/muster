# JACKPOT sensitivity check — preregistration

## Question

Does the synthetic JACKPOT reversal behave as predicted under perturbations
that preserve or break the proposed causal mechanism?

This is not a random robustness sweep. The goal is to test the mechanism.

## Frozen mechanism

In the current JACKPOT fixture, adding early worker isolation succeeds locally
by removing `attacker:worker-access`.

That prevents the later `kubernetes-api-discovery-from-worker` event from
running. The later discovery path would otherwise support escalation and
revocation of `attacker:cluster-credential`.

With that evidence path suppressed, the credential survives, is replayed, and
`attacker:node-access` becomes reachable.

We therefore predict that the reversal depends on two conditions:

1. credential revocation depends on evidence generated after theft by a
   worker-dependent event; and
2. worker isolation occurs before that evidence is generated.

The reversal should not depend on unrelated downstream/topological details.

## Primary outcomes

For each condition, compare a matched pair in the same perturbed world:

- containing composition without added early isolation
- same composition with isolation added

Terminal metrics:

- `credential-held` → `attacker:cluster-credential`
- `node-access` → `attacker:node-access`

Define a **reversal** as:

- both metrics absent without added isolation; and
- both metrics present with added isolation.

## Preregistered conditions

| condition | perturbation | predicted reversal? | reason |
|---|---|:---:|---|
| original | none | **yes** | baseline JACKPOT mechanism |
| unrelated side branch | add topology that does not feed the evidence→revocation path | **yes** | causal mechanism unchanged |
| independent evidence | allow revocation from evidence that survives worker isolation | **no** | revocation no longer depends on suppressed discovery |
| isolation after evidence | move isolation until after revocation-enabling evidence is generated | **no** | isolation arrives too late to suppress the evidence |

## Implementation constraint

Prefer YAML-only changes.

Do not change engine semantics merely to obtain the predicted result. If one of
the perturbations cannot be represented using the existing engine, record that
as a limitation.

Keep this experiment separate from `experiments/toy-jackpot.yaml`.

## Interpretation

If all predictions hold:

> The reversal survives a non-causal perturbation but disappears when either
> preregistered causal condition is removed, supporting the proposed mechanism.

If a prediction fails:

1. preserve the failing result;
2. inspect why;
3. report the failure before modifying the fixture.

Any diagnostic experiment after seeing the failure is post hoc and should be
marked as such.

This experiment establishes only local mechanism consistency. It does not
estimate frequency, prevalence, or robustness across realistic incidents.

## Stop rule

Stop after these four conditions.

Do not add more perturbations because the first four look interesting.
