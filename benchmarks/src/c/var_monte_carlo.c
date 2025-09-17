// FILE: benchmarks/src/c/var_monte_carlo.c
// Purpose: Monte Carlo VaR benchmark (GBM, Box–Muller, xorshift64*).
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <stdint.h>
#include <time.h>

// Fast RNG: xorshift64*
static uint64_t xs_state = 0x9E3779B97F4A7C15ULL;
static inline void xs_seed(uint64_t s) { xs_state = s ? s : 0x9E3779B97F4A7C15ULL; }
static inline uint64_t xs_next_u64(void) {
    uint64_t x = xs_state;
    x ^= x >> 12; x ^= x << 25; x ^= x >> 27;
    xs_state = x;
    return x * 0x2545F4914F6CDD1DULL;
}
static inline double xs_uniform01(void) {
    return (xs_next_u64() >> 11) * (1.0 / 9007199254740992.0); // 53-bit mantissa
}
static inline double xs_normal01(void) {
    double u1 = xs_uniform01(); if (u1 < 1e-300) u1 = 1e-300;
    double u2 = xs_uniform01();
    return sqrt(-2.0 * log(u1)) * cos(2.0 * M_PI * u2);
}

static int cmp_dbl(const void *a, const void *b) {
    double da = *(const double*)a, db = *(const double*)b;
    return (da < db) ? -1 : (da > db);
}

int main(int argc, char **argv) {
    int N     = (argc > 1) ? atoi(argv[1]) : 1000000; // scenarios
    int steps = (argc > 2) ? atoi(argv[2]) : 1;       // time steps
    double alpha = (argc > 3) ? atof(argv[3]) : 0.99; // quantile

    const double S0 = 100.0;
    const double mu = 0.05;
    const double sigma = 0.20;
    const double T = (double)steps / 252.0;
    const double dt = T / (double)steps;

    double *loss = (double*)malloc((size_t)N * sizeof(double));
    if (!loss) { fprintf(stderr, "alloc failed\n"); return 1; }

    xs_seed(123456789u);

    struct timespec start, end;
    clock_gettime(CLOCK_MONOTONIC, &start);

    for (int i = 0; i < N; i++) {
        double S = S0;
        for (int k = 0; k < steps; k++) {
            double z = xs_normal01();
            double drift = (mu - 0.5 * sigma * sigma) * dt;
            double diff  = sigma * sqrt(dt) * z;
            S *= exp(drift + diff);
        }
        double pnl = S - S0;
        loss[i] = -pnl; // losses as positive values
    }

    qsort(loss, (size_t)N, sizeof(double), cmp_dbl);
    int idx = (int)((1.0 - alpha) * N);
    if (idx < 0) idx = 0;
    if (idx >= N) idx = N - 1;
    double var = loss[N - 1 - idx];

    clock_gettime(CLOCK_MONOTONIC, &end);
    long long time_ns = (end.tv_sec - start.tv_sec) * 1000000000LL + (end.tv_nsec - start.tv_nsec);

    printf("TASK=var_mc,N=%d,TIME_NS=%lld,VAR=%.6f\n", N, time_ns, var);
    free(loss);
    return 0;
}