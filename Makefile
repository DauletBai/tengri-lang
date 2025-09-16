# FILE: Makefile

# Tengri Language Project Makefile

# --- Go toolchain ---
GO = go

# --- C toolchain ---
CC = cc
CFLAGS = -O3 -march=native -Iinternal/aotminic/runtime

# --- Project structure ---
BIN_DIR = .bin
INTERNAL_DIR = internal

# --- Go sources ---
CMD_AOT = ./cmd/tengri-aot

# --- AOT runtime ---
AOT_RUNTIME_C = $(INTERNAL_DIR)/aotminic/runtime/runtime.c

# --- Benchmark sources ---
BENCH_SRC_DIR = benchmarks/src
# Go
GO_SORT_SRC_DIR = $(BENCH_SRC_DIR)/go/sort
GO_FIB_ITER_SRC = $(BENCH_SRC_DIR)/go/fib_iter.go
GO_FIB_REC_SRC = $(BENCH_SRC_DIR)/go/fib_rec.go
# C
C_SORT_SRC = $(BENCH_SRC_DIR)/c/sort.c
C_FIB_ITER_SRC = $(BENCH_SRC_DIR)/c/fib_iter.c
C_FIB_REC_SRC = $(BENCH_SRC_DIR)/c/fib_rec.c
# Rust
RUST_SORT_SRC_DIR = $(BENCH_SRC_DIR)/rust/sort
RUST_FIB_ITER_SRC_DIR = $(BENCH_SRC_DIR)/rust/fib_iter
RUST_FIB_REC_SRC_DIR = $(BENCH_SRC_DIR)/rust/fib_rec
# Tengri
TGR_SORT_QS_SRC = $(BENCH_SRC_DIR)/tengri/sort_cli_qs.tgr
TGR_SORT_MS_SRC = $(BENCH_SRC_DIR)/tengri/sort_cli_ms.tgr
TGR_FIB_ITER_SRC = $(BENCH_SRC_DIR)/tengri/fib_iter_cli.tgr
TGR_FIB_REC_SRC = $(BENCH_SRC_DIR)/tengri/fib_rec_cli.tgr

# --- Binaries ---
BIN_AOT = $(BIN_DIR)/tengri-aot
BIN_GO_SORT = $(BIN_DIR)/sort_go
BIN_GO_FIB_ITER = $(BIN_DIR)/fib_iter_go
BIN_GO_FIB_REC = $(BIN_DIR)/fib_rec_go
BIN_C_SORT = $(BIN_DIR)/sort_c
BIN_C_FIB_ITER = $(BIN_DIR)/fib_iter_c
BIN_C_FIB_REC = $(BIN_DIR)/fib_rec_c
BIN_RS_SORT = $(BIN_DIR)/sort_rs
BIN_RS_FIB_ITER = $(BIN_DIR)/fib_iter_rs
BIN_RS_FIB_REC = $(BIN_DIR)/fib_rec_rs
BIN_TGR_SORT_QS = $(BIN_DIR)/sort_cli_qsort
BIN_TGR_SORT_MS = $(BIN_DIR)/sort_cli_msort
BIN_TGR_FIB_ITER = $(BIN_DIR)/fib_cli
BIN_TGR_FIB_REC = $(BIN_DIR)/fib_rec_cli

# --- Main targets ---
.PHONY: all
all: build

.PHONY: build
build: $(BIN_AOT) go_benches c_benches rust_benches aot_benches
	@echo "[build] All targets built."

.PHONY: clean
clean:
	@echo "[clean] .bin removed."
	@rm -rf $(BIN_DIR)

# --- Go targets ---
$(BIN_AOT):
	@echo "[build] aot -> $@"
	@$(GO) build -o $@ $(CMD_AOT)

.PHONY: go_benches
go_benches: $(BIN_GO_SORT) $(BIN_GO_FIB_ITER) $(BIN_GO_FIB_REC)

$(BIN_GO_SORT):
	@echo "[build] sort_go -> $@"
	@$(GO) build -o $@ ./$(GO_SORT_SRC_DIR)

$(BIN_GO_FIB_ITER):
	@echo "[build] fib_iter_go -> $@"
	@$(GO) build -o $@ $(GO_FIB_ITER_SRC)

$(BIN_GO_FIB_REC):
	@echo "[build] fib_rec_go -> $@"
	@$(GO) build -o $@ $(GO_FIB_REC_SRC)

# --- C targets ---
.PHONY: c_benches
c_benches: $(BIN_C_SORT) $(BIN_C_FIB_ITER) $(BIN_C_FIB_REC)

# CORRECTED: Added the AOT runtime to the linker command for sort_c
$(BIN_C_SORT):
	@echo "[build] sort_c -> $@"
	@$(CC) $(CFLAGS) $(C_SORT_SRC) $(AOT_RUNTIME_C) -o $@

$(BIN_C_FIB_ITER):
	@echo "[c] fib_iter -> $@"
	@$(CC) $(CFLAGS) $(C_FIB_ITER_SRC) $(AOT_RUNTIME_C) -o $@

$(BIN_C_FIB_REC):
	@echo "[c] fib_rec -> $@"
	@$(CC) $(CFLAGS) $(C_FIB_REC_SRC) $(AOT_RUNTIME_C) -o $@

# --- Rust targets ---
.PHONY: rust_benches
rust_benches: $(BIN_RS_SORT) $(BIN_RS_FIB_ITER) $(BIN_RS_FIB_REC)

$(BIN_RS_SORT):
	@echo "[build] sort_rs -> $@"
	@cd $(RUST_SORT_SRC_DIR) && cargo build --release
	@cp $(RUST_SORT_SRC_DIR)/target/release/sort_rs $(BIN_RS_SORT)

$(BIN_RS_FIB_ITER):
	@echo "[rust] fib_iter -> $@"
	@cd $(RUST_FIB_ITER_SRC_DIR) && cargo build --release
	@cp $(RUST_FIB_ITER_SRC_DIR)/target/release/fib_iter $@

$(BIN_RS_FIB_REC):
	@echo "[rust] fib_rec -> $@"
	@cd $(RUST_FIB_REC_SRC_DIR) && cargo build --release
	@cp $(RUST_FIB_REC_SRC_DIR)/target/release/fib_rec $@

# --- Tengri AOT targets ---
.PHONY: aot_benches
aot_benches: $(BIN_TGR_FIB_ITER) $(BIN_TGR_FIB_REC) $(BIN_TGR_SORT_QS) $(BIN_TGR_SORT_MS)

$(BIN_TGR_FIB_ITER): $(BIN_AOT)
	@echo "[aot] fib_iter -> $@"
	@./$(BIN_AOT) -o $@.c $(TGR_FIB_ITER_SRC)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TGR_FIB_REC): $(BIN_AOT)
	@echo "[aot] fib_rec -> $@"
	@./$(BIN_AOT) -o $@.c $(TGR_FIB_REC_SRC)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TGR_SORT_QS): $(BIN_AOT)
	@echo "[aot] sort qsort -> $@"
	@./$(BIN_AOT) -o $@.c $(TGR_SORT_QS_SRC)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TGR_SORT_MS): $(BIN_AOT)
	@echo "[aot] sort mergesort -> $@"
	@./$(BIN_AOT) -o $@.c $(TGR_SORT_MS_SRC)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

# --- Benchmarking & Plotting ---
.PHONY: bench_all
bench_all: build
	@./benchmarks/run.sh $(SIZE) $(BENCH_REPS)

.PHONY: plot_csv
plot_csv:
	@./benchmarks/plot.sh