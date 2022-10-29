package calcutron

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// A table containing the memory address of labels in the code
func readSymTable(reader io.Reader) map[string]uint8 {
	scanner := bufio.NewScanner(reader)
	labels := make(map[string]uint8)
	address := 0
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		n := len(line)

		if n == 0 {
			continue
		}

		if i := strings.IndexRune(line, ':'); i >= 0 {
			labels[line[0:i]] = uint8(address)

			// is there anything beyond the label?
			if n == i+1 {
				continue
			}
		}
		address++
	}
	return labels
}

// Assemble single line of code
func AssembleLine(labels map[string]uint8, line string) (int16, error) {
	code := strings.Trim(line, " \t")
	i := len(code)
	if j := strings.Index(code, "//"); j >= 0 {
		i = j
	}
	n := len(code)

	if n == 0 || code[n-1] == ':' {
		return -1, nil
	}

	code = code[0:i]
	if i = strings.IndexRune(code, ' '); i < 0 {
		i = n
	}
	mnemonic := code[0:i]
	var operands []string

	if len(code[i:]) > 0 {
		operands = strings.SplitN(code[i:], ",", 3)
	}

	opcode := ParseOpcode(mnemonic)
	instruction := Instruction{
		opcode: opcode,
	}

	err := instruction.ParseOperands(labels, operands)
	if err != nil {
		return 0, err
	}

	machincode, err := instruction.Machinecode()

	return int16(machincode), err
}

func Assemble(reader io.ReadSeeker, writer io.Writer) error {
	return AssembleWithOptions(reader, writer, AssemblyFlag(0))
}

// Assembler reads assembly code from reader and writes machine code to writer
func AssembleWithOptions(reader io.ReadSeeker, writer io.Writer, options AssemblyFlag) error {
	labels := readSymTable(reader)

	reader.Seek(0, io.SeekStart)

	scanner := bufio.NewScanner(reader)

	var line SourceCodeLine
	line.address = 0
	for line.lineno = 1; scanner.Scan(); line.lineno++ {

		line.sourcecode = strings.Trim(scanner.Text(), " \t")
		var err error
		line.machinecode, err = AssembleLine(labels, line.sourcecode)
		if err != nil {
			return err
		}

		if line.machinecode >= 0 {
			PrintInstruction(writer, line, options|MACHINE_CODE)
			line.address++
		}
	}

	return nil
}

func AssembleFile(filepath string, writer io.Writer) error {
	return AssembleFileWithOptions(filepath, writer, MACHINE_CODE)
}

// AssembleFile reads assembly code from file at path filepath and write machinecode to writer
func AssembleFileWithOptions(filepath string, writer io.Writer, options AssemblyFlag) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return AssembleWithOptions(file, writer, options)
}
