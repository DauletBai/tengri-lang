#!/usr/bin/env bash
set -euo pipefail

# Helper: git mv if source exists and dest dir prepared
mv_if() {
  local src="$1" dst="$2"
  if [ -e "$src" ]; then
    mkdir -p "$(dirname "$dst")"
    git mv "$src" "$dst"
    echo "moved: $src -> $dst"
  else
    echo "skip (not found): $src"
  fi
}

# Create baseline dirs
mkdir -p cmd/benchfast cmd/tengri-aot cmd/tengri-vm cmd/tengri-ast
mkdir -p internal/aotminic/runtime internal/lang/{token,lexer,ast,parser,object,evaluator} internal/vmmini
mkdir -p benchmarks/src/{fib_iter/tengri,fib_iter/go,fib_iter/python,fib_rec}
mkdir -p benchmarks/latest/{results,plots} benchmarks/runs scripts

# AOT: transpiler + runtime + examples
mv_if 06_aot_minic/main.go                       internal/aotminic/transpiler.go
mv_if 06_aot_minic/runtime/runtime.c             internal/aotminic/runtime/runtime.c
mv_if 06_aot_minic/runtime/runtime.h             internal/aotminic/runtime/runtime.h
mv_if 06_aot_minic/examples/fib_cli.tgr          benchmarks/src/fib_iter/tengri/fib_cli.tgr
mv_if 06_aot_minic/examples/fib_rec_cli.tgr      benchmarks/src/fib_rec/tengri/fib_rec_cli.tgr

# VM mini
mv_if 05_vm_mini/main.go                         cmd/tengri-vm/main.go

# AST interpreter (CLI) + compiler packages
mv_if 03_compiler_go/main.go                     cmd/tengri-ast/main.go
mv_if 03_compiler_go/token.go                    internal/lang/token/token.go
mv_if 03_compiler_go/lexer.go                    internal/lang/lexer/lexer.go
mv_if 03_compiler_go/ast.go                      internal/lang/ast/ast.go
mv_if 03_compiler_go/parser.go                   internal/lang/parser/parser.go
mv_if 03_compiler_go/object.go                   internal/lang/object/object.go
mv_if 03_compiler_go/environment.go              internal/lang/object/environment.go
mv_if 03_compiler_go/evaluator.go                internal/lang/evaluator/evaluator.go

# Bench sources
mv_if 04_benchmarks/fibonacci.go                 benchmarks/src/fib_rec/fibonacci.go
mv_if 04_benchmarks/fibonacci.py                 benchmarks/src/fib_rec/fibonacci.py
mv_if 04_benchmarks/fibonacci.tengri             benchmarks/src/fib_rec/fibonacci.tgr
mv_if 04_benchmarks/fibonacci_iter.go            benchmarks/src/fib_iter/go/fibonacci_iter.go
mv_if 04_benchmarks/fibonacci_iter.py            benchmarks/src/fib_iter/python/fibonacci_iter.py

# Bench tool
mv_if tool/benchfast/main.go                     cmd/benchfast/main.go

echo "Done. Now add cmd/tengri-aot/main.go and update Makefile/go.mod."