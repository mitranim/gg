package gg_test

import (
	"math"
	"regexp"
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestTextDat(t *testing.T) {
	defer gtest.Catch(t)

	const init = `hello world`
	var sliced string = init[:0]
	var empty string

	gtest.NotZero(init)
	gtest.Zero(sliced)
	gtest.Zero(empty)

	gtest.NotZero(gg.TextDat(init))
	gtest.NotZero(gg.TextDat(sliced))
	gtest.Zero(gg.TextDat(empty))

	gtest.Eq(gg.TextDat(sliced), gg.TextDat(init))
}

func TestToText(t *testing.T) {
	defer gtest.Catch(t)

	testToString(gg.ToText[string, []byte])
	testToBytes(gg.ToText[[]byte, string])

	t.Run(`between_byte_slices`, func(t *testing.T) {
		defer gtest.Catch(t)

		type Src []byte

		src := Src(`one_two`)[:len(`one`)]
		gtest.TextEq(src, Src(`one`))
		gtest.Len(src, 3)
		gtest.Cap(src, 7)

		type Tar []byte

		gtest.Is(Src(gg.ToText[Tar](src)), src)
		gtest.Is(gg.ToText[Src](src), src)
		gtest.Is(gg.ToText[[]byte](src), []byte(src))
		gtest.Is([]byte(gg.ToText[Tar](src)), []byte(src))
	})
}

func testToString(fun func([]byte) string) {
	test := func(src []byte) {
		tar := fun(src)
		gtest.Eq(string(tar), string(src))
		gtest.Eq(gg.TextDat(src), gg.TextDat(tar))
	}

	test(nil)
	test([]byte{})
	test([]byte(`one`))
	test([]byte(`two`))
	test([]byte(`three`))
}

func testToBytes(fun func(string) []byte) {
	test := func(src string) {
		tar := fun(src)
		gtest.Eq(string(tar), string(src))
		gtest.Eq(gg.TextDat(src), gg.TextDat(tar))
		gtest.Len(tar, len(src))
		gtest.Cap(tar, len(src))
	}

	test(``)
	test(`one`)
	test(`two`)
	test(`three`)
}

func BenchmarkToText_string_to_string(b *testing.B) {
	defer gtest.Catch(b)

	type Src string
	type Out string
	src := Src(`742af97969e845408f0261c213a4c01f`)

	for ind := 0; ind < b.N; ind++ {
		gg.ToText[Out](src)
	}
}

func BenchmarkToText_string_to_bytes(b *testing.B) {
	defer gtest.Catch(b)

	type Src string
	type Out []byte
	src := Src(`742af97969e845408f0261c213a4c01f`)

	for ind := 0; ind < b.N; ind++ {
		gg.ToText[Out](src)
	}
}

func BenchmarkToText_bytes_to_string(b *testing.B) {
	defer gtest.Catch(b)

	type Src []byte
	type Out string
	src := Src(`742af97969e845408f0261c213a4c01f`)

	for ind := 0; ind < b.N; ind++ {
		gg.ToText[Out](src)
	}
}

func BenchmarkToText_bytes_to_bytes(b *testing.B) {
	defer gtest.Catch(b)

	type Src []byte
	type Out []byte
	src := Src(`742af97969e845408f0261c213a4c01f`)

	for ind := 0; ind < b.N; ind++ {
		gg.ToText[Out](src)
	}
}

// TODO dedup with `TestBuf_String`.
func TestToString(t *testing.T) {
	defer gtest.Catch(t)

	testToString(gg.ToString[[]byte])

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

	testToBytes(gg.ToBytes[string])
}

func TestTextPop(t *testing.T) {
	defer gtest.Catch(t)

	rem := `{one,two,,three}`

	gtest.Eq(gg.TextPop(&rem, `,`), `{one`)
	gtest.Eq(rem, `two,,three}`)

	gtest.Eq(gg.TextPop(&rem, `,`), `two`)
	gtest.Eq(rem, `,three}`)

	gtest.Eq(gg.TextPop(&rem, `,`), ``)
	gtest.Eq(rem, `three}`)

	gtest.Eq(gg.TextPop(&rem, `,`), `three}`)
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
	src := gg.ReWord.Get().String()

	for ind := 0; ind < b.N; ind++ {
		regexp.MustCompile(src)
	}
}

func BenchmarkReWord_reuse(b *testing.B) {
	gg.Nop1(gg.ReWord.Get())

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ReWord.Get())
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

func TestSplitLines2(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src, oneExp, twoExp string, sizeExp int) {
		one, two, size := gg.SplitLines2(src)
		gtest.Eq([2]string{one, two}, [2]string{oneExp, twoExp})
		gtest.Eq(size, sizeExp)
	}

	test(``, ``, ``, 0)
	test("one\n", `one`, ``, 1)
	test("one\ntwo", `one`, `two`, 1)
	test("one\r\ntwo", `one`, `two`, 2)
	test("one\rtwo", `one`, `two`, 1)
	test("\ntwo", ``, `two`, 1)
	test("\r\ntwo", ``, `two`, 2)
	test("\rtwo", ``, `two`, 1)
	test("one\ntwo\nthree", `one`, "two\nthree", 1)
}

