package asm

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ordovician/calcutron/prog"
)

func Example_assembleLine() {
	labels := make(prog.SymbolTable)
	lines := [...]string{
		// all the standard instructions
		"ADD   x9, x8, x7",
		"ADDI  x8, 42",
		"SUB   x2, x4, x1",
		"SHIFT x5, x2, 4",
		"LOAD  x1, x2, x3",
		"MOVE  x1, 24",
		"STORE x5, x1, x2",
		"BEQ   x3, x2, x1",
		"BGT   x3, x2, x1",
		"JMP x9, 82",

		// using non standard number of operands
		"ADD   x3, x7",
		"SUB   x4, x3",
	}

	for _, line := range lines {
		instruction, _ := AssembleLine(labels, line, 0)

		machinecode := instruction.MachineCode()
		fmt.Println(machinecode, line)
	}

	// Output:
	// 1987 ADD   x9, x8, x7
	// 2842 ADDI  x8, 42
	// 3241 SUB   x2, x4, x1
	// 4524 SHIFT x5, x2, 4
	// 5123 LOAD  x1, x2, x3
	// 6124 MOVE  x1, 24
	// 7512 STORE x5, x1, x2
	// 8321 BEQ   x3, x2, x1
	// 9321 BGT   x3, x2, x1
	// 982 JMP x9, 82
	// 1337 ADD   x3, x7
	// 3443 SUB   x4, x3
}

func ExampleInstruction_SourceCode() {
	labels := make(prog.SymbolTable)
	lines := [...]string{
		// all the standard instructions
		"ADD   x9, x8, x7",
		"ADDI  x8, 42",
		"SUB   x2, x4, x1",
		"SHIFT x5, x2, 4",
		"LOAD  x1, x2, x3",
		"MOVE  x1, 24",
		"STORE x5, x1, x2",
		"BEQ   x3, x2, x1",
		"BGT   x3, x2, x1",
		"JMP x9, 82",

		// using non standard number of operands
		"ADD   x3, x7",
		"SUB   x4, x3",
	}

	for _, line := range lines {
		instruction, _ := AssembleLine(labels, line, 0)

		sourcecode := instruction.SourceCode()
		fmt.Println(sourcecode)
	}

	// Output:
	// ADD   x9, x8, x7
	// ADDI  x8, 42
	// SUB   x2, x4, x1
	// SHIFT x5, x2, 4
	// LOAD  x1, x2, x3
	// MOVE  x1, 24
	// STORE x5, x1, x2
	// BEQ   x3, x2, x1
	// BGT   x3, x2, x1
	// JMP x9, 82
	// ADD   x3, x7
	// SUB   x4, x3
}

func Example_parseLine() {
	lines := [...]string{
		"ADD   x9, x8, x7",
		"ADDI  x8, 42",
		"SUB   x2, x4, x1",
		"SHIFT x5, x2, 4",
		"LOAD  x1, x2, x3",
		"MOVE  x1, 24",
		"STORE x5, x1, x2",
		"BEQ   x3, x2, x1",
		"BGT   x3, x2, x1",
		"JMP x9, 82",
	}

	for _, line := range lines {
		mnemonic, operands := parseLine(line)

		fmt.Printf("%-5s %s\n", mnemonic, strings.Join(operands, ", "))
	}

	// Output:
	// ADD   x9, x8, x7
	// ADDI  x8, 42
	// SUB   x2, x4, x1
	// SHIFT x5, x2, 4
	// LOAD  x1, x2, x3
	// MOVE  x1, 24
	// STORE x5, x1, x2
	// BEQ   x3, x2, x1
	// BGT   x3, x2, x1
	// JMP x9, 82

}

func Example_assembleFile() {
	program, _ := AssembleFile("../testdata/adder.ct33")

	program.PrintWithOptions(os.Stdout, &prog.PrintOptions{
		SourceCode: true,
	})

	// Output:
	//     MOVE  x9, tape
	// loop:
	//     LOAD  x1, x9
	//     LOAD  x2, x9
	//     ADD   x1, x2
	//     STORE x1, x9
	//     JMP x0, loop
}

func TestHaltAsembly(t *testing.T) {
	AssembleFile("../Examples/memadd.ct33")
}
