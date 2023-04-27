package gg

import (
	"math"
)

/*
Short for "is finite". Missing feature of the standard "math" package.
True if the input is neither NaN nor infinity.
*/
func IsFin[A Float](val A) bool {
	tar := float64(val)
	return !math.IsNaN(tar) && !math.IsInf(tar, 0)
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

// Same as `Add(val, 1)`. Panics on overflow.
func Inc[A Int](val A) A { return Add(val, 1) }

// Same as `Sub(val, 1)`. Panics on underflow.
func Dec[A Int](val A) A { return Sub(val, 1) }

/*
Raises a number to a power. Same as `math.Pow` and calls it under the hood, but
accepts arbitrary numeric types and performs checked conversions via `NumConv`.
Panics on overflow or precision loss. Has minor overhead over `math.Pow`.
Compare `PowUncheck` which runs faster but may overflow.
*/
func Pow[Tar, Pow Num](src Tar, pow Pow) Tar {
	return NumConv[Tar](math.Pow(NumConv[float64](src), NumConv[float64](pow)))
}

/*
Raises a number to a power. Same as `math.Pow` and calls it under the hood, but
accepts arbitrary numeric types. Does not check for overflow or precision loss.
Counterpart to `Pow` which panics on overflow.
*/
func PowUncheck[Tar, Pow Num](src Tar, pow Pow) Tar {
	return Tar(math.Pow(float64(src), float64(pow)))
}

/*
Checked factorial. Panics on overflow. Compare `FacUncheck` which runs faster,
but may overflow.
*/
func Fac[A Uint](src A) A {
	var tar float64 = 1
	mul := NumConv[float64](src)
	for mul > 0 {
		tar *= mul
		mul--
	}
	return NumConv[A](tar)
}

/*
Unchecked factorial. May overflow. Counterpart to `Fac` which panics on
overflow.
*/
func FacUncheck[A Uint](src A) A {
	var out A = 1
	for src > 0 {
		out *= src
		src -= 1
	}
	return out
}

// Checked addition. Panics on overflow/underflow. Has overhead.
func Add[A Int](one, two A) A {
	out := one + two
	if (out > one) == (two > 0) {
		return out
	}
	panic(errAdd(one, two, out))
}

func errAdd[A Int](one, two, out A) Err {
	return Errf(
		`addition overflow for %v: %v + %v = %v`,
		Type[A](), one, two, out,
	)
}

/*
Unchecked addition. Same as Go's `+` operator for numbers, expressed as a
generic function. Does not take strings. May overflow. For integers, prefer
`Add` whenever possible, which has overflow checks.
*/
func AddUncheck[A Num](one, two A) A { return one + two }

// Checked subtraction. Panics on overflow/underflow. Has overhead.
func Sub[A Int](one, two A) A {
	out := one - two
	if (out < one) == (two > 0) {
		return out
	}
	panic(errSub(one, two, out))
}

func errSub[A Int](one, two, out A) Err {
	return Errf(
		`subtraction overflow for %v: %v - %v = %v`,
		Type[A](), one, two, out,
	)
}

/*
Unchecked subtraction. Same as Go's `-` operator, expressed as a generic
function. May overflow. For integers, prefer `Sub` whenever possible, which has
overflow checks.
*/
func SubUncheck[A Num](one, two A) A { return one - two }

// Checked multiplication. Panics on overflow/underflow. Has overhead.
func Mul[A Int](one, two A) A {
	if one == 0 || two == 0 {
		return 0
	}
	out := one * two
	if ((one < 0) == (two < 0)) != (out < 0) && out/two == one {
		return out
	}
	panic(errMul(one, two, out))
}

func errMul[A Int](one, two, out A) Err {
	return Errf(
		`multiplication overflow for %v: %v * %v = %v`,
		Type[A](), one, two, out,
	)
}

/*
Checked numeric conversion. Same as `Out(src)` but with additional assertions.
Panics in case of overflow, underflow, imprecision such as converting large
integers to floats that don't support integers in that range, or anything
involving NaN. Performance overhead is minimal.
*/
func NumConv[Out, Src Num](src Src) Out {
	out := Out(src)
	if !(src == Src(out)) ||
		(src < 0 && out >= 0) ||
		(src >= 0 && out < 0) {
		panic(errNumConv(src, out))
	}
	return out
}

// Uses `String` to avoid the scientific notation for floats.
func errNumConv[Out, Src Num](src Src, out Out) Err {
	return Errf(
		`unable to safely convert %v %v to %v %v due to overflow, underflow, or imprecision`,
		Type[Src](), String(src), Type[Out](), String(out),
	)
}
