// FILE: benchmarks/src/sort/c/sort.c
// Purpose: Sort benchmark with nanosecond timer and unified REPORT line.

#include <stdio.h>
#include <stdlib.h>
#include <time.h>

static int cmp_int(const void* a, const void* b) {
    int ia = *(const int*)a, ib = *(const int*)b;
    return (ia > ib) - (ia < ib);
}

static long getenv_long(const char* key, long defv) {
    const char* v = getenv(key);
    if (!v || !*v) return defv;
    long n = atol(v);
    return n > 0 ? n : defv;
}

static long long now_ns(void) {
#if defined(CLOCK_MONOTONIC)
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (long long)ts.tv_sec*1000000000LL + ts.tv_nsec;
#else
    return 0;
#endif
}

int main(void) {
    long n = getenv_long("SIZE", 100000);
    long reps = getenv_long("BENCH_REPS", 3);
    if (n <= 0) n = 1; if (reps <= 0) reps = 1;

    int* a = (int*)malloc(sizeof(int)*(size_t)n);
    if (!a) { fprintf(stderr, "alloc failed\n"); return 1; }

    // warm-up
    for (long i=0;i<n;i++) a[i] = (int)(n - i);
    qsort(a, (size_t)n, sizeof(int), cmp_int);

    long long t0 = now_ns();
    for (long r=0; r<reps; r++) {
        for (long i=0;i<n;i++) a[i] = (int)(n - i);
        qsort(a, (size_t)n, sizeof(int), cmp_int);
    }
    long long t1 = now_ns();

    long long total = (t1 - t0);
    long long avg = total / reps;

    long long sum = 0;
    for (long i=0;i<n;i++) sum += a[i];
    int first = a[0], last = a[n-1];

    printf("REPORT impl=c task=sort n=%ld reps=%ld time_ns_avg=%lld first=%d last=%d sum=%lld\n",
        n, reps, avg, first, last, sum);

    free(a);
    return 0;
}