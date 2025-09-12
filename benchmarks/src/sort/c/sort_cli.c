// sort_cli.c — Quicksort of N random ints (fixed seed)
// Build: clang -O2 -o .bin/sort_cli benchmarks/src/sort/c/sort_cli.c internal/aotminic/runtime/runtime.c -Iinternal/aotminic/runtime
#include "runtime.h"
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>

static long get_env_long(const char* name, long def) {
    const char* s = getenv(name);
    if (!s || !*s) return def;
    char* end = NULL;
    long v = strtol(s, &end, 10);
    return end && *end == '\0' ? v : def;
}

static uint32_t xorshift32(uint32_t* s) {
    uint32_t x = *s;
    x ^= x << 13; x ^= x >> 17; x ^= x << 5;
    *s = x; return x;
}

static void fill_rand(int* a, int n, uint32_t* seed) {
    for (int i = 0; i < n; ++i) a[i] = (int)(xorshift32(seed) & 0x7fffffff);
}

static void swap(int* a, int* b){ int t=*a; *a=*b; *b=t; }

static void qsort_rec(int* a, int l, int r) {
    while (l < r) {
        int i = l, j = r;
        int p = a[(l + r) >> 1];
        while (i <= j) {
            while (a[i] < p) ++i;
            while (a[j] > p) --j;
            if (i <= j) { swap(&a[i], &a[j]); ++i; --j; }
        }
        // tail recursion elimination
        if (j - l < r - i) { if (l < j) qsort_rec(a, l, j); l = i; }
        else { if (i < r) qsort_rec(a, i, r); r = j; }
    }
}

static uint64_t checksum_i(const int* a, size_t n) {
    uint64_t h = 1469598103934665603ull;
    for (size_t i = 0; i < n; ++i) {
        h ^= (uint64_t)(uint32_t)a[i];
        h *= 1099511628211ull;
    }
    return h;
}

int main(int argc, char** argv) {
    int N = (int)argi(argc, argv, 1, 200000);      // default 2e5
    long REPS = get_env_long("BENCH_REPS", 10L);
    long WARMUP = get_env_long("WARMUP", 2L);
    uint32_t seed0 = (uint32_t)get_env_long("SEED", 42L);

    int* a = (int*)malloc((size_t)N * sizeof(int));
    int* b = (int*)malloc((size_t)N * sizeof(int));
    if (!a || !b) { fprintf(stderr, "OOM\n"); return 1; }

    // warmup
    uint32_t s = seed0;
    fill_rand(a, N, &s);
    for (long w = 0; w < WARMUP; ++w) {
        for (int i = 0; i < N; ++i) b[i] = a[i];
        qsort_rec(b, 0, N - 1);
    }

    long t0 = time_ns();
    for (long r = 0; r < REPS; ++r) {
        s = seed0 + (uint32_t)r;   // vary input across reps but deterministically
        fill_rand(a, N, &s);
        for (int i = 0; i < N; ++i) b[i] = a[i];
        qsort_rec(b, 0, N - 1);
    }
    long t1 = time_ns();
    long dt = t1 - t0;
    if (dt <= 0) dt = 1;
    long avg = dt / (REPS ? REPS : 1);

    uint64_t cs = checksum_i(b, (size_t)N);
    printf("CHECKSUM: %llu\n", (unsigned long long)cs);
    print_time_ns(avg);

    free(a); free(b);
    return 0;
}