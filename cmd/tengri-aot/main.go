package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DauletBai/tengri-lang/internal/aotminic"
)

func main() {
	var out string
	var force bool

	flag.StringVar(&out, "o", "", "Output C file path (required)")
	flag.BoolVar(&force, "force", false, "Allow any .tgr name without warnings")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <source.tgr>\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	src := flag.Arg(0)
	if out == "" {
		fmt.Fprintln(os.Stderr, "error: -o <output.c> is required")
		os.Exit(2)
	}

	// Optional sanity warning for demo names (kept for UX parity).
	if !force {
		base := filepath.Base(src)
		if !(contains(base, "fib_iter_cli") || contains(base, "fib_rec_cli")) {
			fmt.Fprintln(os.Stderr, "warning: mini-AOT demo expects fib_iter_cli / fib_rec_cli; use -force to bypass")
		}
	}

	cCode, err := aotminic.TranspileToC(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: transpile failed: %v\n", err)
		os.Exit(1)
	}

	if err := aotminic.WriteFile(out, cCode); err != nil {
		fmt.Fprintf(os.Stderr, "error: write C file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("C emitted: %s\n", out)
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (filepath.Base(s) == sub || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	// tiny helper to avoid importing strings in this file
outer:
	for i := 0; i+len(sub) <= len(s); i++ {
		for j := 0; j < len(sub); j++ {
			if s[i+j] != sub[j] {
				continue outer
			}
		}
		return i
	}
	return -1
}