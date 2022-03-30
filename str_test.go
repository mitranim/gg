package gg_test

import (
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestStrDat(t *testing.T) {
	defer gtest.Catch(t)

	const init = `hello world`
	var sliced string = init[:0]
	var empty string

	gtest.NotZero(init)
	gtest.Zero(sliced)
	gtest.Zero(empty)

	gtest.NotZero(gg.StrDat(init))
	gtest.NotZero(gg.StrDat(sliced))
	gtest.Zero(gg.StrDat(empty))

	gtest.Eq(gg.StrDat(sliced), gg.StrDat(init))
}

// TODO dedup with `TestBuf_String`.
func TestToString(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.ToString([]byte(nil)), ``)

	test := func(str string) {
		t.Helper()

		src := []byte(str)
		tar := gg.ToString(src)

		gtest.Eq(tar, str)
		gtest.Eq(gg.StrDat(src), gg.StrDat(tar))
	}

	test(``)
	test(`a`)
	test(`ab`)
	test(`abc`)

	t.Run(`mutation`, func(t *testing.T) {
		src := []byte(`abc`)
		tar := gg.ToString(src)
		gtest.Eq(tar, `abc`)

		src[0] = 'd'
		gtest.Eq(tar, `dbc`)
	})
}

func TestToBytes(t *testing.T) {
	defer gtest.Catch(t)

	src := `abc`
	tar := gg.ToBytes(src)

	gtest.Eq(string(tar), `abc`)
	gtest.Eq(gg.StrDat(src), gg.StrDat(tar))
}

func TestStrPop(t *testing.T) {
	defer gtest.Catch(t)

	rem := `{one,two,,three}`

	gtest.Eq(gg.StrPop(&rem, `,`), `{one`)
	gtest.Eq(rem, `two,,three}`)

	gtest.Eq(gg.StrPop(&rem, `,`), `two`)
	gtest.Eq(rem, `,three}`)

	gtest.Eq(gg.StrPop(&rem, `,`), ``)
	gtest.Eq(rem, `three}`)

	gtest.Eq(gg.StrPop(&rem, `,`), `three}`)
	gtest.Eq(rem, ``)
}

// TODO test Unicode.
func TestToWords(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.ToWords(`one two three`), gg.Words{`one`, `two`, `three`})
	gtest.Equal(gg.ToWords(`One Two Three`), gg.Words{`One`, `Two`, `Three`})
	gtest.Equal(gg.ToWords(`OneTwoThree`), gg.Words{`One`, `Two`, `Three`})
	gtest.Equal(gg.ToWords(`One-Two-Three`), gg.Words{`One`, `Two`, `Three`})
	gtest.Equal(gg.ToWords(`One_Two_Three`), gg.Words{`One`, `Two`, `Three`})
}

// TODO test Unicode.
func TestWords(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Words{`one`, `two`, `three`}.Spaced(), `one two three`)
	gtest.Eq(gg.Words{`one`, `two`, `three`}.Snake(), `one_two_three`)
	gtest.Eq(gg.Words{`one`, `two`, `three`}.Kebab(), `one-two-three`)
	gtest.Eq(gg.Words{`one`, `two`, `three`}.Solid(), `onetwothree`)

	gtest.Equal(gg.Words{`ONE`, `TWO`, `THREE`}.Lower(), gg.Words{`one`, `two`, `three`})
	gtest.Equal(gg.Words{`one`, `two`, `three`}.Upper(), gg.Words{`ONE`, `TWO`, `THREE`})
	gtest.Equal(gg.Words{`one`, `two`, `three`}.Title(), gg.Words{`One`, `Two`, `Three`})
	gtest.Equal(gg.Words{`ONE`, `TWO`, `THREE`}.Sentence(), gg.Words{`One`, `two`, `three`})
	gtest.Equal(gg.Words{`one`, `two`, `three`}.Camel(), gg.Words{`one`, `Two`, `Three`})
}

func Benchmark_strings_Join(b *testing.B) {
	val := gg.Map(gg.Span(128), gg.String[int])
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gg.Nop1(strings.Join(val, ` `))
	}
}

func BenchmarkJoinOpt(b *testing.B) {
	val := gg.Map(gg.Span(128), gg.String[int])
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.JoinOpt(val, ` `))
	}
}
