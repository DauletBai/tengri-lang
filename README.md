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


# Tengri Language

Tengri is an experimental programming language designed to explore novel concepts in compilation, runtime efficiency, and language design. Our core hypothesis is that the structural clarity of agglutinative languages can serve as a model for a more performant and intuitive computing paradigm.

The project follows a clear performance path: **AST Interpreter → Bytecode VM → JIT Compiler → AOT Compiler**.

---

## 🚀 Performance Highlights

Our benchmarks demonstrate that Tengri's AOT compiler generates highly optimized native code, **achieving performance parity with established languages like C, Rust, and Go.**

-   **Recursion-Heavy Tasks:** In recursive Fibonacci tests (`fib_rec` for N=34), Tengri-AOT is on par with C and Rust, and slightly faster than Go.
-   **Iteration-Heavy Tasks:** For iterative tasks (`fib_iter`), Tengri-AOT is in the same top tier as C, Rust, and Go, with execution times often too fast to be accurately measured in our micro-benchmark.

These results validate our core hypothesis. For detailed, reproducible results, see our [**Performance & Benchmarks Guide**](README.performance.md) and the raw CSV files in [`benchmarks/latest/results/`](benchmarks/latest/results/).

---

## Project Structure

.
├── benchmarks/     # Benchmark suite and results
├── cmd/            # CLI entry points (tengri-aot, tengri-vm, benchfast)
├── internal/       # Core compiler, runtime, and language internals
├── docs/           # Documentation, philosophy, and language vision
└── scripts/        # Helper scripts


---

## Build & Run

### Requirements
- Go (1.23+)
- C Compiler (Clang or GCC)
- Rust (for full benchmark comparison)
- GNU Make

### Build All Binaries
```bash
make build-all
Run Full Benchmark Suite
Bash

make bench
Contributing
Please read CONTRIBUTING.md before submitting patches. We welcome contributions that improve performance, clarity, and correctness.

License
MIT

## Latest Results (MacBook Air M2)

The following results were obtained from a full run of our benchmark suite. Times are the average per-call duration. Lower is better.

### Recursive Fibonacci (`fib_rec`, N=34)

This test measures the efficiency of function calls and recursion, a critical factor for overall performance.

| Target     | Time (seconds) | Relative to C |
| :--------- | :------------- | :------------ |
| **Tengri AOT** | **0.0268** | **0.99x** |
| C (baseline) | 0.0269         | 1.00x         |
| Rust       | 0.0278         | 1.03x         |
| Go         | 0.0311         | 1.15x         |

**Conclusion**: Tengri-AOT achieves **performance parity with native C and Rust code**, validating the effectiveness of the transpilation strategy.

### Iterative Fibonacci (`fib_iter`)

This test measures the efficiency of tight loops and basic arithmetic. The execution is so fast that results are near the measurement noise floor.

| Target     | Performance Tier |
| :--------- | :--------------- |
| **Tengri AOT** | **Tier 1 (Native Speed)** |
| C          | Tier 1 (Native Speed) |
| Rust       | Tier 1 (Native Speed) |
| Go         | Tier 1 (Native Speed) |
| VM         | Tier 2 (Very Fast Interpreter) |

**Conclusion**: For iterative code, Tengri-AOT is in the same top performance tier as other compiled languages. The VM provides excellent performance for a non-native backend.
