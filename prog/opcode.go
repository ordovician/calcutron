package prog

import (
	"strings"
)

// NOTE: If you are missing the stringer command, you can install it
// $ go install golang.org/x/tools/cmd/stringer@latest

//go:generate stringer -type=Opcode opcode.go

// Actual valid opcodes are numbered from 0 to 9, while the pseudo code instuctions
// map to one of the valid core instructions
type Opcode uint8

const (
	BEQ  Opcode = iota // Branch if Equal
	ADD                // Add registers
	ADDI               // Add Immediate
	SUB                // Subtract registers
	LSH                // Shift digits left, or right if k is negative
	LOAD               // Load
	MOVE               // Load Immediate
	STOR               // Store to memory
	JMP                // Jump to address
	BGT                // Branch if Greater than

	// Pseudo instructions
	DEC  // DECrement
	INC  // INCrement
	SUBI // SUBtract Immediate
	RSH  // Righ SHift
	BRA  // BRAnch
	BLT  // Branch Less Than
	CLR  // Clear
	COPY // COPY from one reg to another
	CALL // CALL subroutine
	NOP  // No Operation
	HLT  // Halt execution
	INP  // INput instruction
	OUT  // OUTput instruction

	// not really instruction
	DAT
)

var AllOpcodes = [...]Opcode{JMP, ADD, ADDI, SUB, LSH, LOAD, MOVE, STOR, BEQ, BGT,
	DEC, INC, SUBI, RSH, BRA, BLT, CLR, COPY, CALL, NOP, HLT, INP, OUT, DAT}
var AllOpcodeStrings []string = make([]string, len(AllOpcodes))

// initialize opcode strings
func init() {
	for i, opcode := range AllOpcodes {
		AllOpcodeStrings[i] = opcode.String()
	}
}

// Turns text string into Opcode
func ParseOpcode(s string) (Opcode, bool) {
	s = strings.ToUpper(s)

	// inefficient to loop but the list is of limited size so it should
	// be acceptable
	for _, opcode := range AllOpcodes {
		if opcode.String() == s {
			return opcode, true
		}
	}
	return HLT, false
}
