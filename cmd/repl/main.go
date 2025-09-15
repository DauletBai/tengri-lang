// FILE: cmd/repl/main.go

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/DauletBai/tengri-lang/internal/lang/evaluator"
	"github.com/DauletBai/tengri-lang/internal/lang/lexer"
	"github.com/DauletBai/tengri-lang/internal/lang/object"
	"github.com/DauletBai/tengri-lang/internal/lang/parser"
)

// PROMPT is the string that appears at the start of each REPL line.
const PROMPT = ">> "

// Start initializes the Read-Eval-Print Loop.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment() // The environment persists across inputs.

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return // Exit on EOF (Ctrl+D)
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

// printParserErrors prints any errors that occurred during parsing.
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func main() {
	// Greet the user.
	fmt.Printf("Tengri Language REPL [v0.1]\n")
	// Start the REPL with standard input and output.
	Start(os.Stdin, os.Stdout)
}