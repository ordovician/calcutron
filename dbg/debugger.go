package dbg

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"

	"github.com/ordovician/calcutron/prog"
	"github.com/ordovician/calcutron/sim"
)

type Completer struct {
	files []string // source code files
}

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
		cmdName := cmd.Name()
		if strings.HasPrefix(cmdName, str) {
			matches = append(matches, []rune(cmdName[n:]))
		}
	}

	for _, opstr := range prog.AllOpcodeStrings {
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

func RunCommand(writer io.Writer, line string, comp *sim.Computer) error {
	args := strings.Fields(line)
	cmd := Lookup(args[0])
	if cmd == nil {
		return fmt.Errorf("no command named '%s' is supported. write 'help' to get list of supported commands", args[0])
	}
	err := cmd.Action(writer, comp, args[1:])
	if err != nil {
		return err
	}
	fmt.Fprintln(writer)
	return nil
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

func CreateReadLine() (*readline.Instance, error) {
	var completer Completer
	completer.files = getSourceCodeFiles()

	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:            green("caluctron> "),
		HistoryFile:       "/tmp/readline.tmp",
		AutoComplete:      &completer,
		HistorySearchFold: true,
	})

	return rl, err
}
