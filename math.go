package gg

import "math"

/*
Short for "power". Missing feature of the standard "math" package. Raises the
input, which may be an arbitrary number, to the given power. Current
limitations: power must be a natural number; no overflow check.
*/
func Pow[A Num](base A, pow int) A {
	if pow == 0 {
		return 1
	}

	out := base
	for range Iter(pow) {
		out *= base
	}
	return out
}

/*
Short for "is finite". Missing feature of the standard "math" package.
True if the input is neither NaN nor infinity.
*/
func IsFin[A Float](val A) bool {
	flo := float64(val)
	return !math.IsNaN(flo) && !math.IsInf(flo, 0)
}
