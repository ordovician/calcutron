package calcutron

import (
	"fmt"
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
