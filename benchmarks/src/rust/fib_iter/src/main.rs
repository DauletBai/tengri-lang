use std::env;
use std::time::Instant;

fn fib_iter(n: u64) -> u64 {
    if n < 2 {
        return n;
    }
    let mut a = 0;
    let mut b = 1;
    for _ in 2..=n {
        let temp = a + b;
        a = b;
        b = temp;
    }
    b
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let n: u64 = if args.len() > 1 {
        args[1].parse().unwrap_or(90)
    } else {
        90
    };

    // Warm-up
    fib_iter(10);

    let start = Instant::now();
    let result = fib_iter(n);
    let duration = start.elapsed();

    // CORRECTED: Output format now matches the unified runtime
    println!("TASK=fib_iter_rs,N={},TIME_NS={}", n, duration.as_nanos());
    eprintln!("Result: {}", result);
}