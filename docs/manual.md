# Artifact manual

## Intended user

This artifact is intended for a reviewer, researcher, or tool user who wants to inspect an end-to-end resilience-analysis workflow without deploying the evaluated systems or rebuilding the three component tools.

## Prerequisite

Install Docker with Compose support and ensure the Docker daemon is running. The normal path does not require Go, Kubernetes, the OpenTelemetry Demo, DeathStarBench, a load generator, or a chaos framework.

The candidate image tag, available only after publication, is:

```text
ghcr.io/mb3r-lab/mb3r-software-complex-repro:v1.3.4
```

The candidate `v1.3.4` uses one version across the Docker tag, documentation, and `toolchain.version` in the lock file. Component versions remain independent pins. `toolchain.reference_release` stays at `1.3.3` because the archived reference outputs are historical. Publication requires the new toolchain to pass the full reproduction checks; already published tags are not repointed.

## Reviewer walkthrough

### 1. Pull and run the short demonstration

```bash
docker compose pull demo
docker compose run --rm demo
```

The command performs these steps on the observed recommendations operation:

1. Procrustes checks the pinned deployment evidence against the recommendations goal and emits a Bering overlay.
2. Bering discovers the operation and three directed RPC edges from recorded healthy traces. The external predicate names mandatory calls without removing either catalog edge.
3. Sheaft checks six observed fault states against that predicate and contrasts reachability-only semantics.
4. Four separate model analyses compare the baseline, reduced configuration, and two one-replica additions under the same assumptions. Birnbaum inspection uses the reduced model and the same profile.
5. The runner checks analytical availability and the stored reference summary. Responses, traces, model inputs, full Sheaft reports, and a readable summary remain inspectable.

The final line must be:

```text
operation-walkthrough-ok states=6 variants=4 endpoint="frontend:GET /api/recommendations"
```

### 2. Inspect the generated evidence

Compose retains new files under `output/operation-walkthrough/` (the direct `docker run` command needs an explicit bind mount to retain them):

```text
output/operation-walkthrough/
├── procrustes/
├── observed-model.json
├── predicate.json
├── baseline/
├── current/
├── add-recommendation/
├── add-catalog/
├── semantic-check/
├── reachability-only-check/
├── result.json
└── summary.md
```

Each configuration directory retains the Bering model/snapshot, replica overlay, analysis configuration, and `sheaft/report.json`. The current report also contains complete component importance. Inputs and actual HTTP/trace evidence are under `cases/operation-semantics/`. See [operation walkthrough](operation-walkthrough.md) for the success criterion and scope.

### 3. Run the complete verification

```bash
docker compose run --rm reproduce
```

This runs:

- `otel-demo-all-blocking`;
- `otel-demo-async`;
- `social-network-norepl`;
- `social-network-repl`.

Every case is executed twice. The runner first checks equality between the two semantic runs and then checks 16 files per case against `reference/manifest.json`. It also regenerates and verifies `comparison.json` and `FULL_COMPARISON.md`.

### 4. Run the Birnbaum criticality walkthrough

```bash
docker compose run --rm criticality
```

This reviewer-facing scenario runs the archived OpenTelemetry async case once, verifies the complete 16-file case surface against the reference manifest, and prints one `criticality-ok` record per endpoint. Each record identifies the top service and edge, their Birnbaum importance estimates, and their Wilson confidence intervals. The full rankings remain in `criticality/report.json`, while `criticality-walkthrough.json` retains the compact walkthrough view when `--work-dir` is bind-mounted.

The final line must be:

```text
criticality-walkthrough-ok toolchain=1.3.3 case=otel-demo-async endpoints=4 components=39 trials=12000 files=16
```

Use `--criticality-walkthrough --case <name>` to run the same scenario on any of the four archived variants.

### 5. Inspect fail-closed behavior

```bash
docker compose run --rm failures
```

This is a separate two-scenario walkthrough, not an additional scientific experiment. It verifies that incomplete static evidence makes Procrustes emit a stored `BLOCKED` assessment and exit with code `2`, and that Sheaft rejects an unsupported Bering contract before creating an analysis report. Inspect the generated canonical evidence under `failure-walkthroughs/` when using a bind-mounted `--work-dir`.

The final line must be:

```text
failure-walkthroughs-ok toolchain=1.3.3 scenarios=2
```

### 6. Exercise an SLO regression and remediation

This is the retained original archival scenario, not the current central demonstration. Its reachability-only predicate omits catalog work. The new operation uses observed traces and mandatory calls, and therefore has different numerical results.

