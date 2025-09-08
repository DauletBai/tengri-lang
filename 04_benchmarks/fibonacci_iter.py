import sys

def fib_iter(n: int) -> int:
    if n < 2:
        return n
    a, b = 0, 1
    for _ in range(n):
        a, b = b, a + b
    return a

n = 300
if len(sys.argv) > 1:
    try:
        n = int(sys.argv[1])
    except:
        pass

print(fib_iter(n))