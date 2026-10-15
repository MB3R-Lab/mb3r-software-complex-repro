# Observed operation and candidate comparison

Every predeclared observed state agrees with the required-call predicate. Both catalog calls remain in the discovered graph.

| Configuration | Availability | Analytical oracle | Grid tolerance | Decision | Reasons |
|---|---:|---:|---:|---|---|
| baseline | 0.979900 | 0.980156 | 0.12 | pass | [gate_pass] |
| current | 0.931200 | 0.933482 | 0.03 | fail | [boundary_below_minimum boundary_regressed endpoint_below_threshold] |
| add-recommendation | 0.934675 | 0.935704 | 0.03 | fail | [boundary_below_minimum boundary_regressed endpoint_below_threshold] |
| add-catalog | 0.979900 | 0.980156 | 0.12 | pass | [gate_pass] |

Endpoint gates use point estimates. Grid certification uses pointwise Wilson intervals. Regression compares discrete certified grid metrics, not a confidence bound on the change of a continuous threshold.

Replica changes are conditional model analyses; the collected requests validate the selected execution predicate, not replication failover or operational failure probabilities.
