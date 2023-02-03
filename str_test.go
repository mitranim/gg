package gg_test

import (
	"regexp"
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

func TestToText(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.ToText[string](``), ``)
	gtest.Eq(gg.ToText[string](`str`), `str`)
	gtest.Equal(gg.ToText[[]byte](``), []byte(nil))
	gtest.Equal(gg.ToText[[]byte](`str`), []byte(`str`))
	gtest.Equal(gg.ToText[string]([]byte(`str`)), `str`)
}

// TODO dedup with `TestBuf_String`.
func TestToString(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.ToString([]byte(nil)), ``)

	test := func(str string) {
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
		defer gtest.Catch(t)

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

func TestJoinAny(t *testing.T) {
	gtest.Catch(t)

	gtest.Zero(gg.JoinAny(nil, ``))
	gtest.Zero(gg.JoinAny([]any{}, ``))
	gtest.Zero(gg.JoinAny([]any{}, `_`))

	gtest.Zero(gg.JoinAny([]any{``}, ``))
	gtest.Zero(gg.JoinAny([]any{``, ``}, ``))
	gtest.Zero(gg.JoinAny([]any{``, ``, ``}, ``))

	gtest.Eq(gg.JoinAny([]any{``}, `_`), ``)
	gtest.Eq(gg.JoinAny([]any{``, ``}, `_`), `_`)
	gtest.Eq(gg.JoinAny([]any{``, ``, ``}, `_`), `__`)

	gtest.Eq(gg.JoinAny([]any{12}, ``), `12`)
	gtest.Eq(gg.JoinAny([]any{12}, `_`), `12`)

	gtest.Eq(gg.JoinAny([]any{12, 34}, ``), `1234`)
	gtest.Eq(gg.JoinAny([]any{12, 34}, `_`), `12_34`)

	gtest.Eq(gg.JoinAny([]any{12, 34, 56}, ``), `123456`)
	gtest.Eq(gg.JoinAny([]any{12, 34, 56}, `_`), `12_34_56`)

	gtest.Eq(gg.JoinAny([]any{12, `str`}, ``), `12str`)
	gtest.Eq(gg.JoinAny([]any{12, `str`}, `_`), `12_str`)

	gtest.Eq(gg.JoinAny([]any{`one`, ``, `two`, ``, `three`}, ``), `onetwothree`)
	gtest.Eq(gg.JoinAny([]any{`one`, ``, `two`, ``, `three`}, `_`), `one__two__three`)
}

func BenchmarkJoinAny(b *testing.B) {
	src := gg.Map(gg.Span(128), gg.ToAny[int])
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.JoinAny(src, ` `))
	}
}

func TestJoinAnyOpt(t *testing.T) {
	gtest.Catch(t)

	gtest.Zero(gg.JoinAnyOpt(nil, ``))
	gtest.Zero(gg.JoinAnyOpt([]any{}, ``))
	gtest.Zero(gg.JoinAnyOpt([]any{}, `_`))

	gtest.Zero(gg.JoinAnyOpt([]any{``}, ``))
	gtest.Zero(gg.JoinAnyOpt([]any{``, ``}, ``))
	gtest.Zero(gg.JoinAnyOpt([]any{``, ``, ``}, ``))

	gtest.Zero(gg.JoinAnyOpt([]any{``}, `_`))
	gtest.Zero(gg.JoinAnyOpt([]any{``, ``}, `_`))
	gtest.Zero(gg.JoinAnyOpt([]any{``, ``, ``}, `_`))

	gtest.Eq(gg.JoinAnyOpt([]any{12}, ``), `12`)
	gtest.Eq(gg.JoinAnyOpt([]any{12}, `_`), `12`)

	gtest.Eq(gg.JoinAnyOpt([]any{12, 34}, ``), `1234`)
	gtest.Eq(gg.JoinAnyOpt([]any{12, 34}, `_`), `12_34`)

	gtest.Eq(gg.JoinAnyOpt([]any{12, 34, 56}, ``), `123456`)
	gtest.Eq(gg.JoinAnyOpt([]any{12, 34, 56}, `_`), `12_34_56`)

	gtest.Eq(gg.JoinAnyOpt([]any{12, `str`}, ``), `12str`)
	gtest.Eq(gg.JoinAnyOpt([]any{12, `str`}, `_`), `12_str`)

	gtest.Eq(gg.JoinAnyOpt([]any{`one`, ``, `two`, ``, `three`}, ``), `onetwothree`)
	gtest.Eq(gg.JoinAnyOpt([]any{`one`, ``, `two`, ``, `three`}, `_`), `one_two_three`)
}

