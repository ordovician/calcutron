package prog

import "golang.org/x/exp/constraints"

// Convert an unsigned integer to a Signed one using 10s complement
// Numbers in range 0 to 9 gets converte to range -5 to 4
func Signed[T constraints.Unsigned](x, max T) int {
	if x >= max/2 {
		return int(x - max)
	} else {
		return int(x)
	}
}

// Calculate Complement of x
// to calculate 10s Complement of x, you would call Complement(x, 10)
// this would return a value between 0 and 9
func Complement[T constraints.Integer](x, max T) uint {
	if x < 0 {
		return uint(max + x)
	} else if x == 0 {
		return 0
	} else {
		return uint(x % max)
	}
}

func abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
