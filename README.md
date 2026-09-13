# muster
incident replay/simulation harness

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
