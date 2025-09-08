.PHONY: bench go-native tengri python clean bench-fast

go-native:
	@echo "--- Native Go ---"
	time go run 04_benchmarks/fibonacci.go

tengri:
	@echo "\n--- Tengri-Lang (Go Interpreter) ---"
	cd 03_compiler_go && time go run . && cd ..

python:
	@echo "\n--- Python Interpreter ---"
	time python3 04_benchmarks/fibonacci.py

bench: go-native tengri python

bench-fast:
	@go run tools/benchfast/main.go

clean:
	rm -f 03_compiler_go/.tmp_fib_ascii.tg