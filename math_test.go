package gg_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

type (
	Byte       uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Uint       uint
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Int        int
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128
)

func TestIsFin(t *testing.T) {
	defer gtest.Catch(t)

	gtest.False(gg.IsFin(math.NaN()))
	gtest.False(gg.IsFin(math.Inf(1)))
	gtest.False(gg.IsFin(math.Inf(-1)))
	gtest.True(gg.IsFin(0.0))
}

func TestIsDivisibleBy(t *testing.T) {
	defer gtest.Catch(t)

	gtest.False(gg.IsDivisibleBy(0, 0))
	gtest.False(gg.IsDivisibleBy(1, 0))
	gtest.False(gg.IsDivisibleBy(2, 0))
	gtest.False(gg.IsDivisibleBy(-1, 0))
	gtest.False(gg.IsDivisibleBy(-2, 0))

	gtest.True(gg.IsDivisibleBy(0, 1))
	gtest.True(gg.IsDivisibleBy(0, 2))
	gtest.True(gg.IsDivisibleBy(0, -1))
	gtest.True(gg.IsDivisibleBy(0, -2))

	gtest.True(gg.IsDivisibleBy(1, 1))
	gtest.True(gg.IsDivisibleBy(2, 1))
	gtest.True(gg.IsDivisibleBy(3, 1))
	gtest.True(gg.IsDivisibleBy(-1, 1))
	gtest.True(gg.IsDivisibleBy(-2, 1))
	gtest.True(gg.IsDivisibleBy(-3, 1))

	gtest.True(gg.IsDivisibleBy(1, -1))
	gtest.True(gg.IsDivisibleBy(2, -1))
	gtest.True(gg.IsDivisibleBy(3, -1))
	gtest.True(gg.IsDivisibleBy(-1, -1))
	gtest.True(gg.IsDivisibleBy(-2, -1))
	gtest.True(gg.IsDivisibleBy(-3, -1))

	gtest.False(gg.IsDivisibleBy(1, -2))
	gtest.False(gg.IsDivisibleBy(1, -3))
	gtest.False(gg.IsDivisibleBy(1, 2))
	gtest.False(gg.IsDivisibleBy(1, 3))

	gtest.False(gg.IsDivisibleBy(4, 0))
	gtest.True(gg.IsDivisibleBy(4, 1))
	gtest.True(gg.IsDivisibleBy(4, 2))
	gtest.False(gg.IsDivisibleBy(4, 3))
	gtest.True(gg.IsDivisibleBy(4, 4))
	gtest.False(gg.IsDivisibleBy(4, 5))
	gtest.False(gg.IsDivisibleBy(4, 6))
	gtest.False(gg.IsDivisibleBy(4, 7))
	gtest.False(gg.IsDivisibleBy(4, 8))
	gtest.False(gg.IsDivisibleBy(4, 9))
	gtest.False(gg.IsDivisibleBy(4, 10))
	gtest.False(gg.IsDivisibleBy(4, 11))
	gtest.False(gg.IsDivisibleBy(4, 12))
	gtest.False(gg.IsDivisibleBy(4, 13))
	gtest.False(gg.IsDivisibleBy(4, 14))
	gtest.False(gg.IsDivisibleBy(4, 15))
	gtest.False(gg.IsDivisibleBy(4, 16))

	gtest.False(gg.IsDivisibleBy(-4, 0))
	gtest.True(gg.IsDivisibleBy(-4, 1))
	gtest.True(gg.IsDivisibleBy(-4, 2))
	gtest.False(gg.IsDivisibleBy(-4, 3))
	gtest.True(gg.IsDivisibleBy(-4, 4))
	gtest.False(gg.IsDivisibleBy(-4, 5))
	gtest.False(gg.IsDivisibleBy(-4, 6))
	gtest.False(gg.IsDivisibleBy(-4, 7))
	gtest.False(gg.IsDivisibleBy(-4, 8))
	gtest.False(gg.IsDivisibleBy(-4, 9))
	gtest.False(gg.IsDivisibleBy(-4, 10))
	gtest.False(gg.IsDivisibleBy(-4, 11))
	gtest.False(gg.IsDivisibleBy(-4, 12))
	gtest.False(gg.IsDivisibleBy(-4, 13))
	gtest.False(gg.IsDivisibleBy(-4, 14))
	gtest.False(gg.IsDivisibleBy(-4, 15))
	gtest.False(gg.IsDivisibleBy(-4, 16))

	gtest.True(gg.IsDivisibleBy(4, -1))
	gtest.True(gg.IsDivisibleBy(4, -2))
	gtest.False(gg.IsDivisibleBy(4, -3))
	gtest.True(gg.IsDivisibleBy(4, -4))
	gtest.False(gg.IsDivisibleBy(4, -5))
	gtest.False(gg.IsDivisibleBy(4, -6))
	gtest.False(gg.IsDivisibleBy(4, -7))
	gtest.False(gg.IsDivisibleBy(4, -8))
	gtest.False(gg.IsDivisibleBy(4, -9))
	gtest.False(gg.IsDivisibleBy(4, -10))
	gtest.False(gg.IsDivisibleBy(4, -11))
	gtest.False(gg.IsDivisibleBy(4, -12))
	gtest.False(gg.IsDivisibleBy(4, -13))
	gtest.False(gg.IsDivisibleBy(4, -14))
	gtest.False(gg.IsDivisibleBy(4, -15))
	gtest.False(gg.IsDivisibleBy(4, -16))

	gtest.True(gg.IsDivisibleBy(-4, -1))
	gtest.True(gg.IsDivisibleBy(-4, -2))
	gtest.False(gg.IsDivisibleBy(-4, -3))
	gtest.True(gg.IsDivisibleBy(-4, -4))
	gtest.False(gg.IsDivisibleBy(-4, -5))
	gtest.False(gg.IsDivisibleBy(-4, -6))
	gtest.False(gg.IsDivisibleBy(-4, -7))
	gtest.False(gg.IsDivisibleBy(-4, -8))
	gtest.False(gg.IsDivisibleBy(-4, -9))
	gtest.False(gg.IsDivisibleBy(-4, -10))
	gtest.False(gg.IsDivisibleBy(-4, -11))
	gtest.False(gg.IsDivisibleBy(-4, -12))
	gtest.False(gg.IsDivisibleBy(-4, -13))
	gtest.False(gg.IsDivisibleBy(-4, -14))
	gtest.False(gg.IsDivisibleBy(-4, -15))
	gtest.False(gg.IsDivisibleBy(-4, -16))
}

