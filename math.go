package gg

import "math"

// Same as input + 1.
func Inc[A Num](val A) A { return val + 1 }

// Same as input - 1.
func Dec[A Num](val A) A { return val - 1 }

/*
Short for "power". Missing feature of the standard "math" package. Raises the
input, which may be an arbitrary number, to the given power. Current
limitations: power must be a natural number; no overflow check.
*/
func Pow[A Num](base A, pow int) A {
	var out A = 1
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

// Factorial without overflow checks. May overflow.
func Fac[A Uint](src A) A {
	var out A = 1
	for src > 0 {
		out *= src
		src -= 1
	}
	return out
}
