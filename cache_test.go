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

type MemTest struct {
	gg.DurMinute
	gg.Mem[[]string]
}

func (self *MemTest) Get() []string { return self.DedupFrom(self) }
func (*MemTest) Make() []string     { return []string{`str`} }

func TestMem(t *testing.T) {
	defer gtest.Catch(t)

	var mem MemTest

	prev := mem.Get()
	gtest.Equal(prev, []string{`str`})
	gtest.SliceIs(prev, mem.Get())
	gtest.SliceIs(prev, mem.Get())

	gtest.NotZero(mem.Timed.Inst)
	mem.Clear()
	gtest.Zero(mem.Timed)

	next := mem.Get()
	gtest.Equal(next, []string{`str`})
	gtest.SliceIs(next, mem.Get())
	gtest.SliceIs(next, mem.Get())

	gtest.Equal(prev, next)
	gtest.NotSliceIs(prev, next)
	gtest.NotZero(mem.Timed.Inst)
}
