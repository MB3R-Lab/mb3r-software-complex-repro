# Full comparison with published results

These tables cover every aggregate row reported in Table 1 of the Social Network article and Table 2 of the OpenTelemetry Demo article. `Δ` is Sheaft minus the published value. Published live values are archival context; no live or chaos experiment is rerun here.

## Social Network — ICSE-NIER 2026

| Mode | p_fail | Sheaft | Published model | Δ model | Published live | Δ live |
|---|---:|---:|---:|---:|---:|---:|
| norepl | 0.1 | 0.418623 | 0.418200 | +0.000423 | 0.553300 | -0.134677 |
| norepl | 0.3 | 0.162548 | 0.161300 | +0.001248 | 0.177500 | -0.014952 |
| norepl | 0.5 | 0.045515 | 0.045400 | +0.000115 | 0.037600 | +0.007915 |
| norepl | 0.7 | 0.001323 | 0.001400 | -0.000077 | 0.006700 | -0.005377 |
| norepl | 0.9 | 0.000000 | 0.000000 | +0.000000 | 0.000000 | +0.000000 |
| repl | 0.1 | 0.630053 | 0.628100 | +0.001953 | 0.696900 | -0.066847 |
| repl | 0.3 | 0.306482 | 0.305400 | +0.001082 | 0.305400 | +0.001082 |
| repl | 0.5 | 0.113918 | 0.114500 | -0.000582 | 0.095800 | +0.018118 |
| repl | 0.7 | 0.012885 | 0.013200 | -0.000315 | 0.015500 | -0.002615 |
| repl | 0.9 | 0.000000 | 0.000000 | +0.000000 | 0.000000 | +0.000000 |

## OpenTelemetry Demo — AINA 2026

| p_fail | Sheaft all-block | Published all-block | Δ all-block | Sheaft async | Published async | Δ async | Published live |
|---:|---:|---:|---:|---:|---:|---:|---:|
| 0.1 | 0.781676 | 0.781000 | +0.000676 | 0.781676 | 0.781000 | +0.000676 | 0.683000 |
| 0.3 | 0.609860 | 0.610000 | -0.000140 | 0.609860 | 0.610000 | -0.000140 | 0.557000 |
| 0.5 | 0.356375 | 0.356000 | +0.000375 | 0.356375 | 0.356000 | +0.000375 | 0.360000 |
| 0.7 | 0.250842 | 0.251000 | -0.000158 | 0.250842 | 0.251000 | -0.000158 | 0.289000 |
| 0.9 | 0.049971 | 0.050000 | -0.000029 | 0.049971 | 0.050000 | -0.000029 | 0.172000 |

Machine-readable row deltas and aggregate error/correlation statistics are in `comparison.json`.
