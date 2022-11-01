// WARNING: A lot of the code in this file is absolutely horrible. Don't use it
// as an inspiration for how to do anything. My focus was on hacking together something
// that would work with all desired features.
// The next step will be a major refactoring and cleanup process.
// Most likely I will create a struct with methods to represent each command, rather
// than maintaining a huge switch-case structure.
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
	files []string // source code files
}

var commands = [...]string{"help", "input", "output", "status", "quit", "set", "next", "run", "list"}

// Returns suggestions based on what user has written thus far
func (completer *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	str := string(line)
	n := len(line)
	var matches [][]rune = make([][]rune, 0, 5)

	if strings.HasPrefix(str, "load") && len(str) > len("load") {
		filearg := strings.TrimSpace(str[len("load"):])
		if len(filearg) == 0 {
			for _, file := range completer.files {
				matches = append(matches, []rune(file))
			}
		}

		for _, file := range completer.files {
			if strings.HasPrefix(file, filearg) {
				match := file[n-len("load")-1:]
				matches = append(matches, []rune(match))
			}
		}
	} else if strings.HasPrefix("load", str) {
		matches = append(matches, []rune("load"[n:]))
	}

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

func getSourceCodeFiles() []string {
	files := make([]string, 0)
	entries, _ := os.ReadDir(".")

	for _, entry := range entries {
		filename := entry.Name()
		if strings.HasSuffix(filename, ".machine") {
			files = append(files, filename)
		}
	}
	return files
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of debugger:\n")
		flag.PrintDefaults()
	}

	var options AssemblyFlag
	options = options.Set(MACHINE_CODE)
	options = options.Set(SOURCE_CODE)
	options = options.Set(COLOR)

	ParseAssemblyOptions(&options)

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(-1)
	}

	filepath := flag.Arg(0)

	//err := AssembleFileWithOptions(filepath, os.Stdout, options)
	var completer Completer
	completer.files = getSourceCodeFiles()

	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	numberColor := NumberColor.SprintFunc()

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
				fmt.Print("Outputs: ")
				JoinFunc(os.Stdout, computer.Outputs, ", ", numberColor)
				fmt.Print("\n\n")
				continue
			case strings.HasPrefix(line, "status"):
				computer.Print(os.Stdout, options.Has(COLOR))
				fmt.Println()
				continue
			case strings.HasPrefix(line, "x"):
				i, err := strconv.Atoi(line[1:])
				if err != nil {
					fmt.Fprintf(os.Stderr, "unable to parse index %s because %v", line[1:], err)
				}
				if i < 0 || i > 9 {
					fmt.Fprintf(os.Stderr, "x0 to x9 are the only valid registers, not x%d\n", i)
				}
				fmt.Printf("x%d: ", i)
				NumberColor.Printf("%02d\n\n", computer.Registers[i])
				continue
			case strings.HasPrefix(line, "quit"):
				goto exit
			case strings.HasPrefix(line, "set"):
				reg, value, err := ParseSetReg(line)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
				computer.Registers[reg] = value
				continue
			case line == "next":
				_ = computer.PrintCurrentInstruction(options | ADDRESS)
				machinecode := computer.Memory[computer.PC]
				instruction := DisassembleInstruction(machinecode)
				computer.Step()

				// Look at impact on registers from running last instruction by printing values
				// of register operands
				if instruction != nil {
					switch instruction.Opcode() {
					case BRA, BRZ, BGT:
						computer.PrintPC(os.Stdout, true)
					case HLT:
						break
					default:
						computer.PrintRegs(os.Stdout, true, instruction.UniqueRegisters()...)
					}
				}
				fmt.Println()
				continue
			case line == "run":
				// TODO: Read in max number of step
				computer.RunStepsWithOptions(40, options|ADDRESS)
				fmt.Println()
				continue
			case strings.HasPrefix(line, "load"):
				filename := strings.TrimSpace(line[4:])
				err := computer.LoadFile(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
				continue
			case line == "list":
				for i := 0; i < 90; i++ {
					machinecode := computer.Memory[i]
					// No need to disassemble numerous consecutive zeroes
					if machinecode == 0 && i > 0 && computer.Memory[i-1] == 0 {
						break
					}
					instruction := DisassembleInstruction(machinecode)
					if instruction != nil {
						codeline := SourceCodeLine{
							Instruction: instruction,
							Address:     i,
						}
						err := codeline.Print(os.Stdout, options|ADDRESS)
						if err != nil {
							fmt.Fprintf(os.Stderr, "%v\n", err)
						}
					}
				}
				fmt.Println()
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
				err := line.Print(os.Stdout, options)
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
