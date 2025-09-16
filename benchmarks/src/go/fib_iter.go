// FILE: benchmarks/src/go/fib_iter.go
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func fibIter(n int) int {
	if n < 2 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func main() {
	n := 90
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			n = v
		}
	}

	_ = fibIter(10) // Warm-up

	start := time.Now()
	result := fibIter(n)
	elapsed := time.Since(start)

	// Output format matching the unified runtime
	fmt.Printf("TASK=fib_iter_go,N=%d,TIME_NS=%d\n", n, elapsed.Nanoseconds())
	fmt.Fprintf(os.Stderr, "Result: %d\n", result) // Prevent optimization
}