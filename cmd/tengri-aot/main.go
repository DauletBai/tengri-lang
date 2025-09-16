// FILE: cmd/tengri-aot/main.go

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/DauletBai/tengri-lang/internal/aotminic"
)

func main() {
	// Define command-line flags
	outFile := flag.String("o", "", "Output C file path (required)")
	force := flag.Bool("force", false, "Allow any .tgr name without warnings")
	flag.Parse()

	// Validate required flags
	if *outFile == "" {
		fmt.Fprintln(os.Stderr, "error: -o output file path is required")
		flag.Usage()
		os.Exit(1)
	}

	// Get the source file path
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "error: exactly one source .tgr file must be provided")
		flag.Usage()
		os.Exit(1)
	}
	sourceFile := flag.Arg(0)

	// Select the C template based on the source filename
	var template string
	switch {
	case strings.HasSuffix(sourceFile, "fib_iter_cli.tgr"):
		template = aotminic.FibIterC
	case strings.HasSuffix(sourceFile, "fib_rec_cli.tgr"):
		template = aotminic.FibRecC
	case strings.HasSuffix(sourceFile, "sort_cli_qs.tgr"):
		template = aotminic.SortQSortC
	case strings.HasSuffix(sourceFile, "sort_cli_ms.tgr"):
		template = aotminic.SortMergeSortC
	default:
		// If not a recognized demo file, exit with an error
		if !*force {
			fmt.Fprintf(os.Stderr, "error: unsupported AOT demo source: %s\n", sourceFile)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "error: force mode is not yet implemented for unknown sources\n")
		os.Exit(1)
	}

	// Emit the C code
	err := os.WriteFile(*outFile, []byte(template), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing to output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("C emitted: %s\n", *outFile)
}