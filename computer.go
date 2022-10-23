package calcutron

import (
	"errors"
	"fmt"
)

// "bufio"

var AllInputRead = errors.New("All inputs read")

type Instruction struct {
	opcode Opcode
	dst    uint8
	addr   uint8
	src    uint8
	offset uint8
}

// Split a number from 0-9999 into its individual parts
func decodeInstruction(instruction uint16) Instruction {
	operands := instruction % 1000
	addr := uint8(operands % 100)
	return Instruction{
		opcode: Opcode(instruction / 1000),
		dst:    uint8(operands / 100),
		addr:   addr,
		src:    uint8(addr / 10),
		offset: addr % 10,
	}
}

type Computer struct {
	PC        uint8      // Program counter 0-99
	Registers [10]uint8  // CPU registers   0-99
	Memory    [99]uint16 // Computer memory 0-9999
	Inputs    []uint8    // Input data to computer 0-99
	Outputs   []uint8    // Ouput from computer    0-99
	inpos     int        // Current position input stream
	// Stdin  io.Reader
	// Stdout io.Writer
	// Stderr io.Writer
}

func (comp *Computer) Step() error {

	return nil
}

func (comp *Computer) ExecuteInstruction(instruction uint16) error {
	if instruction > 9999 {
		return fmt.Errorf("instruction %d not within valid range 0000 - 9999", instruction)
	}

	regs := comp.Registers

	opcode := Opcode(instruction / 1000)
	operands := instruction % 1000
	dst := uint8(operands / 100)
	addr := uint8(operands % 100)
	src := uint8(addr / 10)
	offset := addr % 10
	var rd uint8

	if dst >= 1 && dst <= 9 {
		rd = regs[dst]
	}

	switch opcode {
	case ADD:
		rd = 100 % (regs[src] + regs[offset])
	case SUB:
		rd = 100 % (regs[src] - regs[offset])
	case SUBI:
		rd = 100 % (regs[src] - offset)
	case LSH:
		rd = 100 % (regs[src]*10 ^ offset)
	case RSH:
		rd = regs[src] % (10 ^ offset)
		regs[src] = regs[src] / (10 ^ offset)
	case BRZ:
		if rd == 0 {
			comp.PC = addr
		}
	case BGT:
		if rd > 0 {
			comp.PC = addr
		}
	case LD:
		if addr < 90 {
			rd = uint8(comp.Memory[addr+1])
		} else if addr == 90 {
			if comp.inpos >= len(comp.Inputs) {
				return AllInputRead
			}
		}
	case ST:
		if addr < 90 {
			comp.Memory[addr+1] = uint16(rd)
		} else if addr == 91 {
			comp.Outputs = append(comp.Outputs, rd)
		} else {
			return fmt.Errorf("Writing to address %d is not supported in this version", addr)
		}
	case HLT:
		return nil
	default:
		return fmt.Errorf("Opcode %d, is not supported. Must be between 0-9", opcode)
	}

	if dst >= 1 && dst <= 9 {
		regs[dst] = rd
	}

	return nil
}
