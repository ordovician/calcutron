package calcutron

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

	// not really instruction
	DAT
)

var AllOpcodes = [...]Opcode{HLT, ADD, SUB, SUBI, LSH, RSH, BRZ, BGT, LD, ST, INP, OUT, DEC, INC, ADDI, BRA, CLR, MOV, DAT}
var AllOpcodeStrings []string = make([]string, len(AllOpcodes))

// initialize opcode strings
func init() {
	for i, opcode := range AllOpcodes {
		AllOpcodeStrings[i] = opcode.String()
	}
}

// Turns text string into Opcode
func ParseOpcode(s string) Opcode {
	s = strings.ToUpper(s)

	// inefficient to loop but the list is of limited size so it should
	// be acceptable
	for _, opcode := range AllOpcodes {
		if opcode.String() == s {
			return opcode
		}
	}
	return HLT
}

// Create machine code representation of an instruction with given opcode
// This method also deals with pseudo opcodes
func (opcode Opcode) Machinecode() uint16 {
	switch opcode {
	case HLT:
		return 0
	case ADD:
		return 1000
	case SUB:
		return 2000
	case SUBI:
		return 3000
	case LSH:
		return 4000
	case RSH:
		return 5000
	case BRZ:
		return 6000
	case BGT:
		return 7000
	case LD:
		return 8000
	case ST:
		return 9000

	// pseudo instructions
	case INP:
		return 8090
	case OUT:
		return 9091
	case DEC:
		return 3001
	case INC:
		return 3009
	case ADDI:
		return 3010 // so we can just subtract offset to get negative number
	case BRA:
		return 6000
	case MOV:
		return 1000
	case CLR:
		return 1000

	// non-instuction
	case DAT:
		return 0
	default:
		break
	}
	log.Panicf("unknown opcode %v encountered", opcode)
	return 0
}
