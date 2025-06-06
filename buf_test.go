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
		gtest.Eq(gg.TextDat(buf), gg.TextDat(tar))
	}

	test(``)
	test(`a`)
	test(`ab`)
	test(`abc`)

	t.Run(`mutation`, func(t *testing.T) {
		defer gtest.Catch(t)

		buf := gg.Buf(`abc`)
		tar := buf.String()
		gtest.Eq(tar, `abc`)

		buf[0] = 'd'
		gtest.Eq(tar, `dbc`)
	})
}

func TestBuf_AppendAnys(t *testing.T) {
	defer gtest.Catch(t)

	var buf gg.Buf
	gtest.Zero(buf)

	buf.AppendAnys()
	gtest.Zero(buf)

	buf.AppendAnys(nil)
	gtest.Zero(buf)

	buf.AppendAnys(``, nil, ``)
	gtest.Zero(buf)

	buf.AppendAnys(10)
	gtest.Str(buf, `10`)

	buf.AppendAnys(` `, 20, ` `, 30)
	gtest.Str(buf, `10 20 30`)
}

func TestBuf_AppendAnysln(t *testing.T) {
	defer gtest.Catch(t)

	{
		var buf gg.Buf
		gtest.Zero(buf)
	}

	{
		var buf gg.Buf
		buf.AppendAnysln()
		gtest.Zero(buf)
	}

	{
		var buf gg.Buf
		buf.AppendAnysln(nil)
		gtest.Zero(buf)
	}

	{
		var buf gg.Buf
		buf.AppendAnysln(nil, ``, nil)
		gtest.Zero(buf)

		buf.AppendAnysln()
		gtest.Zero(buf)
	}

	{
		var buf gg.Buf
		buf.AppendAnysln(`one`, `two`+"\n")
		gtest.Str(buf, `onetwo`+"\n")
	}

	{
		var buf gg.Buf
		buf.AppendAnysln(`one`+"\n", `two`+"\n")
		gtest.Str(buf, `one`+"\n"+`two`+"\n")
	}

	{
		var buf gg.Buf
		buf.AppendAnysln(`one`+"\n", `two`)
		gtest.Str(buf, `one`+"\n"+`two`+"\n")
	}

	{
		var buf gg.Buf
		buf.AppendAnysln(`one`)
		gtest.Str(buf, `one`+"\n")

		buf.AppendAnysln()
		gtest.Str(buf, `one`+"\n\n")

		buf.AppendAnysln(`two`)
		gtest.Str(buf, `one`+"\n\n"+`two`+"\n")

		buf.AppendAnysln()
		gtest.Str(buf, `one`+"\n\n"+`two`+"\n\n")
	}
}