func TestIsFrac(t *testing.T) {
	defer gtest.Catch(t)
	testIsFrac[float32]()
	testIsFrac[float64]()
}

func testIsFrac[A gg.Float]() {
	gtest.False(gg.IsFrac(A(math.NaN())))
	gtest.False(gg.IsFrac(A(math.Inf(-1))))
	gtest.False(gg.IsFrac(A(math.Inf(+1))))
	gtest.False(gg.IsFrac(A(-0)))
	gtest.False(gg.IsFrac(A(+0)))
	gtest.False(gg.IsFrac(A(-1)))
	gtest.False(gg.IsFrac(A(+1)))
	gtest.False(gg.IsFrac(A(-2)))
	gtest.False(gg.IsFrac(A(+2)))
	gtest.False(gg.IsFrac(A(-12)))
	gtest.False(gg.IsFrac(A(+12)))
	gtest.False(gg.IsFrac(A(gg.MinSafeIntFloat32)))
	gtest.False(gg.IsFrac(A(gg.MaxSafeIntFloat32)))

	gtest.True(gg.IsFrac(A(-0.000001)))
	gtest.True(gg.IsFrac(A(+0.000001)))
	gtest.True(gg.IsFrac(A(-1.000001)))
	gtest.True(gg.IsFrac(A(+1.000001)))
	gtest.True(gg.IsFrac(A(-2.000001)))
	gtest.True(gg.IsFrac(A(+2.000001)))
	gtest.True(gg.IsFrac(A(-12.000001)))
	gtest.True(gg.IsFrac(A(+12.000001)))

	gtest.True(gg.IsFrac(A(-0.111111)))
	gtest.True(gg.IsFrac(A(+0.111111)))
	gtest.True(gg.IsFrac(A(-1.111111)))
	gtest.True(gg.IsFrac(A(+1.111111)))
	gtest.True(gg.IsFrac(A(-2.111111)))
	gtest.True(gg.IsFrac(A(+2.111111)))
	gtest.True(gg.IsFrac(A(-12.111111)))
	gtest.True(gg.IsFrac(A(+12.111111)))

	gtest.True(gg.IsFrac(A(-0.5)))
	gtest.True(gg.IsFrac(A(+0.5)))
	gtest.True(gg.IsFrac(A(-1.5)))
	gtest.True(gg.IsFrac(A(+1.5)))
	gtest.True(gg.IsFrac(A(-2.5)))
	gtest.True(gg.IsFrac(A(+2.5)))
	gtest.True(gg.IsFrac(A(-12.5)))
	gtest.True(gg.IsFrac(A(+12.5)))

	gtest.True(gg.IsFrac(A(-0.999999)))
	gtest.True(gg.IsFrac(A(+0.999999)))
	gtest.True(gg.IsFrac(A(-1.999999)))
	gtest.True(gg.IsFrac(A(+1.999999)))
	gtest.True(gg.IsFrac(A(-2.999999)))
	gtest.True(gg.IsFrac(A(+2.999999)))
	gtest.True(gg.IsFrac(A(-12.999999)))
	gtest.True(gg.IsFrac(A(+12.999999)))
}

