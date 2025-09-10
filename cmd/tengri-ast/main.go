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

func main() {
	env := object.NewEnvironment()

	if len(os.Args) > 1 {
		if err := runFile(os.Args[1], env); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		return
	}

	repl(os.Stdin, os.Stdout, env)
}

func runFile(path string, env *object.Environment) error {
	src, err := os.ReadFile(path)
	if err != nil { return err }
	l := lexer.New(string(src))
	p := parser.New(l)
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		for _, e := range errs { fmt.Fprintln(os.Stderr, "parser error:", e) }
		return fmt.Errorf("parse failed")
	}
	val := evaluator.Eval(prog, env)
	if val != nil { fmt.Println(val.Inspect()) }
	return nil
}

func repl(in io.Reader, out io.Writer, env *object.Environment) {
	sc := bufio.NewScanner(in)
	fmt.Fprintln(out, "Tengri REPL (AST) — Ctrl+D to exit")
	for {
		fmt.Fprint(out, ">> ")
		if !sc.Scan() { break }
		line := sc.Text()
		l := lexer.New(line)
		p := parser.New(l)
		prog := p.ParseProgram()
		if errs := p.Errors(); len(errs) > 0 {
			for _, e := range errs { fmt.Fprintln(out, "parser error:", e) }
			continue
		}
		val := evaluator.Eval(prog, env)
		if val != nil { fmt.Fprintln(out, val.Inspect()) }
	}
}