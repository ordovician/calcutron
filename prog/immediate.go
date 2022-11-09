package prog

import (
	"bytes"
	"fmt"
	"io"
)

// Instructions of the type INST Rd, k encoded as ?dkk
type ImmediateInstruction struct {
	BaseInstruction
}

func (inst *ImmediateInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
}

func (inst *ImmediateInstruction) DecodeOperands(operands uint) {
	addr := operands % 100

	inst.regIndicies[Rd] = uint(operands / 100)
	inst.constant = Signed(addr, 100)
}

func (inst *ImmediateInstruction) MachineCode() uint {
	regs := inst.regIndicies
	opcodeEncoding := inst.opcode.MachineCode()
	constant := Complement(inst.constant, 100)
	destReg := uint(100 * regs[Rd])
	return opcodeEncoding + destReg + constant
}

func (inst *ImmediateInstruction) printSourceCode(writer io.Writer) {
	printMnemonic(writer, inst.opcode)
	fmt.Fprintf(writer, "x%d, ", inst.regIndicies[Rd])
	if inst.label == "" {
		NumberColor.Fprintf(writer, "%d", inst.constant)
	} else {
		LabelColor.Fprintf(writer, "%s", inst.label)
	}
}

func (inst *ImmediateInstruction) SourceCode() string {
	var buffer bytes.Buffer
	inst.printSourceCode(&buffer)
	return buffer.String()
}
