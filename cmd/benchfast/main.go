package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// -----------------------------
// Config
// -----------------------------

// Ns for benchmark sets
var fibIterNs = []int{40, 60, 90}
var fibRecNs  = []int{30, 32, 34}

// Where to save CSVs
const runsRoot     = "benchmarks/runs"
const latestRoot   = "benchmarks/latest"
const resultsLeaf  = "results"

// Table header note
const timingNote = "TIMING: prefer TIME_NS over wall-clock; fallback to TIME:, then wall-clock"

// -----------------------------
// Target model
// -----------------------------

type Target struct {
	Name      string
	Kind      string // "iter" or "rec"
	Bin       string
	ArgsFn    func(n int) []string
	OnlyIf    func() bool // optional existence/probe check; if nil -> assumed available
	SkipMsg   string      // message to print when skipped
}

// -----------------------------
// Utilities
// -----------------------------

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func ensureDir(p string) error {
	return os.MkdirAll(p, 0o755)
}

func nowStamp() string {
	return time.Now().Format("20060102-150405")
}

type RunOutcome struct {
	WallSec float64
	Ok      bool
	Reason  string  // [OK] / [ERR] / [SKIP] reason
	Output  string  // raw stdout
	Result  string  // parsed "RESULT" or numeric last line (for AOT)
	TimeNS  *int64  // parsed TIME_NS if present
	TimeS   *float64 // parsed TIME: seconds if present
}

var reTimeNS = regexp.MustCompile(`(?m)^\s*TIME_NS:\s*([0-9]+)\s*$`)
var reTimeS  = regexp.MustCompile(`(?m)^\s*TIME:\s*([0-9]*\.?[0-9]+)\s*$`)
var reResult = regexp.MustCompile(`(?m)^\s*(?:RESULT:\s*)?(-?[0-9]+)\s*$`)

// parsePrefTiming extracts TIME_NS (preferred) or TIME:, plus a plausible RESULT-like line for display.
func parsePrefTiming(out string) (timeNS *int64, timeS *float64, result string) {
	if m := reTimeNS.FindStringSubmatch(out); len(m) == 2 {
		v := mustParseInt64(m[1])
		timeNS = &v
	}
	if timeNS == nil { // fallback to TIME:
		if m := reTimeS.FindStringSubmatch(out); len(m) == 2 {
			fv := mustParseFloat(m[1])
			timeS = &fv
		}
	}
	// Try to grab a RESULT-looking number (prefer a line starting with RESULT:, else last standalone number)
	if m := reResult.FindAllStringSubmatch(out, -1); len(m) > 0 {
		result = m[len(m)-1][1]
	}
	return
}

func mustParseInt64(s string) int64 {
	var v int64
	fmt.Sscan(s, &v)
	return v
}

func mustParseFloat(s string) float64 {
	var f float64
	fmt.Sscan(s, &f)
	return f
}

