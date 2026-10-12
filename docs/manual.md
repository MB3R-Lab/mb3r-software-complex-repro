# Artifact manual

## Intended user

This artifact is intended for a reviewer, researcher, or tool user who wants to inspect an end-to-end resilience-analysis workflow without deploying the evaluated systems or rebuilding the three component tools.

## Prerequisite

Install Docker with Compose support and ensure the Docker daemon is running. The normal path does not require Go, Kubernetes, the OpenTelemetry Demo, DeathStarBench, a load generator, or a chaos framework.

The published image is:

```text
ghcr.io/mb3r-lab/mb3r-software-complex-repro:toolchain-1.3.0
```

## Reviewer walkthrough

### 1. Pull and run the short demonstration

```bash
docker compose pull demo
docker compose run --rm demo
```

The command performs these steps on the `social-network-repl` fixture:

1. Procrustes checks whether the archived evidence satisfies the requested analysis goal and emits a Bering overlay.
2. Bering combines the archived topology with that overlay and emits a versioned model plus discovery-quality reports.
3. Sheaft consumes the Bering model and the archived analysis configuration and emits the resilience profile.
4. Sheaft independently runs the opt-in schema 1.3 criticality configuration and ranks service liveness and edge connectivity for every published endpoint in the case family.
5. The artifact runner derives comparison and top-component views from public tool outputs and verifies all semantic files against the checked-in reference manifest.

The final line must be:

```text
reproduce-paper-ok toolchain=1.3.0 cases=1 repeats=1
```

### 2. Inspect the generated evidence

The container verifies generated files internally. Stable examples of every output are available under:

```text
reference/social-network-repl/
├── procrustes/
├── bering/
├── sheaft/
├── criticality/
├── analytics.json
└── summary.md
```

The stage reports show the actual inter-tool handoffs. `criticality/report.json` is the full Sheaft ranking; `analytics.json` is derived after the tools finish and contains inventory, comparison, and compact top-service/top-edge views. Neither is an alternative simulator.

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

### 4. Inspect fail-closed behavior

```bash
docker compose run --rm failures
```

This is a separate two-scenario walkthrough, not an additional scientific experiment. It verifies that incomplete static evidence makes Procrustes emit a stored `BLOCKED` assessment and exit with code `2`, and that Sheaft rejects an unsupported Bering contract before creating an analysis report. Inspect the generated canonical evidence under `failure-walkthroughs/` when using a bind-mounted `--work-dir`.

The final line must be:

```text
failure-walkthroughs-ok toolchain=1.3.0 scenarios=2
```

### 5. Run one alternative case

```bash
docker run --rm \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:toolchain-1.3.0 \
  --case otel-demo-async --repeat 1
```

List accepted names with:

```bash
docker run --rm \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:toolchain-1.3.0 \
  --list-cases
```

## Persist generated outputs

On Linux or macOS:

```bash
mkdir -p output
docker run --rm --user "$(id -u):$(id -g)" \
  --mount "type=bind,source=$PWD/output,target=/output" \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:toolchain-1.3.0 \
  --repeat 2 --work-dir /output
```

In PowerShell with Docker Desktop:

```powershell
New-Item -ItemType Directory -Force output
$artifactOutput = (Resolve-Path output).Path
docker run --rm `
  --mount "type=bind,source=$artifactOutput,target=/output" `
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:toolchain-1.3.0 `
  --repeat 2 --work-dir /output
```

The destination will contain one directory per case plus the two aggregate comparison files.

## Interpretation

`reference/FULL_COMPARISON.md` reports Sheaft minus the published value for every row. Social Network replays the published fixed-size replica-slot protocol. OpenTelemetry uses the archived representative graph, while the article table aggregates archived graph runs; its small row deltas are retained and must not be interpreted as byte-identical reproduction of an aggregate over different graphs.

The published live/chaos values are comparison context only. The artifact never claims to regenerate them.

Criticality is an additional explanatory analysis, not a result transcribed from either prior article. It uses a 0.10 independent-replica failure baseline and 12,000 paired trials. Birnbaum values and Wilson intervals are conditional on the archived topology, predicates, and configured probabilities; overlapping intervals do not establish a reliable ordering, and importance is not used as a release gate.

## Source-build path for maintainers

The source-build path verifies that the current checkouts match the locked commits before it builds anything:

```bash
make reproduce-paper
make reproduce-paper-failures
```

It is intentionally separate from the container path and is not required of reviewers.
