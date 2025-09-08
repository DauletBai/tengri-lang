import sys, os, time, math

MOD = 1_000_000_007

def fib_iter_mod(n: int) -> int:
    if n < 2:
        return n
    a, b = 0, 1
    for _ in range(n):
        a, b = b, (a + b) % MOD
    return a

def pick_reps(n: int) -> int:
    base = 5_000_000
    scale = max(1, 50 // max(1, n))
    reps = base * scale
    return max(reps, 50_000)

if __name__ == "__main__":
    n = 90
    if len(sys.argv) > 1:
        try:
            n = int(sys.argv[1])
        except:
            pass
    reps = pick_reps(n)
    rs = os.getenv("BENCH_REPS")
    if rs:
        try:
            v = int(rs)
            if v > 0:
                reps = v
        except:
            pass

    t0 = time.perf_counter()
    res = 0
    for _ in range(reps):
        res = fib_iter_mod(n)
    t1 = time.perf_counter()
    per_call = (t1 - t0) / reps

    print(f"RESULT: {res}")
    print(f"TIME: {per_call:.6f}")