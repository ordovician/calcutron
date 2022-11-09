package prog

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type DataInstruction struct {
	BaseInstruction
}

func (inst *DataInstruction) ParseOperands(labels SymbolTable, operands []string, address uint) {
	if len(operands) != 1 {
		inst.err = fmt.Errorf("DAT directives can only have a single data point, not %d", len(operands))
	}

	data := operands[0]
	if constant, err := strconv.Atoi(data); err == nil {
		if constant < -9999 || constant > 9999 {
			inst.err = fmt.Errorf("DAT value %d is outside valid range -9999 to 9999", constant)
			return
		}
		inst.constant = constant
	}
}

func (inst *DataInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n != 0 {
		inst.err = fmt.Errorf("DAT directive is not an instruction. Youo cannot specify register operands")
	}
}

func (inst *DataInstruction) Run(comp Machine) bool {
	inst.err = fmt.Errorf("you cannot run a DAT directive. It is not an instruction")
	return false
}

func (inst *DataInstruction) MachineCode() uint {
	return Complement(inst.constant, 1e4)
}

func (inst *DataInstruction) printSourceCode(writer io.Writer) {
	printMnemonic(writer, inst.opcode)
	NumberColor.Fprintf(writer, "%04d", inst.constant)
}

func (inst *DataInstruction) SourceCode() string {
	var buffer bytes.Buffer
	inst.printSourceCode(&buffer)
	return buffer.String()
}
