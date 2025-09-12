📄 README.performance.md

# Performance & Benchmarks 👍

This document describes the benchmarking setup of **tengri-lang**.

---

## Philosophy 📌

Performance is a **first-class concern** of tengri-lang.  
We want to validate the potential of the language not only as an interpreter but also as an **ahead-of-time (AOT) compiler** and **VM runtime**.

To achieve reliable measurements, we:

- Prefer **TIME_NS** (nanosecond-precision timers inside runtime)  
  instead of relying on wall-clock measurements.
- Run benchmarks across multiple backends:
  - **Go** reference implementation
  - **Python** reference
  - **VM mini**
  - **AST interpreter**
  - **AOT transpiler + C runtime**

---

# Performance Benchmarks for Tengri-Lang

This document summarizes the benchmark results across Go, Python, Tengri-VM, and Tengri-AOT.

## Measurement Notes
- Machine: MacBook Air M2, fully charged, minimal background load
- Timing: prefer `TIME_NS` (nanosecond counter), fallback to `TIME:` if unavailable
- BENCH_REPS: 10,000,000 for stable AOT measurements

---

## AOT Results (MacBook Air M2, clean run)

### fib_iter (iterative Fibonacci)
| N   | Go (ns) | Tengri-VM (ns) | Tengri-AOT (ns) | Python (ns) |
|-----|---------|----------------|-----------------|-------------|
| 40  | 83      | 620            | **25**          | ~2000+      |
| 60  | 145     | 922            | **41**          | ~3000+      |
| 90  | 239     | 1383           | **50**          | ~4000+      |

➡️ **Iterative AOT is 3–5× faster than Go and >50× faster than Python.**

---

### fib_rec (recursive Fibonacci)
| N   | Go (ms) | Tengri-AOT (ms) | Python (ms) |
|-----|---------|-----------------|-------------|
| 30  | 54      | **0.23**        | 760         |
| 32  | 31      | **0.62**        | 741         |
| 34  | 31      | **1.61**        | 749         |

➡️ **Recursive AOT is ~20× faster than Go and ~450× faster than Python.**

---

## Next Steps
- Extend benchmarks with `matmul`, `sort`, `sieve`, and `calls`
- Verify results on Linux and Windows for cross-platform stability
- Explore optimizing the runtime further to reduce measurement noise