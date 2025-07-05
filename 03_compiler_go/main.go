// 03_compiler_go/main.go
package main

import (
	"fmt"
	"tengri-lang/03_compiler_go/lexer"
	"tengri-lang/03_compiler_go/parser"
)

func main() {
	input := `
        Π qosw (□ a, □ b) (
            → a + b
        )
        Λ x : qosw(5, 10)
    `

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("Ошибки парсера:")
		for _, msg := range p.Errors() {
			fmt.Println("\t" + msg)
		}
		return
	}

	fmt.Println("--- Древо Мысли (AST) ---")
	fmt.Println(program.String())
}