#include "runtime_cbench.h"

static long fib_iter(long n) {
  if (n <= 1) return n;
  long a = 0, b = 1;
  for (long i = 2; i <= n; ++i) {
    long c = a + b;
    a = b; b = c;
  }
  return b;
}

int main(int argc, char** argv) {
  long n = tgr_argi(argc, argv, 1, 40);
  int64_t reps = tgr_bench_reps();

  // короткая разминка
  (void)fib_iter(10);

  int64_t t0 = tgr_time_ns();
  long res = 0;
  for (int64_t i = 0; i < reps; ++i) {
    res = fib_iter(n);
  }
  int64_t t1 = tgr_time_ns();

  tgr_print_long(res);
  int64_t per = (reps > 0) ? (t1 - t0) / reps : 0;
  tgr_print_time_ns(per);
  return 0;
}
