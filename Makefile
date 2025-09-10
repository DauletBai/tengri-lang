# -------------------------------
# Tengri-lang top-level Makefile
# -------------------------------

SHELL := /bin/sh

BIN_DIR := .bin

# Prebuilt binaries we expect for benchfast
GO_FIB_REC   := $(BIN_DIR)/fib_rec_go
GO_FIB_ITER  := $(BIN_DIR)/fib_iter_go
VM_BIN       := $(BIN_DIR)/vm
AOT_BIN      := $(BIN_DIR)/tengri-aot
AOT_FIB_ITER := $(BIN_DIR)/fib_cli
AOT_FIB_REC  := $(BIN_DIR)/fib_rec_cli

# Sources
SRC_GO_FIB_REC   := benchmarks/src/fib_rec/fibonacci.go
SRC_GO_FIB_ITER  := benchmarks/src/fib_iter/go/fibonacci_iter.go
SRC_PY_FIB_REC   := benchmarks/src/fib_rec/fibonacci.py
SRC_PY_FIB_ITER  := benchmarks/src/fib_iter/python/fibonacci_iter.py
SRC_TGR_FIB_ITER := benchmarks/src/fib_iter/tengri/fib_cli.tgr
SRC_TGR_FIB_REC  := benchmarks/src/fib_rec/tengri/fib_rec_cli.tgr

AOT_RUNTIME_C := internal/aotminic/runtime/runtime.c
CMD_VM_MAIN   := cmd/tengri-vm/main.go
CMD_AOT_MAIN  := cmd/tengri-aot/main.go
CMD_BENCHFAST := cmd/benchfast/main.go

# Utilities
CLANG := clang

.PHONY: all setup build rebuild aot-examples bench bench-rebuild clean tidy

all: bench

setup: tidy build

tidy:
	@echo "[go mod tidy]"
	@go mod tidy

# -------------------------------
# Build (idempotent)
# -------------------------------

build: $(GO_FIB_REC) $(GO_FIB_ITER) opt-vm opt-aot aot-examples
	@echo "[build] done."

# Build Go baselines
$(GO_FIB_REC): $(SRC_GO_FIB_REC)
	@mkdir -p $(BIN_DIR)
	@echo "[build] go fib_rec -> $@"
	@go build -o $@ $<

$(GO_FIB_ITER): $(SRC_GO_FIB_ITER)
	@mkdir -p $(BIN_DIR)
	@echo "[build] go fib_iter -> $@"
	@go build -tags=iter -o $@ $<

# Optional VM build (skip if sources absent)
opt-vm:
	@if [ -f "$(CMD_VM_MAIN)" ]; then \
	  echo "[build] vm -> $(VM_BIN)"; \
	  go build -o $(VM_BIN) $(CMD_VM_MAIN); \
	else \
	  echo "[build] vm skipped (cmd/tengri-vm/main.go not found)"; \
	fi

# Optional AOT transpiler build (skip if sources absent)
opt-aot:
	@if [ -f "$(CMD_AOT_MAIN)" ]; then \
	  echo "[build] aot -> $(AOT_BIN)"; \
	  GO111MODULE=on go build -o $(AOT_BIN) $(CMD_AOT_MAIN); \
	else \
	  echo "[build] aot skipped (cmd/tengri-aot/main.go not found)"; \
	fi

# AOT examples (TGR -> C -> ELF) if AOT is present
aot-examples:
	@if [ -x "$(AOT_BIN)" ] && [ -f "$(AOT_RUNTIME_C)" ]; then \
	  echo "[aot] transpile+link iterative -> $(AOT_FIB_ITER)"; \
	  $(AOT_BIN) $(SRC_TGR_FIB_ITER) -o $(BIN_DIR)/fib_cli.c; \
	  $(CLANG) -O2 -o $(AOT_FIB_ITER) $(BIN_DIR)/fib_cli.c $(AOT_RUNTIME_C); \
	  echo "[aot] transpile+link recursive -> $(AOT_FIB_REC)"; \
	  $(AOT_BIN) $(SRC_TGR_FIB_REC) -o $(BIN_DIR)/fib_rec_cli.c; \
	  $(CLANG) -O2 -o $(AOT_FIB_REC) $(BIN_DIR)/fib_rec_cli.c $(AOT_RUNTIME_C); \
	else \
	  echo "[aot] skipped (no $(AOT_BIN) or no $(AOT_RUNTIME_C))"; \
	fi

# -------------------------------
# Rebuild (force rebuild all .bin)
# -------------------------------

rebuild:
	@echo "[rebuild] removing $(BIN_DIR)"
	@rm -rf $(BIN_DIR)
	@$(MAKE) build

# -------------------------------
# Bench commands
# -------------------------------

bench:
	@echo "Task = fib_rec / fib_iter"
	@echo "TIMING: prefer TIME_NS over wall-clock (enforced in benchfast)"
	@go run $(CMD_BENCHFAST)

bench-rebuild:
	@echo "[bench-rebuild] full rebuild, then bench"
	@$(MAKE) rebuild
	@go run $(CMD_BENCHFAST) -rebuild

# -------------------------------
# Clean
# -------------------------------

clean:
	@echo "[clean] removing .bin and latest CSVs"
	@rm -rf $(BIN_DIR)
	@rm -rf benchmarks/latest/results/*.csv