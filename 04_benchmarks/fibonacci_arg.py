import sys

def fib(n):
    if n < 2:
        return n
    return fib(n-1) + fib(n-2)

n = 35
if len(sys.argv) > 1:
    try:
        n = int(sys.argv[1])
    except:
        pass

print(fib(n))