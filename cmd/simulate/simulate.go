package main

import (
	"flag"
	"fmt"
	"os"

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

	flag.Parse()

	filepath := flag.Arg(0)

	var comp Computer
	comp.LoadFile(filepath)

	comp.LoadInputs(os.Stdin)

	comp.RunSteps(maxsteps)

	fmt.Printf("\nCPU Status\n%v\n", &comp)
}
