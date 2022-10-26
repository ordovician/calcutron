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

// Disassemble machinecode read from reader
func Disassemble(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)

	for lineno := 1; scanner.Scan(); lineno++ {
		line := scanner.Text()
		machinecode, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("%d: unable to disassemble because: %w", lineno, err)
		}
		if machinecode < 0 {
			log.Panicf("%d: something went from in parsing code. Machine code instruction should never be less than 0", lineno)
		}

		instruction := decodeInstruction(uint16(machinecode))
		fmt.Printf("%02d: %04d; %v\n", lineno-1, machinecode, instruction)
	}

	return nil
}

// Disassemble file and write output to stdout
func DisassembleFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return Disassemble(file)
}
