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

// Reads source code to determine address of labels.
// Labels starting with a dot '.' are treated as offsets
// from a base address. The base address is a non dot based label preceeding
// dot based addresses. The purpose of this is to be able to load
// a base address into a register and use immediate value offsets
// to get to specific addresses
func ReadSymTable(reader io.Reader) SymbolTable {
	scanner := bufio.NewScanner(reader)
	labels := make(SymbolTable)
	address := 0
	baseAddress := 0
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		n := len(line)

		if n == 0 {
			continue
		}

		if i := strings.IndexRune(line, ':'); i >= 0 {

			// check if we should record an offset or absolute address
			if strings.HasPrefix(line, ".") {
				address -= baseAddress
			} else {
				baseAddress = address
			}

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
