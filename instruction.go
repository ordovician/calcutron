package calcutron

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Based on parsing mnemonic and operands of an instruction
type Instruction struct {
	opcode   Opcode
	regs     []uint8
	constant int8
	label    string // redundant, but helps with source code visualization
}

func DisassembleInstruction(machinecode uint16) *Instruction {
	var inst Instruction
	opcode := Opcode(machinecode / 1000)

	operands := machinecode % 1000
	addr := int8(operands % 100)

	regs := []uint8{uint8(operands / 100)}

	switch opcode {
	case LD, ST, BRZ, BGT, BRA:
		inst.constant = addr
	case SUBI, LSH, RSH:
		inst.constant = addr % 10
		regs = append(regs, uint8(addr/10))
	case ADD, SUB:
		regs = append(regs, uint8(addr/10))
		regs = append(regs, uint8(addr%10))
	}

	inst.opcode = opcode
	inst.regs = regs

	return &inst
}

func (inst *Instruction) ParseOperands(labels map[string]uint8, operands []string) error {
	registers := make([]uint8, 0)

	for _, operand := range operands {
		operand = strings.TrimSpace(operand)
		if addr, ok := labels[operand]; ok {
			inst.constant = int8(addr)
		} else if constant, err := strconv.Atoi(operand); err == nil {
			if constant < -99 || constant > 99 {
				return fmt.Errorf("constant %d is outside valid range -99 to 99", constant)
			}
			inst.constant = int8(constant)
		} else if len(operand) > 0 && operand[0] == 'x' {
			i, err := strconv.Atoi(operand[1:])
			if err != nil {
				return fmt.Errorf("unable to parse index %s because %w", operand[1:], err)
			}
			if i < 0 || i > 9 {
				return fmt.Errorf("x0 to x9 are the only valid registers, not x%d", i)
			}
			registers = append(registers, uint8(i))
		}
	}
	inst.regs = registers

	return nil
}

// Get the destination register machine code
// if there is no destination register we'll return 0 as that will have
// no affect on how machine code instruction is made
// It will be in range 0 - 900
func (inst *Instruction) DestRegCode() uint16 {
	return 100 * uint16(inst.DestReg())
}

func (inst *Instruction) DestReg() uint8 {
	if len(inst.regs) == 0 {
		return 0
	}
	return inst.regs[0]
}

// Get the Rs1 source registers machine code
// It will be in range 0 - 90
func (inst *Instruction) FirstSourceRegCode() uint16 {
	return uint16(10 * inst.FirstSourceReg())
}

func (inst *Instruction) FirstSourceReg() uint8 {
	regs := inst.regs
	var machinecode uint8

	switch inst.opcode {
	case ADD, SUB:
		switch len(regs) {
		case 3:
			machinecode = regs[1]
		case 2:
			machinecode = regs[0]
		}
	case SUBI, LSH, RSH:
		switch len(regs) {
		case 2:
			machinecode = regs[1]
		case 1:
			machinecode = regs[0]
		}
	case DEC, INC:
		machinecode = regs[0]
	case MOV:
		machinecode = 0
	default:
		log.Panicf("mnemonic %v has no rs1 operand", inst.opcode)
	}
	return machinecode
}

// Get the Rs2 source registers machine code
// It will be in range 0 - 9
func (inst *Instruction) SecondSourceRegCode() uint16 {
	return uint16(inst.SecondSourceReg())
}

func (inst *Instruction) SecondSourceReg() uint8 {
	regs := inst.regs
	switch inst.opcode {
	case ADD, SUB:
		switch len(regs) {
		case 3:
			return regs[2]
		case 2:
			return regs[1]
		}
	case MOV:
		return regs[1]
	default:
		log.Panicf("mnemonic %v has no rs1 operand", inst.opcode)
	}

	return 0
}

func (inst *Instruction) Constant() uint8 {
	return uint8(inst.constant)
}

