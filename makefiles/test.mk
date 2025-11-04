.PHONY: test
test:
	$(GO_ENV) go test ./... -count=1

.PHONY: test-race
test-race:
	$(GO_ENV) go test ./... -count=1 -race
