# tenge Language Project Makefile

GO = go
CC = cc
CFLAGS = -O3 -march=native -Iinternal/aotminic/runtime -lm

BIN_DIR        = .bin
BIN_DIR_ABS    = $(abspath $(BIN_DIR))
AOT_RUNTIME_C  = internal/aotminic/runtime/runtime.c

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

CMD_COMPILER = ./cmd/tenge
BIN_COMPILER = $(BIN_DIR)/tenge

BENCH_SRC_DIR = benchmarks/src

# C
C_SORT_SRC      = $(BENCH_SRC_DIR)/c/sort.c
C_FIB_ITER_SRC  = $(BENCH_SRC_DIR)/c/fib_iter.c
C_FIB_REC_SRC   = $(BENCH_SRC_DIR)/c/fib_rec.c
C_VAR_MC_SRC    = $(BENCH_SRC_DIR)/c/var_monte_carlo.c

# Go
GO_SORT_SRC_DIR = $(BENCH_SRC_DIR)/go/sort
GO_FIB_ITER_SRC = $(BENCH_SRC_DIR)/go/fib_iter.go
GO_FIB_REC_SRC  = $(BENCH_SRC_DIR)/go/fib_rec.go
GO_VAR_MC_SRC   = $(BENCH_SRC_DIR)/go/var_mc.go

# Rust
RUST_SORT_SRC_DIR     = $(BENCH_SRC_DIR)/rust/sort
RUST_FIB_ITER_SRC_DIR = $(BENCH_SRC_DIR)/rust/fib_iter
RUST_FIB_REC_SRC_DIR  = $(BENCH_SRC_DIR)/rust/fib_rec
RUST_VAR_MC_SRC_DIR   = $(BENCH_SRC_DIR)/rust/var_mc

# Tenge AOT demos
TNG_SORT_QS_SRC    = $(BENCH_SRC_DIR)/tenge/sort_cli_qs.tng
TNG_SORT_MS_SRC    = $(BENCH_SRC_DIR)/tenge/sort_cli_ms.tng
TNG_FIB_ITER_SRC   = $(BENCH_SRC_DIR)/tenge/fib_iter_cli.tng
TNG_FIB_REC_SRC    = $(BENCH_SRC_DIR)/tenge/fib_rec_cli.tng
TNG_VAR_MC_SORT    = $(BENCH_SRC_DIR)/tenge/var_mc_sort_cli.tng
TNG_VAR_MC_ZIG     = $(BENCH_SRC_DIR)/tenge/var_mc_zig_cli.tng
TNG_VAR_MC_QSEL    = $(BENCH_SRC_DIR)/tenge/var_mc_qsel_cli.tng
TNG_SORT_PDQ_SRC   = $(BENCH_SRC_DIR)/tenge/sort_cli_pdq.tng
TNG_SORT_RADIX_SRC = $(BENCH_SRC_DIR)/tenge/sort_cli_radix.tng

# Binaries
BIN_C_SORT      = $(BIN_DIR)/sort_c
BIN_C_FIB_ITER  = $(BIN_DIR)/fib_iter_c
BIN_C_FIB_REC   = $(BIN_DIR)/fib_rec_c
BIN_C_VAR_MC    = $(BIN_DIR)/var_mc_c

BIN_GO_SORT     = $(BIN_DIR)/sort_go
BIN_GO_FIB_ITER = $(BIN_DIR)/fib_iter_go
BIN_GO_FIB_REC  = $(BIN_DIR)/fib_rec_go
BIN_GO_VAR_MC   = $(BIN_DIR)/var_mc_go

BIN_RS_SORT     = $(BIN_DIR)/sort_rs
BIN_RS_FIB_ITER = $(BIN_DIR)/fib_iter_rs
BIN_RS_FIB_REC  = $(BIN_DIR)/fib_rec_rs
BIN_RS_VAR_MC   = $(BIN_DIR)/var_mc_rs

