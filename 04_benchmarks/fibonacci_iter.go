//go:build iter

package main

import (
	"fmt"
	"os"
	"strconv"
)

func fibIter(n int) int {
	if n < 2 {
		return n
	}
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b, a+b
	}
	return a
}

func main() {
	n := 300
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			n = v
		}
	}
	fmt.Println(fibIter(n))
}