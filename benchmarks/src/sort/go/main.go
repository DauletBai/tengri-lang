// FILE: benchmarks/src/sort/go/main.go
// Purpose: Sort benchmark with nanosecond timer and unified REPORT line.

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"sort"
)

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

func getenvInt64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			return n
		}
	}
	return def
}

func main() {
	n := getenvInt("SIZE", 100000)
	reps := getenvInt64("BENCH_REPS", 3)

	arr := make([]int, n)
	// warm-up
	for i := 0; i < n; i++ { arr[i] = n - i }
	sort.Ints(arr)

	start := time.Now()
	for r := int64(0); r < reps; r++ {
		for i := 0; i < n; i++ { arr[i] = n - i }
		sort.Ints(arr)
	}
	dur := time.Since(start).Nanoseconds()
	avg := dur
	if reps > 0 { avg = dur / reps }

	sum := 0
	for _, x := range arr { sum += x }
	first := arr[0]
	last := arr[len(arr)-1]

	fmt.Printf("REPORT impl=go task=sort n=%d reps=%d time_ns_avg=%d first=%d last=%d sum=%d\n",
		n, reps, avg, first, last, sum)
}