package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	//"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// -----------------------------
// Config & flags
// -----------------------------

var (
	flagRebuild   = flag.Bool("rebuild", false, "Rebuild all .bin targets before running benchmarks")
	flagCSVDir    = flag.String("csvdir", "benchmarks/latest/results", "Directory to write CSV results")
	flagRunsDir   = flag.String("runsdir", "benchmarks/runs", "Directory to store timestamped results")
	flagNoColor   = flag.Bool("nocolor", false, "Disable ANSI colors in output")
)

var (
	// Recursive task Ns (kept short because recursive fib explodes)
	NsRec  = []int{30, 32, 34}
	// Iterative task Ns (quick & scalable)
	NsIter = []int{40, 60, 90}
)

const (
	// Binaries we expect (built by -rebuild)
	binFibRecGo  = ".bin/fib_rec_go"
	binFibIterGo = ".bin/fib_iter_go"
	binVM        = ".bin/vm"
	binAOT       = ".bin/tengri-aot"

	// AOT-produced C/ELF binaries (also made by -rebuild)
	binFibCLI     = ".bin/fib_cli"      // iterative
	binFibRecCLI  = ".bin/fib_rec_cli"  // recursive

	// Sources
	srcGoFibRec      = "benchmarks/src/fib_rec/fibonacci.go"
	srcGoFibIter     = "benchmarks/src/fib_iter/go/fibonacci_iter.go"
	srcPyFibRec      = "benchmarks/src/fib_rec/fibonacci.py"
	srcPyFibIter     = "benchmarks/src/fib_iter/python/fibonacci_iter.py"
	srcTgrFibIter    = "benchmarks/src/fib_iter/tengri/fib_cli.tgr"
	srcTgrFibRec     = "benchmarks/src/fib_rec/tengri/fib_rec_cli.tgr"
	aotRuntimeC      = "internal/aotminic/runtime/runtime.c"
	cmdVMMain        = "cmd/tengri-vm/main.go"
	cmdAOTMain       = "cmd/tengri-aot/main.go"
)

type row struct {
	target   string
	seconds  float64
	status   string // [OK]/[ERR]/[SKIP]
	output   string // short preview or reason
}

type table struct {
	taskName string
	rows     []row
}

func main() {
	flag.Parse()

	ts := time.Now().Format("20060102-150405")
	runDir := filepath.Join(*flagRunsDir, ts, "results")
	must(os.MkdirAll(runDir, 0o755))
	must(os.MkdirAll(*flagCSVDir, 0o755))

	if *flagRebuild {
		fmt.Println(gray("Rebuilding .bin targets…"))
		if err := rebuildAll(); err != nil {
			failf("rebuild failed: %v", err)
		}
	}

	// --- fib_rec ---
	recTable := runFibRec()
	saveCSV(filepath.Join(runDir, "fib_rec.csv"), recTable)
	saveCSV(filepath.Join(*flagCSVDir, "fib_rec.csv"), recTable)

	// --- fib_iter ---
	iterTable := runFibIter()
	saveCSV(filepath.Join(runDir, "fib_iter.csv"), iterTable)
	saveCSV(filepath.Join(*flagCSVDir, "fib_iter.csv"), iterTable)
}

// -----------------------------
// Rebuild pipeline
// -----------------------------

func rebuildAll() error {
	// 1) Build Go baselines
	if err := sh("go", "build", "-o", binFibRecGo, srcGoFibRec); err != nil {
		return fmt.Errorf("build fib_rec_go: %w", err)
	}
	if err := sh("go", "build", "-tags=iter", "-o", binFibIterGo, srcGoFibIter); err != nil {
		return fmt.Errorf("build fib_iter_go: %w", err)
	}

	// 2) Build VM mini (optional for iter timing)
	if fileExists(cmdVMMain) {
		if err := sh("go", "build", "-o", binVM, cmdVMMain); err != nil {
			return fmt.Errorf("build vm: %w", err)
		}
	}

	// 3) Build AOT transpiler
	if fileExists(cmdAOTMain) {
		if err := sh("go", "build", "-o", binAOT, cmdAOTMain); err != nil {
			return fmt.Errorf("build tengri-aot: %w", err)
		}
		// 4) Transpile TGR → C and compile C with runtime
		if err := aotProduce(binFibCLI, srcTgrFibIter); err != nil {
			return err
		}
		if err := aotProduce(binFibRecCLI, srcTgrFibRec); err != nil {
			return err
		}
	}
	return nil
}

func aotProduce(outBin, srcTgr string) error {
	cFile := strings.TrimSuffix(outBin, filepath.Ext(outBin)) + ".c"
	if err := sh(binAOT, srcTgr, "-o", cFile); err != nil {
		return fmt.Errorf("aot transpile %s: %w", srcTgr, err)
	}
	if err := sh("clang", "-O2", "-o", outBin, cFile, aotRuntimeC); err != nil {
		return fmt.Errorf("clang link %s: %w", outBin, err)
	}
	return nil
}

