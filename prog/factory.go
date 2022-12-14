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
	case LODI:
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
	case MOVE:
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
	// cannot make this as an instruction as will produce many entries
	case STR:
		panic("STR is not an instruction. Bug in software. Circumvent by not using STR directive.")
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
