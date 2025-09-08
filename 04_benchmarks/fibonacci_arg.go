//go:build arg

package main

import (
	"fmt"
	"os"
	"strconv"
)

func fib(n int) int {
	if n < 2 { return n }
	return fib(n-1) + fib(n-2)
}

func main() {
	n := 35
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			n = v
		}
	}
	fmt.Println(fib(n))
}