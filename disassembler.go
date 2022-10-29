package calcutron

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func Disassemble(reader io.Reader) error {
	return DisassembleWithOptions(reader, SOURCE_CODE)
}

// Disassemble machinecode read from reader
func DisassembleWithOptions(reader io.Reader, option AssemblyFlag) error {
	scanner := bufio.NewScanner(reader)

	var line SourceCodeLine
	line.address = 0
	for line.lineno = 1; scanner.Scan(); line.lineno++ {
		line.sourcecode = scanner.Text()
		machinecode, err := strconv.Atoi(line.sourcecode)
		line.machinecode = int16(machinecode)
		if err != nil {
			return fmt.Errorf("%d: unable to disassemble because: %w", line.lineno, err)
		}
		if line.machinecode < 0 {
			log.Panicf("%d: something went from in parsing code. Machine code instruction should never be less than 0", line.lineno)
		}

		instruction := DisassembleInstruction(uint16(machinecode))

		fmt.Printf("%02d: %04d; %v\n", line.lineno-1, line.machinecode, instruction)
	}

	return nil
}

func DisassembleFile(filepath string) error {
	return DisassembleFileWithOptions(filepath, SOURCE_CODE)
}

// Disassemble file and write output to stdout
func DisassembleFileWithOptions(filepath string, options AssemblyFlag) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return DisassembleWithOptions(file, options)
}
