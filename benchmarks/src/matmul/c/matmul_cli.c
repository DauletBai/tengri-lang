// matmul_cli.c — NxN int matrix multiply (C = A*B)
// Build: clang -O2 -o .bin/matmul_cli benchmarks/src/matmul/c/matmul_cli.c internal/aotminic/runtime/runtime.c -Iinternal/aotminic/runtime
#include "runtime.h"
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

static long get_env_long(const char* name, long def) {
    const char* s = getenv(name);
    if (!s || !*s) return def;
    char* end = NULL;
    long v = strtol(s, &end, 10);
    return end && *end == '\0' ? v : def;
}

static uint32_t xorshift32(uint32_t* s) {
    uint32_t x = *s;
    x ^= x << 13;
    x ^= x >> 17;
    x ^= x << 5;
    *s = x;
    return x;
}

static uint64_t checksum_i(int* a, size_t n) {
    uint64_t h = 1469598103934665603ull;
    for (size_t i = 0; i < n; ++i) {
        h ^= (uint64_t)(uint32_t)a[i];
        h *= 1099511628211ull;
    }
    return h;
}

static void fill_rand(int* a, int n, uint32_t* seed) {
    for (int i = 0; i < n*n; ++i) {
        a[i] = (int)(xorshift32(seed) % 1000) - 500;
    }
}

static void matmul(int* C, const int* A, const int* B, int n) {
    // naive i-k-j (better cache than i-j-k)
    for (int i = 0; i < n; ++i) {
        int* Ci = &C[i*n];
        for (int j = 0; j < n; ++j) Ci[j] = 0;
        for (int k = 0; k < n; ++k) {
            int aik = A[i*n + k];
            const int* Bk = &B[k*n];
            for (int j = 0; j < n; ++j) {
                Ci[j] += aik * Bk[j];
            }
        }
    }
}

int main(int argc, char** argv) {
    int n = (int)argi(argc, argv, 1, 64);              // default 64x64
    long REPS = get_env_long("BENCH_REPS", 50L);
    long WARMUP = get_env_long("WARMUP", 2L);
    uint32_t seed = (uint32_t)get_env_long("SEED", 42L);

    size_t bytes = (size_t)n * (size_t)n * sizeof(int);
    int* A = (int*)malloc(bytes);
    int* B = (int*)malloc(bytes);
    int* C = (int*)malloc(bytes);
    if (!A || !B || !C) { fprintf(stderr, "OOM\n"); return 1; }

    fill_rand(A, n, &seed);
    fill_rand(B, n, &seed);

    // warmup
    for (long w = 0; w < WARMUP; ++w) matmul(C, A, B, n);

    long t0 = time_ns();
    for (long r = 0; r < REPS; ++r) matmul(C, A, B, n);
    long t1 = time_ns();

    long dt = t1 - t0;
    if (dt <= 0) dt = 1;
    long avg = dt / (REPS ? REPS : 1);

    uint64_t cs = checksum_i(C, (size_t)n*(size_t)n);
    printf("CHECKSUM: %llu\n", (unsigned long long)cs);
    print_time_ns(avg);

    free(A); free(B); free(C);
    return 0;
}