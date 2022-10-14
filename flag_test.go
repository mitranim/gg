package gg_test

import (
	r "reflect"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

type Flags struct {
	Args   []string   `flag:""`
	Str    string     `flag:"-s"`
	Strs   []string   `flag:"-ss"`
	Bool   bool       `flag:"-b"`
	Bools  []bool     `flag:"-bs"`
	Num    float64    `flag:"-n"`
	Nums   []float64  `flag:"-ns"`
	Parser StrsParser `flag:"-p"`
}

type FlagsWithInit struct {
	Args   []string   `flag:""`
	Str    string     `flag:"-s"  init:"one"`
	Strs   []string   `flag:"-ss" init:"two"`
	Bool   bool       `flag:"-b"  init:"true"`
	Bools  []bool     `flag:"-bs" init:"true"`
	Num    float64    `flag:"-n"  init:"12.34"`
	Nums   []float64  `flag:"-ns" init:"56.78"`
	Parser StrsParser `flag:"-p"  init:"three"`
}

type FlagsWithDesc struct {
	Args   []string   `flag:""`
	Str    string     `flag:"-s"  desc:"Str flag"`
	Strs   []string   `flag:"-ss" desc:"Strs flag"`
	Bool   bool       `flag:"-b"  desc:"Bool flag"`
	Bools  []bool     `flag:"-bs" desc:"Bools flag"`
	Num    float64    `flag:"-n"  desc:"Num flag"`
	Nums   []float64  `flag:"-ns" desc:"Nums flag"`
	Parser StrsParser `flag:"-p"  desc:"Parser flag"`
}

type FlagsPart struct {
	Args   []string   `flag:""`
	Str    string     `flag:"-s"  init:"one"   desc:"Str flag"   `
	Strs   []string   `flag:"-ss"              desc:"Strs flag"  `
	Bool   bool       `flag:"-b"  init:"true"  desc:"Bool flag"  `
	Bools  []bool     `flag:"-bs"                                `
	Num    float64    `flag:"-n"  init:"12.34" desc:"Num flag"   `
	Nums   []float64  `flag:"-ns" init:"56.78"                   `
	Parser StrsParser `flag:"-p"               desc:"Parser flag"`
}

type FlagsFull struct {
	Args   []string   `flag:""`
	Str    string     `flag:"-s"  init:"one"   desc:"Str flag"`
	Strs   []string   `flag:"-ss" init:"two"   desc:"Strs flag"`
	Bool   bool       `flag:"-b"  init:"true"  desc:"Bool flag"`
	Bools  []bool     `flag:"-bs" init:"true"  desc:"Bools flag"`
	Num    float64    `flag:"-n"  init:"12.34" desc:"Num flag"`
	Nums   []float64  `flag:"-ns" init:"56.78" desc:"Nums flag"`
	Parser StrsParser `flag:"-p"  init:"three" desc:"Parser flag"`
}

var argsMixed = []string{
	`-s=one`,
	`-ss=two`, `-ss=three`, `-ss`, `four`,
	`-b=false`,
	`-bs=true`, `-bs=false`, `-bs=true`,
	`-n=12`,
	`-ns=23`, `-ns=34`, `-ns`, `45`,
	`-p=five`, `-p=six`, `-p`, `seven`,
	`eight`, `-nine`, `--ten`,
}

var flagsMixed = Flags{
	Str:    `one`,
	Strs:   []string{`two`, `three`, `four`},
	Bool:   false,
	Bools:  []bool{true, false, true},
	Num:    12,
	Nums:   []float64{23, 34, 45},
	Parser: StrsParser{`five`, `six`, `seven`},
	Args:   []string{`eight`, `-nine`, `--ten`},
}

func TestFlagDef(t *testing.T) {
	defer gtest.Catch(t)

	typ := gg.Type[FlagsFull]()
	fields := gg.StructDeepPublicFieldCache.Get(typ)

	gtest.Equal(
		gg.FlagDefCache.Get(typ),
		gg.FlagDef{
			Type:  typ,
			Args:  makeFlagDefField(gg.Head(fields)),
			Flags: gg.Map(gg.Tail(fields), makeFlagDefField),
			Index: map[string]int{
				`-s`:  0,
				`-ss`: 1,
				`-b`:  2,
				`-bs`: 3,
				`-n`:  4,
				`-ns`: 5,
				`-p`:  6,
			},
		},
	)
}

func TestFlagHelp(t *testing.T) {
	defer gtest.Catch(t)

	testFlagHelp[SomeModel](gg.Newline)

	testFlagHelp[Flags](`
flag
----
-s
-ss
-b
-bs
-n
-ns
-p
`)

	testFlagHelp[FlagsWithInit](`
flag    init
-------------
-s      one
-ss     two
-b      true
-bs     true
-n      12.34
-ns     56.78
-p      three
`)

	testFlagHelp[FlagsWithDesc](`
flag    desc
-------------------
-s      Str flag
-ss     Strs flag
-b      Bool flag
-bs     Bools flag
-n      Num flag
-ns     Nums flag
-p      Parser flag
`)

	testFlagHelp[FlagsPart](`
flag    init     desc
----------------------------
-s      one      Str flag
-ss              Strs flag
-b      true     Bool flag
-bs
-n      12.34    Num flag
-ns     56.78
-p               Parser flag
`)

	t.Run(`partial_without_head`, func(t *testing.T) {
		defer gtest.Catch(t)
		defer gg.Swap(&gg.FlagFmtDefault.Head, false).Done()

		testFlagHelp[FlagsPart](`
-s     one      Str flag
-ss             Strs flag
-b     true     Bool flag
-bs
-n     12.34    Num flag
-ns    56.78
-p              Parser flag
`)
	})

	testFlagHelp[FlagsFull](`
flag    init     desc
----------------------------
-s      one      Str flag
-ss     two      Strs flag
-b      true     Bool flag
-bs     true     Bools flag
-n      12.34    Num flag
-ns     56.78    Nums flag
-p      three    Parser flag
`)

	t.Run(`full_without_head_under`, func(t *testing.T) {
		defer gtest.Catch(t)
		defer gg.Swap(&gg.FlagFmtDefault.HeadUnder, ``).Done()

		testFlagHelp[FlagsFull](`
flag    init     desc
-s      one      Str flag
-ss     two      Strs flag
-b      true     Bool flag
-bs     true     Bools flag
-n      12.34    Num flag
-ns     56.78    Nums flag
-p      three    Parser flag
`)
	})

	t.Run(`full_without_head`, func(t *testing.T) {
		defer gtest.Catch(t)
		defer gg.Swap(&gg.FlagFmtDefault.Head, false).Done()

		testFlagHelp[FlagsFull](`
-s     one      Str flag
-ss    two      Strs flag
-b     true     Bool flag
-bs    true     Bools flag
-n     12.34    Num flag
-ns    56.78    Nums flag
-p     three    Parser flag
`)
	})
}

func testFlagHelp[A any](exp string) {
	gtest.Eq(
		trimLines(gg.Newline+gg.FlagHelp[A]()),
		trimLines(exp),
	)
}

func trimLines(src string) string {
	return gg.JoinLines(gg.Map(gg.SplitLines(src), gg.TrimSpaceSuffix[string])...)
}

func BenchmarkFlagHelp(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.FlagHelp[FlagsFull]())
	}
}

