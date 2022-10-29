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

func DisassembleInstruction(machinecode uint16) *ParsedInstruction {
	var inst ParsedInstruction
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

func (inst *ParsedInstruction) ParseOperands(labels map[string]uint8, operands []string) error {
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
	case LD, ST, BRZ, BGT, BRA:
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

	if machinecode > 9999 {
		log.Panicf("%d isn't a valid machine code instruction as it cannot be longer than 4 digits", machinecode)
	}
	return uint16(machinecode), nil
}

// A table containing the memory address of labels in the code
func readSymTable(reader io.Reader) map[string]uint8 {
	scanner := bufio.NewScanner(reader)
	labels := make(map[string]uint8)
	address := 0
	for scanner.Scan() {
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
		address++
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
	var operands []string

	if len(code[i:]) > 0 {
		operands = strings.SplitN(code[i:], ",", 3)
	}

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

func Assemble(reader io.ReadSeeker, writer io.Writer) error {
	return AssembleWithOptions(reader, writer, AssemblyFlag(0))
}

// Assembler reads assembly code from reader and writes machine code to writer
func AssembleWithOptions(reader io.ReadSeeker, writer io.Writer, options AssemblyFlag) error {
	labels := readSymTable(reader)

	reader.Seek(0, io.SeekStart)

	scanner := bufio.NewScanner(reader)

	var line SourceCodeLine
	line.address = 0
	for line.lineno = 1; scanner.Scan(); line.lineno++ {

		line.sourcecode = strings.Trim(scanner.Text(), " \t")
		var err error
		line.machinecode, err = AssembleLine(labels, line.sourcecode)
		if err != nil {
			return err
		}

		if line.machinecode >= 0 {
			PrintInstruction(writer, line, options|MACHINE_CODE)
			line.address++
		}
	}

	return nil
}

func AssembleFile(filepath string, writer io.Writer) error {
	return AssembleFileWithOptions(filepath, writer, MACHINE_CODE)
}

// AssembleFile reads assembly code from file at path filepath and write machinecode to writer
func AssembleFileWithOptions(filepath string, writer io.Writer, options AssemblyFlag) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return AssembleWithOptions(file, writer, options)
}
