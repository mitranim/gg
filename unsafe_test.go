package gg_test

import (
	"math"
	r "reflect"
	"testing"
	u "unsafe"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestAnyNoEscUnsafe(t *testing.T) {
	defer gtest.Catch(t)

	testAnyNoEscUnsafe(any(nil))
	testAnyNoEscUnsafe(``)
	testAnyNoEscUnsafe(`str`)
	testAnyNoEscUnsafe(0)
	testAnyNoEscUnsafe(10)
	testAnyNoEscUnsafe(SomeModel{})
	testAnyNoEscUnsafe((func())(nil))
}

/*
This doesn't verify that the value doesn't escape, because it's tricky to
implement for different types. Instead, various benchmarks serve as indirect
indicators.
*/
func testAnyNoEscUnsafe[A any](src A) {
	tar := gg.AnyNoEscUnsafe(src)
	gtest.Equal(r.TypeOf(tar), r.TypeOf(src))
	gtest.Equal(tar, any(src))
}

func BenchmarkAnyNoEscUnsafe(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		val := []int{ind}
		gg.Nop1(esc(gg.AnyNoEscUnsafe(val)))
	}
}

func BenchmarkSizeof(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Sizeof[string]())
	}
}

func TestAsBytes(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.AsBytes[struct{}](nil))
	gtest.Zero(gg.AsBytes[bool](nil))
	gtest.Zero(gg.AsBytes[uint64](nil))

	{
		var src struct{}
		tar := gg.AsBytes(&src)

		gtest.Equal(tar, []byte{})
		gtest.Eq(u.Pointer(u.SliceData(tar)), u.Pointer(&src))
		gtest.Eq(u.Pointer(u.SliceData(tar)), u.Pointer(u.SliceData([]struct{}{})), `zerobase`)
		gtest.Len(tar, 0)
		gtest.Cap(tar, 0)
	}

	{
		var src bool
		gtest.False(src)

		tar := gg.AsBytes(&src)

		gtest.Equal(tar, []byte{0})
		gtest.Eq(u.Pointer(u.SliceData(tar)), u.Pointer(&src))
		gtest.Len(tar, 1)
		gtest.Cap(tar, 1)

		tar[0] = 1
		gtest.True(src)
	}

	{
		var src uint64
		gtest.Eq(src, 0)

		tar := gg.AsBytes(&src)

		gtest.Equal(tar, make([]byte, 8))
		gtest.Eq(u.Pointer(u.SliceData(tar)), u.Pointer(&src))
		gtest.Len(tar, 8)
		gtest.Cap(tar, 8)

		for ind := range tar {
			tar[ind] = 255
		}

		gtest.Eq(src, math.MaxUint64)
	}
}
