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
	case JUMP:
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
	case BRA:
		inst = &UnconditionalBranchInstruction{}
		inst.setOpcode(BEQ)
	case BLT:
		inst = &BranchLessThanInstruction{}
		inst.setOpcode(BGT)
	case COPY:
		inst = &CopyInstruction{}
		inst.setOpcode(ADD)
	case CLEAR:
		inst = &ClearInstruction{}
		inst.setOpcode(ADD)
	case CALL:
		inst = &CallInstruction{}
		inst.setOpcode(JUMP)
	case NOP:
		inst = &NoOperationInstruction{}
		inst.setOpcode(ADD)
	case HALT:
		inst = &HaltInstruction{}
		inst.setOpcode(BEQ)
	case IN:
		inst = &InputInstruction{}
		inst.setOpcode(LOAD)
	case OUT:
		inst = &OutputInstruction{}
		inst.setOpcode(STORE)

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
