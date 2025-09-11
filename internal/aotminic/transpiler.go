package aotminic

// Minimal, example-only AOT transpiler that recognizes two demo sources
// (fib_iter_cli.tgr / fib_rec_cli.tgr) and emits a small C program with
// a matching fib implementation, runtime glue, and normalized timing.
//
// The emitted TIME_NS is PER CALL: total_ns / REPS, where
// REPS = env("BENCH_REPS", default=5_000_000).

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TranspileToC emits a single C file at outC based on the source name.
// It supports two demo patterns:
//   * fib_iter_cli.tgr → iterative Fibonacci
//   * fib_rec_cli.tgr  → recursive Fibonacci
//
// If 'force' is true, any other filename will be treated as iterative.
func TranspileToC(srcPath, outC string, force bool) error {
	base := filepath.Base(srcPath)
	kind := detectKind(base, force)
	if kind == kindUnknown {
		return fmt.Errorf("cannot detect program kind (expected fib_iter(...) or fib_rec(...))")
	}

	c := buildC(kind)
	if err := writeFile(outC, c); err != nil {
		return err
	}
	return nil
}

type programKind int

const (
	kindUnknown programKind = iota
	kindIter
	kindRec
)

func detectKind(name string, force bool) programKind {
	n := strings.ToLower(name)
	switch {
	case strings.Contains(n, "fib_iter"):
		return kindIter
	case strings.Contains(n, "fib_rec"):
		return kindRec
	default:
		if force {
			return kindIter
		}
		return kindUnknown
	}
}

func writeFile(path, data string) error {
	if path == "" {
		return errors.New("output path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(data), 0o644)
}

func buildC(kind programKind) string {
	header := `#include "runtime.h"
#include <stdlib.h>

/* BENCH_REPS controls how many times the kernel runs.
 * We normalize to per-call TIME_NS by dividing total_ns / reps.
 */
static long get_reps(void) {
    const char* env = getenv("BENCH_REPS");
    if (!env || !env[0]) return 5000000L; /* default */
    char* end = NULL;
    long v = strtol(env, &end, 10);
    if (end == env || v <= 0) return 5000000L;
    return v;
}
`

	var fib string
	switch kind {
	case kindIter:
		fib = `
long long fib_iter(long long n) {
    if (n <= 1) return n;
    long long a = 0, b = 1;
    for (long long i = 2; i <= n; i++) {
        long long t = a + b;
        a = b;
        b = t;
    }
    return b;
}
`
	case kindRec:
		fib = `
long long fib_rec(long long n) {
    if (n <= 1) return n;
    return fib_rec(n - 1) + fib_rec(n - 2);
}
`
	}

	mainBodyIter := `
int main(int argc, char** argv) {
    long n = argi(argc, argv, 1, 40); /* default N=40 */
    long reps = get_reps();
    volatile long long acc = 0;

    long long start = time_ns();
    for (long i = 0; i < reps; i++) {
        acc += fib_iter(n);
    }
    long long end = time_ns();

    print(acc ? fib_iter(n) : 0); /* print single-call RESULT */

    long long total = end - start;
    long long per   = (reps > 0) ? (total / reps) : 0;
    print_time_ns(per);
    return 0;
}
`

	mainBodyRec := `
int main(int argc, char** argv) {
    long n = argi(argc, argv, 1, 34); /* default N=34 */
    long reps = get_reps();
    volatile long long acc = 0;

    long long start = time_ns();
    for (long i = 0; i < reps; i++) {
        acc += fib_rec(n);
    }
    long long end = time_ns();

    print(acc ? fib_rec(n) : 0); /* print single-call RESULT */

    long long total = end - start;
    long long per   = (reps > 0) ? (total / reps) : 0;
    print_time_ns(per);
    return 0;
}
`

	switch kind {
	case kindIter:
		return header + fib + mainBodyIter
	case kindRec:
		return header + fib + mainBodyRec
	default:
		// Should not happen; detectKind filters this.
		return header + "\nint main(){ return 1; }\n"
	}
}