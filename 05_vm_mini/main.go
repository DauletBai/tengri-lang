package main

import (
	"fmt"
	"os"
	"strconv"
)

/*
Итеративный fib на простейшем байткоде:

a = 0; b = 1; i = 0;
while i < n {
    a, b = b, a+b
    i++
}
print(a)

Опкоды работают со стеком int. Чтобы упростить, используем
4 "псевдо-регистра" в вершине стека в порядке: [a b i n]
и будем их читать/писать через адресные операции.
*/

const (
	OpPushConst byte = iota // arg1: const value (signed int8)
	OpDup                   // дублировать вершину
	OpSwap                  // свапнуть два верхних
	OpPop
	OpLoad                  // arg1: addr(0=a,1=b,2=i,3=n) -> push value
	OpStore                 // arg1: addr -> pop to that slot
	OpAdd
	OpIncI                  // i++
	OpCLT                   // сравнить: (top-1) < (top) -> push 1/0; pop2
	OpJmpIfZero             // arg2: u16 offset (big-endian)
	OpJmp                   // arg2: u16 offset
	OpPrint                 // вывести вершину (оставим на стеке)
	OpHalt
)

type VM struct {
	code []byte
	ip   int
	// регистры: a,b,i,n
	reg [4]int
	out int
}

func (vm *VM) push(v int) {
	// уже держим регистры вне стека: стек нам не нужен, оставим минимум
	vm.out = v
}

func (vm *VM) pop() int {
	// out как единственная вершина для нужд операций
	return vm.out
}

func (vm *VM) Run() int {
	for vm.ip < len(vm.code) {
		op := vm.code[vm.ip]
		vm.ip++
		switch op {
		case OpPushConst:
			val := int(int8(vm.code[vm.ip]))
			vm.ip++
			vm.push(val)

		case OpDup:
			vm.push(vm.out)

		case OpSwap:
			// Не требуется для текущей программы (оставлено для примера)
			// Ничего не делаем
		case OpPop:
			vm.out = 0

		case OpLoad:
			addr := vm.code[vm.ip]
			vm.ip++
			vm.push(vm.reg[addr])

		case OpStore:
			addr := vm.code[vm.ip]
			vm.ip++
			vm.reg[addr] = vm.pop()

		case OpAdd:
			// используем out как аккумулятор, прибавим b к a и положим в out
			// Для простоты: out = a + b
			vm.push(vm.reg[0] + vm.reg[1])

		case OpIncI:
			vm.reg[2]++

		case OpCLT:
			// сравнение (i < n): push 1 или 0
			// предполагаем, что out содержит n (или i), но сделаем честно
			// тут читаем i и n из регистров
			if vm.reg[2] < vm.reg[3] {
				vm.push(1)
			} else {
				vm.push(0)
			}

		case OpJmpIfZero:
			off := int(vm.code[vm.ip])<<8 | int(vm.code[vm.ip+1])
			vm.ip += 2
			if vm.pop() == 0 {
				vm.ip = off
			}

		case OpJmp:
			off := int(vm.code[vm.ip])<<8 | int(vm.code[vm.ip+1])
			vm.ip = off

		case OpPrint:
			// вывод через fmt.Println в конце (вернём значение)
			// здесь просто оставим out как есть
			_ = vm.out

		case OpHalt:
			return vm.out

		default:
			panic(fmt.Sprintf("unknown opcode %d", op))
		}
	}
	return vm.out
}

func buildProgram(n int) []byte {
	// Программа:
	// a=0; b=1; i=0; n=arg
	// LOOP:
	// if i < n else -> END
	// tmp = a + b
	// a = b
	// b = tmp
	// i++
	// jmp LOOP
	// END: print a; halt

	code := []byte{}

	emit := func(bs ...byte) { code = append(code, bs...) }
	j2 := func(op byte, off int) { emit(op, byte(off>>8), byte(off&0xff)) }

	// init
	emit(OpPushConst, 0)  // 0
	emit(OpStore, 0)      // a = 0
	emit(OpPushConst, 1)  // 1
	emit(OpStore, 1)      // b = 1
	emit(OpPushConst, 0)  // 0
	emit(OpStore, 2)      // i = 0
	emit(OpPushConst, byte(n))
	emit(OpStore, 3)      // n = arg

	loop := len(code)
	emit(OpCLT)                 // out = (i<n)?1:0
	// if zero -> END
	endPatch := len(code)
	j2(OpJmpIfZero, 0)          // заполнится позже

	emit(OpAdd)                 // out = a+b
	emit(OpStore, 2)            // (временно положим в i, чтобы не плодить опкоды) tmp := out -> рег i
	// a = b
	emit(OpLoad, 1)
	emit(OpStore, 0)
	// b = tmp
	emit(OpLoad, 2)
	emit(OpStore, 1)
	// i++
	emit(OpIncI)
	// jmp loop
	j2(OpJmp, loop)

	end := len(code)
	// заполняем адрес END
	code[endPatch+1] = byte(end >> 8)
	code[endPatch+2] = byte(end & 0xff)

	// print a; halt
	emit(OpLoad, 0)
	emit(OpPrint)
	emit(OpHalt)
	return code
}

func main() {
	n := 35
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			n = v
		}
	}
	vm := &VM{code: buildProgram(n)}
	res := vm.Run()
	fmt.Println(res)
}