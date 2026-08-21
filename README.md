# Software-complex paper reproduction

This standalone repository contains the workflow and evidence for reproducing the software-toolchain results of Procrustes → Bering → Sheaft on two archived systems:

- OpenTelemetry Demo, using the archived 16-service trace-discovered graph and a representative checkout operation from the [published AINA 2026 study](https://doi.org/10.1007/978-3-032-23304-2_24);
- DeathStarBench Social Network, using the archived 12-service Jaeger dependency graph, replica map, endpoint predicates, and workload weights from the [published ICSE-NIER 2026 study](https://doi.org/10.1145/3786582.3786823).

It deliberately performs no live deployment, trace collection, load test, or chaos experiment. Static evidence is inspected by Procrustes; the resulting overlay is consumed by Bering; the Bering model is consumed by Sheaft. Analytics are derived by the Go reproduction program from tool outputs, not by standalone Python simulation scripts.

## Run

The canonical run is the GitHub Actions workflow in `.github/workflows/reproduce.yml`. It reads exact revisions from `toolchain-lock.json`, checks out all three software repositories, builds their CLIs, runs both cases twice, and compares the results with the checked-in semantic references.

For an optional developer run, prerequisites are Go and sibling checkouts named `Procrustes`, `Bering`, and `Sheaft` next to this repository. Their exact revisions must match `toolchain-lock.json`.

```bash
make reproduce-paper
```

The command verifies source revisions, builds all three CLIs, runs both cases twice, checks semantic byte stability after removal of explicitly documented runtime fields, and compares the outputs with `reference/manifest.json`.

Override sibling locations with `PROCRUSTES_REPO`, `BERING_REPO`, or `SHEAFT_REPO`. Use `make reproduce-paper-update` only when intentionally accepting a new checked-in reference set after reviewing the toolchain diff.

Maintainers can accept outputs downloaded from a completed CI run without rerunning the tools locally:

```bash
go run ./reproduce.go --accept-work-dir <downloaded-artifact-root>
```

Generated CI working files live under `.tmp/paper-reproduction/` and are uploaded as workflow artifacts. Checked-in references contain complete Procrustes reports and handoff artifacts, Bering models and quality reports, Sheaft reports, plus compact paper analytics.

This repository contains no Procrustes, Bering, or Sheaft product source. The product repositories contain no paper reproduction workflow or paper outputs.

## Evidence provenance

The topology inputs are normalized copies of existing archived results. Absolute workstation paths were replaced by stable `archive://` references; topology, replica counts, endpoint targets, and workload weights were not re-estimated here. The static manifests and tracing snippets are minimal, self-contained evidence fixtures derived from the corresponding deployments and are used only to exercise Procrustes' preflight contract.

Each case contains `publication.json` with the publisher DOI, proceedings location, artifact DOI, and the reported table values used as publication context. The Social Network `p=0.30` replicated aggregate is compared directly with the published Table 1 value. The current OpenTelemetry case is explicitly marked as a representative-endpoint run rather than a direct comparison with the four-endpoint aggregate in published Table 2.

See [CONTRIBUTION_BOUNDARY.md](CONTRIBUTION_BOUNDARY.md) for the publication claim boundary.
