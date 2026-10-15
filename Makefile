GO ?= go
DOCKER_COMPOSE ?= docker compose
BENCHMARK_REPEATS ?= 3

.PHONY: demo criticality-walkthrough failure-walkthroughs decision-walkthrough validation-matrix cost-profile evaluation-suite reproduce reproduce-paper reproduce-paper-criticality reproduce-paper-failures reproduce-paper-decision reproduce-paper-validation reproduce-paper-cost reproduce-paper-evaluation reproduce-paper-update verify-runner

demo:
	$(DOCKER_COMPOSE) run --rm demo

.PHONY: operation-walkthrough reproduce-paper-operation
operation-walkthrough:
	$(DOCKER_COMPOSE) run --rm operation

reproduce-paper-operation:
	$(GO) run . --operation-walkthrough

criticality-walkthrough:
	$(DOCKER_COMPOSE) run --rm criticality

failure-walkthroughs:
	$(DOCKER_COMPOSE) run --rm failures

decision-walkthrough:
	$(DOCKER_COMPOSE) run --rm decision

validation-matrix:
	$(DOCKER_COMPOSE) run --rm validation

cost-profile:
	$(DOCKER_COMPOSE) run --rm cost --benchmark-repeats $(BENCHMARK_REPEATS)

evaluation-suite:
	$(DOCKER_COMPOSE) run --rm evaluation --benchmark-repeats $(BENCHMARK_REPEATS)

reproduce:
	$(DOCKER_COMPOSE) run --rm reproduce

reproduce-paper:
	$(GO) run .

reproduce-paper-criticality:
	$(GO) run . --criticality-walkthrough

reproduce-paper-failures:
	$(GO) run . --failure-walkthroughs

reproduce-paper-decision:
	$(GO) run . --decision-walkthrough

reproduce-paper-validation:
	$(GO) run . --validation-matrix

reproduce-paper-cost:
	$(GO) run . --cost-profile --benchmark-repeats $(BENCHMARK_REPEATS)

reproduce-paper-evaluation:
	$(GO) run . --evaluation-suite --benchmark-repeats $(BENCHMARK_REPEATS)

reproduce-paper-update:
	$(GO) run . --update-reference

verify-runner:
	$(GO) test ./...
