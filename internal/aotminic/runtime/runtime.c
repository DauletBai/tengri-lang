#include "runtime.h"

#include <stdint.h>
#include <string.h>
#include <time.h>

#ifdef __APPLE__
#include <mach/mach_time.h>
static mach_timebase_info_data_t g_timebase = {0, 0};
int64_t tgr_time_ns(void) {
    if (g_timebase.denom == 0) {
        (void)mach_timebase_info(&g_timebase);
    }
    uint64_t t = mach_absolute_time();
    // Convert to ns
    return (int64_t)((t * (uint64_t)g_timebase.numer) / (uint64_t)g_timebase.denom);
}
#else
int64_t tgr_time_ns(void) {
    struct timespec ts;
    // CLOCK_MONOTONIC или MONOTONIC_RAW — оба дают ns (если доступны)
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (int64_t)ts.tv_sec * 1000000000LL + (int64_t)ts.tv_nsec;
}
#endif

int64_t tgr_bench_reps(void) {
    const char* s = getenv("BENCH_REPS");
    if (!s || !*s) return 1;
    long long v = atoll(s);
    if (v < 1) v = 1;
    return (int64_t)v;
}

void tgr_print_time_ns(int64_t elapsed_ns, int64_t reps) {
    if (reps < 1) reps = 1;
    // Усредняем в int64, без потери точности
    int64_t avg = elapsed_ns / reps;
    // benchfast ожидает строки ровно в таком формате:
    //   TIME_NS: <number>
    // (дополнительно можно печатать и "TIME: seconds" при желании)
    printf("TIME_NS: %lld\n", (long long)avg);
}

long tgr_argi(int argc, char** argv, int index, long def) {
    if (index < argc && index >= 1) {
        char* endp = NULL;
        long v = strtol(argv[index], &endp, 10);
        if (endp && *endp == '\0') return v;
    }
    return def;
}

void tgr_print_long(long v) {
    printf("%ld\n", v);
}