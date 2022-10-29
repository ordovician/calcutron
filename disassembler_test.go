package calcutron

import (
	"fmt"
	"testing"
)

func ExampleInstruction_String() {
	instructions := []uint16{1234, 8213, 9999, 1010, 0000}

	for _, inst := range instructions {
		instruction := DisassembleInstruction(inst)
		fmt.Println(instruction)
	}

	// Output:
	// ADD  x2, x3, x4
	// LD   x2, 13
	// ST   x9, 99
	// ADD  x0, x1, x0
	// HLT
}

func TestLoadInstruction(t *testing.T) {
	var comp Computer
	comp.Inputs = []uint8{2, 3, 8, 4}

	comp.ExecuteInstruction(8190) // LD   x1, 90
	comp.ExecuteInstruction(8290) // LD   x2, 90

	if comp.Registers[1] != comp.Inputs[0] {
		t.Errorf("Register x1 value %d not equal expected value %d", comp.Registers[1], comp.Inputs[0])
	}

	if comp.Registers[2] != comp.Inputs[1] {
		t.Errorf("Register x1 value %d not equal expected value %d", comp.Registers[1], comp.Inputs[1])
	}

	comp.ExecuteInstruction(1112) // ADD  x1, x1, x2

	expected := comp.Inputs[0] + comp.Inputs[1]
	if comp.Registers[1] != expected {
		t.Errorf("Register x1 value %d not equal expected value %d", comp.Registers[1], expected)
	}
}

// Make sure disassembly and assembly works roundtrip for a single instruction
func TestDisassembleInstruction(t *testing.T) {
	for machinecode := 1000; machinecode <= 9999; machinecode++ {
		inst := DisassembleInstruction(uint16(machinecode))
		got, err := inst.Machinecode()
		if err != nil {
			t.Errorf("Unable to encode machinecode instruction: %v", err)
		}
		if got != uint16(machinecode) {
			t.Errorf("Expected %d got %d", machinecode, got)
		}

	}
}
