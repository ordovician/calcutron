package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	. "github.com/ordovician/calcutron"
)

type Completer struct {
}

// Returns suggestions based on what user has written thus far
func (completer *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	str := string(line)
	n := len(line)
	var matches [][]rune = make([][]rune, 0, 5)

	for _, opstr := range AllOpcodeStrings {
		if strings.HasPrefix(opstr, str) {
			matches = append(matches, []rune(opstr[n:]))
		}

		// Allow askin for help on a command
		helpStr := "?" + opstr
		if strings.HasPrefix(helpStr, str) {
			matches = append(matches, []rune(helpStr[n:]))
		}
	}

	return matches, n
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of debugger:\n")
		flag.PrintDefaults()
	}

	var options AssemblyFlag
	options = options.Set(MACHINE_CODE)
	options = options.Set(COLOR)

	ParseAssemblyOptions(&options)

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(-1)
	}

	// filepath := flag.Arg(0)

	//err := AssembleFileWithOptions(filepath, os.Stdout, options)

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	var completer Completer
	rl.Config.AutoComplete = &completer

	labels := make(map[string]uint8)
	//var computer Computer

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}

		if len(line) == 0 {
			continue
		}

		if line == "?" {
			fmt.Println(strings.Join(AllOpcodeStrings, " "))
		}

		if line[0] != '?' {
			instruction, err := AssembleLine(labels, line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to assemble instruction: %v\n", err)
			}

			machinecode, err := instruction.Machinecode()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			instruction = DisassembleInstruction(machinecode)

			if instruction != nil {
				line := SourceCodeLine{
					Instruction: instruction,
					Address:     0,
				}
				err := line.Print(os.Stdout, options|SOURCE_CODE)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
			fmt.Println()
		}
	}

	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to debug %s: %v", filepath, err)
	// }

}
