package calcutron

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"golang.org/x/exp/constraints"
)

// SplitAt returns a bufio.SplitFunc closure, splitting at a substring
// scanner.Split(SplitAt("\n# "))
func SplitAt(substring string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		// Return nothing if at end of file and no data passed
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Find the index of the input of the separator substring
		if i := strings.Index(string(data), substring); i >= 0 {
			return i + len(substring), data[0:i], nil
		}

		// If at end of file with data return the data
		if atEOF {
			return len(data), data, nil
		}

		return
	}
}

// Joins elements with separator sep and write output to writer
func Join[T any](writer io.Writer, elements []T, sep string) {
	if len(elements) == 0 {
		return
	}

	fmt.Fprintf(writer, "%v", elements[0])

	for _, elem := range elements[1:] {
		fmt.Fprintf(writer, "%s%v", sep, elem)
	}
}

// Calculate complement of x
// to calculate 10s complement of x, you would call complement(x, 10)
// this would return a value between 0 and 9
func complement(x, max int8) uint16 {
	if x < 0 {
		return uint16(max + x)
	} else if x == 0 {
		return 0
	} else {
		return uint16(x % max)
	}
}

// Check if all runes in string are digits
func AllDigits(s string) bool {
	for _, ch := range s {
		if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}

func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

type AssemblyFlag uint8

// controls output of Assemble function
const (
	MACHINE_CODE AssemblyFlag = 1 << iota
	LINE_NO                   // show line number
	SOURCE_CODE               // show source code
	ADDRESS                   // show address of machine code instruction
	COLOR                     // Colorize output
)

// Set bit for flag
func (flag AssemblyFlag) Set(b AssemblyFlag) AssemblyFlag { return flag | b }

func (flag AssemblyFlag) TurnOn(b AssemblyFlag, on bool) AssemblyFlag {
	if on {
		return flag.Set(b)
	} else {
		return flag.Clear(b)
	}
}

// Clear bit for flag
func (flag AssemblyFlag) Clear(b AssemblyFlag) AssemblyFlag { return flag &^ b }

// Toggle bit for flat
func (flag AssemblyFlag) Toggle(b AssemblyFlag) AssemblyFlag { return flag ^ b }

// Check if bit is set
func (flag AssemblyFlag) Has(b AssemblyFlag) bool { return flag&b != 0 }

type SourceCodeLine struct {
	address     int
	machinecode int16
	sourcecode  string
	lineno      int
}

func PrintInstruction(writer io.Writer, line SourceCodeLine, options AssemblyFlag) {

	if options.Has(ADDRESS) {
		fmt.Fprintf(writer, "%02d: ", line.address)
	}

	fmt.Fprintf(writer, "%04d", line.machinecode)

	if options.Has(SOURCE_CODE) {
		if options.Has(LINE_NO) {
			fmt.Fprintf(writer, "; %-18s", line.sourcecode)
		} else {
			fmt.Fprintf(writer, "; %s", line.sourcecode)
		}
	}

	if options.Has(LINE_NO) {
		fmt.Fprintf(writer, " // Line %2d ", line.lineno)
	}

	fmt.Fprintln(writer)
}
