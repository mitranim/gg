package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

type IniterStr string

func (self *IniterStr) Init() {
	if *self != `` {
		panic(`redundant init`)
	}
	*self = `inited`
}

// Incomplete: doesn't verify concurrency safety.
func TestLazyIniter(t *testing.T) {
	defer gtest.Catch(t)

	var tar gg.LazyIniter[IniterStr, *IniterStr]

	test := func(exp gg.Opt[IniterStr]) {
		// Note: `gg.CastUnsafe[A](tar)` would have been simpler, but technically
		// involves passing `tar` by value, which is invalid due to inner mutex.
		// Wouldn't actually matter.
		gtest.Eq(*gg.CastUnsafe[*gg.Opt[IniterStr]](&tar), exp)
	}

	test(gg.Zero[gg.Opt[IniterStr]]())

	gtest.Eq(tar.Get(), `inited`)
	gtest.Eq(tar.Get(), tar.Get())
	test(gg.OptVal(IniterStr(`inited`)))

	tar.Clear()
	test(gg.Zero[gg.Opt[IniterStr]]())

	gtest.Eq(tar.Get(), `inited`)
	gtest.Eq(tar.Get(), tar.Get())
	test(gg.OptVal(IniterStr(`inited`)))

	tar.Reset(`inited_1`)
	gtest.Eq(tar.Get(), `inited_1`)
	gtest.Eq(tar.Get(), tar.Get())
	test(gg.OptVal(IniterStr(`inited_1`)))
}
