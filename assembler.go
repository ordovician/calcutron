package calcutron

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// Based on parsing mnemonic and operands of an instruction
type ParsedInstruction struct {
	opcode   Opcode
	regs     []uint8
	constant int8
}

func (inst *ParsedInstruction) ParseOperands(labels map[string]uint8, operands []string) error {
	registers := make([]uint8, 0)

	for _, operand := range operands {
		operand = strings.TrimSpace(operand)
		if addr, ok := labels[operand]; ok {
			inst.constant = int8(addr)
		} else if AllDigits(operand) {
			constant, err := strconv.Atoi(operand)
			if err != nil {
				return fmt.Errorf("unable to parse constant %s because %w", operand, err)
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
func (inst *ParsedInstruction) DestReg() uint16 {
	if len(inst.regs) == 0 {
		return 0
	}
	return 100 * uint16(inst.regs[0])
}

// Get the Rs1 source registers machine code
// It will be in range 0 - 90
func (inst *ParsedInstruction) FirstSourceReg() uint16 {
	regs := inst.regs
	var machinecode uint16

	switch inst.opcode {
	case ADD, SUB:
		switch len(regs) {
		case 3:
			machinecode = uint16(regs[1])
		case 2:
			machinecode = uint16(regs[0])
		}
	case SUBI, LSH, RSH:
		switch len(regs) {
		case 2:
			machinecode = uint16(regs[1])
		case 1:
			machinecode = uint16(regs[0])
		}
	case DEC, INC:
		machinecode = uint16(regs[0])
	case MOV:
		machinecode = 0
	default:
		log.Panicf("mnemonic %v has no rs1 operand", inst.opcode)
	}

	return 10 * machinecode
}

// Get the Rs2 source registers machine code
// It will be in range 0 - 9
func (inst *ParsedInstruction) SecondSourceReg() uint16 {
	regs := inst.regs
	switch inst.opcode {
	case ADD, SUB:
		switch len(regs) {
		case 3:
			return uint16(regs[2])
		case 2:
			return uint16(regs[1])
		}
	case MOV:
		return uint16(regs[1])
	default:
		log.Panicf("mnemonic %v has no rs1 operand", inst.opcode)
	}

	return 0
}

func (inst *ParsedInstruction) MachineInstruction() (uint16, error) {
	opcode := inst.opcode
	constant := inst.constant

	// For instructions without destination register a zero will be returned
	// Some instructions aren't real instructions such as DAT, and will thus return zero
	var machinecode uint16 = opcode.Machinecode() + inst.DestReg()

	switch opcode {
	case ADD, SUB:
		machinecode += inst.FirstSourceReg() + inst.SecondSourceReg()
	case SUBI, LSH, RSH:
		if Abs(constant) > 9 {
			return 0, fmt.Errorf("constant %d is too many digits. Instruction %v only allows single digit constants", constant, opcode)
		}
		machinecode += inst.FirstSourceReg() + complement(constant, 10)
	case LD, ST, BRZ, BGT:
		if constant < 0 {
			return 0, fmt.Errorf("cannot use negative address %d with %v instruction", constant, opcode)
		}
		machinecode += uint16(constant)
	case DAT:
		machinecode += complement(constant, 100)
	case DEC:
		machinecode += inst.FirstSourceReg()
	case MOV:
		machinecode += inst.SecondSourceReg()
	case HLT, INP, OUT, CLR:
		break
	default:
		log.Panicf("%v is an unknown opcode", opcode)
	}

	if machinecode > 9999 || machinecode < 0 {
		log.Panicf("%d isn't a valid machine code instructin as it cannot be longer than 4 digits", machinecode)
	}
	return uint16(machinecode), nil
}

// A table containing the memory address of labels in the code
func readSymTable(reader io.Reader) map[string]uint8 {
	scanner := bufio.NewScanner(reader)
	labels := make(map[string]uint8)
	for address := 0; scanner.Scan(); address++ {
		line := strings.Trim(scanner.Text(), " \t")
		n := len(line)

		if n == 0 {
			continue
		}

		if i := strings.IndexRune(line, ':'); i >= 0 {
			labels[line[0:i]] = uint8(address)

			// is there anything beyond the label?
			if n == i+1 {
				continue
			}
		}
	}
	return labels
}

// Assemble single line of code
func AssembleLine(labels map[string]uint8, line string) (int16, error) {
	code := strings.Trim(line, " \t")
	i := len(code)
	if j := strings.Index(code, "//"); j >= 0 {
		i = j
	}
	n := len(code)

	if n == 0 || code[n-1] == ':' {
		return -1, nil
	}

	code = code[0:i]
	if i = strings.IndexRune(code, ' '); i < 0 {
		i = n
	}
	mnemonic := code[0:i]
	var operands []string = strings.SplitN(code[i:], ",", 3)
	opcode := ParseOpcode(mnemonic)
	instruction := ParsedInstruction{
		opcode: opcode,
	}

	err := instruction.ParseOperands(labels, operands)
	if err != nil {
		return 0, err
	}

	machincode, err := instruction.MachineInstruction()

	return int16(machincode), err
}

// Assemble data from input
func Assemble(reader io.ReadSeeker) error {
	labels := readSymTable(reader)

	reader.Seek(0, io.SeekStart)

	scanner := bufio.NewScanner(reader)

	for lineno := 1; scanner.Scan(); lineno++ {
		machinecode, err := AssembleLine(labels, scanner.Text())
		if err != nil {
			return err
		}

		if machinecode > 0 {
			fmt.Printf("%04d\n", machinecode)
		}
	}

	return nil
}

// Assemble file and write output to stdout
func AssembleFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return Assemble(file)
}
