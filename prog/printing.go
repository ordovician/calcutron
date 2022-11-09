package prog

import (
	"fmt"
	"io"
)

func printRegisterOperands(writer io.Writer, regs []uint) {
	for i, r := range regs {
		if i > 0 {
			fmt.Fprintf(writer, ", ")
		}
		fmt.Fprintf(writer, "x%d", r)
	}
}

func printMnemonic(writer io.Writer, opcode Opcode) {
	MnemonicColor.Fprintf(writer, "%-6v", opcode)
}
