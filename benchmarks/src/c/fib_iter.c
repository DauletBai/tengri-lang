#include "runtime.h"

// The core logic is now in its own function.
long long run_fib_iter(int n) {
    if (n < 2) {
        return n;
    }
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
    
    // The macro now receives a simple, single function call.
    TIME_IT_NS(
        (void)run_fib_iter(n);,
        "fib_iter_c",
        n
    );
    
    return 0;
}