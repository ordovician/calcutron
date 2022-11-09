package disasm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ordovician/calcutron/prog"
	"github.com/ordovician/calcutron/utils"
)

func DisassembleInstruction(machinecode uint) prog.Instruction {
	opcode := prog.Opcode(machinecode / 1000)

	operands := machinecode % 1000
	inst := prog.NewInstruction(opcode)
	inst.DecodeOperands(operands)

	return inst
}

// Disassemble a machine code program read from reader
func Disassemble(reader io.Reader) (*prog.Program, error) {
	machineprogram := make([]uint, 0, 10)
	scanner := bufio.NewScanner(reader)

	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := strings.TrimSpace(scanner.Text())
		if !utils.AllDigits(line) {
			return nil, fmt.Errorf("machine code must be all digits. Cannot disassemble '%s'", line)
		}
		machinecode, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("%d: unable to disassemble because: %w", lineNo, err)
		}
		machineprogram = append(machineprogram, uint(machinecode))
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("unable to disassemble file: %w", scanner.Err())
	}

	return DisassembleMemory(machineprogram)
}

// Disassemble a section of our made up computer memory
func DisassembleMemory(machineprogram []uint) (*prog.Program, error) {
	labels := make(prog.SymbolTable)
	labels.AddIOLabels() // so we got labels like input and output
	program := prog.Program{
		Labels:       labels,
		Instructions: make([]prog.Instruction, 0, 10),
	}

	var addr uint = 0
	for _, machinecode := range machineprogram {
		instruction := DisassembleInstruction(uint(machinecode))

		if instruction != nil {
			program.Add(instruction)
			addr++
		}
	}
	return &program, nil
}

// DisassembleFile reads assembly code from file at path filepath and returns program representing code read
func DisassembleFile(filepath string) (*prog.Program, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	file.Seek(0, io.SeekStart)

	if !strings.HasSuffix(filepath, ".machine") && !utils.AllDigits(line) {
		return nil, fmt.Errorf("file '%s' doesn't look like a machine code file.\nFirst line is '%s', which is not a number", filepath, line)
	}

	return Disassemble(file)
}