func TestTextCut_ours(t *testing.T) {
	defer gtest.Catch(t)
	testTextCut(gg.TextCut[string])
}

func TestTextCut_alternate(t *testing.T) {
	defer gtest.Catch(t)
	testTextCut(TextCutRuneSlice[string])
}

func testTextCut(fun func(string, int, int) string) {
	const src = `🐒🐴🦖🦔🐲🐈`

	gtest.Zero(fun(src, 0, 0))
	gtest.Eq(fun(src, 0, 1), `🐒`)
	gtest.Eq(fun(src, 0, 2), `🐒🐴`)
	gtest.Eq(fun(src, 0, 3), `🐒🐴🦖`)
	gtest.Eq(fun(src, 0, 4), `🐒🐴🦖🦔`)
	gtest.Eq(fun(src, 0, 5), `🐒🐴🦖🦔🐲`)
	gtest.Eq(fun(src, 0, 6), `🐒🐴🦖🦔🐲🐈`)
	gtest.Eq(fun(src, 0, 7), `🐒🐴🦖🦔🐲🐈`)
	gtest.Eq(fun(src, 0, 8), `🐒🐴🦖🦔🐲🐈`)

	gtest.Eq(fun(src, -1, 1), `🐒`)
	gtest.Eq(fun(src, -1, 6), `🐒🐴🦖🦔🐲🐈`)

	gtest.Zero(fun(src, 1, 0))
	gtest.Zero(fun(src, 1, 1))
	gtest.Eq(fun(src, 1, 2), `🐴`)
	gtest.Eq(fun(src, 1, 6), `🐴🦖🦔🐲🐈`)

	gtest.Eq(fun(`one two three four`, 4, 13), `two three`)
}

// Alternate implementation for comparison with ours.
func TextCutRuneSlice[A ~string](src A, start, end int) A {
	if !(end > start) {
		return ``
	}

	runes := []rune(src)
	size := len(runes)

	if start < 0 {
		start = 0
	} else if start > size {
		start = size
	}
	if end < 0 {
		end = 0
	} else if end > size {
		end = size
	}

	return A(runes[start:end])
}

func BenchmarkTextCut_ours(b *testing.B) {
	const src = `🐒🐴🦖🦔🐲🐈`

	for ind := 0; ind < b.N; ind++ {
		for start := range gg.Iter(3) {
			start--
			for end := range gg.Iter(6) {
				gg.Nop1(gg.TextCut(src, start, end))
			}
		}
	}
}

func BenchmarkTextCut_alternate(b *testing.B) {
	const src = `🐒🐴🦖🦔🐲🐈`

	for ind := 0; ind < b.N; ind++ {
		for start := range gg.Iter(3) {
			start--
			for end := range gg.Iter(6) {
				gg.Nop1(TextCutRuneSlice(src, start, end))
			}
		}
	}
}

