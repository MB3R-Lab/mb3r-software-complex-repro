# Sheaft Report Summary

- Decision: **report**
- Mode: `report`
- Overall availability: `0.4927`
- Weighted overall availability: `0.4927`
- Cross-profile availability: `0.4927`
- Cross-profile weighted availability: `0.4927`
- Risk score: `0.5073`
- Confidence: `1.00`

## Profiles

- `archived-p30`: decision=`report`, weighted=`0.4927`, unweighted=`0.4927`, below-threshold=`1`

## Why

- `endpoint_below_threshold` profile=`archived-p30` endpoint=`frontend:POST /api/checkout` delta=`-0.4973`: endpoint "frontend:POST /api/checkout" availability 0.4927 is below threshold 0.9900

## Failure-tolerance sweeps

### `checkout-failure-curve`

- Base profile: `archived-p30`
- Axis: `independent_replica_failure_probability`
- Trials per point: `50000`
- Confidence level: `0.950`
- Endpoint `frontend:POST /api/checkout`: SLO=`0.5000`, status=`crossed`, certification=`certified`, certified-tolerance=`0.1000` (lower=`0.8040`), last-meeting=`0.1000` (availability=`0.8075`), first-violating=`0.3000` (availability=`0.4853`), bracket=`[0.1000, 0.3000]`

## Endpoint results

- `archived-p30` / `frontend:POST /api/checkout`: availability=`0.4927`, threshold=`0.9900`, delta=`-0.4973`, status=`warn`
