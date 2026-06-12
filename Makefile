MAIN_PATH = "cmd/pgxcli/main.go"
BUILD_PATH = "bin"
TIMEOUT = 60

.PHONY: build clean run update runc lint test precommit fmt vet test-race test-short test-verbose test-bench test-integration coverage go-mod-tidy

build:
	@mkdir -p $(BUILD_PATH)
	@CGO_ENABLED=0 go build -o $(BUILD_PATH)/app $(MAIN_PATH)
	@echo "✓ Build complete: $(BUILD_PATH)/app"

fmt:
	@echo "Formatting code..."
	@gofmt -w ./cmd ./internal
	@echo "✓ Code formatted"

vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet passed"

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run
	@echo "✓ Lint passed"

# Test targets with different configurations
test-default: ARGS=-race
test-short: ARGS=-short
test-verbose: ARGS=-v -race
test-bench: ARGS=-run=xxxxxMatchNothingxxxxx -bench=.
test-integration: ARGS=-tags=integration
test-default test-short test-verbose test-bench test-integration: test
test-race: test-default

test:
	@echo "Running tests $(ARGS)..."
	@go test -timeout $(TIMEOUT)s $(ARGS) ./...
	@echo "✓ Tests passed"

coverage:
	@echo "Generating coverage report..."
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"
	@go tool cover -func=coverage.out | tail -1

go-mod-tidy:
	@echo "Tidying go.mod..."
	@go mod tidy -compat=1.26
	@go mod verify
	@echo "✓ Go modules verified"

precommit: fmt vet lint test
	@echo "✓ All precommit checks passed"

runc: build
	@./bin/app $(DB)

run:
	@./bin/app $(DB)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_PATH) coverage.out coverage.html
	@go clean -cache -testcache
	@echo "✓ Clean complete"

update:
	@go get -u ./...


DOCS_DIR := ./docs-site
VERSION ?= v0.1.1

docs-init:
	git worktree add $(DOCS_DIR) docs

docs-status:
	cd $(DOCS_DIR) && git status

docs-clean:
	git worktree remove --force $(DOCS_DIR)
	rm -rf $(DOCS_DIR)
	git worktree prune
