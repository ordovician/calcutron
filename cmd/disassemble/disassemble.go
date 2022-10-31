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
	var options AssemblyFlag
	options = options.Set(SOURCE_CODE)

	ParseAssemblyOptions(&options)

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(-1)
	}

	filepath := flag.Arg(0)

	err := DisassembleFileWithOptions(filepath, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to disassemble %s: %v", filepath, err)
	}

}
