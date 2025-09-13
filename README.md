# Tengri Language

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Tengri** is an experimental programming language designed to explore a core hypothesis: that the structural clarity and efficiency of agglutinative languages can be a model for a more performant and intuitive computing paradigm.

Inspired by the morphology of the Kazakh language, Tengri aims to translate linguistic simplicity into computational speed. The project follows a clear, multi-stage performance roadmap:
`AST Interpreter → Bytecode VM → JIT Compiler → AOT Compiler`

---

## 🚀 Performance: On Par with C & Rust

Our benchmarks validate the core hypothesis. Tengri's Ahead-of-Time (AOT) compiler generates highly optimized C code that achieves **performance parity with native C and Rust**.

The latest results were captured on a MacBook Air M2. For the full methodology and raw data, see our [Performance Guide](README.performance.md) and the result files in `benchmarks/latest/results/`.

#### Recursive Benchmark (`fib_rec`, N=34)
This test highlights the efficiency of function call overhead.

| Target       | Time (avg)    | Relative to C |
| :---         | :---          | :---  |
| Tengri AOT   | 26,894,750 ns | 0.99x |
| C (baseline) | 26,957,000 ns | 1.00x |
| Rust         | 27,806,459 ns | 1.03x |
| Go           | 31,128,292 ns | 1.15x |

#### Iterative Benchmark (`fib_iter`, N=40)
This test highlights the performance of tight loops and arithmetic.

| Target |Time (avg)| Note          |
| :---   | :---     | :---          |
| Tengri | 41 ns    | Champion** 🏆 |
| C      | 42 ns    | Native Speed  |
| Go     | 146 ns   | Native Speed  |
| Rust   | 125 ns   | Native Speed  |
| VM     | 2,443 ns | Excellent for a VM |

---

## ✨ Core Philosophy

-   **Expressive Minimalism:** A small, orthogonal set of features that compose into powerful patterns.
-   **Performance by Design:** A clear path from high-level semantics to efficient machine code.
-   **Determinism:** Predictable evaluation with visible costs.
-   **International by Default:** Unicode-native syntax and standard library.

---

## 🛠️ Getting Started

### Prerequisites
- Go (1.23+)
- A C compiler (Clang or GCC)
- Rust (for full benchmark comparison)
- GNU Make

### Build All Binaries
To compile the Tengri toolchain and all benchmark targets, run:
```bash
make build-all
```

### Run the Benchmark Suite
To run all benchmarks and generate fresh results, use:
```bash
make bench
```

---

## 🚧 Project Status

The project is currently in the **proof-of-concept** stage.

-   ✅ **AOT Compiler:** The transpiler successfully generates high-performance C code for our benchmark set, proving the viability of the core concept.
-   🔶 **VM & Interpreter:** A functional bytecode VM and AST interpreter exist as proofs-of-concept.
-   🔴 **Parser & Language Features:** The language parser is under development and is the current focus for stabilization and feature expansion (loops, variables, new types).

---

## 🤝 Contributing

We welcome contributions! Please read our [**Contributing Guidelines**](CONTRIBUTING.md) to get started.

## 📄 License

Tengri is open source and licensed under the [MIT License](LICENSE).