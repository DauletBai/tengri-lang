BIN_DIR=./bin
BIN=$(BIN_DIR)/tengri-aot
CC=clang

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	GO111MODULE=on go build -o $(BIN) .

hello: build
	$(BIN) examples/hello.tgr -o $(BIN_DIR)/hello.c
	$(CC) -O2 -o $(BIN_DIR)/hello $(BIN_DIR)/hello.c runtime/runtime.c
	$(BIN_DIR)/hello

fib: build
	$(BIN) examples/fib_cli.tgr -o $(BIN_DIR)/fib_cli.c
	$(CC) -O2 -o $(BIN_DIR)/fib_cli $(BIN_DIR)/fib_cli.c runtime/runtime.c
	$(BIN_DIR)/fib_cli 40

clean:
	rm -rf $(BIN_DIR)