func TestInc(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Inc(-3), -2)
	gtest.Eq(gg.Inc(-2), -1)
	gtest.Eq(gg.Inc(-1), 0)
	gtest.Eq(gg.Inc(0), 1)
	gtest.Eq(gg.Inc(1), 2)
	gtest.Eq(gg.Inc(2), 3)
	gtest.Eq(gg.Inc(3), 4)

	gtest.PanicStr(
		`addition overflow for uint8: 255 + 1 = 0`,
		func() { gg.Inc[uint8](math.MaxUint8) },
	)

	gtest.PanicStr(
		`addition overflow for int8: 127 + 1 = -128`,
		func() { gg.Inc[int8](math.MaxInt8) },
	)

	/**
	TODO restore support for floats.

	gtest.Eq(gg.Inc(-3.5), -2.5)
	gtest.Eq(gg.Inc(-2.5), -1.5)
	gtest.Eq(gg.Inc(-1.5), -0.5)
	gtest.Eq(gg.Inc(-0.5), 0.5)
	gtest.Eq(gg.Inc(0.0), 1)
	gtest.Eq(gg.Inc(0.5), 1.5)
	gtest.Eq(gg.Inc(1.5), 2.5)
	gtest.Eq(gg.Inc(2.5), 3.5)
	gtest.Eq(gg.Inc(3.5), 4.5)

	gtest.PanicStr(...)
	*/
}

func TestDec(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Dec(-3), -4)
	gtest.Eq(gg.Dec(-2), -3)
	gtest.Eq(gg.Dec(-1), -2)
	gtest.Eq(gg.Dec(0), -1)
	gtest.Eq(gg.Dec(1), 0)
	gtest.Eq(gg.Dec(2), 1)
	gtest.Eq(gg.Dec(3), 2)

	gtest.PanicStr(
		`subtraction overflow for uint8: 0 - 1 = 255`,
		func() { gg.Dec[uint8](0) },
	)

	gtest.PanicStr(
		`subtraction overflow for int8: -128 - 1 = 127`,
		func() { gg.Dec[int8](math.MinInt8) },
	)

	/**
	TODO restore support for floats.

	gtest.Eq(gg.Dec(-3.5), -4.5)
	gtest.Eq(gg.Dec(-2.5), -3.5)
	gtest.Eq(gg.Dec(-1.5), -2.5)
	gtest.Eq(gg.Dec(-0.5), -1.5)
	gtest.Eq(gg.Dec(0.0), -1.0)
	gtest.Eq(gg.Dec(0.5), -0.5)
	gtest.Eq(gg.Dec(1.5), 0.5)
	gtest.Eq(gg.Dec(2.5), 1.5)
	gtest.Eq(gg.Dec(3.5), 2.5)

	gtest.PanicStr(...)
	*/
}

func TestPow(t *testing.T) {
	defer gtest.Catch(t)

	testPowInt(gg.Pow[int, int])
	testPowFloat(gg.Pow[float64, float64])

	gtest.PanicStr(
		`unable to safely convert float64 1162261467 to uint8`,
		func() { gg.Pow[uint8](3, 19) },
	)

	gtest.PanicStr(
		`unable to safely convert float64 1162261467 to int8`,
		func() { gg.Pow[int8](3, 19) },
	)
}

/*
TODO test fractional and negative powers.
We're simply calling `math.Pow`, but we do need a sanity check.
*/
func testPowInt(fun func(int, int) int) {
	src := []int{-12, -1, 0, 1, 12}
	for _, val := range src {
		testPow0(val, fun)
		testPow1(val, fun)
		testPow2(val, fun)
		testPow3(val, fun)
	}
}

func testPowFloat(fun func(float64, float64) float64) {
	src := []float64{-12.23, -1.2, -1.0, 0.0, 1.0, 1.2, 12.23}
	for _, val := range src {
		testPow0(val, fun)
		testPow1(val, fun)
		testPow2(val, fun)
		testPow3(val, fun)
	}
}

