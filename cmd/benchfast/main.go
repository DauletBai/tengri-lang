// tool/benchfast/main.go
package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Target struct {
	Name string
	Run  func(task string, n int) (string, time.Duration, error)
}

var (
	reTimeNS = regexp.MustCompile(`TIME_NS[:\s]+(\d+)`)
	reTimeS  = regexp.MustCompile(`TIME[:\s]+([0-9.]+)`)
)

func main() {
	plot := flag.Bool("plot", false, "make plots (no-op stub)")
	flag.Parse()

	// наборы
	tasks := []string{"fib_rec", "fib_iter"}
	NsRec := []int{30, 32, 34}
	NsIter := []int{40, 60, 90}

	// REPS для внутренних измерений
	repsIter := 2000000
	repsRec := 2000

	root, _ := os.Getwd()
	runID := time.Now().Format("20060102-150405")
	outDir := filepath.Join(root, "benchmarks", "runs", runID)
	latestDir := filepath.Join(root, "benchmarks", "latest")
	_ = os.MkdirAll(filepath.Join(outDir, "results"), 0755)
	_ = os.MkdirAll(filepath.Join(outDir, "plots"), 0755)
	_ = os.MkdirAll(filepath.Join(latestDir, "results"), 0755)
	_ = os.MkdirAll(filepath.Join(latestDir, "plots"), 0755)
	_ = os.MkdirAll(".bin", 0755)

	// подготовка бинарников один раз
	prepGo()
	prepVM()
	prepAOT()

	targets := []Target{
		{"go", runGoBin},
		{"tengri-ast", runTengriAST},
		{"python", runPy},
		{"vm", runVMBin},
		{"tengri-aot", func(task string, n int) (string, time.Duration, error) {
			if task == "fib_iter" {
				return runAOTBinIter(n, repsIter)
			}
			return runAOTBinRec(n, repsRec)
		}},
	}

	for _, task := range tasks {
		var Ns []int
		if task == "fib_rec" {
			Ns = NsRec
		} else {
			Ns = NsIter
		}

		records := [][]string{{"task", "target", "N", "time_s", "output"}}
		for _, n := range Ns {
			fmt.Printf("\nTask = %s, N = %d\n", task, n)
			fmt.Println("──────────────────────────────────────────────────────────────")
			fmt.Printf("%-20s %-10s  %s\n", "Target", "Time (s)", "Output")
			fmt.Println("------------------------------------------------------------------")
			for _, t := range targets {
				out, dur, err := t.Run(task, n)
				sec := preferInternalTime(out, dur)
				status := markStatus(t.Name, out, err)
				if err != nil && out == "" {
					out = err.Error()
				}
				fmt.Printf("%-20s %-10.6f  [%s] %s\n", t.Name, sec, status, summarize(out))
				records = append(records, []string{task, t.Name, strconv.Itoa(n), fmt.Sprintf("%.6f", sec), strings.TrimSpace(out)})
			}
		}
		writeCSV(filepath.Join(outDir, "results", task+".csv"), records)
		writeCSV(filepath.Join(latestDir, "results", task+".csv"), records)
		fmt.Println("CSV saved:",
			filepath.Join(outDir, "results", task+".csv"),
			"and",
			filepath.Join(latestDir, "results", task+".csv"))
	}

	if *plot {
		_ = plotWithGo(filepath.Join(outDir, "plots"))
		_ = plotWithGo(filepath.Join(latestDir, "plots"))
	}
}

/* ------------ helpers ------------- */

func has(path string) bool { _, err := os.Stat(path); return err == nil }

func summarize(s string) string {
	ss := strings.TrimSpace(s)
	if len(ss) > 64 { return ss[:61] + "..." }
	return ss
}

func writeCSV(path string, rows [][]string) {
	f, err := os.Create(path)
	if err != nil { fmt.Println("writeCSV:", err); return }
	defer f.Close()
	w := csv.NewWriter(f); _ = w.WriteAll(rows); w.Flush()
}

func sh(cmd string) (string, time.Duration, error) {
	start := time.Now()
	c := exec.Command("bash", "-lc", cmd)
	var buf bytes.Buffer
	c.Stdout, c.Stderr = &buf, &buf
	err := c.Run()
	return buf.String(), time.Since(start), err
}

func run(cmd string, args ...string) (string, time.Duration, error) {
	start := time.Now()
	c := exec.Command(cmd, args...)
	var buf bytes.Buffer
	c.Stdout, c.Stderr = &buf, &buf
	err := c.Run()
	return buf.String(), time.Since(start), err
}

// если в выводе есть TIME_NS/TIME — используем его; иначе wall time
func preferInternalTime(out string, wall time.Duration) float64 {
	if m := reTimeNS.FindStringSubmatch(out); len(m) == 2 {
		if ns, err := strconv.ParseFloat(m[1], 64); err == nil {
			return ns / 1e9
		}
	}
	if m := reTimeS.FindStringSubmatch(out); len(m) == 2 {
		if s, err := strconv.ParseFloat(m[1], 64); err == nil {
			return s
		}
	}
	return wall.Seconds()
}

