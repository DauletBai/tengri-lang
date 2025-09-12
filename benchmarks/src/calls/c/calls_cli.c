// calls_cli.c — micro: function call overhead
// Build: clang -O2 -o .bin/calls_cli benchmarks/src/calls/c/calls_cli.c internal/aotminic/runtime/runtime.c -Iinternal/aotminic/runtime
#include "runtime.h"
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>

#if defined(__GNUC__) || defined(__clang__)
__attribute__((noinline))
#endif
static int callee(int x) {
    // do something trivial that can’t be folded away
    return x + 1;
}

static long get_env_long(const char* name, long def) {
    const char* s = getenv(name);
    if (!s || !*s) return def;
    char* end = NULL;
    long v = strtol(s, &end, 10);
    return end && *end == '\0' ? v : def;
}

int main(int argc, char** argv) {
    // N = iterations inside one rep (default 1_000)
    long N = argi(argc, argv, 1, 1000);
    // REPS = how many times we repeat the whole inner loop (default from env)
    long REPS = get_env_long("BENCH_REPS", 5000000L);
    long WARMUP = get_env_long("WARMUP", 5L);

    volatile int sink = 0;

    // warmup
    for (long w = 0; w < WARMUP; ++w) {
        for (long i = 0; i < N; ++i) sink += callee((int)i);
    }

    long t0 = time_ns();
    for (long r = 0; r < REPS; ++r) {
        for (long i = 0; i < N; ++i) sink += callee((int)i);
    }
    long t1 = time_ns();
    long dt = t1 - t0;
    if (dt <= 0) dt = 1;
    long avg = dt / (REPS ? REPS : 1);

    printf("CHECKSUM: %d\n", sink);
    print_time_ns(avg);
    return 0;
}