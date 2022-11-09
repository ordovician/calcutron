package prog

func NewInstruction(opcode Opcode) (inst Instruction) {

	switch opcode {
	case ADD:
		inst = &AddInstruction{}
	case ADDI:
		inst = &AddImmediateInstruction{}
	case SUB:
		inst = &SubInstruction{}
	case SHIFT:
		inst = &ShiftInstruction{}
	case LOAD:
		inst = &LoadInstruction{}
	case MOVE:
		inst = &MoveInstruction{}
	case STORE:
		inst = &StoreInstruction{}
	case BEQ:
		inst = &BranchEqualInstruction{}
	case BGT:
		inst = &BranchGreaterThanInstruction{}
	case RJUMP:
		inst = &RJumpInstruction{}

	// pseudo instructions
	case DEC:
		inst = &DecInstruction{}
	case INC:
		inst = &IncInstruction{}
	case SUBI:
		inst = &SubImmediateInstruction{} // so we can just subtract offset to get negative number
	case BRA:
		inst = &UnconditionalBranchInstruction{}
	case BLT:
		inst = &BranchLessThanInstruction{}
	case COPY:
		inst = &CopyInstruction{}
	case CLEAR:
		inst = &ClearInstruction{}
	case CALL:
		inst = &CallInstruction{}
	case RET:
		inst = &ReturnInstruction{}
	case NOP:
		inst = &NoOperationInstruction{}
	case HALT:
		inst = &HaltInstruction{}

	// non-instruction
	case DAT:
		inst = &DataInstruction{}
	default:
		break
	}
	inst.setOpcode(opcode)

	return
}
