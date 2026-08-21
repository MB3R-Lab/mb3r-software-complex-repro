# Sheaft Report Summary

- Decision: **report**
- Mode: `report`
- Overall availability: `0.7817`
- Weighted overall availability: `0.7817`
- Cross-profile availability: `0.4097`
- Cross-profile weighted availability: `0.4097`
- Risk score: `0.2183`
- Confidence: `1.00`

## Profiles

- `p10`: decision=`report`, weighted=`0.7817`, unweighted=`0.7817`, below-threshold=`4`
- `p30`: decision=`report`, weighted=`0.6099`, unweighted=`0.6099`, below-threshold=`4`
- `p50`: decision=`report`, weighted=`0.3564`, unweighted=`0.3564`, below-threshold=`4`
- `p70`: decision=`report`, weighted=`0.2508`, unweighted=`0.2508`, below-threshold=`4`
- `p90`: decision=`report`, weighted=`0.0500`, unweighted=`0.0500`, below-threshold=`4`

## Why

- `endpoint_below_threshold` profile=`p10` endpoint=`frontend:GET /api/cart` delta=`-0.1224`: endpoint "frontend:GET /api/cart" availability 0.8676 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p10` endpoint=`frontend:GET /api/products` delta=`-0.1227`: endpoint "frontend:GET /api/products" availability 0.8673 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p10` endpoint=`frontend:GET /api/recommendations` delta=`-0.1230`: endpoint "frontend:GET /api/recommendations" availability 0.8670 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p10` endpoint=`frontend:POST /api/checkout` delta=`-0.4652`: endpoint "frontend:POST /api/checkout" availability 0.5248 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`frontend:GET /api/cart` delta=`-0.2576`: endpoint "frontend:GET /api/cart" availability 0.7324 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`frontend:GET /api/products` delta=`-0.2567`: endpoint "frontend:GET /api/products" availability 0.7333 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`frontend:GET /api/recommendations` delta=`-0.2582`: endpoint "frontend:GET /api/recommendations" availability 0.7318 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p30` endpoint=`frontend:POST /api/checkout` delta=`-0.7480`: endpoint "frontend:POST /api/checkout" availability 0.2420 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`frontend:GET /api/cart` delta=`-0.5236`: endpoint "frontend:GET /api/cart" availability 0.4664 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`frontend:GET /api/products` delta=`-0.5225`: endpoint "frontend:GET /api/products" availability 0.4675 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`frontend:GET /api/recommendations` delta=`-0.5241`: endpoint "frontend:GET /api/recommendations" availability 0.4659 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p50` endpoint=`frontend:POST /api/checkout` delta=`-0.9643`: endpoint "frontend:POST /api/checkout" availability 0.0257 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`frontend:GET /api/cart` delta=`-0.6569`: endpoint "frontend:GET /api/cart" availability 0.3331 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`frontend:GET /api/products` delta=`-0.6570`: endpoint "frontend:GET /api/products" availability 0.3330 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`frontend:GET /api/recommendations` delta=`-0.6565`: endpoint "frontend:GET /api/recommendations" availability 0.3335 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p70` endpoint=`frontend:POST /api/checkout` delta=`-0.9863`: endpoint "frontend:POST /api/checkout" availability 0.0037 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`frontend:GET /api/cart` delta=`-0.9235`: endpoint "frontend:GET /api/cart" availability 0.0665 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`frontend:GET /api/products` delta=`-0.9228`: endpoint "frontend:GET /api/products" availability 0.0672 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`frontend:GET /api/recommendations` delta=`-0.9238`: endpoint "frontend:GET /api/recommendations" availability 0.0662 is below threshold 0.9900
- `endpoint_below_threshold` profile=`p90` endpoint=`frontend:POST /api/checkout` delta=`-0.9900`: endpoint "frontend:POST /api/checkout" availability 0.0000 is below threshold 0.9900

## Endpoint results

- `p10` / `frontend:GET /api/cart`: availability=`0.8676`, threshold=`0.9900`, delta=`-0.1224`, status=`warn`
- `p10` / `frontend:GET /api/products`: availability=`0.8673`, threshold=`0.9900`, delta=`-0.1227`, status=`warn`
- `p10` / `frontend:GET /api/recommendations`: availability=`0.8670`, threshold=`0.9900`, delta=`-0.1230`, status=`warn`
- `p10` / `frontend:POST /api/checkout`: availability=`0.5248`, threshold=`0.9900`, delta=`-0.4652`, status=`warn`
