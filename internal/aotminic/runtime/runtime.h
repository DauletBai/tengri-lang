#ifndef TENGRI_MINI_RUNTIME_H
#define TENGRI_MINI_RUNTIME_H

#include <stdio.h>
#include <time.h>
#include <stdlib.h>

static inline long time_ns(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (long)ts.tv_sec * 1000000000L + (long)ts.tv_nsec;
}

static inline void print_time_ns(long ns) {
    printf("TIME_NS: %ld\n", ns);
}

static inline void print(long v) {
    printf("%ld\n", v);
}

static inline long argi(int argc, char** argv, int idx, long defv) {
    if (argc > idx) return strtoll(argv[idx], NULL, 10);
    return defv;
}

#endif