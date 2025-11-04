.PHONY: lint
lint:
	@echo "golangci-lint recomendado; configure .golangci.yml e instale o binário"

.PHONY: fmt
fmt:
	$(GO_ENV) gofmt -w $(shell go list -f '{{.Dir}}' ./...)
