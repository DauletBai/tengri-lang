// FILE: benchmarks/src/go/fib_rec.go
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func fibRec(n int) int {
	if n < 2 {
		return n
	}
	return fibRec(n-1) + fibRec(n-2)
}

func main() {
	n := 34
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			n = v
		}
	}

	_ = fibRec(10) // Warm-up

	start := time.Now()
	result := fibRec(n)
	elapsed := time.Since(start)

	// Output format matching the unified runtime
	fmt.Printf("TASK=fib_rec_go,N=%d,TIME_NS=%d\n", n, elapsed.Nanoseconds())
	fmt.Fprintf(os.Stderr, "Result: %d\n", result) // Prevent optimization
}