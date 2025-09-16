// FILE: internal/aotminic/runtime/runtime.c

#include "runtime.h"

int get_n(int argc, char** argv, int default_n) {
    if (argc > 1) {
        return atoi(argv[1]);
    }
    return default_n;
}

int* create_array(int n) {
    int* arr = (int*)malloc(n * sizeof(int));
    if (arr == NULL) {
        fprintf(stderr, "Failed to allocate memory\n");
        exit(1);
    }
    for (int i = 0; i < n; i++) {
        arr[i] = i + 1;
    }
    return arr;
}