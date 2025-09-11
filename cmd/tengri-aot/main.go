package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DauletBai/tengri-lang/internal/aotminic"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] <source.tgr>\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "  -o string   Output C file path (required)")
	fmt.Fprintln(os.Stderr, "  -force      Allow any .tgr name without warnings")
}

func main() {
	out := flag.String("o", "", "Output C file path (required)")
	force := flag.Bool("force", false, "Allow any .tgr name without warnings")

	flag.Usage = usage
	flag.Parse()

	if *out == "" || flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	src := flag.Arg(0)

	if err := aotminic.TranspileToC(src, *out, *force); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "C emitted: %s\n", *out)
}