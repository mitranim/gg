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

func BenchmarkSize(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Size[string]())
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

	{
		type Tar struct {
			One   uint64
			Two   uint64
			Three uint64
		}

		src := Tar{10, 20, 30}
		bytes := gg.AsBytes(&src)
		tar := *gg.CastUnsafe[*Tar](bytes)

		gtest.Eq(src, tar)
		gtest.Eq(tar, Tar{10, 20, 30})

		gg.CastUnsafe[*Tar](bytes).Two = 40
		tar = *gg.CastUnsafe[*Tar](bytes)

		gtest.Eq(src, tar)
		gtest.Eq(src, Tar{10, 40, 30})
	}
}

func TestCast(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(`size mismatch: uint8 (size 1) vs int64 (size 8)`, func() {
		gg.Cast[int64](byte(0))
	})

	gtest.PanicStr(`size mismatch: int64 (size 8) vs uint8 (size 1)`, func() {
		gg.Cast[byte](int64(0))
	})

	gtest.PanicStr(`size mismatch: string (size 16) vs []uint8 (size 24)`, func() {
		gg.Cast[[]byte](string(``))
	})

	gtest.PanicStr(`size mismatch: []uint8 (size 24) vs string (size 16)`, func() {
		gg.Cast[string]([]byte(nil))
	})

	gtest.Zero(gg.Cast[struct{}]([0]struct{}{}))
	gtest.Eq(gg.Cast[int8](uint8(math.MaxUint8)), -1)
	gtest.Eq(gg.Cast[uint8](int8(math.MaxInt8)), 127)
	gtest.Eq(gg.Cast[uint8](int8(math.MinInt8)), 128)

	{
		type Src [16]byte
		type Tar struct{ Src }

		src := Src([]byte(`ef1e7d2249dc45fc`))
		gtest.Eq(string(src[:]), `ef1e7d2249dc45fc`)

		tar := gg.Cast[Tar](src)
		gtest.Eq(tar.Src, src)
		gtest.Eq(gg.Cast[Src](tar), src)
	}
}

func TestCastSlice(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(`size mismatch: uint8 (size 1) vs int64 (size 8)`, func() {
		gg.CastSlice[int64, byte](nil)
	})

	gtest.PanicStr(`size mismatch: int64 (size 8) vs uint8 (size 1)`, func() {
		gg.CastSlice[byte, int64](nil)
	})

	gtest.PanicStr(`size mismatch: string (size 16) vs []uint8 (size 24)`, func() {
		gg.CastSlice[[]byte, string](nil)
	})

	gtest.PanicStr(`size mismatch: []uint8 (size 24) vs string (size 16)`, func() {
		gg.CastSlice[string, []byte](nil)
	})

	gtest.Zero(gg.CastSlice[uint8, int8](nil))
	gtest.Zero(gg.CastSlice[int8, uint8](nil))

	{
		src := []int8{-128, -127, -1, 0, 1, 127}
		tar := gg.CastSlice[uint8](src)

		gtest.Equal(tar, []uint8{128, 129, 255, 0, 1, 127})
		gtest.Eq(u.Pointer(u.SliceData(tar)), u.Pointer(u.SliceData(src)))
		gtest.Len(tar, len(src))
		gtest.Cap(tar, cap(src))
	}
}
