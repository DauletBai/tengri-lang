package aotminic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TranspileToC is a tiny demo "transpiler" that accepts two canonical sources:
//   - .../fib_iter_cli.tgr  → iterative Fibonacci CLI in C
//   - .../fib_rec_cli.tgr   → recursive Fibonacci CLI in C
//
// It emits a single C file that includes runtime timers and prints TIME_NS
// in a unified format that benchfast understands.
//
// NOTE: This is a minimal PoC, not a full Tengri frontend. The .tgr contents
// are not parsed; the filename is used as a switch.
func TranspileToC(sourcePath string) (string, error) {
	base := filepath.Base(sourcePath)
	switch {
	case strings.Contains(base, "fib_iter_cli"):
		return emitFibIterC(), nil
	case strings.Contains(base, "fib_rec_cli"):
		return emitFibRecC(), nil
	default:
		return "", errors.New("unsupported source pattern (mini-AOT demo expects fib_iter_cli / fib_rec_cli)")
	}
}

// emitFibIterC generates an iterative Fibonacci CLI.
// Includes are intentionally flat: #include "runtime.h".
// The Makefile provides -Iinternal/aotminic/runtime so headers resolve regardless of CWD.
func emitFibIterC() string {
	return strings.TrimLeft(`
#include "runtime.h"
#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>

static unsigned long long fib_iter(unsigned long long n) {
    if (n <= 1) return n;
    unsigned long long a = 0, b = 1;
    for (unsigned long long i = 2; i <= n; i++) {
        unsigned long long t = a + b;
        a = b;
        b = t;
    }
    return b;
}

int main(int argc, char** argv) {
    if (argc < 2) {
        printf("usage: %s <N>\n", argv[0]);
        return 1;
    }
    unsigned long long n = strtoull(argv[1], NULL, 10);

    long start = time_ns();
    unsigned long long result = fib_iter(n);
    long end = time_ns();

    // Unified benchfast-friendly output:
    printf("%" PRIu64 "\n", (uint64_t)result);
    print_time_ns(end - start);
    return 0;
}
`, "\n")
}

// emitFibRecC generates a naive recursive Fibonacci CLI.
func emitFibRecC() string {
	return strings.TrimLeft(`
#include "runtime.h"
#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>

static unsigned long long fib_rec(unsigned long long n) {
    if (n < 2) return n;
    return fib_rec(n - 1) + fib_rec(n - 2);
}

int main(int argc, char** argv) {
    if (argc < 2) {
        printf("usage: %s <N>\n", argv[0]);
        return 1;
    }
    unsigned long long n = strtoull(argv[1], NULL, 10);

    long start = time_ns();
    unsigned long long result = fib_rec(n);
    long end = time_ns();

    // Unified benchfast-friendly output:
    printf("%" PRIu64 "\n", (uint64_t)result);
    print_time_ns(end - start);
    return 0;
}
`, "\n")
}

// WriteFile is a small helper for tests/tools that want a file on disk.
func WriteFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	return os.WriteFile(path, []byte(content), 0o644)
}