#include "runtime.h"
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#ifdef __APPLE__
  #include <mach/mach_time.h>
#elif defined(_WIN32)
  #define WIN32_LEAN_AND_MEAN
  #include <windows.h>
#else
  #include <time.h>
#endif

int __argc = 0;
char **__argv = NULL;

long print(long x) {
    printf("%ld\n", x);
    return x;
}

long argi(long idx) {
    if (idx < 0 || idx >= __argc) return 0;
    char *end = NULL;
    long v = strtol(__argv[idx], &end, 10);
    return (end && *end == '\0') ? v : 0;
}

long time_ns(void) {
#ifdef __APPLE__
    static mach_timebase_info_data_t info = {0,0};
    if (info.denom == 0) { mach_timebase_info(&info); }
    uint64_t t = mach_absolute_time();
    uint64_t ns = (t * info.numer) / info.denom;
    return (long)ns;
#elif defined(_WIN32)
    static LARGE_INTEGER freq = {0};
    if (freq.QuadPart == 0) {
        QueryPerformanceFrequency(&freq);
    }
    LARGE_INTEGER now;
    QueryPerformanceCounter(&now);
    long double seconds = (long double)now.QuadPart / (long double)freq.QuadPart;
    uint64_t ns = (uint64_t)(seconds * 1000000000.0L);
    return (long)ns;
#else
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    uint64_t ns = (uint64_t)ts.tv_sec * 1000000000ull + (uint64_t)ts.tv_nsec;
    return (long)ns;
#endif
}

long print_time_ns(long ns) {
    printf("TIME_NS: %ld\n", ns);
    return ns;
}