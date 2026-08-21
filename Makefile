GO ?= go

.PHONY: reproduce-paper reproduce-paper-update verify-runner

reproduce-paper:
	$(GO) run ./reproduce.go

reproduce-paper-update:
	$(GO) run ./reproduce.go --update-reference

verify-runner:
	$(GO) test ./reproduce.go
