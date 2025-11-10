.PHONY: test
# Run all unit tests in the module. You can pass extra flags via TESTFLAGS, e.g.
#   make test TESTFLAGS='-run TestFoo -v'
test:
	@echo "[TEST] running go test ./..."
	$(GO_ENV) go test ./... -count=1 $(TESTFLAGS)

.PHONY: test-race
# Run tests with the race detector enabled
test-race:
	@echo "[TEST] running go test ./... -race"
	$(GO_ENV) go test ./... -count=1 -race $(TESTFLAGS)

.PHONY: test-pkg
# Run tests for a specific package. Set PKG to a package path (default: ./...)
# Example: make test-pkg PKG=./internal/chat -v
test-pkg:
	@echo "[TEST] running go test for package: $(if $(PKG),$(PKG),./...)"
	$(GO_ENV) go test $(if $(PKG),$(PKG),./...) -count=1 $(TESTFLAGS)

.PHONY: coverage
# Run tests and produce a coverage report (coverage.out)
coverage:
	@echo "[TEST] running tests and generating coverage.out"
	$(GO_ENV) go test ./... -coverprofile=coverage.out $(TESTFLAGS)
	@echo "Coverage summary:"
	@go tool cover -func=coverage.out | sed -n '1,200p'
