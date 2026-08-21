# Procrustes static fit

- Goal: `POST /wrk2-api/post/compose` on `nginx-web-server`
- Goal status: **EXPECTED**
- Fingerprint: `sha256:468a3e2513b15187bd6e6aeb2399b92d2739d3ca916d62a96b3ecb7e5ceb04b1`
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
| `compose-post-service` | CONFIRMED | DECLARED | DECLARED | DECLARED | CONFIRMED |
| `nginx-web-server` | CONFIRMED | DECLARED | DECLARED | DECLARED | CONFIRMED |

## Detected tracing systems

| System | Role | Service | Status | Source |
|---|---|---|---|---|
| `jaeger` | backend | `—` | DECLARED | `compose/docker-compose.yaml` |
| `jaeger` | configuration | `—` | DECLARED | `compose/jaeger-config.yml` |
| `jaeger` | configuration | `compose-post-service` | DECLARED | `compose/src/compose-post-service/jaeger-config.yml` |
| `jaeger` | configuration | `nginx-web-server` | DECLARED | `compose/src/nginx-web-server/jaeger-config.yml` |
| `jaeger` | instrumentation | `compose-post-service` | CONFIRMED | `compose/src/compose-post-service/ComposePostService.cpp` |
| `jaeger` | instrumentation | `nginx-web-server` | CONFIRMED | `compose/src/nginx-web-server/jaeger-config.json` |
| `jaeger` | propagation | `nginx-web-server` | DECLARED | `compose/src/nginx-web-server/nginx.conf` |
| `opentracing` | instrumentation | `compose-post-service` | CONFIRMED | `compose/src/compose-post-service/ComposePostService.cpp` |
| `opentracing` | propagation | `compose-post-service` | DECLARED | `compose/src/compose-post-service/TraceContext.cpp` |

## Limitations

- Backend-specific tracing may require a source adapter before Bering ingestion; EXPECTED means that the static collection path is complete, not that runtime trace quality is proven.
- DECLARED instrumentation and propagation remain runtime assumptions until matching traces are observed.
- Procrustes evaluates static eligibility only; Bering must validate runtime evidence, coverage, and signal quality.
