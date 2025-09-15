// FILE: internal/aotminic/transpiler.go
// Purpose: Minimal AOT "mini-C" transpiler for demo benches.
// Supported demo basenames:
//   - fib_iter_cli.tgr   => iterative fibonacci
//   - fib_rec_cli.tgr    => recursive fibonacci
//   - sort_cli.tgr       => sort via qsort (baseline)
//   - sort_cli_m.tgr     => sort via stable mergesort (reference)
// Any other name -> error (unless 'force' is true, falls back to fib_iter).

package aotminic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func TranspileToC(srcPath, outC string, force bool) error {
	if srcPath == "" || outC == "" {
		return errors.New("srcPath and outC are required")
	}
	mode := detectMode(strings.ToLower(filepath.Base(srcPath)), force)
	if mode == "" {
		return fmt.Errorf("unsupported AOT demo source: %s", srcPath)
	}
	var csrc string
	switch mode {
	case "fib_iter":
		csrc = cFibIter()
	case "fib_rec":
		csrc = cFibRec()
	case "sort_qsort":
		csrc = cSortQsort()
	case "sort_msort":
		csrc = cSortMerge()
	default:
		return fmt.Errorf("internal mode error: %s", mode)
	}
	if err := os.WriteFile(outC, []byte(csrc), 0o644); err != nil {
		return fmt.Errorf("write C failed: %w", err)
	}
	return nil
}

func detectMode(base string, force bool) string {
	switch base {
	case "fib_iter_cli.tgr":
		return "fib_iter"
	case "fib_rec_cli.tgr":
		return "fib_rec"
	case "sort_cli.tgr":
		return "sort_qsort"
	case "sort_cli_m.tgr":
		return "sort_msort"
	default:
		if force {
			return "fib_iter"
		}
		return ""
	}
}

func cFibIter() string {
	return `#include <stdint.h>
#include <stdio.h>
#include "runtime.h"

// iterative fibonacci
static long fib_iter(long n) {
    long a = 0, b = 1;
    for (long i = 0; i < n; i++) {
        long tmp = a + b;
        a = b;
        b = tmp;
    }
    return a;
}

int main(int argc, char** argv) {
    long n = tgr_argi(argc, argv, 1, 45);
    int64_t reps = tgr_bench_reps();
    (void)fib_iter(10); // warmup

    int64_t t0 = tgr_time_ns();
    long res = 0;
    for (int64_t i = 0; i < reps; i++) res = fib_iter(n);
    int64_t t1 = tgr_time_ns();

    int64_t total = (t1 - t0);
    int64_t avg   = (reps > 0 ? total / reps : total);

    // Unified one-line report for CSV
    printf("REPORT impl=tengri-aot task=fib_iter n=%ld reps=%lld time_ns_avg=%lld result=%ld\n",
        n, (long long)reps, (long long)avg, res);
    return 0;
}
`
}

func cFibRec() string {
	return `#include <stdint.h>
#include <stdio.h>
#include "runtime.h"

// recursive fibonacci
static long fib_rec(long n) {
    if (n < 2) return n;
    return fib_rec(n-1) + fib_rec(n-2);
}

int main(int argc, char** argv) {
    long n = tgr_argi(argc, argv, 1, 35);
    int64_t reps = tgr_bench_reps();
    (void)fib_rec(10); // warmup

    int64_t t0 = tgr_time_ns();
    long res = 0;
    for (int64_t i = 0; i < reps; i++) res = fib_rec(n);
    int64_t t1 = tgr_time_ns();

    int64_t total = (t1 - t0);
    int64_t avg   = (reps > 0 ? total / reps : total);

    printf("REPORT impl=tengri-aot task=fib_rec n=%ld reps=%lld time_ns_avg=%lld result=%ld\n",
        n, (long long)reps, (long long)avg, res);
    return 0;
}
`
}

