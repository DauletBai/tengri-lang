SHELL := /bin/bash
GO := GO111MODULE=on go

# Бинарники/пути
BIN_DIR := .bin
AOT_BIN := $(BIN_DIR)/tengri-aot
VM_BIN  := $(BIN_DIR)/vm
GO_ITER := $(BIN_DIR)/fib_iter_go

# Файлы AOT-примеров
AOT_FIB_ITER_TGR := 06_aot_minic/examples/fib_cli.tgr
AOT_FIB_REC_TGR  := 06_aot_minic/examples/fib_rec_cli.tgr

# Инструменты
BENCHFAST := tool/benchfast/main.go

.PHONY: all setup build clean bench-fast bench-plot aot aot-examples go-vm

all: build

setup:
	@$(GO) mod tidy
	@$(GO) get gonum.org/v1/plot@latest

build: aot go-vm

aot:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(AOT_BIN) 06_aot_minic

aot-examples: aot
	@$(AOT_BIN) $(AOT_FIB_ITER_TGR) -o $(BIN_DIR)/fib_cli.c
	@clang -O2 -o $(BIN_DIR)/fib_cli $(BIN_DIR)/fib_cli.c 06_aot_minic/runtime/runtime.c
	@$(AOT_BIN) $(AOT_FIB_REC_TGR) -o $(BIN_DIR)/fib_rec_cli.c
	@clang -O2 -o $(BIN_DIR)/fib_rec_cli $(BIN_DIR)/fib_rec_cli.c 06_aot_minic/runtime/runtime.c

go-vm:
	@mkdir -p $(BIN_DIR)
	@go build -tags=iter -o $(GO_ITER) 04_benchmarks/fibonacci_iter.go || true
	@go build -o $(VM_BIN) 05_vm_mini/main.go || true

bench-fast: build aot-examples
	@go run $(BENCHFAST)

bench-plot: build aot-examples
	@go run $(BENCHFAST) -plot

clean:
	@rm -rf $(BIN_DIR)