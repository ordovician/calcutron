package prog

type BranchEqualInstruction struct {
	ShortImmInstruction
}

func (inst *BranchEqualInstruction) Run(comp Machine) bool {
	left := inst.RegValue(comp, Rd)
	right := inst.RegValue(comp, Ra)
	addr := int(comp.PC()) + inst.constant

	if left == right {
		comp.SetPC(uint(addr))
	}
	return true
}

type BranchGreaterThanInstruction struct {
	ShortImmInstruction
}

func (inst *BranchGreaterThanInstruction) Run(comp Machine) bool {
	left := inst.RegValue(comp, Rd)
	right := inst.RegValue(comp, Ra)
	addr := int(comp.PC()) + inst.constant

	if left > right {
		comp.SetPC(uint(addr))
	}
	return true
}

type JUMPInstruction struct {
	LongImmInstruction
}

func (inst *JUMPInstruction) Run(comp Machine) bool {
	// Set return address
	inst.SetRegValue(comp, Rd, int(comp.PC()+1))

	// jumping back to same instruction will create an infinite loop
	// hence this is a terminating instruction
	if inst.constant == 0 {
		return false
	}

	addr := uint(inst.constant)

	comp.SetPC(addr)
	return true
}
