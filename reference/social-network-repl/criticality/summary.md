# Sheaft Report Summary

- Decision: **report**
- Mode: `report`
- Overall availability: `0.7279`
- Weighted overall availability: `0.7279`
- Cross-profile availability: `0.7279`
- Cross-profile weighted availability: `0.7279`
- Risk score: `0.2721`
- Confidence: `1.00`

## Profiles

- `importance-baseline`: decision=`report`, weighted=`0.7279`, unweighted=`0.7279`, below-threshold=`0`

## Why

- `gate_pass`: all evaluated endpoints, assertions, and aggregates satisfy the configured gate

## Component criticality (Birnbaum importance)

### `published-endpoint-criticality`

Base profile: `importance-baseline`; paired trials: `12000`; confidence: `0.950`; evaluated components: `31`.

Endpoint `compose-post`; sampled baseline availability: `0.4741`.

| Rank | Component | Importance | Confidence interval | If up | If down | Perfect-component gain | Forced-down impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `service/unique-id-service` | 0.5294 | [0.5205, 0.5383] | 0.5294 | 0.0000 | 0.0529 | 0.4765 |
| 2 | `service/user-mention-service` | 0.5291 | [0.5201, 0.5380] | 0.5291 | 0.0000 | 0.0529 | 0.4762 |
| 3 | `service/nginx-web-server` | 0.5277 | [0.5187, 0.5366] | 0.5277 | 0.0000 | 0.0528 | 0.4749 |
| 3 | `service/post-storage-service` | 0.5277 | [0.5187, 0.5366] | 0.5277 | 0.0000 | 0.0528 | 0.4749 |
| 5 | `service/social-graph-service` | 0.5272 | [0.5182, 0.5361] | 0.5272 | 0.0000 | 0.0527 | 0.4745 |
| 6 | `service/url-shorten-service` | 0.5249 | [0.5160, 0.5338] | 0.5249 | 0.0000 | 0.0525 | 0.4724 |
| 7 | `service/user-service` | 0.5242 | [0.5152, 0.5331] | 0.5242 | 0.0000 | 0.0524 | 0.4718 |
| 8 | `service/media-service` | 0.4747 | [0.4657, 0.4836] | 0.4747 | 0.0000 | 0.0005 | 0.4742 |
| 9 | `service/home-timeline-service` | 0.4745 | [0.4656, 0.4834] | 0.4745 | 0.0000 | 0.0005 | 0.4740 |
| 9 | `service/text-service` | 0.4745 | [0.4656, 0.4834] | 0.4745 | 0.0000 | 0.0005 | 0.4740 |
| 11 | `service/user-timeline-service` | 0.4744 | [0.4655, 0.4834] | 0.4744 | 0.0000 | 0.0005 | 0.4739 |
| 12 | `service/compose-post-service` | 0.4743 | [0.4653, 0.4832] | 0.4743 | 0.0000 | 0.0005 | 0.4738 |
| 13 | `edge/compose-post-service\|media-service\|sync\|true` | 0.4741 | [0.4652, 0.4830] | 0.4741 | 0.0000 | 0.0000 | 0.4741 |
| 13 | `edge/compose-post-service\|text-service\|sync\|true` | 0.4741 | [0.4652, 0.4830] | 0.4741 | 0.0000 | 0.0000 | 0.4741 |
| 13 | `edge/compose-post-service\|unique-id-service\|sync\|true` | 0.4741 | [0.4652, 0.4830] | 0.4741 | 0.0000 | 0.0000 | 0.4741 |
| 13 | `edge/nginx-web-server\|compose-post-service\|sync\|true\|protocol=http\|operation=POST\|route=%2Fwrk2-api%2Fpost%2Fcompose` | 0.4741 | [0.4652, 0.4830] | 0.4741 | 0.0000 | 0.0000 | 0.4741 |
| 13 | `edge/text-service\|url-shorten-service\|sync\|true` | 0.4741 | [0.4652, 0.4830] | 0.4741 | 0.0000 | 0.0000 | 0.4741 |
| 13 | `edge/text-service\|user-mention-service\|sync\|true` | 0.4741 | [0.4652, 0.4830] | 0.4741 | 0.0000 | 0.0000 | 0.4741 |
| 19 | `edge/compose-post-service\|home-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/compose-post-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/compose-post-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/compose-post-service\|user-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/home-timeline-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/home-timeline-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/nginx-web-server\|home-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/nginx-web-server\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/nginx-web-server\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/nginx-web-server\|user-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/social-graph-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/user-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |
| 19 | `edge/user-timeline-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.4741 | 0.4741 | 0.0000 | 0.0000 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering. Service importance refers to aggregate service liveness, including replicated services.

