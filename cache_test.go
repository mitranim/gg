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

	once := gg.Lazy(func() int {
		count++
		if count > 1 {
			panic(gg.Errf(`excessive count %v`, count))
		}
		return count
	})

	gtest.NoPanic(func() {
		gtest.Eq(once(), 1)
		gtest.Eq(once(), 1)
		gtest.Eq(once(), 1)
		gtest.Eq(once(), 1)
	})
}

func BenchmarkLazy(b *testing.B) {
	once := gg.Lazy(gg.Cwd)

	for i := 0; i < b.N; i++ {
		gg.Nop1(once())
	}
}
