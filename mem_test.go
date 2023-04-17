package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

type IniterStrs []string

func (self *IniterStrs) Init() {
	if *self != nil {
		panic(`redundant init`)
	}
	*self = []string{`one`, `two`, `three`}
}

func TestMem(t *testing.T) {
	defer gtest.Catch(t)

	{
		// Note: `IniterStr` panics if init is called on a non-zero value.
		// This is part of our contract.
		var mem gg.Mem[gg.DurSecond, IniterStr, *IniterStr]

		test := func() {
			testMemZero[IniterStr](&mem)
			gtest.Eq(mem.Get(), `inited`)
			gtest.Eq(mem.Get(), `inited`)
			gtest.Eq(mem.Peek(), `inited`)
			gtest.Eq(mem.Peek(), `inited`)
			testMemNotZero[IniterStr](&mem)
		}

		test()
		mem.Clear()
		test()
	}

	{
		// Note: `IniterStrs` panics if init is called on a non-zero value.
		// This is part of our contract.
		var mem gg.Mem[gg.DurSecond, IniterStrs, *IniterStrs]

		test := func() {
			testMemZero[IniterStrs](&mem)

			prev := mem.Get()
			gtest.Equal(prev, []string{`one`, `two`, `three`})
			gtest.SliceIs(mem.Get(), prev)
			gtest.SliceIs(mem.Get(), prev)
			gtest.NotZero(mem.Timed().Inst)
			gtest.NotZero(mem.Timed().Inst)

			testMemNotZero[IniterStrs](&mem)
		}

		test()
		prev := mem.Get()
		mem.Clear()
		test()
		next := mem.Get()
		gtest.NotSliceIs(next, prev)
	}
}

type Peeker[Tar any] interface {
	gg.Peeker[Tar]
	PeekTimed() gg.Timed[Tar]
}

func testMemZero[Tar any, Mem Peeker[Tar]](mem Mem) {
	gtest.Zero(mem.Peek())
	gtest.Zero(mem.Peek())
	gtest.Zero(mem.PeekTimed())
	gtest.Zero(mem.PeekTimed())
}

func testMemNotZero[Tar any, Mem Peeker[Tar]](mem Mem) {
	gtest.NotZero(mem.Peek())
	gtest.NotZero(mem.Peek())
	gtest.NotZero(mem.PeekTimed())
	gtest.NotZero(mem.PeekTimed())
	gtest.NotZero(mem.PeekTimed().Inst)
	gtest.NotZero(mem.PeekTimed().Inst)
}
