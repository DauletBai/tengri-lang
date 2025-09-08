// node 18+
const n = process.argv[2] ? parseInt(process.argv[2], 10) : 35;
function fib(k) { return k < 2 ? k : fib(k-1) + fib(k-2); }
console.log(fib(n));