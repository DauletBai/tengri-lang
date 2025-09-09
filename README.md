# Tengri Language (Tengri-lang)

[![Status](https://img.shields.io/badge/status-in_development-orange)](#)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

A research-first, open-source programming language inspired by the **agglutinative morphology of Kazakh**. Our goal is to turn *linguistic clarity* into *computational clarity* and deliver **predictable performance** across a staged toolchain:

**AST → VM → JIT → AOT**

> **Preliminary evidence**: our prototype VM already shows multi‑× speedups over Python on numeric kernels and approaches Go on selected microbenchmarks. Reproducible CSV/plots live under `benchmarks/` (see below).

---

## ✨ Design in a nutshell
- **Expressive minimalism.** Small, orthogonal core; explicit effects; visible costs.
- **Determinism & safety.** Predictable evaluation; explicit mutation; simple errors.
- **Performance path.** One semantics mapped consistently from AST to VM to JIT to AOT.
- **International by default.** Unicode‑ready tooling, English/Russian/Kazakh docs.

See also:
- `01_philosophy/mission_i18n.html` (Mission & Philosophy, EN/RU/KZ)
- `01_philosophy/governance_roadmap_i18n.html` (Governance & Roadmap, EN/RU/KZ)

---

## 📦 Repository layout

```
.
├── 01_philosophy/              # Mission, Governance (EN/RU/KZ), book materials
│   └── site/                   # Mini-site for GitHub Pages (index.html, mission, governance)
├── 02_prototype_python/        # Reference prototype in Python
├── 03_compiler_go/             # Go implementation (lexer, parser, AST evaluator)
├── 04_benchmarks/              # Standalone benchmark programs (Go/Python/Tengri)
├── 05_vm_mini/                 # Minimal register VM prototype
├── tools/benchfast/            # Cross-runtime benchmark runner (tables, CSV, plots)
├── benchmarks/                 # Results (versioned runs and "latest")
│   ├── latest/
│   │   ├── results/*.csv
│   │   └── plots/*.png
│   └── runs/YYYYmmdd-HHMMSS/...
├── Makefile                    # Common shortcuts (bench, bench-fast, bench-plot, bench-commit)
├── README.md                   # You are here
├── CONTRIBUTING.md             # Contributing guide
├── CODE_OF_CONDUCT.md          # CoC
└── LICENSE                     # MIT
```

---

## 🚀 Quick start

### Requirements
- Go (1.21+ recommended)
- Python 3.9+
- Optional (for plots): `gonum.org/v1/plot`
  ```bash
  go get gonum.org/v1/plot@latest
  go mod tidy
  ```

### Run sanity checks
```bash
# Fibonacci (recursive/iterative) reference programs
go run 04_benchmarks/fibonacci.go
go run -tags=iter 04_benchmarks/fibonacci_iter.go 60

python3 04_benchmarks/fibonacci.py
python3 04_benchmarks/fibonacci_iter.py 60

# Minimal VM
go run 05_vm_mini/main.go 60
```

### Cross-runtime benchmarks
```bash
# Tables + CSV
make bench-fast

# Tables + CSV + plots (PNG)
make bench-plot

# Save plots & CSV to repo and commit
make bench-commit
```
Outputs are stored under:
- `benchmarks/latest/results/*.csv`
- `benchmarks/latest/plots/*.png`
- and versioned in `benchmarks/runs/<timestamp>/...`

> The runner recognizes parse failures in the current Go interpreter and marks them as `ERR` (see `tools/benchfast/main.go`, `markStatus`).

---

## 🧠 Current status & known issues

- **Parser (Go implementation):** messages like “не найдена функция для разбора токена ')'” indicate a missing Pratt entry. Check `03_compiler_go/parser.go`:
  - Ensure `registerInfix(token.RPAREN, ...)` is **not** needed (right paren should usually terminate a sub‑expression).
  - Verify that *all* infix operators used in programs (e.g. `PLUS`, `MINUS`, `ASTERISK`, `SLASH`, `LT`, `GT`, `EQ`, `NOT_EQ`, custom `ARROW`, `SEMICOLON` handling) are registered via `registerInfix(tok, p.parseInfixExpression)` and that `parseInfixExpression` is implemented.
  - If you introduced a new token (e.g. `ARROW`), add it to `token/token.go`, teach the **lexer**, assign a **precedence**, and register a **prefix/infix parse function**.
  - Make sure prefix calls (`registerPrefix`) cover identifiers, integers, strings, booleans, unary `!`/`-`, grouped `(`expr`)`, function literals, call expressions, and indexing.

- **Bench runner:** keep the N-sets small by default to avoid long runs. You can override in `tools/benchfast/main.go`:
  ```go
  NsRec := []int{30, 32, 34}
  NsIter := []int{40, 60, 90}
  ```

- **Reproducibility:** each run writes CSV/plots under `benchmarks/runs/<timestamp>/...` and updates `benchmarks/latest/...` for Git diffs.

---

## 📈 Interpreting early numbers (rule of thumb)

- **VM vs Python:** expect ×3–×10 on numeric kernels (tight loops), depending on Python’s implementation and I/O.
- **VM vs Go:** VM will trail native Go; closing the gap requires simple peephole passes and hot‑path fusion.
- **JIT/AOT outlook:** reducing the gap to within ×2 of Go on kernel workloads is a realistic medium‑term target.

> Microbenchmarks are **not** the whole story. We’ll add string/JSON, maps, recursion, and IO-bound suites to get a balanced view.

---

## 🤝 Contributing

We welcome issues and PRs in English, Russian, or Kazakh. Please read:
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`

Good first issues:
- Parser: fill Pratt tables and precedence; fix error messages.
- VM: add opcodes, constant folding, peephole rules.
- Bench: add tasks (maps/reduce, string ops, JSON parse, matrix mul).

---

## 📚 Papers & community

We plan to publish English-language writeups and submit to PL venues (ICFP/PLDI/POPL) as the prototype matures. Follow the mini‑site:
- `01_philosophy/site/index.html` (or GitHub Pages if enabled)

---

## 📝 License

MIT © Tengri Language contributors
