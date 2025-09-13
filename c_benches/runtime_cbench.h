#ifndef TGR_CBENCH_RUNTIME_H
#define TGR_CBENCH_RUNTIME_H

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>

#ifdef __APPLE__
  #include <mach/mach_time.h>
  static inline int64_t tgr_time_ns(void) {
    static mach_timebase_info_data_t info = {0};
    if (info.denom == 0) mach_timebase_info(&info);
    uint64_t t = mach_absolute_time();
    return (int64_t)((t * info.numer) / info.denom);
  }
#else
  #include <time.h>
  static inline int64_t tgr_time_ns(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (int64_t)ts.tv_sec * 1000000000LL + ts.tv_nsec;
  }
#endif

static inline long tgr_argi(int argc, char** argv, int idx, long defv) {
  if (argc > idx) return strtol(argv[idx], NULL, 10);
  return defv;
}

static inline int64_t tgr_bench_reps(void) {
  const char* s = getenv("BENCH_REPS");
  if (!s || !*s) return 1;
  // поддержка подчеркиваний: "10_000_000"
  char buf[64]; size_t j = 0;
  for (size_t i = 0; s[i] && j+1 < sizeof(buf); ++i) if (s[i] != '_') buf[j++] = s[i];
  buf[j] = 0;
  return strtoll(buf, NULL, 10);
}

static inline void tgr_print_long(long v) {
  printf("RESULT: %ld\n", v);
}

static inline void tgr_print_time_ns(int64_t per_ns) {
  printf("TIME_NS: %lld\n", (long long)per_ns);
}

#endif // TGR_CBENCH_RUNTIME_H