# otel-demo-all-blocking reference summary

- Published comparison mode: `all-blocking`.
- Procrustes goal status: `EXPECTED`; assessed services: 5; blockers: 0.
- Bering model: 16 services, 23 edges (0 async), 4 endpoints.
- Published source: Refining Resilience Model Discovery: A Case Study on the Limited Role of Asynchronous Edges in the OpenTelemetry Demo, DOI `10.1007/978-3-032-23304-2_24`.

| p_fail | k | Sheaft | Published model | Δ model | Published live | Δ live |
|---:|---:|---:|---:|---:|---:|---:|
| 0.1 | 2 | 0.781676 | 0.781000 | +0.000676 | 0.683000 | +0.098676 |
| 0.3 | 4 | 0.609860 | 0.610000 | -0.000140 | 0.557000 | +0.052860 |
| 0.5 | 8 | 0.356375 | 0.356000 | +0.000375 | 0.360000 | -0.003625 |
| 0.7 | 10 | 0.250842 | 0.251000 | -0.000158 | 0.289000 | -0.038158 |
| 0.9 | 14 | 0.049971 | 0.050000 | -0.000029 | 0.172000 | -0.122029 |
