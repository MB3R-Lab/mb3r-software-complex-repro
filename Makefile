GO ?= go
DOCKER_COMPOSE ?= docker compose

.PHONY: demo criticality-walkthrough failure-walkthroughs reproduce reproduce-paper reproduce-paper-criticality reproduce-paper-failures reproduce-paper-update verify-runner

demo:
	$(DOCKER_COMPOSE) run --rm demo

criticality-walkthrough:
	$(DOCKER_COMPOSE) run --rm criticality

failure-walkthroughs:
	$(DOCKER_COMPOSE) run --rm failures

reproduce:
	$(DOCKER_COMPOSE) run --rm reproduce

reproduce-paper:
	$(GO) run ./reproduce.go

reproduce-paper-criticality:
	$(GO) run ./reproduce.go --criticality-walkthrough

reproduce-paper-failures:
	$(GO) run ./reproduce.go --failure-walkthroughs

reproduce-paper-update:
	$(GO) run ./reproduce.go --update-reference

verify-runner:
	$(GO) test ./...
