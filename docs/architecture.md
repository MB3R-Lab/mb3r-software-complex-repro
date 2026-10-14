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
              ├─ published-result replay
              └─ endpoint component criticality
              │
              ▼
     comparison and verification
      ├─ analytics.json
      ├─ comparison.json
      └─ semantic SHA-256 manifest
```

The artifact runner invokes pinned public CLI contracts. It does not import internal packages from the three tools and does not calculate the resilience model independently.

The separate fail-closed path stops at the boundary under test: incomplete evidence is retained as a Procrustes `BLOCKED` report, while an unsupported Bering contract is rejected by Sheaft before simulation. Canonical `result.json` records make both outcomes inspectable without treating an expected process failure as an artifact failure.

The extended evaluation keeps three concerns separate. The decision walkthrough uses the normal Procrustes → Bering → Sheaft handoffs for three explicit replica maps. The validation matrix compares public CLI outputs with fixed diagnoses, exact trace-derived sets, and closed-form numeric oracles. The cost profile measures the same commands but does not participate in semantic reference hashing because timing and memory are machine-dependent.

## Component responsibilities

| Component | Input in this artifact | Responsibility | Output consumed downstream |
|---|---|---|---|
| Procrustes | analysis goal and minimal static evidence | conservative evidence and capability preflight | Bering overlay and traceable assessment report |
| Bering | archived topology plus Procrustes overlay | construct and validate the versioned system model | model, discovery snapshot, and quality reports |
| Sheaft | Bering model plus replay or criticality configuration | evaluate endpoint profiles and paired component counterfactuals | reports, model views, summaries, sweeps, and rankings |
| Artifact runner | tool outputs plus published table metadata | orchestration, normalization, hashing, and comparison | reference manifest and human/machine-readable comparisons |

## Contract boundaries

The important handoffs are stored rather than hidden in an in-process integration:

1. `procrustes/bering.overlay.yaml` carries the preflight-derived overlay into Bering.
2. `bering/model.json` carries services, replicas, edges, endpoints, predicates, provenance, and failure eligibility into Sheaft.
3. `sheaft/report.json` carries tool-produced endpoint and aggregate results into the comparison layer.
4. `criticality/report.json` carries Sheaft's complete endpoint-specific service/edge Birnbaum rankings; the runner derives only compact top-component views from it.

`toolchain-lock.json` pins the exact source revision and declared version of every component. `reference/manifest.json` binds that toolchain release to normalized SHA-256 hashes for each expected semantic output.

## Determinism boundary

Repeated runs compare semantic normal forms. The following runtime-only differences are removed or canonicalized before hashing:

- generation timestamps and recomputation duration;
- machine-local paths and file-derived identifiers;
- case-equivalent Bering discovery references;
- the non-semantic selection of one OpenTelemetry evidence field by Procrustes;
- JSON object ordering and whitespace.

Topology, contracts, predicates, replica counts, simulation parameters, profile results, criticality rankings, and published comparisons remain within the hash boundary. The complete rules are recorded in `reference/manifest.json`.

## Container trust boundary

The published image contains prebuilt binaries from the three commits in `toolchain-lock.json`, the artifact runner, cases, and reference outputs. GitHub Actions also performs an independent source-build verification from those exact commits. Thus the fast reviewer path and the source-verification path converge on the same checked-in semantic references.
