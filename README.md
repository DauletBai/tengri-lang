# 👍 Tengri-Lang Benchmarks (AOT vs C / Go / Rust)

## Status (2025-09-15)

We added **honest AOT benchmarks** for Fibonacci and Sort:

- **fib_iter (n=45)**
- **fib_rec (n=35)**
- **sort (SIZE=100000)**

All benchmarks are executed via:

```bash
make bench_all SIZE=100000 BENCH_REPS=5

The results are logged into CSV files under
benchmarks/results/suite_YYYYmmdd_HHMMSS.csv.

Plots can be generated with:

make plot_csv
make plot_csv PLOT_LOG=1           # logarithmic Y scale
make plot_csv PLOT_REL=1           # normalized to best implementation
make plot_csv PLOT_LOG=1 PLOT_REL=1

(requires gnuplot installed)

⸻

Results

Example run (SIZE=100000, REPS=5):

Impl	            Task  time_ns_avg	Notes
Rust	            sort	  ~91k ns	Fastest, specialized sort for int
Go	                sort	 ~199k ns	Specialized sort.Ints, inline comparisons
C         (qsort)	sort	~2.55M ns	Slow due to comparator function pointer
Tengri-AOT qsort	sort	~2.55M ns	Matches C, proves correct runtime integration
Tengri-AOT msort	sort	~2.64M ns	Slightly slower than qsort, but stable
fib_iter   (all)fib_iter	~50–80 ns	Too small, dominated by timer noise
fib_rec    (all)fib_rec	 ~8.7–10.3 ms	All implementations align, expected exponential recursion cost

⸻

Conclusions
	•	Tengri-AOT integration is correct — overhead is negligible, runtime matches C performance.
	•	To be competitive with Go/Rust, we need a specialized integer sort without function pointers (inline comparisons).
	•	Go and Rust outperform qsort on primitives because their libraries avoid function pointers and allow full inlining/vectorization.

⸻

Roadmap
	•	Add int introsort/timsort in AOT (specialized, no function pointers).
	•	Add radix/counting sort for O(n) demonstration on integers.
	•	Benchmark larger input sizes (1e6, 1e7) with adjusted repetitions.
	•	Normalize and visualize results with make plot_csv (PLOT_LOG, PLOT_REL).
	•	Publish benchmark results and plots on GitHub Pages.

⸻

✅ Usage

# Clean and rebuild everything
make clean && make build

# Run all benchmarks with configurable size and reps
make bench_all SIZE=100000 BENCH_REPS=5

# Generate plots (requires gnuplot)
make plot_csv
make plot_csv PLOT_LOG=1
make plot_csv PLOT_REL=1
make plot_csv PLOT_LOG=1 PLOT_REL=1
