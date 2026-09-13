# Muster Lab v0 UX contract

Primary loop:

> choose defenses → replay → inspect → change composition → compare first divergence

## Toy incident

- Frozen six-event incident.
- Defense budget: 4 human-facing controls.
- Vedette / Picket / Reserve choices.
- Presets:
  - Start clean
  - Strong but brittle
  - Diverse prevention
  - Want something counterintuitive? (JACKPOT)

A human-facing Reserve may expand to the review/escalation controls required to
make that response path executable. Those controls are still explicit Go engine
controls; Python only selects the package.

## HF-inspired public-record abstraction

Uses the existing `examples/hf-july-2026/scenario.yaml`.

No artificial budget. The interactive control pool exposes:
- correlated worker-compromise observation;
- privileged hostPath admission control;
- service-connector credential scoping;
- critical page + worker isolation response path.

## Boundary

Python does not decide replay semantics; it selects controls and presents returned traces.
