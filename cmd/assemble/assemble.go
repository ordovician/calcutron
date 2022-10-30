package main

import (
	"flag"
	"fmt"
	"os"

	. "github.com/ordovician/calcutron"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of assemble:\n")
		flag.PrintDefaults()
	}

	var (
		AddressOn     bool
		LineNoOn      bool
		MachinecodeOn bool
		SourceCodeOn  bool
		ColorOn       bool
	)

	flag.BoolVar(&AddressOn, "address", false, "show address of each instruction")
	flag.BoolVar(&LineNoOn, "lineno", false, "show source code line number of each instruction")
	flag.BoolVar(&MachinecodeOn, "machinecode", true, "show address of each instruction")
	flag.BoolVar(&SourceCodeOn, "sourcecode", false, "show address of each instruction")
	flag.BoolVar(&ColorOn, "color", false, "colorize output")

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(-1)
	}

	filepath := flag.Arg(0)

	var options AssemblyFlag
	options = options.TurnOn(ADDRESS, AddressOn)
	options = options.TurnOn(LINE_NO, LineNoOn)
	options = options.TurnOn(MACHINE_CODE, MachinecodeOn)
	options = options.TurnOn(SOURCE_CODE, SourceCodeOn)
	options = options.TurnOn(COLOR, ColorOn)

	err := AssembleFileWithOptions(filepath, os.Stdout, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to assemble %s: %v", filepath, err)
	}

}
