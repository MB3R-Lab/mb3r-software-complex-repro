# Sheaft Report Summary

- Decision: **report**
- Mode: `report`
- Overall availability: `0.8374`
- Weighted overall availability: `0.8374`
- Cross-profile availability: `0.8374`
- Cross-profile weighted availability: `0.8374`
- Risk score: `0.1626`
- Confidence: `1.00`

## Profiles

- `importance-baseline`: decision=`report`, weighted=`0.8374`, unweighted=`0.8374`, below-threshold=`0`

## Why

- `gate_pass`: all evaluated endpoints, assertions, and aggregates satisfy the configured gate

## Component criticality (Birnbaum importance)

### `published-endpoint-criticality`

Base profile: `importance-baseline`; paired trials: `12000`; confidence: `0.950`; evaluated components: `39`.

Endpoint `frontend:GET /api/cart`; sampled baseline availability: `0.9032`.

| Rank | Component | Importance | Confidence interval | If up | If down | Perfect-component gain | Forced-down impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `service/cart` | 1.0000 | [0.9997, 1.0000] | 1.0000 | 0.0000 | 0.1000 | 0.9000 |
| 2 | `service/frontend` | 0.9032 | [0.8978, 0.9084] | 0.9032 | 0.0000 | 0.0000 | 0.9032 |
| 3 | `edge/frontend\|cart\|sync\|true` | 0.0905 | [0.0855, 0.0958] | 0.9032 | 0.8127 | 0.0000 | 0.0905 |
| 4 | `edge/ad\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/cart\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|cart\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|currency\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|email\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|kafka\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|payment\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|shipping\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/fraud-detection\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|ad\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|checkout\|sync\|true\|protocol=http\|operation=POST\|route=%2Fapi%2Fcheckout` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|product-reviews\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|recommendation\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/kafka\|accounting\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/kafka\|fraud-detection\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/payment\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/product-reviews\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/recommendation\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/recommendation\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `edge/shipping\|quote\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/accounting` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/ad` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/checkout` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/currency` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/email` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/flagd` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/fraud-detection` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/kafka` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/payment` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/product-catalog` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/product-reviews` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/quote` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/recommendation` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |
| 4 | `service/shipping` | 0.0000 | [0.0000, 0.0003] | 0.9032 | 0.9032 | 0.0000 | 0.0000 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering. Service importance refers to aggregate service liveness, including replicated services.

Endpoint `frontend:GET /api/products`; sampled baseline availability: `0.8973`.

| Rank | Component | Importance | Confidence interval | If up | If down | Perfect-component gain | Forced-down impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `service/product-catalog` | 1.0000 | [0.9997, 1.0000] | 1.0000 | 0.0000 | 0.1000 | 0.9000 |
| 2 | `service/frontend` | 0.8973 | [0.8918, 0.9026] | 0.8973 | 0.0000 | 0.0000 | 0.8973 |
| 3 | `edge/frontend\|product-catalog\|sync\|true` | 0.0096 | [0.0080, 0.0115] | 0.8973 | 0.8878 | 0.0000 | 0.0096 |
| 4 | `edge/ad\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/cart\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|cart\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|currency\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|email\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|kafka\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|payment\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|shipping\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/fraud-detection\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|ad\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|cart\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|checkout\|sync\|true\|protocol=http\|operation=POST\|route=%2Fapi%2Fcheckout` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|product-reviews\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|recommendation\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/kafka\|accounting\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/kafka\|fraud-detection\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/payment\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/product-reviews\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/recommendation\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/recommendation\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `edge/shipping\|quote\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/accounting` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/ad` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/cart` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/checkout` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/currency` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/email` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/flagd` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/fraud-detection` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/kafka` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/payment` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/product-reviews` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/quote` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/recommendation` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |
| 4 | `service/shipping` | 0.0000 | [0.0000, 0.0003] | 0.8973 | 0.8973 | 0.0000 | 0.0000 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering. Service importance refers to aggregate service liveness, including replicated services.

Endpoint `frontend:GET /api/recommendations`; sampled baseline availability: `0.8992`.