// -----------------------------
// Fib (recursive)
// -----------------------------

func runFibRec() table {
	tbl := table{taskName: "fib_rec"}
	fmt.Println()
	fmt.Println(bold("Task = fib_rec"))
	fmt.Println(line())
	fmt.Println(gray("TIMING: prefer TIME_NS over wall-clock; fallback to TIME:, then wall-clock"))

	// Header
	fmt.Println(bandHeader())

	for _, N := range NsRec {
		fmt.Println()
		fmt.Printf("N = %d\n", N)
		fmt.Println(line())
		// go (prebuilt)
		{
			sec, out, status := runAndParse(binFibRecGo, nil)
			tbl.rows = append(tbl.rows, row{"go", sec, status, out})
			fmt.Println(formatRow("go", sec, status, out))
		}
		// tengri-ast (disabled in new layout)
		{
			tbl.rows = append(tbl.rows, row{"tengri-ast", 0, "SKIP", "SKIP: AST is disabled in the new layout"})
			fmt.Println(formatRow("tengri-ast", 0, "SKIP", "SKIP: AST is disabled in the new layout"))
		}
		// python
		{
			if fileExists(srcPyFibRec) {
				sec, out, status := runAndParse("python3", []string{srcPyFibRec})
				tbl.rows = append(tbl.rows, row{"python", sec, status, out})
				fmt.Println(formatRow("python", sec, status, out))
			} else {
				tbl.rows = append(tbl.rows, row{"python", 0, "SKIP", "SKIP: " + srcPyFibRec + " missing"})
				fmt.Println(formatRow("python", 0, "SKIP", "SKIP: "+srcPyFibRec+" missing"))
			}
		}
		// vm (not applicable)
		{
			tbl.rows = append(tbl.rows, row{"vm", 0, "SKIP", "SKIP: vm supports fib_iter only"})
			fmt.Println(formatRow("vm", 0, "SKIP", "SKIP: vm supports fib_iter only"))
		}
		// tengri-aot (prebuilt fib_rec_cli expects N)
		{
			if fileExists(binFibRecCLI) {
				sec, out, status := runAndParse(binFibRecCLI, []string{fmt.Sprint(N)})
				tbl.rows = append(tbl.rows, row{"tengri-aot", sec, status, out})
				fmt.Println(formatRow("tengri-aot", sec, status, out))
			} else {
				tbl.rows = append(tbl.rows, row{"tengri-aot", 0, "ERR", "AOT binary missing: " + binFibRecCLI})
				fmt.Println(formatRow("tengri-aot", 0, "ERR", "AOT binary missing: "+binFibRecCLI))
			}
		}
	}
	fmt.Println(saveNote("fib_rec"))
	return tbl
}

// -----------------------------
// Fib (iterative)
// -----------------------------

func runFibIter() table {
	tbl := table{taskName: "fib_iter"}
	fmt.Println()
	fmt.Println(bold("Task = fib_iter"))
	fmt.Println(line())
	fmt.Println(gray("TIMING: prefer TIME_NS over wall-clock; fallback to TIME:, then wall-clock"))
	fmt.Println(bandHeader())

	for _, N := range NsIter {
		fmt.Println()
		fmt.Printf("N = %d\n", N)
		fmt.Println(line())

		// go (prebuilt, accepts N)
		{
			sec, out, status := runAndParse(binFibIterGo, []string{fmt.Sprint(N)})
			tbl.rows = append(tbl.rows, row{"go", sec, status, out})
			fmt.Println(formatRow("go", sec, status, out))
		}
		// tengri-ast (disabled)
		{
			tbl.rows = append(tbl.rows, row{"tengri-ast", 0, "SKIP", "SKIP: AST is disabled in the new layout"})
			fmt.Println(formatRow("tengri-ast", 0, "SKIP", "SKIP: AST is disabled in the new layout"))
		}
		// python (accepts N)
		{
			if fileExists(srcPyFibIter) {
				sec, out, status := runAndParse("python3", []string{srcPyFibIter, fmt.Sprint(N)})
				tbl.rows = append(tbl.rows, row{"python", sec, status, out})
				fmt.Println(formatRow("python", sec, status, out))
			} else {
				tbl.rows = append(tbl.rows, row{"python", 0, "SKIP", "SKIP: " + srcPyFibIter + " missing"})
				fmt.Println(formatRow("python", 0, "SKIP", "SKIP: "+srcPyFibIter+" missing"))
			}
		}
		// vm (accepts N)
		{
			if fileExists(binVM) {
				sec, out, status := runAndParse(binVM, []string{fmt.Sprint(N)})
				tbl.rows = append(tbl.rows, row{"vm", sec, status, out})
				fmt.Println(formatRow("vm", sec, status, out))
			} else {
				tbl.rows = append(tbl.rows, row{"vm", 0, "ERR", "vm binary missing: " + binVM})
				fmt.Println(formatRow("vm", 0, "ERR", "vm binary missing: "+binVM))
			}
		}
		// tengri-aot (accepts N)
		{
			if fileExists(binFibCLI) {
				sec, out, status := runAndParse(binFibCLI, []string{fmt.Sprint(N)})
				tbl.rows = append(tbl.rows, row{"tengri-aot", sec, status, out})
				fmt.Println(formatRow("tengri-aot", sec, status, out))
			} else {
				tbl.rows = append(tbl.rows, row{"tengri-aot", 0, "ERR", "AOT binary missing: " + binFibCLI})
				fmt.Println(formatRow("tengri-aot", 0, "ERR", "AOT binary missing: "+binFibCLI))
			}
		}
	}
	fmt.Println(saveNote("fib_iter"))
	return tbl
}

