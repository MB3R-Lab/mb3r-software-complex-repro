FROM golang:1.26.2-bookworm AS builder
WORKDIR /src
COPY go.mod ./
COPY cmd/operation-probe ./cmd/operation-probe
RUN CGO_ENABLED=0 go build -o /proxy ./cmd/operation-probe
FROM scratch
COPY --from=builder /proxy /proxy
ENTRYPOINT ["/proxy", "proxy"]
