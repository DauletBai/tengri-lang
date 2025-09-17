#!/bin/bash
# FILE: benchmarks/run.sh
# Purpose: Run all available benchmarks and aggregate results to CSV.

set -euo pipefail

BIN_DIR=".bin"
RESULTS_DIR="benchmarks/results"
mkdir -p "$RESULTS_DIR"

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
CSV_FILE="${RESULTS_DIR}/suite_${TIMESTAMP}.csv"

SIZE=${SIZE:-100000}
REPS=${REPS:-5}
FIB_ITER_N=${FIB_ITER_N:-90}
FIB_REC_N=${FIB_REC_N:-35}

VAR_N=${VAR_N:-1000000}
VAR_STEPS=${VAR_STEPS:-1}
VAR_ALPHA=${VAR_ALPHA:-0.99}

echo "[bench] Running all benchmarks with REPS=${REPS}"
echo "[bench] Sort tasks use SIZE=${SIZE}. Fib tasks use fixed N values."
echo "[bench] Results will be saved to ${CSV_FILE}"

echo "ts,task,impl,paramN,time_ns,extra" > "$CSV_FILE"

run_and_parse() {
  local bin="$1"
  local task="$2"
  local impl="$3"
  shift 3
  local args=( "$@" )
  if [[ ! -x "$bin" ]]; then
    echo " -> Skipping ${bin} (not found)"
    return 0
  fi
  local acc=0
  local extra_last=""
  for ((i=1; i<=REPS; i++)); do
    local out
    if ! out=$("$bin" "${args[@]}" 2>/dev/null); then
      echo " -> ERROR running ${bin}" >&2
      return 1
    fi
    local time_ns
    time_ns=$(echo "$out" | sed -n 's/.*TIME_NS=\([0-9][0-9]*\).*/\1/p' | head -n1)
    local nparam
    nparam=$(echo "$out" | sed -n 's/.*N=\([0-9][0-9]*\).*/\1/p' | head -n1)
    [[ -z "$nparam" ]] && nparam="${args[0]}"
    local varfield
    varfield=$(echo "$out" | sed -n 's/.*VAR=\([0-9.][0-9.]*\).*/VAR=\1/p' | head -n1)
    extra_last="$varfield"
    [[ -z "$time_ns" ]] && time_ns=0
    acc=$((acc + time_ns))
    echo "$(date +%s),${task},${impl},${nparam},${time_ns},${extra_last}" >> "$CSV_FILE"
  done
  local avg=$(( acc / REPS ))
  echo "    ... OK! avg=${avg} ns"
}

# C
echo " -> Running sort for c with N=${SIZE}..."
run_and_parse "${BIN_DIR}/sort_c" "sort" "c" "$SIZE"
echo " -> Running fib_iter for c with N=${FIB_ITER_N}..."
run_and_parse "${BIN_DIR}/fib_iter_c" "fib_iter" "c" "$FIB_ITER_N"
echo " -> Running fib_rec for c with N=${FIB_REC_N}..."
run_and_parse "${BIN_DIR}/fib_rec_c" "fib_rec" "c" "$FIB_REC_N"

# Go
echo " -> Running sort for go with N=${SIZE}..."
run_and_parse "${BIN_DIR}/sort_go" "sort" "go" "$SIZE"
echo " -> Running fib_iter for go with N=${FIB_ITER_N}..."
run_and_parse "${BIN_DIR}/fib_iter_go" "fib_iter" "go" "$FIB_ITER_N"
echo " -> Running fib_rec for go with N=${FIB_REC_N}..."
run_and_parse "${BIN_DIR}/fib_rec_go" "fib_rec" "go" "$FIB_REC_N"

# Rust
echo " -> Running sort for rust with N=${SIZE}..."
run_and_parse "${BIN_DIR}/sort_rs" "sort" "rust" "$SIZE"
echo " -> Running fib_iter for rust with N=${FIB_ITER_N}..."
run_and_parse "${BIN_DIR}/fib_iter_rs" "fib_iter" "rust" "$FIB_ITER_N"
echo " -> Running fib_rec for rust with N=${FIB_REC_N}..."
run_and_parse "${BIN_DIR}/fib_rec_rs" "fib_rec" "rust" "$FIB_REC_N"

# Tenge (impl='tenge')
echo " -> Running sort_qsort for tenge with N=${SIZE}..."
run_and_parse "${BIN_DIR}/sort_cli_qsort" "sort_qsort" "tenge" "$SIZE"
echo " -> Running sort_msort for tenge with N=${SIZE}..."
run_and_parse "${BIN_DIR}/sort_cli_msort" "sort_msort" "tenge" "$SIZE"
echo " -> Running sort_pdq for tenge with N=${SIZE}..."
run_and_parse "${BIN_DIR}/sort_cli_pdq" "sort_pdq" "tenge" "$SIZE"
echo " -> Running sort_radix for tenge with N=${SIZE}..."
run_and_parse "${BIN_DIR}/sort_cli_radix" "sort_radix" "tenge" "$SIZE"
echo " -> Running fib_iter for tenge with N=${FIB_ITER_N}..."
run_and_parse "${BIN_DIR}/fib_cli" "fib_iter" "tenge" "$FIB_ITER_N"
echo " -> Running fib_rec for tenge with N=${FIB_REC_N}..."
run_and_parse "${BIN_DIR}/fib_rec_cli" "fib_rec" "tenge" "$FIB_REC_N"

# Monte Carlo
echo " -> Running var_mc_sort for tenge with N=${VAR_N}..."
run_and_parse "${BIN_DIR}/var_mc_tng_sort" "var_mc_sort" "tenge" "$VAR_N" "$VAR_STEPS" "$VAR_ALPHA"
echo " -> Running var_mc_zig for tenge with N=${VAR_N}..."
run_and_parse "${BIN_DIR}/var_mc_tng_zig" "var_mc_zig" "tenge" "$VAR_N" "$VAR_STEPS" "$VAR_ALPHA"
echo " -> Running var_mc_qsel for tenge with N=${VAR_N}..."
run_and_parse "${BIN_DIR}/var_mc_tng_qsel" "var_mc_qsel" "tenge" "$VAR_N" "$VAR_STEPS" "$VAR_ALPHA"

echo " -> Running var_mc for c with N=${VAR_N}..."
run_and_parse "${BIN_DIR}/var_mc_c" "var_mc" "c" "$VAR_N" "$VAR_STEPS" "$VAR_ALPHA"
echo " -> Running var_mc for go with N=${VAR_N}..."
run_and_parse "${BIN_DIR}/var_mc_go" "var_mc" "go" "$VAR_N" "$VAR_STEPS" "$VAR_ALPHA"
echo " -> Running var_mc for rust with N=${VAR_N}..."
run_and_parse "${BIN_DIR}/var_mc_rs" "var_mc" "rust" "$VAR_N" "$VAR_STEPS" "$VAR_ALPHA"

echo "[bench] All benchmarks finished successfully."