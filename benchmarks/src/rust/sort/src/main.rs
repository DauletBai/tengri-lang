use std::env;
use std::time::Instant;

fn quicksort<T: Ord>(arr: &mut [T]) {
    if arr.len() <= 1 {
        return;
    }
    let pivot_index = arr.len() / 2;
    arr.swap(0, pivot_index);
    let mut i = 1;
    for j in 1..arr.len() {
        if arr[j] < arr[0] {
            arr.swap(i, j);
            i += 1;
        }
    }
    arr.swap(0, i - 1);
    quicksort(&mut arr[0..i - 1]);
    quicksort(&mut arr[i..]);
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let n: usize = if args.len() > 1 {
        args[1].parse().unwrap_or(100)
    } else {
        100
    };

    let mut numbers: Vec<i64> = (1..=n as i64).collect();
    
    // Warm-up run
    let mut warm_up_numbers = numbers.clone();
    quicksort(&mut warm_up_numbers);

    let start = Instant::now();
    quicksort(&mut numbers);
    let duration = start.elapsed();

    // CORRECTED: Output format now matches the unified runtime
    println!("TASK=sort_rs,N={},TIME_NS={}", n, duration.as_nanos());
    // Print result to stderr to prevent compiler from optimizing it away
    eprintln!("Result (last element): {}", numbers[n - 1]);
}