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
			testPeekerZero[IniterStr](&mem)
			gtest.Eq(mem.Get(), `inited`)
			gtest.Eq(mem.Get(), `inited`)
			gtest.Eq(mem.Peek(), `inited`)
			gtest.Eq(mem.Peek(), `inited`)
			testPeekerNotZero[IniterStr](&mem)
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
			testPeekerZero[IniterStrs](&mem)

			prev := mem.Get()
			gtest.Equal(prev, []string{`one`, `two`, `three`})
			gtest.SliceIs(mem.Get(), prev)
			gtest.SliceIs(mem.Get(), prev)

			testPeekerNotZero[IniterStrs](&mem)
		}

		test()
		prev := mem.Get()
		mem.Clear()
		test()
		next := mem.Get()
		gtest.NotSliceIs(next, prev)
	}
}

func testPeekerZero[Tar any, Src gg.Peeker[Tar]](src Src) {
	gtest.Zero(src.Peek())
	gtest.Zero(src.Peek())
}

func testPeekerNotZero[Tar any, Src gg.Peeker[Tar]](src Src) {
	gtest.NotZero(src.Peek())
	gtest.NotZero(src.Peek())
}
