# Data and fixture provenance

## Scope

The repository contains archived, static evidence needed to exercise the released toolchain. It does not contain or regenerate live chaos windows, load-test measurements, a running OpenTelemetry Demo, or a running DeathStarBench deployment.

## Published sources

| Case family | Published source | Archived artifact | Rows represented |
|---|---|---|---|
| Social Network | [ICSE-NIER 2026, DOI 10.1145/3786582.3786823](https://doi.org/10.1145/3786582.3786823) | [socialnet-resilience](https://github.com/a-a-k/socialnet-resilience), DOI 10.5281/zenodo.15332249 | Table 1, both replica modes and five failure fractions |
| OpenTelemetry Demo | [AINA 2026, DOI 10.1007/978-3-032-23304-2_24](https://doi.org/10.1007/978-3-032-23304-2_24) | [otel-demo-resilience](https://github.com/a-a-k/otel-demo-resilience), DOI 10.5281/zenodo.17703953 | Table 2, both edge semantics and five failure fractions |

The exact published numbers are transcribed into each case family's `publication.json` together with the paper title, venue, pages, DOI, table identifier, metric, and artifact location.

## Committed case inputs

### Social Network

- `topology-api.yaml`: normalized archived 12-service graph with the published replica interpretation;
- `topology-norepl.yaml`: the corresponding non-replicated interpretation;
- `goal.yaml`: evidence requirements supplied to Procrustes;
- `analysis.yaml`: endpoint predicates, weights, failure fractions, and fixed-size protocol supplied to Sheaft;
- `static/` and `static-norepl/`: minimal configuration and tracing fixtures used to exercise Procrustes' preflight contract.

### OpenTelemetry Demo

- `topology-api.yaml`: normalized representative 16-service graph with asynchronous edge annotations;
- `topology-all-blocking.yaml`: the same representative graph under the all-blocking interpretation;
- `goal.yaml`: evidence requirements supplied to Procrustes;
- `analysis.yaml`: the four endpoint predicates, equal workload mix, and published fixed-size failure schedule;
- `static/`: minimal Kubernetes and tracing-test fixtures used to exercise Procrustes' preflight contract.

## Normalization performed for this repository

Absolute workstation paths were replaced with stable `archive://` references. Machine-local identifiers and byte-derived digests are normalized only when reference hashes are calculated. The archived topology, replica counts, endpoint targets, workload weights, failure fractions, and published table values were not re-estimated.

The small static evidence directories are self-contained preflight fixtures derived from the corresponding deployment configurations. They are not asserted to be complete deployable copies of either system.

## Output provenance

Each reference case retains the stage boundary:

- `procrustes/`: assessment, collection plan, and Bering overlay;
- `bering/`: model, discovery snapshot, and quality reports;
- `sheaft/`: consumed model view, analysis report, and summary;
- `analytics.json`: inventory and published-result comparison derived from those outputs;
- `summary.md`: a readable rendering of the same case-level comparison.

The workflow is deterministic with respect to these archived inputs and the pinned toolchain. This is a software-complex reproduction claim, not a new empirical validation of the underlying mathematical model.
