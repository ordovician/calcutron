package calcutron

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
)

func Example_assembleLine() {
	labels := make(map[string]uint8)

	lines := [...]string{
		"SUBI x9, x8, 7",
		"ADD x1, x3, x2",
		"SUB x2, x4, x1",
		"INP x1",
		"INP x2",
		"CLR x3",
		"OUT x3",
		"ADD x3, x1",
		"CLR x3",
		"DEC x2",
		"MOV x9, x8",
	}

	for _, line := range lines {
		instruction, _ := AssembleLine(labels, line)
		machinecode, _ := instruction.Machinecode()
		fmt.Println(machinecode, line)
	}

	// Output:
	// 3987 SUBI x9, x8, 7
	// 1132 ADD x1, x3, x2
	// 2241 SUB x2, x4, x1
	// 8190 INP x1
	// 8290 INP x2
	// 1300 CLR x3
	// 9391 OUT x3
	// 1331 ADD x3, x1
	// 1300 CLR x3
	// 3221 DEC x2
	// 1908 MOV x9, x8
}

func Example_readSymTable() {
	file, _ := os.Open("testdata/labels-nocode.ct33")
	defer file.Close()

	labels := readSymTable(file)

	for key, value := range labels {
		fmt.Println(value, key)
	}

	// Unordered Output:
	// 0 epsilon
	// 0 alpha
	// 0 beta
	// 0 gamma
	// 0 delta
}

// this is a regression test
func TestRoundTrip(t *testing.T) {
	entries, _ := os.ReadDir("testdata")

	for _, entry := range entries {
		filename := entry.Name()
		// skip files without assembly code, such as machine code files
		if !strings.HasSuffix(filename, ".ct33") {
			continue
		}

		sourcecodepath := path.Join("testdata", filename)
		binarypath := strings.Replace(sourcecodepath, ".ct33", ".machine", 1)

		expected, err := os.ReadFile(binarypath)
		if err != nil {
			continue // Not all .ct33 files have .machine files, so we skip those cases
			//t.Errorf("%s: unable to read machine code file: %v", filename, err)
		}

		buffer := &bytes.Buffer{}
		err = AssembleFile(sourcecodepath, buffer)
		if err != nil {
			panic(err)
		}
		got := buffer.String()

		expectedLines := strings.Split(string(expected), "\n")
		gottenLines := strings.Split(got, "\n")

		for i, expected := range expectedLines {
			if i >= len(gottenLines) {
				t.Errorf("%s: Expected %d lines but got %d lines", filename, len(expectedLines), len(gottenLines))
				break
			}

			got := gottenLines[i]
			if expected != got {
				t.Errorf("%s %d: Expected '%s' got '%s'", filename, i, expected, got)
				mexpected, _ := strconv.Atoi(expected)
				mgot, _ := strconv.Atoi(got)
				t.Errorf("%s => %s and %s => %s", expected, DisassembleInstruction(uint16(mexpected)), got, DisassembleInstruction(uint16(mgot)))
			}
		}

	}
}

// Testing the assembly of a single file
// primary interest is debug specific files which fail to compile
func TestSingleFileAssembly(t *testing.T) {
	err := AssembleFileWithOptions("testdata/isa.ct33", os.Stdout, SOURCE_CODE)
	if err != nil {
		panic(err)
	}
}

func ExampleAssembleFile() {
	err := AssembleFileWithOptions("testdata/isa.ct33", os.Stdout, SOURCE_CODE)
	if err != nil {
		panic(err)
	}

	// Output:
	// 1987; ADD  x9, x8, x7
	// 2987; SUB  x9, x8, x7
	// 3987; SUBI x9, x8, 7
	// 4987; LSH  x9, x8, 7
	// 5987; RSH  x9, x8, 7
	// 6910; BRZ  x9, pseudo
	// 7910; BGT  x9, pseudo
	// 8916; LD   x9, data
	// 9916; ST   x9, data
	// 0000; HLT
	// 8990; INP  x9
	// 9991; OUT  x9
	// 1908; MOV  x9, x8
	// 1900; CLR  x9
	// 3991; DEC  x9
	// 6010; BRA  pseudo
	// 0098; DAT  98
}

func ExampleInstruction_PrintSourceCode() {
	dat := SourceCodeLine{
		Instruction: &Instruction{
			opcode:   DAT,
			constant: 98,
		},
	}
	dat.Print(os.Stdout, SOURCE_CODE|MACHINE_CODE)

	subi := SourceCodeLine{
		Instruction: &Instruction{
			opcode:   SUBI,
			Regs:     []uint8{3, 4},
			constant: 8,
		},
	}

	subi.Print(os.Stdout, SOURCE_CODE|MACHINE_CODE)

	// Output:
	// 0098; DAT  98
	// 3348; SUBI x3, x4, 8
}
