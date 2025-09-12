#ifndef TENGRI_RUNTIME_H
#define TENGRI_RUNTIME_H

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

// Monotonic high-resolution timestamp in nanoseconds.
int64_t tgr_time_ns(void);

// Read BENCH_REPS from env (default 1, clamp >= 1).
int64_t tgr_bench_reps(void);

// Print timing as required by benchfast.
// elapsed_ns — полное время одного прогона (или суммарное) в наносекундах.
// Если вы подаете суммарное за REP прогона, функция усреднит.
void tgr_print_time_ns(int64_t elapsed_ns, int64_t reps);

// Helpers used by transpiled code
long tgr_argi(int argc, char** argv, int index, long def);
void tgr_print_long(long v);

#ifdef __cplusplus
}
#endif

#endif // TENGRI_RUNTIME_H