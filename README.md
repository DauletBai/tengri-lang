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
├── cmd/                 # CLI entry points (benchfast, VM, AOT)
├── internal/            # Compiler, runtime, language internals
├── docs/                # Documentation and philosophy
├── scripts/             # Helper scripts
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


