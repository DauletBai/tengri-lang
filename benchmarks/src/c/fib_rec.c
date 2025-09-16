#include "runtime.h"

long long fib_rec(int n) {
    if (n < 2) {
        return n;
    }
    return fib_rec(n - 1) + fib_rec(n - 2);
}

int main(int argc, char** argv) {
    int n = get_n(argc, argv, 35);

    // The macro call is now simple and clean.
    TIME_IT_NS(
        (void)fib_rec(n);,
        "fib_rec_c",
        n
    );

    return 0;
}