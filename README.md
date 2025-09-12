📄 README.md 🚀

# tengri-lang

**tengri-lang** is an experimental programming language and research project.  
It explores new approaches to **compilation, runtime efficiency, and digital sovereignty**, with the long-term vision of forming the basis for:

- Operating Systems
- Databases
- Communication Protocols
- Digital Economy and Security

---

## Project Structure

.
├── benchmarks/          # Benchmarks and reference implementations
│   └── src/             # Language-specific sources
│       ├── fib_iter/    # Iterative Fibonacci (Go, Python, Tengri)
│       └── fib_rec/     # Recursive Fibonacci (Go, Python, Tengri)
├── cmd/                 # CLI entry points
│   ├── benchfast/       # Benchmark runner
│   ├── tengri-aot/      # AOT transpiler CLI
│   └── tengri-vm/       # VM CLI
├── internal/            # Compiler, runtime, language internals
│   ├── aotminic/        # AOT backend + runtime (C)
│   └── lang/            # Lexer, parser, AST, evaluator
├── docs/                # Documentation
│   └── philosophy/      # Mission and vision
├── scripts/             # Helper scripts (restructure, CI helpers, etc.)
└── .bin/                # Built binaries (ignored in VCS)

---

## Build & Run

### Requirements
- Go 1.23+
- Python 3.10+
- Clang
- GNU Make

### Setup
```bash
make setup

Run Benchmarks

make bench


⸻

Performance & Benchmarks

Performance is a core focus of tengri-lang.
We provide strict benchmarking tools and prefer TIME_NS (nanosecond-precision internal timers) over wall-clock timing for reproducibility.

See the detailed Performance & Benchmarks Guide.

⸻

Contributing

Please read CONTRIBUTING.md before submitting patches.
We welcome contributions that improve performance, structure, and clarity.

⸻

License

MIT

---


