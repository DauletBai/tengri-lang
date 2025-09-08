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

bench-plot:
	@go run tools/benchfast/main.go --plot

clean:
	rm -f 03_compiler_go/.tmp_fib_ascii.tg

bench-fast:
	@go run tools/benchfast/main.go

bench-plot:
	@go run tools/benchfast/main.go --plot

bench-commit:
	@go run tools/benchfast/main.go --plot
	@git add benchmarks
	@git commit -m "bench: $(shell date +%Y-%m-%dT%H:%M:%S) run saved under benchmarks/runs and updated benchmarks/latest" || true