# otel-demo-async reference summary

- Published comparison mode: `async`.
- Procrustes goal status: `EXPECTED`; assessed services: 5; blockers: 0.
- Bering model: 16 services, 23 edges (3 async), 4 endpoints.
- Published source: Refining Resilience Model Discovery: A Case Study on the Limited Role of Asynchronous Edges in the OpenTelemetry Demo, DOI `10.1007/978-3-032-23304-2_24`.

| p_fail | k | Sheaft | Published model | Δ model | Published live | Δ live |
|---:|---:|---:|---:|---:|---:|---:|
| 0.1 | 2 | 0.781676 | 0.781000 | +0.000676 | 0.683000 | +0.098676 |
| 0.3 | 4 | 0.609860 | 0.610000 | -0.000140 | 0.557000 | +0.052860 |
| 0.5 | 8 | 0.356375 | 0.356000 | +0.000375 | 0.360000 | -0.003625 |
| 0.7 | 10 | 0.250842 | 0.251000 | -0.000158 | 0.289000 | -0.038158 |
| 0.9 | 14 | 0.049971 | 0.050000 | -0.000029 | 0.172000 | -0.122029 |

## Endpoint component criticality

Birnbaum importance uses `importance-baseline` with 12000 paired trials; values are conditional on the archived model and configured independent-replica baseline.

| Endpoint | Top service (I_B) | Top edge (I_B) | Baseline availability |
|---|---|---|---:|
| `frontend:GET /api/cart` | `cart` (1.0000) | `frontend|cart|sync|true` (0.0905) | 0.9032 |
| `frontend:GET /api/products` | `product-catalog` (1.0000) | `frontend|product-catalog|sync|true` (0.0096) | 0.8973 |
| `frontend:GET /api/recommendations` | `recommendation` (1.0000) | `frontend|recommendation|sync|true` (0.8992) | 0.8992 |
| `frontend:POST /api/checkout` | `shipping` (0.7320) | `checkout|payment|sync|true` (0.6573) | 0.6573 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering.
