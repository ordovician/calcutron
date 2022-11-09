package prog

import (
	"fmt"
	"io"
)

// Maps addresses to symbols
type AddressTable map[uint]string

type PrintContext struct {
	labels    SymbolTable
	addresses AddressTable
	options   *PrintOptions
}

func NewPrintContext(labels SymbolTable, options *PrintOptions) *PrintContext {
	addrToLabel := make(AddressTable)
	for label, addr := range labels {
		addrToLabel[addr] = label
	}

	return &PrintContext{
		labels:    labels,
		addresses: addrToLabel,
		options:   options,
	}
}

type AddressInstruction struct {
	Addr uint
	Inst Instruction
}

func (ctx *PrintContext) Print(writer io.Writer, input <-chan AddressInstruction) {
	addrToLabel := ctx.addresses
	options := ctx.options

	for addrInst := range input {
		addr := addrInst.Addr
		inst := addrInst.Inst

		if options.Address {
			AddressColor.Fprintf(writer, "%02d ", addr)
		} else if label, ok := addrToLabel[uint(addr)]; ok && options.SourceCode {
			LabelColor.Fprint(writer, label, ":")
			fmt.Fprintln(writer)
		}

		if options.MachineCode {
			GrayColor.Fprintf(writer, "%04d ", inst.MachineCode())
		} else {
			fmt.Fprintf(writer, "    ")
		}

		if options.SourceCode {
			fmt.Fprint(writer, inst.SourceCode())
		}
		fmt.Fprintln(writer)
	}
}
