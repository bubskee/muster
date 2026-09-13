# Muster

**A small incident-replay harness for asking where defensive interventions actually change an attack.**

Muster takes a compressed attack trajectory, applies different defensive controls, and shows which downstream attacker capabilities remain reachable.

It was built alongside **Watchline**, a role-based model of layered defense.

```text
                  ATTACK TRAJECTORY
                         │
                         ▼
             ─────── transition ───────
                         │
             🔭 VEDETTE  │ observe
             🛡️ PICKET   │ intervene
             🔧 RESERVE  │ recover
                         │
                         ▼
             ─────── transition ───────
                         │
                        ...
```

> **You cannot defend every transition. Where do you place your watchline?**

## Why?

Post-incident analysis often asks:

* What happened?
* What failed?
* What would have stopped it?

That last question is harder than it looks.

One control may stop a single attack path while leaving another open. Excellent detection may produce no effective response. A late intervention may still prevent the outcome that matters.

Muster makes those counterfactuals executable.

## Watchline

Watchline groups defenses by role rather than product or infrastructure layer:

**🔭 Vedette — observe**
Detect, correlate, classify, or escalate evidence.

**🛡️ Picket — intervene**
Prevent or interrupt attacker transitions.

**🔧 Reserve — recover**
Contain compromise, revoke access, rebuild, or restore state.

These roles are **interleaved**, not a pipeline. A system may have all three at many different layers.

Muster lets you test those placements against an incident path.

## Try Muster

The easiest way to understand Muster is to **build a Watchline and replay an incident**.

Open [`notebooks/muster_lab.ipynb`](notebooks/muster_lab.ipynb).

Start with the **Toy incident**. Pick defenses, run the replay, and see what the attacker can still reach. Change your Watchline and try again.

No Go or Python knowledge is required to use the lab.

### Running the notebook

You’ll need **Go**, **Python 3**, and Jupyter.

<details>
<summary><strong>Need to install Go or Python?</strong></summary>

### macOS

With [Homebrew](https://brew.sh/):

```bash
brew install go python
```

### Windows

With `winget`:

```powershell
winget install GoLang.Go
winget install Python.Python.3
```

### Ubuntu / Debian

```bash
sudo apt update
sudo apt install golang-go python3 python3-pip
```

Check installation:

```bash
go version
python3 --version
```

On Windows, use `python` if `python3` is unavailable.

</details>

Clone Muster:

```bash
git clone https://github.com/bubskee/muster
cd muster
```

Install Jupyter and launch the notebook:

```bash
python3 -m pip install jupyter ipywidgets
python3 -m jupyter notebook
```

Navigate to:

```text
notebooks/muster_lab.ipynb
```

Then choose **Run All**.

Once the toy model feels intuitive, switch to the **HF-inspired public-record abstraction**.

## Hugging Face example

The main real-world example is a compressed abstraction of the **July 2026 Hugging Face intrusion**, based on Hugging Face's published technical timeline.

It is **not a forensic reconstruction**. It preserves enough structure to compare defensive counterfactuals across two independent attack paths.

| Watchline                            | Node root | Internal network | Cluster admin |
| ------------------------------------ | --------: | ---------------: | ------------: |
| observed signal, no page             |       yes |              yes |           yes |
| privileged/hostPath admission policy |        no |               no |           yes |
| cluster-scoped connector identity    |       yes |              yes |            no |
| critical page + worker isolation     |        no |               no |            no |

Different controls buy different pieces of the incident.

Blocking privileged workloads cuts one branch but leaves the connector path open. Scoping the connector credential cuts that branch but leaves node compromise possible. Early escalation and worker isolation cuts both downstream paths.

See [`examples/hf-july-2026/`](examples/hf-july-2026/) for provenance and modeling caveats.

## What Muster is — and isn't

Muster is intentionally small.

It is **not** a cyber range, attacker emulator, or full threat model.

It throws away most incident detail to preserve one question:

> **If this defense succeeded here, which later attacker capabilities would still be reachable?**

Think of it as a **counterfactual reasoning aid**, not a realistic simulator.

## CLI and reproducible experiments

The notebook sits on top of the Go replay engine. Technical users can run experiments directly:

```bash
go run . experiment --file experiments/hf-july-2026.yaml
```

JSON and CSV output are also available:

```bash
go run . experiment \
  --file experiments/hf-july-2026.yaml \
  --format json

go run . experiment \
  --file experiments/hf-july-2026.yaml \
  --format csv
```

Individual scenario/control pairs can be replayed directly:

```bash
go run . replay \
  --scenario examples/gate2-between.yaml \
  --controls examples/gate2-broad-correlated.yaml
```

Scenario configuration lives in YAML; replay semantics remain in Go.

## Direction

Existing frameworks describe attacker techniques, defensive controls, and response playbooks.

Muster explores one layer between them:

```text
incident / playbook
        ↓
defensive actions
        ↓
Watchline roles
        ↓
counterfactual replay
        ↓
what did this defense actually buy us?
```

Longer term, the question is whether this helps compare defensive portfolios, allocate scarce attention, and decide where the next watchline belongs.

For now, it is a small executable argument.

---

> “Everything is very simple in war, but the simplest thing is difficult.”
>
> — Carl von Clausewitz, *On War*

---

Built for the **Apart Research Incident Response Sprint**.
