# 06_aot_minic — Minimal AOT (C-emitting) compiler

Tiny line-oriented transpiler from a minimal Tengri subset to C.
For early AOT demos & benchmarks (not a full parser).

## Build & Run
```bash
make -C 06_aot_minic build
make -C 06_aot_minic hello   # prints 42
make -C 06_aot_minic fib     # runs fibonacci (iter) with N=40
```

## Benchfast integration
Add a target in your benchfast to call:
- transpile: `.bin/tengri-aot 06_aot_minic/examples/fib_cli.tgr -o .bin/fib_cli.c`
- compile: `clang -O2 -o .bin/fib_cli .bin/fib_cli.c 06_aot_minic/runtime/runtime.c`
- run: `.bin/fib_cli <N>`