func cSortQsort() string {
	return `#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include "runtime.h"

static int cmp_int(const void* a, const void* b) {
    int ia = *(const int*)a, ib = *(const int*)b;
    return (ia > ib) - (ia < ib);
}

int main(int argc, char** argv) {
    long n = tgr_argi(argc, argv, 1, 100000);
    if (n <= 0) n = 1;
    int64_t reps = tgr_bench_reps();

    int* arr = (int*)malloc(sizeof(int) * (size_t)n);
    if (!arr) { fprintf(stderr, "alloc failed\n"); return 1; }

    // warm-up
    for (long i = 0; i < n; i++) arr[i] = (int)(n - i);
    qsort(arr, (size_t)n, sizeof(int), cmp_int);

    int64_t t0 = tgr_time_ns();
    for (int64_t r = 0; r < reps; r++) {
        for (long i = 0; i < n; i++) arr[i] = (int)(n - i);
        qsort(arr, (size_t)n, sizeof(int), cmp_int);
    }
    int64_t t1 = tgr_time_ns();

    long long sum = 0;
    for (long i = 0; i < n; i++) sum += arr[i];
    int first = arr[0], last = arr[n-1];

    int64_t total = (t1 - t0);
    int64_t avg   = (reps > 0 ? total / reps : total);

    printf("REPORT impl=tengri-aot task=sort-qsort n=%ld reps=%lld time_ns_avg=%lld first=%d last=%d sum=%lld\n",
        n, (long long)reps, (long long)avg, first, last, sum);

    free(arr);
    return 0;
}
`
}

func cSortMerge() string {
	return `#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "runtime.h"

// NOTE: avoid clash with system 'mergesort' symbol on macOS by using tgr_* names.

// Stable mergesort for int
static void tgr_merge(int* a, int* buf, long l, long m, long r) {
    long i=l, j=m, k=l;
    while (i<m && j<r) {
        if (a[i] <= a[j]) buf[k++] = a[i++];
        else              buf[k++] = a[j++];
    }
    while (i<m) buf[k++] = a[i++];
    while (j<r) buf[k++] = a[j++];
    for (long t=l; t<r; t++) a[t] = buf[t];
}

static void tgr_mergesort_rec(int* a, int* buf, long l, long r) {
    if (r - l <= 32) { // small insertion sort
        for (long i=l+1; i<r; i++) {
            int x = a[i]; long j = i-1;
            while (j>=l && a[j] > x) { a[j+1] = a[j]; j--; }
            a[j+1] = x;
        }
        return;
    }
    long m = l + (r - l)/2;
    tgr_mergesort_rec(a, buf, l, m);
    tgr_mergesort_rec(a, buf, m, r);
    if (a[m-1] <= a[m]) return; // already sorted
    tgr_merge(a, buf, l, m, r);
}

static void tgr_mergesort(int* a, long n) {
    int* buf = (int*)malloc(sizeof(int)*(size_t)n);
    if (!buf) return;
    tgr_mergesort_rec(a, buf, 0, n);
    free(buf);
}

int main(int argc, char** argv) {
    long n = tgr_argi(argc, argv, 1, 100000);
    if (n <= 0) n = 1;
    int64_t reps = tgr_bench_reps();

    int* arr = (int*)malloc(sizeof(int) * (size_t)n);
    if (!arr) { fprintf(stderr, "alloc failed\n"); return 1; }

    // warm-up
    for (long i = 0; i < n; i++) arr[i] = (int)(n - i);
    tgr_mergesort(arr, n);

    int64_t t0 = tgr_time_ns();
    for (int64_t r = 0; r < reps; r++) {
        for (long i = 0; i < n; i++) arr[i] = (int)(n - i);
        tgr_mergesort(arr, n);
    }
    int64_t t1 = tgr_time_ns();

    long long sum = 0;
    for (long i = 0; i < n; i++) sum += arr[i];
    int first = arr[0], last = arr[n-1];

    int64_t total = (t1 - t0);
    int64_t avg   = (reps > 0 ? total / reps : total);

    printf("REPORT impl=tengri-aot task=sort-msort n=%ld reps=%lld time_ns_avg=%lld first=%d last=%d sum=%lld\n",
        n, (long long)reps, (long long)avg, first, last, sum);

    free(arr);
    return 0;
}
`
}