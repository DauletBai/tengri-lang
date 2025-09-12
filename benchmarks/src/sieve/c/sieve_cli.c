// sieve_cli.c — Sieve of Eratosthenes (count primes <= N)
// Build: clang -O2 -o .bin/sieve_cli benchmarks/src/sieve/c/sieve_cli.c internal/aotminic/runtime/runtime.c -Iinternal/aotminic/runtime
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

static uint64_t checksum_u64(const uint8_t* a, size_t n) {
    // FNV-1a
    uint64_t h = 1469598103934665603ull;
    for (size_t i = 0; i < n; ++i) {
        h ^= a[i];
        h *= 1099511628211ull;
    }
    return h;
}

static long sieve_count(long N) {
    if (N < 2) return 0;
    size_t sz = (size_t)(N + 1);
    uint8_t* is_comp = (uint8_t*)malloc(sz);
    if (!is_comp) return -1;
    memset(is_comp, 0, sz);
    long cnt = 0;

    for (long p = 2; p * p <= N; ++p) {
        if (!is_comp[p]) {
            for (long q = p * p; q <= N; q += p) is_comp[q] = 1;
        }
    }
    for (long i = 2; i <= N; ++i) if (!is_comp[i]) cnt++;
    // keep array alive briefly (avoid DCE)
    volatile uint64_t cs = checksum_u64(is_comp, sz);
    (void)cs;
    free(is_comp);
    return cnt;
}

int main(int argc, char** argv) {
    long N = argi(argc, argv, 1, 100000);          // default 1e5
    long REPS = get_env_long("BENCH_REPS", 100L);  // heavier than micro
    long WARMUP = get_env_long("WARMUP", 3L);

    long ans = 0;
    for (long w = 0; w < WARMUP; ++w) ans ^= sieve_count(N);

    long t0 = time_ns();
    for (long r = 0; r < REPS; ++r) ans ^= sieve_count(N);
    long t1 = time_ns();
    long dt = t1 - t0;
    if (dt <= 0) dt = 1;
    long avg = dt / (REPS ? REPS : 1);

    printf("RESULT: %ld\n", ans);
    printf("CHECKSUM: %ld\n", ans);
    print_time_ns(avg);
    return 0;
}