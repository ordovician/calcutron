package prog

import (
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

func (inst *ShiftInstruction) Run(comp Machine) bool {
	var value int
	regValue := inst.RegValue(comp, Ra)
	multiplier := int(math.Pow10(abs(inst.constant)))
	if inst.constant >= 0 {
		value = regValue * multiplier
		inst.SetRegValue(comp, Rd, value/(1e4))
	} else {
		value = regValue / multiplier
		inst.SetRegValue(comp, Rd, regValue%multiplier)
	}
	inst.SetRegValue(comp, Ra, value%(1e4))

	return true
}
