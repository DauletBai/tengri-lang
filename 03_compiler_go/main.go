// 03_compiler_go/main.go
package main

import (
	"fmt"
	"tengri-lang/03_compiler_go/lexer"
	"tengri-lang/03_compiler_go/parser"
	"tengri-lang/03_compiler_go/evaluator"
	"tengri-lang/03_compiler_go/object"
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

	fmt.Println("--- AST ---")
	fmt.Println(program.String())

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	fmt.Println("--- Результат ---")
	if result != nil {
		fmt.Println(result.Inspect())
	} else {
		fmt.Println("null")
	}
}
