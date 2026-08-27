# Procrustes-Bering-Sheaft demonstration artifact

[![Reproduce software-complex results](https://github.com/MB3R-Lab/mb3r-software-complex-repro/actions/workflows/reproduce.yml/badge.svg)](https://github.com/MB3R-Lab/mb3r-software-complex-repro/actions/workflows/reproduce.yml)

This repository is the standalone demonstration and reproduction artifact for the Procrustes → Bering → Sheaft resilience-analysis toolchain. It packages pinned builds, archived evidence fixtures, published baselines, reference outputs, and the workflow that connects the three tools.

The artifact covers two published studies:

- OpenTelemetry Demo: the representative 16-service trace-discovered graph, four endpoint predicates, both edge semantics, and all five failure fractions from the [AINA 2026 article](https://doi.org/10.1007/978-3-032-23304-2_24);
- DeathStarBench Social Network: the archived 12-service Jaeger dependency graph, both replica modes, endpoint predicates, workload weights, and all five failure fractions from the [ICSE-NIER 2026 article](https://doi.org/10.1145/3786582.3786823).

No deployment, trace collection, load generation, or chaos experiment is performed. Procrustes inspects static evidence, Bering constructs the resilience model, and Sheaft calculates the reported profiles. The Go artifact runner only orchestrates those CLIs, checks their outputs, and calculates comparison statistics; it does not reimplement the resilience analysis.

## Quick start

Prerequisite: Docker with Compose support. No Go installation or source build is required.

```bash
git clone https://github.com/MB3R-Lab/mb3r-software-complex-repro.git
cd mb3r-software-complex-repro
docker compose run --rm demo
```

The demo pulls the versioned container and runs the replicated Social Network case once through all three tools. Success ends with:

```text
reproduce-paper-ok toolchain=1.2.0 cases=1 repeats=1
```

The equivalent command without Compose is:

```bash
docker run --rm --pull=always \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:toolchain-1.2.0
```

## Full reproduction

Run all four case variants twice, verify deterministic semantic outputs, and compare every aggregate row reported by both articles:

```bash
docker compose run --rm reproduce
```

The complete readable comparison is checked in as [`reference/FULL_COMPARISON.md`](reference/FULL_COMPARISON.md). Machine-readable row deltas, MAE, RMSE, maximum and signed error, and Pearson correlation are in [`reference/comparison.json`](reference/comparison.json).

Available cases can be listed or selected directly:

```bash
docker run --rm \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:toolchain-1.2.0 \
  --list-cases

docker run --rm \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:toolchain-1.2.0 \
  --case otel-demo-async --repeat 1
```

## What is verified

- the exact component revisions in [`toolchain-lock.json`](toolchain-lock.json);
- the Procrustes preflight report and Procrustes-to-Bering overlay;
- the Bering model, discovery snapshot, and quality reports;
- the Sheaft model, report, endpoint results, and sweeps;
- semantic stability across repeated executions after documented normalization of runtime-only fields;
- complete row-level comparison with Social Network Table 1 and OpenTelemetry Demo Table 2.

The canonical source-build run is the [GitHub Actions workflow](https://github.com/MB3R-Lab/mb3r-software-complex-repro/actions/workflows/reproduce.yml). It checks out the locked revisions of [Procrustes](https://github.com/MB3R-Lab/Procrustes), [Bering](https://github.com/MB3R-Lab/Bering), and [Sheaft](https://github.com/MB3R-Lab/Sheaft), then independently verifies the same checked-in references.

## Documentation

- [`docs/manual.md`](docs/manual.md): commands, output locations, and reviewer walkthrough;
- [`docs/architecture.md`](docs/architecture.md): component responsibilities and handoff contracts;
- [`docs/data-provenance.md`](docs/data-provenance.md): origin and transformation of case inputs and published baselines;
- [`docs/troubleshooting.md`](docs/troubleshooting.md): common Docker, platform, and verification failures;
- [`CONTRIBUTION_BOUNDARY.md`](CONTRIBUTION_BOUNDARY.md): claims supported by this artifact and claims deliberately excluded.

## Maintainer source run

A source run is optional and is not the reviewer path. It requires Go plus sibling checkouts named `Procrustes`, `Bering`, and `Sheaft` at the exact commits in `toolchain-lock.json`.

```bash
make reproduce-paper
```

Override checkout locations with `PROCRUSTES_REPO`, `BERING_REPO`, and `SHEAFT_REPO`. Use `make reproduce-paper-update` only after intentionally changing the locked toolchain or case inputs and reviewing the resulting reference diff.

Generated files are written below `.tmp/paper-reproduction/`. The product repositories contain neither this workflow nor paper-specific outputs.

## License

The artifact runner and repository documentation are available under the [MIT License](LICENSE). Fixture provenance and upstream notices are documented in [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md).
