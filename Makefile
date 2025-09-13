# ---------------------------------------------
# Tengri-Lang — Makefile (portable, macOS/Linux)
# All comments are for repository documentation.
# ---------------------------------------------

# Tools
GO       := go
CC       := cc
CARGO    := cargo

# Directories
BIN_DIR      := .bin
RUNTIME_DIR  := internal/aotminic/runtime
CBENCH_DIR   := c_benches
RS_DIR       := rust_benches

# C flags (override via environment if needed)
CFLAGS   ?= -O3 -std=c11

# Core binaries
BIN_AOT        := $(BIN_DIR)/tengri-aot
BIN_VM         := $(BIN_DIR)/vm
BIN_GO_FIB_IT  := $(BIN_DIR)/fib_iter_go
BIN_GO_FIB_REC := $(BIN_DIR)/fib_rec_go

# AOT demo sources and outputs
FIB_ITER_TGR := benchmarks/src/fib_iter/tengri/fib_iter_cli.tgr
FIB_REC_TGR  := benchmarks/src/fib_rec/tengri/fib_rec_cli.tgr

BIN_AOT_ITER := $(BIN_DIR)/fib_cli
BIN_AOT_REC  := $(BIN_DIR)/fib_rec_cli

# C bench binaries
BIN_FIB_ITER_C := $(BIN_DIR)/fib_iter_c
BIN_FIB_REC_C  := $(BIN_DIR)/fib_rec_c

# Rust bench binaries
BIN_FIB_ITER_RS := $(BIN_DIR)/fib_iter_rs
BIN_FIB_REC_RS  := $(BIN_DIR)/fib_rec_rs

# Phony
.PHONY: all build build-all aot-examples cbenches rustbenches bench bench-rebuild clean veryclean

# Default
all: build

# ---------------------------------------------
# Core build
# ---------------------------------------------

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

# AOT CLI
$(BIN_AOT): | $(BIN_DIR)
	@echo "[build] aot -> $(BIN_AOT)"
	$(GO) build -o $(BIN_AOT) ./cmd/tengri-aot

# VM CLI
$(BIN_VM): | $(BIN_DIR)
	@echo "[build] vm -> $(BIN_VM)"
	$(GO) build -o $(BIN_VM) ./cmd/tengri-vm

# Go baseline benches (iter/rec)
$(BIN_GO_FIB_IT): | $(BIN_DIR)
	@echo "[build] go fib_iter -> $(BIN_GO_FIB_IT)"
	$(GO) build -o $(BIN_GO_FIB_IT) ./benchmarks/src/fib_iter/go

$(BIN_GO_FIB_REC): | $(BIN_DIR)
	@echo "[build] go fib_rec -> $(BIN_GO_FIB_REC)"
	$(GO) build -o $(BIN_GO_FIB_REC) ./benchmarks/src/fib_rec

build: $(BIN_AOT) $(BIN_VM) $(BIN_GO_FIB_IT) $(BIN_GO_FIB_REC)
	@echo "[build] done."

# ---------------------------------------------
# AOT demo (transpile + link)
# ---------------------------------------------

# Iterative CLI: transpile to C then link with runtime
$(BIN_AOT_ITER): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(FIB_ITER_TGR) | $(BIN_DIR)
	@echo "[aot] transpile+link iterative -> $(BIN_AOT_ITER)"
	@$(BIN_AOT) -o $(BIN_DIR)/fib_cli.c $(FIB_ITER_TGR)
	@echo "C emitted: $(BIN_DIR)/fib_cli.c"
	$(CC) $(CFLAGS) -I$(RUNTIME_DIR) $(BIN_DIR)/fib_cli.c $(RUNTIME_DIR)/runtime.c -o $(BIN_AOT_ITER)

# Recursive CLI: transpile to C then link with runtime
$(BIN_AOT_REC): $(BIN_AOT) $(RUNTIME_DIR)/runtime.c $(RUNTIME_DIR)/runtime.h $(FIB_REC_TGR) | $(BIN_DIR)
	@echo "[aot] transpile+link recursive -> $(BIN_AOT_REC)"
	@$(BIN_AOT) -o $(BIN_DIR)/fib_rec_cli.c $(FIB_REC_TGR)
	@echo "C emitted: $(BIN_DIR)/fib_rec_cli.c"
	$(CC) $(CFLAGS) -I$(RUNTIME_DIR) $(BIN_DIR)/fib_rec_cli.c $(RUNTIME_DIR)/runtime.c -o $(BIN_AOT_REC)

aot-examples: $(BIN_AOT_ITER) $(BIN_AOT_REC)

# ---------------------------------------------
# C benches
# ---------------------------------------------

$(BIN_FIB_ITER_C): $(CBENCH_DIR)/fib_iter.c $(CBENCH_DIR)/runtime_cbench.h | $(BIN_DIR)
	@echo "[cbench] build fib_iter -> $(BIN_FIB_ITER_C)"
	$(CC) $(CFLAGS) -o $@ $(CBENCH_DIR)/fib_iter.c

$(BIN_FIB_REC_C): $(CBENCH_DIR)/fib_rec.c $(CBENCH_DIR)/runtime_cbench.h | $(BIN_DIR)
	@echo "[cbench] build fib_rec -> $(BIN_FIB_REC_C)"
	$(CC) $(CFLAGS) -o $@ $(CBENCH_DIR)/fib_rec.c

cbenches: $(BIN_FIB_ITER_C) $(BIN_FIB_REC_C)
	@echo "[cbenches] done."

# ---------------------------------------------
# Rust benches
# ---------------------------------------------

rustbenches: | $(BIN_DIR)
	@echo "[rustbench] build fib_iter (release)"
	cd $(RS_DIR)/fib_iter && $(CARGO) build --release
	@cp $(RS_DIR)/fib_iter/target/release/fib_iter $(BIN_FIB_ITER_RS)
	@echo "[rustbench] build fib_rec (release)"
	cd $(RS_DIR)/fib_rec && $(CARGO) build --release
	@cp $(RS_DIR)/fib_rec/target/release/fib_rec $(BIN_FIB_REC_RS)
	@echo "[rustbenches] done."

# ---------------------------------------------
# Aggregate builds
# ---------------------------------------------

build-all: build aot-examples cbenches rustbenches
	@echo "[build-all] done."

# ---------------------------------------------
# Bench orchestration
# ---------------------------------------------

bench:
	@echo "[build] done."
	$(GO) run ./cmd/benchfast/main.go

bench-rebuild:
	@echo "[bench-rebuild] full rebuild, then bench"
	$(MAKE) clean
	$(MAKE) build
	$(MAKE) aot-examples
	$(MAKE) bench

# ---------------------------------------------
# Housekeeping
# ---------------------------------------------

clean:
	@rm -rf $(BIN_DIR)
	@echo "[clean] $(BIN_DIR) removed."

veryclean: clean
	@find benchmarks/latest -type f -name '*.csv' -delete || true
	@echo "[veryclean] benchmark artifacts removed."