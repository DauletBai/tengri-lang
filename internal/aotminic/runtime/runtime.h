// FILE: internal/aotminic/runtime/runtime.h

#ifndef tenge_RUNTIME_H
#define tenge_RUNTIME_H

#include <stdio.h>
#include <stdlib.h>
#include <time.h>

// --- Helper functions ---
int get_n(int argc, char** argv, int default_n);
int* create_array(int n);

// --- Timing macro ---
// This macro is now simpler and expects the calling code to provide
// a single statement (which can be a do-while(0) block).
#define TIME_IT_NS(block, task_name, n)                                      \
    do {                                                                     \
        struct timespec start, end;                                          \
        clock_gettime(CLOCK_MONOTONIC, &start);                              \
        block; /* Execute the provided code block directly */                \
        clock_gettime(CLOCK_MONOTONIC, &end);                                \
        long long time_ns = (end.tv_sec - start.tv_sec) * 1000000000LL +      \
                            (end.tv_nsec - start.tv_nsec);                   \
        printf("TASK=%s,N=%d,TIME_NS=%lld\n", task_name, n, time_ns);         \
    } while (0)

#endif // tenge_RUNTIME_H