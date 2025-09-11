#include "runtime.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

long argi(int argc, char **argv, int index, long def) {
    if (index < argc) {
        char *end = NULL;
        long val = strtol(argv[index], &end, 10);
        if (end != argv[index]) return val;
    }
    return def;
}

void print(long long v) {
    /* Print the numeric result on its own line */
    printf("%lld\n", v);
}

/* On macOS / Linux we use CLOCK_MONOTONIC for stable measurement */
long long time_ns(void) {
    struct timespec ts;
#ifdef CLOCK_MONOTONIC
    clock_gettime(CLOCK_MONOTONIC, &ts);
#else
    clock_gettime(CLOCK_REALTIME, &ts);
#endif
    return (long long)ts.tv_sec * 1000000000LL + (long long)ts.tv_nsec;
}

void print_time_ns(long long ns) {
    printf("TIME_NS: %lld\n", ns);
}