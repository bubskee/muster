# Muster

**A small incident-replay harness for asking where defensive interventions actually change an attack.**

Muster takes a compressed incident trajectory, applies different defensive controls to it, and shows which downstream attacker capabilities remain reachable.

It was built for the Apart Research Incident Response Sprint alongside **Watchline**, a role-based way of thinking about layered defense.

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

The motivating question is:

> **You cannot defend every transition. Where do you place your watchline?**

## Why?

Security incidents are usually analyzed after the fact:

* What happened?
* What failed?
* What control might have stopped it?

The last question is deceptively difficult.

A control may stop one branch of an attack while leaving another untouched. An excellent detection system may observe an intrusion without producing an effective intervention. A defense that appears late in an incident may still prevent the outcome we actually care about.

Muster makes those claims executable.

Instead of asking whether a defense is generically “good,” we replay the same incident trajectory under different defensive placements and inspect **what each intervention actually buys us**.

## Watchline

Watchline describes defensive mechanisms by the role they play rather than by product or layer:

**Vedette — observe.**
Detect, classify, correlate, or escalate evidence that something is wrong.

**Picket — intervene.**
Prevent or interrupt an attacker transition: sandboxing, admission policy, identity scope, network isolation, and similar controls.

**Reserve — recover.**
Respond after compromise: revoke credentials, rebuild workloads, restore clean state, rotate secrets, or contain affected systems.

These roles can recur at every layer. They are not a pipeline.

```text
vedette → picket → reserve
    ↘       ↘       ↘
      interleaved across the system
```

Muster is the executable counterpart: a way to test Watchline placements against concrete incident paths.

## Hugging Face example

The main example is a deliberately compressed abstraction of the **July 2026 Hugging Face intrusion**, based on Hugging Face's published technical timeline.

It is **not a forensic reconstruction**.

The scenario preserves the relationships needed to ask defensive counterfactuals, including two independent paths from an initially compromised worker toward broader infrastructure compromise.

We replay that trajectory under four Watchlines:

| Watchline                            | Node root | Internal network | Cluster admin |
| ------------------------------------ | --------: | ---------------: | ------------: |
| observed signal, no page             |       yes |              yes |           yes |
| privileged/hostPath admission policy |        no |               no |           yes |
| cluster-scoped connector identity    |       yes |              yes |            no |
| critical page + worker isolation     |        no |               no |            no |

The interesting result is not that one control “wins.”

It is that **different controls buy different pieces of the incident**.

Blocking privileged pods kills one path but leaves the connector path available. Scoping the connector credential kills that path but leaves node compromise possible. Correctly escalating an early signal and isolating the worker cuts both downstream branches.

That boundary is what Muster is designed to make inspectable.

See [`examples/hf-july-2026/`](examples/hf-july-2026/) for the scenario, controls, provenance, and modeling caveats.

## Run it

Requires Go.

```bash
git clone https://github.com/bubskee/muster
cd muster

go run . experiment --file experiments/hf-july-2026.yaml
```

Machine-readable output is also available:

```bash
go run . experiment \
  --file experiments/hf-july-2026.yaml \
  --format json

go run . experiment \
  --file experiments/hf-july-2026.yaml \
  --format csv
```

Individual scenarios and control sets are ordinary YAML:

```bash
go run . replay \
  --scenario examples/gate2-between.yaml \
  --controls examples/gate2-broad-correlated.yaml
```

YAML describes the experiment; replay semantics remain in Go.

## What Muster is — and isn't

Muster is intentionally small.

It is **not** a cyber range, attacker emulator, probabilistic threat model, or claim that an incident can be faithfully reduced to ten transitions.

The abstraction throws away most real-world detail on purpose.

The goal is to preserve enough causal structure to ask a narrower question:

> **If this defensive action had succeeded here, which later attacker capabilities would still have been reachable?**

That makes Muster closer to a counterfactual reasoning aid than a full security simulator.

## Direction

Existing frameworks are very good at specifying adversary techniques, defensive controls, and incident-response actions.

Muster explores a complementary layer:

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

