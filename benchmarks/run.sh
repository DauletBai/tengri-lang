#!/bin/bash

# --- Tengri Benchmark Runner ---
# This script runs all compiled benchmark binaries and collects their results.

# Exit immediately if a command fails
set -e

# --- Configuration ---
BIN_DIR=".bin"
RESULTS_DIR="benchmarks/results"
SIZE=${1:-100000} # Default size for tasks like sorting
REPS=${2:-5}     # Default reps
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
CSV_FILE="${RESULTS_DIR}/suite_${TIMESTAMP}.csv"

# --- Specific N values for different tasks ---
FIB_ITER_N=90  # A large N is fine for iterative fib
FIB_REC_N=35   # A small N is CRITICAL for recursive fib

# --- Setup ---
mkdir -p "$RESULTS_DIR"
echo "impl,task,n,reps,time_ns_avg" > "$CSV_FILE" # CSV Header
echo "[bench] Running all benchmarks with REPS=${REPS}"
echo "[bench] Sort tasks use SIZE=${SIZE}. Fib tasks use fixed N values."
echo "[bench] Results will be saved to ${CSV_FILE}"

# --- Helper function to run a binary and parse its output ---
run_and_parse() {
    local binary_path=$1
    local task_name=$2
    local implementation_name=$3
    local n_value=$4 # The specific N to pass to the binary

    if [ ! -f "$binary_path" ]; then
        echo " -> Skipping ${binary_path} (not found)"
        return
    fi

    echo " -> Running ${task_name} for ${implementation_name} with N=${n_value}..."

    total_ns=0
    for (( i=0; i<$REPS; i++ )); do
        # Capture the line with TASK=...,N=...,TIME_NS=...
        output=$($binary_path $n_value 2>/dev/null | grep 'TASK=')
        
        # Extract the time in nanoseconds
        time_ns=$(echo "$output" | sed -n 's/.*TIME_NS=\([0-9]*\).*/\1/p')
        total_ns=$((total_ns + time_ns))
    done

    # Calculate the average
    avg_ns=$((total_ns / REPS))

    # Write to CSV
    echo "${implementation_name},${task_name},${n_value},${REPS},${avg_ns}" >> "$CSV_FILE"
    echo "    ... OK! avg=${avg_ns} ns"
}


# --- Run all benchmarks ---

# C
run_and_parse "${BIN_DIR}/sort_c" "sort" "c" "$SIZE"
run_and_parse "${BIN_DIR}/fib_iter_c" "fib_iter" "c" "$FIB_ITER_N"
run_and_parse "${BIN_DIR}/fib_rec_c" "fib_rec" "c" "$FIB_REC_N"

# Go
run_and_parse "${BIN_DIR}/sort_go" "sort" "go" "$SIZE"
run_and_parse "${BIN_DIR}/fib_iter_go" "fib_iter" "go" "$FIB_ITER_N"
run_and_parse "${BIN_DIR}/fib_rec_go" "fib_rec" "go" "$FIB_REC_N"

# Rust
run_and_parse "${BIN_DIR}/sort_rs" "sort" "rust" "$SIZE"
run_and_parse "${BIN_DIR}/fib_iter_rs" "fib_iter" "rust" "$FIB_ITER_N"
run_and_parse "${BIN_DIR}/fib_rec_rs" "fib_rec" "rust" "$FIB_REC_N"

# Tengri AOT
run_and_parse "${BIN_DIR}/sort_cli_qsort" "sort_qsort" "tengri-aot" "$SIZE"
run_and_parse "${BIN_DIR}/sort_cli_msort" "sort_msort" "tengri-aot" "$SIZE"
run_and_parse "${BIN_DIR}/fib_cli" "fib_iter" "tengri-aot" "$FIB_ITER_N"
run_and_parse "${BIN_DIR}/fib_rec_cli" "fib_rec" "tengri-aot" "$FIB_REC_N"

echo "[bench] All benchmarks finished successfully."