func TestFlagParseTo(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`unknown`, func(t *testing.T) {
		defer gtest.Catch(t)

		testFlagUnknown(`-one`)
		testFlagUnknown(`--one`)

		testFlagUnknown(`-one`, ``)
		testFlagUnknown(`--one`, ``)

		testFlagUnknown(`-one`, `two`)
		testFlagUnknown(`--one`, `two`)
	})

	t.Run(`defaults`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(gg.FlagParseTo[FlagsFull](nil), FlagsFull{
			Str:    `one`,
			Strs:   []string{`two`},
			Bool:   true,
			Bools:  []bool{true},
			Num:    12.34,
			Nums:   []float64{56.78},
			Parser: StrsParser{`three`},
		})
	})

	t.Run(`just_args`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src ...string) {
			gtest.Equal(gg.FlagParseTo[Flags](src), Flags{Args: src})
		}

		test()
		test(`one`)
		test(`one`, `two`)
		test(`one`, `two`, `three`)
		test(`one`, `two`, `three`)
		test(`one`, `--two`)
		test(`one`, `--two`, `three`)
		test(`one`, `--two`, `three`, `--four`)
	})

	t.Run(`string_scalar`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src []string, exp string) {
			gtest.Equal(gg.FlagParseTo[Flags](src), Flags{Str: exp})
		}

		test([]string{`-s`, ``}, ``)
		test([]string{`-s`, `one`}, `one`)
		test([]string{`-s`, `one`, `-s`, `two`}, `two`)

		test([]string{`-s=`}, ``)
		test([]string{`-s=one`}, `one`)
		test([]string{`-s=one`, `-s=two`}, `two`)
	})

	t.Run(`string_slice`, func(t *testing.T) {
		defer gtest.Catch(t)

		type Src = [2][]string

		test := func(src Src) {
			gtest.Equal(gg.FlagParseTo[Flags](src[0]), Flags{Strs: src[1]})
		}

		test(Src{{`-ss`, ``}, {``}})
		test(Src{{`-ss`, `one`}, {`one`}})
		test(Src{{`-ss`, `one`, `-ss`, `two`}, {`one`, `two`}})
		test(Src{{`-ss`, `one`, `-ss`, ``, `-ss`, `two`}, {`one`, ``, `two`}})

		test(Src{{`-ss=`}, {``}})
		test(Src{{`-ss=one`}, {`one`}})
		test(Src{{`-ss=one`, `-ss=two`}, {`one`, `two`}})
		test(Src{{`-ss=one`, `-ss=`, `-ss=two`}, {`one`, ``, `two`}})
	})

	var boolInvalid = []string{
		` `, `FALSE`, `TRUE`, `0`, `1`, `123`, `false `, `true `, ` false`, ` true`,
	}

	t.Run(`bool_scalar`, func(t *testing.T) {
		defer gtest.Catch(t)

		t.Run(`invalid`, func(t *testing.T) {
			defer gtest.Catch(t)
			testFlagInvalids(`-b`, boolInvalid)
		})

		test := func(src []string, exp bool) {
			gtest.Equal(gg.FlagParseTo[Flags](src).Bool, exp)
		}

		test([]string{}, false)
		test([]string{`-b`}, true)

		// Bool flags don't support `-flag value`, only `-flag=value`.
		test([]string{`-b`, ``}, true)
		test([]string{`-b`, `arg`}, true)
		test([]string{`-b`, `true`}, true)
		test([]string{`-b`, `false`}, true)

		test([]string{`-b=`}, true)
		test([]string{`-b=false`}, false)
		test([]string{`-b=true`}, true)

		test([]string{`-b`, `-b`}, true)
		test([]string{`-b=false`, `-b`}, true)
		test([]string{`-b=true`, `-b`}, true)

		test([]string{`-b`, `-b=`}, true)
		test([]string{`-b=false`, `-b=false`}, false)
		test([]string{`-b=true`, `-b=true`}, true)
	})

	t.Run(`bool_slice`, func(t *testing.T) {
		defer gtest.Catch(t)

		t.Run(`invalid`, func(t *testing.T) {
			defer gtest.Catch(t)
			testFlagInvalids(`-bs`, boolInvalid)
		})

		test := func(src []string, exp []bool) {
			gtest.Equal(gg.FlagParseTo[Flags](src).Bools, exp)
		}

		test([]string{`-bs`}, []bool{true})
		test([]string{`-bs`, `-bs`}, []bool{true, true})

		test([]string{`-bs=false`, `-bs`}, []bool{false, true})
		test([]string{`-bs=true`, `-bs`}, []bool{true, true})

		test([]string{`-bs=false`, `-bs=true`}, []bool{false, true})
		test([]string{`-bs=true`, `-bs=false`}, []bool{true, false})

		test([]string{`-bs`, `-bs=false`}, []bool{true, false})
		test([]string{`-bs`, `-bs=true`}, []bool{true, true})

		test([]string{`-bs`, ``, `-bs`}, []bool{true})
		test([]string{`-bs`, `arg`, `-bs`}, []bool{true})

		test([]string{`-bs=false`, ``, `-bs`}, []bool{false})
		test([]string{`-bs=false`, `arg`, `-bs`}, []bool{false})
	})

	var numInvalid = []string{` `, `false`, `true`, `±1`, ` 1 `, `1a`, `a1`}

	t.Run(`num_scalar`, func(t *testing.T) {
		defer gtest.Catch(t)

		t.Run(`invalid`, func(t *testing.T) {
			defer gtest.Catch(t)

			testFlagMissingValue(`-n`)
			testFlagMissingValue(`-n`, `-0`)
			testFlagMissingValue(`-n`, `-12.34`)

			testFlagInvalids(`-n`, numInvalid)
		})

		test := func(src []string, exp float64) {
			gtest.Equal(gg.FlagParseTo[Flags](src), Flags{Num: exp})
		}

		test([]string{}, 0)

		test([]string{`-n`, `0`}, 0)
		test([]string{`-n`, `+0`}, 0)

		test([]string{`-n`, `12.34`}, 12.34)
		test([]string{`-n`, `+12.34`}, 12.34)

		test([]string{`-n=-0`}, 0)
		test([]string{`-n=0`}, 0)
		test([]string{`-n=+0`}, 0)

		test([]string{`-n=-12.34`}, -12.34)
		test([]string{`-n=12.34`}, 12.34)
		test([]string{`-n=+12.34`}, 12.34)

		test([]string{`-n`, `34.56`, `-n`, `12.34`}, 12.34)
		test([]string{`-n`, `34.56`, `-n`, `+12.34`}, 12.34)

		test([]string{`-n`, `34.56`, `-n=-12.34`}, -12.34)
		test([]string{`-n`, `34.56`, `-n=12.34`}, 12.34)
		test([]string{`-n`, `34.56`, `-n=+12.34`}, 12.34)
	})

	t.Run(`num_slice`, func(t *testing.T) {
		defer gtest.Catch(t)

		t.Run(`invalid`, func(t *testing.T) {
			defer gtest.Catch(t)
			testFlagMissingValue(`-ns`)
			testFlagInvalids(`-ns`, numInvalid)
		})

		test := func(src []string, exp []float64) {
			gtest.Equal(gg.FlagParseTo[Flags](src), Flags{Nums: exp})
		}

		test([]string{`-ns`, `0`}, []float64{0})
		test([]string{`-ns`, `+0`}, []float64{0})

		test([]string{`-ns`, `12.34`}, []float64{12.34})
		test([]string{`-ns`, `+12.34`}, []float64{12.34})

		test([]string{`-ns=12.34`}, []float64{12.34})
		test([]string{`-ns=+12.34`}, []float64{12.34})

		test([]string{`-ns`, `12.34`, `-ns`, `56.78`}, []float64{12.34, 56.78})
		test([]string{`-ns=12.34`, `-ns=56.78`}, []float64{12.34, 56.78})
	})

	t.Run(`Parser`, func(t *testing.T) {
		defer gtest.Catch(t)

		type Src []string
		type Out = StrsParser

		type Tar struct {
			Val Out `flag:"-v"`
		}

		test := func(src Src, out Out) {
			gtest.Equal(gg.FlagParseTo[Tar](src), Tar{out})
		}

		test(nil, nil)

		test(Src{`-v`, ``}, Out{``})
		test(Src{`-v=`}, Out{``})

		test(Src{`-v`, `10`}, Out{`10`})
		test(Src{`-v=10`}, Out{`10`})

		test(Src{`-v`, `10`, `-v`, `20`}, Out{`10`, `20`})
		test(Src{`-v=10`, `-v=20`}, Out{`10`, `20`})
	})

	t.Run(`flag.Value`, func(t *testing.T) {
		defer gtest.Catch(t)

		type Src []string
		type Out = IntsValue

		type Tar struct {
			Val Out `flag:"-v"`
		}

		test := func(src Src, out Out) {
			gtest.Equal(gg.FlagParseTo[Tar](src), Tar{out})
		}

		test(nil, nil)

		test(Src{`-v`, `10`}, Out{10})
		test(Src{`-v=10`}, Out{10})

		test(Src{`-v`, `10`, `-v`, `20`}, Out{10, 20})
		test(Src{`-v=10`, `-v=20`}, Out{10, 20})
	})

	t.Run(`[]flag.Value`, func(t *testing.T) {
		defer gtest.Catch(t)

		type Src []string
		type Out = []IntValue

		type Tar struct {
			Val Out `flag:"-v"`
		}

		test := func(src Src, out Out) {
			gtest.Equal(gg.FlagParseTo[Tar](src), Tar{out})
		}

		test(nil, nil)

		test(Src{`-v`, `10`}, Out{{10}})
		test(Src{`-v=10`}, Out{{10}})

		test(Src{`-v`, `10`, `-v`, `20`}, Out{{10}, {20}})
		test(Src{`-v=10`, `-v=20`}, Out{{10}, {20}})
	})

	t.Run(`mixed`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(gg.FlagParseTo[Flags](argsMixed), flagsMixed)
	})
}

func BenchmarkFlagParseTo_empty(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.FlagParseTo[Flags](nil))
	}
}

// Kinda slow (microseconds) but not anyone's bottleneck.
func BenchmarkFlagParseTo_full(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.FlagParseTo[Flags](argsMixed))
	}
}

func testFlagInvalid(src ...string) {
	gtest.PanicStr(`unable to decode`, func() {
		gg.FlagParseTo[Flags](src)
	})
}

func testFlagInvalids(flag string, src []string) {
	for _, val := range src {
		testFlagInvalid(flag + `=` + val)
	}
}

func testFlagUnknown(src ...string) {
	gtest.PanicStr(`unable to find flag`, func() {
		gg.FlagParseTo[Flags](src)
	})
}

func testFlagMissingValue(src ...string) {
	gtest.PanicStr(`missing value for trailing flag`, func() {
		gg.FlagParseTo[Flags](src)
	})
}

func makeFlagDefField(src r.StructField) (out gg.FlagDefField) {
	out.Set(src)
	return
}
