use std::env;
use std::time::Instant;

fn fib_iter(n: i64) -> i64 {
    if n <= 1 { return n; }
    let (mut a, mut b) = (0_i64, 1_i64);
    for _ in 2..=n {
        let c = a + b;
        a = b; b = c;
    }
    b
}

fn parse_arg_or(idx: usize, default_v: i64) -> i64 {
    std::env::args().nth(idx).and_then(|s| s.parse::<i64>().ok()).unwrap_or(default_v)
}

fn bench_reps() -> i64 {
    let raw = env::var("BENCH_REPS").unwrap_or_default();
    if raw.is_empty() { return 1; }
    let filtered: String = raw.chars().filter(|&c| c != '_').collect();
    filtered.parse::<i64>().unwrap_or(1)
}

fn main() {
    let n = parse_arg_or(1, 40);
    let reps = bench_reps().max(1);

    // warm-up
    let _ = fib_iter(10);

    let t0 = Instant::now();
    let mut res = 0_i64;
    for _ in 0..reps {
        res = fib_iter(n);
    }
    let elapsed = t0.elapsed();
    let total_ns = elapsed.as_secs_f64() * 1e9;
    let per_ns = (total_ns / reps as f64).round() as i64;

    println!("RESULT: {}", res);
    println!("TIME_NS: {}", per_ns);
}
