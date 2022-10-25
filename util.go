package calcutron

import (
	"fmt"
	"io"
	"strings"
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
func complement(x, max int8) uint8 {
	if x < 0 {
		return uint8(max + x)
	} else if x == 0 {
		return 0
	} else {
		return uint8(x % max)
	}
}
