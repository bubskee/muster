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
