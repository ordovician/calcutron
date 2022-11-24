package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ordovician/calcutron/asm"
	"github.com/ordovician/calcutron/dbg"
	"github.com/ordovician/calcutron/disasm"
	"github.com/ordovician/calcutron/prog"
	"github.com/ordovician/calcutron/sim"
	// "github.com/ordovician/calcutron/utils"
	"github.com/urfave/cli/v2"
)

var printOptions prog.PrintOptions

func assemble(ctx *cli.Context) error {
	filepath := ctx.Args().First()
	program, err := asm.AssembleFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return nil
	}

	program.PrintWithOptions(os.Stdout, &printOptions)
	return nil
}

func disassemble(ctx *cli.Context) error {
	filepath := ctx.Args().First()
	program, err := disasm.DisassembleFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return nil
	}

	program.PrintWithOptions(os.Stdout, &printOptions)
	return nil
}

func runCode(ctx *cli.Context) error {
	filepath := ctx.Args().First()
	verbose := ctx.Bool("verbose")

	var program *prog.Program
	var err error
	if strings.HasSuffix(filepath, ".ct33") {
		program, err = asm.AssembleFile(filepath)
	} else if strings.HasSuffix(filepath, ".machine") {
		program, err = disasm.DisassembleFile(filepath)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return nil
	}

	comp := sim.NewComputer(program)
	comp.LoadInputs(os.Stdin)

	if verbose {
		pContext := prog.NewPrintContext(program.Labels, &printOptions)

		var group sync.WaitGroup
		group.Add(1)
		channel := make(chan prog.AddressInstruction)
		go func() {
			pContext.Print(os.Stdout, channel)
			group.Done()
		}()

		// NOTE: I am arbitrarily restricting program from running more than 5000 instructions
		// that is just to not end up in some endless loop
		comp.RunChannel(5000, channel)
		close(channel)

		// wait until executed instuctions have been printed out to consol
		group.Wait()

		// Print out status of computer
		fmt.Println()
		fmt.Println(comp.String())
	} else {
		comp.Run(5000)
		for _, n := range comp.Outputs() {
			fmt.Println(n)
		}
	}

	return nil
}

func debug(ctx *cli.Context) error {
	errorColor := color.New(color.FgRed)

	var comp sim.Computer
	args := ctx.Args()
	if args.Len() > 0 {
		comp.LoadFile(args.First())
	}

	rl, err := dbg.CreateReadLine()
	if err != nil {
		return err
	}

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			// errorColor.Fprintf(os.Stderr, "Error: ")
			// fmt.Fprintf(os.Stderr, "%v\n\n", err)
			break
		}

		if len(line) == 0 {
			continue
		}

		err = dbg.RunCommand(os.Stdout, line, &comp)
		if err != nil {
			errorColor.Fprintf(os.Stderr, "Error: ")
			fmt.Fprintf(os.Stderr, "%v\n\n", err)
		}
	}

	return nil
}

type CommandType int

const (
	ASSEMBLY CommandType = iota
	DISASSEMBLY
	SIMULATION
)

func createFlags(cmdType CommandType) []cli.Flag {

	addressFlag := cli.BoolFlag{
		Name:        "address",
		Usage:       "show memory address of instruction",
		Value:       cmdType == SIMULATION,
		Destination: &printOptions.Address,
	}

	sourceCodeFlag := cli.BoolFlag{
		Name:        "sourcecode",
		Usage:       "show source code for each instruction",
		Value:       cmdType == DISASSEMBLY || cmdType == SIMULATION,
		Destination: &printOptions.SourceCode,
	}

	machineCodeFlag := cli.BoolFlag{
		Name:        "machinecode",
		Usage:       "show machine code for each instruction",
		Value:       cmdType == ASSEMBLY || cmdType == SIMULATION,
		Destination: &printOptions.MachineCode,
	}

	lineNoFlag := cli.BoolFlag{
		Name:        "lineno",
		Usage:       "show line number for each instruction",
		Destination: &printOptions.LineNo,
	}

	return []cli.Flag{
		&addressFlag,
		&sourceCodeFlag,
		&machineCodeFlag,
		&lineNoFlag,
	}
}

func main() {

	disassembleCmd := cli.Command{
		Name:    "disassemble",
		Aliases: []string{"disasm"},
		Usage:   "disassemble a file containing calcutron-33 machine code",
		Action:  disassemble,
		Flags:   createFlags(DISASSEMBLY),
	}

	assembleCmd := cli.Command{
		Name:    "assemble",
		Aliases: []string{"asm"},
		Usage:   "assemble a calcutron-33 assembly code file",
		Action:  assemble,
		Flags:   createFlags(ASSEMBLY),
	}

	runFlags := createFlags(SIMULATION)
	verboseFlag := cli.BoolFlag{
		Name:  "verbose",
		Usage: "show each instruction executed",
	}
	runFlags = append(runFlags, &verboseFlag)

	runCmd := cli.Command{
		Name:    "run",
		Aliases: []string{"simulate", "sim"},
		Usage:   "run a calcutron-33 machine code file",
		Action:  runCode,
		Flags:   runFlags,
	}

	dbgCmd := cli.Command{
		Name:    "debug",
		Aliases: []string{"dbg"},
		Usage:   "debug calcutron program",
		Action:  debug,
	}

	app := &cli.App{
		Usage: "Tool to assemble, disassemble and run Calcutron-33 assembly code",
		Commands: []*cli.Command{
			&assembleCmd,
			&disassembleCmd,
			&runCmd,
			&dbgCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
