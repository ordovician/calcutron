package calcutron

import (
	"fmt"
	"testing"
)

func TestAssembleLine(t *testing.T) {

}

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
		machinecode, _ := assembleLine(labels, line)
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
