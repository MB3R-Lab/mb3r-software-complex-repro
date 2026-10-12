# social-network-repl reference summary

- Published comparison mode: `repl`.
- Procrustes goal status: `EXPECTED`; assessed services: 2; blockers: 0.
- Bering model: 12 services, 19 edges (0 async), 4 endpoints.
- Published source: Model Discovery and Graph Simulation: A Lightweight Gateway to Chaos Engineering, DOI `10.1145/3786582.3786823`.

| p_fail | k | Sheaft | Published model | Δ model | Published live | Δ live |
|---:|---:|---:|---:|---:|---:|---:|
| 0.1 | 3 | 0.630053 | 0.628100 | +0.001953 | 0.696900 | -0.066847 |
| 0.3 | 7 | 0.306482 | 0.305400 | +0.001082 | 0.305400 | +0.001082 |
| 0.5 | 11 | 0.113918 | 0.114500 | -0.000582 | 0.095800 | +0.018118 |
| 0.7 | 16 | 0.012885 | 0.013200 | -0.000315 | 0.015500 | -0.002615 |
| 0.9 | 20 | 0.000000 | 0.000000 | +0.000000 | 0.000000 | +0.000000 |

## Endpoint component criticality

Birnbaum importance uses `importance-baseline` with 12000 paired trials; values are conditional on the archived model and configured independent-replica baseline.

| Endpoint | Top service (I_B) | Top edge (I_B) | Baseline availability |
|---|---|---|---:|
| `compose-post` | `unique-id-service` (0.5294) | `compose-post-service|media-service|sync|true` (0.4741) | 0.4741 |
| `home-timeline` | `post-storage-service` (0.8099) | `nginx-web-server|home-timeline-service|sync|true` (0.0003) | 0.7275 |
| `user-timeline` | `post-storage-service` (0.8993) | `nginx-web-server|user-timeline-service|sync|true` (0.0004) | 0.8077 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering.