func testPow0[A gg.Num](src A, fun func(A, A) A) { gtest.Eq(fun(src, 0), 1) }
func testPow1[A gg.Num](src A, fun func(A, A) A) { gtest.Eq(fun(src, 1), src) }
func testPow2[A gg.Num](src A, fun func(A, A) A) { gtest.Eq(fun(src, 2), src*src) }
func testPow3[A gg.Num](src A, fun func(A, A) A) { gtest.Eq(fun(src, 3), src*src*src) }

func BenchmarkPow(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.Pow(79, 7)
	}
}

func TestPowUncheck(t *testing.T) {
	defer gtest.Catch(t)

	testPowInt(gg.PowUncheck[int, int])
	testPowFloat(gg.PowUncheck[float64, float64])
}

func BenchmarkPowUncheck(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.PowUncheck(79, 7)
	}
}

func TestFac(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src, exp uint64) { gtest.Eq(gg.Fac(src), exp) }

	test(0, 1)
	test(1, 1)
	test(2, 2)
	test(3, 6)
	test(4, 24)
	test(5, 120)
	test(6, 720)
	test(7, 5040)
	test(8, 40320)
	test(9, 362880)
	test(10, 3628800)
	test(11, 39916800)
	test(12, 479001600)
	test(13, 6227020800)
	test(14, 87178291200)
	test(15, 1307674368000)
	test(16, 20922789888000)
	test(17, 355687428096000)
	test(18, 6402373705728000)

	gtest.PanicStr(
		`unable to safely convert float64 720 to uint8`,
		func() { gg.Fac[uint8](6) },
	)
}

func BenchmarkFac(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.Fac[uint64](19)
	}
}

func TestFacUncheck(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src, exp uint64) { gtest.Eq(gg.FacUncheck(src), exp) }

	test(0, 1)
	test(1, 1)
	test(2, 2)
	test(3, 6)
	test(4, 24)
	test(5, 120)
	test(6, 720)
	test(7, 5040)
	test(8, 40320)
	test(9, 362880)
	test(10, 3628800)
	test(11, 39916800)
	test(12, 479001600)
	test(13, 6227020800)
	test(14, 87178291200)
	test(15, 1307674368000)
	test(16, 20922789888000)
	test(17, 355687428096000)
	test(18, 6402373705728000)
	test(19, 121645100408832000)
	test(20, 2432902008176640000)

	gtest.Eq(gg.FacUncheck[uint8](6), 208, `expecting overflow`)
}

func BenchmarkFacUncheck(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.FacUncheck[uint64](19)
	}
}

/*
Supplementary for `gg.Add`.

Definitions:

	A = addend
	B = addend
	V = valid output
	O = overflow output

States for unsigned integers:

			        -->
	0+++++++++++ | 0+++++
	AB           | V
	A     B      |    V
	B     A      |    V
	      AB     | O  OV
	   A     B   |    OV
	   B     A   |    OV

States for signed integers:

			               -->
	---------0+++++++++ | ---------0+++++++++
	         AB         |          V
	    A    B          |     V
	    B    A          |     V
	         A    B     |               V
	         B    A     |               V
	    A         B     |     V    V    V
	    B         A     |     V    V    V
	  AB                |     V    O    O
	  A   B             |     V         O
	  B   A             |     V         O
	            AB      |     O         V
	            A   B   |     O         V
	            B   A   |     O         V
*/

func TestAdd_uint8(t *testing.T) {
	defer gtest.Catch(t)

	type Type = uint8
	fun := gg.Add[Type]
	enum := gg.RangeIncl[Type](0, math.MaxUint8)

	/**
	This should cover all possible cases. The hardcoded assertions below serve as
	a sanity check and documentation.
	*/
	for _, one := range enum {
		for _, two := range enum {
			if int(one+two) == int(one)+int(two) {
				gtest.Eq(fun(one, two), one+two, one, two)
				continue
			}

			gtest.PanicStr(
				fmt.Sprintf(`addition overflow for uint8: %v + %v`, one, two),
				func() { fun(one, two) },
				one, two,
			)
		}
	}

	gtest.Eq(fun(0, 0), 0)
	gtest.Eq(fun(3, 5), 8)
	gtest.Eq(fun(13, 7), 20)
	gtest.Eq(fun(103, 152), 255)

	gtest.PanicStr(
		`addition overflow for uint8: 255 + 255 = 254`,
		func() { fun(255, 255) },
	)

	gtest.PanicStr(
		`addition overflow for uint8: 103 + 153 = 0`,
		func() { fun(103, 153) },
	)

	gtest.PanicStr(
		`addition overflow for uint8: 199 + 239 = 182`,
		func() { fun(199, 239) },
	)
}

