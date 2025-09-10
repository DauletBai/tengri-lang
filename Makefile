# ─────────────────────────────────────────────────────────────
# Tengri-Lang — unified build & benchmark Makefile
# ─────────────────────────────────────────────────────────────

SHELL := /bin/sh

# Directories
BIN_DIR        := .bin
RUNTIME_DIR    := internal/aotminic/runtime
BENCH_DIR      := benchmarks
BENCH_SRC_DIR  := $(BENCH_DIR)/src
BENCH_RUNS_DIR := $(BENCH_DIR)/runs
BENCH_LATEST   := $(BENCH_DIR)/latest

# Binaries
BIN_AOT        := $(BIN_DIR)/tengri-aot
BIN_VM         := $(BIN_DIR)/vm
BIN_FIB_REC_GO := $(BIN_DIR)/fib_rec_go
BIN_FIB_ITER_GO:= $(BIN_DIR)/fib_iter_go

# AOT outputs (C/ELF)
BIN_AOT_ITER_C := $(BIN_DIR)/fib_cli.c
BIN_AOT_ITER   := $(BIN_DIR)/fib_cli
BIN_AOT_REC_C  := $(BIN_DIR)/fib_rec_cli.c
BIN_AOT_REC    := $(BIN_DIR)/fib_rec_cli

# Sources
FIB_REC_GO     := $(BENCH_SRC_DIR)/fib_rec/fibonacci.go
FIB_REC_PY     := $(BENCH_SRC_DIR)/fib_rec/fibonacci.py
FIB_REC_TGR    := $(BENCH_SRC_DIR)/fib_rec/tengri/fib_rec_cli.tgr

FIB_ITER_GO    := $(BENCH_SRC_DIR)/fib_iter/go/fibonacci_iter.go
FIB_ITER_PY    := $(BENCH_SRC_DIR)/fib_iter/python/fibonacci_iter.py
FIB_ITER_TGR   := $(BENCH_SRC_DIR)/fib_iter/tengri/fib_iter_cli.tgr

# Compiler tools
GO     ?= go
PY     ?= python3
CC     ?= clang

# Default target
.PHONY: all
all: build

# ─────────────────────────────────────────────────────────────
# Build binaries
# ─────────────────────────────────────────────────────────────

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

$(BIN_AOT): | $(BIN_DIR)
	@echo "[build] aot -> $(BIN_AOT)"
	@$(GO) build -o $(BIN_AOT) ./cmd/tengri-aot

$(BIN_VM): | $(BIN_DIR)
	@echo "[build] vm -> $(BIN_VM)"
	@$(GO) build -o $(BIN_VM) ./cmd/tengri-vm

$(BIN_FIB_REC_GO): | $(BIN_DIR)
	@echo "[build] go fib_rec -> $(BIN_FIB_REC_GO)"
	@$(GO) build -o $(BIN_FIB_REC_GO) $(FIB_REC_GO)

$(BIN_FIB_ITER_GO): | $(BIN_DIR)
	@echo "[build] go fib_iter -> $(BIN_FIB_ITER_GO)"
	@$(GO) build -tags=iter -o $(BIN_FIB_ITER_GO) $(FIB_ITER_GO)

.PHONY: build
build: $(BIN_AOT) $(BIN_VM) $(BIN_FIB_REC_GO) $(BIN_FIB_ITER_GO)
	@echo "[build] done."

# ─────────────────────────────────────────────────────────────
# AOT examples (transpile to C, then link with runtime)
# ─────────────────────────────────────────────────────────────

# NOTE: flag order matters for Go's flag package. Put flags before the source.
# Was:  $(BIN_AOT) $(FIB_ITER_TGR) -o $(BIN_DIR)/fib_cli.c  (WRONG)
# Now:  $(BIN_AOT) -o $(BIN_DIR)/fib_cli.c $(FIB_ITER_TGR) (CORRECT)
$(BIN_AOT_ITER): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(FIB_ITER_TGR) | $(BIN_DIR)
	@echo "[aot] transpile+link iterative -> $(BIN_AOT_ITER)"
	@$(BIN_AOT) -o $(BIN_AOT_ITER_C) $(FIB_ITER_TGR)
	@$(CC) -O2 -I$(RUNTIME_DIR) -o $(BIN_AOT_ITER) $(BIN_AOT_ITER_C) $(RUNTIME_DIR)/runtime.c

$(BIN_AOT_REC): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(FIB_REC_TGR) | $(BIN_DIR)
	@echo "[aot] transpile+link recursive -> $(BIN_AOT_REC)"
	@$(BIN_AOT) -o $(BIN_AOT_REC_C) $(FIB_REC_TGR)
	@$(CC) -O2 -I$(RUNTIME_DIR) -o $(BIN_AOT_REC) $(BIN_AOT_REC_C) $(RUNTIME_DIR)/runtime.c

.PHONY: aot-examples
aot-examples: $(BIN_AOT_ITER) $(BIN_AOT_REC)

# ─────────────────────────────────────────────────────────────
# Benchmarks
# ─────────────────────────────────────────────────────────────

.PHONY: bench
bench: build
	@$(GO) run ./cmd/benchfast

.PHONY: bench-rebuild
bench-rebuild:
	@echo "[bench-rebuild] full rebuild, then bench"
	@$(MAKE) -s build
	@$(MAKE) -s aot-examples
	@$(GO) run ./cmd/benchfast

# ─────────────────────────────────────────────────────────────
# Housekeeping
# ─────────────────────────────────────────────────────────────

.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)

.PHONY: veryclean
veryclean: clean
	@rm -rf $(BENCH_DIR)/runs $(BENCH_LATEST)

# Convenience
.PHONY: setup
setup: tidy build aot-examples bench