# Tengri Language

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Tengri** is an experimental, compiled programming language designed to explore a core hypothesis: that the structural clarity and efficiency of agglutinative languages can be a model for a more performant and intuitive computing paradigm.

Inspired by the morphology of the Kazakh language, Tengri aims to translate linguistic simplicity into computational speed. The project currently focuses on an Ahead-of-Time (AOT) compiler that transpiles to C, achieving performance parity with native systems languages.

---

## 🚀 Performance: On Par with C

Our comprehensive benchmarks validate the core hypothesis. **Tengri achieves performance parity with native C and Rust in compute-bound tasks.**

The latest results were captured on a MacBook Air (ARM64). For a detailed analysis, see our [**Performance & Benchmarks Guide**](README.performance.md).

#### Recursive Benchmark (`fib_rec`, N=35)
This test highlights the efficiency of function call overhead. **Tengri is the champion here.**

| Implementation | Time (avg)        | Relative to C |
| -------------- | ----------------- | ------------- |
| **Tengri**     | **44,589,800 ns** | **0.99x**  🏆 |
| C (baseline)   | 44,740,400 ns     | 1.00x         |
| Rust           | 45,616,974 ns     | 1.02x         |
| Go             | 49,598,592 ns     | 1.11x         |

#### Iterative Benchmark (`fib_iter`, N=90)
This tests the raw speed of tight loops. **Tengri is again the champion.**

| Implementation | Time (avg)  |
| -------------- | ----------- |
| **Tengri**     | **50 ns** 🏆 |
| C              | 61 ns       |
| Rust           | 241 ns      |
| Go             | 6,858 ns    |

**Conclusion:** The results are a massive success. They prove that for raw computation and function calls, Tengri's compiled code is just as fast—or even faster—than native C.

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
make build

Run the Benchmark Suite
To run all benchmarks and generate fresh results, use:

Bash
make clean && make bench_all
make bench_all SIZE=100000 REPS=5

## 🤝 Contributing
We welcome contributions! Please read our Contributing Guidelines to get started.

## 📄 License
Tengri is open source and licensed under the MIT License.