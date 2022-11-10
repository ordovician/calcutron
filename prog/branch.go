package prog

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type BranchInstruction struct {
	BaseInstruction
}

func (inst *BranchInstruction) AssignRegisters() {
	if inst.err != nil {
		return
	}
	n := len(inst.parsedRegIndicies)
	if n < 3 {
		inst.err = fmt.Errorf("conditional branch instructions take 3 operands not %d", n)
	} else {
		inst.regIndicies[Rd] = inst.parsedRegIndicies[0]
		inst.regIndicies[Ra] = inst.parsedRegIndicies[1]
		inst.regIndicies[Rb] = inst.parsedRegIndicies[2]
	}
}

type BranchEqualInstruction struct {
	BranchInstruction
}

func (inst *BranchEqualInstruction) Run(comp Machine) bool {
	left := inst.RegValue(comp, Rd)
	right := inst.RegValue(comp, Ra)
	addr := inst.RegValue(comp, Rb)

	if left == right {
		comp.SetPC(uint(addr))
	}
	return true
}

type BranchGreaterThanInstruction struct {
	BranchInstruction
}

func (inst *BranchGreaterThanInstruction) Run(comp Machine) bool {
	left := inst.RegValue(comp, Rd)
	right := inst.RegValue(comp, Ra)
	addr := inst.RegValue(comp, Rb)

	if left > right {
		comp.SetPC(uint(addr))
	}
	return true
}

type RJumpInstruction struct {
	ImmediateInstruction
}

func (inst *RJumpInstruction) Run(comp Machine) bool {
	// Set return address
	inst.SetRegValue(comp, Rd, int(comp.PC()+1))

	// jumping back to same instruction will create an infinite loop
	// hence this is a terminating instruction
	if inst.constant == 0 {
		return false
	}

	addr := int(comp.PC()) + inst.constant

	comp.SetPC(uint(addr))
	return true
}

func (inst *RJumpInstruction) ParseOperands(labels SymbolTable, operands []string, offset uint) {
	var regIndex uint // x0 is the default
	var jmpAddr string
	switch len(operands) {
	case 2:
		regstr := operands[0]
		i := strings.IndexFunc(regstr, func(r rune) bool {
			return unicode.IsDigit(r)
		})
		if i <= 0 {
			inst.err = fmt.Errorf("couldn't find register index in string '%s'", regstr)
		}
		if regstr[0:i] != "x" {
			inst.err = fmt.Errorf("register names start with an 'x' not '%s'", regstr[0:i])
			return
		}
		i, err := strconv.Atoi(regstr[i:])
		if err != nil {
			inst.err = fmt.Errorf("unable to parse index %s because %w", regstr[1:], err)
			return
		}
		if i < 0 || i > 9 {
			inst.err = fmt.Errorf("x0 to x9 are the only valid registers, not x%d", i)
			return
		}
		regIndex = uint(i)

		jmpAddr = operands[1]
	case 1:
		jmpAddr = operands[0]
	default:
		inst.err = fmt.Errorf("RJUMP takes one or two operands not %d", len(operands))
		return
	}

	if addr, ok := labels[jmpAddr]; ok {
		inst.constant = int(addr) - int(offset) // Calculate relative address
		inst.label = jmpAddr
	} else if constant, err := strconv.Atoi(jmpAddr); err == nil {
		if constant < -99 || constant > 99 {
			inst.err = fmt.Errorf("address %d is outside valid range -99 to 99", constant)
			return
		}
		inst.constant = constant
	}
	inst.parsedRegIndicies = []uint{regIndex}
}
