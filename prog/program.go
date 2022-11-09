package prog

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type PrintOptions struct {
	LineNo      bool
	Address     bool
	MachineCode bool
	SourceCode  bool
}

type Program struct {
	Labels       SymbolTable
	Instructions []Instruction
}

func (prog *Program) Add(inst Instruction) {
	prog.Instructions = append(prog.Instructions, inst)
}

func (prog *Program) Print(writer io.Writer) {
	for _, inst := range prog.Instructions {
		fmt.Fprintln(writer, inst.SourceCode())
	}
}

func (prog *Program) PrintWithOptions(writer io.Writer, options *PrintOptions) {

	ctx := NewPrintContext(prog.Labels, options)
	channel := make(chan AddressInstruction)

	var group sync.WaitGroup
	group.Add(1)
	go func() {
		go ctx.Print(writer, channel)
		group.Done()
	}()

	for addr, inst := range prog.Instructions {
		channel <- AddressInstruction{
			Addr: uint(addr),
			Inst: inst,
		}
	}
	close(channel)
	group.Wait()
}

func (prog *Program) String() string {
	var buffer bytes.Buffer
	prog.Print(&buffer)
	return buffer.String()
}
