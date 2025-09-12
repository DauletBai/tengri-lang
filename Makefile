# ============================================================
# Tengri-lang — unified Makefile
# ============================================================
# Tools
GO      := go
CLANG   := clang

# Dirs & paths
BIN_DIR        := .bin
RUNTIME_DIR    := internal/aotminic/runtime

# AOT transpiler & examples
BIN_AOT        := $(BIN_DIR)/tengri-aot
AOT_SRC_ITER   := benchmarks/src/fib_iter/tengri/fib_iter_cli.tgr
AOT_SRC_REC    := benchmarks/src/fib_rec/tengri/fib_rec_cli.tgr
BIN_AOT_ITER   := $(BIN_DIR)/fib_cli
BIN_AOT_REC    := $(BIN_DIR)/fib_rec_cli

# VM / Go Fibonacci binaries
BIN_VM         := $(BIN_DIR)/vm
BIN_GO_ITER    := $(BIN_DIR)/fib_iter_go
BIN_GO_REC     := $(BIN_DIR)/fib_rec_go

# Bench tool (Go)
BENCHFAST_MAIN := ./cmd/benchfast/main.go

# C micro/mid benches (compiled with Tengri runtime)
BIN_CALLS      := $(BIN_DIR)/calls_cli
BIN_SIEVE      := $(BIN_DIR)/sieve_cli
BIN_MATMUL     := $(BIN_DIR)/matmul_cli
BIN_SORT       := $(BIN_DIR)/sort_cli

# Default target
.PHONY: all
all: build

# ------------------------------------------------------------
# Utilities
# ------------------------------------------------------------
$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)
	@echo "[clean] .bin removed."

.PHONY: tidy
tidy:
	@$(GO) mod tidy

# ------------------------------------------------------------
# Core builds (AOT transpiler, VM, Go fib)
# ------------------------------------------------------------
.PHONY: build
build: $(BIN_AOT) $(BIN_VM) $(BIN_GO_ITER) $(BIN_GO_REC)
	@echo "[build] done."

$(BIN_AOT): | $(BIN_DIR)
	@echo "[build] aot -> $@"
	@$(GO) build -o $@ ./cmd/tengri-aot

$(BIN_VM): | $(BIN_DIR)
	@echo "[build] vm -> $@"
	@$(GO) build -o $@ ./cmd/tengri-vm

$(BIN_GO_ITER): benchmarks/src/fib_iter/go/fibonacci_iter.go | $(BIN_DIR)
	@echo "[build] go fib_iter -> $@"
	@$(GO) build -o $@ $<

$(BIN_GO_REC): benchmarks/src/fib_rec/fibonacci.go | $(BIN_DIR)
	@echo "[build] go fib_rec -> $@"
	@$(GO) build -o $@ $<

# ------------------------------------------------------------
# AOT example apps (iterative / recursive Fibonacci)
# ------------------------------------------------------------
.PHONY: aot-examples
aot-examples: $(BIN_AOT_ITER) $(BIN_AOT_REC)

$(BIN_AOT_ITER): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(AOT_SRC_ITER) | $(BIN_DIR)
	@echo "[aot] transpile+link iterative -> $@"
	@$(BIN_AOT) -o $(BIN_DIR)/fib_cli.c $(AOT_SRC_ITER)
	@$(CLANG) -O2 -o $@ $(BIN_DIR)/fib_cli.c $(RUNTIME_DIR)/runtime.c -I$(RUNTIME_DIR)

$(BIN_AOT_REC): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(AOT_SRC_REC) | $(BIN_DIR)
	@echo "[aot] transpile+link recursive -> $@"
	@$(BIN_AOT) -o $(BIN_DIR)/fib_rec_cli.c $(AOT_SRC_REC)
	@$(CLANG) -O2 -o $@ $(BIN_DIR)/fib_rec_cli.c $(RUNTIME_DIR)/runtime.c -I$(RUNTIME_DIR)

# ------------------------------------------------------------
# Native C benches (calls / sieve / matmul / sort)
# ------------------------------------------------------------
.PHONY: cbenches
cbenches: $(BIN_CALLS) $(BIN_SIEVE) $(BIN_MATMUL) $(BIN_SORT)
	@echo "[cbenches] done."

$(BIN_CALLS): benchmarks/src/calls/c/calls_cli.c $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h | $(BIN_DIR)
	@echo "[cbench] build calls -> $@"
	@$(CLANG) -O2 -o $@ $< $(RUNTIME_DIR)/runtime.c -I$(RUNTIME_DIR)

$(BIN_SIEVE): benchmarks/src/sieve/c/sieve_cli.c $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h | $(BIN_DIR)
	@echo "[cbench] build sieve -> $@"
	@$(CLANG) -O2 -o $@ $< $(RUNTIME_DIR)/runtime.c -I$(RUNTIME_DIR)

$(BIN_MATMUL): benchmarks/src/matmul/c/matmul_cli.c $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h | $(BIN_DIR)
	@echo "[cbench] build matmul -> $@"
	@$(CLANG) -O2 -o $@ $< $(RUNTIME_DIR)/runtime.c -I$(RUNTIME_DIR)

$(BIN_SORT): benchmarks/src/sort/c/sort_cli.c $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h | $(BIN_DIR)
	@echo "[cbench] build sort -> $@"
	@$(CLANG) -O2 -o $@ $< $(RUNTIME_DIR)/runtime.c -I$(RUNTIME_DIR)

# ------------------------------------------------------------
# Bench runs
# ------------------------------------------------------------
.PHONY: bench
bench: build
	@echo "[build] done."
	@$(GO) run $(BENCHFAST_MAIN)

# Полный цикл: пересборка всего + AOT-примеры + запуск benchfast с -rebuild
.PHONY: bench-rebuild
bench-rebuild:
	@echo "[bench-rebuild] full rebuild, then bench"
	@$(MAKE) --no-print-directory clean
	@$(MAKE) --no-print-directory build
	@$(MAKE) --no-print-directory aot-examples
	@echo "[bench-rebuild] invoking benchfast -rebuild"
	@$(GO) run $(BENCHFAST_MAIN) -rebuild

# Удобная связка: всё собрать (включая C-benches)
.PHONY: build-all
build-all: build aot-examples cbenches
	@echo "[build-all] all artifacts built."

# ------------------------------------------------------------
# Help
# ------------------------------------------------------------
.PHONY: help
help:
	@echo "Targets:"
	@echo "  tidy             - go mod tidy"
	@echo "  build            - build AOT, VM, Go Fibonacci binaries"
	@echo "  aot-examples     - build .bin/fib_cli and .bin/fib_rec_cli"
	@echo "  cbenches         - build native C benches (calls/sieve/matmul/sort)"
	@echo "  build-all        - build everything above"
	@echo "  bench            - run benchfast (uses existing .bin/*)"
	@echo "  bench-rebuild    - clean + rebuild + aot-examples + benchfast -rebuild"
	@echo "  clean            - remove .bin"