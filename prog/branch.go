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
