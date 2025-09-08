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
	Task   string
	Target string
	Cmd    []string
	Env    []string
	Dir    string
	N      int
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

func parseInlineTime(out string) (float64, bool) {
	// Ищем строку вида: TIME: <float>
	lines := strings.Split(out, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "TIME:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "TIME:"))
			if f, err := strconv.ParseFloat(val, 64); err == nil && f >= 0 {
				return f, true
			}
		}
	}
	return 0, false
}

// ── Генерация Tengri-скрипта (рекурсивный fib, круглые блоки)
func writeTengriScriptRec(n int) (string, error) {
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

type dirs struct {
	repoRoot    string
	tengriDir   string
	runsRoot    string
	latestRoot  string
	runStamp    string
	runResults  string
	runPlots    string
	latestRes   string
	latestPlots string
	binDir      string
}

func ensureDirs() (d dirs, err error) {
	d.repoRoot, err = os.Getwd()
	if err != nil {
		return
	}
	d.tengriDir = filepath.Join(d.repoRoot, "03_compiler_go")

	d.runStamp = time.Now().Format("20060102-150405")

	bmarks := filepath.Join(d.repoRoot, "benchmarks")
	d.runsRoot = filepath.Join(bmarks, "runs", d.runStamp)
	d.latestRoot = filepath.Join(bmarks, "latest")

	d.runResults = filepath.Join(d.runsRoot, "results")
	d.runPlots = filepath.Join(d.runsRoot, "plots")
	d.latestRes = filepath.Join(d.latestRoot, "results")
	d.latestPlots = filepath.Join(d.latestRoot, "plots")

	d.binDir = filepath.Join(d.repoRoot, ".bin")

	for _, p := range []string{
		filepath.Join(d.repoRoot, "benchmarks", "runs"),
		d.runsRoot, d.runResults, d.runPlots,
		d.latestRoot, d.latestRes, d.latestPlots,
		d.binDir,
	} {
		if mkErr := os.MkdirAll(p, 0o755); mkErr != nil && err == nil {
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

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

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
	// Отключим «optional build tengri» чтобы не шумел логом
	return "", fmt.Errorf("skip optional build")
}

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
		if i == 0 {
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

func main() {
	plotFlag := flag.Bool("plot", false, "Render plots to benchmarks/{runs/<ts>|latest}/plots/*.png")
	flag.Parse()

	D, err := ensureDirs()
	if err != nil {
		fmt.Println("init error:", err)
		return
	}

	// Наборы N (оставим в безопасной зоне int8 для VM demo и достаточной нагрузке)
	NsRec := []int{30, 32, 34}
	NsIter := []int{40, 60, 90}

	goFibArg := filepath.Join(D.repoRoot, "04_benchmarks", "fibonacci_arg.go")
	goFibIter := filepath.Join(D.repoRoot, "04_benchmarks", "fibonacci_iter.go")
	pyFibArg := filepath.Join(D.repoRoot, "04_benchmarks", "fibonacci_arg.py")
	pyFibIter := filepath.Join(D.repoRoot, "04_benchmarks", "fibonacci_iter.py")
	vmSrc := filepath.Join(D.repoRoot, "05_vm_mini", "main.go")
	tengriSrcDir := filepath.Join(D.repoRoot, "03_compiler_go")

	goRecBin := filepath.Join(D.binDir, "fib_rec_go")
	goIterBin := filepath.Join(D.binDir, "fib_iter_go")
	vmBin := filepath.Join(D.binDir, "vm")
	// tengri как бинарник не собираем (шумело в логе), сразу run .
	_, _ = tengriSrcDir, vmSrc

	if _, err := ensureBuiltGo(goRecBin, []string{"arg"}, goFibArg, D.repoRoot); err != nil {
		fmt.Println("build go(rec):", err)
	}
	if _, err := ensureBuiltGo(goIterBin, []string{"iter"}, goFibIter, D.repoRoot); err != nil {
		fmt.Println("build go(iter):", err)
	}
	if _, err := ensureBuiltGo(vmBin, nil, vmSrc, D.repoRoot); err != nil {
		fmt.Println("build vm:", err)
	}

	// CSV (run + latest)
	fibRecCSV := filepath.Join(D.runResults, "fib_rec.csv")
	fibIterCSV := filepath.Join(D.runResults, "fib_iter.csv")
	_ = writeCSVHeader(fibRecCSV)
	_ = writeCSVHeader(fibIterCSV)

	fibRecLatest := filepath.Join(D.latestRes, "fib_rec.csv")
	fibIterLatest := filepath.Join(D.latestRes, "fib_iter.csv")
	_ = writeCSVHeader(fibRecLatest)
	_ = writeCSVHeader(fibIterLatest)

	// ── fib_rec
	for _, N := range NsRec {
		scriptRec, err := writeTengriScriptRec(N)
		if err != nil {
			fmt.Println("prep-rec error:", err)
			continue
		}
		fmt.Printf("\nTask = %s, N = %d\n", "fib_rec", N)
		fmt.Println("──────────────────────────────────────────────────────────────")
		fmt.Printf("%-20s  %-10s  %s\n", "Target", "Time (s)", "Output")
		fmt.Println(strings.Repeat("-", 66))

		tests := []Case{
			{Task: "fib_rec", Target: "go",         Cmd: []string{goRecBin, fmt.Sprintf("%d", N)},                               Dir: D.repoRoot,  N: N},
			{Task: "fib_rec", Target: "tengri-ast", Cmd: []string{"go", "run", "."}, Env: []string{"TENGRI_SCRIPT=" + scriptRec}, Dir: D.tengriDir, N: N},
			{Task: "fib_rec", Target: "python",     Cmd: []string{"python3", pyFibArg, fmt.Sprintf("%d", N)},                     Dir: D.repoRoot,  N: N},
		}
		for _, t := range tests {
			d, out, err := runCase(t)
			status := markStatus(t.Target, out, err)
			// если есть TIME:, берём внутреннее время, иначе wall-clock
			itime, ok := parseInlineTime(out)
			use := d.Seconds()
			if ok {
				use = itime
			}
			last := ""
			if parts := strings.Split(out, "\n"); len(parts) > 0 {
				last = parts[len(parts)-1]
			}
			fmt.Printf("%-20s  %8.6f  [%s] %s\n", t.Target, use, status, last)
			_ = writeCSVRow(fibRecCSV, []string{t.Task, t.Target, fmt.Sprintf("%d", t.N), fmt.Sprintf("%.6f", use), status, last})
			_ = writeCSVRow(fibRecLatest, []string{t.Task, t.Target, fmt.Sprintf("%d", t.N), fmt.Sprintf("%.6f", use), status, last})
		}
	}

	// ── fib_iter
	for _, N := range NsIter {
		fmt.Printf("\nTask = %s, N = %d\n", "fib_iter", N)
		fmt.Println("──────────────────────────────────────────────────────────────")
		fmt.Printf("%-20s  %-10s  %s\n", "Target", "Time (s)", "Output")
		fmt.Println(strings.Repeat("-", 66))
		tests := []Case{
			{Task: "fib_iter", Target: "go",     Cmd: []string{goIterBin, fmt.Sprintf("%d", N)},                 Dir: D.repoRoot, N: N},
			{Task: "fib_iter", Target: "python", Cmd: []string{"python3", pyFibIter, fmt.Sprintf("%d", N)},      Dir: D.repoRoot, N: N},
			{Task: "fib_iter", Target: "vm",     Cmd: []string{vmBin, fmt.Sprintf("%d", N)},                      Dir: D.repoRoot, N: N},
		}
		for _, t := range tests {
			d, out, err := runCase(t)
			status := markStatus(t.Target, out, err)
			itime, ok := parseInlineTime(out)
			use := d.Seconds()
			if ok {
				use = itime
			}
			last := ""
			if parts := strings.Split(out, "\n"); len(parts) > 0 {
				last = parts[len(parts)-1]
			}
			fmt.Printf("%-20s  %8.6f  [%s] %s\n", t.Target, use, status, last)
			_ = writeCSVRow(fibIterCSV, []string{t.Task, t.Target, fmt.Sprintf("%d", t.N), fmt.Sprintf("%.6f", use), status, last})
			_ = writeCSVRow(fibIterLatest, []string{t.Task, t.Target, fmt.Sprintf("%d", t.N), fmt.Sprintf("%.6f", use), status, last})
		}
	}

	fmt.Printf("\nCSV saved (run %s):\n- %s\n- %s\n", D.runStamp, filepath.Join(D.runResults, "fib_rec.csv"), filepath.Join(D.runResults, "fib_iter.csv"))
	fmt.Printf("CSV saved (latest):\n- %s\n- %s\n", filepath.Join(D.latestRes, "fib_rec.csv"), filepath.Join(D.latestRes, "fib_iter.csv"))

	if *plotFlag {
		recPNG := filepath.Join(D.runPlots, "fib_rec.png")
		iterPNG := filepath.Join(D.runPlots, "fib_iter.png")
		if err := plotFromCSV(filepath.Join(D.runResults, "fib_rec.csv"), recPNG, "Fibonacci (recursive)"); err != nil {
			fmt.Println("plot fib_rec:", err)
		} else {
			fmt.Println("plot saved:", recPNG)
			_ = copyFile(recPNG, filepath.Join(D.latestPlots, "fib_rec.png"))
		}
		if err := plotFromCSV(filepath.Join(D.runResults, "fib_iter.csv"), iterPNG, "Fibonacci (iterative, mod M)"); err != nil {
			fmt.Println("plot fib_iter:", err)
		} else {
			fmt.Println("plot saved:", iterPNG)
			_ = copyFile(iterPNG, filepath.Join(D.latestPlots, "fib_iter.png"))
		}
	}
}