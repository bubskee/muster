# Hugging Face July 2026 — compressed public-record example

This example is a deliberately small **public-record abstraction**, not a
forensic reconstruction of the July 2026 Hugging Face intrusion.

Primary source:

- Hugging Face, *Anatomy of a Frontier Lab Agent Intrusion: A Technical
  Timeline of the July 2026 Incident* (2026-07-27):
  https://huggingface.co/blog/agent-intrusion-technical-timeline

Useful independent defensive mapping:

- Elastic Security Labs, *Exploring the Hugging Face Breach: mapping AI agent
  tactics to Elastic Defend* (2026-07-31):
  https://www.elastic.co/security-labs/threat-command/ai-agent-attack-detection-hugging-face-breach

## What the scenario preserves

The ten semantic transitions preserve several public-record relationships that
matter for defensive counterfactuals:

1. Malicious dataset configuration produced local-file disclosure from the
   production worker, exposing environment credentials and worker source.
2. A separate template-injection vector produced arbitrary code execution in
   the worker pod.
3. The worker foothold enabled Kubernetes and cloud discovery.
4. A privileged pod with a host filesystem mount enabled node root. Hugging
   Face states that no admission policy rejected privileged/hostPath pods.
5. Node root enabled cluster-secret theft and then mesh-network enrollment.
6. Separately from the mesh path, the compromised worker could reach an
   internal service connector. A shared connector credential bound to
   `system:masters` yielded cluster-admin across multiple clusters.
7. Hugging Face reports that security signals were correlated into a coherent
   attack signal but failed to receive the correct criticality and on-call
   escalation.

The model intentionally omits action counts, exact timestamps, C2 mechanics,
exfiltration, supply-chain exploration, denied probes, persistence, and many
other details. Those details are important to the incident but unnecessary for
this demonstration.

## The four runs

`observed-no-page.yaml` abstracts the disclosed detection/escalation failure: a
Vedette records worker compromise, but no response follows. It does **not**
claim that the real signal came from the exact event modeled here.

The other three Watchlines are explicitly counterfactual:

- `admission-policy.yaml` blocks the privileged/hostPath pod. It should stop
  node root and the secret/mesh branch while leaving the separate connector
  path available.
- `connector-scope.yaml` blocks cross-cluster replay of the shared connector
  identity. It should stop cluster-admin through that branch while leaving
  node-root/mesh compromise possible.
- `page-and-isolate.yaml` models a correctly prioritized alert followed by
  worker isolation. It should cut both downstream branches after the initial
  worker RCE.

This is the point of the example: **different controls buy different pieces of
the incident**, and a replay makes those counterfactual boundaries inspectable
without pretending to reproduce every attacker action.

## Run it

From the repository root:

```bash
go run . experiment --file experiments/hf-july-2026.yaml
```

For machine-readable output:

```bash
go run . experiment --file experiments/hf-july-2026.yaml --format json \
  > runs/hf-july-2026.json

go run . experiment --file experiments/hf-july-2026.yaml --format csv \
  > runs/hf-july-2026.csv
```

Expected high-level shape:

| Watchline | node root | internal network | cluster admin |
|---|---:|---:|---:|
| observed signal, no page | yes | yes | yes |
| privileged/hostPath admission policy | no | no | yes |
| cluster-scoped connector identity | yes | yes | no |
| critical page + worker isolation | no | no | no |

These are outputs of the **compressed model**, not empirical claims that any
single intervention would certainly have produced the same outcome in the
real incident.
