package calcutron

import (
	"fmt"
	"os"
)

// Just to understand Go directory and file APIs better
// might want to throw away this test later
func Example_openDirectory() {
	dir, err := os.Open("testdata")
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	entries, err := dir.ReadDir(-1)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}

	// Output:
	// roundtrip.ct33
	// simplemult.ct33
	// labels-nocode.machine
	// maximizer.ct33
	// adder.ct33
	// adder.machine
	// doubler.machine
	// isa.machine
	// isa.ct33
	// simplemult.machine
	// doubler.ct33
	// labels-nocode.ct33
	// maximizer.machine
}
