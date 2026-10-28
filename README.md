# Procrustes-Bering-Sheaft demonstration artifact

[![Reproduce software-complex results](https://github.com/MB3R-Lab/mb3r-software-complex-repro/actions/workflows/reproduce.yml/badge.svg)](https://github.com/MB3R-Lab/mb3r-software-complex-repro/actions/workflows/reproduce.yml)

This repository is the standalone demonstration and reproduction artifact for the Procrustes → Bering → Sheaft resilience-analysis toolchain. It packages pinned builds, archived evidence fixtures, published baselines, reference outputs, endpoint-specific component-criticality rankings, and the workflow that connects the three tools.

The primary demonstration uses one observed HTTP operation from a pinned application revision, checks its execution predicate, and compares two replica changes. It also preserves numerical replay of two published studies:

- OpenTelemetry Demo: the representative 16-service trace-discovered graph, four endpoint predicates, both edge semantics, and all five failure fractions from the [AINA 2026 article](https://doi.org/10.1007/978-3-032-23304-2_24);
- DeathStarBench Social Network: the archived 12-service Jaeger dependency graph, both replica modes, endpoint predicates, workload weights, and all five failure fractions from the [ICSE-NIER 2026 article](https://doi.org/10.1145/3786582.3786823).

The normal reviewer command replays recorded evidence without deploying the application. A separate GitHub Actions workflow collected real HTTP responses and Jaeger traces under six predeclared conditions; its inputs and evidence are retained in this repository. Procrustes inspects static evidence for the same operation, Bering constructs the model from observed healthy traces, and Sheaft evaluates it. The Go runner orchestrates the public CLIs and checks their outputs against response semantics, analytical oracles, and references; it is not another resilience simulator.

## Quick start

Prerequisite: Docker with Compose support. No Go installation or source build is required.

```bash
git clone https://github.com/MB3R-Lab/mb3r-software-complex-repro.git
cd mb3r-software-complex-repro
docker compose run --rm demo
```

The demo pulls the versioned container and runs the recommendations operation through all three tools. It compares a baseline, a replica reduction, and two one-replica additions to relevant components. Generated reports are retained in `output/operation-walkthrough/`. Success ends with:

```text
operation-walkthrough-ok states=6 variants=4 endpoint="frontend:GET /api/recommendations"
```

The equivalent command without Compose is:

```bash
docker run --rm --pull=always \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:v1.3.4
```

The candidate `v1.3.4` package pins its image and three component commits in `toolchain-lock.json`. Its `toolchain.reference_release` keeps the archived reference outputs attributed to `v1.3.3`; the new toolchain must reproduce their semantic and numerical fields before release.

## CI/CD example: the same tools in an application pipeline

Two complementary entry points use the same versioned Procrustes, Bering and Sheaft CLIs:

- **This demonstration package:** `docker compose run --rm demo` replays recorded evidence and evaluates the fixed model variants. It needs no application deployment and is the main entry point linked from the paper.
- **Application CI integration:** the [OpenTelemetry Demo example](https://github.com/MB3R-Lab/opentelemetry-demo-mb3r-stack/blob/main/mb3r/ci/README.md) deploys baseline and candidate configurations on a GitHub runner, collects fresh healthy HTTP/Jaeger evidence, and invokes the same tools. It does not run this package's reproduction entrypoint.

Use [Actions → MB3R parameterized demo](https://github.com/MB3R-Lab/opentelemetry-demo-mb3r-stack/actions/workflows/mb3r-demo.yml) to select replica counts, the SLO, failure assumptions, trials and seed; see its [manual and parameter reference](https://github.com/MB3R-Lab/opentelemetry-demo-mb3r-stack/blob/main/mb3r/ci/DEMO.md). A separate required check, **MB3R / recommendations gate**, controls PR acceptance. Manual exploration does not replace that check.

The [live HTML report](https://mb3r-lab.github.io/opentelemetry-demo-mb3r-stack/) shows the latest published manual run. For a permanent example, the [frozen report and JSON from run 38788504177](reference/ci-demo/38788504177/README.md) retain the original files, application revisions, parameters and checksums. Download the archived HTML and open it in a browser; compare run IDs before pairing a live report with a GitHub Summary.

The archived CI report documents a separate application run, with its own evidence and provenance. The operation replay below uses pinned reference inputs retained from `v1.3.3`. Both entry points estimate availability on a model; fresh healthy traces are not live failover measurements or calibrated failure probabilities.

## Observed operation and two candidate changes

The operation is `GET /api/recommendations?currencyCode=USD` on application commit `be466eb395294e0b5c2f91f7e76bed475953eaa6`, with fixed catalog data and the recommendation cache disabled. Success requires HTTP 200 and a nonempty product array. Six conditions (healthy, two service outages, three directed-channel outages) have two recorded requests each. The model retains both catalog calls; an external predicate contract 1.1 declares all three RPC channels mandatory. Reachability alone is not treated as execution success.

Both candidates add one replica to the reduced configuration: either recommendation or product-catalog. Both components have nonzero Birnbaum importance. All candidates use the same failure distribution, SLO, trial budget, seed, and grid. These are conditional model-based replica changes, not live replication-failover measurements or an automatically synthesized correction.

See [`docs/operation-walkthrough.md`](docs/operation-walkthrough.md), the preregistered [`plan.json`](cases/operation-semantics/plan.json), and the retained [`evidence/`](cases/operation-semantics/evidence/). The original archived scenario and its `0.20 → 0.04 → 0.20` results remain unchanged under `decision`; they are not results for the new predicate.

## Full reproduction

Run all four case variants twice, verify deterministic semantic outputs, compare every aggregate row reported by both articles, and generate Birnbaum service/edge rankings for all seven published endpoints:

```bash
docker compose run --rm reproduce
```

The complete readable comparison is checked in as [`reference/FULL_COMPARISON.md`](reference/FULL_COMPARISON.md). Machine-readable row deltas, MAE, RMSE, maximum and signed error, and Pearson correlation are in [`reference/comparison.json`](reference/comparison.json).

Available cases can be listed or selected directly:

```bash
docker run --rm \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:v1.3.4 \
  --list-cases

docker run --rm \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:v1.3.4 \
  --case otel-demo-async --repeat 1
```

## Birnbaum criticality walkthrough

Run the endpoint-specific component-importance scenario independently:

```bash
docker compose run --rm criticality
```

The walkthrough uses the archived OpenTelemetry async case by default, executes Procrustes and Bering, runs Sheaft's schema 1.3 Birnbaum analysis, verifies all 16 semantic case files, and prints the top service and edge with confidence intervals for each of the four endpoints. It ends with:

```text
criticality-walkthrough-ok toolchain=1.3.3 case=otel-demo-async endpoints=4 components=39 trials=12000 files=16
```

Select another archived case with `--case`, for example:

```bash
docker run --rm \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:v1.3.4 \
  --criticality-walkthrough --case social-network-repl
```

## Fail-closed walkthroughs

Run the two short negative demonstrations independently of the four scientific-result cases:

```bash
docker compose run --rm failures
```

The first supplies incomplete static evidence to Procrustes and requires exit code `2`, a `BLOCKED` goal, and explicit blocker rules. The second supplies an unsupported Bering model contract to Sheaft and requires exit code `1`, the precise compatibility rejection, and no generated analysis report. The wrapper returns success only when both expected failures are observed and ends with:

```text
failure-walkthroughs-ok toolchain=1.3.3 scenarios=2
```

## Extended evaluation

Run the retained controlled/archival evaluation suite (the primary observed-operation demonstration is `demo` or `operation`):

```bash
docker compose run --rm evaluation
```

It contains three separate scenarios:

1. `decision`: retain the original archived reachability-only predicate and recommendations/ad example. Its `0.20 → 0.04 → 0.20` sequence remains a conditional numerical regression fixture, not handler-semantic validation. It is superseded as the central demonstration by `operation`.
2. `validation`: run 15 controlled checks covering Procrustes positive/negative diagnoses, exact Bering service/operation/edge discovery from constructed traces, analytical Sheaft availability and Birnbaum oracles, confidence-certified SLO `pass`, `fail`, and `indeterminate` outcomes, and an unsupported handoff contract.
3. `cost`: report full-pipeline and saved-model-analysis wall time plus peak child-process RSS where the operating system exposes it. These are descriptive measurements for the stated machine, not a cross-system performance benchmark.

Run them independently with:

```bash
docker compose run --rm decision
docker compose run --rm validation
docker compose run --rm cost --benchmark-repeats 5
```

Checked-in deterministic summaries are in [`reference/evaluation/`](reference/evaluation/). Cost measurements are generated on demand and are deliberately not treated as reference values.

## What is verified

- the exact component revisions in [`toolchain-lock.json`](toolchain-lock.json);
- agreement between actual responses and the required-call predicate on six predeclared states, plus a reachability-only countercheck on the same discovered graph;
- two separately evaluated, same-budget replica candidates affecting relevant components;
- the Procrustes preflight report and Procrustes-to-Bering overlay;
- the Bering model, discovery snapshot, and quality reports;
- the Sheaft model, report, endpoint results, and sweeps;
- separate Sheaft 1.3 Birnbaum importance reports with complete service/edge rankings, paired-trial intervals, forced-down impact, and perfect-component gain;
- semantic stability across repeated executions after documented normalization of runtime-only fields;
- complete row-level comparison with Social Network Table 1 and OpenTelemetry Demo Table 2.
- criticality coverage for three Social Network endpoints and four OpenTelemetry endpoints in every represented deployment/edge-semantics variant.
- fail-closed behavior for incomplete Procrustes evidence and an unsupported Bering-to-Sheaft contract.
- an end-to-end SLO decision sequence with a compatible baseline, an explicitly rejected regression, a Birnbaum-nominated component, and an independently verified remediation;
- controlled oracle checks with fixed distributions, failure populations, trials, seeds, expected values, observed values, and complete discrepancies;
- exact operation-aware Bering discovery on constructed traces with zero missing and zero extra services, edges, and endpoints;
- descriptive full-pass and saved-model-analysis execution cost.

The canonical source-build run is the [GitHub Actions workflow](https://github.com/MB3R-Lab/mb3r-software-complex-repro/actions/workflows/reproduce.yml). It checks out the locked revisions of [Procrustes](https://github.com/MB3R-Lab/Procrustes), [Bering](https://github.com/MB3R-Lab/Bering), and [Sheaft](https://github.com/MB3R-Lab/Sheaft), then independently verifies the same checked-in references.

## Documentation

- [`docs/manual.md`](docs/manual.md): commands, output locations, and reviewer walkthrough;
- [`docs/operation-walkthrough.md`](docs/operation-walkthrough.md): current central scenario, assumptions, evidence, and interpretation;
- [`docs/architecture.md`](docs/architecture.md): component responsibilities and handoff contracts;
- [`docs/data-provenance.md`](docs/data-provenance.md): origin and transformation of case inputs and published baselines;
- [`docs/troubleshooting.md`](docs/troubleshooting.md): common Docker, platform, and verification failures;
- [`CONTRIBUTION_BOUNDARY.md`](CONTRIBUTION_BOUNDARY.md): claims supported by this artifact and claims deliberately excluded.

## Maintainer source run

A source run is optional and is not the reviewer path. It requires Go plus sibling checkouts named `Procrustes`, `Bering`, and `Sheaft` at the exact commits in `toolchain-lock.json`.

```bash
make reproduce-paper
make reproduce-paper-operation
make reproduce-paper-criticality
make reproduce-paper-failures
make reproduce-paper-decision
make reproduce-paper-validation
make reproduce-paper-cost BENCHMARK_REPEATS=5
```

Override checkout locations with `PROCRUSTES_REPO`, `BERING_REPO`, and `SHEAFT_REPO`. Use `make reproduce-paper-update` only after intentionally changing the locked toolchain or case inputs and reviewing the resulting reference diff.

Generated files are written below `.tmp/paper-reproduction/`. The product repositories contain neither this workflow nor paper-specific outputs.

## License

The artifact runner and repository documentation are available under the [MIT License](LICENSE). Fixture provenance and upstream notices are documented in [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md).
