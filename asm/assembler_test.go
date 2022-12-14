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
		"ADD  x9, x8, x7",
		"ADDI x8, 42",
		"SUB  x2, x4, x1",
		"LSH x5, x2, 4",
		"LOAD x1, x2, 3",
		"LODI x1, 24",
		"STOR x5, x1, 2",
		"BEQ  x3, x2, 1",
		"BGT  x3, x2, 1",
		"JMP  x9, 82",

		// using non standard number of operands
		"ADD  x3, x7",
		"SUB  x4, x3",
	}

	for _, line := range lines {
		instruction, _ := AssembleLine(labels, line, 0)

		machinecode := instruction.MachineCode()
		fmt.Printf("%04d %s\n", machinecode, line)
	}

	// Output:
	// 1987 ADD  x9, x8, x7
	// 2842 ADDI x8, 42
	// 3241 SUB  x2, x4, x1
	// 4524 LSH x5, x2, 4
	// 5123 LOAD x1, x2, 3
	// 6124 LODI x1, 24
	// 7512 STOR x5, x1, 2
	// 0321 BEQ  x3, x2, 1
	// 9321 BGT  x3, x2, 1
	// 8982 JMP  x9, 82
	// 1337 ADD  x3, x7
	// 3443 SUB  x4, x3
}

func ExampleInstruction_SourceCode() {
	labels := make(prog.SymbolTable)
	lines := [...]string{
		// all the standard instructions
		"ADD  x9, x8, x7",
		"ADDI x8, 42",
		"SUB  x2, x4, x1",
		"LSH  x5, x2, 4",
		"LOAD x1, x2, 3",
		"LODI x1, 24",
		"STOR x5, x1, 2",
		"BEQ  x3, x2, 1",
		"BGT  x3, x2, 1",
		"JMP  x9, 82",

		// using non standard number of operands
		"ADD  x3, x7",
		"SUB  x4, x3",
	}

	for _, line := range lines {
		instruction, err := AssembleLine(labels, line, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		sourcecode := instruction.SourceCode()
		fmt.Println(sourcecode)
	}

	// Output:
	// ADD  x9, x8, x7
	// ADDI x8, 42
	// SUB  x2, x4, x1
	// LSH  x5, x2, 4
	// LOAD x1, x2, 3
	// LODI x1, 24
	// STOR x5, x1, 2
	// BEQ  x3, x2, 1
	// BGT  x3, x2, 1
	// JMP  x9, 82
	// ADD  x3, x7
	// SUB  x4, x3
}

func Example_parseLine() {
	lines := [...]string{
		"ADD  x9, x8, x7",
		"ADDI x8, 42",
		"SUB  x2, x4, x1",
		"LSH  x5, x2, 4",
		"LOAD x1, x2, 3",
		"LODI x1, 24",
		"STOR x5, x1, 2",
		"BEQ  x3, x2, 1",
		"BGT  x3, x2, 1",
		"JMP  x9, 82",
	}

	for _, line := range lines {
		mnemonic, operands := parseLine(line)

		fmt.Printf("%-4s %s\n", mnemonic, strings.Join(operands, ", "))
	}

	// Output:
	// ADD  x9, x8, x7
	// ADDI x8, 42
	// SUB  x2, x4, x1
	// LSH  x5, x2, 4
	// LOAD x1, x2, 3
	// LODI x1, 24
	// STOR x5, x1, 2
	// BEQ  x3, x2, 1
	// BGT  x3, x2, 1
	// JMP  x9, 82

}

func Example_assembleFile() {
	program, _ := AssembleFile("../examples/adder.ct33")

	program.PrintWithOptions(os.Stdout, &prog.PrintOptions{
		SourceCode: true,
	})

	// Output:
	// loop:
	//     LOAD x1, x0, -1
	//     LOAD x2, x0, -1
	//     ADD  x3, x1, x2
	//     STOR x3, x0, -1
	//     JMP  x0, loop
	//     HLT
}

// func TestHaltAsembly(t *testing.T) {
// 	srcfile := "adder.ct33"
// 	program, err := AssembleFile("../Examples/" + srcfile)
// 	if err != nil {
// 		t.Errorf("failed to assemble %s because %v", srcfile, err)
// 	}

// 	for _, inst := range program.Instructions {
// 		fmt.Println(inst.MachineCode(), inst.SourceCode())
// 	}
// }

func TestImmediateRange(t *testing.T) {
	labels := make(prog.SymbolTable)
	var err error

	// Here we use the range 0-99 which is legal for JMP
	_, err = AssembleLine(labels, "JMP  x9, 82", 0)
	if err != nil {
		t.Errorf("failed to assemble 'JMP x9, 82' becase %v", err)
	}

	_, err = AssembleLine(labels, "JMP  x9, -20", 0)
	if err == nil {
		t.Errorf("A 'JMP  x9, -20' should fail to assemble as jump is negative")
	}

	_, err = AssembleLine(labels, "LODI x5, 52", 0)
	if err == nil {
		t.Errorf("A 'LODI x5, 52' should fail to assemble as immediate value is outside range -50 to 49")
	}

	_, err = AssembleLine(labels, "LODI x5, 48", 0)
	if err != nil {
		t.Errorf("A 'LODI x5, 49' didn't assemble because %v", err)
	}

	// negative immediate values should be allowed
	_, err = AssembleLine(labels, "LODI x5, -20", 0)
	if err != nil {
		t.Errorf("A 'LODI x5, -20' didn't assemble because %v", err)
	}

	// negative immediate values should be allowed
	_, err = AssembleLine(labels, "ADDI x5, -20", 0)
	if err != nil {
		t.Errorf("A 'ADDI x5, -20' didn't assemble because %v", err)
	}

	_, err = AssembleLine(labels, "ADDI x5, 52", 0)
	if err == nil {
		t.Errorf("A 'ADDI x5, 52' should fail to assemble as immediate value is outside range -50 to 49")
	}

}
