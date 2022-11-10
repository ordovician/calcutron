package asm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ordovician/calcutron/prog"
	"github.com/ordovician/calcutron/utils"
)

// Get the mnemonic and operands of a source code line
func parseLine(line string) (mnemonic string, operands []string) {
	operands = make([]string, 0)

	code := strings.Trim(line, " \t")
	i := len(code)
	if j := strings.Index(code, "//"); j >= 0 {
		i = j
	}
	n := len(code)

	if n == 0 || code[n-1] == ':' {
		return
	}

	code = code[0:i]
	if i = strings.IndexRune(code, ' '); i < 0 {
		i = n
	}
	mnemonic = code[0:i]

	if len(code[i:]) > 0 {
		operands = strings.SplitN(code[i:], ",", 3)

		// Cleanup white space around operands so function further
		// down the chain don't have to deal with thm
		for i, oper := range operands {
			operands[i] = strings.TrimSpace(oper)
		}
	}

	return
}

// When we assemble an instruction the address in the program of the instruction can affect the machine code generated
// because some instructions such as JMP use relative jumps. Thus the address part of the JMP depends on
// where the JMP instruction is assembled. If you don't care about the address, just set the address to zero.
// Then relative positions will look like absolute positions
func AssembleLine(labels prog.SymbolTable, line string, address uint) (prog.Instruction, error) {
	mnemonic, operands := parseLine(line)
	if mnemonic == "" {
		return nil, nil
	}

	opcode, ok := prog.ParseOpcode(mnemonic)
	if !ok {
		return nil, fmt.Errorf("%s is not a legal mnemonic", mnemonic)
	}

	inst := prog.NewInstruction(opcode)
	inst.ParseOperands(labels, operands, address)
	inst.AssignRegisters()

	return inst, inst.Error()
}

// Assembler reads assembly code from reader and writes machine code to writer
func Assemble(reader io.ReadSeeker) (*prog.Program, error) {
	labels := prog.ReadSymTable(reader)
	labels.AddIOLabels() // so we got labels like input and output
	program := prog.Program{
		Labels:       labels,
		Instructions: make([]prog.Instruction, 0, 10),
	}

	reader.Seek(0, io.SeekStart)

	scanner := bufio.NewScanner(reader)

	var addr uint = 0
	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := scanner.Text()
		instruction, err := AssembleLine(program.Labels, line, addr)
		if err != nil {
			return nil, fmt.Errorf("%2d %s: unable to assemble: %w", lineNo, strings.TrimSpace(line), err)
		}
		if instruction != nil {
			program.Add(instruction)
			addr++
		}
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("unable to assemble file: %w", scanner.Err())
	}

	return &program, nil
}

// AssembleFile reads assembly code from file at path filepath and write machinecode to writer
func AssembleFile(filepath string) (*prog.Program, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	file.Seek(0, io.SeekStart)

	if !strings.HasSuffix(filepath, ".ct33") && utils.AllDigits(line) {
		return nil, fmt.Errorf("file '%s' doesn't look like an assembly code file.\nFirst line, '%s', is a number. Are you sure this isn't a machine code file?", filepath, line)
	}

	return Assemble(file)
}
