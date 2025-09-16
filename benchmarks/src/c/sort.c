// FILE: benchmarks/src/c/sort.c

#include "runtime.h" // CORRECTED: Include the unified runtime header
#include <stdlib.h>

// Standard qsort comparison function
int compare(const void *a, const void *b) {
    return (*(int*)a - *(int*)b);
}

// The core logic is now wrapped in a function for clarity
void run_sort(int n, int* arr) {
    qsort(arr, n, sizeof(int), compare);
}

int main(int argc, char** argv) {
    int n = get_n(argc, argv, 100000);
    int* arr = create_array(n);
    
    // The macro now receives a simple, single function call,
    // ensuring its output matches what the run.sh script expects.
    TIME_IT_NS(
        run_sort(n, arr);,
        "sort_c",
        n
    );

    free(arr);
    return 0;
}