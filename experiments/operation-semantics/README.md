# Recommendations execution semantics

This is a new, targeted experiment on application commit
`be466eb395294e0b5c2f91f7e76bed475953eaa6`, not a reconstruction of the application
revision used by either archived publication. Run it in the GitHub workflow
**Validate recommendations execution semantics**.

Predeclared request: `GET /api/recommendations?currencyCode=USD`, with no excluded
product IDs, the pinned PostgreSQL seed data, and `recommendationCacheFailure=false`.
Success means HTTP 200 with a nonempty JSON array of products, each with a nonempty
ID. Empty or degraded responses do not satisfy this scenario's requirement.
PostgreSQL, flagd, tracing infrastructure, and frontend are held healthy; the
failure population for analysis consists of recommendation and product-catalog
replicas and the three application RPC channels.

The conditions are healthy execution; stopped recommendation; stopped catalog;
disabled frontend-to-recommendation; disabled recommendation-to-catalog; disabled
frontend-to-catalog. Each fault is applied for the entire request, including all
RPCs and reconnects. A recovery request is required between conditions. Two
recorded requests are collected per condition. An edge is a channel state shared
by repeated RPCs within that request, not a fresh independent error for every RPC.

The proxy forwards bytes without rewriting RPCs or generating traces. Disabling
a channel closes existing connections and rejects new ones. The workflow stores
HTTP responses, actual Jaeger traces, normalized spans for Bering, container
state, resolved deployment configuration, and application provenance. The
application handlers are not replaced by mocks. A timeout is inconclusive and
fails the experiment rather than being counted as confirmed application failure.

The semantic acceptance criterion is equality between observed success and the
declared predicate for every predeclared condition. These observations do not
estimate production failure probabilities or establish universal equivalence.
