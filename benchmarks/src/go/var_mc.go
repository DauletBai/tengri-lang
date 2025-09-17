// FILE: benchmarks/src/go/var_mc.go
// Purpose: Monte Carlo VaR benchmark (GBM, Box–Muller, xorshift64*).
package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
)

type rng struct{ s uint64 }

func (r *rng) Seed(seed uint64) {
	if seed == 0 {
		seed = 0x9E3779B97F4A7C15
	}
	r.s = seed
}
func (r *rng) U64() uint64 {
	x := r.s
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	r.s = x
	return x * 0x2545F4914F6CDD1D
}
func (r *rng) Uniform() float64 {
	return float64(r.U64()>>11) * (1.0 / 9007199254740992.0)
}
func (r *rng) Normal() float64 {
	u1 := r.Uniform()
	if u1 < 1e-300 {
		u1 = 1e-300
	}
	u2 := r.Uniform()
	return math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
}

func main() {
	N := 1000000
	steps := 1
	alpha := 0.99
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			N = v
		}
	}
	if len(os.Args) > 2 {
		if v, err := strconv.Atoi(os.Args[2]); err == nil {
			steps = v
		}
	}
	if len(os.Args) > 3 {
		if v, err := strconv.ParseFloat(os.Args[3], 64); err == nil {
			alpha = v
		}
	}

	S0, mu, sigma := 100.0, 0.05, 0.20
	T := float64(steps) / 252.0
	dt := T / float64(steps)

	loss := make([]float64, N)
	var r rng
	r.Seed(123456789)

	start := time.Now()
	for i := 0; i < N; i++ {
		S := S0
		for k := 0; k < steps; k++ {
			z := r.Normal()
			drift := (mu - 0.5*sigma*sigma) * dt
			diff := sigma * math.Sqrt(dt) * z
			S *= math.Exp(drift + diff)
		}
		pnl := S - S0
		loss[i] = -pnl
	}
	sort.Float64s(loss)
	idx := int((1.0 - alpha) * float64(N))
	if idx < 0 {
		idx = 0
	}
	if idx >= N {
		idx = N - 1
	}
	vaR := loss[N-1-idx]
	elapsed := time.Since(start)
	fmt.Printf("TASK=var_mc,N=%d,TIME_NS=%d,VAR=%.6f\n", N, elapsed.Nanoseconds(), vaR)
}