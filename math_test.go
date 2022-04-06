package gg_test

import (
	"math"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestInc(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Inc(-3), -2)
	gtest.Eq(gg.Inc(-2), -1)
	gtest.Eq(gg.Inc(-1), 0)
	gtest.Eq(gg.Inc(0), 1)
	gtest.Eq(gg.Inc(1), 2)
	gtest.Eq(gg.Inc(2), 3)
	gtest.Eq(gg.Inc(3), 4)

	gtest.Eq(gg.Inc(-3.5), -2.5)
	gtest.Eq(gg.Inc(-2.5), -1.5)
	gtest.Eq(gg.Inc(-1.5), -0.5)
	gtest.Eq(gg.Inc(-0.5), 0.5)
	gtest.Eq(gg.Inc(0.0), 1)
	gtest.Eq(gg.Inc(0.5), 1.5)
	gtest.Eq(gg.Inc(1.5), 2.5)
	gtest.Eq(gg.Inc(2.5), 3.5)
	gtest.Eq(gg.Inc(3.5), 4.5)
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

	gtest.Eq(gg.Dec(-3.5), -4.5)
	gtest.Eq(gg.Dec(-2.5), -3.5)
	gtest.Eq(gg.Dec(-1.5), -2.5)
	gtest.Eq(gg.Dec(-0.5), -1.5)
	gtest.Eq(gg.Dec(0.0), -1.0)
	gtest.Eq(gg.Dec(0.5), -0.5)
	gtest.Eq(gg.Dec(1.5), 0.5)
	gtest.Eq(gg.Dec(2.5), 1.5)
	gtest.Eq(gg.Dec(3.5), 2.5)
}

func TestPow(t *testing.T) {
	defer gtest.Catch(t)

	testInts := []int{0, 1, 12, -1, -12}
	testFloats := []float64{0.0, 1.0, 1.2, 12.23, -1.0, -1.2, -12.23}

	gg.Each(testInts, testPow0[int])
	gg.Each(testFloats, testPow0[float64])

	gg.Each(testInts, testPow1[int])
	gg.Each(testFloats, testPow1[float64])

	gg.Each(testInts, testPow2[int])
	gg.Each(testFloats, testPow2[float64])

	gg.Each(testInts, testPow3[int])
	gg.Each(testFloats, testPow3[float64])
}

func testPow0[A gg.Num](src A) { gtest.Eq(gg.Pow(src, 0), 1) }

func testPow1[A gg.Num](src A) { gtest.Eq(gg.Pow(src, 1), src) }

func testPow2[A gg.Num](src A) {
	gtest.Eq(gg.Pow(src, 2), src*src)
}

func testPow3[A gg.Num](src A) {
	gtest.Eq(gg.Pow(src, 3), src*src*src)
}

func TestIsFin(t *testing.T) {
	defer gtest.Catch(t)

	gtest.False(gg.IsFin(math.NaN()))
	gtest.False(gg.IsFin(math.Inf(1)))
	gtest.False(gg.IsFin(math.Inf(-1)))
	gtest.True(gg.IsFin(0.0))
	gtest.True(gg.IsFin(-0.0))
}
