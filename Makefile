# --- Makefile (clean) ---

.PHONY: bench bench-fast bench-plot tools deps

TOOLS_DIR := tools/benchfast
BENCH     := $(TOOLS_DIR)/main.go

bench: bench-fast

bench-fast:
	@echo ">> benchfast (csv only)"
	@go run $(BENCH)

bench-plot:
	@echo ">> benchfast + plots"
	@go run $(BENCH) -plot

deps:
	@echo ">> deps (gonum/plot)"
	@go get gonum.org/v1/plot@latest
	@go mod tidy