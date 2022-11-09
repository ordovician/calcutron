package prog

import "fmt"

type IncInstruction struct {
	AddImmediateInstruction
}

func (inst *IncInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
	inst.constant = 1
}

type DecInstruction struct {
	AddImmediateInstruction
}

func (inst *DecInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	regIndex := inst.parsedRegIndicies[0]
	inst.regIndicies[Rd] = regIndex
	inst.constant = -1
}

type SubImmediateInstruction struct {
	AddImmediateInstruction
}

func (inst *SubImmediateInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
	inst.constant = -inst.constant
}

type BranchLessThanInstruction struct {
	BranchGreaterThanInstruction
}

func (inst *BranchLessThanInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n < 3 {
		inst.err = fmt.Errorf("conditional branch instructions take 3 operands not %d", n)
	} else {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
		inst.regIndicies[Ra] = inst.parsedRegIndicies[2]
		inst.regIndicies[Rb] = inst.parsedRegIndicies[1]
	}
}

type UnconditionalBranchInstruction struct {
	BranchEqualInstruction
}

func (inst *UnconditionalBranchInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n > 1 {
		inst.err = fmt.Errorf("unconditional branch instructions take 1 operands not %d", n)
	} else {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
		inst.regIndicies[Ra] = 0
		inst.regIndicies[Rb] = 0
	}
}

type CopyInstruction struct {
	AddInstruction
}

func (inst *CopyInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n != 2 {
		inst.err = fmt.Errorf("MOVE instructions takes 2 register operands not %d", n)
	} else {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
		inst.regIndicies[Ra] = inst.parsedRegIndicies[1]
		inst.regIndicies[Rb] = 0
	}
}

type ClearInstruction struct {
	AddInstruction
}

func (inst *ClearInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n != 1 {
		inst.err = fmt.Errorf("CLEAR instructions takes 1 register operands not %d", n)
	} else {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
		inst.regIndicies[Ra] = 0
		inst.regIndicies[Rb] = 0
	}
}

type CallInstruction struct {
	RJumpInstruction
}

func (inst *CallInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n != 0 {
		inst.err = fmt.Errorf("CALL instructions takes 0 register operands not %d", n)
	} else {
		inst.regIndicies[Rd] = 9
		inst.regIndicies[Ra] = 0
		inst.regIndicies[Rb] = 0
	}
}

type ReturnInstruction struct {
	BranchEqualInstruction
}

func (inst *ReturnInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n != 0 {
		inst.err = fmt.Errorf("RET instructions takes 0 register operands not %d", n)
	} else {
		inst.regIndicies[Rd] = 9
		inst.regIndicies[Ra] = 0
		inst.regIndicies[Rb] = 0
	}
}

type NoOperationInstruction struct {
	AddInstruction
}

func (inst *NoOperationInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n != 0 {
		inst.err = fmt.Errorf("NOP instructions takes 0 register operands not %d", n)
	} else {
		inst.regIndicies[Rd] = 0
		inst.regIndicies[Ra] = 0
		inst.regIndicies[Rb] = 0
	}
}

type HaltInstruction struct {
	RJumpInstruction
}

func (inst *HaltInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n != 0 {
		inst.err = fmt.Errorf("HALT instructions takes 0 register operands not %d", n)
	} else {
		inst.regIndicies[Rd] = 0
		inst.regIndicies[Ra] = 0
		inst.regIndicies[Rb] = 0
	}
}

func (inst *HaltInstruction) ParseOperands(labels SymbolTable, operands []string, offset uint) {
	if len(operands) != 0 {
		inst.err = fmt.Errorf("HALT should not have any operands")
	}
}
