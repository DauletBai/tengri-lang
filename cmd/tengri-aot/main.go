package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DauletBai/tengri-lang/internal/aotminic"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <input.tgr> -o <out.c>\n", filepath.Base(os.Args[0]))
	os.Exit(2)
}

func main() {
	if len(os.Args) < 4 {
		usage()
	}
	inPath := os.Args[1]
	if os.Args[2] != "-o" {
		usage()
	}
	outPath := os.Args[3]

	src, err := os.ReadFile(inPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	out, err := aotminic.TranspileToC(string(src), filepath.Base(inPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outPath, []byte(out), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("C emitted:", outPath)
}