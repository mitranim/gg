package gg_test

import (
	r "reflect"
	"testing"

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
