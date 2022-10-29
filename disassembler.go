package calcutron

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type MachineInstruction struct {
	opcode Opcode
	dst    uint8
	addr   uint8
	src    uint8
	offset uint8
}

// Split a number from 0-9999 into its individual parts
func decodeInstruction(instruction uint16) *MachineInstruction {
	operands := instruction % 1000
	addr := uint8(operands % 100)
	return &MachineInstruction{
		opcode: Opcode(instruction / 1000),
		dst:    uint8(operands / 100),
		addr:   addr,
		src:    uint8(addr / 10),
		offset: addr % 10,
	}
}

// Disassembles instruction
func (inst *MachineInstruction) String() string {
	switch inst.opcode {
	case ADD, SUB:
		return fmt.Sprintf("%-4v x%d, x%d, x%d", inst.opcode, inst.dst, inst.src, inst.offset)
	case SUBI, LSH, RSH:
		return fmt.Sprintf("%-4v x%d, x%d, %d", inst.opcode, inst.dst, inst.src, inst.offset)
	case LD, ST, BRZ, BGT:
		return fmt.Sprintf("%-4v x%d, %d", inst.opcode, inst.dst, inst.addr)
	case HLT:
		return fmt.Sprintf("%-4v", inst.opcode)
	default:
		break
	}
	return ""
}

func Disassemble(reader io.Reader) error {
	return DisassembleWithOptions(reader, SOURCE_CODE)
}

// Disassemble machinecode read from reader
func DisassembleWithOptions(reader io.Reader, option AssemblyFlag) error {
	scanner := bufio.NewScanner(reader)

	var line SourceCodeLine
	line.address = 0
	for line.lineno = 1; scanner.Scan(); line.lineno++ {
		line.sourcecode = scanner.Text()
		machinecode, err := strconv.Atoi(line.sourcecode)
		line.machinecode = int16(machinecode)
		if err != nil {
			return fmt.Errorf("%d: unable to disassemble because: %w", line.lineno, err)
		}
		if line.machinecode < 0 {
			log.Panicf("%d: something went from in parsing code. Machine code instruction should never be less than 0", line.lineno)
		}

		instruction := decodeInstruction(uint16(machinecode))
		fmt.Printf("%02d: %04d; %v\n", line.lineno-1, line.machinecode, instruction)
	}

	return nil
}

func DisassembleFile(filepath string) error {
	return DisassembleFileWithOptions(filepath, SOURCE_CODE)
}

// Disassemble file and write output to stdout
func DisassembleFileWithOptions(filepath string, options AssemblyFlag) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return DisassembleWithOptions(file, options)
}
