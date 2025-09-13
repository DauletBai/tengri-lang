#include "runtime_cbench.h"

static long fib_rec(long n) {
  if (n <= 1) return n;
  return fib_rec(n-1) + fib_rec(n-2);
}

int main(int argc, char** argv) {
  long n = tgr_argi(argc, argv, 1, 34);
  int64_t reps = tgr_bench_reps();

  (void)fib_rec(10); // разминка

  int64_t t0 = tgr_time_ns();
  long res = 0;
  for (int64_t i = 0; i < reps; ++i) {
    res = fib_rec(n);
  }
  int64_t t1 = tgr_time_ns();

  tgr_print_long(res);
  int64_t per = (reps > 0) ? (t1 - t0) / reps : 0;
  tgr_print_time_ns(per);
  return 0;
}