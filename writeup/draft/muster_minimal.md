# Muster: where do defensive interventions actually change an attack?

Bobby Faber — [affiliation]
With Apart Research

## Abstract

Post-incident analysis asks what would have stopped an attack. That is harder than it looks: one control may stop a single path while leaving another open, excellent detection may produce no response, and a late intervention may still prevent the outcome that matters. Muster is a small replay engine that holds a compressed incident trace fixed, applies different defensive controls, and shows which downstream attacker capabilities remain reachable. On a public-record abstraction of the July 2026 Hugging Face intrusion, an admission policy and a scoped connector credential each cut one of two escalation paths; a correctly paged worker isolation cuts both. On a synthetic incident, adding a functioning early-isolation control to a composition that was already containing the attack causes the attack to succeed, because isolation removes the access needed to generate the evidence a later revocation depends on. Muster is a counterfactual reasoning aid, not a simulator. It evaluates controls against one modeled trace.

## 1. Introduction

Post-incident analysis asks three questions: what happened, what failed, what would have stopped it. The last is harder than it looks. For Hugging Face, many interventions look obvious after reconstructing the incident, and they are not the same intervention.

Muster makes those counterfactuals executable. It takes a compressed attack trajectory, applies a chosen set of controls, and reports which attacker capabilities remain reachable. Controls are grouped by role, not by product or layer:

- **Vedette** — observe: detect, correlate, classify, escalate.
- **Picket** — intervene: prevent or interrupt a transition.
- **Reserve** — recover: contain, revoke, rebuild.

The roles are interleaved, not a pipeline. A Vedette changes what the defender knows; a Picket or Reserve changes what the attacker can do. Knowing is not stopping.

Contributions:

1. A ten-transition abstraction of the Hugging Face intrusion, replayed under four control sets.
2. A non-monotone defensive interaction on a synthetic incident.
3. The same replay on Anthropic's PyPI incident, where responsibility crosses three organizations.

## 2. Related work

Existing frameworks describe attacker techniques, defensive controls, and response playbooks. Muster sits between them: incident → defensive actions → roles → counterfactual replay → what did this defense buy. It does not have to prove that attack graphs or fault trees cannot represent the same incident. It has to justify why the argument was made executable: the paper explains a result, the replay lets you discover it.

## 3. Methods

**Replay.** A scenario is an ordered list of events with preconditions and effects on a fact set. For each event: preconditions are checked (unmet → skipped); Vedettes observe and may add facts; Pickets may block, in which case effects do not occur; otherwise effects apply; review, criticality, and escalation controls fire in sequence, each gated on the previous; Reserves respond and may remove attacker facts. Output is a per-event trace and a terminal state. Deterministic, no attacker agent, no timing beyond order.

**Hugging Face scenario.** Ten transitions from the published technical timeline: two initial vectors (dataset local-file read exposing worker secrets; template injection giving worker RCE), then two independent paths — privileged/hostPath pod → node root → cluster secrets → mesh → internal network, and service-connector discovery → shared `system:masters` credential replay → cluster-admin. Hugging Face reports no admission policy rejected privileged pods and that correlated signals never received correct criticality or escalation. Both are modeled. Action counts, timestamps, C2, exfiltration, persistence are omitted. Not a forensic reconstruction.

**Toy scenario.** Six events, ten-control pool, four-control budget in the notebook.

**PyPI scenario.** Eight events from Anthropic's report, crossing evaluation operator, public registry, and third-party vendor.

**Development gates.** The first HF fixture was a linear chain, on which every result collapsed to the earliest Picket. A second gate asked whether stateful replay exposes anything a static dependency model cannot; default verdict was collapse. One result earned replay: a persisted alert flipped the outcome between two runs identical in every other respect, and a memoryless-reviewer ablation removed the effect. Verdict: keep tiny.

**Reproduce.** `go run . experiment --file experiments/hf-july-2026.yaml`

## 4. Results

**Table 1. Hugging Face.**

| Watchline | Node root | Internal network | Cluster admin |
|---|:---:|:---:|:---:|
| observed signal, no page | yes | yes | yes |
| privileged/hostPath admission policy | no | no | yes |
| cluster-scoped connector identity | yes | yes | no |
| critical page + worker isolation | no | no | no |

Blocking privileged workloads cuts one branch and leaves the connector path open. Scoping the connector credential cuts that branch and leaves node compromise possible. Early escalation and worker isolation cuts both. Different controls buy different pieces of the incident. An admission Picket can kill the node-root branch while the independent connector path still reaches cluster-admin.

**Table 2. Toy: adding a defense makes it worse.**

| Composition | Outcome |
|---|---|
| no defense | UNDETECTED & UNCONTAINED |
| Vedettes only | DETECTED, NOT CONTAINED |
| matched Picket | PREVENTED |
| V1 + V2 + R1 (revoke credential) | RESPONDED & CONTAINED |
| V1 + V2 + R1 + R2 (add early isolation) | ADVERSE DESPITE RESPONSE |

Mechanism: early worker isolation removes the access needed for the later discovery event. That event would have generated the evidence required to revoke the stolen credential. The credential survives and is replayed. Containment is not causally independent of investigation.

**Table 3. PyPI: one incident, three owners.**

| Owner | Watchline | Third-party execution | Vendor DB |
|---|---|:---:|:---:|
| eval operator | monitor only | yes | yes |
| eval operator | seal egress | no | no |
| vendor | sandbox scanner | yes | no |
| vendor | detect + revoke | yes | no |

The Watchline does not belong to one organization. It follows the incident.

## 5. Discussion and Limitations

A single frozen incident makes Pickets look almost unfairly powerful: once the path is known, place a control on the exact transition. "What would have stopped this?" is the retrospective question and identifies the Picket. "What will the next event look like?" is the prospective question and asks whether the composition still works when you guessed the Picket wrong. Vedettes maximize what you can notice; Pickets cover the transitions you think most likely; Reserves preserve recovery when the prediction is wrong.

**Limitations.** Muster evaluates defenses conditional on a modeled trajectory; a single replay does not estimate a control's value under uncertainty over future attacks. No attacker adaptation. Compression choices are load-bearing. No timing, cost, or probability. The non-monotone result is a constructed fixture, not a finding about Hugging Face.

**Future work.** Single trace → what does this control buy against this incident. Suite of plausible traces → how robust is this Watchline across attacker choices.

## 6. Conclusion

Hold an incident trace fixed, vary the controls, watch the trajectory move. Different controls buy different pieces; a stronger-looking composition can be worse; the defense can span organizations. A small executable argument.

## Code and Data

https://github.com/bubskee/muster (MIT). Scenarios in `examples/`, experiment in `experiments/`, notebook in `notebooks/`.

## References

- Hugging Face (2026). Anatomy of a Frontier Lab Agent Intrusion. https://huggingface.co/blog/agent-intrusion-technical-timeline
- Elastic Security Labs (2026). Exploring the Hugging Face Breach. https://www.elastic.co/security-labs/threat-command/ai-agent-attack-detection-hugging-face-breach
- Anthropic (2026). Investigating three real-world incidents in our cybersecurity evaluations. https://www.anthropic.com/news/investigating-incidents-cybersecurity-evals
- Anthropic (2026). Alignment assessment: cybersecurity incidents. https://www.anthropic.com/research/alignment-assessment-cybersecurity-incidents

## Appendix: Limitations and Dual-Use Considerations

[Required. Expand from §5 limitations. Dual-use: no exploit code or non-public detail; scenarios are compressions of published reports.]

## LLM Usage Statement

[Required. Write accurately.]
