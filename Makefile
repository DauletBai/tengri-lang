SHELL := /bin/bash
GO    := GO111MODULE=on go

BIN_DIR := .bin

# Binaries
AOT_CLI := $(BIN_DIR)/tengri-aot
VM_CLI  := $(BIN_DIR)/tengri-vm
GO_ITER := $(BIN_DIR)/fib_iter_go

# Bench tool
BENCHFAST := ./cmd/benchfast

# Runtime
RUNTIME_C := internal/aotminic/runtime/runtime.c

# Bench sources (AOT examples already moved under benchmarks/src)
TGR_FIB_ITER := benchmarks/src/fib_iter/tengri/fib_cli.tgr
TGR_FIB_REC  := benchmarks/src/fib_rec/tengri/fib_rec_cli.tgr

.PHONY: all setup build clean bench-fast bench-plot aot aot-examples vm go-iter

all: build

setup:
	@$(GO) mod tidy
	@$(GO) get gonum.org/v1/plot@latest

build: aot vm go-iter

aot:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(AOT_CLI) ./cmd/tengri-aot

aot-examples: aot
	@$(AOT_CLI) $(TGR_FIB_ITER) -o $(BIN_DIR)/fib_cli.c
	@clang -O2 -o $(BIN_DIR)/fib_cli $(BIN_DIR)/fib_cli.c $(RUNTIME_C)
	@$(AOT_CLI) $(TGR_FIB_REC) -o $(BIN_DIR)/fib_rec_cli.c
	@clang -O2 -o $(BIN_DIR)/fib_rec_cli $(BIN_DIR)/fib_rec_cli.c $(RUNTIME_C)

vm:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(VM_CLI) ./cmd/tengri-vm || true

go-iter:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -tags=iter -o $(GO_ITER) ./benchmarks/src/fib_iter/go/fibonacci_iter.go || true

bench-fast: build aot-examples
	@$(GO) run ./cmd/benchfast

bench-plot: build aot-examples
	@$(GO) run ./cmd/benchfast -plot

clean:
	@rm -rf $(BIN_DIR)