Endpoint `home-timeline`; sampled baseline availability: `0.7275`.

| Rank | Component | Importance | Confidence interval | If up | If down | Perfect-component gain | Forced-down impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `service/post-storage-service` | 0.8099 | [0.8028, 0.8168] | 0.8099 | 0.0000 | 0.0810 | 0.7289 |
| 2 | `service/nginx-web-server` | 0.8093 | [0.8022, 0.8163] | 0.8093 | 0.0000 | 0.0809 | 0.7284 |
| 3 | `service/social-graph-service` | 0.8074 | [0.8003, 0.8144] | 0.8074 | 0.0000 | 0.0807 | 0.7267 |
| 4 | `service/home-timeline-service` | 0.7282 | [0.7202, 0.7361] | 0.7282 | 0.0000 | 0.0007 | 0.7275 |
| 5 | `edge/nginx-web-server\|home-timeline-service\|sync\|true` | 0.0003 | [0.0001, 0.0009] | 0.7275 | 0.7272 | 0.0000 | 0.0003 |
| 6 | `edge/compose-post-service\|home-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|media-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|text-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|unique-id-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|user-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/home-timeline-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/home-timeline-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/nginx-web-server\|compose-post-service\|sync\|true\|protocol=http\|operation=POST\|route=%2Fwrk2-api%2Fpost%2Fcompose` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/nginx-web-server\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/nginx-web-server\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/nginx-web-server\|user-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/social-graph-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/text-service\|url-shorten-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/text-service\|user-mention-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/user-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `edge/user-timeline-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `service/compose-post-service` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `service/media-service` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `service/text-service` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `service/unique-id-service` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `service/url-shorten-service` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `service/user-mention-service` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `service/user-service` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |
| 6 | `service/user-timeline-service` | 0.0000 | [0.0000, 0.0003] | 0.7275 | 0.7275 | 0.0000 | 0.0000 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering. Service importance refers to aggregate service liveness, including replicated services.

Endpoint `user-timeline`; sampled baseline availability: `0.8077`.

| Rank | Component | Importance | Confidence interval | If up | If down | Perfect-component gain | Forced-down impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `service/post-storage-service` | 0.8993 | [0.8938, 0.9046] | 0.8993 | 0.0000 | 0.0899 | 0.8094 |
| 2 | `service/nginx-web-server` | 0.8986 | [0.8931, 0.9039] | 0.8986 | 0.0000 | 0.0899 | 0.8087 |
| 3 | `service/user-timeline-service` | 0.8082 | [0.8010, 0.8151] | 0.8082 | 0.0000 | 0.0008 | 0.8074 |
| 4 | `edge/nginx-web-server\|user-timeline-service\|sync\|true` | 0.0004 | [0.0002, 0.0010] | 0.8077 | 0.8073 | 0.0000 | 0.0004 |
| 5 | `edge/compose-post-service\|home-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/compose-post-service\|media-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/compose-post-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/compose-post-service\|text-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/compose-post-service\|unique-id-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/compose-post-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/compose-post-service\|user-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/home-timeline-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/home-timeline-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/nginx-web-server\|compose-post-service\|sync\|true\|protocol=http\|operation=POST\|route=%2Fwrk2-api%2Fpost%2Fcompose` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/nginx-web-server\|home-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/nginx-web-server\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/nginx-web-server\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/social-graph-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/text-service\|url-shorten-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/text-service\|user-mention-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/user-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `edge/user-timeline-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `service/compose-post-service` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `service/home-timeline-service` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `service/media-service` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `service/social-graph-service` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `service/text-service` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `service/unique-id-service` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `service/url-shorten-service` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `service/user-mention-service` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |
| 5 | `service/user-service` | 0.0000 | [0.0000, 0.0003] | 0.8077 | 0.8077 | 0.0000 | 0.0000 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering. Service importance refers to aggregate service liveness, including replicated services.

## Endpoint results

- `importance-baseline` / `compose-post`: availability=`0.4777`, threshold=`0.0100`, delta=`0.4677`, status=`pass`
- `importance-baseline` / `home-timeline`: availability=`0.7273`, threshold=`0.0100`, delta=`0.7173`, status=`pass`
- `importance-baseline` / `nginx-web-server:POST /wrk2-api/post/compose`: availability=`0.8998`, threshold=`0.0100`, delta=`0.8898`, status=`pass`
- `importance-baseline` / `user-timeline`: availability=`0.8069`, threshold=`0.0100`, delta=`0.7969`, status=`pass`