func TestAdd_int8(t *testing.T) {
	defer gtest.Catch(t)

	type Type = int8
	fun := gg.Add[Type]
	enum := gg.RangeIncl[Type](math.MinInt8, math.MaxInt8)

	/**
	This should cover all possible cases. The hardcoded assertions below serve as
	a sanity check and documentation.
	*/
	for _, one := range enum {
		for _, two := range enum {
			if int(one+two) == int(one)+int(two) {
				gtest.Eq(fun(one, two), one+two, one, two)
				continue
			}

			gtest.PanicStr(
				fmt.Sprintf(`addition overflow for int8: %v + %v`, one, two),
				func() { fun(one, two) },
				one, two,
			)
		}
	}

	gtest.Eq(fun(0, 0), 0)

	gtest.Eq(fun(3, 5), 8)
	gtest.Eq(fun(13, 7), 20)
	gtest.Eq(fun(79, 48), 127)

	gtest.Eq(fun(-3, 5), 2)
	gtest.Eq(fun(-13, 7), -6)
	gtest.Eq(fun(-79, 48), -31)

	gtest.Eq(fun(3, -5), -2)
	gtest.Eq(fun(13, -7), 6)
	gtest.Eq(fun(79, -48), 31)

	gtest.Eq(fun(-3, -5), -8)
	gtest.Eq(fun(-13, -7), -20)
	gtest.Eq(fun(-79, -49), -128)

	gtest.Eq(fun(127, -128), -1)
	gtest.Eq(fun(-128, 127), -1)

	gtest.PanicStr(
		`addition overflow for int8: 127 + 127 = -2`,
		func() { fun(127, 127) },
	)

	gtest.PanicStr(
		`addition overflow for int8: -128 + -128 = 0`,
		func() { fun(-128, -128) },
	)

	gtest.PanicStr(
		`addition overflow for int8: 79 + 97 = -80`,
		func() { fun(79, 97) },
	)

	gtest.PanicStr(
		`addition overflow for int8: -79 + -97 = 80`,
		func() { fun(-79, -97) },
	)
}

func TestAdd_uint16(t *testing.T) {
	defer gtest.Catch(t)

	type Type = uint16
	fun := gg.Add[Type]

	gtest.Eq(fun(0, 0), 0)
	gtest.Eq(fun(3, 5), 8)
	gtest.Eq(fun(13, 7), 20)
	gtest.Eq(fun(21963, 43572), 65535)

	gtest.PanicStr(
		`addition overflow for uint16: 65535 + 65535 = 65534`,
		func() { fun(65535, 65535) },
	)

	gtest.PanicStr(
		`addition overflow for uint16: 21963 + 43573 = 0`,
		func() { fun(21963, 43573) },
	)

	gtest.PanicStr(
		`addition overflow for uint16: 43573 + 39571 = 17608`,
		func() { fun(43573, 39571) },
	)
}

func TestAdd_int16(t *testing.T) {
	defer gtest.Catch(t)

	type Type = int16
	fun := gg.Add[Type]

	gtest.Eq(fun(0, 0), 0)

	gtest.Eq(fun(3, 5), 8)
	gtest.Eq(fun(13, 7), 20)
	gtest.Eq(fun(79, 48), 127)
	gtest.Eq(fun(21963, 10804), 32767)

	gtest.Eq(fun(-3, 5), 2)
	gtest.Eq(fun(-13, 7), -6)
	gtest.Eq(fun(-79, 48), -31)
	gtest.Eq(fun(-21963, 10804), -11159)

	gtest.Eq(fun(3, -5), -2)
	gtest.Eq(fun(13, -7), 6)
	gtest.Eq(fun(79, -48), 31)
	gtest.Eq(fun(21963, -10804), 11159)

	gtest.Eq(fun(-3, -5), -8)
	gtest.Eq(fun(-13, -7), -20)
	gtest.Eq(fun(-79, -49), -128)
	gtest.Eq(fun(-21963, -10804), -32767)

	gtest.Eq(fun(32767, -32768), -1)
	gtest.Eq(fun(-32768, 32767), -1)

	gtest.PanicStr(
		`addition overflow for int16: 32767 + 32767 = -2`,
		func() { fun(32767, 32767) },
	)

	gtest.PanicStr(
		`addition overflow for int16: -32768 + -32768 = 0`,
		func() { fun(-32768, -32768) },
	)

	gtest.PanicStr(
		`addition overflow for int16: 21963 + 28436 = -15137`,
		func() { fun(21963, 28436) },
	)

	gtest.PanicStr(
		`addition overflow for int16: -21963 + -28436 = 15137`,
		func() { fun(-21963, -28436) },
	)
}

