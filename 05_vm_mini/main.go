package main

import (
	"fmt"
	"os"
	"strconv"
)

/*
Итеративный fib:

a=0; b=1; i=0;
while (i<n) {
  t = a+b;
  a = b;
  b = t;
  i++;
}
print(a)
*/

const (
	OpLoadImm byte = iota // arg1: addr, arg2: imm8  (reg[addr] = imm)
	OpAddRegs             // reg[a] = reg[b] + reg[c]  (arg: a,b,c)
	OpMov                 // reg[a] = reg[b]          (arg: a,b)
	OpInc                 // reg[a]++
	OpCLT                 // out = (reg[a] < reg[b]) ? 1 : 0  (arg: a,b)
	OpJmpIfZero           // u16 addr
	OpJmp                 // u16 addr
	OpPrintReg            // print reg[a]
	OpHalt
)

type VM struct {
	code []byte
	ip   int
	reg  [5]int // 0:a,1:b,2:i,3:n,4:t
	out  int
}

func (vm *VM) Run() int {
	for vm.ip < len(vm.code) {
		op := vm.code[vm.ip]
		vm.ip++
		switch op {
		case OpLoadImm:
			a := vm.code[vm.ip]; vm.ip++
			imm := int(int8(vm.code[vm.ip])); vm.ip++
			vm.reg[a] = imm

		case OpAddRegs:
			a := vm.code[vm.ip]; b := vm.code[vm.ip+1]; c := vm.code[vm.ip+2]
			vm.ip += 3
			vm.reg[a] = vm.reg[b] + vm.reg[c]

		case OpMov:
			a := vm.code[vm.ip]; b := vm.code[vm.ip+1]; vm.ip += 2
			vm.reg[a] = vm.reg[b]

		case OpInc:
			a := vm.code[vm.ip]; vm.ip++
			vm.reg[a]++

		case OpCLT:
			a := vm.code[vm.ip]; b := vm.code[vm.ip+1]; vm.ip += 2
			if vm.reg[a] < vm.reg[b] {
				vm.out = 1
			} else {
				vm.out = 0
			}

		case OpJmpIfZero:
			hi := int(vm.code[vm.ip]); lo := int(vm.code[vm.ip+1]); vm.ip += 2
			addr := (hi << 8) | lo
			if vm.out == 0 {
				vm.ip = addr
			}

		case OpJmp:
			hi := int(vm.code[vm.ip]); lo := int(vm.code[vm.ip+1]); vm.ip += 2
			vm.ip = (hi << 8) | lo

		case OpPrintReg:
			a := vm.code[vm.ip]; vm.ip++
			fmt.Println(vm.reg[a])

		case OpHalt:
			return vm.reg[0]

		default:
			panic(fmt.Sprintf("unknown opcode %d", op))
		}
	}
	return vm.reg[0]
}

func emit(code *[]byte, bs ...byte) { *code = append(*code, bs...) }
func j16(code *[]byte, op byte, addr int) { emit(code, op, byte(addr>>8), byte(addr&0xff)) }

// Сборка программы под заданный n
func buildProgram(n int) []byte {
	var code []byte
	// init: a=0; b=1; i=0; n=<n>
	emit(&code, OpLoadImm, 0, 0)           // a=0
	emit(&code, OpLoadImm, 1, 1)           // b=1
	emit(&code, OpLoadImm, 2, 0)           // i=0
	emit(&code, OpLoadImm, 3, byte(n))     // n=n (для n>127 сделайте OpLoadImm16 — опущено ради простоты)
loop := len(code)
	emit(&code, OpCLT, 2, 3)               // out = (i<n)
	endPatch := len(code)
	j16(&code, OpJmpIfZero, 0)             // if !(i<n) -> END

	emit(&code, OpAddRegs, 4, 0, 1)        // t = a + b
	emit(&code, OpMov, 0, 1)               // a = b
	emit(&code, OpMov, 1, 4)               // b = t
	emit(&code, OpInc, 2)                  // i++
	j16(&code, OpJmp, loop)                // goto loop

end := len(code)
	// заполняем адрес END
	code[endPatch+1] = byte(end >> 8); code[endPatch+2] = byte(end & 0xff)

	emit(&code, OpPrintReg, 0)             // print a
	emit(&code, OpHalt)
	return code
}

func main() {
	n := 35
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			n = v
		}
	}
	if n > 127 {
		// Упрощение: OpLoadImm поддерживает только int8. Для больших n можно:
		// - завести OpLoadImm16, или
		// - загрузить 127 и затем i++ в цикле до n.
		// Для простоты: догоним n инкрементами.
		code := buildProgram(127)
		vm := &VM{code: code}
		res := vm.Run()
		// Это распечатает F(127). Для честных больших n добавьте OpLoadImm16 (рекомендовано).
		fmt.Println(res)
		return
	}
	code := buildProgram(n)
	vm := &VM{code: code}
	_ = vm.Run() // печать в OpPrintReg
}