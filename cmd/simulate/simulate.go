package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	. "github.com/ordovician/calcutron"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of simulate:\n")
		flag.PrintDefaults()
	}

	var (
		maxsteps int
		inputs   string
	)

	flag.IntVar(&maxsteps, "maxsteps", 1000, "Max number of instruction to execute")
	flag.StringVar(&inputs, "inputs", "", "Input numbers for program to read")

	var options AssemblyFlag
	options = options.Set(MACHINE_CODE)
	options = options.Set(SOURCE_CODE)
	options = options.Set(ADDRESS)

	ParseAssemblyOptions(&options)

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(-1)
	}

	filepath := flag.Arg(0)

	var comp Computer
	comp.LoadFile(filepath)

	comp.LoadInputs(os.Stdin)

	comp.RunStepsWithOptions(maxsteps, options)

	if options.Has(COLOR) {
		white := color.New(color.FgWhite, color.Bold)
		white.Printf("\nCPU Status\n")
	} else {
		fmt.Printf("\nCPU Status\n")
	}
	comp.Print(os.Stdout, options.Has(COLOR))
	fmt.Println()
}
