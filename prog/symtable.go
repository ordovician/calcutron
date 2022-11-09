package prog

import (
	"bufio"
	"io"
	"strings"
)

type SymbolTable map[string]uint

// Add labels such as input and output to make it easier to write code
// read and writing to output and input
func (labels SymbolTable) AddIOLabels() {
	labels["tape"] = 99

}

// A table containing the memory address of labels in the code
func ReadSymTable(reader io.Reader) SymbolTable {
	scanner := bufio.NewScanner(reader)
	labels := make(SymbolTable)
	address := 0
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		n := len(line)

		if n == 0 {
			continue
		}

		if i := strings.IndexRune(line, ':'); i >= 0 {
			labels[line[0:i]] = uint(address)

			// is there anything beyond the label?
			if n == i+1 {
				continue
			}
		}
		address++
	}
	return labels
}
