#ifndef TENGRI_RUNTIME_H
#define TENGRI_RUNTIME_H

#include <stddef.h>

/* CLI helpers */
long argi(int argc, char **argv, int index, long def);

/* I/O */
void print(long long v);

/* Timing (nanoseconds, monotonic) */
long long time_ns(void);
void print_time_ns(long long ns);

/*
 * These are only declarations to silence implicit-prototype warnings.
 * Definitions for fib_* come from the transpiled C (emitted by the AOT tool).
 */
long long fib_iter(long long n);
long long fib_rec(long long n);

#endif /* TENGRI_RUNTIME_H */