func TestStr(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.Str())
	gtest.Zero(gg.Str(nil))
	gtest.Zero(gg.Str(nil, nil))
	gtest.Zero(gg.Str(``))
	gtest.Zero(gg.Str(``, nil))
	gtest.Zero(gg.Str(``, nil, ``))
	gtest.Zero(gg.Str(``, nil, ``, nil))

	gtest.Eq(gg.Str(0), `0`)
	gtest.Eq(gg.Str(0, 0), `00`)
	gtest.Eq(gg.Str(0, 10), `010`)
	gtest.Eq(gg.Str(0, 10, 20), `01020`)
	gtest.Eq(gg.Str(`one`), `one`)
	gtest.Eq(gg.Str(`one`, ``), `one`)
	gtest.Eq(gg.Str(`one`, ``, `two`), `onetwo`)
	gtest.Eq(gg.Str(`one`, `_`, `two`), `one_two`)
	gtest.Eq(gg.Str(10, `_`, 20), `10_20`)
}

func BenchmarkStr_0(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Str()
	}
}

func BenchmarkStr_1(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Str(`one`)
	}
}

func BenchmarkStr_2(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Str(`one`, `two`)
	}
}

func BenchmarkStr_3(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Str(`one`, `two`, `three`)
	}
}

func TestEllipsis(t *testing.T) {
	defer gtest.Catch(t)

	const src = `🐒🐴🦖🦔🐲🐈`

	gtest.Zero(gg.Ellipsis(src, 0))
	gtest.Eq(gg.Ellipsis(src, 1), `…`)
	gtest.Eq(gg.Ellipsis(src, 2), `🐒…`)
	gtest.Eq(gg.Ellipsis(src, 3), `🐒🐴…`)
	gtest.Eq(gg.Ellipsis(src, 4), `🐒🐴🦖…`)
	gtest.Eq(gg.Ellipsis(src, 5), `🐒🐴🦖🦔…`)
	gtest.Eq(gg.Ellipsis(src, 6), `🐒🐴🦖🦔🐲🐈`)
	gtest.Eq(gg.Ellipsis(src, 7), `🐒🐴🦖🦔🐲🐈`)
	gtest.Eq(gg.Ellipsis(src, 8), `🐒🐴🦖🦔🐲🐈`)
	gtest.Eq(gg.Ellipsis(src, 9), `🐒🐴🦖🦔🐲🐈`)
	gtest.Eq(gg.Ellipsis(src, math.MaxUint), `🐒🐴🦖🦔🐲🐈`)
}

func BenchmarkEllipsis_changed(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Ellipsis(`🐒🐴🦖🦔🐲🐈`, 5)
	}
}

func BenchmarkEllipsis_unchanged(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Ellipsis(`🐒🐴🦖🦔🐲🐈`, 6)
	}
}

func TestTextHeadChar(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src string, valExp rune, sizeExp int) {
		val, size := gg.TextHeadChar(src)
		gtest.Eq(val, valExp, `matching char`)
		gtest.Eq(size, sizeExp, `matching size`)
	}

	test(``, 0, 0)
	test(`one`, 'o', 1)
	test(`🐒🐴🦖🦔🐲🐈`, '🐒', 4)
}

func TestAppendNewlineOpt(t *testing.T) {
	defer gtest.Catch(t)

	testAppendNewlineOptSame(``)
	testAppendNewlineOptSame([]byte(nil))
	testAppendNewlineOptSame([]byte{})
	testAppendNewlineOptSame("\n")
	testAppendNewlineOptSame("one\n")
	testAppendNewlineOptSame("one\r")
	testAppendNewlineOptSame("one\r\n")
	testAppendNewlineOptSame([]byte("\n"))
	testAppendNewlineOptSame([]byte("one\n"))
	testAppendNewlineOptSame([]byte("one\r"))
	testAppendNewlineOptSame([]byte("one\r\n"))

	gtest.TextEq(gg.AppendNewlineOpt(`one`), "one\n")
	gtest.TextEq(gg.AppendNewlineOpt([]byte(`one`)), []byte("one\n"))

	{
		src := []byte("one two")
		tar := src[:len(`one`)]
		out := gg.AppendNewlineOpt(tar)

		gtest.TextEq(src, []byte("one\ntwo"))
		gtest.TextEq(out, []byte("one\n"))
		gtest.Eq(gg.SliceDat(src), gg.SliceDat(tar))
		gtest.Eq(gg.SliceDat(src), gg.SliceDat(out))
	}
}

func testAppendNewlineOptSame[A gg.Text](src A) {
	gtest.Is(gg.AppendNewlineOpt(src), src)
}
