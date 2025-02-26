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

/*
Technical note: because `Mem` contains a synchronization primitive, its method
`.MarshalJSON` cannot be safely implemented on the value type, and we cannot
safely pass the value to JSON-encoding functions. This is a known inconsistency
with other types that implement JSON encoding.
*/
func TestMem_MarshalJSON(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`primitive`, func(t *testing.T) {
		defer gtest.Catch(t)

		type Type = gg.Mem[gg.DurSecond, IniterStr, *IniterStr]

		gtest.Eq(gg.ToString(gg.Try1((*Type)(nil).MarshalJSON())), `null`)
		gtest.Eq(gg.JsonString((*Type)(nil)), `null`)

		var mem Type
		gtest.Eq(gg.JsonString(&mem), `null`)

		mem.Get()
		gtest.Eq(gg.JsonString(&mem), `"inited"`)

		mem.Clear()
		gtest.Eq(gg.JsonString(&mem), `null`)
	})

	t.Run(`slice`, func(t *testing.T) {
		defer gtest.Catch(t)

		type Type = gg.Mem[gg.DurSecond, IniterStrs, *IniterStrs]

		gtest.Eq(gg.ToString(gg.Try1((*Type)(nil).MarshalJSON())), `null`)
		gtest.Eq(gg.JsonString((*Type)(nil)), `null`)

		var mem Type
		gtest.Eq(gg.JsonString(&mem), `null`)

		mem.Get()
		gtest.Eq(gg.JsonString(&mem), `["one","two","three"]`)

		mem.Clear()
		gtest.Eq(gg.JsonString(&mem), `null`)
	})
}

func TestMem_UnmarshalJSON(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`primitive`, func(t *testing.T) {
		defer gtest.Catch(t)

		var mem gg.Mem[gg.DurSecond, IniterStr, *IniterStr]
		gtest.NoErr(mem.UnmarshalJSON([]byte(`null`)))
		gtest.Zero(mem.Peek())

		gtest.NoErr(mem.UnmarshalJSON([]byte(`""`)))
		gtest.Eq(mem.Peek(), ``)

		gtest.NoErr(mem.UnmarshalJSON([]byte(`"test"`)))
		gtest.Eq(mem.Peek(), `test`)

		gtest.ErrStr(`invalid character 'i'`, mem.UnmarshalJSON([]byte(`invalid`)))
		gtest.Eq(mem.Peek(), `test`)
	})

	t.Run(`slice`, func(t *testing.T) {
		defer gtest.Catch(t)

		var mem gg.Mem[gg.DurSecond, IniterStrs, *IniterStrs]
		gtest.NoErr(mem.UnmarshalJSON([]byte(`null`)))
		gtest.Zero(mem.Peek())

		gtest.NoErr(mem.UnmarshalJSON([]byte(`[]`)))
		gtest.Equal(mem.Peek(), IniterStrs{})

		gtest.NoErr(mem.UnmarshalJSON([]byte(`["four","five","six"]`)))
		gtest.Equal(mem.Peek(), IniterStrs{`four`, `five`, `six`})

		gtest.ErrStr(`invalid character 'i'`, mem.UnmarshalJSON([]byte(`invalid`)))
		gtest.Equal(mem.Peek(), IniterStrs{`four`, `five`, `six`})
	})
}

func testPeekerZero[Tar any, Src gg.Peeker[Tar]](src Src) {
	gtest.Zero(src.Peek())
	gtest.Zero(src.Peek())
}

func testPeekerNotZero[Tar any, Src gg.Peeker[Tar]](src Src) {
	gtest.NotZero(src.Peek())
	gtest.NotZero(src.Peek())
}
