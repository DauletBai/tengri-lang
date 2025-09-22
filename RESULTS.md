# Benchmark Results (3 runs average)

| task        | impl   |   paramN |      avg_time_ns |
|:------------|:-------|---------:|-----------------:|
| fib_iter    | c      |       90 |     41.00        |
| fib_iter    | go     |       90 |   5169.53        |
| fib_iter    | rust   |       90 |    197.133       |
| fib_iter    | tenge  |       90 |     50.0667      |
| fib_rec     | c      |       35 |      3.72818e+07 |
| fib_rec     | go     |       35 |      4.42755e+07 |
| fib_rec     | rust   |       35 |      4.06111e+07 |
| fib_rec     | tenge  |       35 |      3.8633e+07  |
| sort        | c      |   100000 | 494667.0         |
| sort        | go     |   100000 | 112847.0         |
| sort        | rust   |   100000 |      1.66633e+06 |
| sort_msort  | tenge  |   100000 |      3.49707e+06 |
| sort_pdq    | tenge  |   100000 | 816600.0         |
| sort_qsort  | tenge  |   100000 |      2.43993e+06 |
| sort_radix  | tenge  |   100000 |      1.54027e+06 |
| var_mc      | c      |  1000000 |      1.85065e+08 |
| var_mc      | go     |  1000000 |      1.76879e+08 |
| var_mc      | rust   |  1000000 |      8.45282e+07 |
| var_mc_qsel | tenge  |  1000000 |      3.0872e+07  |
| var_mc_sort | tenge  |  1000000 |      1.84507e+08 |
| var_mc_zig  | tenge  |  1000000 |      1.00724e+08 |

Link to the new project [tenge](https://github.com/DauletBai/tenge) and its benchmarks [tenge/benchmarks/results](https://github.com/DauletBai/tenge/tree/main/benchmarks/results)