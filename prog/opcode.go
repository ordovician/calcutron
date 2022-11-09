package prog

import (
	"log"
	"strings"
)

// NOTE: If you are missing the stringer command, you can install it
// $ go install golang.org/x/tools/cmd/stringer@latest

//go:generate stringer -type=Opcode opcode.go

// Actual valid opcodes are numbered from 0 to 9, while the pseudo code instuctions
// map to one of the valid core instructions
type Opcode uint8

const (
	RJUMP Opcode = iota // HALT execution
	ADD                 // Add registers
	ADDI                // Add Immediate
	SUB                 // Subtract registers
	SHIFT               // Shift digits left, or right if k is negative
	LOAD                // Load
	MOVE                // Load Immediate
	STORE               // Store to memory
	BEQ                 // Branch if Equal
	BGT                 // Branch if Greater than

	// Pseudo instructions
	DEC   // DECrement
	INC   // INCrement
	SUBI  // SUBtract Immediate
	BRA   // BRAnch
	BLT   // Branch Less Than
	CLEAR // Clear
	COPY  // COPY from one reg to another
	CALL  // CALL subroutine
	RET   // Return from subroutine
	NOP   // No Operation
	HALT  // Halt execution

	// not really instruction
	DAT
)

var AllOpcodes = [...]Opcode{RJUMP, ADD, ADDI, SUB, SHIFT, LOAD, MOVE, STORE, BEQ, BGT,
	DEC, INC, SUBI, BRA, BLT, CLEAR, COPY, CALL, RET, NOP, HALT, DAT}
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
	return HALT, false
}

// Create machine code representation of an instruction with given opcode
// This method also deals with pseudo opcodes
func (opcode Opcode) MachineCode() uint {
	switch opcode {
	case RJUMP, HALT:
		return 0
	case ADD:
		return 1000
	case ADDI:
		return 2000
	case SUB:
		return 3000
	case SHIFT:
		return 4000
	case LOAD:
		return 5000
	case MOVE:
		return 6000
	case STORE:
		return 7000
	case BEQ:
		return 8000
	case BGT:
		return 9000

	// pseudo instructions
	case DEC:
		return 2000
	case INC:
		return 2000
	case SUBI:
		return 2100 // so we can just subtract offset to get negative number
	case BRA:
		return 8000
	case BLT:
		return 9000
	case COPY:
		return 1000
	case CLEAR:
		return 1000
	case CALL:
		return 8000

	// non-instuction
	case DAT:
		return 0
	default:
		break
	}
	log.Panicf("unknown opcode %v encountered", opcode)
	return 0
}
