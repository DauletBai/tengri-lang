// 03_compiler_go/main.go
package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"tengri-lang/03_compiler_go/evaluator"
	"tengri-lang/03_compiler_go/lexer"
	"tengri-lang/03_compiler_go/object"
	"tengri-lang/03_compiler_go/parser"
)

func loadScriptFromEnv() (string, bool) {
    path := os.Getenv("TENGRI_SCRIPT")
    if path == "" {
        return "", false
    }
    b, err := os.ReadFile(path)
    if err != nil {
        fmt.Fprintf(os.Stderr, "cannot read %s: %v\n", path, err)
        return "", false
    }
    return string(b), true
}

func main() {
	// Читаем код из файла для бенчмарка
	inputBytes, err := ioutil.ReadFile("../04_benchmarks/fibonacci.tengri")
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %s\n", err)
		return
	}
	input := string(inputBytes)

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

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil {
		fmt.Println(result.Inspect()) // Выводим результат вычислений
	}
}