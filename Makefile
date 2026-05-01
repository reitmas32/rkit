.PHONY: test test-coverage test-coverage-html serve-coverage help clean

# Variables
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html
PORT := 8001
GO := $(shell which go)

# Colores para output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
BLUE := \033[0;34m
NC := \033[0m # No Color

help: ## Mostrar esta ayuda
	@echo "$(GREEN)Base Module - Makefile$(NC)"
	@echo ""
	@echo "Targets disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-25s$(NC) %s\n", $$1, $$2}'

test: ## Ejecutar todos los tests del módulo base
	@echo "$(GREEN)Running tests...$(NC)"
	@$(GO) test ./core/... ./persistence/... -v
	@echo "$(GREEN)✓ Tests completed$(NC)"

test-coverage: ## Ejecutar tests y generar coverage
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@$(GO) test ./core/... ./persistence/... -coverprofile=$(COVERAGE_OUT) -covermode=atomic 2>&1 | grep -v "no such tool" || true
	@if [ -f $(COVERAGE_OUT) ]; then \
		echo ""; \
		echo "$(GREEN)Coverage Summary:$(NC)"; \
		$(GO) tool cover -func=$(COVERAGE_OUT) 2>/dev/null | tail -1 || true; \
		echo "$(GREEN)✓ Coverage report generated: $(COVERAGE_OUT)$(NC)"; \
	else \
		echo "$(RED)Error: Coverage file not generated$(NC)"; \
		exit 1; \
	fi

test-coverage-html: ## Ejecutar tests y generar reporte HTML de coverage
	@echo "$(GREEN)Running tests with coverage and generating HTML report...$(NC)"
	@$(GO) test ./core/... ./persistence/... -coverprofile=$(COVERAGE_OUT) -covermode=atomic 2>&1 | grep -E "(PASS|FAIL|coverage:)" | grep -v "no such tool" || true
	@if [ ! -f $(COVERAGE_OUT) ]; then \
		echo "$(RED)Error: Coverage file not generated$(NC)"; \
		exit 1; \
	fi
	@echo ""
	@echo "$(GREEN)Generating HTML coverage report...$(NC)"
	@$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML) 2>/dev/null || (echo "$(RED)Error generating HTML report$(NC)" && exit 1)
	@echo "$(GREEN)✓ HTML coverage report generated: $(COVERAGE_HTML)$(NC)"
	@echo ""
	@echo "$(YELLOW)Coverage Summary:$(NC)"
	@$(GO) tool cover -func=$(COVERAGE_OUT) 2>/dev/null | tail -1 || true
	@echo ""
	@echo "$(YELLOW)To view the HTML report:$(NC)"
	@echo "$(YELLOW)  - Open $(COVERAGE_HTML) in your browser$(NC)"
	@echo "$(YELLOW)  - Or run: $(BLUE)make serve-coverage$(NC) to serve it on port $(PORT)"
	@echo "$(YELLOW)  Coverage visualization:$(NC)"
	@echo "$(YELLOW)    - $(GREEN)Green$(NC) = covered by tests"
	@echo "$(YELLOW)    - $(RED)Red$(NC) = not covered by tests"

serve-coverage: ## Servir el reporte HTML de coverage en el puerto 8001
	@if [ ! -f $(COVERAGE_HTML) ]; then \
		echo "$(RED)Error: $(COVERAGE_HTML) not found. Run 'make test-coverage-html' first.$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Serving coverage report on http://localhost:$(PORT)$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop the server$(NC)"
	@echo ""
	@cp $(COVERAGE_HTML) index.html
	@sh -c 'trap "rm -f index.html" EXIT INT TERM; \
		python3 -m http.server $(PORT) || \
		python -m SimpleHTTPServer $(PORT) || \
		(echo "$(RED)Error: Python not found. Please install Python to serve the HTML file.$(NC)" && rm -f index.html && exit 1)'

clean: ## Limpiar archivos de coverage generados
	@echo "$(YELLOW)Cleaning coverage files...$(NC)"
	@rm -f $(COVERAGE_OUT) $(COVERAGE_HTML) index.html
	@echo "$(GREEN)✓ Cleaned$(NC)"

release: ## Crear y publicar una nueva versión (uso: make release VERSION=1.2.3)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Error: VERSION is required. Usage: make release VERSION=1.2.3$(NC)"; \
		exit 1; \
	fi
	@./scripts/new-version $(VERSION)
