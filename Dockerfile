# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.26.2
FROM golang:${GO_VERSION}-bookworm AS builder

ARG PROCRUSTES_COMMIT
ARG BERING_COMMIT
ARG SHEAFT_COMMIT

RUN test -n "${PROCRUSTES_COMMIT}" \
    && test -n "${BERING_COMMIT}" \
    && test -n "${SHEAFT_COMMIT}"

RUN set -eux; \
    git init /src/Procrustes; \
    git -C /src/Procrustes remote add origin https://github.com/MB3R-Lab/Procrustes.git; \
    git -C /src/Procrustes fetch --depth 1 origin "${PROCRUSTES_COMMIT}"; \
    git -C /src/Procrustes checkout --detach FETCH_HEAD; \
    git init /src/Bering; \
    git -C /src/Bering remote add origin https://github.com/MB3R-Lab/Bering.git; \
    git -C /src/Bering fetch --depth 1 origin "${BERING_COMMIT}"; \
    git -C /src/Bering checkout --detach FETCH_HEAD; \
    git init /src/Sheaft; \
    git -C /src/Sheaft remote add origin https://github.com/MB3R-Lab/Sheaft.git; \
    git -C /src/Sheaft fetch --depth 1 origin "${SHEAFT_COMMIT}"; \
    git -C /src/Sheaft checkout --detach FETCH_HEAD

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd /src/Procrustes \
    && CGO_ENABLED=0 go build -trimpath -o /out/procrustes ./cmd/procrustes

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd /src/Bering \
    && CGO_ENABLED=0 go build -trimpath -o /out/bering ./cmd/bering

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd /src/Sheaft \
    && CGO_ENABLED=0 go build -trimpath -o /out/sheaft ./cmd/sheaft

COPY go.mod reproduce.go /src/reproduction/
RUN cd /src/reproduction \
    && CGO_ENABLED=0 go build -trimpath -o /out/reproduce .

FROM alpine:3.22

ARG TOOLCHAIN_VERSION
ARG IMAGE_REVISION

LABEL org.opencontainers.image.title="Procrustes-Bering-Sheaft demonstration artifact" \
      org.opencontainers.image.description="Executable ICSE tool-demonstration artifact for the MB3R resilience-analysis toolchain" \
      org.opencontainers.image.source="https://github.com/MB3R-Lab/mb3r-software-complex-repro" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.version="${TOOLCHAIN_VERSION}" \
      org.opencontainers.image.revision="${IMAGE_REVISION}"

RUN addgroup -S artifact && adduser -S -G artifact artifact \
    && mkdir -p /artifact /work \
    && chown artifact:artifact /artifact /work

COPY --from=builder /out/procrustes /usr/local/bin/procrustes
COPY --from=builder /out/bering /usr/local/bin/bering
COPY --from=builder /out/sheaft /usr/local/bin/sheaft
COPY --from=builder /out/reproduce /usr/local/bin/reproduce
COPY --chown=artifact:artifact cases /artifact/cases
COPY --chown=artifact:artifact reference /artifact/reference
COPY --chown=artifact:artifact schemas /artifact/schemas
COPY --chown=artifact:artifact toolchain-lock.json /artifact/toolchain-lock.json

USER artifact
WORKDIR /artifact

ENV REPRODUCTION_ROOT=/artifact \
    PROCRUSTES_BIN=/usr/local/bin/procrustes \
    BERING_BIN=/usr/local/bin/bering \
    SHEAFT_BIN=/usr/local/bin/sheaft

ENTRYPOINT ["/usr/local/bin/reproduce", "--prebuilt", "--work-dir", "/work"]
CMD ["--case", "social-network-repl", "--repeat", "1"]
