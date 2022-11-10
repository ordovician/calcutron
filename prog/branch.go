package prog

type BranchEqualInstruction struct {
	ShortImmInstruction
}

func (inst *BranchEqualInstruction) Run(comp Machine) bool {
	// since jumps are relative using 0 will get us stuck in endless loop
	// thus we treat endless loop as HLT or termination
	if inst.constant == 0 {
		return false
	}

	left := inst.RegValue(comp, Rd)
	right := inst.RegValue(comp, Ra)
	addr := int(comp.PC()) + inst.constant

	// IMPORTANT: Choosing to not deal with negative numbers here
	// that means any code dealing with actual negative numbers in a comparison will fail
	if Complement(left, 1e4) == Complement(right, 1e4) {
		comp.SetPC(uint(addr))
	}
	return true
}

type BranchGreaterThanInstruction struct {
	ShortImmInstruction
}

func (inst *BranchGreaterThanInstruction) Run(comp Machine) bool {
	// since jumps are relative using 0 will get us stuck in endless loop
	// thus we treat endless loop as HLT or termination
	if inst.constant == 0 {
		return false
	}

	left := inst.RegValue(comp, Rd)
	right := inst.RegValue(comp, Ra)
	addr := int(comp.PC()) + inst.constant

	// IMPORTANT: Choosing to not deal with negative numbers here
	// which wil cause problems in some cases.
	// however stuff like shifting numbers becomes really awkward if we handle numbers
	// as negative
	if Complement(left, 1e4) > Complement(right, 1e4) {
		comp.SetPC(uint(addr))
	}
	return true
}
