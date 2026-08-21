# Procrustes static fit

- Goal: `POST /api/checkout` on `frontend`
- Goal status: **EXPECTED**
- Fingerprint: `sha256:53817d39a4ee0416e3c1ac294afa83a1e87d2f765122c0f943013b3e1dde039c`
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
| `cart` | CONFIRMED | DECLARED | CONFIRMED | CONFIRMED | CONFIRMED |
| `checkout` | CONFIRMED | DECLARED | CONFIRMED | CONFIRMED | CONFIRMED |
| `frontend` | CONFIRMED | DECLARED | CONFIRMED | CONFIRMED | CONFIRMED |
| `payment` | CONFIRMED | DECLARED | CONFIRMED | CONFIRMED | CONFIRMED |
| `shipping` | CONFIRMED | DECLARED | CONFIRMED | CONFIRMED | CONFIRMED |

## Limitations

- DECLARED instrumentation and propagation remain runtime assumptions until matching traces are observed.
- Procrustes evaluates static eligibility only; Bering must validate runtime evidence, coverage, and signal quality.
