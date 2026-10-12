# Sheaft Report Summary

- Decision: **report**
- Mode: `report`
- Overall availability: `0.6151`
- Weighted overall availability: `0.6151`
- Cross-profile availability: `0.6151`
- Cross-profile weighted availability: `0.6151`
- Risk score: `0.3849`
- Confidence: `1.00`

## Profiles

- `importance-baseline`: decision=`report`, weighted=`0.6151`, unweighted=`0.6151`, below-threshold=`0`

## Why

- `gate_pass`: all evaluated endpoints, assertions, and aggregates satisfy the configured gate

## Component criticality (Birnbaum importance)

### `published-endpoint-criticality`

Base profile: `importance-baseline`; paired trials: `12000`; confidence: `0.950`; evaluated components: `31`.

Endpoint `compose-post`; sampled baseline availability: `0.2849`.

| Rank | Component | Importance | Confidence interval | If up | If down | Perfect-component gain | Forced-down impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `service/user-mention-service` | 0.3201 | [0.3118, 0.3285] | 0.3201 | 0.0000 | 0.0320 | 0.2881 |
| 2 | `service/user-service` | 0.3188 | [0.3106, 0.3272] | 0.3188 | 0.0000 | 0.0319 | 0.2870 |
| 3 | `service/unique-id-service` | 0.3172 | [0.3089, 0.3256] | 0.3172 | 0.0000 | 0.0317 | 0.2854 |
| 4 | `service/media-service` | 0.3171 | [0.3088, 0.3255] | 0.3171 | 0.0000 | 0.0317 | 0.2854 |
| 4 | `service/text-service` | 0.3171 | [0.3088, 0.3255] | 0.3171 | 0.0000 | 0.0317 | 0.2854 |
| 6 | `service/nginx-web-server` | 0.3169 | [0.3087, 0.3253] | 0.3169 | 0.0000 | 0.0317 | 0.2852 |
| 7 | `service/home-timeline-service` | 0.3165 | [0.3082, 0.3249] | 0.3165 | 0.0000 | 0.0316 | 0.2848 |
| 7 | `service/social-graph-service` | 0.3165 | [0.3082, 0.3249] | 0.3165 | 0.0000 | 0.0316 | 0.2848 |
| 9 | `service/compose-post-service` | 0.3161 | [0.3078, 0.3245] | 0.3161 | 0.0000 | 0.0316 | 0.2845 |
| 10 | `service/url-shorten-service` | 0.3159 | [0.3077, 0.3243] | 0.3159 | 0.0000 | 0.0316 | 0.2843 |
| 11 | `service/post-storage-service` | 0.3149 | [0.3067, 0.3233] | 0.3149 | 0.0000 | 0.0315 | 0.2834 |
| 12 | `service/user-timeline-service` | 0.3130 | [0.3048, 0.3214] | 0.3130 | 0.0000 | 0.0313 | 0.2817 |
| 13 | `edge/compose-post-service\|media-service\|sync\|true` | 0.2849 | [0.2769, 0.2931] | 0.2849 | 0.0000 | 0.0000 | 0.2849 |
| 13 | `edge/compose-post-service\|text-service\|sync\|true` | 0.2849 | [0.2769, 0.2931] | 0.2849 | 0.0000 | 0.0000 | 0.2849 |
| 13 | `edge/compose-post-service\|unique-id-service\|sync\|true` | 0.2849 | [0.2769, 0.2931] | 0.2849 | 0.0000 | 0.0000 | 0.2849 |
| 13 | `edge/nginx-web-server\|compose-post-service\|sync\|true\|protocol=http\|operation=POST\|route=%2Fwrk2-api%2Fpost%2Fcompose` | 0.2849 | [0.2769, 0.2931] | 0.2849 | 0.0000 | 0.0000 | 0.2849 |
| 13 | `edge/text-service\|url-shorten-service\|sync\|true` | 0.2849 | [0.2769, 0.2931] | 0.2849 | 0.0000 | 0.0000 | 0.2849 |
| 13 | `edge/text-service\|user-mention-service\|sync\|true` | 0.2849 | [0.2769, 0.2931] | 0.2849 | 0.0000 | 0.0000 | 0.2849 |
| 19 | `edge/compose-post-service\|home-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/compose-post-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/compose-post-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/compose-post-service\|user-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/home-timeline-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/home-timeline-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/nginx-web-server\|home-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/nginx-web-server\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/nginx-web-server\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/nginx-web-server\|user-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/social-graph-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/user-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |
| 19 | `edge/user-timeline-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.2849 | 0.2849 | 0.0000 | 0.0000 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering. Service importance refers to aggregate service liveness, including replicated services.

Endpoint `home-timeline`; sampled baseline availability: `0.6636`.

