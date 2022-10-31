package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	. "github.com/ordovician/calcutron"
)

type Completer struct {
}

var commands = [...]string{"help", "input", "output", "status", "quit", "set", "next", "run"}

// Returns suggestions based on what user has written thus far
func (completer *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	str := string(line)
	n := len(line)
	var matches [][]rune = make([][]rune, 0, 5)

	for _, cmd := range commands {
		if strings.HasPrefix(cmd, str) {
			matches = append(matches, []rune(cmd[n:]))
		}
	}

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

	filepath := flag.Arg(0)

	//err := AssembleFileWithOptions(filepath, os.Stdout, options)
	var completer Completer

	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	numberColor := color.New(color.FgHiRed).SprintFunc()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:            green("caluctron> "),
		HistoryFile:       "/tmp/readline.tmp",
		AutoComplete:      &completer,
		HistorySearchFold: true,
	})

	if err != nil {
		panic(err)
	}
	defer rl.Close()

	labels := make(map[string]uint8)
	var computer Computer
	computer.LoadFile(filepath)

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
			switch {
			case strings.HasPrefix(line, "help"):
				fmt.Println(strings.Join(commands[:], " "))
				continue
			case strings.HasPrefix(line, "input"):
				computer.StringInputs(line[5:])
				continue
			case line == "output":
				// TODO: Figure out why this doesn't work
				JoinFunc(os.Stdout, computer.Outputs, " ", numberColor)
				continue
			case strings.HasPrefix(line, "status"):
				computer.Print(os.Stdout, options.Has(COLOR))
				fmt.Println()
				continue
			case strings.HasPrefix(line, "quit"):
				goto exit
			case strings.HasPrefix(line, "set"):
				// TODO: Figure out register and what value it is set to
				args := strings.TrimSpace(line[3:])
				if len(args) > 0 && args[0] == 'x' {

					i, err := strconv.Atoi(args[1:])
					if err != nil {
						fmt.Fprintf(os.Stderr, "unable to parse index %s because %v", args[1:], err)
					}
					if i < 0 || i > 9 {
						fmt.Fprintf(os.Stderr, "x0 to x9 are the only valid registers, not x%d\n", i)
					}

				}
				continue
			case line == "next":
				_ = computer.PrintCurrentInstruction(options)
				computer.Step()
				continue
			case line == "run":
				// TODO: Read in max number of step
				computer.RunStepsWithOptions(40, options)
				continue
			default:
				break
			}

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

			err = computer.ExecuteInstruction(machinecode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to execute instructon: %v", err)
			}

			fmt.Println()
		}
	}
exit:
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to debug %s: %v", filepath, err)
	// }

}
