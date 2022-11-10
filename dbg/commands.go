package dbg

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/ordovician/calcutron/asm"
	"github.com/ordovician/calcutron/disasm"
	"github.com/ordovician/calcutron/prog"
	"github.com/ordovician/calcutron/sim"
	"github.com/ordovician/calcutron/utils"
)

type Command interface {
	Name() string
	Help(writer io.Writer)
	Action(writer io.Writer, comp *sim.Computer, args []string) error
}

type HelpCmd struct{}
type InputCmd struct{}
type OutputCmd struct{}
type LoadCmd struct{}
type ListCmd struct{}
type StatusCmd struct{}
type NextCmd struct{}
type RunCmd struct{}
type PrintCmd struct{}
type SetCmd struct{}
type AsmCmd struct{}
type ResetCmd struct{}
type MemoryCmd struct{}

func (cmd *HelpCmd) Name() string {
	return "help"
}

func (cmd *HelpCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	help -- get help about supported commands
SYNOPSIS
	help
DESCRIPTION
	get help about supported commands`)
}

func (cmd *HelpCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	if len(args) == 0 {
		for _, cmd := range commands {
			fmt.Fprintln(writer, cmd.Name())
		}
	} else if len(args) == 1 {
		helpCmd := Lookup(args[0])
		if helpCmd != nil {
			helpCmd.Help(writer)
		} else {
			fmt.Fprintf(writer, "could not find command with name %s\n", args[0])
		}
	} else {
		return fmt.Errorf("help takes none or one argument, not %d", len(args))
	}

	return nil
}

func (cmd *InputCmd) Name() string {
	return "inputs"
}

func (cmd *InputCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	inputs -- set inputs to program
SYNOPSIS
	inputs
DESCRIPTION
	load input data`)
}

func (cmd *InputCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	if len(args) == 0 {
		comp.PrintInputs(writer)
		return nil
	}

	elements := make([]uint, len(args))
	for i, arg := range args {
		x, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("cannnot parse argument %d to inputs because %w", i, err)
		}
		elements[i] = prog.Complement(x, 1e4)
	}
	comp.SetInputs(elements)
	return nil
}

func (cmd *OutputCmd) Name() string {
	return "outputs"
}

func (cmd *OutputCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	outputs -- show output from program
SYNOPSIS
	outputs
DESCRIPTION
	show all numbes written to output on address 99`)
}

func (cmd *OutputCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	comp.PrintOutputs(writer)
	return nil
}

func (cmd *LoadCmd) Name() string {
	return "load"
}

func (cmd *LoadCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	load -- loads program into simulator
SYNOPSIS
	load filepath
DESCRIPTION
	loads either machine code or source code into memory
	ready for execution.`)
}

func (cmd *LoadCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("cannot load code, missing file argument")
	}
	return comp.LoadFile(args[0])
}

func (cmd *ListCmd) Name() string {
	return "list"
}

func (cmd *ListCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	list -- list code of loaded program
SYNOPSIS
	list
DESCRIPTION
	disassembles code in memory and shows it`)
}

func (cmd *ListCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	memory := comp.ProgramSlice()
	program, err := disasm.DisassembleMemory(memory)
	if err != nil {
		return err
	}
	program.PrintWithOptions(writer, &prog.PrintOptions{
		Address:     true,
		MachineCode: true,
		SourceCode:  true,
	})
	return nil
}

func (cmd *StatusCmd) Name() string {
	return "status"
}

func (cmd *StatusCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	status -- show value of registers, input and output
SYNOPSIS
	status
DESCRIPTION
	shows internal status of computer.`)
}

func (cmd *StatusCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	comp.Print(writer)
	return nil
}

func (cmd *NextCmd) Name() string {
	return "next"
}

func (cmd *NextCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	next -- make one step into code
SYNOPSIS
	next
DESCRIPTION
	execute a single assembly code instruction`)
}

func (cmd *NextCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	inst := comp.Instruction()

	prog.AddressColor.Fprintf(writer, "%02d ", comp.PC())
	prog.GrayColor.Fprintf(writer, "%04d ", inst.MachineCode())
	fmt.Fprintln(writer, inst.SourceCode())

	comp.Step()

	switch inst.Opcode() {
	case prog.BRA, prog.BEQ, prog.BGT, prog.BLT, prog.JUMP:
		comp.PrintPC(writer)
	case prog.HALT:
		break
	default:
		comp.PrintRegs(writer, inst.UniqueRegisters()...)
	}

	return nil
}

func (cmd *RunCmd) Name() string {
	return "run"
}

func (cmd *RunCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	run -- run loaded program or specified instruction
SYNOPSIS
	run [instruction]
DESCRIPTION
	Execute loaded from from start until a halting instruction is met or input is exausted.
	If an assembly code instruction is specified it will be parsed and run`)
}

func (cmd *RunCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	// check if we should run an individual instruction
	line := strings.TrimSpace(strings.Join(args, " "))
	if len(line) > 0 && utils.AllDigits(line) {
		machinecode, _ := strconv.Atoi(line)
		inst := disasm.DisassembleInstruction(uint(machinecode))
		inst.Run(comp)
		return nil
	} else if len(line) > 0 {
		labels := make(prog.SymbolTable)
		inst, err := asm.AssembleLine(labels, line, 0)
		if err != nil {
			return nil
		}
		if inst != nil {
			inst.Run(comp)
		}
		return nil
	}

	memory := comp.ProgramSlice()
	program, err := disasm.DisassembleMemory(memory)
	if err != nil {
		return err
	}

	pContext := prog.NewPrintContext(program.Labels, &prog.PrintOptions{
		Address:     true,
		MachineCode: true,
		SourceCode:  true,
	})

	var group sync.WaitGroup
	group.Add(1)
	channel := make(chan prog.AddressInstruction)
	go func() {
		pContext.Print(writer, channel)
		group.Done()
	}()

	comp.RunChannel(512, channel)
	close(channel)

	// wait until executed instuctions have been printed out to consol
	group.Wait()

	// Print out status of computer
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, comp.String())
	return nil
}

