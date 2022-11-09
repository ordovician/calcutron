package dbg

import (
	"os"
	"testing"

	"github.com/ordovician/calcutron/sim"
)

func TestRunning(t *testing.T) {
	var comp sim.Computer
	//comp.LoadFile("../examples/maximizer.ct33")
	comp.LoadFile("../maximizer.machine")
	RunCommand(os.Stdout, "inputs 1 2 3 4", &comp)
	RunCommand(os.Stdout, "run", &comp)
}

func TestReadline(t *testing.T) {
	var comp sim.Computer
	//comp.LoadFile("../examples/maximizer.ct33")
	comp.LoadFile("../maximizer.machine")
	RunCommand(os.Stdout, "inputs 1 2 3 4", &comp)
	//RunCommand(os.Stdout, "run", &comp)

	rl, err := CreateReadLine()
	if err != nil {
		return
	}

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}

		if len(line) == 0 {
			continue
		}

		RunCommand(os.Stdout, line, &comp)
	}
}
