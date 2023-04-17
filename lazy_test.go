package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

// TODO: test with concurrency.
func TestLazy(t *testing.T) {
	defer gtest.Catch(t)

	var count int

	once := gg.NewLazy(func() int {
		count++
		if count > 1 {
			panic(gg.Errf(`excessive count %v`, count))
		}
		return count
	})

	gtest.Eq(*gg.CastUnsafe[*int](once), 0)
	gtest.Eq(once.Get(), 1)
	gtest.Eq(once.Get(), once.Get())
	gtest.Eq(*gg.CastUnsafe[*int](once), 1)
}

func BenchmarkLazy(b *testing.B) {
	once := gg.NewLazy(gg.Cwd)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(once.Get())
	}
}
