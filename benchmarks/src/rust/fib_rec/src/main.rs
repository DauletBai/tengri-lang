use std::env;
use std::time::Instant;

fn fib_rec(n: u64) -> u64 {
    if n < 2 {
        return n;
    }
    fib_rec(n - 1) + fib_rec(n - 2)
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let n: u64 = if args.len() > 1 {
        args[1].parse().unwrap_or(35)
    } else {
        35
    };

    // Warm-up
    fib_rec(10);

    let start = Instant::now();
    let result = fib_rec(n);
    let duration = start.elapsed();

    // CORRECTED: Output format now matches the unified runtime
    println!("TASK=fib_rec_rs,N={},TIME_NS={}", n, duration.as_nanos());
    eprintln!("Result: {}", result);
}