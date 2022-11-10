package prog

type JumpInstruction struct {
	LongImmInstruction
}

// Make it possible to turn JMP Rd into JMP Rd, 0 and JMP k into JMP x0, k
func (inst *JumpInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	if len(inst.parsedRegIndicies) == 0 {
		inst.regIndicies[Rd] = 0
	} else {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
	}
}

// A JMP Rd, k instruction computes destination addres with Rd + k and
// stores a return address in Rd. Rd = PC + 1
func (inst *JumpInstruction) Run(comp Machine) bool {
	// jumping back to same instruction will create an infinite loop
	// hence this is a terminating instruction
	if inst.constant == 0 {
		return false
	}

	// calculate destination address
	addr := uint(inst.RegValue(comp, Rd) + inst.constant)

	// Set return address
	inst.SetRegValue(comp, Rd, int(comp.PC()+1))

	comp.SetPC(addr)
	return true
}
