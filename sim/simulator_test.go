package sim

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/ordovician/calcutron/asm"
	"github.com/ordovician/calcutron/prog"
)

// for testing if a given register has given value
type RegisterValues struct {
	reg uint
	val int
}

func TestSimpleAdder(t *testing.T) {
	var labels prog.SymbolTable
	lines := []string{
		"MOVE  x1, 3",
		"MOVE  x2, 4",
		"ADD   x3, x1, x2",
		"SUB   x4, x2, x1",
		"SUB   x5, x1, x2",
		"MOVE  x6, 2",
		"SHFT x7, x6, 1",
		"MOVE  x8, 48",
		"SHFT x9, x8, -1",
		"JMP x0, 0",
	}

	var comp Computer

	for addr, line := range lines {
		inst, err := asm.AssembleLine(labels, line, uint(addr))
		if err != nil {
			t.Errorf("failed to assemble '%s' because %v", line, err)
			return
		}
		inst.Run(&comp)
	}

	data := []RegisterValues{
		{1, 3},
		{2, 4},
		{3, 7},
		{4, 1},
		{5, -1},
		{6, 20},
		{7, 0},
		{8, 4},
		{9, 8},
	}

	for _, expected := range data {
		if comp.Register(expected.reg) != expected.val {
			t.Errorf("Expected %d got %d", expected.val, comp.Register(expected.reg))
		}
	}

}

func TestCounter(t *testing.T) {
	sourceCode := `
		MOVE x2, 4
		MOVE x1, 1
	loop:
		ADD   x3, x3, x1
		BGT   x2, x3, loop
		HLT		
	`

	buffer := bytes.NewReader([]byte(sourceCode))

	program, err := asm.Assemble(buffer)
	if err != nil {
		t.Errorf("Failed to assemble  because %v", err)
		return
	}

	comp := NewComputer(program)
	comp.Run(20)

	data := []RegisterValues{
		{1, 1},
		{2, 4},
		{3, 4},
	}

	for _, expected := range data {
		if comp.Register(expected.reg) != expected.val {
			t.Errorf("Expected %d got %d", expected.val, comp.Register(expected.reg))
		}
	}

}

func Example_maximizer() {
	comp, err := NewComputerFile("../Examples/maximizer.ct33")
	comp.inputs = []uint{2, 3, 8, 4}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	comp.Run(50)
	comp.Print(os.Stdout)
	// Output:
	// PC: 00    Steps: 50
	//
	// x1: 0008, x4: 0000, x7: 0000
	// x2: 0004, x5: 0000, x8: 0000
	// x3: 0000, x6: 0000, x9: 0000
	//
	// Inputs:  2, 3, 8, 4
	// Outputs: 3, 8
}

func Example_doubler() {
	comp, err := NewComputerFile("../Examples/doubler.ct33")
	comp.inputs = []uint{2, 3, 8, 4}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	comp.Run(50)
	comp.Print(os.Stdout)
	// Output:
	// PC: 00
	//
	// x1: 0004, x4: 0000, x7: 0000
	// x2: 0000, x5: 0000, x8: 0000
	// x3: 0008, x6: 0000, x9: 0000
	//
	// Inputs:  2, 3, 8, 4
	// Outputs: 4, 6, 16, 8
}

func Example_simpleMult() {
	comp, err := NewComputerFile("../Examples/simplemult.ct33")
	comp.inputs = []uint{2, 3, 8, 4}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	comp.Run(50)
	comp.Print(os.Stdout)
	// Output:
	// PC: 00
	//
	// x1: 0008, x4: 0000, x7: 0000
	// x2: 0000, x5: 0000, x8: 0000
	// x3: 0032, x6: 0000, x9: 0000
	//
	// Inputs:  2, 3, 8, 4
	// Outputs: 6, 32
}