func markStatus(target, out string, err error) string {
	if strings.HasPrefix(out, "SKIP:") { return "SKIP" }
	if err != nil { return "ERR" }
	if target == "tengri-ast" {
		if strings.Contains(out, "Ошибка парсера") ||
			strings.Contains(out, "не найдена функция для разбора токена") {
			return "ERR"
		}
	}
	return "OK"
}

/* ---------- prepare (build once) ---------- */

func prepGo() {
	if has("04_benchmarks/fibonacci_iter.go") && !has(".bin/fib_iter_go") {
		_, _, _ = run("go", "build", "-tags=iter", "-o", ".bin/fib_iter_go", "04_benchmarks/fibonacci_iter.go")
	}
}

func prepVM() {
	if has("05_vm_mini/main.go") && !has(".bin/vm") {
		_, _, _ = run("go", "build", "-o", ".bin/vm", "05_vm_mini/main.go")
	}
}

func prepAOT() {
	if has("06_aot_minic/main.go") && !has(".bin/tengri-aot") {
		_, _, _ = sh("cd 06_aot_minic && GO111MODULE=on go build -o ../.bin/tengri-aot .")
	}
	if has(".bin/tengri-aot") {
		if has("06_aot_minic/examples/fib_cli.tgr") && !has(".bin/fib_cli") {
			_, _, _ = run(".bin/tengri-aot", "06_aot_minic/examples/fib_cli.tgr", "-o", ".bin/fib_cli.c")
			_, _, _ = run("clang", "-O2", "-o", ".bin/fib_cli", ".bin/fib_cli.c", "06_aot_minic/runtime/runtime.c")
		}
		if has("06_aot_minic/examples/fib_rec_cli.tgr") && !has(".bin/fib_rec_cli") {
			_, _, _ = run(".bin/tengri-aot", "06_aot_minic/examples/fib_rec_cli.tgr", "-o", ".bin/fib_rec_cli.c")
			_, _, _ = run("clang", "-O2", "-o", ".bin/fib_rec_cli", ".bin/fib_rec_cli.c", "06_aot_minic/runtime/runtime.c")
		}
	}
}

/* --------------- targets ---------------- */

func runGoBin(task string, n int) (string, time.Duration, error) {
	switch task {
	case "fib_rec":
		if !has("04_benchmarks/fibonacci.go") { return "SKIP: 04_benchmarks/fibonacci.go missing", 0, nil }
		return run("go", "run", "04_benchmarks/fibonacci.go")
	case "fib_iter":
		if has(".bin/fib_iter_go") { return run(".bin/fib_iter_go", strconv.Itoa(n)) }
		if !has("04_benchmarks/fibonacci_iter.go") { return "SKIP: 04_benchmarks/fibonacci_iter.go missing", 0, nil }
		return run("go", "run", "-tags=iter", "04_benchmarks/fibonacci_iter.go", strconv.Itoa(n))
	}
	return "SKIP: unknown task", 0, nil
}

func runPy(task string, n int) (string, time.Duration, error) {
	switch task {
	case "fib_rec":
		if !has("04_benchmarks/fibonacci.py") { return "SKIP: 04_benchmarks/fibonacci.py missing", 0, nil }
		return run("python3", "04_benchmarks/fibonacci.py")
	case "fib_iter":
		if !has("04_benchmarks/fibonacci_iter.py") { return "SKIP: 04_benchmarks/fibonacci_iter.py missing", 0, nil }
		return run("python3", "04_benchmarks/fibonacci_iter.py", strconv.Itoa(n))
	}
	return "SKIP: unknown task", 0, nil
}

func runVMBin(task string, n int) (string, time.Duration, error) {
	if task != "fib_iter" { return "SKIP: vm supports fib_iter only", 0, nil }
	if has(".bin/vm") { return run(".bin/vm", strconv.Itoa(n)) }
	if !has("05_vm_mini/main.go") { return "SKIP: 05_vm_mini/main.go missing", 0, nil }
	return run("go", "run", "05_vm_mini/main.go", strconv.Itoa(n))
}

func runTengriAST(task string, n int) (string, time.Duration, error) {
	if !has("03_compiler_go") { return "SKIP: 03_compiler_go missing", 0, nil }
	return sh("cd 03_compiler_go && go run . || true")
}

func runAOTBinIter(n int, reps int) (string, time.Duration, error) {
	if !has(".bin/fib_cli") { return "SKIP: .bin/fib_cli missing", 0, nil }
	return run(".bin/fib_cli", strconv.Itoa(n), strconv.Itoa(reps))
}

func runAOTBinRec(n int, reps int) (string, time.Duration, error) {
	if !has(".bin/fib_rec_cli") { return "SKIP: .bin/fib_rec_cli missing", 0, nil }
	return run(".bin/fib_rec_cli", strconv.Itoa(n), strconv.Itoa(reps))
}

/* plotting stub (no-op) */
func plotWithGo(outDir string) error { return nil }