package calcutron

import (
	"fmt"
)

func ExampleInstruction_String() {
	instructions := []uint16{1234, 8213, 9999, 1010, 0000}

	for _, inst := range instructions {
		var instruction *Instruction = decodeInstruction(inst)
		fmt.Println(instruction)
	}

	// Output:
	// ADD  x2, x3, x4
	// LD   x2, 13
	// ST   x9, 99
	// ADD  x0, x1, x0
	// HLT
}

func Example_decodeInstruction() {
	instructions := []uint16{1234, 8213, 9999, 0000, 1010}

	for _, inst := range instructions {
		var instruction *Instruction = decodeInstruction(inst)

		fmt.Println("Decoding instruction: ", inst)
		fmt.Println("  Opcode: ", instruction.opcode)
		fmt.Println("  Dest:   ", instruction.dst)
		fmt.Println("  Source: ", instruction.src)
		fmt.Println("  Address:", instruction.addr)
		fmt.Println("  Offset: ", instruction.offset)
	}

	// Output:
	// Decoding instruction:  1234
	//   Opcode:  ADD
	//   Dest:    2
	//   Source:  3
	//   Address: 34
	//   Offset:  4
	// Decoding instruction:  8213
	//   Opcode:  LD
	//   Dest:    2
	//   Source:  1
	//   Address: 13
	//   Offset:  3
	// Decoding instruction:  9999
	//   Opcode:  ST
	//   Dest:    9
	//   Source:  9
	//   Address: 99
	//   Offset:  9
	// Decoding instruction:  0
	//   Opcode:  HLT
	//   Dest:    0
	//   Source:  0
	//   Address: 0
	//   Offset:  0
	// Decoding instruction:  1010
	//   Opcode:  ADD
	//   Dest:    0
	//   Source:  1
	//   Address: 10
	//   Offset:  0
}
