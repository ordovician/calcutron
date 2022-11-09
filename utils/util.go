package utils

import (
	"fmt"
	"io"
	"strconv"
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

func JoinFunc[T any](writer io.Writer, elements []T, sep string, fn func(...any) string) {
	if len(elements) == 0 {
		return
	}

	fmt.Fprintf(writer, "%v", fn(elements[0]))

	for _, elem := range elements[1:] {
		fmt.Fprintf(writer, "%s%v", sep, fn(elem))
	}
}

func RemoveDuplicates[T constraints.Ordered](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
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

// Parse a string such as 'set x2 42'. It should return that index of register is 2
// and value is 42
func ParseSetReg(line string) (reg, value uint8, erro error) {
	fields := strings.Fields(line)

	if len(fields) != 3 {
		erro = fmt.Errorf("expected command on form 'set x2 42' but got %s", line)
		return
	}

	regstr := fields[1]

	if !strings.HasPrefix(regstr, "x") {
		erro = fmt.Errorf("registers must start with an 'x' you start it with '%c'", regstr[1])
		return
	}

	i, err := strconv.Atoi(regstr[1:])
	if err != nil {
		erro = fmt.Errorf("unable to parse register index %s because %w", regstr[1:], err)
		return
	}
	if i < 0 || i > 9 {
		erro = fmt.Errorf("x0 to x9 are the only valid registers, not x%d", i)
		return
	}
	reg = uint8(i)

	valuestr := fields[2]
	val, err := strconv.Atoi(valuestr)
	if err != nil {
		erro = fmt.Errorf("unable to parse register value %s because %w", valuestr, err)
		return
	}
	value = uint8(val)

	return
}
