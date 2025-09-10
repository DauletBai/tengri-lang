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

## Benchfast Tool

All benchmarks are orchestrated by the `benchfast` tool:

```bash
go run cmd/benchfast/main.go

Options
	•	-plot → also generate plots in benchmarks/runs/*/plots/
	•	-rebuild → rebuilds .bin/* binaries before running benchmarks

⸻

Makefile Targets

For convenience, common benchmark workflows are scripted:

# Run standard benchmarks
make bench

# Run benchmarks and regenerate plots
make bench-plot

# Force rebuild of all binaries before running
make bench-rebuild


⸻

Benchmark Sources

Benchmark source files are located in:

benchmarks/src/
├── fib_rec/      # Recursive Fibonacci
│   ├── go/
│   ├── python/
│   └── tengri/
└── fib_iter/     # Iterative Fibonacci
    ├── go/
    ├── python/
    └── tengri/


⸻

Interpreting Results
	•	TIMING line in reports indicates:
TIMING: prefer TIME_NS over wall-clock
	•	Each task prints outputs for all runtimes side by side.
	•	CSV results are saved in:
	•	benchmarks/latest/results/
	•	benchmarks/runs/<timestamp>/results/

Plots are saved in benchmarks/runs/<timestamp>/plots/.

⸻

Contribution Guidelines

When adding new benchmarks:
	1.	Place sources in benchmarks/src/<task>/<lang>/.
	2.	Ensure benchfast can run them via the Targets table.
	3.	Use TIME_NS if the runtime supports it.

---

