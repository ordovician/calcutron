package prog

func NewInstruction(opcode Opcode) (inst Instruction) {

	switch opcode {
	case ADD:
		inst = &AddInstruction{}
	case ADDI:
		inst = &AddImmediateInstruction{}
	case SUB:
		inst = &SubInstruction{}
	case LSH:
		inst = &ShiftInstruction{}
	case LOAD:
		inst = &LoadInstruction{}
	case MOVE:
		inst = &MoveInstruction{}
	case STOR:
		inst = &StoreInstruction{}
	case BEQ:
		inst = &BranchEqualInstruction{}
	case BGT:
		inst = &BranchGreaterThanInstruction{}
	case JMP:
		inst = &JumpInstruction{}

	// pseudo instructions
	case DEC:
		inst = &DecInstruction{}
		inst.setOpcode(ADDI)
	case INC:
		inst = &IncInstruction{}
		inst.setOpcode(ADDI)
	case SUBI:
		inst = &SubImmediateInstruction{} // so we can just subtract offset to get negative number
		inst.setOpcode(ADDI)
	case RSH:
		inst = &RightShiftInstruction{} // so we can just subtract offset to get negative number
		inst.setOpcode(LSH)
	case BRA:
		inst = &UnconditionalBranchInstruction{}
		inst.setOpcode(BEQ)
	case BLT:
		inst = &BranchLessThanInstruction{}
		inst.setOpcode(BGT)
	case COPY:
		inst = &CopyInstruction{}
		inst.setOpcode(ADD)
	case CLR:
		inst = &ClearInstruction{}
		inst.setOpcode(ADD)
	case CALL:
		inst = &CallInstruction{}
		inst.setOpcode(JMP)
	case NOP:
		inst = &NoOperationInstruction{}
		inst.setOpcode(ADD)
	case HLT:
		inst = &HaltInstruction{}
		inst.setOpcode(BEQ)
	case INP:
		inst = &InputInstruction{}
		inst.setOpcode(LOAD)
	case OUT:
		inst = &OutputInstruction{}
		inst.setOpcode(STOR)

	// non-instruction
	case DAT:
		inst = &DataInstruction{}
	default:
		break
	}

	// pseudo instructions must get opcode set explicitly
	if opcode <= 9 {
		inst.setOpcode(opcode)
	}
	inst.setPseudoCode(opcode)

	return
}
