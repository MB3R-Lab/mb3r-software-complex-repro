# Demonstration artifact 1.2.0 — toolchain 1.3.2

The default reviewer command now exercises one observed operation, with a matching
preflight goal, real trace discovery, explicit required calls, SLO/regression
decisions, and two independently evaluated one-replica corrections.

```bash
docker compose run --rm demo
```

Outputs persist under `output/operation-walkthrough/`. The image is
`ghcr.io/mb3r-lab/mb3r-software-complex-repro:artifact-1.2.0`; exact component
commits are in `toolchain-lock.json`, including Sheaft v1.3.1. This artifact lock
is authoritative; the separately linked mb3r-stack distribution is not an
additional executable dependency or a claim that its release was changed.

## Validated scope

- New pinned-application experiment: six predeclared conditions, twelve actual
  HTTP responses, raw Jaeger traces, and recorded deployment/container state.
- Required-call predicate matches all six observed states. Reachability-only
  semantics disagree with both separately disabled catalog channels on the same
  discovered graph.
- Both one-replica candidates affect components with nonzero importance.
  Adding recommendation retains a failing 0.03 grid tolerance; adding catalog
  restores 0.12 and passes. Parameters and all numeric discrepancies are retained.
- The prior 15 controlled checks, 20 published-result comparisons, and archived
  `0.20 → 0.04 → 0.20` scenario are preserved under their original semantics.
- Source and container CI check the deterministic operation reference, legacy
  numerical references, rejection behavior, and packaged evaluation. A fresh job
  pulls the published image anonymously and runs the documented Compose command.

This is selected execution-state validation plus conditional model what-ifs,
not calibrated production probabilities, live replica-failover measurement,
universal handler equivalence, or automatically deployed remediation. See
`docs/operation-walkthrough.md` for assumptions and the three distinct statistical
decision mechanisms. No manuscript changes are part of this release.
