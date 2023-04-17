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

	gtest.Eq(tar.Get(), `inited`)
	gtest.Eq(tar.Get(), tar.Get())

	// Note: `gg.CastUnsafe[string](tar)` would have been simpler, but
	// technically involves passing `tar` by value, which is invalid due
	// to inner mutex. Wouldn't actually matter.
	gtest.Eq(*gg.CastUnsafe[*string](&tar), `inited`)

	gtest.Eq(tar.Ptr(), gg.CastUnsafe[*IniterStr](&tar))
	gtest.Eq(tar.Ptr(), tar.Ptr())
}
