package main

import (
	"flag"
	"fmt"
	"os"

	. "github.com/ordovician/calcutron"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of disassemble:\n")
		flag.PrintDefaults()
	}

	// var (
	// 	maxsteps int
	// 	inputs   string
	// )

	// flag.IntVar(&maxsteps, "maxsteps", 1000, "Max number of instruction to execute")
	// flag.StringVar(&inputs, "inputs", "", "Input numbers for program to read")

	flag.Parse()

	filepath := flag.Arg(0)
	err := DisassembleFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to disassemble %s: %v", filepath, err)
	}

}