/*
Supplementary for `gg.Sub`.

Definitions:

	A = minuend
	B = subtrahend
	V = valid output
	O = overflow output

States for unsigned integers:

			        -->
	0+++++++++++ | 0+++++
	AB           | V
	A     B      |    O
	B     A      |    V
	      AB     | V
	   A     B   |    O
	   B     A   |    V

States for signed integers:

			               -->
	---------0+++++++++ | ---------0+++++++++
	         AB         |          V
	    A    B          |     V
	    B    A          |     O         V
	         A    B     |     V
	         B    A     |               V
	    A         B     |     V         O
	    B         A     |     O         V
	  AB                |          V
	  A   B             |     V
	  B   A             |               V
	            AB      |          V
	            A   B   |     V
	            B   A   |               V
*/

// TODO tests for wider types.
func TestSub_uint8(t *testing.T) {
	defer gtest.Catch(t)

	type Type = uint8
	fun := gg.Sub[Type]
	enum := gg.RangeIncl[Type](0, math.MaxUint8)

	/**
	This should cover all possible cases. The hardcoded assertions below serve as
	a sanity check and documentation.
	*/
	for _, one := range enum {
		for _, two := range enum {
			if int(one-two) == int(one)-int(two) {
				gtest.Eq(fun(one, two), one-two, one, two)
				continue
			}

			gtest.PanicStr(
				fmt.Sprintf(`subtraction overflow for uint8: %v - %v`, one, two),
				func() { fun(one, two) },
				one, two,
			)
		}
	}

	gtest.Eq(fun(0, 0), 0)
	gtest.Eq(fun(1, 1), 0)
	gtest.Eq(fun(1, 0), 1)
	gtest.Eq(fun(5, 3), 2)
	gtest.Eq(fun(13, 7), 6)
	gtest.Eq(fun(152, 103), 49)
	gtest.Eq(fun(255, 0), 255)
	gtest.Eq(fun(255, 1), 254)
	gtest.Eq(fun(255, 254), 1)
	gtest.Eq(fun(255, 255), 0)

	gtest.PanicStr(
		`subtraction overflow for uint8: 0 - 1 = 255`,
		func() { fun(0, 1) },
	)

	gtest.PanicStr(
		`subtraction overflow for uint8: 0 - 255 = 1`,
		func() { fun(0, 255) },
	)

	gtest.PanicStr(
		`subtraction overflow for uint8: 103 - 153 = 206`,
		func() { fun(103, 153) },
	)

	gtest.PanicStr(
		`subtraction overflow for uint8: 79 - 255 = 80`,
		func() { fun(79, 255) },
	)
}

// TODO tests for wider types.
func TestSub_int8(t *testing.T) {
	defer gtest.Catch(t)

	type Type = int8
	fun := gg.Sub[Type]
	enum := gg.RangeIncl[Type](math.MinInt8, math.MaxInt8)

	/**
	This should cover all possible cases. The hardcoded assertions below serve as
	a sanity check and documentation.
	*/
	for _, one := range enum {
		for _, two := range enum {
			if int(one-two) == int(one)-int(two) {
				gtest.Eq(fun(one, two), one-two, one, two)
				continue
			}

			gtest.PanicStr(
				fmt.Sprintf(`subtraction overflow for int8: %v - %v`, one, two),
				func() { fun(one, two) },
				one, two,
			)
		}
	}

	gtest.Eq(fun(0, 0), 0)

	gtest.Eq(fun(3, 5), -2)
	gtest.Eq(fun(13, 7), 6)
	gtest.Eq(fun(79, 48), 31)

	gtest.Eq(fun(-3, 5), -8)
	gtest.Eq(fun(-13, 7), -20)
	gtest.Eq(fun(-79, 48), -127)
	gtest.Eq(fun(-79, 49), -128)

	gtest.Eq(fun(3, -5), 8)
	gtest.Eq(fun(13, -7), 20)
	gtest.Eq(fun(79, -48), 127)

	gtest.Eq(fun(-3, -5), 2)
	gtest.Eq(fun(-13, -7), -6)
	gtest.Eq(fun(-79, -49), -30)

	gtest.Eq(fun(127, 0), 127)
	gtest.Eq(fun(127, 1), 126)
	gtest.Eq(fun(127, 126), 1)
	gtest.Eq(fun(127, 127), 0)

	gtest.Eq(fun(-128, 0), -128)
	gtest.Eq(fun(-128, -1), -127)
	gtest.Eq(fun(-128, -127), -1)
	gtest.Eq(fun(-128, -128), 0)

	gtest.PanicStr(
		`subtraction overflow for int8: -128 - 1 = 127`,
		func() { fun(-128, 1) },
	)

	gtest.PanicStr(
		`subtraction overflow for int8: -128 - 2 = 126`,
		func() { fun(-128, 2) },
	)

	gtest.PanicStr(
		`subtraction overflow for int8: -128 - 127 = 1`,
		func() { fun(-128, 127) },
	)

	gtest.PanicStr(
		`subtraction overflow for int8: 127 - -1 = -128`,
		func() { fun(127, -1) },
	)

	gtest.PanicStr(
		`subtraction overflow for int8: 127 - -2 = -127`,
		func() { fun(127, -2) },
	)

	gtest.PanicStr(
		`subtraction overflow for int8: 127 - -128 = -1`,
		func() { fun(127, -128) },
	)

	gtest.PanicStr(
		`subtraction overflow for int8: -79 - 50 = 127`,
		func() { fun(-79, 50) },
	)

	gtest.PanicStr(
		`subtraction overflow for int8: 79 - -49 = -128`,
		func() { fun(79, -49) },
	)
}

