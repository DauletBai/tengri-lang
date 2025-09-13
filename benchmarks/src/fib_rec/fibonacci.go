// benchmarks/src/fib_rec/fibonacci.go

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func main() {
	n := 34 // Значение по умолчанию
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			n = v
		}
	}

	reps := 1 // По умолчанию 1 прогон для рекурсии
	if rs := os.Getenv("BENCH_REPS"); rs != "" {
		if v, err := strconv.Atoi(rs); err == nil && v > 0 {
			reps = v
		}
	}

	// Короткий прогрев, чтобы JIT-компилятор Go "проснулся"
	_ = fib(10)

	start := time.Now()
	var res int
	for i := 0; i < reps; i++ {
		res = fib(n)
	}
	elapsed := time.Since(start)

	// Усредняем время на один вызов
	perCallNs := float64(elapsed.Nanoseconds()) / float64(reps)

	fmt.Printf("RESULT: %d\n", res)
	fmt.Printf("TIME_NS: %.0f\n", perCallNs)
}