package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Case struct {
	Task   string   // "fib_rec" | "fib_iter"
	Target string   // "go" | "python" | "tengri-ast" | "vm"
	Cmd    []string // команда запуска (путь к бинарнику или интерпретатор + скрипт)
	Env    []string // доп. переменные окружения
	Dir    string   // рабочая директория
	N      int      // размер входа
}

func runCase(c Case) (time.Duration, string, error) {
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

// ───────────────────────────────────────────────────────────────
// Генерация скриптов для Tengri (круглые блоки, без while)
// ───────────────────────────────────────────────────────────────

func writeTengriScriptRec(n int) (string, error) {
	// Рекурсивный fib, совместимый с вашим парсером (круглые блоки)
	s := fmt.Sprintf(`let fib = fn(n) (
  if (n < 2) ( return n; )
  return fib(n - 1) + fib(n - 2);
);

let n = %d;
fib(n);
`, n)
	tmp := filepath.Join("03_compiler_go", ".tg_fib_rec.tg")
	if err := os.WriteFile(tmp, []byte(s), 0o644); err != nil {
		return "", err
	}
	return tmp, nil
}

// ───────────────────────────────────────────────────────────────
// Инфраструктура: директории, CSV, печать таблиц, build
// ───────────────────────────────────────────────────────────────

func ensureDirs() (repoRoot, tengriDir, resultsDir, plotsDir, binDir string, err error) {
	repoRoot, err = os.Getwd()
	if err != nil {
		return
	}
	tengriDir = filepath.Join(repoRoot, "03_compiler_go")
	resultsDir = filepath.Join(repoRoot, "benchmarks", "results")
	plotsDir = filepath.Join(repoRoot, "benchmarks", "plots")
	binDir = filepath.Join(repoRoot, ".bin")
	for _, d := range []string{resultsDir, plotsDir, binDir} {
		if mkErr := os.MkdirAll(d, 0o755); mkErr != nil && err == nil {
			err = mkErr
		}
	}
	return
}

func tableHeader(task string, n int) {
	fmt.Printf("\nTask = %s, N = %d\n", task, n)
	fmt.Println("──────────────────────────────────────────────────────────────")
	fmt.Printf("%-20s  %-10s  %s\n", "Target", "Time (s)", "Output")
	fmt.Println(strings.Repeat("-", 66))
}

func writeCSVHeader(path string) error {
	return os.WriteFile(path, []byte("task,target,n,time,status,output\n"), 0o644)
}

func writeCSVRow(path string, row []string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	if err := w.Write(row); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

// ───────────────────────────────────────────────────────────────
// Детектор статуса: честно помечаем tengri-ast как ERR по тексту
// ───────────────────────────────────────────────────────────────

func markStatus(target, out string, err error) string {
	if err != nil {
		return "ERR"
	}
	if target == "tengri-ast" {
		if strings.Contains(out, "Ошибка парсера") ||
			strings.Contains(out, "не найдена функция для разбора токена") {
			return "ERR"
		}
	}
	return "OK"
}

// ───────────────────────────────────────────────────────────────
// Build cache: собираем один раз бинарники, потом запускаем
// ───────────────────────────────────────────────────────────────

func runBuild(cmd []string, dir string) error {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Dir = dir
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %v\n%s", err, string(out))
	}
	return nil
}

func ensureBuiltGo(binPath string, buildTags []string, src string, repoRoot string) (string, error) {
	if _, err := os.Stat(binPath); err == nil {
		return binPath, nil
	}
	args := []string{"build"}
	if len(buildTags) > 0 {
		args = append(args, "-tags="+strings.Join(buildTags, ","))
	}
	args = append(args, "-o", binPath, src)
	if err := runBuild(append([]string{"go"}, args...), repoRoot); err != nil {
		return "", err
	}
	return binPath, nil
}

func ensureBuiltDir(binPath string, srcDir string, repoRoot string) (string, error) {
	if _, err := os.Stat(binPath); err == nil {
		return binPath, nil
	}
	// go build ./<dir>
	args := []string{"build", "-o", binPath, srcDir}
	if err := runBuild(append([]string{"go"}, args...), repoRoot); err != nil {
		return "", err
	}
	return binPath, nil
}

// ───────────────────────────────────────────────────────────────
// Plotting (gonum/plot)
// ───────────────────────────────────────────────────────────────

type rec struct {
	task   string
	target string
	n      int
	sec    float64
	status string
}

func loadCSV(path string) ([]rec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	var out []rec
	for i, row := range rows {
		if i == 0 { // header
			continue
		}
		if len(row) < 6 {
			continue
		}
		nv, _ := strconv.Atoi(row[2])
		tv, _ := strconv.ParseFloat(row[3], 64)
		out = append(out, rec{
			task:   row[0],
			target: row[1],
			n:      nv,
			sec:    tv,
			status: row[4],
		})
	}
	return out, nil
}

func uniqueTargets(recs []rec) []string {
	set := map[string]struct{}{}
	for _, r := range recs {
		set[r.target] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	return out
}

func nsSorted(recs []rec) []int {
	set := map[int]struct{}{}
	for _, r := range recs {
		set[r.n] = struct{}{}
	}
	out := make([]int, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	// простой insertion sort для маленького множества
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1] > out[j]; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}

func plotFromCSV(csvPath, pngPath, title string) error {
	data, err := loadCSV(csvPath)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("no data in %s", csvPath)
	}
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = "N"
	p.Y.Label.Text = "Time (s)"

	targets := uniqueTargets(data)
	ns := nsSorted(data)

	lines := []interface{}{}
	for _, t := range targets {
		pts := make(plotter.XYs, 0, len(ns))
		for _, n := range ns {
			found := false
			val := 0.0
			for _, r := range data {
				if r.target == t && r.n == n {
					val = r.sec
					found = true
					break
				}
			}
			if found {
				pts = append(pts, plotter.XY{X: float64(n), Y: val})
			}
		}
		if len(pts) == 0 {
			continue
		}
		lines = append(lines, t)
		lines = append(lines, pts)
	}

	if err := plotutil.AddLines(p, lines...); err != nil {
		return err
	}
	if err := p.Save(800*vg.Points(1), 500*vg.Points(1), pngPath); err != nil {
		return err
	}
	return nil
}

// ───────────────────────────────────────────────────────────────

func main() {
	plotFlag := flag.Bool("plot", false, "Render plots to benchmarks/plots/*.png using gonum/plot")
	flag.Parse()

	repoRoot, tengriDir, resultsDir, plotsDir, binDir, err := ensureDirs()
	if err != nil {
		fmt.Println("init error:", err)
		return
	}

	// Быстрые наборы N
	NsRec := []int{30, 32, 34}  // рекурсивные (малые)
	NsIter := []int{30, 35, 40} // итеративные (чтобы не ждать)

	// Пути к исходникам
	goFibArg := filepath.Join(repoRoot, "04_benchmarks", "fibonacci_arg.go")   // рекурсивный (Go, build tag: arg)
	goFibIter := filepath.Join(repoRoot, "04_benchmarks", "fibonacci_iter.go") // итеративный (Go, build tag: iter)
	pyFibArg := filepath.Join(repoRoot, "04_benchmarks", "fibonacci_arg.py")   // рекурсивный (Py)
	pyFibIter := filepath.Join(repoRoot, "04_benchmarks", "fibonacci_iter.py") // итеративный (Py)
	vmSrc := filepath.Join(repoRoot, "05_vm_mini", "main.go")                  // VM (Go)
	tengriSrcDir := filepath.Join(repoRoot, "03_compiler_go")                  // интерпретатор Tengri

	// Пути к бинарникам (кэш сборки)
	goRecBin := filepath.Join(binDir, "fib_rec_go")
	goIterBin := filepath.Join(binDir, "fib_iter_go")
	vmBin := filepath.Join(binDir, "vm")
	tengriBin := filepath.Join(binDir, "tengri")

	// Сборка бинарников (если их нет)
	if _, err := ensureBuiltGo(goRecBin, []string{"arg"}, goFibArg, repoRoot); err != nil {
		fmt.Println("build go(rec):", err)
	}
	if _, err := ensureBuiltGo(goIterBin, []string{"iter"}, goFibIter, repoRoot); err != nil {
		fmt.Println("build go(iter):", err)
	}
	if _, err := ensureBuiltGo(vmBin, nil, vmSrc, repoRoot); err != nil {
		fmt.Println("build vm:", err)
	}
	// Попробуем собрать интерпретатор Tengri (если проект позволяет). Если не соберётся — benchfast будет делать go run .
	if _, err := ensureBuiltDir(tengriBin, tengriSrcDir, repoRoot); err != nil {
		fmt.Println("build tengri (optional):", err)
	}

	// CSV файлы
	fibRecCSV := filepath.Join(resultsDir, "fib_rec.csv")
	fibIterCSV := filepath.Join(resultsDir, "fib_iter.csv")
	_ = writeCSVHeader(fibRecCSV)
	_ = writeCSVHeader(fibIterCSV)

	// ─────────────── fib_rec ───────────────
	for _, N := range NsRec {
		// подготовим скрипт для Tengri
		scriptRec, err := writeTengriScriptRec(N)
		if err != nil {
			fmt.Println("prep-rec error:", err)
			continue
		}
		tableHeader("fib_rec", N)

		// Prefer running tengri as a built binary; if it doesn't exist, fallback to `go run .`
		tengriCmd := []string{}
		if info, err := os.Stat(tengriBin); err == nil && info.Mode().Perm()&0o111 != 0 {
			tengriCmd = []string{tengriBin}
		} else {
			tengriCmd = []string{"go", "run", "."}
		}

		tests := []Case{
			{Task: "fib_rec", Target: "go",         Cmd: []string{goRecBin, fmt.Sprintf("%d", N)},                               Dir: repoRoot,  N: N},
			{Task: "fib_rec", Target: "tengri-ast", Cmd: tengriCmd, Env: []string{"TENGRI_SCRIPT=" + scriptRec},                Dir: tengriDir, N: N},
			{Task: "fib_rec", Target: "python",     Cmd: []string{"python3", pyFibArg, fmt.Sprintf("%d", N)},                    Dir: repoRoot,  N: N},
		}

		for _, t := range tests {
			d, out, err := runCase(t)
			status := markStatus(t.Target, out, err)
			last := ""
			if parts := strings.Split(out, "\n"); len(parts) > 0 {
				last = parts[len(parts)-1]
			}
			fmt.Printf("%-20s  %8.3f  [%s] %s\n", t.Target, d.Seconds(), status, last)
			_ = writeCSVRow(fibRecCSV, []string{t.Task, t.Target, fmt.Sprintf("%d", t.N), fmt.Sprintf("%.6f", d.Seconds()), status, last})
		}
	}

	// ─────────────── fib_iter ───────────────
	for _, N := range NsIter {
		tableHeader("fib_iter", N)
		tests := []Case{
			{Task: "fib_iter", Target: "go",     Cmd: []string{goIterBin, fmt.Sprintf("%d", N)}, Dir: repoRoot, N: N},
			{Task: "fib_iter", Target: "python", Cmd: []string{"python3", pyFibIter, fmt.Sprintf("%d", N)}, Dir: repoRoot, N: N},
			{Task: "fib_iter", Target: "vm",     Cmd: []string{vmBin, fmt.Sprintf("%d", N)},     Dir: repoRoot, N: N},
		}
		for _, t := range tests {
			d, out, err := runCase(t)
			status := markStatus(t.Target, out, err)
			last := ""
			if parts := strings.Split(out, "\n"); len(parts) > 0 {
				last = parts[len(parts)-1]
			}
			fmt.Printf("%-20s  %8.3f  [%s] %s\n", t.Target, d.Seconds(), status, last)
			_ = writeCSVRow(fibIterCSV, []string{t.Task, t.Target, fmt.Sprintf("%d", t.N), fmt.Sprintf("%.6f", d.Seconds()), status, last})
		}
	}

	fmt.Printf("\nCSV saved:\n- %s\n- %s\n", fibRecCSV, fibIterCSV)

	// ─────────────── Графики ───────────────
	if *plotFlag {
		recPNG := filepath.Join(plotsDir, "fib_rec.png")
		iterPNG := filepath.Join(plotsDir, "fib_iter.png")
		if err := plotFromCSV(fibRecCSV, recPNG, "Fibonacci (recursive)"); err != nil {
			fmt.Println("plot fib_rec:", err)
		} else {
			fmt.Println("plot saved:", recPNG)
		}
		if err := plotFromCSV(fibIterCSV, iterPNG, "Fibonacci (iterative)"); err != nil {
			fmt.Println("plot fib_iter:", err)
		} else {
			fmt.Println("plot saved:", iterPNG)
		}
	}
}