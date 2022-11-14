package prog

import (
	"bytes"
	"fmt"
	"io"
)

type LoadStoreInstruction struct {
	BaseInstruction
}

func (inst *LoadStoreInstruction) printSourceCode(writer io.Writer) {
	printMnemonic(writer, inst.opcode)
	printRegisterOperands(writer, inst.regIndicies[0:2])
	fmt.Fprintf(writer, ", ")
	NumberColor.Fprintf(writer, "%d", inst.constant)
}

// Will return colorized source code but this can be turned off with
// color.NoColor = true, or individual colors can be turned of such as with LabelColor.DisableColor() and LabelColor.EnableColor()
func (inst *LoadStoreInstruction) SourceCode() string {
	var buffer bytes.Buffer
	inst.printSourceCode(&buffer)
	return buffer.String()
}

func (inst *LoadStoreInstruction) DecodeOperands(operands uint) {
	addr := operands % 100

	inst.regIndicies[Rd] = uint(operands / 100)
	inst.regIndicies[Ra] = uint(addr / 10)
	inst.constant = int(addr % 10) // not using signed offsets for LOAD STORE instructions
	if inst.constant >= 8 {
		inst.constant = 10 - inst.constant
	}
}

func (inst *LoadStoreInstruction) MachineCode() uint {
	regs := inst.regIndicies
	constant := inst.constant
	if constant < 0 {
		constant = 10 + constant
	}
	operands := uint(100*regs[Rd] + 10*regs[Ra] + uint(constant))
	machineOpcode := uint(inst.opcode) * 1000

	code := machineOpcode + operands
	return code
}

func (inst *LoadStoreInstruction) ParseOperands(labels SymbolTable, operands []string, programCounter uint) {
	inst.BaseInstruction.ParseOperands(labels, operands, programCounter)
	if inst.err != nil {
		return
	}

	if inst.constant < -2 || inst.constant > 7 {
		inst.err = fmt.Errorf("offset %d is outside valid range -2 to 7. This can happen if label is too far away from address zero or base address", inst.constant)
	}
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