func (cmd *PrintCmd) Name() string {
	return "print"
}

func (cmd *PrintCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	print -- print value of register
SYNOPSIS
	print register
DESCRIPTION
	print value of specified register from x1 to x9 or program counter PC`)
}

func (cmd *PrintCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("must specify at least one register to print")
	}

	regStr := args[0]
	if strings.ToUpper(regStr) == "PC" {
		fmt.Printf("PC: ")
		prog.NumberColor.Fprintf(writer, "%02d\n", comp.PC())
		return nil
	}

	regIndex, err := strconv.Atoi(regStr[1:])
	if err != nil {
		return fmt.Errorf("cannot parse index of register %s because %w", regStr, err)
	}
	if regIndex < 0 || regIndex > 9 {
		return fmt.Errorf("register index %d is outside of valid range 1 to 9", regIndex)
	}

	fmt.Printf("x%d: ", regIndex)
	prog.NumberColor.Fprintf(writer, "%d\n", comp.Register(uint(regIndex)))

	return nil
}

func (cmd *SetCmd) Name() string {
	return "set"
}

func (cmd *SetCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	set -- set value of register
SYNOPSIS
	set register value
DESCRIPTION
	set value of specified register from x1 to x9 or program counter PC`)
}

func (cmd *SetCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("command must be of the form 'set x3 42'.\nyou need a register and a value to set")
	}

	value, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("unable to parse value '%s` because %w", args[1], err)
	}

	regStr := args[0]
	if strings.ToUpper(regStr) == "PC" {
		if value < 0 || value > 99 {
			return fmt.Errorf("value %d is more than two digits. Program counter can only hold two digit values", value)
		}

		comp.SetPC(uint(value))

		fmt.Printf("PC: ")
		prog.NumberColor.Fprintf(writer, "%02d\n", comp.PC())
		return nil
	}

	regIndex, err := strconv.Atoi(regStr[1:])
	if err != nil {
		return fmt.Errorf("cannot parse index of register %s because %w", regStr, err)
	}
	if regIndex < 1 || regIndex > 9 {
		return fmt.Errorf("register index %d is outside of valid range 1 to 9", regIndex)
	}

	if value < -9999 || value > 9999 {
		return fmt.Errorf("value %d is more than four digits. Registers can only hold four digit values", value)
	}

	comp.SetRegister(uint(regIndex), value)

	fmt.Printf("x%d: ", regIndex)
	prog.NumberColor.Fprintf(writer, "%d\n", comp.Register(uint(regIndex)))

	return nil
}

func (cmd *AsmCmd) Name() string {
	return "asm"
}

func (cmd *AsmCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	asm -- assemble instruction
SYNOPSIS
	asm
DESCRIPTION
	give assembly code for given instruction`)
}

func (cmd *AsmCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	labels := make(prog.SymbolTable)
	line := strings.Join(args, " ")
	inst, err := asm.AssembleLine(labels, line, 0)
	if err != nil {
		return err
	}
	// do a roundtrip so we can see how instruction is interpreted
	inst = disasm.DisassembleInstruction(inst.MachineCode())

	program := prog.Program{
		Instructions: []prog.Instruction{inst},
	}

	program.PrintWithOptions(writer, &prog.PrintOptions{
		MachineCode: true,
		SourceCode:  true,
	})

	return nil
}

func (cmd *ResetCmd) Name() string {
	return "reset"
}

func (cmd *ResetCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	reset -- reset inputs, output, program counter and registers
SYNOPSIS
	reset
DESCRIPTION
	reset computer state so program can run over again`)
}

func (cmd *ResetCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	comp.Reset()
	return nil
}

func (cmd *MemoryCmd) Name() string {
	return "memory"
}

func (cmd *MemoryCmd) Help(writer io.Writer) {
	fmt.Fprintln(writer,
		`NAME
	memory -- read contents of memory at given location
SYNOPSIS
	memory address [value]
DESCRIPTION
	print contents of memory at given address`)
}

func (cmd *MemoryCmd) Action(writer io.Writer, comp *sim.Computer, args []string) error {
	var addr int
	var err error
	if len(args) > 0 {
		addr, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("unable to parse address for memory instruction because %w", err)
		}
	}

	switch len(args) {
	case 1:
		prog.AddressColor.Fprintf(writer, "%02d ", addr)
		prog.NumberColor.Fprintf(writer, "%04d\n", comp.Memory(uint(addr)))
	case 2:
		value, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("unable to parse value for memory instruction because %w", err)
		}
		comp.SetMemory(uint(addr), prog.Complement(value, 1e4))
	default:
		return fmt.Errorf("you need to provide address of memory cell you want to view or modify. That mean 1 or 2 arguments, not %d", len(args))
	}
	return nil
}

var commands = [...]Command{
	new(HelpCmd),
	new(InputCmd),
	new(OutputCmd),
	new(ListCmd),
	new(LoadCmd),
	new(StatusCmd),
	new(NextCmd),
	new(RunCmd),
	new(PrintCmd),
	new(SetCmd),
	new(AsmCmd),
	new(ResetCmd),
	new(MemoryCmd),
}

// Lookup command with given name. Returns nil if command
// of that name is not found
func Lookup(cmdName string) Command {
	for _, cmd := range commands {
		if cmd.Name() == cmdName {
			return cmd
		}
	}
	return nil
}
