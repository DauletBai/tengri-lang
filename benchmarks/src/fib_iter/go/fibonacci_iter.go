package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

const MOD = 1000000007

func fibIterMod(n int) int {
	if n < 2 {
		return n
	}
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b, (a+b)%MOD
	}
	return a
}

func pickReps(n int) int {
	// База увеличена, чтобы TIME был в мкс/нс, а не округлялся до нуля
	base := 5_000_000
	scale := int(math.Max(1, float64(50/max(1, n))))
	reps := base * scale
	if reps < 500_000 {
		reps = 500_000
	}
	return reps
}

func max(a, b int) int { if a > b { return a }; return b }

func main() {
	n := 90
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			n = v
		}
	}
	reps := pickReps(n)
	if rs := os.Getenv("BENCH_REPS"); rs != "" {
		if v, err := strconv.Atoi(rs); err == nil && v > 0 {
			reps = v
		}
	}

	start := time.Now()
	var res int
	for i := 0; i < reps; i++ {
		res = fibIterMod(n)
	}
	elapsed := time.Since(start)
	perCall := float64(elapsed) / float64(reps)

	fmt.Printf("RESULT: %d\n", res)
	fmt.Printf("TIME: %.9f\n", perCall/1e9)          
	fmt.Printf("TIME_NS: %.0f\n", perCall)         
}