func BenchmarkJoinAnyOpt(b *testing.B) {
	src := gg.Map(gg.Span(128), gg.ToAny[int])
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.JoinAnyOpt(src, ` `))
	}
}

func Benchmark_strings_Join(b *testing.B) {
	src := gg.Map(gg.Span(128), gg.String[int])
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(strings.Join(src, ` `))
	}
}

func BenchmarkJoin(b *testing.B) {
	src := gg.Map(gg.Span(128), gg.String[int])
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Join(src, ` `))
	}
}

func BenchmarkJoinOpt(b *testing.B) {
	src := gg.Map(gg.Span(128), gg.String[int])
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.JoinOpt(src, ` `))
	}
}

// TODO test Unicode.
func TestToWords(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src string, exp gg.Words) { gtest.Equal(gg.ToWords(src), exp) }

	test(``, nil)
	test(` `, nil)

	test(`one`, gg.Words{`one`})
	test(`one two`, gg.Words{`one`, `two`})
	test(`one two three`, gg.Words{`one`, `two`, `three`})
	test(`one  two  three`, gg.Words{`one`, `two`, `three`})
	test(`One Two Three`, gg.Words{`One`, `Two`, `Three`})
	test(`ONE TWO THREE`, gg.Words{`ONE`, `TWO`, `THREE`})
	test(`one12 two34 three56`, gg.Words{`one12`, `two34`, `three56`})
	test(`One12 Two34 Three56`, gg.Words{`One12`, `Two34`, `Three56`})
	test(`ONE12 TWO34 THREE56`, gg.Words{`ONE12`, `TWO34`, `THREE56`})

	test(`one_two_three`, gg.Words{`one`, `two`, `three`})
	test(`one_Two_Three`, gg.Words{`one`, `Two`, `Three`})
	test(`One_Two_Three`, gg.Words{`One`, `Two`, `Three`})
	test(`ONE_TWO_THREE`, gg.Words{`ONE`, `TWO`, `THREE`})
	test(`one12_two34_three56`, gg.Words{`one12`, `two34`, `three56`})
	test(`one12_Two34_Three56`, gg.Words{`one12`, `Two34`, `Three56`})
	test(`One12_Two34_Three56`, gg.Words{`One12`, `Two34`, `Three56`})
	test(`ONE12_TWO34_THREE56`, gg.Words{`ONE12`, `TWO34`, `THREE56`})

	test(`oneTwoThree`, gg.Words{`one`, `Two`, `Three`})
	test(`OneTwoThree`, gg.Words{`One`, `Two`, `Three`})
	test(`one12Two34Three56`, gg.Words{`one12`, `Two34`, `Three56`})
	test(`One12Two34Three56`, gg.Words{`One12`, `Two34`, `Three56`})

	test(`one-two-three`, gg.Words{`one`, `two`, `three`})
	test(`one-Two-Three`, gg.Words{`one`, `Two`, `Three`})
	test(`One-Two-Three`, gg.Words{`One`, `Two`, `Three`})
	test(`ONE-TWO-THREE`, gg.Words{`ONE`, `TWO`, `THREE`})
	test(`one12-two34-three56`, gg.Words{`one12`, `two34`, `three56`})
	test(`one12-Two34-Three56`, gg.Words{`one12`, `Two34`, `Three56`})
	test(`One12-Two34-Three56`, gg.Words{`One12`, `Two34`, `Three56`})
	test(`ONE12-TWO34-THREE56`, gg.Words{`ONE12`, `TWO34`, `THREE56`})
}

