package main

import (
	"fmt"
	"os"
	"io/ioutil"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: <program> file")
		os.Exit(1)
	}
	
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	memory := make([]uint16, len(data)/2)
	for i := range memory {
		k := i*2
		memory[i] = uint16(data[k]) | (uint16(data[k+1]) << 8)
	}

	pgm := program{
		halted : false,
		pointer : 0,
		memory : memory,
		stack : []uint16{},
		registers: make([]uint16, 8),
	}

	pgm.Run()

}

type program struct {
	halted bool
	pointer uint16
	memory []uint16
	stack []uint16
	registers []uint16
}

var operations []func(*program) = []func(*program){
	(*program).Halt,
	(*program).Set,
	(*program).Push,
	(*program).Pop,
	(*program).Eq,
	(*program).Gt,
	(*program).Jmp,
	(*program).Jt,
	(*program).Jf,
	(*program).Add,
	(*program).Mult,
	(*program).Mod,
	(*program).And,
	(*program).Or,
	(*program).Not,
	(*program).Rmem,
	(*program).Wmem,
	(*program).Call,
	(*program).Ret,
	(*program).Out,
	(*program).In,
	(*program).Noop,
}

func (pgm *program) Run() {

	for int(pgm.pointer) < len(pgm.memory) {
		// fmt.Println("pointer:",pgm.pointer)
		opcode := pgm.Next()
		// if opcode != 21 { fmt.Println(opcode) }
		if int(opcode) >= len(operations) {
			panic("Unknown Opcode!")
		}
		operations[opcode](pgm)
		if pgm.halted { break }
	}
	
	fmt.Println("--------------\nProgram exited")
}

func (pgm *program) Halt() {
	pgm.halted = true
}

func (pgm *program) Set() {
	a, b := pgm.Next(), pgm.Next()
	value := pgm.OperandGet(b)
	pgm.OperandSet(a, value)
}

func (pgm *program) Push() {
	a := pgm.Next()
	value := pgm.OperandGet(a)
	pgm.stack = append(pgm.stack, value)
}

func (pgm *program) Pop() {
	if len(pgm.stack) <= 0 { panic("Pop: empty stack!!") }
	a := pgm.Next()
	slen := len(pgm.stack)
	value := pgm.stack[slen-1]
	pgm.stack = pgm.stack[:slen-1]
	pgm.OperandSet(a, value)
}

func (pgm *program) Eq() {
	a, b, c := pgm.Next(), pgm.Next(), pgm.Next()
	bValue, cValue := pgm.OperandGet(b), pgm.OperandGet(c)
	value := uint16(0)
	if bValue == cValue {
		value = 1
	}
	pgm.OperandSet(a, value)
}

func (pgm *program) Gt() {
	a, b, c := pgm.Next(), pgm.Next(), pgm.Next()
	bValue, cValue := pgm.OperandGet(b), pgm.OperandGet(c)
	value := uint16(0)
	if bValue > cValue {
		value = 1
	}
	pgm.OperandSet(a, value)
}

func (pgm *program) Jmp() {
	a := pgm.Next()
	pgm.pointer = pgm.OperandGet(a)
}

func (pgm *program) Jt() {
	a, b := pgm.Next(), pgm.Next()
	aValue := pgm.OperandGet(a)
	if aValue != 0 {
		bValue := pgm.OperandGet(b)
		pgm.pointer = bValue
	}
}

func (pgm *program) Jf() {
	a, b := pgm.Next(), pgm.Next()
	aValue := pgm.OperandGet(a)
	if aValue == 0 {
		bValue := pgm.OperandGet(b)
		pgm.pointer = bValue
	}
}

func (pgm *program) Add() {
	a, b, c := pgm.Next(), pgm.Next(), pgm.Next()
	bValue, cValue := pgm.OperandGet(b), pgm.OperandGet(c)
	value := (bValue + cValue) & 0x7fff
	pgm.OperandSet(a, value)
}

func (pgm *program) Mult() {
	a, b, c := pgm.Next(), pgm.Next(), pgm.Next()
	bValue, cValue := pgm.OperandGet(b), pgm.OperandGet(c)
	value := (bValue * cValue) & 0x7fff
	pgm.OperandSet(a, value)
}

func (pgm *program) Mod() {
	a, b, c := pgm.Next(), pgm.Next(), pgm.Next()
	bValue, cValue := pgm.OperandGet(b), pgm.OperandGet(c)
	value := bValue % cValue
	pgm.OperandSet(a, value)
}

func (pgm *program) And() {
	a, b, c := pgm.Next(), pgm.Next(), pgm.Next()
	bValue, cValue := pgm.OperandGet(b), pgm.OperandGet(c)
	value := bValue & cValue
	pgm.OperandSet(a, value)
}

func (pgm *program) Or() {
	a, b, c := pgm.Next(), pgm.Next(), pgm.Next()
	bValue, cValue := pgm.OperandGet(b), pgm.OperandGet(c)
	value := bValue | cValue
	pgm.OperandSet(a, value)
}

func (pgm *program) Not() {
	a, b := pgm.Next(), pgm.Next()
	bValue := pgm.OperandGet(b)
	value := (^ bValue) & 0x7fff
	pgm.OperandSet(a, value)
}

func (pgm *program) Rmem() {
	a, b := pgm.Next(), pgm.Next()
	bValue := pgm.OperandGet(b)
	value := pgm.memory[bValue]
	pgm.OperandSet(a, value)
}

func (pgm *program) Wmem() {
	a, b := pgm.Next(), pgm.Next()
	aValue, bValue := pgm.OperandGet(a), pgm.OperandGet(b)
	pgm.memory[aValue] = bValue
}

func (pgm *program) Call() {
	a := pgm.Next()
	aValue := pgm.OperandGet(a)
	pgm.stack = append(pgm.stack, pgm.pointer)
	pgm.pointer = aValue
}

func (pgm *program) Ret() {
	if len(pgm.stack) <= 0 {
		pgm.halted = true
		return
	}
	slen := len(pgm.stack)
	pgm.pointer = pgm.stack[slen-1]
	pgm.stack = pgm.stack[:slen-1]
}

func (pgm *program) Out() {
	a := pgm.Next()
	aValue := pgm.OperandGet(a)
	fmt.Printf("%c", aValue)
}

func (pgm *program) In() {
	a := pgm.Next()
	value := uint16(readStdinChar())
	pgm.OperandSet(a, value)
}

func (pgm *program) Noop() {}





func (pgm *program) Next() uint16 {
	num := pgm.memory[pgm.pointer]
	pgm.pointer = pgm.pointer + 1
	return num
}

func (pgm *program) OperandSet(num, value uint16) {
	if num >= 0x8000 {
		pgm.registers[num-0x8000] = value
	}	else {
		pgm.memory[num] = value
	}
}

func (pgm *program) OperandGet(num uint16) uint16 {
	if num >= 0x8000 {
		return pgm.registers[num-0x8000]
	}
	return num
}