GO ?= go
DOCKER_COMPOSE ?= docker compose

.PHONY: demo reproduce reproduce-paper reproduce-paper-update verify-runner

demo:
	$(DOCKER_COMPOSE) run --rm demo

reproduce:
	$(DOCKER_COMPOSE) run --rm reproduce

reproduce-paper:
	$(GO) run ./reproduce.go

reproduce-paper-update:
	$(GO) run ./reproduce.go --update-reference

verify-runner:
	$(GO) test ./...