// TODO test Unicode.
func TestWords(t *testing.T) {
	defer gtest.Catch(t)

	src := func() gg.Words { return gg.Words{`one`, `two`, `three`} }

	gtest.Equal(gg.ToWords(`one two`).Lower(), gg.Words{`one`, `two`})
	gtest.Equal(gg.ToWords(`One Two`).Lower(), gg.Words{`one`, `two`})
	gtest.Equal(gg.ToWords(`ONE TWO`).Lower(), gg.Words{`one`, `two`})

	gtest.Eq(src().Spaced(), `one two three`)
	gtest.Eq(src().Snake(), `one_two_three`)
	gtest.Eq(src().Kebab(), `one-two-three`)
	gtest.Eq(src().Dense(), `onetwothree`)

	gtest.Equal(src().Lower(), gg.Words{`one`, `two`, `three`})
	gtest.Equal(src().Upper(), gg.Words{`ONE`, `TWO`, `THREE`})
	gtest.Equal(src().Title(), gg.Words{`One`, `Two`, `Three`})
	gtest.Equal(src().Sentence(), gg.Words{`One`, `two`, `three`})
	gtest.Equal(src().Camel(), gg.Words{`one`, `Two`, `Three`})

	gtest.Eq(src().Lower().Spaced(), `one two three`)
	gtest.Eq(src().Lower().Snake(), `one_two_three`)
	gtest.Eq(src().Lower().Kebab(), `one-two-three`)
	gtest.Eq(src().Lower().Dense(), `onetwothree`)

	gtest.Eq(src().Upper().Spaced(), `ONE TWO THREE`)
	gtest.Eq(src().Upper().Snake(), `ONE_TWO_THREE`)
	gtest.Eq(src().Upper().Kebab(), `ONE-TWO-THREE`)
	gtest.Eq(src().Upper().Dense(), `ONETWOTHREE`)

	gtest.Eq(src().Title().Spaced(), `One Two Three`)
	gtest.Eq(src().Title().Snake(), `One_Two_Three`)
	gtest.Eq(src().Title().Kebab(), `One-Two-Three`)
	gtest.Eq(src().Title().Dense(), `OneTwoThree`)

	gtest.Eq(src().Sentence().Spaced(), `One two three`)
	gtest.Eq(src().Sentence().Snake(), `One_two_three`)
	gtest.Eq(src().Sentence().Kebab(), `One-two-three`)
	gtest.Eq(src().Sentence().Dense(), `Onetwothree`)

	gtest.Eq(src().Camel().Spaced(), `one Two Three`)
	gtest.Eq(src().Camel().Snake(), `one_Two_Three`)
	gtest.Eq(src().Camel().Kebab(), `one-Two-Three`)
	gtest.Eq(src().Camel().Dense(), `oneTwoThree`)
}

func BenchmarkReWord_init(b *testing.B) {
	src := gg.ReWord().String()

	for ind := 0; ind < b.N; ind++ {
		regexp.MustCompile(src)
	}
}

func BenchmarkReWord_reuse(b *testing.B) {
	gg.Nop1(gg.ReWord())

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ReWord())
	}
}

func TestSplitLines(t *testing.T) {
	defer gtest.Catch(t)

	var Split = gg.SplitLines[string]

	gtest.Zero(Split(""))
	gtest.Equal(Split(" "), []string{` `})
	gtest.Equal(Split("\n"), []string{``, ``})
	gtest.Equal(Split("one"), []string{`one`})
	gtest.Equal(Split("one\n"), []string{`one`, ``})
	gtest.Equal(Split("\none"), []string{``, `one`})
	gtest.Equal(Split("\none\n"), []string{``, `one`, ``})
	gtest.Equal(Split("one\ntwo"), []string{`one`, `two`})
	gtest.Equal(Split("one\ntwo\n"), []string{`one`, `two`, ``})
	gtest.Equal(Split("\none\ntwo"), []string{``, `one`, `two`})
	gtest.Equal(Split("\none\ntwo\n"), []string{``, `one`, `two`, ``})
}
