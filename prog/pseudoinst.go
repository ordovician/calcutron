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
	if n < 2 {
		inst.err = fmt.Errorf("conditional branch instructions take 2 register operands not %d", n)
	} else {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[1]
		inst.regIndicies[Ra] = inst.parsedRegIndicies[0]
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
	if n > 0 {
		inst.err = fmt.Errorf("unconditional branch instructions take 0 register operands not %d", n)
	} else {
		inst.regIndicies[Rd] = 0
		inst.regIndicies[Ra] = 0
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
		inst.err = fmt.Errorf("move instructions takes 2 register operands not %d", n)
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
		inst.err = fmt.Errorf("the clear instruction takes 1 register operands not %d", n)
	} else {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
		inst.regIndicies[Ra] = 0
		inst.regIndicies[Rb] = 0
	}
}

type CallInstruction struct {
	JumpInstruction
}

func (inst *CallInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n != 0 {
		inst.err = fmt.Errorf("the call instruction takes 0 register operands not %d", n)
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
		inst.err = fmt.Errorf("A NOP instructions takes 0 register operands not %d", n)
	} else {
		inst.regIndicies[Rd] = 0
		inst.regIndicies[Ra] = 0
		inst.regIndicies[Rb] = 0
	}
}

type HaltInstruction struct {
	UnconditionalBranchInstruction
}

func (inst *HaltInstruction) ParseOperands(labels SymbolTable, operands []string, offset uint) {
	if len(operands) != 0 {
		inst.err = fmt.Errorf("HLT should not have any operands")
	}
}

type InputInstruction struct {
	LoadInstruction
}

func (inst *InputInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n != 1 {
		inst.err = fmt.Errorf("the input instruction takes 1 register operand not %d", n)
	} else {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
		inst.regIndicies[Ra] = 0
		inst.constant = -1
	}
}

type OutputInstruction struct {
	StoreInstruction
}

func (inst *OutputInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n != 1 {
		inst.err = fmt.Errorf("the out instruction takes 1 register operand not %d", n)
	} else {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
		inst.regIndicies[Ra] = 0
		inst.constant = -1
	}
}
