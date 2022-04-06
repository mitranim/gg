package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

// Placeholder, needs a concurrency test.
func TestAtom(t *testing.T) {
	defer gtest.Catch(t)

	var ref gg.Atom[string]
	gtest.Zero(ref.Load())
	gtest.Eq(
		gg.Tuple2(ref.Loaded()),
		gg.Tuple2(``, false),
	)

	ref.Store(``)
	gtest.Zero(ref.Load())
	gtest.Eq(
		gg.Tuple2(ref.Loaded()),
		gg.Tuple2(``, true),
	)

	ref.Store(`one`)
	gtest.Eq(ref.Load(), `one`)
	gtest.Eq(
		gg.Tuple2(ref.Loaded()),
		gg.Tuple2(`one`, true),
	)

	gtest.False(ref.CompareAndSwap(`three`, `two`))
	gtest.Eq(ref.Load(), `one`)

	gtest.True(ref.CompareAndSwap(`one`, `two`))
	gtest.Eq(ref.Load(), `two`)
}

func BenchmarkAtom_Store(b *testing.B) {
	var ref gg.Atom[string]

	for i := 0; i < b.N; i++ {
		ref.Store(`str`)
	}
}

func BenchmarkAtom_Load(b *testing.B) {
	var ref gg.Atom[string]
	ref.Store(`str`)

	for i := 0; i < b.N; i++ {
		gg.Nop1(ref.Load())
	}
}
