# social-network reference summary

- Procrustes goal status: `EXPECTED`; assessed services: 2; blockers: 0.
- Bering model: 12 services, 19 edges (0 async), 4 endpoints.
- Sheaft p=0.30 weighted availability: 0.305890.

- Published source: Model Discovery and Graph Simulation: A Lightweight Gateway to Chaos Engineering, DOI `10.1145/3786582.3786823`.
- Published comparison: actual `0.305890`, published `0.305400`, absolute delta `0.000490` (tolerance `0.015000`).

| Endpoint | Availability |
|---|---:|
| compose-post | 0.022450 |
| home-timeline | 0.285100 |
| nginx-web-server:POST /wrk2-api/post/compose | 0.662150 |
| user-timeline | 0.441950 |
