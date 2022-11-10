package prog

import "fmt"

type LoadStoreInstruction struct {
	ShortImmInstruction
}

// We want register assignments for  load and store to work as follows:
// Rd, k -> Rd, x0, k
// which is different from how arithmetic operations work
func (inst *LoadStoreInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	switch n {
	case 1:
		inst.regIndicies[Ra] = 0
	case 2:
		inst.regIndicies[Ra] = inst.parsedRegIndicies[1]
	default:
		inst.err = fmt.Errorf("load and store instruction require 1 or 2 register operands not %d", n)
		return
	}
	inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
}

type LoadInstruction struct {
	LoadStoreInstruction
}

func (inst *LoadInstruction) Run(comp Machine) bool {
	signedAddr := inst.RegValue(comp, Ra) + inst.constant
	addr := Complement(signedAddr, 1e4)

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
	LongImmInstruction
}

func (inst *MoveInstruction) Run(comp Machine) bool {
	inst.SetRegValue(comp, Rd, inst.constant)
	return true
}

type StoreInstruction struct {
	LoadStoreInstruction
}

func (inst *StoreInstruction) Run(comp Machine) bool {
	address := inst.RegValue(comp, Ra) + inst.constant
	value := inst.RegValue(comp, Rd)

	if address == -1 {
		comp.PushOutput(value)
	} else {
		comp.SetMemory(Complement(address, 1e4), Complement(value, 1e4))
	}
	return true
}
