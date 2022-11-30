package prog

import (
	"fmt"
	"math"
)

type AddInstruction struct {
	BaseInstruction
}

func (inst *AddInstruction) Run(comp Machine) bool {
	value := inst.RegValue(comp, Ra) + inst.RegValue(comp, Rb)
	inst.SetRegValue(comp, Rd, value)
	return true
}

type AddImmediateInstruction struct {
	LongImmInstruction
}

func (inst *AddImmediateInstruction) Run(comp Machine) bool {
	value := inst.RegValue(comp, Rd)
	inst.SetRegValue(comp, Rd, value+inst.constant)
	return true
}

// Base implementation is setup for unsigned numbers such as JMP
func (inst *AddImmediateInstruction) ParseOperands(labels SymbolTable, operands []string, address uint) {
	inst.BaseInstruction.ParseOperands(labels, operands, address)
	if inst.err != nil {
		return
	}
	if inst.constant > 49 || inst.constant < -50 {
		inst.err = fmt.Errorf("constant %d is outside valid range -50 to 49", inst.constant)
	}
}

type SubInstruction struct {
	BaseInstruction
}

func (inst *SubInstruction) Run(comp Machine) bool {
	inst.SetRegValue(comp, Rd, inst.RegValue(comp, Ra)-inst.RegValue(comp, Rb))
	return true
}

type ShiftInstruction struct {
	ShortImmInstruction
}

func (inst *ShiftInstruction) ParseOperands(labels SymbolTable, operands []string, programCounter uint) {
	inst.ShortImmInstruction.ParseOperands(labels, operands, programCounter)
	if inst.err != nil {
		return
	}

	// if shift isn't given, use 1 as the shift
	if inst.constant == 0 {
		inst.constant = 1
	}
}

func (inst *ShiftInstruction) Run(comp Machine) bool {
	var value uint
	regValue := Complement(inst.RegValue(comp, Ra), 1e4)
	multiplier := uint(math.Pow10(abs(inst.constant)))
	if inst.constant >= 0 {
		value = regValue * multiplier
		inst.SetRegValue(comp, Rd, Signed(value/(1e4), 1e4))
	} else {
		value = regValue / multiplier
		inst.SetRegValue(comp, Rd, Signed(regValue%multiplier, 1e4))
	}
	inst.SetRegValue(comp, Ra, Signed(value%(1e4), 1e4))

	return true
}
