📄 CONTRIBUTING.md

# Contributing Guidelines

Thank you for your interest in contributing to **tengri-lang**!  
This document explains how the project is structured and how you can participate effectively.

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

## Build and Run

### Requirements
- Go 1.23+
- Python 3.10+
- Clang (for AOT runtime builds)
- GNU Make

### Setup
```bash
make setup

Run Benchmarks

make bench           # Run existing benchmarks
make bench-rebuild   # Force rebuild all .bin/* then run

Reports are saved into:
	•	benchmarks/latest/results/*.csv
	•	benchmarks/runs/<timestamp>/

⸻

Coding Guidelines
	•	Go code must follow gofmt and idiomatic Go practices.
	•	Comments in English, concise, for project documentation only.
	•	Commit messages use Conventional Commits:
	•	feat: add new parser rule
	•	fix: correct VM runtime stack
	•	refactor: reorganize benchmarks

⸻

Contribution Workflow
	1.	Fork the repo and create a feature branch:

git checkout -b feat/my-feature


	2.	Make your changes and run tests/benchmarks.
	3.	Commit with a descriptive message.
	4.	Open a Pull Request with details on motivation and design.

⸻

Bench Philosophy

We rely on strict timing (TIME_NS) over wall-clock, to ensure reproducibility across environments.
This project emphasizes performance transparency and comparability between implementations (Go, Python, VM, AOT).

---

### 📄 `README.performance.md`

```markdown
# Performance Benchmarks — tengri-lang

This document summarizes how performance is measured and where results are stored.

---

## Benchfast Tool

The main driver for benchmarks is [`cmd/benchfast`](../cmd/benchfast).

Features:
- Runs Fibonacci (iterative + recursive) in **Go**, **Python**, **VM**, and **Tengri AOT**.
- Supports CSV export (`benchmarks/latest/results/*.csv`).
- Supports plotting (`-plot`) via `gonum/plot`.
- Optionally rebuilds all binaries with `-rebuild`.

---

## Timing Methodology

⚡ **TIMING: prefer TIME_NS over wall-clock**  
All core benchmarks report **nanosecond-precision timing** collected internally.  
Wall-clock values are reported for reference but are not used for performance comparisons.

---

## Usage

### Quick run
```bash
make bench

Rebuild + run

make bench-rebuild

Plot results

go run cmd/benchfast/main.go -plot


⸻

Results 👍
	•	Latest CSV: benchmarks/latest/results/
	•	Historical runs: benchmarks/runs/<timestamp>/
	•	Plots: benchmarks/runs/<timestamp>/plots/

⸻

Notes for Contributors 📌
	•	If you add new benchmarks, place sources in benchmarks/src/<task>/<lang>/.
	•	Ensure outputs match across implementations (Go, Python, Tengri).
	•	Keep runtime environment consistent to avoid skew in performance data.