package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Case struct {
	Name string
	Cmd  []string
	Env  []string
	Dir  string
}

func run(c Case) (time.Duration, string, error) {
	cmd := exec.Command(c.Cmd[0], c.Cmd[1:]...)
	if c.Dir != "" {
		cmd.Dir = c.Dir
	}
	if len(c.Env) > 0 {
		cmd.Env = append(os.Environ(), c.Env...)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	start := time.Now()
	err := cmd.Run()
	dur := time.Since(start)
	return dur, strings.TrimSpace(out.String()), err
}

func writeTengriScript(n int) (string, error) {
	s := fmt.Sprintf(`let fib = fn(n) {
  if (n < 2) { return n; }
  return fib(n - 1) + fib(n - 2);
};

let n = %d;
fib(n);
`, n)
	tmp := filepath.Join("03_compiler_go", ".tmp_fib_ascii.tg")
	if err := os.WriteFile(tmp, []byte(s), 0o644); err != nil {
		return "", err
	}
	return tmp, nil
}

func main() {
	Ns := []int{30, 32, 34}

	repoRoot, _ := os.Getwd()
	tengriDir := filepath.Join(repoRoot, "03_compiler_go")

	for _, N := range Ns {
		fmt.Printf("\nN = %d\n", N)
		fmt.Println("──────────────────────────────────────────────")
		fmt.Printf("%-22s  %-10s  %s\n", "Target", "Time", "Output")
		fmt.Println(strings.Repeat("-", 66))

		// Подготовим скрипт для интерпретатора
		scriptPath, err := writeTengriScript(N)
		if err != nil {
			fmt.Printf("%-22s  %8.3fs  [ERR] %v\n", "prep-script", 0.0, err)
			continue
		}

		tests := []Case{
			{
				Name: "Go native",
				Cmd:  []string{"go", "run", "04_benchmarks/fibonacci_arg.go", fmt.Sprintf("%d", N)},
				Dir:  repoRoot,
			},
			{
				Name: "Tengri (AST interp)",
				Cmd:  []string{"go", "run", "."},
				Env:  []string{"TENGRI_SCRIPT=" + scriptPath},
				Dir:  tengriDir, // ВАЖНО: запускать из 03_compiler_go
			},
			{
				Name: "Python",
				Cmd:  []string{"python3", "04_benchmarks/fibonacci_arg.py", fmt.Sprintf("%d", N)},
				Dir:  repoRoot,
			},
			{
				Name: "VM (bytecode mini)",
				Cmd:  []string{"go", "run", "05_vm_mini/main.go", fmt.Sprintf("%d", N)},
				Dir:  repoRoot,
			},
		}

		for _, t := range tests {
			d, out, err := run(t)
			status := "OK"
			if err != nil {
				status = "ERR"
			}
			lines := strings.Split(out, "\n")
			last := ""
			if len(lines) > 0 {
				last = lines[len(lines)-1]
			}
			fmt.Printf("%-22s  %8.3fs  [%s] %s\n", t.Name, d.Seconds(), status, last)
		}
	}
}