// -----------------------------
// Runner & parsing
// -----------------------------

var (
	reTimeNS = regexp.MustCompile(`(?m)TIME_NS:\s*([0-9]+)`)
	reTime   = regexp.MustCompile(`(?m)TIME:\s*([0-9]*\.?[0-9]+)`)
)

func runAndParse(cmd string, args []string) (seconds float64, preview string, status string) {
	if !fileExists(cmd) && !isInPath(cmd) {
		return 0, "missing: " + cmd, "SKIP"
	}
	start := time.Now()
	out, err := run(cmd, args...)
	wall := time.Since(start).Seconds()

	// Prefer TIME_NS
	if m := reTimeNS.FindStringSubmatch(out); len(m) == 2 {
		ns, _ := parseInt64(m[1])
		return float64(ns) / 1e9, clip(out), markStatus(out, err)
	}
	// Fallback to TIME:
	if m := reTime.FindStringSubmatch(out); len(m) == 2 {
		f, _ := parseFloat(m[1])
		return f, clip(out), markStatus(out, err)
	}
	// Else wall-clock
	return wall, clip(out), markStatus(out, err)
}

func markStatus(out string, err error) string {
	if err != nil {
		return "ERR"
	}
	// Parser diagnostic auto-detect (ru messages)
	if strings.Contains(out, "Ошибка парсера") ||
		strings.Contains(out, "не найдена функция для разбора токена") ||
		strings.Contains(out, "ошибка: ожидался") {
		return "ERR"
	}
	return "OK"
}

// -----------------------------
// CSV
// -----------------------------

func saveCSV(path string, tbl table) {
	must(os.MkdirAll(filepath.Dir(path), 0o755))
	f, err := os.Create(path)
	must(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	must(w.Write([]string{"task", "target", "seconds", "status"}))
	for _, r := range tbl.rows {
		must(w.Write([]string{
			tbl.taskName, r.target, fmt.Sprintf("%.9f", r.seconds), r.status,
		}))
	}
	fmt.Println(gray(fmt.Sprintf("CSV saved: %s", path)))
}

// -----------------------------
// Utils: exec, fmt, etc.
// -----------------------------

func run(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	var buf bytes.Buffer
	c.Stdout = &buf
	c.Stderr = &buf
	err := c.Run()
	return buf.String(), err
}

func sh(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func must(err error) {
	if err != nil {
		failf("%v", err)
	}
}

func failf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "benchfast: "+format+"\n", a...)
	os.Exit(1)
}

func parseInt64(s string) (int64, error) {
	var x int64
	_, err := fmt.Sscan(s, &x)
	return x, err
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscan(s, &f)
	return f, err
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func isInPath(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}

func clip(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 70 {
		return s[:70] + "..."
	}
	return s
}

func bandHeader() string {
	return fmt.Sprintf("%-20s %-10s  %s", "Target", "Time (s)", "Output")
}

func formatRow(target string, sec float64, status, out string) string {
	statusTag := "[" + status + "]"
	return fmt.Sprintf("%-20s %-10.6f  %s %s", target, sec, statusTag, out)
}

func line() string {
	return strings.Repeat("─", 58)
}

func bold(s string) string {
	if *flagNoColor {
		return s
	}
	return "\033[1m" + s + "\033[0m"
}

func gray(s string) string {
	if *flagNoColor {
		return s
	}
	return "\033[90m" + s + "\033[0m"
}

func saveNote(task string) string {
	return gray(fmt.Sprintf("CSV saved: benchmarks/runs/<TS>/results/%s.csv and benchmarks/latest/results/%s.csv", task, task))
}