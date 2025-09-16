// FILE: benchmarks/src/go/sort/main.go
package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	n := 100000
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			n = v
		}
	}

	numbers := make([]int, n)
	for i := 0; i < n; i++ {
		numbers[i] = i + 1
	}

	// A short warm-up run
	warmUp := make([]int, 100)
	sort.Ints(warmUp)

	start := time.Now()
	sort.Ints(numbers)
	elapsed := time.Since(start)

	// Output format matching the unified runtime
	fmt.Printf("TASK=sort_go,N=%d,TIME_NS=%d\n", n, elapsed.Nanoseconds())
	// Print result to stderr to prevent compiler from optimizing it away
	fmt.Fprintf(os.Stderr, "Result (last element): %d\n", numbers[n-1])
}