# One observed operation, two replica corrections

Run `docker compose run --rm demo` (or `operation`). New files remain under
`output/operation-walkthrough/`. The reviewer command does not deploy OTel Demo;
it replays committed observations and invokes the pinned public tool CLIs.

## Evidence and execution condition

Application revision: `be466eb395294e0b5c2f91f7e76bed475953eaa6` of
`MB3R-Lab/opentelemetry-demo-mb3r-stack`. This is a new targeted experiment, not
an asserted reconstruction of either prior publication's application revision.
The goal, static scenario, Bering declaration, and Sheaft SLO all refer to
`GET /api/recommendations?currencyCode=USD`. Success is HTTP 200 with a nonempty
array of products with nonempty IDs. Empty/degraded results are not successful.
The seed catalog data are fixed and `recommendationCacheFailure` is disabled.

The actual frontend handler calls recommendation, which obtains catalog IDs via
`ListProducts`; frontend subsequently retrieves up to four products via
`GetProduct`. Both catalog edges are retained. PostgreSQL, flagd, frontend, and
telemetry are held healthy. The three channel states are shared by all RPCs in a
request, not resampled for each individual RPC.

Live collection ran only in [GitHub Actions 38718363369](https://github.com/MB3R-Lab/mb3r-software-complex-repro/actions/runs/38718363369).
Each condition has two requests and a preceding recovery check. The proxy closes
existing connections and prevents reconnects for the selected directed channel.
Service-outage states stop the actual container. Transport timeouts would make
the experiment inconclusive, not count as confirmed application failures.

| Predeclared condition | HTTP responses | Reachability-only predicate | Required-call predicate |
|---|---|---|---|
| Healthy | 200, 200; nonempty products | success | success |
| Recommendation stopped | 500, 500 | failure | failure |
| Catalog stopped | 500, 500 | failure | failure |
| Frontend → recommendation disabled | 500, 500 | failure | failure |
| Recommendation → catalog disabled | 500, 500 | **success (mismatch)** | failure |
| Frontend → catalog disabled | 500, 500 | **success (mismatch)** | failure |

The comparison uses the same discovered graph and mandatory targets. Changing
only the external predicate removes the two unexplained state mismatches. This
supports the narrowly scoped addition of `required_edges` in predicate contract
1.1, not removal of observed edges. Contract 1.0 and archived inputs are unchanged.
The 12 responses, raw Jaeger traces, 42 healthy normalized spans, deployment
configuration, container states, and image provenance are in
[`cases/operation-semantics/evidence/`](../cases/operation-semantics/evidence/).

## Fixed comparison protocol

[`plan.json`](../cases/operation-semantics/plan.json) fixes the candidates and
parameters before numerical results: independent eligible-replica failures with
probability 0.05, independent channel liveness 0.995, SLO 0.95, 40,000 trials,
seed 2719, the complete 15-point grid, minimum certified tolerance 0.10, and maximum
regression 0.02. Recommendation and catalog replicas form the service failure
population. Other infrastructure is held healthy. The tool's derived per-profile
seeds are also retained in full reports.

Each correction adds one replica to the reduced configuration. Both affected
services have nonzero importance. This equal replica-count budget does not imply
equal monetary/resource cost. The observed deployment had one replica per
service; these four configurations are explicit model what-ifs, not observed
replicated deployments or measured failover tests.

| Configuration (recommendation/catalog) | Availability | Pointwise 95% Wilson interval | Certified grid tolerance | Gate |
|---|---:|---|---:|---|
| Baseline 2/2 | 0.979900 | [0.978478, 0.981230] | 0.12 | pass |
| Reduced 2/1 | 0.931200 | [0.928678, 0.933639] | 0.03 | fail |
| Add recommendation 3/1 | 0.934675 | [0.932211, 0.937055] | 0.03 | fail |
| Add catalog 2/2 | 0.979900 | [0.978478, 0.981230] | 0.12 | pass |

The same-profile Birnbaum estimates on the reduced configuration are 0.934925
for recommendation (95% interval [0.932466, 0.937301]) and 0.983175 for catalog
([0.981867, 0.984390]). Perfect-component gains are 0.002337 and 0.049159.
Importance concerns aggregate service liveness: it does not itself establish the
benefit of adding a replica. Separate analysis of both candidates establishes
their conditional effects. The catalog addition restores the baseline rather
than saving a replica relative to it.

Example retained report excerpt for the recommendation addition:

```json
{
  "decision": "fail",
  "reason_ids": ["boundary_below_minimum", "boundary_regressed", "endpoint_below_threshold"],
  "certified_tolerance": 0.03,
  "reference_tolerance": 0.12,
  "regression": -0.09,
  "max_regression": 0.02
}
```

The fields above are a compact view of the actual decision and boundary rule,
not a claim that all fields occur in one object in the tool's full report.
The independent series oracle is
`A(f,r_rec,r_cat) = 0.995^3 (1-f^r_rec) (1-f^r_cat)`.
Exact availability is 0.9801556573 (2/2), 0.9334815784 (2/1), and
0.9357041536 (3/1). All discrepancies are retained; maximum absolute error is
0.0022815784 against the predeclared tolerance 0.006. Outputs were accepted from
[replay 38720567632](https://github.com/MB3R-Lab/mb3r-software-complex-repro/actions/runs/38720567632).

## What the decisions mean

- The ordinary endpoint threshold uses a **point estimate**.
- Certified grid tolerance uses **pointwise Wilson intervals** and a contiguous
  certified prefix of the fixed grid, not an exact continuous crossing.
- `max_regression` compares **two certified grid metrics**, not a confidence bound
  for the difference between true continuous thresholds. Separate minimum and
  regression rules preserve both diagnoses instead of short-circuiting the second.

There is no claimed simultaneous 95% guarantee for the complete release decision.
The live checks establish agreement only on the specified operation, data, cache
state, and six selected states. They do not calibrate probabilities, establish
general handler equivalence, model payload/order/cache behavior, or measure
replicated failover. The runner evaluates prepared changes; it does not synthesize
or deploy corrections.

## Reference and legacy boundaries

[`reference/operation/result.json`](../reference/operation/result.json) is checked
on every replay; the full four candidate reports and the two semantic-check
reports are retained alongside it for inspection. Those full reports preserve CI
paths/timing; the exact repeatability assertion is on the deterministic summary,
not byte equality of every raw report. Runtime handoffs are regenerated each run.

The old `decision` scenario and its `0.20 → 0.04 → 0.20` results remain tied to
their original reachability-only predicate. The 15 controlled checks and 20
archived aggregate comparisons remain separate. New operation evidence does not
retroactively validate archived handler assumptions or reproduce prior live
measurements.