func (inst *Instruction) Machinecode() (uint16, error) {
	opcode := inst.opcode
	constant := inst.constant

	// For instructions without destination register a zero will be returned
	// Some instructions aren't real instructions such as DAT, and will thus return zero
	var machinecode uint16 = opcode.Machinecode() + inst.DestRegCode()

	switch opcode {
	case ADD, SUB:
		machinecode += inst.FirstSourceRegCode() + inst.SecondSourceRegCode()
	case SUBI, LSH, RSH:
		if Abs(constant) > 9 {
			return 0, fmt.Errorf("constant %d is too many digits. Instruction %v only allows single digit constants", constant, opcode)
		}
		machinecode += inst.FirstSourceRegCode() + complement(constant, 10)
	case LD, ST, BRZ, BGT, BRA:
		if constant < 0 {
			return 0, fmt.Errorf("cannot use negative address %d with %v instruction", constant, opcode)
		}
		machinecode += uint16(constant)
	case DAT:
		machinecode += complement(constant, 100)
	case DEC:
		machinecode += inst.FirstSourceRegCode()
	case MOV:
		machinecode += inst.SecondSourceRegCode()
	case HLT, INP, OUT, CLR:
		break
	default:
		log.Panicf("%v is an unknown opcode", opcode)
	}

	if machinecode > 9999 {
		log.Panicf("%d isn't a valid machine code instruction as it cannot be longer than 4 digits", machinecode)
	}
	return uint16(machinecode), nil
}

// Used to print something that look like the original source code which
// was parsed to produce the Instruction object
func (inst *Instruction) PrintSourceCode(writer io.Writer) {
	opcode := inst.opcode
	constant := inst.Constant()

	fmt.Fprintf(writer, "%-4v", opcode)

	for i, r := range inst.regs {
		if i > 0 {
			fmt.Fprintf(writer, ", ")
		}
		fmt.Fprintf(writer, "x%v", r)
	}

	switch opcode {
	case SUBI, LSH, RSH:
		fmt.Fprintf(writer, ", %d", constant)
	case LD, ST, BRZ, BGT, BRA:
		if inst.label == "" {
			fmt.Fprintf(writer, ", %d", constant)
		} else {
			fmt.Fprintf(writer, ", %s", inst.label)
		}
	default:
		break
	}
}

func (inst *Instruction) PrintColoredSourceCode(writer io.Writer) {
	cyan := color.New(color.FgCyan, color.Bold)
	pink := color.New(color.FgHiRed)

	opcode := inst.opcode
	constant := inst.Constant()

	cyan.Fprintf(writer, "%-4v", opcode)

	for i, r := range inst.regs {
		if i > 0 {
			fmt.Fprintf(writer, ", ")
		}
		fmt.Fprintf(writer, "x%v", r)
	}

	switch opcode {
	case SUBI, LSH, RSH:
		fmt.Fprintf(writer, ", ")
		pink.Fprintf(writer, "%d", constant)
	case LD, ST, BRZ, BGT, BRA:
		if inst.label == "" {
			fmt.Fprintf(writer, ", ")
			pink.Fprintf(writer, "%d", constant)
		} else {
			fmt.Fprintf(writer, ", %s", inst.label)
		}
	default:
		break
	}
}

// What a disassembled instruction looks like in text form without colors
func (inst *Instruction) String() string {
	opcode := inst.opcode
	dst := inst.DestReg()
	constant := inst.Constant()

	switch opcode {
	case ADD, SUB:
		return fmt.Sprintf("%-4v x%d, x%d, x%d", opcode, dst, inst.FirstSourceReg(), inst.SecondSourceReg())
	case SUBI, LSH, RSH:
		return fmt.Sprintf("%-4v x%d, x%d, %d", opcode, dst, inst.FirstSourceReg(), constant)
	case LD, ST, BRZ, BGT:
		return fmt.Sprintf("%-4v x%d, %d", opcode, dst, constant)
	case HLT:
		return fmt.Sprintf("%-4v", opcode)
	default:
		break
	}
	return ""
}
