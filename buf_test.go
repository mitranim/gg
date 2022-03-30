package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

// TODO dedup with `TestToString`.
func TestBuf_String(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.Buf(nil).String())

	test := func(src string) {
		buf := gg.Buf(src)
		tar := buf.String()

		gtest.Eq(tar, src)
		gtest.Eq(gg.StrDat(buf), gg.StrDat(tar))
	}

	test(``)
	test(`a`)
	test(`ab`)
	test(`abc`)

	t.Run(`mutation`, func(t *testing.T) {
		buf := gg.Buf(`abc`)
		tar := buf.String()
		gtest.Eq(tar, `abc`)

		buf[0] = 'd'
		gtest.Eq(tar, `dbc`)
	})
}
