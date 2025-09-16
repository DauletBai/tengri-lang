// FILE: benchmarks/src/c/fib_iter.c

#include "runtime.h"
#include <time.h> // Required for manual timing

long long run_fib_iter(int n) {
    if (n < 2) {
        return n;
    }
    long long a = 0, b = 1;
    for (int i = 2; i <= n; i++) {
        long long temp = a + b;
        a = b;
        b = temp;
    }
    return b;
}

int main(int argc, char** argv) {
    int n = get_n(argc, argv, 90);
    int inner_reps = 10000;
    
    // CORRECTED: Manual timing to calculate the average per operation.
    struct timespec start, end;
    clock_gettime(CLOCK_MONOTONIC, &start);

    for (int i = 0; i < inner_reps; i++) {
        // We use a volatile variable to prevent the compiler
        // from optimizing the loop away.
        volatile long long result = run_fib_iter(n);
        (void)result; // Suppress unused variable warning
    }

    clock_gettime(CLOCK_MONOTONIC, &end);
    long long total_ns = (end.tv_sec - start.tv_sec) * 1000000000LL + (end.tv_nsec - start.tv_nsec);
    long long avg_ns = total_ns / inner_reps; // The crucial division step.

    printf("TASK=fib_iter_c,N=%d,TIME_NS=%lld\n", n, avg_ns);
    
    return 0;
}