// runOne runs the binary with args, captures stdout/stderr, and returns timing info.
// It prefers embedded TIME_NS / TIME: from the program output; otherwise uses wall clock.
func runOne(bin string, args []string, passEnv map[string]string) RunOutcome {
	var outB, errB bytes.Buffer
	cmd := exec.Command(bin, args...)
	cmd.Stdout = &outB
	cmd.Stderr = &errB

	// Inherit env + overlay passEnv
	env := os.Environ()
	for k, v := range passEnv {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	t0 := time.Now()
	err := cmd.Run()
	elapsed := time.Since(t0)
	out := outB.String() + errB.String()

	if err != nil {
		return RunOutcome{
			WallSec: elapsed.Seconds(),
			Ok:      false,
			Reason:  fmt.Sprintf("ERR: %v", err),
			Output:  out,
		}
	}

	tNS, tS, res := parsePrefTiming(out)

	return RunOutcome{
		WallSec: elapsed.Seconds(),
		Ok:      true,
		Reason:  "OK",
		Output:  out,
		Result:  res,
		TimeNS:  tNS,
		TimeS:   tS,
	}
}

// --- ИЗМЕНЕНИЕ ЗДЕСЬ ---
func fmtTimeOutcome(o RunOutcome) (timeCol string, outCol string) {
	// Compose the timing column and the "Output" column text (result/diagnostic)
	switch {
	case o.TimeNS != nil:
		// Отображаем наносекунды
		timeCol = fmt.Sprintf("%d ns", *o.TimeNS)
	case o.TimeS != nil:
		// Если ns нет, показываем секунды (как раньше)
		timeCol = fmt.Sprintf("%.6f s", *o.TimeS)
	default:
		// В крайнем случае — wall-clock
		timeCol = fmt.Sprintf("%.6f s (wall)", o.WallSec)
	}
	tag := "[OK]"
	if !o.Ok {
		tag = "[ERR]"
	}
	if !o.Ok && strings.HasPrefix(o.Reason, "SKIP") {
		tag = "[SKIP]"
	}
	if o.Result != "" && tag == "[OK]" {
		outCol = fmt.Sprintf("%s %s", tag, o.Result)
	} else {
		// show short reason (OK/ERR/SKIP + maybe short message)
		outCol = tag
		if o.Reason != "" && tag != "[OK]" {
			outCol += " " + o.Reason
		}
	}
	return
}

// CSV writer helper
type csvRow struct {
	Target string
	N      int
	TimeNS string // keep as string; if empty — NA
	TimeS  string // same
	WallS  string
	Result string
	Status string // OK/ERR/SKIP
}

func writeCSV(rows []csvRow, taskName string, ts string) error {
	// runs/<TS>/results and latest/results
	destRun := filepath.Join(runsRoot, ts, resultsLeaf)
	destLatest := filepath.Join(latestRoot, resultsLeaf)
	for _, d := range []string{destRun, destLatest} {
		if err := ensureDir(d); err != nil {
			return err
		}
	}

	filename := func(root string) string {
		return filepath.Join(root, fmt.Sprintf("%s.csv", taskName))
	}

	writeTo := func(root string) error {
		f, err := os.Create(filename(root))
		if err != nil {
			return err
		}
		defer f.Close()

		w := csv.NewWriter(f)
		defer w.Flush()

		_ = w.Write([]string{"target", "N", "time_ns", "time_s", "wall_s", "result", "status"})
		for _, r := range rows {
			if err := w.Write([]string{
				r.Target,
				fmt.Sprintf("%d", r.N),
				r.TimeNS,
				r.TimeS,
				r.WallS,
				r.Result,
				r.Status,
			}); err != nil {
				return err
			}
		}
		return w.Error()
	}

	if err := writeTo(destRun); err != nil {
		return err
	}
	if err := writeTo(destLatest); err != nil {
		return err
	}
	fmt.Printf("CSV saved: %s and %s\n", filepath.Join(runsRoot, "<TS>", resultsLeaf, taskName+".csv"), filepath.Join(latestRoot, resultsLeaf, taskName+".csv"))
	fmt.Printf("CSV saved: %s\n", filepath.Join(runsRoot, ts, resultsLeaf, taskName+".csv"))
	fmt.Printf("CSV saved: %s\n", filepath.Join(latestRoot, resultsLeaf, taskName+".csv"))
	return nil
}

// -----------------------------
// Targets registry
// -----------------------------

func availableOr(path, skip string) (func() bool, string) {
	return func() bool { return fileExists(path) }, skip
}

func targetsFor(kind string) []Target {
	var t []Target

	// Go
	if kind == "iter" {
		t = append(t, Target{
			Name: "go",
			Kind: "iter",
			Bin:  ".bin/fib_iter_go",
			ArgsFn: func(n int) []string { return []string{fmt.Sprintf("%d", n)} },
			OnlyIf: func() bool { return fileExists(".bin/fib_iter_go") },
			SkipMsg: "go fib_iter binary missing",
		})
	} else {
		t = append(t, Target{
			Name: "go",
			Kind: "rec",
			Bin:  ".bin/fib_rec_go",
			ArgsFn: func(n int) []string { return []string{fmt.Sprintf("%d", n)} },
			OnlyIf: func() bool { return fileExists(".bin/fib_rec_go") },
			SkipMsg: "go fib_rec binary missing",
		})
	}

	// VM (iter only)
	if kind == "iter" {
		t = append(t, Target{
			Name: "vm",
			Kind: "iter",
			Bin:  ".bin/vm",
			ArgsFn: func(n int) []string {
				// existing vm expects fib_iter only; pass N as arg
				return []string{"fib_iter", fmt.Sprintf("%d", n)}
			},
			OnlyIf: func() bool { return fileExists(".bin/vm") },
			SkipMsg: "vm supports fib_iter only",
		})
	}

	// AOT
	if kind == "iter" {
		t = append(t, Target{
			Name: "tenge-aot",
			Kind: "iter",
			Bin:  ".bin/fib_cli",
			ArgsFn: func(n int) []string { return []string{fmt.Sprintf("%d", n)} },
			OnlyIf: func() bool { return fileExists(".bin/fib_cli") },
			SkipMsg: "AOT binary missing: .bin/fib_cli",
		})
	} else {
		t = append(t, Target{
			Name: "tenge-aot",
			Kind: "rec",
			Bin:  ".bin/fib_rec_cli",
			ArgsFn: func(n int) []string { return []string{fmt.Sprintf("%d", n)} },
			OnlyIf: func() bool { return fileExists(".bin/fib_rec_cli") },
			SkipMsg: "AOT binary missing: .bin/fib_rec_cli",
		})
	}

	// Optional C/C++
	if kind == "iter" {
		t = append(t, Target{
			Name: "c",
			Kind: "iter",
			Bin:  ".bin/fib_iter_c",
			ArgsFn: func(n int) []string { return []string{fmt.Sprintf("%d", n)} },
			OnlyIf: func() bool { return fileExists(".bin/fib_iter_c") },
			SkipMsg: "C binary missing: .bin/fib_iter_c",
		})
	} else {
		t = append(t, Target{
			Name: "c",
			Kind: "rec",
			Bin:  ".bin/fib_rec_c",
			ArgsFn: func(n int) []string { return []string{fmt.Sprintf("%d", n)} },
			OnlyIf: func() bool { return fileExists(".bin/fib_rec_c") },
			SkipMsg: "C binary missing: .bin/fib_rec_c",
		})
	}

	// Optional Rust
	if kind == "iter" {
		t = append(t, Target{
			Name: "rust",
			Kind: "iter",
			Bin:  ".bin/fib_iter_rs",
			ArgsFn: func(n int) []string { return []string{fmt.Sprintf("%d", n)} },
			OnlyIf: func() bool { return fileExists(".bin/fib_iter_rs") },
			SkipMsg: "Rust binary missing: .bin/fib_iter_rs",
		})
	} else {
		t = append(t, Target{
			Name: "rust",
			Kind: "rec",
			Bin:  ".bin/fib_rec_rs",
			ArgsFn: func(n int) []string { return []string{fmt.Sprintf("%d", n)} },
			OnlyIf: func() bool { return fileExists(".bin/fib_rec_rs") },
			SkipMsg: "Rust binary missing: .bin/fib_rec_rs",
		})
	}

	// Stable order in output
	sort.SliceStable(t, func(i, j int) bool {
		order := map[string]int{
			"go": 0, "vm": 1, "tenge-aot": 2, "c": 3, "rust": 4,
		}
		return order[t[i].Name] < order[t[j].Name]
	})
	return t
}

// -----------------------------
// Driver
// -----------------------------

// --- ИЗМЕНЕНИЕ ЗДЕСЬ ---
func runTask(title string, ns []int, kind string) ([]csvRow, error) {
	fmt.Printf("\nTask = %s\n", title)
	fmt.Println("──────────────────────────────────────────────────────────")
	fmt.Println(timingNote)
	// Обновляем заголовок таблицы
	fmt.Printf("%-20s %-15s %-12s\n\n", "Target", "Time", "Output")

	rows := make([]csvRow, 0, len(ns)*6)
	tgts := targetsFor(kind)

	for _, n := range ns {
		fmt.Printf("N = %d\n", n)
		fmt.Println("──────────────────────────────────────────────────────────")
		for _, t := range tgts {
			if t.Kind != kind {
				continue
			}
			available := true
			if t.OnlyIf != nil {
				available = t.OnlyIf()
			}
			if !available {
				fmt.Printf("%-20s %-15s [SKIP] %s\n\n", t.Name, "0 ns", t.SkipMsg)
				rows = append(rows, csvRow{
					Target: t.Name, N: n, TimeNS: "", TimeS: "", WallS: "0.000000", Result: "", Status: "SKIP",
				})
				continue
			}
			args := []string{}
			if t.ArgsFn != nil {
				args = t.ArgsFn(n)
			}
			// Pass through BENCH_REPS if user set it
			passEnv := map[string]string{}
			if v := os.Getenv("BENCH_REPS"); v != "" {
				passEnv["BENCH_REPS"] = v
			}

			out := runOne(t.Bin, args, passEnv)
			timeCol, outCol := fmtTimeOutcome(out)
			fmt.Printf("%-20s %-15s %s\n\n", t.Name, timeCol, outCol)

			row := csvRow{
				Target: t.Name,
				N:      n,
				Result: out.Result,
				Status: "OK",
				WallS:  fmt.Sprintf("%0.6f", out.WallSec),
			}
			if !out.Ok {
				row.Status = "ERR"
			}
			if out.TimeNS != nil {
				row.TimeNS = fmt.Sprintf("%d", *out.TimeNS)
			}
			if out.TimeS != nil {
				row.TimeS = fmt.Sprintf("%0.9f", *out.TimeS)
			}
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func main() {
	// Ensure roots exist
	_ = ensureDir(filepath.Join(latestRoot, resultsLeaf))

	ts := nowStamp()

	var allErrs []error

	// fib_rec
	recRows, err := runTask("fib_rec", fibRecNs, "rec")
	if err != nil {
		allErrs = append(allErrs, err)
	}
	if err := writeCSV(recRows, "fib_rec", ts); err != nil {
		allErrs = append(allErrs, err)
	}

	// fib_iter
	iterRows, err := runTask("fib_iter", fibIterNs, "iter")
	if err != nil {
		allErrs = append(allErrs, err)
	}
	if err := writeCSV(iterRows, "fib_iter", ts); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) > 0 {
		var b strings.Builder
		for _, e := range allErrs {
			b.WriteString(e.Error())
			b.WriteString("; ")
		}
		// return non-zero to signal overall issues
		fmt.Fprintln(os.Stderr, "benchfast finished with errors:", b.String())
		os.Exit(1)
	}
}