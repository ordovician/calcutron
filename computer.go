package calcutron

import (
	"errors"
	"fmt"
)

// "bufio"

var ErrAllInputRead = errors.New("all inputs read")
var ErrProgramHalt = errors.New("reach halt instruction")

type Computer struct {
	PC        uint8      // Program counter 0-99
	Registers [10]uint8  // CPU registers   0-99
	Memory    [99]uint16 // Computer memory 0-9999
	Inputs    []uint8    // Input data to computer 0-99
	Outputs   []uint8    // Output from computer    0-99
	inpos     int        // Current position input stream
	// Stdin  io.Reader
	// Stdout io.Writer
	// Stderr io.Writer
}

func (comp *Computer) Step() error {
	pc := comp.PC
	ir := comp.Memory[pc]
	err := comp.ExecuteInstruction(ir)
	if err != nil {
		return fmt.Errorf("unable to execute instruction at address %d because: %w", pc, err)
	}

	// Make sure we didn't execute a branch instruction before updating Program counter
	if pc == comp.PC {
		comp.PC += 1
	}

	return nil
}

func (comp *Computer) RunSteps(nsteps int) error {
	for i := 0; i < nsteps; i++ {
		comp.PrintCurrentInstruction()
		err := comp.Step()
		if err != nil {
			return fmt.Errorf("could not run program because %w", err)
		}
	}
	return nil
}

func (comp *Computer) ExecuteInstruction(instruction uint16) error {
	if instruction > 9999 {
		return fmt.Errorf("instruction %d not within valid range 0000 - 9999", instruction)
	}

	regs := comp.Registers

	inst := decodeInstruction(instruction)

	var rd uint8

	if inst.dst >= 1 && inst.dst <= 9 {
		rd = regs[inst.dst]
	}

	switch inst.opcode {
	case ADD:
		rd = regs[inst.src] + regs[inst.offset]
	case SUB:
		rd = regs[inst.src] - regs[inst.offset]
	case SUBI:
		rd = regs[inst.src] - inst.offset
	case LSH:
		rd = regs[inst.src]*10 ^ inst.offset
	case RSH:
		rd = regs[inst.src] % (10 ^ inst.offset)
		regs[inst.src] = regs[inst.src] / (10 ^ inst.offset)
	case BRZ:
		if rd == 0 {
			comp.PC = inst.addr
		}
	case BGT:
		if rd > 0 {
			comp.PC = inst.addr
		}
	case LD:
		if inst.addr < 90 {
			rd = uint8(comp.Memory[inst.addr+1])
		} else if inst.addr == 90 {
			if comp.inpos >= len(comp.Inputs) {
				return ErrAllInputRead
			}
			rd = comp.Inputs[comp.inpos]
			comp.inpos += 1
		}
	case ST:
		if inst.addr < 90 {
			comp.Memory[inst.addr+1] = uint16(rd)
		} else if inst.addr == 91 {
			comp.Outputs = append(comp.Outputs, rd)
		} else {
			return fmt.Errorf("writing to address %d is not supported in this version", inst.addr)
		}
	case HLT:
		return ErrProgramHalt
	default:
		return fmt.Errorf("opcode %d, is not supported. Must be between 0-9", inst.opcode)
	}

	// Make sure register values stay within range 0-99
	// Act as if a register can hold two digits
	if rd != 0 {
		rd = 100 % rd
	}

	if inst.dst >= 1 && inst.dst <= 9 {
		regs[inst.dst] = rd
	}

	return nil
}

// Print instruction at program counter disassembled with address and machine code
// Example:
//
//	02: 1432; ADD  x4, x3, x2
func (comp *Computer) PrintCurrentInstruction() {
	machinecode := comp.Memory[comp.PC]
	instruction := decodeInstruction(machinecode)
	fmt.Printf("%02d: %04d; %v\n", comp.PC, machinecode, instruction)
}
