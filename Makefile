GO ?= go
DOCKER_COMPOSE ?= docker compose

.PHONY: demo failure-walkthroughs reproduce reproduce-paper reproduce-paper-failures reproduce-paper-update verify-runner

demo:
	$(DOCKER_COMPOSE) run --rm demo

failure-walkthroughs:
	$(DOCKER_COMPOSE) run --rm failures

reproduce:
	$(DOCKER_COMPOSE) run --rm reproduce

reproduce-paper:
	$(GO) run ./reproduce.go

reproduce-paper-failures:
	$(GO) run ./reproduce.go --failure-walkthroughs

reproduce-paper-update:
	$(GO) run ./reproduce.go --update-reference

verify-runner:
	$(GO) test ./...
