package aotminic

// C code templates for AOT compilation demos.
// These templates now also follow the best practice of separating logic into functions.

const FibIterC = `
#include "runtime.h"

long long run_fib_iter(int n) {
    if (n < 2) { return n; }
    long long a = 0, b = 1;
    for (int i = 2; i <= n; i++) {
        long long temp = a + b;
        a = b;
        b = temp;
    }
    return b;
}

int main(int argc, char** argv) {
    int n = get_n(argc, argv, 45);
    TIME_IT_NS(
        (void)run_fib_iter(n);,
        "fib_iter_tengri_aot",
        n
    );
    return 0;
}
`

const FibRecC = `
#include "runtime.h"

long long fib(int n) {
    if (n < 2) { return n; }
    return fib(n-1) + fib(n-2);
}

int main(int argc, char** argv) {
    int n = get_n(argc, argv, 35);
    TIME_IT_NS(
        (void)fib(n);,
        "fib_rec_tengri_aot",
        n
    );
    return 0;
}
`

const SortQSortC = `
#include "runtime.h"
#include <stdlib.h>

int compare(const void *a, const void *b) {
    return (*(int*)a - *(int*)b);
}

void run_qsort(int n, int* arr) {
    qsort(arr, n, sizeof(int), compare);
}

int main(int argc, char** argv) {
    int n = get_n(argc, argv, 100000);
    int* arr = create_array(n);
    TIME_IT_NS(
        run_qsort(n, arr);,
        "sort_qsort_tengri_aot",
        n
    );
    free(arr);
    return 0;
}
`

const SortMergeSortC = `
#include "runtime.h"
#include <stdlib.h>

void merge(int arr[], int l, int m, int r) {
    int i, j, k;
    int n1 = m - l + 1;
    int n2 = r - m;
    int *L = malloc(n1 * sizeof(int));
    int *R = malloc(n2 * sizeof(int));
    for (i = 0; i < n1; i++) L[i] = arr[l + i];
    for (j = 0; j < n2; j++) R[j] = arr[m + 1 + j];
    i = 0; j = 0; k = l;
    while (i < n1 && j < n2) {
        if (L[i] <= R[j]) arr[k++] = L[i++];
        else arr[k++] = R[j++];
    }
    while (i < n1) arr[k++] = L[i++];
    while (j < n2) arr[k++] = R[j++];
    free(L);
    free(R);
}

void mergeSort(int arr[], int l, int r) {
    if (l < r) {
        int m = l + (r - l) / 2;
        mergeSort(arr, l, m);
        mergeSort(arr, m + 1, r);
        merge(arr, l, m, r);
    }
}

void run_msort(int n, int* arr) {
    mergeSort(arr, 0, n - 1);
}

int main(int argc, char** argv) {
    int n = get_n(argc, argv, 100000);
    int* arr = create_array(n);
    TIME_IT_NS(
        run_msort(n, arr);,
        "sort_msort_tengri_aot",
        n
    );
    free(arr);
    return 0;
}
`