/*
Supplementary for `gg.Mul`.

Definitions:

	A = multiplicand
	B = multiplicand
	V = valid output
	O = overflow output

States for unsigned integers:

			        -->
	0+++++++++++ | 0++++++
	AB           | V
	A     B      | V
	B     A      | V
	      AB     | O  OV
	   A     B   | O  OV
	   B     A   | O  OV

States for signed integers:

			               -->
	---------0+++++++++ | ---------0+++++++++
	         AB         |          V
	    A    B          |          V
	    B    A          |          V
	         A    B     |          V
	         B    A     |          V
	    A         B     |     OV   O    O
	    B         A     |     OV   O    O
	  AB                |     OV   O    O
	  A   B             |     O    O    OV
	  B   A             |     O    O    OV
	            AB      |               OV
	            A   B   |     O    O    OV
	            B   A   |     O    O    OV
*/

// TODO tests for wider types.
func TestMul_uint8(t *testing.T) {
	defer gtest.Catch(t)

	type Type = uint8
	fun := gg.Mul[Type]
	enum := gg.RangeIncl[Type](0, math.MaxUint8)

	/**
	This should cover all possible cases. The hardcoded assertions below serve as
	a sanity check and documentation.
	*/
	for _, one := range enum {
		for _, two := range enum {
			if int(one*two) == int(one)*int(two) {
				gtest.Eq(fun(one, two), one*two, one, two)
				continue
			}

			gtest.PanicStr(
				fmt.Sprintf(`multiplication overflow for uint8: %v * %v`, one, two),
				func() { fun(one, two) },
				one, two,
			)
		}
	}

	gtest.Eq(fun(3, 5), 15)
	gtest.Eq(fun(5, 3), 15)

	gtest.Eq(fun(5, 7), 35)
	gtest.Eq(fun(7, 5), 35)

	gtest.Eq(fun(7, 11), 77)
	gtest.Eq(fun(11, 7), 77)

	gtest.Eq(fun(11, 13), 143)
	gtest.Eq(fun(13, 11), 143)

	gtest.Eq(fun(13, 17), 221)
	gtest.Eq(fun(17, 13), 221)

	gtest.Eq(fun(17, 15), 255)
	gtest.Eq(fun(15, 17), 255)

	gtest.PanicStr(
		`multiplication overflow for uint8: 255 * 255 = 1`,
		func() { fun(255, 255) },
	)

	gtest.PanicStr(
		`multiplication overflow for uint8: 17 * 19 = 67`,
		func() { fun(17, 19) },
	)

	gtest.PanicStr(
		`multiplication overflow for uint8: 19 * 17 = 67`,
		func() { fun(19, 17) },
	)

	gtest.PanicStr(
		`multiplication overflow for uint8: 2 * 128 = 0`,
		func() { fun(2, 128) },
	)

	gtest.PanicStr(
		`multiplication overflow for uint8: 128 * 2 = 0`,
		func() { fun(128, 2) },
	)
}

