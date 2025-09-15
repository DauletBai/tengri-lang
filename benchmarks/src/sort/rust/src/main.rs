// FILE: benchmarks/src/sort/rust/src/main.rs
// Purpose: Sort benchmark with nanosecond timer and unified REPORT line.

use std::env;
use std::time::Instant;

fn getenv_usize(key: &str, defv: usize) -> usize {
    env::var(key).ok()
        .and_then(|s| s.parse::<usize>().ok())
        .filter(|&n| n > 0)
        .unwrap_or(defv)
}
fn getenv_u64(key: &str, defv: u64) -> u64 {
    env::var(key).ok()
        .and_then(|s| s.parse::<u64>().ok())
        .filter(|&n| n > 0)
        .unwrap_or(defv)
}

fn main() {
    let n = getenv_usize("SIZE", 100_000);
    let reps = getenv_u64("BENCH_REPS", 3);

    let mut v: Vec<i64> = (0..n as i64).map(|i| n as i64 - i).collect();
    // warm-up
    v.sort();

    let t0 = Instant::now();
    for _ in 0..reps {
        for i in 0..n { v[i] = (n - i) as i64; }
        v.sort();
    }
    let total_ns = t0.elapsed().as_nanos() as u128;
    let avg_ns = total_ns / (reps as u128);

    let first = v[0];
    let last = v[n-1];
    let sum: i128 = v.iter().map(|&x| x as i128).sum();

    println!(
        "REPORT impl=rust task=sort n={} reps={} time_ns_avg={} first={} last={} sum={}",
        n, reps, avg_ns, first, last, sum
    );
}