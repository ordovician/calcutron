package prog

type LoadInstruction struct {
	BaseInstruction
}

func (inst *LoadInstruction) Run(comp Machine) bool {
	signedAddr := inst.RegValue(comp, Ra) + inst.RegValue(comp, Rb)
	addr := Complement(signedAddr, 10e5)

	if signedAddr == -1 {
		// check if we have exhausted input. Program should terminate in that case
		value, ok := comp.PopInput()
		if !ok {
			return false
		}
		inst.SetRegValue(comp, Rd, value)
	} else {
		value := Signed(comp.Memory(addr), 1e4)
		inst.SetRegValue(comp, Rd, value)
	}
	return true
}

type MoveInstruction struct {
	ImmediateInstruction
}

func (inst *MoveInstruction) Run(comp Machine) bool {
	inst.SetRegValue(comp, Rd, inst.constant)
	return true
}

type StoreInstruction struct {
	BaseInstruction
}

func (inst *StoreInstruction) Run(comp Machine) bool {
	address := inst.RegValue(comp, Ra) + inst.RegValue(comp, Rb)
	value := inst.RegValue(comp, Rd)

	if address == -1 {
		comp.PushOutput(value)
	} else {
		comp.SetMemory(Complement(address, 1e4), Complement(value, 1e4))
	}
	return true
}
