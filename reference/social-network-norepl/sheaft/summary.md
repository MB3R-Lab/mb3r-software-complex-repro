# Sheaft Report Summary

- Decision: **report**
- Mode: `report`
- Overall availability: `0.4132`
- Weighted overall availability: `0.4186`
- Cross-profile availability: `0.1438`
- Cross-profile weighted availability: `0.1256`
- Risk score: `0.5814`
- Confidence: `1.00`

## Profiles

- `p10`: decision=`report`, weighted=`0.4186`, unweighted=`0.4132`, below-threshold=`4`
- `p30`: decision=`report`, weighted=`0.1625`, unweighted=`0.2060`, below-threshold=`4`
- `p50`: decision=`report`, weighted=`0.0455`, unweighted=`0.0871`, below-threshold=`4`
- `p70`: decision=`report`, weighted=`0.0013`, unweighted=`0.0125`, below-threshold=`4`
- `p90`: decision=`report`, weighted=`0.0000`, unweighted=`0.0000`, below-threshold=`4`

## Why

- `endpoint_below_threshold` profile=`p10` endpoint=`compose-post` delta=`-0.9900`: endpoint "compose-post" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p10` endpoint=`home-timeline` delta=`-0.5659`: endpoint "home-timeline" availability 0.4241 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p10` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.3084`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.6816 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p10` endpoint=`user-timeline` delta=`-0.4427`: endpoint "user-timeline" availability 0.5473 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`compose-post` delta=`-0.9900`: endpoint "compose-post" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`home-timeline` delta=`-0.8473`: endpoint "home-timeline" availability 0.1427 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.5652`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.4248 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`user-timeline` delta=`-0.7337`: endpoint "user-timeline" availability 0.2563 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`compose-post` delta=`-0.9900`: endpoint "compose-post" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`home-timeline` delta=`-0.9596`: endpoint "home-timeline" availability 0.0304 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.7631`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.2269 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`user-timeline` delta=`-0.8990`: endpoint "user-timeline" availability 0.0910 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`compose-post` delta=`-0.9900`: endpoint "compose-post" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`home-timeline` delta=`-0.9900`: endpoint "home-timeline" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.9442`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.0458 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`user-timeline` delta=`-0.9856`: endpoint "user-timeline" availability 0.0044 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`compose-post` delta=`-0.9900`: endpoint "compose-post" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`home-timeline` delta=`-0.9900`: endpoint "home-timeline" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.9900`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`user-timeline` delta=`-0.9900`: endpoint "user-timeline" availability 0.0000 is below threshold 0.9900

## Endpoint results

- `p10` / `compose-post`: availability=`0.0000`, threshold=`0.9900`, delta=`-0.9900`, status=`warn`
- `p10` / `home-timeline`: availability=`0.4241`, threshold=`0.9900`, delta=`-0.5659`, status=`warn`
- `p10` / `nginx-web-server:POST /wrk2-api/post/compose`: availability=`0.6816`, threshold=`0.9900`, delta=`-0.3084`, status=`warn`
- `p10` / `user-timeline`: availability=`0.5473`, threshold=`0.9900`, delta=`-0.4427`, status=`warn`
