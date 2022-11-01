package calcutron

import (
	"fmt"
	"testing"
)

func ExampleAssemblyFlag_TurnOn() {
	var (
		MachinecodeOn bool = true
		LineNoOn      bool = true
		SourceCodeOn  bool
		AddressOn     bool = true
		ColorOn       bool
	)

	var options AssemblyFlag

	options = options.TurnOn(MACHINE_CODE, MachinecodeOn)
	options = options.TurnOn(LINE_NO, LineNoOn)
	options = options.TurnOn(SOURCE_CODE, SourceCodeOn)
	options = options.TurnOn(ADDRESS, AddressOn)
	options = options.TurnOn(COLOR, ColorOn)

	fmt.Printf("%b\n", options)

	// Output:
	// 1011
}

func TestParseSetReg(t *testing.T) {
	reg, value, err := ParseSetReg("set x3 42")
	if err != nil {
		t.Errorf("Could not parse 'set x3 42' because %v", err)
	}

	if reg != 3 || value != 42 {
		t.Errorf("Register index and value should be 3 and 42 respectively, but we got %d and %d", reg, value)
	}
}
