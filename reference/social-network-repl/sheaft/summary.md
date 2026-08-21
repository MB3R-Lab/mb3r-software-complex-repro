# Sheaft Report Summary

- Decision: **report**
- Mode: `report`
- Overall availability: `0.6321`
- Weighted overall availability: `0.6301`
- Cross-profile availability: `0.2441`
- Cross-profile weighted availability: `0.2127`
- Risk score: `0.3699`
- Confidence: `1.00`

## Profiles

- `p10`: decision=`report`, weighted=`0.6301`, unweighted=`0.6321`, below-threshold=`4`
- `p30`: decision=`report`, weighted=`0.3065`, unweighted=`0.3539`, below-threshold=`4`
- `p50`: decision=`report`, weighted=`0.1139`, unweighted=`0.1823`, below-threshold=`4`
- `p70`: decision=`report`, weighted=`0.0129`, unweighted=`0.0489`, below-threshold=`4`
- `p90`: decision=`report`, weighted=`0.0000`, unweighted=`0.0032`, below-threshold=`4`

## Why

- `endpoint_below_threshold` profile=`p10` endpoint=`compose-post` delta=`-0.6974`: endpoint "compose-post" availability 0.2926 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p10` endpoint=`home-timeline` delta=`-0.3594`: endpoint "home-timeline" availability 0.6306 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p10` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.1261`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.8639 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p10` endpoint=`user-timeline` delta=`-0.2485`: endpoint "user-timeline" availability 0.7415 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`compose-post` delta=`-0.9659`: endpoint "compose-post" availability 0.0241 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`home-timeline` delta=`-0.7040`: endpoint "home-timeline" availability 0.2860 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.3262`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.6638 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`user-timeline` delta=`-0.5484`: endpoint "user-timeline" availability 0.4416 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`compose-post` delta=`-0.9900`: endpoint "compose-post" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`home-timeline` delta=`-0.9019`: endpoint "home-timeline" availability 0.0881 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.5524`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.4376 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`user-timeline` delta=`-0.7865`: endpoint "user-timeline" availability 0.2035 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`compose-post` delta=`-0.9900`: endpoint "compose-post" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`home-timeline` delta=`-0.9848`: endpoint "home-timeline" availability 0.0052 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.8321`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.1579 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`user-timeline` delta=`-0.9574`: endpoint "user-timeline" availability 0.0326 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`compose-post` delta=`-0.9900`: endpoint "compose-post" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`home-timeline` delta=`-0.9900`: endpoint "home-timeline" availability 0.0000 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`nginx-web-server:POST /wrk2-api/post/compose` delta=`-0.9770`: endpoint "nginx-web-server:POST /wrk2-api/post/compose" availability 0.0130 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`user-timeline` delta=`-0.9900`: endpoint "user-timeline" availability 0.0000 is below threshold 0.9900

## Endpoint results

- `p10` / `compose-post`: availability=`0.2926`, threshold=`0.9900`, delta=`-0.6974`, status=`warn`
- `p10` / `home-timeline`: availability=`0.6306`, threshold=`0.9900`, delta=`-0.3594`, status=`warn`
- `p10` / `nginx-web-server:POST /wrk2-api/post/compose`: availability=`0.8639`, threshold=`0.9900`, delta=`-0.1261`, status=`warn`
- `p10` / `user-timeline`: availability=`0.7415`, threshold=`0.9900`, delta=`-0.2485`, status=`warn`
