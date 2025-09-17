// FILE: benchmarks/src/rust/var_mc/src/main.rs
// Purpose: Monte Carlo VaR benchmark (GBM, Box–Muller, xorshift64*).

use std::env;
use std::time::Instant;

struct Rng { s: u64 }
impl Rng {
    fn seed(seed: u64) -> Self {
        Self { s: if seed == 0 { 0x9E3779B97F4A7C15 } else { seed } }
    }
    fn u64(&mut self) -> u64 {
        let mut x = self.s;
        x ^= x >> 12;
        x ^= x << 25;
        x ^= x >> 27;
        self.s = x;
        x.wrapping_mul(0x2545F4914F6CDD1D)
    }
    fn uniform(&mut self) -> f64 {
        ((self.u64() >> 11) as f64) * (1.0 / 9007199254740992.0)
    }
    fn normal(&mut self) -> f64 {
        let mut u1 = self.uniform();
        if u1 < 1e-300 { u1 = 1e-300; }
        let u2 = self.uniform();
        (-2.0 * u1.ln()).sqrt() * (2.0 * std::f64::consts::PI * u2).cos()
    }
}

fn main() {
    let mut args = env::args().skip(1);
    let n: usize = args.next().and_then(|s| s.parse().ok()).unwrap_or(1_000_000);
    let steps: usize = args.next().and_then(|s| s.parse().ok()).unwrap_or(1);
    let alpha: f64 = args.next().and_then(|s| s.parse().ok()).unwrap_or(0.99);

    let s0 = 100.0_f64;
    let mu = 0.05_f64;
    let sigma = 0.20_f64;
    let t = (steps as f64) / 252.0;
    let dt = t / (steps as f64);

    let mut loss = vec![0.0_f64; n];
    let mut r = Rng::seed(123456789);

    let start = Instant::now();
    for i in 0..n {
        let mut s = s0;
        for _ in 0..steps {
            let z = r.normal();
            let drift = (mu - 0.5 * sigma * sigma) * dt;
            let diff  = sigma * dt.sqrt() * z;
            s *= (drift + diff).exp();
        }
        let pnl = s - s0;
        loss[i] = -pnl;
    }
    loss.sort_by(|a,b| a.partial_cmp(b).unwrap());

    let mut idx = ((1.0 - alpha) * (n as f64)) as isize;
    if idx < 0 { idx = 0; }
    if idx >= n as isize { idx = n as isize - 1; }
    let var = loss[n - 1 - (idx as usize)];

    let elapsed = start.elapsed();
    let time_ns = (elapsed.as_secs() as u128) * 1_000_000_000u128 + (elapsed.subsec_nanos() as u128);
    println!("TASK=var_mc,N={},TIME_NS={},VAR={:.6}", n, time_ns, var);
}