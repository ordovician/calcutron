package calcutron

import "strings"

// NOTE: If you are missing the stringer command, you can install it
// $ go install golang.org/x/tools/cmd/stringer@latest

//go:generate stringer -type=Opcode opcode.go

// Actual valid opcodes are numbered from 0 to 9, while the pseudo code instuctions
// map to one of the valid core instructions
type Opcode uint8

const (
	HLT  Opcode = iota // HaLT execution
	ADD                // ADD registers
	SUB                // SUBtract registers
	SUBI               // SUBtract Immediate
	LSH                // Left SHift
	RSH                // Right SHift
	BRZ                // BRanch if Zero
	BGT                // Branch if Greater Than
	LD                 // LoaD
	ST                 // STore

	// Pseudo instructions
	INP  // INPut
	OUT  // OUTput
	DEC  // DECrement
	INC  // INCrement
	ADDI // ADD Immediate
	BRA  // BRAnch
	CLR  // CLearR
	MOV  // MOVe from one reg to another
)

var opcodes = [...]Opcode{HLT, ADD, SUB, SUBI, LSH, RSH, BRZ, BGT, LD, ST, INP, OUT, DEC, INC, ADDI, BRA, CLR, MOV}

// Turns text string into Opcode
func ParseOpcode(s string) Opcode {
	s = strings.ToUpper(s)

	// inefficient to loop but the list is of limited size so it should
	// be acceptable
	for _, opcode := range opcodes {
		if opcode.String() == s {
			return opcode
		}
	}
	return HLT
}
