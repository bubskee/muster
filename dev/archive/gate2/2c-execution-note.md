# Gate 2C Execution Resolution

The frozen Gate 2 design specifies the 2×2 architecture matrix but does not
specify a treatment matrix.

Before executing 2C, we resolve that ambiguity as follows:

- hold telemetry failure timing at BETWEEN, because the BETWEEN history effect
  in 2A is the evidence that triggered conditional Gate 2C;
- cross all four sensing/review architectures with all three already-frozen
  review treatments: NONE, SOC-FAILURE, AGENT-FAILURE;
- introduce no new control or failure semantics.

Static-null prediction:

| Observation | Review | NONE | SOC-FAILURE | AGENT-FAILURE |
|---|---|---:|---:|---:|
| correlated V1+V2 | correlated H1+H2 | secure | NODE | secure |
| correlated V1+V2 | diverse H1+A1 | secure | NODE | secure |
| diverse V2+V3 | correlated H1+H2 | secure | NODE | secure |
| diverse V2+V3 | diverse H1+A1 | secure | secure | NODE |

Prediction: any crossover is expected from ordinary dependency propagation.
If Muster reproduces this table, Gate 2C adds architectural insight but no
additional evidence for stateful replay.