BIN_TNG_SORT_QS    = $(BIN_DIR)/sort_cli_qsort
BIN_TNG_SORT_MS    = $(BIN_DIR)/sort_cli_msort
BIN_TNG_FIB_ITER   = $(BIN_DIR)/fib_cli
BIN_TNG_FIB_REC    = $(BIN_DIR)/fib_rec_cli
BIN_TNG_VAR_MC_S   = $(BIN_DIR)/var_mc_tng_sort
BIN_TNG_VAR_MC_Z   = $(BIN_DIR)/var_mc_tng_zig
BIN_TNG_VAR_MC_Q   = $(BIN_DIR)/var_mc_tng_qsel
BIN_TNG_SORT_PDQ   = $(BIN_DIR)/sort_cli_pdq
BIN_TNG_SORT_RADIX = $(BIN_DIR)/sort_cli_radix

.PHONY: all build clean c_benches go_benches rust_benches aot_benches bench_all plot

all: build

build: $(BIN_DIR) $(BIN_COMPILER) c_benches go_benches rust_benches aot_benches
	@echo "[build] All targets built."

clean:
	@echo "[clean] .bin removed."
	@rm -rf $(BIN_DIR)

$(BIN_COMPILER): | $(BIN_DIR)
	@echo "[build] compiler -> $@"
	@$(GO) build -o $@ $(CMD_COMPILER)

# C
c_benches: $(BIN_C_SORT) $(BIN_C_FIB_ITER) $(BIN_C_FIB_REC) $(BIN_C_VAR_MC)

$(BIN_C_SORT): | $(BIN_DIR)
	@echo "[build] c_sort -> $@"
	@$(CC) $(CFLAGS) $(C_SORT_SRC) internal/aotminic/runtime/runtime.c -o $@

$(BIN_C_FIB_ITER): | $(BIN_DIR)
	@echo "[build] c_fib_iter -> $@"
	@$(CC) $(CFLAGS) $(C_FIB_ITER_SRC) internal/aotminic/runtime/runtime.c -o $@

$(BIN_C_FIB_REC): | $(BIN_DIR)
	@echo "[build] c_fib_rec -> $@"
	@$(CC) $(CFLAGS) $(C_FIB_REC_SRC) internal/aotminic/runtime/runtime.c -o $@

$(BIN_C_VAR_MC): | $(BIN_DIR)
	@echo "[build] c_var_mc -> $@"
	@$(CC) $(CFLAGS) $(C_VAR_MC_SRC) -o $@

# Go
go_benches: $(BIN_GO_SORT) $(BIN_GO_FIB_ITER) $(BIN_GO_FIB_REC) $(BIN_GO_VAR_MC)

$(BIN_GO_SORT): | $(BIN_DIR)
	@echo "[build] go_sort -> $@"
	@$(GO) build -o $@ ./$(GO_SORT_SRC_DIR)

$(BIN_GO_FIB_ITER): | $(BIN_DIR)
	@echo "[build] go_fib_iter -> $@"
	@$(GO) build -o $@ $(GO_FIB_ITER_SRC)

$(BIN_GO_FIB_REC): | $(BIN_DIR)
	@echo "[build] go_fib_rec -> $@"
	@$(GO) build -o $@ $(GO_FIB_REC_SRC)

$(BIN_GO_VAR_MC): | $(BIN_DIR)
	@echo "[build] go_var_mc -> $@"
	@$(GO) build -o $@ $(GO_VAR_MC_SRC)

# Rust
rust_benches: $(BIN_RS_SORT) $(BIN_RS_FIB_ITER) $(BIN_RS_FIB_REC) $(BIN_RS_VAR_MC)

$(BIN_RS_SORT): | $(BIN_DIR)
	@echo "[build] rust_sort -> $@"
	@cd $(RUST_SORT_SRC_DIR) && cargo build --release && cp target/release/sort_rs "$(BIN_DIR_ABS)/sort_rs"

