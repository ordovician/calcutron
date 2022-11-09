package sim

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ordovician/calcutron/asm"
	"github.com/ordovician/calcutron/disasm"
	"github.com/ordovician/calcutron/prog"
	"github.com/ordovician/calcutron/utils"
)

var ErrAllInputRead = errors.New("all inputs read")
var ErrProgramHalt = errors.New("reach halt instruction")

type Computer struct {
	pc        uint     // Program counter 0-99
	registers [10]uint // CPU registers   0-99
	memory    [99]uint // Computer memory 0-9999
	inputs    []uint   // Input data to computer 0-99
	outputs   []uint   // Output from computer    0-99
	inpos     int      // Current position input stream
}

// valid registers are in range 0 to 9, but register 0 will always contains 0
func (comp *Computer) Register(i uint) int {
	if i > 0 && i <= 9 {
		return prog.Signed(comp.registers[i], 1e4)
	} else {
		return 0
	}
}

// valid registers are in range 0 to 9, but register 0 will never get altered
func (comp *Computer) SetRegister(i uint, value int) {
	if i > 0 && i <= 9 {
		comp.registers[i] = prog.Complement(value, 1e4)
	}
}

func (comp *Computer) PC() uint {
	return comp.pc
}

func (comp *Computer) SetPC(address uint) {
	comp.pc = address
}

func (comp *Computer) Memory(address uint) uint {
	return comp.memory[address]
}

// slice of memory representing a program
func (comp *Computer) ProgramSlice() []uint {
	n := len(comp.memory) - 1
	for ; n > 1; n-- {
		if comp.memory[n-1] != 0 {
			break
		}
	}

	return comp.memory[:n+1]
}

func (comp *Computer) SetMemory(address uint, value uint) {
	comp.memory[address] = value
}

func (comp *Computer) PopInput() (int, bool) {
	if comp.inpos >= len(comp.inputs) {
		return 0, false
	}
	input := prog.Signed(comp.inputs[comp.inpos], 1e4)
	comp.inpos++
	return input, true
}

func (comp *Computer) PushOutput(value int) {
	comp.outputs = append(comp.outputs, prog.Complement(value, 1e4))
}

func NewComputer(program *prog.Program) *Computer {
	if len(program.Instructions) > 99 {
		panic("programs cannot be longer than 99 instructions")
	}

	var comp Computer
	comp.LoadProgram(program)

	return &comp
}

func NewComputerFile(filepath string) (*Computer, error) {
	var comp Computer
	err := comp.LoadFile(filepath)
	return &comp, err
}

// Doesn't erase installed program of set input but resets everything so
// program can be run over again and give same result
func (comp *Computer) Reset() {
	comp.pc = 0
	comp.inpos = 0
	comp.outputs = make([]uint, 0)
	for i := range comp.registers {
		comp.registers[i] = 0
	}
}

func (comp *Computer) LoadProgram(program *prog.Program) {
	memory := comp.memory[:]
	for i, inst := range program.Instructions {
		machinecode := inst.MachineCode()
		memory[i] = machinecode
	}
}

func (comp *Computer) LoadFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("could not load program from file: %w", err)
	}
	defer file.Close()

	if strings.HasSuffix(filepath, ".ct33") {
		return comp.LoadSourceCode(file)
	} else if strings.HasSuffix(filepath, ".machine") {
		return comp.LoadMachineCode(file)
	}
	return fmt.Errorf("unknown file suffix")
}

// Load program into computer from reader
func (comp *Computer) LoadMachineCode(reader io.Reader) error {
	program, err := disasm.Disassemble(reader)
	if err != nil {
		return err
	}
	comp.LoadProgram(program)
	return nil
}

func (comp *Computer) LoadSourceCode(reader io.ReadSeeker) error {
	program, err := asm.Assemble(reader)
	if err != nil {
		return err
	}
	comp.LoadProgram(program)
	return nil
}

func (comp *Computer) SetInputs(elements []uint) {
	comp.inputs = elements
}

// Load data input to program from file
func (comp *Computer) LoadInputsFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("could not inputs from file: %w", err)
	}
	defer file.Close()

	return comp.LoadInputs(file)
}

