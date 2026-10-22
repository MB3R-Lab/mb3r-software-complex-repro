# Frozen application CI example — run 38788504177

Static archive of the public report produced by [manual application CI run 38788504177](https://github.com/MB3R-Lab/opentelemetry-demo-mb3r-stack/actions/runs/38788504177), attempt 1. The HTML and JSON were captured from the public report site on 22 October 2026.

## Open the report

Download [index.html](index.html) and open it in a browser. It is self-contained: the chart is an inline SVG and no external plotting service or JavaScript is required. Download [result.json](result.json) into the same directory if you also want the report's relative "Report JSON" link to work. GitHub's normal HTML-file view displays source code, not the rendered report.

The original HTML and JSON are retained byte-for-byte. Their SHA-256 digests, source URLs, source modification dates, run identity, application commits and selected inputs are in [provenance.json](provenance.json). The local `.gitattributes` disables line-ending conversion for these two files so their hashes survive Windows checkouts. The original page's "latest published run" notice describes the source website; this archived copy does not update.

The [live Pages report](https://mb3r-lab.github.io/opentelemetry-demo-mb3r-stack/) can change after another manual run. Always compare its displayed run ID with the GitHub Summary before presenting them as one result.

## What this example shows

| Configuration | Recommendation/catalog replicas | Model availability | Certified grid tolerance | Decision |
| --- | --- | --- | --- | --- |
| Baseline | 2/2 | 0.979900 | 0.12 | PASS |
| Candidate | 2/1 | 0.931200 | 0.03 | FAIL |

Both application configurations use commit `607de3071ef594d4e0149a38887d6b8e9f8d87d9`. The workflow collects two fresh healthy requests and 42 normalized spans per configuration. It then invokes the public Procrustes 0.1.0, Bering 1.0.0 and Sheaft 1.3.1 CLIs from the same immutable v1.3.3 toolchain image.

The policy has SLO 0.95, independent eligible-replica failure probability 0.05, channel liveness 0.995, 40,000 trials and seed 2719. It requires a minimum certified tolerance of 0.10 and permits a tolerance drop of at most 0.02. The candidate violates the endpoint threshold, minimum tolerance and regression rules. Full parameters and all evaluated points are in the preserved JSON.

This manual workflow compares one baseline/candidate pair. It does not run all four configurations from the paper in one report. It does not satisfy or replace the required pull-request gate.

## Relationship to the main demonstration

- The main artifact's `docker compose run --rm demo` replays recorded observations from application revision `be466eb395294e0b5c2f91f7e76bed475953eaa6` and evaluates four declared model variants.
- This consumer CI example deploys its own baseline and candidate revisions, collects fresh healthy evidence, and invokes the same public tools. It does not invoke the artifact's reproduction entrypoint.
- Both paths use the same selected operation and explicit required-call semantics, but they have different evidence and provenance. Component-importance values need not be identical.
- Availability remains model-based. This CI run does not inject chaos faults, measure replicated failover or calibrate failure probabilities.
- This CI report is archived separately from the package's v1.3.3 operation references; it is not used as a replay baseline.

See the [application CI guide](https://github.com/MB3R-Lab/opentelemetry-demo-mb3r-stack/blob/main/mb3r/ci/README.md), [manual-demo guide](https://github.com/MB3R-Lab/opentelemetry-demo-mb3r-stack/blob/main/mb3r/ci/DEMO.md), and [main artifact README](../../../README.md).

Only the already public derived report and JSON are archived here. Raw responses, full trace files and run artifacts are not recreated or claimed to be retained by this snapshot.