$(BIN_RS_FIB_ITER): | $(BIN_DIR)
	@echo "[build] rust_fib_iter -> $@"
	@cd $(RUST_FIB_ITER_SRC_DIR) && cargo build --release && cp target/release/fib_iter "$(BIN_DIR_ABS)/fib_iter_rs"

$(BIN_RS_FIB_REC): | $(BIN_DIR)
	@echo "[build] rust_fib_rec -> $@"
	@cd $(RUST_FIB_REC_SRC_DIR) && cargo build --release && cp target/release/fib_rec "$(BIN_DIR_ABS)/fib_rec_rs"

$(BIN_RS_VAR_MC): | $(BIN_DIR)
	@echo "[build] rust_var_mc -> $@"
	@cd $(RUST_VAR_MC_SRC_DIR) && cargo build --release && cp target/release/var_mc "$(BIN_DIR_ABS)/var_mc_rs"

# Tenge AOT
aot_benches: $(BIN_TNG_FIB_ITER) $(BIN_TNG_FIB_REC) $(BIN_TNG_SORT_QS) $(BIN_TNG_SORT_MS) $(BIN_TNG_VAR_MC_S) $(BIN_TNG_VAR_MC_Z) $(BIN_TNG_VAR_MC_Q) $(BIN_TNG_SORT_PDQ) $(BIN_TNG_SORT_RADIX)

$(BIN_TNG_FIB_ITER): $(BIN_COMPILER) | $(BIN_DIR)
	@echo "[aot] fib_iter -> $@"
	@./$(BIN_COMPILER) -o $@.c $(TNG_FIB_ITER_SRC)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TNG_FIB_REC): $(BIN_COMPILER) | $(BIN_DIR)
	@echo "[aot] fib_rec -> $@"
	@./$(BIN_COMPILER) -o $@.c $(TNG_FIB_REC_SRC)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TNG_SORT_QS): $(BIN_COMPILER) | $(BIN_DIR)
	@echo "[aot] sort_qsort -> $@"
	@./$(BIN_COMPILER) -o $@.c $(TNG_SORT_QS_SRC)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TNG_SORT_MS): $(BIN_COMPILER) | $(BIN_DIR)
	@echo "[aot] sort_msort -> $@"
	@./$(BIN_COMPILER) -o $@.c $(TNG_SORT_MS_SRC)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TNG_VAR_MC_S): $(BIN_COMPILER) | $(BIN_DIR)
	@echo "[aot] var_mc_sort -> $@"
	@./$(BIN_COMPILER) -o $@.c $(TNG_VAR_MC_SORT)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TNG_VAR_MC_Z): $(BIN_COMPILER) | $(BIN_DIR)
	@echo "[aot] var_mc_zig -> $@"
	@./$(BIN_COMPILER) -o $@.c $(TNG_VAR_MC_ZIG)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TNG_VAR_MC_Q): $(BIN_COMPILER) | $(BIN_DIR)
	@echo "[aot] var_mc_qsel -> $@"
	@./$(BIN_COMPILER) -o $@.c $(TNG_VAR_MC_QSEL)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TNG_SORT_PDQ): $(BIN_COMPILER) | $(BIN_DIR)
	@echo "[aot] sort_pdq -> $@"
	@./$(BIN_COMPILER) -o $@.c $(TNG_SORT_PDQ_SRC)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

$(BIN_TNG_SORT_RADIX): $(BIN_COMPILER) | $(BIN_DIR)
	@echo "[aot] sort_radix -> $@"
	@./$(BIN_COMPILER) -o $@.c $(TNG_SORT_RADIX_SRC)
	@$(CC) $(CFLAGS) $@.c $(AOT_RUNTIME_C) -o $@

bench_all: build
	@./benchmarks/run.sh

plot:
	@./benchmarks/plot.sh