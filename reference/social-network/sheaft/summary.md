# Sheaft Report Summary

- Decision: **report**
- Mode: `report`
- Overall availability: `0.3529`
- Weighted overall availability: `0.3059`
- Cross-profile availability: `0.3529`
- Cross-profile weighted availability: `0.3059`
- Risk score: `0.6941`
- Confidence: `1.00`

## Profiles

- `archived-p30`: decision=`report`, weighted=`0.3059`, unweighted=`0.3529`, below-threshold=`4`

## Why

- `endpoint_below_threshold` profile=`archived-p30` endpoint=`compose-post` delta=`-0.9676`: endpoint "compose-post" availability 0.0225 is below threshold 0.9900
- `endpoint_below_threshold` profile=`archived-p30` endpoint=`home-timeline` delta=`-0.7049`: endpoint "home-timeline" availability 0.2851 is below threshold 0.9900
- `endpoint_below_threshold` profile=`archived-p30` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.3278`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.6622 is below threshold 0.9900
- `endpoint_below_threshold` profile=`archived-p30` endpoint=`user-timeline` delta=`-0.5480`: endpoint "user-timeline" availability 0.4420 is below threshold 0.9900

## Failure-tolerance sweeps

### `social-endpoint-failure-curves`

- Base profile: `archived-p30`
- Axis: `independent_replica_failure_probability`
- Trials per point: `20000`
- Confidence level: `0.950`
- Endpoint `compose-post`: SLO=`0.5000`, status=`crossed`, certification=`certified`, certified-tolerance=`0.0000` (lower=`0.9998`), last-meeting=`0.0000` (availability=`1.0000`), first-violating=`0.1000` (availability=`0.4805`), bracket=`[0.0000, 0.1000]`
- Endpoint `home-timeline`: SLO=`0.5000`, status=`crossed`, certification=`certified`, certified-tolerance=`0.1000` (lower=`0.7235`), last-meeting=`0.1000` (availability=`0.7297`), first-violating=`0.3000` (availability=`0.3367`), bracket=`[0.1000, 0.3000]`
- Endpoint `user-timeline`: SLO=`0.5000`, status=`crossed`, certification=`certified`, certified-tolerance=`0.1000` (lower=`0.8022`), last-meeting=`0.1000` (availability=`0.8077`), first-violating=`0.3000` (availability=`0.4760`), bracket=`[0.1000, 0.3000]`

## Endpoint results

- `archived-p30` / `compose-post`: availability=`0.0225`, threshold=`0.9900`, delta=`-0.9676`, status=`warn`
- `archived-p30` / `home-timeline`: availability=`0.2851`, threshold=`0.9900`, delta=`-0.7049`, status=`warn`
- `archived-p30` / `nginx-web-server:POST /wrk2-api/post/compose`: availability=`0.6622`, threshold=`0.9900`, delta=`-0.3278`, status=`warn`
- `archived-p30` / `user-timeline`: availability=`0.4420`, threshold=`0.9900`, delta=`-0.5480`, status=`warn`