// TODO tests for wider types.
func TestMul_int8(t *testing.T) {
	defer gtest.Catch(t)

	type Type = int8
	fun := gg.Mul[Type]
	enum := gg.RangeIncl[Type](math.MinInt8, math.MaxInt8)

	/**
	This should cover all possible cases. The hardcoded assertions below serve as
	a sanity check and documentation.
	*/
	for _, one := range enum {
		for _, two := range enum {
			if int(one*two) == int(one)*int(two) {
				gtest.Eq(fun(one, two), one*two, one, two)
				continue
			}

			gtest.PanicStr(
				fmt.Sprintf(`multiplication overflow for int8: %v * %v`, one, two),
				func() { fun(one, two) },
				one, two,
			)
		}
	}

	gtest.Eq(fun(0, 0), 0)

	gtest.Eq(fun(0, 3), 0)
	gtest.Eq(fun(3, 0), 0)

	gtest.Eq(fun(0, -3), 0)
	gtest.Eq(fun(-3, 0), 0)

	gtest.Eq(fun(1, 3), 3)
	gtest.Eq(fun(3, 1), 3)

	gtest.Eq(fun(1, 127), 127)
	gtest.Eq(fun(127, 1), 127)

	gtest.Eq(fun(1, -128), -128)
	gtest.Eq(fun(-128, 1), -128)

	gtest.Eq(fun(1, -3), -3)
	gtest.Eq(fun(-3, 1), -3)

	gtest.Eq(fun(-1, 3), -3)
	gtest.Eq(fun(3, -1), -3)

	gtest.Eq(fun(-1, 127), -127)
	gtest.Eq(fun(127, -1), -127)

	gtest.Eq(fun(-1, -3), 3)
	gtest.Eq(fun(-3, -1), 3)

	gtest.Eq(fun(3, 5), 15)
	gtest.Eq(fun(5, 3), 15)

	gtest.Eq(fun(3, -5), -15)
	gtest.Eq(fun(-5, 3), -15)

	gtest.Eq(fun(-3, 5), -15)
	gtest.Eq(fun(5, -3), -15)

	gtest.Eq(fun(-3, -5), 15)
	gtest.Eq(fun(-5, -3), 15)

	gtest.Eq(fun(9, 14), 126)
	gtest.Eq(fun(14, 9), 126)

	gtest.Eq(fun(9, -14), -126)
	gtest.Eq(fun(-14, 9), -126)

	gtest.Eq(fun(-9, 14), -126)
	gtest.Eq(fun(14, -9), -126)

	gtest.Eq(fun(-9, -14), 126)
	gtest.Eq(fun(-14, -9), 126)

	gtest.PanicStr(
		`multiplication overflow for int8: 127 * 127 = 1`,
		func() { fun(127, 127) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: 127 * 126 = -126`,
		func() { fun(127, 126) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: 126 * 127 = -126`,
		func() { fun(126, 127) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: 126 * 126 = 4`,
		func() { fun(126, 126) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: 126 * 126 = 4`,
		func() { fun(126, 126) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -128 * -128 = 0`,
		func() { fun(-128, -128) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -128 * -127 = -128`,
		func() { fun(-128, -127) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -127 * -128 = -128`,
		func() { fun(-127, -128) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -127 * -127 = 1`,
		func() { fun(-127, -127) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: 127 * -128 = -128`,
		func() { fun(127, -128) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -128 * 127 = -128`,
		func() { fun(-128, 127) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -127 * 127 = -1`,
		func() { fun(-127, 127) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -126 * 127 = 126`,
		func() { fun(-126, 127) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -1 * -128 = -128`,
		func() { fun(-1, -128) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -128 * -1 = -128`,
		func() { fun(-128, -1) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: 11 * 13 = -113`,
		func() { fun(11, 13) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: 13 * 11 = -113`,
		func() { fun(13, 11) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -11 * -13 = -113`,
		func() { fun(-11, -13) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: -13 * -11 = -113`,
		func() { fun(-13, -11) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: 2 * 127 = -2`,
		func() { fun(2, 127) },
	)

	gtest.PanicStr(
		`multiplication overflow for int8: 127 * 2 = -2`,
		func() { fun(127, 2) },
	)
}

//go:noinline
func safePairInt8() (int8, int8) { return 5, 7 }

func Benchmark_mul_int8_native(b *testing.B) {
	defer gtest.Catch(b)
	one, two := safePairInt8()
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(one * two)
	}
}

func Benchmark_mul_int8_ours(b *testing.B) {
	defer gtest.Catch(b)
	one, two := safePairInt8()
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Mul[int8](one, two)
	}
}
