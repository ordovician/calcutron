package prog

import (
	"bytes"
	"fmt"
	"io"
)

// Instructions of the type INST Rd, k encoded as ?dkk
type LongImmInstruction struct {
	BaseInstruction
}

func (inst *LongImmInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n == 1 {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
	} else {
		inst.err = fmt.Errorf("instruction expects 1 register operand but you gave %d", n)
	}
}

func (inst *LongImmInstruction) DecodeOperands(operands uint) {
	addr := operands % 100

	inst.regIndicies[Rd] = uint(operands / 100)
	inst.constant = Signed(addr, 100)
}

func (inst *LongImmInstruction) MachineCode() uint {
	regs := inst.regIndicies
	machineOpcode := uint(inst.opcode) * 1000

	constant := Complement(inst.constant, 100)
	destReg := uint(100 * regs[Rd])
	return machineOpcode + destReg + constant
}

func (inst *LongImmInstruction) printSourceCode(writer io.Writer) {
	printMnemonic(writer, inst.pseudoCode)
	fmt.Fprintf(writer, "x%d, ", inst.regIndicies[Rd])
	if inst.label == "" {
		NumberColor.Fprintf(writer, "%d", inst.constant)
	} else {
		LabelColor.Fprintf(writer, "%s", inst.label)
	}
}

func (inst *LongImmInstruction) SourceCode() string {
	var buffer bytes.Buffer
	inst.printSourceCode(&buffer)
	return buffer.String()
}

func (inst *LongImmInstruction) String() string {
	return inst.SourceCode()
}

// Instructions of the type INST Rd, Ra, k encoded as ?dak
type ShortImmInstruction struct {
	BaseInstruction
}

func (inst *ShortImmInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n == 2 {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
		inst.regIndicies[Ra] = inst.parsedRegIndicies[1]
	} else {
		inst.err = fmt.Errorf("two register operands were expected but you specified %d register operands", n)
	}
}

func (inst *ShortImmInstruction) printSourceCode(writer io.Writer) {
	// NOTE: A bit of a hack to treat infinite branches as HLT operation
	if (inst.opcode == BEQ || inst.opcode == BGT) && inst.constant == 0 {
		printMnemonic(writer, HLT)
		return
	}
	printMnemonic(writer, inst.opcode)
	printRegisterOperands(writer, inst.regIndicies[0:2])
	fmt.Fprintf(writer, ", ")
	NumberColor.Fprintf(writer, "%d", inst.constant)
}

// Will return colorized source code but this can be turned off with
// color.NoColor = true, or individual colors can be turned of such as with LabelColor.DisableColor() and LabelColor.EnableColor()
func (inst *ShortImmInstruction) SourceCode() string {
	var buffer bytes.Buffer
	inst.printSourceCode(&buffer)
	return buffer.String()
}

func (inst *ShortImmInstruction) String() string {
	return inst.SourceCode()
}

func (inst *ShortImmInstruction) DecodeOperands(operands uint) {
	addr := operands % 100

	inst.regIndicies[Rd] = uint(operands / 100)
	inst.regIndicies[Ra] = uint(addr / 10)
	inst.constant = Signed(addr%10, 10)
}

func (inst *ShortImmInstruction) MachineCode() uint {
	regs := inst.regIndicies
	operands := uint(100*regs[Rd] + 10*regs[Ra] + Complement(inst.constant, 10))
	machineOpcode := uint(inst.opcode) * 1000

	code := machineOpcode + operands
	return code
}

func (inst *ShortImmInstruction) ParseOperands(labels SymbolTable, operands []string, programCounter uint) {
	inst.BaseInstruction.ParseOperands(labels, operands, programCounter)
	if inst.err != nil {
		return
	}

	// Check if a label was used. In this case we need to calculate
	// a relative address from PC to that label. It is potentially too far
	if inst.label != "" {
		// in constant we want to store a relative address to label
		// hence we subtract from the label address the program counter (address of instruction)
		labelAddress := inst.constant
		inst.constant = labelAddress - int(programCounter)
	}

	if inst.constant < -5 || inst.constant > 4 {
		inst.err = fmt.Errorf("constant %d is outside valid range -5 to 4. Can happen if you try to jump to a label too far away with a branch instruction", inst.constant)
	}
}
