package calcutron

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
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
func AssembleLine(labels map[string]uint8, line string) (*Instruction, error) {
	code := strings.Trim(line, " \t")
	i := len(code)
	if j := strings.Index(code, "//"); j >= 0 {
		i = j
	}
	n := len(code)

	if n == 0 || code[n-1] == ':' {
		return nil, nil
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
		return nil, err
	}

	return &instruction, err
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
		instruction, err := AssembleLine(labels, scanner.Text())
		if err != nil {
			return err
		}

		if instruction != nil {
			line.instruction = instruction
			err := line.Print(writer, options|MACHINE_CODE)
			if err != nil {
				return err
			}
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

type SourceCodeLine struct {
	address     int
	instruction *Instruction
	lineno      int
}

func (line *SourceCodeLine) Machinecode() (uint16, error) {
	return line.instruction.Machinecode()
}

func (line *SourceCodeLine) Registers() []uint8 {
	return line.instruction.regs
}

func (line *SourceCodeLine) Constant() uint8 {
	return line.instruction.Constant()
}

func (line *SourceCodeLine) Print(writer io.Writer, options AssemblyFlag) error {
	if options.Has(COLOR) {
		return line.printWithColor(writer, options)
	} else {
		return line.printWithoutColor(writer, options)
	}
}

func (line *SourceCodeLine) printWithoutColor(writer io.Writer, options AssemblyFlag) error {
	if options.Has(ADDRESS) {
		fmt.Fprintf(writer, "%02d: ", line.address)
	}

	machinecode, err := line.Machinecode()
	if err != nil {
		return err
	}
	fmt.Fprintf(writer, "%04d", machinecode)

	if options.Has(SOURCE_CODE) {
		var buffer bytes.Buffer
		line.instruction.PrintSourceCode(&buffer)

		if options.Has(LINE_NO) {
			fmt.Fprintf(writer, "; %-18s", buffer.String())
		} else {

			fmt.Fprintf(writer, "; %s", buffer.String())
		}
	}

	if options.Has(LINE_NO) {
		fmt.Fprintf(writer, " // Line %2d ", line.lineno)
	}

	fmt.Fprintln(writer)
	return nil
}

func (line *SourceCodeLine) printWithColor(writer io.Writer, options AssemblyFlag) error {
	gray := color.New(color.FgHiBlack)
	yellow := color.New(color.FgYellow)

	if options.Has(ADDRESS) {
		yellow.Fprintf(writer, "%02d: ", line.address)
	}

	machinecode, err := line.Machinecode()
	if err != nil {
		return err
	}
	gray.Fprintf(writer, "%04d", machinecode)

	if options.Has(SOURCE_CODE) {
		var buffer bytes.Buffer
		line.instruction.PrintColoredSourceCode(&buffer)
		gray.Fprint(writer, ";")
		if options.Has(LINE_NO) {
			fmt.Fprintf(writer, " %-30s", buffer.String())
		} else {
			fmt.Fprintf(writer, " %s", buffer.String())
		}
	}

	if options.Has(LINE_NO) {
		gray.Fprintf(writer, " // Line %2d ", line.lineno)
	}

	fmt.Fprintln(writer)
	return nil
}
