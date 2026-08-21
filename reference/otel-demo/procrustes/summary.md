# Procrustes static fit

- Goal: `POST /api/checkout` on `frontend`
- Goal status: **EXPECTED**
- Fingerprint: `sha256:8c2dcf0af8803299aba1e5d20b80361fb594e7e8d765c6ecfab0b59f18069f8e`
- Required actions: 0

## Capabilities

| Capability | Status |
|---|---|
| `service_inventory` | **EXPECTED** |
| `candidate_topology` | **EXPECTED** |
| `operation_specific_model` | **EXPECTED** |
| `replica_aware_model` | **EXPECTED** |

## Services

| Service | Identity | Instrumentation | Exporter | Propagation | Replicas |
|---|---|---|---|---|---|
| `checkout` | CONFIRMED | DECLARED | CONFIRMED | CONFIRMED | CONFIRMED |
| `frontend` | CONFIRMED | DECLARED | CONFIRMED | CONFIRMED | CONFIRMED |

## Limitations

- DECLARED instrumentation and propagation remain runtime assumptions until matching traces are observed.
- Procrustes evaluates static eligibility only; Bering must validate runtime evidence, coverage, and signal quality.