| Rank | Component | Importance | Confidence interval | If up | If down | Perfect-component gain | Forced-down impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `service/nginx-web-server` | 0.7381 | [0.7301, 0.7459] | 0.7381 | 0.0000 | 0.0738 | 0.6643 |
| 2 | `service/home-timeline-service` | 0.7362 | [0.7282, 0.7440] | 0.7362 | 0.0000 | 0.0736 | 0.6625 |
| 3 | `service/social-graph-service` | 0.7349 | [0.7269, 0.7427] | 0.7349 | 0.0000 | 0.0735 | 0.6614 |
| 4 | `service/post-storage-service` | 0.7319 | [0.7239, 0.7398] | 0.7319 | 0.0000 | 0.0732 | 0.6587 |
| 5 | `edge/nginx-web-server\|home-timeline-service\|sync\|true` | 0.0632 | [0.0590, 0.0677] | 0.6636 | 0.6004 | 0.0000 | 0.0632 |
| 6 | `edge/home-timeline-service\|post-storage-service\|sync\|true` | 0.0062 | [0.0049, 0.0077] | 0.6636 | 0.6574 | 0.0000 | 0.0062 |
| 7 | `edge/compose-post-service\|home-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/compose-post-service\|media-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/compose-post-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/compose-post-service\|text-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/compose-post-service\|unique-id-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/compose-post-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/compose-post-service\|user-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/home-timeline-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/nginx-web-server\|compose-post-service\|sync\|true\|protocol=http\|operation=POST\|route=%2Fwrk2-api%2Fpost%2Fcompose` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/nginx-web-server\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/nginx-web-server\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/nginx-web-server\|user-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/social-graph-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/text-service\|url-shorten-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/text-service\|user-mention-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/user-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `edge/user-timeline-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `service/compose-post-service` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `service/media-service` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `service/text-service` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `service/unique-id-service` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `service/url-shorten-service` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `service/user-mention-service` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `service/user-service` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |
| 7 | `service/user-timeline-service` | 0.0000 | [0.0000, 0.0003] | 0.6636 | 0.6636 | 0.0000 | 0.0000 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering. Service importance refers to aggregate service liveness, including replicated services.

Endpoint `user-timeline`; sampled baseline availability: `0.7367`.

| Rank | Component | Importance | Confidence interval | If up | If down | Perfect-component gain | Forced-down impact |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `service/user-timeline-service` | 0.8163 | [0.8092, 0.8231] | 0.8163 | 0.0000 | 0.0816 | 0.7346 |
| 2 | `service/nginx-web-server` | 0.8157 | [0.8086, 0.8225] | 0.8157 | 0.0000 | 0.0816 | 0.7341 |
| 3 | `service/post-storage-service` | 0.8120 | [0.8049, 0.8189] | 0.8120 | 0.0000 | 0.0812 | 0.7308 |
| 4 | `edge/nginx-web-server\|user-timeline-service\|sync\|true` | 0.0684 | [0.0640, 0.0731] | 0.7367 | 0.6683 | 0.0000 | 0.0684 |
| 5 | `edge/user-timeline-service\|post-storage-service\|sync\|true` | 0.0065 | [0.0052, 0.0081] | 0.7367 | 0.7302 | 0.0000 | 0.0065 |
| 6 | `edge/compose-post-service\|home-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|media-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|text-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|unique-id-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/compose-post-service\|user-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/home-timeline-service\|post-storage-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/home-timeline-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/nginx-web-server\|compose-post-service\|sync\|true\|protocol=http\|operation=POST\|route=%2Fwrk2-api%2Fpost%2Fcompose` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/nginx-web-server\|home-timeline-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/nginx-web-server\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/nginx-web-server\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/social-graph-service\|user-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/text-service\|url-shorten-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/text-service\|user-mention-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `edge/user-service\|social-graph-service\|sync\|true` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `service/compose-post-service` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `service/home-timeline-service` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `service/media-service` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `service/social-graph-service` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `service/text-service` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `service/unique-id-service` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `service/url-shorten-service` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `service/user-mention-service` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |
| 6 | `service/user-service` | 0.0000 | [0.0000, 0.0003] | 0.7367 | 0.7367 | 0.0000 | 0.0000 |

Ranks use point estimates; overlapping intervals do not establish a reliable ordering. Service importance refers to aggregate service liveness, including replicated services.

## Endpoint results

- `importance-baseline` / `compose-post`: availability=`0.2825`, threshold=`0.0100`, delta=`0.2725`, status=`pass`
- `importance-baseline` / `home-timeline`: availability=`0.6446`, threshold=`0.0100`, delta=`0.6346`, status=`pass`
- `importance-baseline` / `nginx-web-server:POST /wrk2-api/post/compose`: availability=`0.8067`, threshold=`0.0100`, delta=`0.7967`, status=`pass`
- `importance-baseline` / `user-timeline`: availability=`0.7264`, threshold=`0.0100`, delta=`0.7164`, status=`pass`
