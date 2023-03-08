package gg

import "math"

/*
Short for "is finite". Missing feature of the standard "math" package.
True if the input is neither NaN nor infinity.
*/
func IsFin[A Float](val A) bool {
	flo := float64(val)
	return !math.IsNaN(flo) && !math.IsInf(flo, 0)
}

// Short for "is natural". True if >= 0. Also see `IsPos`.
func IsNat[A Num](val A) bool { return val >= 0 }

// Short for "is positive". True if > 0. Also see `IsNat`.
func IsPos[A Num](val A) bool { return val > 0 }

// Short for "is negative". True if < 0. Also see `IsNat`.
func IsNeg[A Num](val A) bool { return val < 0 }

/*
True if the remainder of dividing the first argument by the second argument is
zero. If the divisor is zero, does not attempt the division and returns false.
Note that the result is unaffected by the signs of either the dividend or the
divisor.
*/
func IsDivisibleBy[A Int](dividend, divisor A) bool {
	return divisor != 0 && dividend%divisor == 0
}

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

// Factorial without overflow checks. May overflow.
func Fac[A Uint](src A) A {
	var out A = 1
	for src > 0 {
		out *= src
		src -= 1
	}
	return out
}
