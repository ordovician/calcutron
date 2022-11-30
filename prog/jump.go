package prog

import "fmt"

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
	// calculate destination address
	addr := uint(inst.RegValue(comp, Rd) + inst.constant)

	// Set return address
	inst.SetRegValue(comp, Rd, int(comp.PC()+1))

	comp.SetPC(addr)
	return true
}

func (inst *JumpInstruction) ParseOperands(labels SymbolTable, operands []string, address uint) {
	inst.BaseInstruction.ParseOperands(labels, operands, address)
	if inst.err != nil {
		return
	}
	if inst.constant > 99 || inst.constant < 0 {
		inst.err = fmt.Errorf("constant %d is outside valid range 0 to 99", inst.constant)
	}
}
