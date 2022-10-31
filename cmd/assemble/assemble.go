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

	var options AssemblyFlag
	options = options.Set(MACHINE_CODE)

	ParseAssemblyOptions(&options)

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(-1)
	}

	filepath := flag.Arg(0)

	err := AssembleFileWithOptions(filepath, os.Stdout, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to assemble %s: %v", filepath, err)
	}

}
