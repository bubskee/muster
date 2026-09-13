# Suggestions (not in Bobby's material)

Everything here is optional. Nothing in `muster_minimal.md` depends on it. Full audit with provenance per claim is in `AUDIT_for_bobby.md`.

## Corrections to the repo's own prose

1. **Ablation claim.** `results.md` says "neither L1 nor R2 alone produced the adverse outcome." On the notebook toy, isolate-worker *alone* does reach node access (there is no revocation to lose). Accurate statement, verified: starting from the containing composition, adding isolate-worker without its escalation link changes nothing; the full early-containment path is required. Declaration order does not matter (verified by reversing IDs).
2. **PyPI metric.** `public:malicious-package` is removed by the final `registry-removes-package` event in every run, so it can never be a terminal-state metric. Use execution and vendor-DB access.
3. **README table vs. paper.** Keep them character-identical; a judge will diff them mentally.

## Claims to confirm before printing

- Gate 2 was preregistered (predictions written before runs). The prediction tables and tests exist; sequence is unverifiable from the repo. If not true, describe the tests, drop "before."
- "Keep tiny / expand not available" was your rule, not ChatGPT's framing.
- The `dev/archive` narrative (linear-chain Gate 1, Picket tripwire origin) is what actually happened.

## Additions worth considering

- **One sentence mapping roles to NIST CSF** (observe/block/respond ≈ Detect/Protect/Respond–Recover). Judges are detection engineers; without it the vocabulary reads as renaming.
- **Related work, two lines:** attack graphs (Sheyner et al. 2002) and attack–defense trees (Kordy et al. 2014) model prevention, not evidence-gated response; SecureLayer7's "which action number does your control fire on" is the closest prior proposal (sprint Resources tab).
- **Assumption-failure list** for the appendix: (a) both HF paths pass through worker access — if not, the isolation row changes; (b) stolen credentials persist until removed — if they rotate, the Table 2 interaction vanishes; (c) skipped events never retry — blocked-path rows are upper bounds.
- **Cost column** for Table 1 (rough tiers), since Track 1 examples ask for it repeatedly.
- **State the containment standard as a Table 1**, not a control list: modeled escape path, controls, what each leaves reachable. Auditable without network access, which is the Track 1 judging criterion.

## Fifteen-minute gradeability

- README first screen, in your words: what it is, one command, Table 1.
- One command per table. `suggested-repo-additions/` has `experiments/toy-jackpot.yaml`, `experiments/anthropic-pypi-2026.yaml`, and the `subsets/` control files; verified to run and reproduce Tables 2–3. Otherwise Tables 2–3 need hand-typed `--only-controls` lists.
- Make the CLI the primary path; the notebook is optional (Go + Jupyter + ipywidgets is its own ten minutes).
- `dev/archive/` is visible, filenames read `opus-gate-redteam.md`, contents read generated. Add a one-line README inside ("working notes, LLM-assisted") and disclose in the LLM statement.
- Tag the commit cited in the paper.

## LLM statement must cover

ChatGPT-drafted prose (README, `dev/archive`, `UX.md`, writeup notes); LLM red-teaming of gate designs; Claude regenerating tables and drafting this skeleton at Josh's request; that the report was rewritten by you and every number checked.

## Style

Cut the aphorisms in the final ("Pickets win battles…", "It follows the incident", "watch the trajectory move"). They are the most recognizable LLM tell and the guidelines say generated-reading reports won't be scored. The minimal draft keeps a few because they are in your notes; drop any you wouldn't say.
