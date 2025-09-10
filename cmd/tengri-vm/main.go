package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const MOD = 1000000007

// Opcodes
const (
	OpLoadImm byte = iota // reg[a] = imm8 (демо; для N>127 нужно будет добавить Imm16)
	OpAddMod              // reg[a] = (reg[b] + reg[c]) % MOD
	OpMov                 // reg[a] = reg[b]
	OpInc                 // reg[a]++
	OpCLT                 // out = (reg[a] < reg[b]) ? 1 : 0
	OpJmpIfZero           // addr u16
	OpJmp                 // addr u16
	OpReadReg             // out = reg[a]
	OpHalt
)

type VM struct {
	code []byte
	ip   int
	reg  [5]int // 0:a,1:b,2:i,3:n,4:t
	out  int
}

func (vm *VM) Run() {
	vm.ip = 0
	for vm.ip < len(vm.code) {
		op := vm.code[vm.ip]
		vm.ip++
		switch op {
		case OpLoadImm:
			a := vm.code[vm.ip]; vm.ip++
			imm := int(int8(vm.code[vm.ip])); vm.ip++
			vm.reg[a] = ((imm % MOD) + MOD) % MOD

		case OpAddMod:
			a := vm.code[vm.ip]; b := vm.code[vm.ip+1]; c := vm.code[vm.ip+2]; vm.ip += 3
			vm.reg[a] = (vm.reg[b] + vm.reg[c]) % MOD

		case OpMov:
			a := vm.code[vm.ip]; b := vm.code[vm.ip+1]; vm.ip += 2
			vm.reg[a] = vm.reg[b]

		case OpInc:
			a := vm.code[vm.ip]; vm.ip++
			vm.reg[a] = (vm.reg[a] + 1) % MOD

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

		case OpReadReg:
			a := vm.code[vm.ip]; vm.ip++
			vm.out = vm.reg[a]

		case OpHalt:
			return

		default:
			panic(fmt.Sprintf("unknown opcode %d", op))
		}
	}
}

func emit(code *[]byte, bs ...byte) { *code = append(*code, bs...) }
func j16(code *[]byte, op byte, addr int) { emit(code, op, byte(addr>>8), byte(addr&0xff)) }

// a=0; b=1; i=0; while(i<n){ t=a+b; a=b; b=t; i++; } out=a
func buildProgram(n int) []byte {
	var code []byte
	emit(&code, OpLoadImm, 0, 0)       // a=0
	emit(&code, OpLoadImm, 1, 1)       // b=1
	emit(&code, OpLoadImm, 2, 0)       // i=0
	if n > 127 { n = 127 }             // для демо Imm8; при желании добавим Imm16
	emit(&code, OpLoadImm, 3, byte(n)) // n=n

loop := len(code)
	emit(&code, OpCLT, 2, 3)           // out = (i<n)
	endPatch := len(code)
	j16(&code, OpJmpIfZero, 0)         // if !(i<n) -> END

	emit(&code, OpAddMod, 4, 0, 1)     // t = (a+b)%MOD
	emit(&code, OpMov, 0, 1)           // a = b
	emit(&code, OpMov, 1, 4)           // b = t
	emit(&code, OpInc, 2)              // i++
	j16(&code, OpJmp, loop)            // goto loop

end := len(code)
	code[endPatch+1] = byte(end >> 8); code[endPatch+2] = byte(end & 0xff)

	emit(&code, OpReadReg, 0)          // out = a
	emit(&code, OpHalt)
	return code
}

func pickReps(n int) int {
	base := 5_000_000
	scale := 50 / max(1, n)
	reps := base * scale
	if reps < 500_000 { reps = 500_000 }
	return reps
}
func max(a, b int) int { if a > b { return a }; return b }

func main() {
	n := 90
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			n = v
		}
	}
	reps := pickReps(n)
	if rs := os.Getenv("BENCH_REPS"); rs != "" {
		if v, err := strconv.Atoi(rs); err == nil && v > 0 {
			reps = v
		}
	}

	code := buildProgram(n)
	vm := &VM{code: code}

	t0 := time.Now()
	for i := 0; i < reps; i++ {
		vm.Run()
	}
	t1 := time.Now()
	perCall := float64(t1.Sub(t0).Nanoseconds()) / float64(reps) // ns/call

	fmt.Printf("RESULT: %d\n", vm.out)
	fmt.Printf("TIME: %.9f\n", perCall/1e9) 
	fmt.Printf("TIME_NS: %.0f\n", perCall)  
}