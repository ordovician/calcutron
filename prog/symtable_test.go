package prog

import (
	"fmt"
	"os"
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
