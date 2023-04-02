package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

type IniterOnceStr string

func (self *IniterOnceStr) Init() {
	if *self == `` {
		*self = `inited`
		return
	}
	panic(`redundant init`)
}

// Incomplete: doesn't verify concurrency safety.
func TestLazyIniter(t *testing.T) {
	defer gtest.Catch(t)

	var tar gg.LazyIniter[IniterOnceStr, *IniterOnceStr]

	gtest.Eq(tar.Get(), `inited`)
	gtest.Eq(tar.Get(), tar.Get())

	// Note: `gg.CastUnsafe[string](tar)` would have been simpler, but
	// technically involves passing `tar` by value, which is invalid due
	// to inner mutex. Wouldn't actually matter.
	gtest.Eq(*gg.CastUnsafe[*string](&tar), `inited`)

	gtest.Eq(tar.Ptr(), gg.CastUnsafe[*IniterOnceStr](&tar))
	gtest.Eq(tar.Ptr(), tar.Ptr())
}
