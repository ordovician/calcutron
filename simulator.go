package calcutron

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// "bufio"

var ErrAllInputRead = errors.New("all inputs read")
var ErrProgramHalt = errors.New("reach halt instruction")

type Computer struct {
	PC        uint8      // Program counter 0-99
	Registers [10]uint8  // CPU registers   0-99
	Memory    [99]uint16 // Computer memory 0-9999
	Inputs    []uint8    // Input data to computer 0-99
	Outputs   []uint8    // Output from computer    0-99
	inpos     int        // Current position input stream
}

func NewComputer(program []uint16) *Computer {
	if len(program) > 99 {
		panic("programs cannot be longer than 99 instructions")
	}

	var comp Computer
	copy(comp.Memory[0:], program)

	return &comp
}

// Load program into computer from file
func (comp *Computer) LoadFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("could not load program from file: %w", err)
	}
	defer file.Close()

	return comp.Load(file)
}

// Load program into computer from reader
func (comp *Computer) Load(reader io.Reader) error {

	scanner := bufio.NewScanner(reader)
	for addr := 0; scanner.Scan(); addr++ {
		// Ignore comments
		line := scanner.Text()
		if i := strings.Index(line, "//"); i >= 0 {
			line = line[0:i]
		}

		instruction, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("failed to parse machine code because %w", err)
		}

		comp.Memory[addr] = uint16(instruction)
	}
	return nil
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
		comp.Inputs = append(comp.Inputs, uint8(input))
	}
	return nil
}

// Doesn't erase installed program of set input but resets everything so
// program can be run over again and give same result
func (comp *Computer) Reset() {
	comp.PC = 0
	comp.inpos = 0
	comp.Outputs = []uint8{}
	for i := range comp.Registers {
		comp.Registers[i] = 0
	}
}

// Execute instruction at current address. Where the program counter (PC) is.
func (comp *Computer) Step() error {
	pc := comp.PC
	ir := comp.Memory[pc]
	err := comp.ExecuteInstruction(ir)
	if err != nil {
		return fmt.Errorf("unable to execute instruction at address %d because: %w", pc, err)
	}

	// Make sure we didn't execute a branch instruction before updating Program counter
	if pc == comp.PC {
		comp.PC += 1
	}

	return nil
}

// Execute N instructions
func (comp *Computer) RunSteps(nsteps int) error {
	for i := 0; i < nsteps; i++ {
		comp.PrintCurrentInstruction()
		err := comp.Step()
		if err != nil {
			return fmt.Errorf("could not run program because %w", err)
		}
	}
	return nil
}

// Execute given instruction without altering the program counter (PC)
// unless running a branch instruction which performs a branch
func (comp *Computer) ExecuteInstruction(instruction uint16) error {
	if instruction > 9999 {
		return fmt.Errorf("instruction %d not within valid range 0000 - 9999", instruction)
	}

	regs := comp.Registers[0:]

	//inst := decodeInstruction(instruction)

	pinst := DisassembleInstruction(instruction)

	var rd uint8

	if pinst.DestReg() <= 9 {
		rd = regs[pinst.DestReg()]
	}

	switch pinst.opcode {
	case ADD:
		rd = regs[pinst.FirstSourceReg()] + regs[pinst.SecondSourceReg()]
	case SUB:
		rd = regs[pinst.FirstSourceReg()] - regs[pinst.SecondSourceReg()]
	case SUBI:
		rd = regs[pinst.FirstSourceReg()] - pinst.Constant()
	case LSH:
		rd = regs[pinst.FirstSourceReg()] * (10 ^ pinst.Constant())
	case RSH:
		rd = regs[pinst.FirstSourceReg()] % (10 ^ pinst.Constant())
		regs[pinst.FirstSourceReg()] = regs[pinst.FirstSourceReg()] / (10 ^ pinst.Constant())
	case BRZ:
		if rd == 0 {
			comp.PC = pinst.Constant()
		}
	case BGT:
		if rd > 0 {
			comp.PC = pinst.Constant()
		}
	case LD:
		if pinst.Constant() < 90 {
			rd = uint8(comp.Memory[pinst.Constant()])
		} else if pinst.Constant() == 90 {
			if comp.inpos >= len(comp.Inputs) {
				return ErrAllInputRead
			}
			rd = comp.Inputs[comp.inpos]
			comp.inpos += 1
		}
	case ST:
		if pinst.Constant() < 90 {
			comp.Memory[pinst.Constant()] = uint16(rd)
		} else if pinst.Constant() == 91 {
			comp.Outputs = append(comp.Outputs, rd)
		} else {
			return fmt.Errorf("writing to address %d is not supported in this version", pinst.Constant())
		}
	case HLT:
		return ErrProgramHalt
	default:
		return fmt.Errorf("opcode %d, is not supported. Must be between 0-9", pinst.opcode)
	}

	// Make sure register values stay within range 0-99
	// Act as if a register can hold two digits
	if rd != 0 {
		rd = rd % 100
	}

	// Cannot write to register 0, so it is excluded
	if pinst.DestReg() >= 1 && pinst.DestReg() <= 9 {
		regs[pinst.DestReg()] = rd
	}

	return nil
}

// Print instruction at program counter disassembled with address and machine code
// Example:
//
//	02: 1432; ADD  x4, x3, x2
func (comp *Computer) PrintCurrentInstruction() {
	machinecode := comp.Memory[comp.PC]
	instruction := decodeInstruction(machinecode)
	fmt.Printf("%02d: %04d; %v\n", comp.PC, machinecode, instruction)
}

func (comp *Computer) String() string {
	buffer := bytes.NewBufferString("")

	fmt.Fprintln(buffer, "PC:", comp.PC)
	for i, reg := range comp.Registers {
		fmt.Fprintf(buffer, "x%d: %d, ", i, reg)
	}
	fmt.Fprintln(buffer)

	fmt.Fprintf(buffer, "Inputs:  ")
	Join(buffer, comp.Inputs, ", ")
	fmt.Fprintln(buffer)

	fmt.Fprintf(buffer, "Outputs: ")
	Join(buffer, comp.Outputs, ", ")
	fmt.Fprintln(buffer)

	// for _, output := range comp.Outputs {
	// 	fmt.Fprintf(buffer, "%d, ", output)
	// }
	// fmt.Fprintln(buffer)

	return buffer.String()
}
