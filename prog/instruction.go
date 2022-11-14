package prog

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/ordovician/calcutron/utils"
)

// Is it the 1st, 2nd or 3rd operand to an instruction
type OperandIndex uint

const (
	Rd OperandIndex = iota // Index of destination register
	Ra                     // First source register index
	Rb                     // Second source register index
)

// The interface a simulator of a computer has to adhere in order to be able to supply
// an instruction word with what it needs when it is executed
type Machine interface {
	Register(i uint) int
	SetRegister(i uint, value int)
	PC() uint
	SetPC(address uint)
	Memory(address uint) uint
	SetMemory(address uint, value uint)
	PopInput() (int, bool)
	PushOutput(value int)
}

type Instruction interface {
	setOpcode(opcode Opcode)
	setPseudoCode(pseudoCode Opcode)

	Opcode() Opcode
	UniqueRegisters() []uint
	Run(comp Machine) bool
	MachineCode() uint
	SourceCode() string

	ParseOperands(labels SymbolTable, operands []string, address uint)
	DecodeOperands(machinecode uint)
	AssignRegisters()
	Error() error
}

type BaseInstruction struct {
	opcode      Opcode
	pseudoCode  Opcode  // fake opcode outside the 0 to 9 range
	regIndicies [3]uint // machine code would set this directly
	constant    int     // signed constant. How to convert this depends on whether we deal with single of double digit constant
	label       string

	parsedRegIndicies []uint // set from parsed source code
	err               error  // sticky error
}

func (inst *BaseInstruction) setOpcode(opcode Opcode) {
	inst.opcode = opcode
}

func (inst *BaseInstruction) setPseudoCode(pseudoCode Opcode) {
	inst.pseudoCode = pseudoCode
}

func (inst *BaseInstruction) Opcode() Opcode {
	return inst.opcode
}

func (inst *BaseInstruction) UniqueRegisters() []uint {
	return utils.RemoveDuplicates(inst.regIndicies[:])
}

func (inst *BaseInstruction) Run(comp Machine) bool {
	return true
}

func (inst *BaseInstruction) RegValue(comp Machine, operand OperandIndex) int {
	return comp.Register(inst.regIndicies[operand])
}

func (inst *BaseInstruction) SetRegValue(comp Machine, operand OperandIndex, value int) {
	comp.SetRegister(inst.regIndicies[operand], value)
}

func (inst *BaseInstruction) MachineCode() uint {
	regs := inst.regIndicies
	operands := uint(100*regs[Rd] + 10*regs[Ra] + regs[Rb])
	machineOpcode := uint(inst.opcode) * 1000

	code := machineOpcode + operands
	return code
}

func (inst *BaseInstruction) printSourceCode(writer io.Writer) {
	printMnemonic(writer, inst.pseudoCode)

	regs := inst.regIndicies[:]
	if len(inst.parsedRegIndicies) > 0 {
		regs = inst.parsedRegIndicies
	}

	printRegisterOperands(writer, regs)
}

// Will return colorized source code but this can be turned off with
// color.NoColor = true, or individual colors can be turned of such as with LabelColor.DisableColor() and LabelColor.EnableColor()
func (inst *BaseInstruction) SourceCode() string {
	var buffer bytes.Buffer
	inst.printSourceCode(&buffer)
	return buffer.String()
}

func (inst *BaseInstruction) String() string {
	return inst.SourceCode()
}

// Fills in the parsedRegIndicies, constant and label fields
// operands should not contain uncessesary whitespace. Caller is responsible for cleaning up operands before calling
// ParseOperands
func (inst *BaseInstruction) ParseOperands(labels SymbolTable, operands []string, address uint) {
	registers := make([]uint, 0)

	for _, operand := range operands {
		if addr, ok := labels[operand]; ok {
			inst.constant = int(addr)
			inst.label = operand
		} else if constant, err := strconv.Atoi(operand); err == nil {
			if constant < -50 || constant > 99 {
				inst.err = fmt.Errorf("constant %d is outside valid range -50 to 99", constant)
				return
			}
			inst.constant = constant
		} else if len(operand) > 0 && operand[0] == 'x' {
			i, err := strconv.Atoi(operand[1:])
			if err != nil {
				inst.err = fmt.Errorf("unable to parse index %s because %w", operand[1:], err)
				return
			}
			if i < 0 || i > 9 {
				inst.err = fmt.Errorf("x0 to x9 are the only valid registers, not x%d", i)
				return
			}
			registers = append(registers, uint(i))
		}
	}
	inst.parsedRegIndicies = registers
}

// Figures out how parsedRegIndicies should be assigned to regIndicies.
// Why are these operations not done in one go? Because different instructions do it differently while
// ParseOperands which fills in parsedRegIndicies is the same for all instructions and can thus be reused.
// Register assignment will differ between instructions
// The most common register assigment works as follows:
//
//	Rd, Ra -> Rd, Rd, Ra
//	Rd, k -> Rd, Rd, k
func (inst *BaseInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
	n := len(inst.parsedRegIndicies)
	switch n {
	case 2:
		inst.regIndicies[Ra] = inst.parsedRegIndicies[0]
		inst.regIndicies[Rb] = inst.parsedRegIndicies[1]
	case 3:
		inst.regIndicies[Ra] = inst.parsedRegIndicies[1]
		inst.regIndicies[Rb] = inst.parsedRegIndicies[2]
	default:
		inst.err = fmt.Errorf("instruction expects 2 or 3 register operands but you gave %d", n)
	}
}

func (inst *BaseInstruction) DecodeOperands(operands uint) {
	addr := operands % 100

	inst.regIndicies[Rd] = uint(operands / 100)
	inst.regIndicies[Ra] = uint(addr / 10)
	inst.regIndicies[Rb] = uint(addr % 10)
}

func (inst *BaseInstruction) Error() error {
	return inst.err
}