```bash
docker compose run --rm decision
```

The selected OpenTelemetry Demo operation is `frontend:GET /api/recommendations` with SLO `0.95`. The controlled baseline has two `recommendation` and two `ad` replicas. The proposed cost reduction removes one `recommendation` replica; the gate must report a compatible boundary regression from `0.20` to `0.04`. Birnbaum importance nominates `recommendation` for inspection but is not itself treated as a causal or release verdict. The corrected candidate preserves its two replicas and instead removes one `ad` replica. A separate Sheaft run must recover the `0.20` certified boundary and pass the gate.

The final line must be:

```text
decision-walkthrough-ok toolchain=1.3.3 endpoint="frontend:GET /api/recommendations" tolerance=0.20->0.04->0.20 decisions=pass,fail,pass
```

### 7. Run the controlled validation matrix

```bash
docker compose run --rm validation
```

This scenario checks controlled static evidence and diagnoses in Procrustes, exact service/edge/endpoint sets from constructed Bering traces, small Sheaft models with closed-form answers, and strict rejection of an unsupported contract. The stochastic oracle fixes the independent fail-stop distribution, failure population, 50,000 trials, and seed `2711`. Its output records each expected value, observed value, absolute error, and tolerance. Success requires all 15 checks to pass.

### 8. Measure descriptive application cost

```bash
docker compose run --rm cost --benchmark-repeats 5
```

The result reports min/median/max wall time for a complete Procrustes → Bering → Sheaft pass and for repeating Sheaft analysis on the saved Bering model. Peak child-process RSS is included on Linux. Treat these values as measurements of the stated machine, not as a comparative performance claim.

Run all three scenarios together with `docker compose run --rm evaluation --benchmark-repeats 5`.

### 9. Run one alternative case

```bash
docker run --rm \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:v1.3.4 \
  --case otel-demo-async --repeat 1
```

List accepted names with:

```bash
docker run --rm \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:v1.3.4 \
  --list-cases
```

## Persist generated outputs

On Linux or macOS:

```bash
mkdir -p output
docker run --rm --user "$(id -u):$(id -g)" \
  --mount "type=bind,source=$PWD/output,target=/output" \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:v1.3.4 \
  --repeat 2 --work-dir /output
```

In PowerShell with Docker Desktop:

```powershell
New-Item -ItemType Directory -Force output
$artifactOutput = (Resolve-Path output).Path
docker run --rm `
  --mount "type=bind,source=$artifactOutput,target=/output" `
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:v1.3.4 `
  --repeat 2 --work-dir /output
```

The destination will contain one directory per case plus the two aggregate comparison files.

## Interpretation

`reference/FULL_COMPARISON.md` reports Sheaft minus the published value for every row. Social Network replays the published fixed-size replica-slot protocol. OpenTelemetry uses the archived representative graph, while the article table aggregates archived graph runs; its small row deltas are retained and must not be interpreted as byte-identical reproduction of an aggregate over different graphs.

The published live/chaos values are comparison context only. The artifact never claims to regenerate them.

Separately, the new operation's live collection validates its predicate on six specified states. It does not calibrate failure probabilities, measure replicated failover, or establish universal handler equivalence. Ordinary endpoint thresholds use point estimates; sweep certification uses pointwise Wilson intervals; `max_regression` compares discrete certified-grid metrics. Their difference is not a confidence bound for a change in a true continuous threshold, and the whole release decision has no claimed simultaneous 95% coverage.

Criticality is an additional explanatory analysis, not a result transcribed from either prior article. It uses a 0.10 independent-replica failure baseline and 12,000 paired trials. Birnbaum values and Wilson intervals are conditional on the archived topology, predicates, and configured probabilities; overlapping intervals do not establish a reliable ordering, and importance is not used as a release gate.

The controlled validation matrix tests implementation behavior against constructed or closed-form oracles. It is not empirical validation against a new production system and does not establish that the mathematical model is generally valid. In the decision walkthrough, only the re-analysis of the modified candidate establishes its modelled effect; the Birnbaum ranking merely nominates a component for inspection.

## Source-build path for maintainers

The source-build path verifies that the current checkouts match the locked commits before it builds anything:

```bash
make reproduce-paper
make reproduce-paper-operation
make reproduce-paper-criticality
make reproduce-paper-failures
make reproduce-paper-evaluation BENCHMARK_REPEATS=3
```

It is intentionally separate from the container path and is not required of reviewers.
