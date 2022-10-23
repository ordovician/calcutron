package calcutron

import "fmt"

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

// Disassembles instruction
func (inst *Instruction) String() string {
	switch inst.opcode {
	case ADD, SUB:
		return fmt.Sprintf("%4v x%d, x%d, x%d", inst.opcode, inst.dst, inst.src, inst.offset)
	case SUBI, LSH, RSH:
		return fmt.Sprintf("%4v x%d, x%d, %d", inst.opcode, inst.dst, inst.src, inst.offset)
	case LD, ST, BRZ, BGT:
		return fmt.Sprintf("%4v x%d, %d", inst.opcode, inst.dst, inst.addr)
	case HLT:
		return fmt.Sprintf("%4v", inst.opcode)
	default:
		break
	}
	return ""
}
