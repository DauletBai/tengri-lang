# Makefile (root)

# ----------------------------
# Paths
# ----------------------------
BIN_DIR        := .bin
RUNTIME_DIR    := internal/aotminic/runtime

# Tools (allow override via env)
GO             ?= go
CC             ?= clang

# OS-specific flags (macOS/Linux)
UNAME_S        := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	OS_DEFS   := -D_DARWIN_C_SOURCE
else
	OS_DEFS   := -D_POSIX_C_SOURCE=200809L
endif

CSTD           := -std=c11
CFLAGS         := -O2 $(CSTD) $(OS_DEFS) -I$(RUNTIME_DIR)
LDFLAGS        :=

# ----------------------------
# Binaries
# ----------------------------
BIN_AOT        := $(BIN_DIR)/tengri-aot
BIN_VM         := $(BIN_DIR)/vm
BIN_GO_REC     := $(BIN_DIR)/fib_rec_go
BIN_GO_ITER    := $(BIN_DIR)/fib_iter_go

BIN_AOT_ITER   := $(BIN_DIR)/fib_cli
BIN_AOT_REC    := $(BIN_DIR)/fib_rec_cli

# ----------------------------
# Sources
# ----------------------------
AOT_CMD_PKG    := ./cmd/tengri-aot
VM_CMD_PKG     := ./cmd/tengri-vm
GO_ITER_SRC    := benchmarks/src/fib_iter/go/fibonacci_iter.go
GO_REC_SRC     := benchmarks/src/fib_rec/fibonacci.go

# Tengri programs (transpiler inputs)
FIB_ITER_TGR   := benchmarks/src/fib_iter/tengri/fib_iter_cli.tgr
FIB_REC_TGR    := benchmarks/src/fib_rec/tengri/fib_rec_cli.tgr

# ----------------------------
# Phony
# ----------------------------
.PHONY: all help build aot-examples bench bench-rebuild clean

# ----------------------------
# Default/help
# ----------------------------
all: help

help:
	@echo "Targets:"
	@echo "  build           - build .bin/tengri-aot, .bin/vm, Go bench binaries"
	@echo "  aot-examples    - transpile+link AOT examples: $(BIN_AOT_ITER), $(BIN_AOT_REC)"
	@echo "  bench           - run benchmarks (requires AOT examples built)"
	@echo "  bench-rebuild   - full rebuild of binaries + AOT examples, then run benchmarks"
	@echo "  clean           - remove build artifacts in $(BIN_DIR)"
	@echo ""
	@echo "Env:"
	@echo "  BENCH_REPS      - repetitions inside AOT runners (default: iter=5e6, rec=1)"
	@echo ""
	@echo "Notes:"
	@echo "  TIMING: prefer TIME_NS over wall-clock; fallback to TIME:, then wall-clock"

# ----------------------------
# Build binaries
# ----------------------------
$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

$(BIN_AOT): | $(BIN_DIR)
	@echo "[build] aot -> $@"
	@$(GO) build -o $@ $(AOT_CMD_PKG)

$(BIN_VM): | $(BIN_DIR)
	@echo "[build] vm -> $@"
	@$(GO) build -o $@ $(VM_CMD_PKG)

$(BIN_GO_ITER): $(GO_ITER_SRC) | $(BIN_DIR)
	@echo "[build] go fib_iter -> $@"
	@$(GO) build -o $@ $(GO_ITER_SRC)

$(BIN_GO_REC): $(GO_REC_SRC) | $(BIN_DIR)
	@echo "[build] go fib_rec -> $@"
	@$(GO) build -o $@ $(GO_REC_SRC)

build: $(BIN_AOT) $(BIN_VM) $(BIN_GO_ITER) $(BIN_GO_REC)
	@echo "[build] done."

# ----------------------------
# AOT examples (transpile + link)
# ----------------------------
# Iterative
$(BIN_AOT_ITER): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(FIB_ITER_TGR) | $(BIN_DIR)
	@echo "[aot] transpile+link iterative -> $@"
	@$(BIN_AOT) -o $(BIN_DIR)/fib_cli.c $(FIB_ITER_TGR)
	@echo "C emitted: $(BIN_DIR)/fib_cli.c"
	@$(CC) $(CFLAGS) -o $@ $(BIN_DIR)/fib_cli.c $(RUNTIME_DIR)/runtime.c $(LDFLAGS)

# Recursive
$(BIN_AOT_REC): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(FIB_REC_TGR) | $(BIN_DIR)
	@echo "[aot] transpile+link recursive -> $@"
	@$(BIN_AOT) -o $(BIN_DIR)/fib_rec_cli.c $(FIB_REC_TGR)
	@echo "C emitted: $(BIN_DIR)/fib_rec_cli.c"
	@$(CC) $(CFLAGS) -o $@ $(BIN_DIR)/fib_rec_cli.c $(RUNTIME_DIR)/runtime.c $(LDFLAGS)

aot-examples: $(BIN_AOT_ITER) $(BIN_AOT_REC)

# ----------------------------
# Benchmarks
# ----------------------------
bench:
	@echo "[build] done."
	@$(GO) run ./cmd/benchfast/main.go

bench-rebuild: build aot-examples
	@echo "[bench-rebuild] full rebuild, then bench"
	@$(GO) run ./cmd/benchfast/main.go -rebuild

# ----------------------------
# Clean
# ----------------------------
clean:
	@rm -rf $(BIN_DIR)
	@echo "[clean] $(BIN_DIR) removed."