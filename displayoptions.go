// options for what should be displayed in output of assemble and disassemble output
package calcutron

import "flag"

// Setup what flags exist on command line
func ParseAssemblyOptions(options *AssemblyFlag) {
	var (
		AddressOn     bool
		LineNoOn      bool
		MachinecodeOn bool
		SourceCodeOn  bool
		ColorOn       bool
	)

	flag.BoolVar(&AddressOn, "address", options.Has(ADDRESS), "show address of each instruction")
	flag.BoolVar(&LineNoOn, "lineno", options.Has(LINE_NO), "show source code line number of each instruction")
	flag.BoolVar(&MachinecodeOn, "machinecode", options.Has(MACHINE_CODE), "show address of each instruction")
	flag.BoolVar(&SourceCodeOn, "sourcecode", options.Has(SOURCE_CODE), "show address of each instruction")
	flag.BoolVar(&ColorOn, "color", options.Has(COLOR), "colorize output")

	flag.Parse()

	*options = options.TurnOn(ADDRESS, AddressOn)
	*options = options.TurnOn(LINE_NO, LineNoOn)
	*options = options.TurnOn(MACHINE_CODE, MachinecodeOn)
	*options = options.TurnOn(SOURCE_CODE, SourceCodeOn)
	*options = options.TurnOn(COLOR, ColorOn)
}
