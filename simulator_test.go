package calcutron

import "testing"

func ExampleComputer_RunSteps() {
	// adder program
	instructions := [99]uint16{8190, 8290, 1112, 9191, 6000}

	comp := Computer{
		Memory: instructions,
		Inputs: []uint8{2, 3, 8, 4},
	}

	comp.RunSteps(40)

	// Output:
	// 00: 8190; LD   x1, 90
	// 01: 8290; LD   x2, 90
	// 02: 1112; ADD  x1, x1, x2
	// 03: 9191; ST   x1, 91
	// 04: 6000; BRZ  x0, 0
	// 00: 8190; LD   x1, 90
	// 01: 8290; LD   x2, 90
	// 02: 1112; ADD  x1, x1, x2
	// 03: 9191; ST   x1, 91
	// 04: 6000; BRZ  x0, 0
	// 00: 8190; LD   x1, 90
}

func ExampleComputer_LoadFile() {
	var comp Computer
	comp.LoadFile("testdata/adder.machine")
	comp.Inputs = []uint8{2, 3, 8, 4}
	comp.RunSteps(40)

	// Output:
	// 00: 8190; LD   x1, 90
	// 01: 8290; LD   x2, 90
	// 02: 1112; ADD  x1, x1, x2
	// 03: 9191; ST   x1, 91
	// 04: 6000; BRZ  x0, 0
	// 00: 8190; LD   x1, 90
	// 01: 8290; LD   x2, 90
	// 02: 1112; ADD  x1, x1, x2
	// 03: 9191; ST   x1, 91
	// 04: 6000; BRZ  x0, 0
	// 00: 8190; LD   x1, 90
}

func TestAdder(t *testing.T) {
	var comp Computer
	comp.LoadFile("testdata/adder.machine")
	comp.Inputs = []uint8{2, 3, 8, 4}
	comp.RunSteps(40)

}
