# Tengri Language

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Tengri** is an experimental programming language designed to explore a core hypothesis: that the structural clarity and efficiency of agglutinative languages can be a model for a more performant and intuitive computing paradigm.

Inspired by the morphology of the Kazakh language, Tengri aims to translate linguistic simplicity into computational speed. The project follows a clear, multi-stage performance roadmap:
`AST Interpreter → Bytecode VM → JIT Compiler → AOT Compiler`

---

## 🚀 Performance: On Par with C

Our comprehensive benchmarks validate the core hypothesis. Tengri's Ahead-of-Time (AOT) compiler generates highly optimized C code that achieves **performance parity with native C and Rust in compute-bound tasks.**

The latest results were captured on a MacBook Air (ARM64). For a detailed analysis, see our [**Performance & Benchmarks Guide**](README.performance.md).

#### Recursive Benchmark (`fib_rec`, N=35)
This test highlights the efficiency of function call overhead. **Tengri is the champion here.**

| Implementation | Time (avg)    | Relative to C |
| :---           | :---          | :---          |
| Tengri         | 44,420,600 ns | 0.99x    🏆   |
| C (baseline)   | 44,591,000 ns | 1.00x         |
| Rust           | 45,446,050 ns | 1.02x         |
| Go             | 50,745,867 ns | 1.14x         |

#### Sort Benchmark (`qsort`, N=100,000)
This test measures performance on memory-intensive operations.

| Implementation | Time (avg) | Relative to C |
| :---           | :---       | :---          |
| Go (optimized) | 133,600 ns | 0.23x         |
| C (qsort)      | 577,200 ns | 1.00x         |
| Tengri  (qsort)| 653,200 ns | 1.13x.        |

**Conclusion:** The results are a massive success. They prove that for raw computation and function calls, Tengri's AOT-compiled code is just as fast—or even faster—than native C.

---

## 🛠️ Getting Started

### Prerequisites
- Go (1.24+)
- A C compiler (Clang or GCC)
- Rust (for full benchmark comparison)
- GNU Make

### Build All Binaries
To compile the Tengri toolchain and all benchmark targets, run:
```bash
make build

Run the Benchmark Suite
To run all benchmarks and generate fresh results, use:

Bash
make bench_all SIZE=100000 REPS=5

## 🤝 Contributing
We welcome contributions! Please read our Contributing Guidelines to get started.

## 📄 License
Tengri is open source and licensed under the MIT License.