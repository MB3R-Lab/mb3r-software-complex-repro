# Artifact architecture

## End-to-end flow

```text
archived configuration and topology
              │
              ▼
      Procrustes preflight
      ├─ report.json
      ├─ collection-plan.json
      └─ bering.overlay.yaml
              │
              ▼
        Bering discovery
      ├─ model.json
      ├─ snapshot.json
      └─ quality reports
              │
              ▼
         Sheaft analysis
      ├─ model.json
      ├─ report.json
      └─ summary.md
              │
              ▼
     comparison and verification
      ├─ analytics.json
      ├─ comparison.json
      └─ semantic SHA-256 manifest
```

The artifact runner invokes pinned public CLI contracts. It does not import internal packages from the three tools and does not calculate the resilience model independently.

## Component responsibilities

| Component | Input in this artifact | Responsibility | Output consumed downstream |
|---|---|---|---|
| Procrustes | analysis goal and minimal static evidence | conservative evidence and capability preflight | Bering overlay and traceable assessment report |
| Bering | archived topology plus Procrustes overlay | construct and validate the versioned system model | model, discovery snapshot, and quality reports |
| Sheaft | Bering model plus analysis configuration | evaluate endpoint and aggregate resilience profiles | report, model view, summaries, and sweeps |
| Artifact runner | tool outputs plus published table metadata | orchestration, normalization, hashing, and comparison | reference manifest and human/machine-readable comparisons |

## Contract boundaries

The important handoffs are stored rather than hidden in an in-process integration:

1. `procrustes/bering.overlay.yaml` carries the preflight-derived overlay into Bering.
2. `bering/model.json` carries services, replicas, edges, endpoints, predicates, provenance, and failure eligibility into Sheaft.
3. `sheaft/report.json` carries tool-produced endpoint and aggregate results into the comparison layer.

`toolchain-lock.json` pins the exact source revision and declared version of every component. `reference/manifest.json` binds that toolchain release to normalized SHA-256 hashes for each expected semantic output.

## Determinism boundary

Repeated runs compare semantic normal forms. The following runtime-only differences are removed or canonicalized before hashing:

- generation timestamps and recomputation duration;
- machine-local paths and file-derived identifiers;
- case-equivalent Bering discovery references;
- the non-semantic selection of one OpenTelemetry evidence field by Procrustes;
- JSON object ordering and whitespace.

Topology, contracts, predicates, replica counts, simulation parameters, results, and published comparisons remain within the hash boundary. The complete rules are recorded in `reference/manifest.json`.

## Container trust boundary

The published image contains prebuilt binaries from the three commits in `toolchain-lock.json`, the artifact runner, cases, and reference outputs. GitHub Actions also performs an independent source-build verification from those exact commits. Thus the fast reviewer path and the source-verification path converge on the same checked-in semantic references.
