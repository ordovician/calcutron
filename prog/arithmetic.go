package prog

import (
	"bytes"
	"fmt"
	"io"
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
	ImmediateInstruction
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
	BaseInstruction
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

func (inst *ShiftInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
	inst.regIndicies[Ra] = inst.parsedRegIndicies[1]
}

func (inst *ShiftInstruction) printSourceCode(writer io.Writer) {
	printMnemonic(writer, inst.opcode)
	printRegisterOperands(writer, inst.regIndicies[0:2])
	fmt.Fprintf(writer, ", ")
	NumberColor.Fprintf(writer, "%d", inst.constant)
}

// Will return colorized source code but this can be turned off with
// color.NoColor = true, or individual colors can be turned of such as with LabelColor.DisableColor() and LabelColor.EnableColor()
func (inst *ShiftInstruction) SourceCode() string {
	var buffer bytes.Buffer
	inst.printSourceCode(&buffer)
	return buffer.String()
}

func (inst *ShiftInstruction) DecodeOperands(operands uint) {
	addr := operands % 100

	inst.regIndicies[Rd] = uint(operands / 100)
	inst.regIndicies[Ra] = uint(addr / 10)
	inst.constant = Signed(addr%10, 10)
}

func (inst *ShiftInstruction) MachineCode() uint {
	regs := inst.regIndicies
	operands := uint(100*regs[Rd] + 10*regs[Ra] + Complement(inst.constant, 10))
	code := inst.opcode.MachineCode() + operands
	return code
}
