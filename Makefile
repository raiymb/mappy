# ================================
# 📦 VARIABLES
# ================================
LINTER = golangci-lint
APP_NAME = mappy-backend
MAIN = ./cmd/server/main.go
OUT_DIR = ./bin
BUILD_OUT = $(OUT_DIR)/$(APP_NAME)

# ================================
# 🆘 HELP
# ================================
.DEFAULT_GOAL := help

help: ## Show help for each command
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ================================
# 🧹 FORMATTING & LINTING
# ================================

fmt: ## Format code with gofmt
	@gofmt -l -w .

lint-check: ## Check code using golangci-lint (quiet mode)
	@$(LINTER) run --quiet

lint: ## Run golangci-lint
	@$(LINTER) run

lint-fix: ## Fix linter issues
	@$(LINTER) run --fix

gosec: ## Run gosec for security checks
	@gosec ./...

check: fmt lint ## Run all static checks
	@echo "✅ All checks passed!"

# ================================
# 🚀 RUN / BUILD
# ================================

run: ## Run the application (main.go)
	@go run $(MAIN)

build: ## Build the binary
	@mkdir -p $(OUT_DIR)
	@go build -o $(BUILD_OUT) $(MAIN)
	@echo "✅ Built binary at $(BUILD_OUT)"

start: build ## Run the built binary
	@$(BUILD_OUT)

clean: ## Clean generated files
	@rm -rf $(OUT_DIR)

# ================================
# 🧪 TESTING (optional)
# ================================
test: ## Run all tests
	@go test ./...