| Rank | Component | Importance | Confidence interval | If up | If down | Perfect-component gain | Forced-down impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `service/recommendation` | 1.0000 | [0.9997, 1.0000] | 1.0000 | 0.0000 | 0.1000 | 0.9000 |
| 2 | `edge/frontend\|recommendation\|sync\|true` | 0.8992 | [0.8937, 0.9044] | 0.8992 | 0.0000 | 0.0000 | 0.8992 |
| 2 | `service/frontend` | 0.8992 | [0.8937, 0.9044] | 0.8992 | 0.0000 | 0.0000 | 0.8992 |
| 4 | `edge/ad\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/cart\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|cart\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|currency\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|email\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|kafka\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|payment\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/checkout\|shipping\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/fraud-detection\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|ad\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|cart\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|checkout\|sync\|true\|protocol=http\|operation=POST\|route=%2Fapi%2Fcheckout` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/frontend\|product-reviews\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/kafka\|accounting\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/kafka\|fraud-detection\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/payment\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/product-reviews\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/recommendation\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/recommendation\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `edge/shipping\|quote\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/accounting` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/ad` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/cart` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/checkout` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/currency` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/email` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/flagd` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/fraud-detection` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/kafka` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/payment` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/product-catalog` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/product-reviews` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/quote` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |
| 4 | `service/shipping` | 0.0000 | [0.0000, 0.0003] | 0.8992 | 0.8992 | 0.0000 | 0.0000 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering. Service importance refers to aggregate service liveness, including replicated services.

Endpoint `frontend:POST /api/checkout`; sampled baseline availability: `0.6573`.

| Rank | Component | Importance | Confidence interval | If up | If down | Perfect-component gain | Forced-down impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `service/shipping` | 0.7320 | [0.7240, 0.7398] | 0.7320 | 0.0000 | 0.0732 | 0.6588 |
| 2 | `service/checkout` | 0.7302 | [0.7222, 0.7381] | 0.7302 | 0.0000 | 0.0730 | 0.6572 |
| 3 | `service/payment` | 0.7296 | [0.7216, 0.7375] | 0.7296 | 0.0000 | 0.0730 | 0.6566 |
| 4 | `service/cart` | 0.7272 | [0.7191, 0.7351] | 0.7272 | 0.0000 | 0.0727 | 0.6544 |
| 5 | `edge/checkout\|payment\|sync\|true` | 0.6573 | [0.6488, 0.6658] | 0.6573 | 0.0000 | 0.0000 | 0.6573 |
| 5 | `edge/checkout\|shipping\|sync\|true` | 0.6573 | [0.6488, 0.6658] | 0.6573 | 0.0000 | 0.0000 | 0.6573 |
| 5 | `edge/frontend\|checkout\|sync\|true\|protocol=http\|operation=POST\|route=%2Fapi%2Fcheckout` | 0.6573 | [0.6488, 0.6658] | 0.6573 | 0.0000 | 0.0000 | 0.6573 |
| 5 | `service/frontend` | 0.6573 | [0.6488, 0.6658] | 0.6573 | 0.0000 | 0.0000 | 0.6573 |
| 9 | `edge/ad\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/cart\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/checkout\|cart\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/checkout\|currency\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/checkout\|email\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/checkout\|kafka\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/checkout\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/fraud-detection\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/frontend\|ad\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/frontend\|cart\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/frontend\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/frontend\|product-reviews\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/frontend\|recommendation\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/kafka\|accounting\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/kafka\|fraud-detection\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/payment\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/product-reviews\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/recommendation\|flagd\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/recommendation\|product-catalog\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `edge/shipping\|quote\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/accounting` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/ad` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/currency` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/email` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/flagd` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/fraud-detection` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/kafka` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/product-catalog` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/product-reviews` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/quote` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |
| 9 | `service/recommendation` | 0.0000 | [0.0000, 0.0003] | 0.6573 | 0.6573 | 0.0000 | 0.0000 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering. Service importance refers to aggregate service liveness, including replicated services.

## Endpoint results

- `importance-baseline` / `frontend:GET /api/cart`: availability=`0.9000`, threshold=`0.0100`, delta=`0.8900`, status=`pass`
- `importance-baseline` / `frontend:GET /api/products`: availability=`0.8976`, threshold=`0.0100`, delta=`0.8876`, status=`pass`
- `importance-baseline` / `frontend:GET /api/recommendations`: availability=`0.8976`, threshold=`0.0100`, delta=`0.8876`, status=`pass`
- `importance-baseline` / `frontend:POST /api/checkout`: availability=`0.6544`, threshold=`0.0100`, delta=`0.6444`, status=`pass`