// Load data input to program from reader
func (comp *Computer) LoadInputs(reader io.Reader) error {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("unable to read input data because %w", err)
	}

	strings.Fields(string(bytes))
	for _, word := range strings.Fields(string(bytes)) {
		input, err := strconv.Atoi(word)
		if err != nil {
			return fmt.Errorf("failed to parse machine code because %w", err)
		}
		comp.inputs = append(comp.inputs, uint(input))
	}
	return nil
}

func (comp *Computer) StringInputs(s string) error {
	buffer := bytes.NewBufferString(s)
	return comp.LoadInputs(buffer)
}

func (comp *Computer) Instruction() prog.Instruction {
	pc := comp.pc
	machinecode := comp.memory[pc]
	instruction := disasm.DisassembleInstruction(machinecode)
	return instruction
}

func (comp *Computer) StepChannel(out chan<- prog.AddressInstruction) bool {
	pc := comp.pc

	inst := comp.Instruction()

	out <- prog.AddressInstruction{
		Addr: pc,
		Inst: inst,
	}
	// Check if we have reached a terminating instruction
	if !inst.Run(comp) {
		return false
	}

	// Make sure we didn't execute a branch instruction before updating Program counter
	if pc == comp.pc {
		comp.pc += 1
	}
	return true
}

func (comp *Computer) RunChannel(nsteps int, out chan<- prog.AddressInstruction) {

	for i := 0; i < nsteps; i++ {
		if !comp.StepChannel(out) {
			break
		}
	}
}

func (comp *Computer) Step() bool {
	pc := comp.pc
	inst := comp.Instruction()

	// Check if we have reached a terminating instruction
	if !inst.Run(comp) {
		return false
	}

	// Make sure we didn't execute a branch instruction before updating Program counter
	if pc == comp.pc {
		comp.pc += 1
	}
	return true
}

func (comp *Computer) Run(nsteps int) {
	for i := 0; i < nsteps; i++ {
		if !comp.Step() {
			break
		}
	}
}

func (comp *Computer) PrintRegs(writer io.Writer, indices ...uint) {

	for i, index := range indices {
		// skip the zero register, since it is always zero
		if index == 0 {
			continue
		}

		if i > 0 {
			fmt.Fprint(writer, ", ")
		}
		fmt.Fprintf(writer, "x%d: ", index)
		prog.NumberColor.Fprintf(writer, "%04d", prog.Complement(comp.Register(index), 1e4))
	}
	fmt.Fprintln(writer)
}

func (comp *Computer) PrintPC(writer io.Writer) {
	fmt.Fprint(writer, "PC: ")
	prog.NumberColor.Fprintf(writer, "%02d\n", comp.pc)
}
func (comp *Computer) PrintInputs(writer io.Writer) {
	numberColor := prog.NumberColor.SprintFunc()
	grayColor := prog.GrayColor.SprintFunc()

	fmt.Fprintf(writer, "Inputs:  ")
	// using gray and pink to separate consumed and non-consumed input data
	utils.JoinFunc(writer, comp.inputs[:comp.inpos], ", ", grayColor)
	if comp.inpos > 0 && comp.inpos < len(comp.inputs) {
		fmt.Fprintf(writer, ", ")
	}
	utils.JoinFunc(writer, comp.inputs[comp.inpos:], ", ", numberColor)
	fmt.Fprintln(writer)
}

func (comp *Computer) PrintOutputs(writer io.Writer) {
	numberColor := prog.NumberColor.SprintFunc()

	fmt.Fprintf(writer, "Outputs: ")
	utils.JoinFunc(writer, comp.outputs, ", ", numberColor)
	fmt.Fprintln(writer)
}

func (comp *Computer) Print(writer io.Writer) {

	comp.PrintPC(writer)
	fmt.Fprintln(writer)
	comp.PrintRegs(writer, 1, 4, 7)
	comp.PrintRegs(writer, 2, 5, 8)
	comp.PrintRegs(writer, 3, 6, 9)
	fmt.Fprintln(writer)

	comp.PrintInputs(writer)
	comp.PrintOutputs(writer)
}

func (comp *Computer) String() string {
	buffer := bytes.NewBufferString("")
	comp.Print(buffer)
	return buffer.String()
}
