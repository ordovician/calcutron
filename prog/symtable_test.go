package prog

import (
	"fmt"
	"os"
	"testing"
)

func Example_readSymTable() {
	file, _ := os.Open("testdata/labels-nocode.ct33")
	defer file.Close()

	labels := ReadSymTable(file)

	for key, value := range labels {
		fmt.Println(value, key)
	}

	// Unordered Output:
	// 0 epsilon
	// 0 alpha
	// 0 beta
	// 0 gamma
	// 0 delta
}

func TestReadMemAddSymbols(t *testing.T) {
	file, _ := os.Open("../examples/memcalc.ct33")
	defer file.Close()

	labels := ReadSymTable(file)

	for key, value := range labels {
		fmt.Println(value, key)
	}

	checkSymbol := func(label string, expected uint) {
		got := labels[label]
		if expected != got {
			t.Errorf("label '%s' expected %d got %d", label, expected, got)
		}
	}

	checkSymbol("first", 3)
	checkSymbol("second", 4)
	checkSymbol("added", 5)
	checkSymbol("shiftresult", 1)
	checkSymbol(".shifted", 0)
	checkSymbol(".overflow", 1)
}