The longer-term question is whether this becomes useful for comparing defensive portfolios, reasoning about scarce attention and response capacity, and identifying where additional watchlines have the highest strategic value.

For now, it is a small executable argument.

**You cannot defend every transition. Where do you place your watchline?**

## Try Muster

The easiest way to understand Muster is to **build a Watchline and replay an incident**.

Open [`notebooks/muster_lab.ipynb`](notebooks/muster_lab.ipynb).

You’ll choose from:

* 🔭 **Vedettes** — observe
* 🛡️ **Pickets** — intervene
* 🔧 **Reserves** — respond / recover

Start with the **Toy incident**. Pick defenses, run the replay, and see what the attacker can still reach. Then change your Watchline and try again.

### Running the notebook

You’ll need **Go**, **Python 3**, and Jupyter.

<details>
<summary><strong>Need to install Go or Python?</strong></summary>

#### macOS

If you have [Homebrew](https://brew.sh/):

```bash
brew install go python
```

#### Windows

With `winget`:

```powershell
winget install GoLang.Go
winget install Python.Python.3
```

#### Ubuntu / Debian Linux

```bash
sudo apt update
sudo apt install golang-go python3 python3-pip
```

Check that both are available:

```bash
go version
python3 --version
```

On Windows, `python` may work instead of `python3`.

</details>

From the Muster repository directory:

```bash
python3 -m pip install jupyter ipywidgets
python3 -m jupyter notebook
```

Your browser should open automatically. Navigate to:

```text
notebooks/muster_lab.ipynb
```

Then choose **Run All**.

No Go or Python knowledge is required to use the lab.

Once the toy model makes sense, switch to the **HF-inspired public-record abstraction** and try defending the compressed Hugging Face incident.

> **You cannot defend every transition. Where do you place your watchline?**


```md
## Passing Muster: Role-Interleaved Defense for AI Infrastructure
_Lessons from the Hugging Face compromise_
```

### 🚀 Sprint Project Note

Muster was built as part of the Apart Research Incident Response Hackathon / Sprint.

## YAML replay

Muster can load scenarios and control sets from YAML without changing replay semantics.

```bash
go get gopkg.in/yaml.v3@v3.0.1

go run . replay \
  --scenario examples/gate2-between.yaml \
  --controls examples/gate2-broad-correlated.yaml
```

Emit the full trace as JSON for notebooks or other presentation tooling:

```bash
mkdir -p runs

go run . replay \
  --scenario examples/gate2-early.yaml \
  --controls examples/gate2-broad-correlated.yaml \
  --json > runs/gate2-early-broad.json

go run . replay \
  --scenario examples/gate2-between.yaml \
  --controls examples/gate2-broad-correlated.yaml \
  --json > runs/gate2-between-broad.json
```

YAML is intentionally declarative. It maps to existing `engine.Scenario` and
`engine.EventControl` values; replay semantics remain in Go.

## Experiments

Run a batch of scenario/control-set pairs and render the result as Markdown,
JSON, or CSV:

```bash
go run . experiment --file experiments/hf-july-2026.yaml
go run . experiment --file experiments/hf-july-2026.yaml --format json
go run . experiment --file experiments/hf-july-2026.yaml --format csv
```

Experiment files are presentation/orchestration only. Each run still loads a
normal scenario and control set and calls the existing replay engine. Metrics
are terminal-state facts selected for the report.

### Public-record Hugging Face example

`examples/hf-july-2026/` contains a deliberately compressed abstraction of the
July 2026 Hugging Face intrusion, based on Hugging Face's published technical
timeline. It is not a forensic reconstruction. The example preserves two
independent lateral-movement paths and compares three defensive
counterfactuals against an observed-signal/no-page baseline.

See `examples/hf-july-2026/README.md` for scope and provenance.

## writeup thoughts

> Existing frameworks specify what responders should do. Muster is a small executable layer for asking how those actions compose against a concrete incident trajectory.

one-month-more target:
```
CACAO / IR playbook
       ↓
extract defensive actions
       ↓
Watchline roles
       ↓
Muster counterfactual replay
       ↓
"what does this playbook actually buy us
 against this incident trace?"
 ```
