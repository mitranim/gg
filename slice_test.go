package gg_test

import (
	"fmt"
	"strconv"
	"testing"
	u "unsafe"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func ExampleSlice() {
	values := []string{`one`, `two`, `three`}
	indexes := []int{0, 2}
	result := gg.Map(indexes, gg.SliceOf(values...).Get)

	fmt.Println(gg.GoString(result))

	// Output:
	// []string{"one", "three"}
}

func TestLens(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Lens[[]int](), 0)
	gtest.Eq(gg.Lens([]int{}, []int{10}), 1)
	gtest.Eq(gg.Lens([]int{}, []int{10}, []int{20, 30}), 3)
}

func BenchmarkLens(b *testing.B) {
	val := [][]int{{}, {10}, {20, 30}, {40, 50, 60}}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Lens(val...))
	}
}

func TestGrowLen(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`from_empty`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(
			gg.GrowLen([]int(nil), 3),
			[]int{0, 0, 0},
		)
	})

	t.Run(`within_capacity`, func(t *testing.T) {
		defer gtest.Catch(t)

		src := []int{10, 20, 30, 40, 0, 0, 0, 0}
		cur := src[:2]

		test := func(size int, expTar, expSrc []int) {
			tar := gg.GrowLen(cur, size)
			gtest.Equal(src, expSrc)
			gtest.Equal(tar, expTar)
			gtest.Eq(cap(tar), cap(src))
			gtest.Eq(u.SliceData(src), u.SliceData(tar))
		}

		test(0, []int{10, 20}, []int{10, 20, 30, 40, 0, 0, 0, 0})
		test(1, []int{10, 20, 0}, []int{10, 20, 0, 40, 0, 0, 0, 0})
		test(2, []int{10, 20, 0, 0}, []int{10, 20, 0, 0, 0, 0, 0, 0})
		test(3, []int{10, 20, 0, 0, 0}, []int{10, 20, 0, 0, 0, 0, 0, 0})
	})

	t.Run(`exceeding_capacity`, func(t *testing.T) {
		defer gtest.Catch(t)

		src := []int{10, 20, 30, 40}[:2]
		tar := gg.GrowLen(src, 3)

		gtest.Equal(src, []int{10, 20})
		gtest.Equal(tar, []int{10, 20, 0, 0, 0})
		gtest.Equal(src[:4], []int{10, 20, 30, 40})
		gtest.NotEq(u.SliceData(src), u.SliceData(tar))
	})
}

func BenchmarkGrowLen(b *testing.B) {
	buf := make([]byte, 0, 1024)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.GrowLen(buf, 128))
	}
}

func TestGrowCap(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`from_empty`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.True(gg.GrowCap([]int(nil), 0) == nil)

		{
			tar := gg.GrowCap([]int(nil), 8)
			gtest.Eq(len(tar), 0)
			gtest.Eq(cap(tar), 8)
		}
	})

	t.Run(`within_capacity`, func(t *testing.T) {
		defer gtest.Catch(t)

		src := []int{10, 20, 30, 40}[:2]
		tar := gg.GrowCap(src, 2)

		gtest.Equal(tar, src)
		gtest.Equal(tar, []int{10, 20})

		gtest.Eq(u.SliceData(src), u.SliceData(tar))
		gtest.Eq(len(tar), len(src))
		gtest.Eq(cap(tar), cap(src))
	})
}

func TestTruncLen(t *testing.T) {
	defer gtest.Catch(t)

	type Slice = []int

	test := func(src Slice, size int, exp Slice) {
		out := gg.TruncLen(src, size)
		gtest.Equal(out, exp)

		gtest.Eq(
			u.SliceData(out),
			u.SliceData(src),
			`reslicing must preserve data pointer`,
		)

		gtest.Eq(
			cap(out),
			cap(src),
			`reslicing must preserve capacity`,
		)
	}

	test(nil, -1, nil)
	test(nil, 0, nil)
	test(nil, 1, nil)
	test(nil, 2, nil)

	test(Slice{}, -1, Slice{})
	test(Slice{}, 0, Slice{})
	test(Slice{}, 1, Slice{})
	test(Slice{}, 2, Slice{})

	src := Slice{10, 20, 30, 40}
	test(src, -1, Slice{})
	test(src, 0, Slice{})
	test(src, 1, Slice{10})
	test(src, 2, Slice{10, 20})
	test(src, 3, Slice{10, 20, 30})
	test(src, 4, src)
	test(src, 5, src)
	test(src, 6, src)
}

func TestTake(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.Take([]int(nil), -2), []int(nil))
	gtest.Equal(gg.Take([]int(nil), -1), []int(nil))
	gtest.Equal(gg.Take([]int(nil), 0), []int(nil))
	gtest.Equal(gg.Take([]int(nil), 1), []int(nil))
	gtest.Equal(gg.Take([]int(nil), 2), []int(nil))

	gtest.Equal(gg.Take([]int{}, -2), []int{})
	gtest.Equal(gg.Take([]int{}, -1), []int{})
	gtest.Equal(gg.Take([]int{}, 0), []int{})
	gtest.Equal(gg.Take([]int{}, 1), []int{})
	gtest.Equal(gg.Take([]int{}, 2), []int{})

	gtest.Equal(gg.Take([]int{10}, -2), []int{})
	gtest.Equal(gg.Take([]int{10}, -1), []int{})
	gtest.Equal(gg.Take([]int{10}, 0), []int{})
	gtest.Equal(gg.Take([]int{10}, 1), []int{10})
	gtest.Equal(gg.Take([]int{10}, 2), []int{10})

	gtest.Equal(gg.Take([]int{10, 20}, -2), []int{})
	gtest.Equal(gg.Take([]int{10, 20}, -1), []int{})
	gtest.Equal(gg.Take([]int{10, 20}, 0), []int{})
	gtest.Equal(gg.Take([]int{10, 20}, 1), []int{10})
	gtest.Equal(gg.Take([]int{10, 20}, 2), []int{10, 20})

	gtest.Equal(gg.Take([]int{10, 20, 30}, -2), []int{})
	gtest.Equal(gg.Take([]int{10, 20, 30}, -1), []int{})
	gtest.Equal(gg.Take([]int{10, 20, 30}, 0), []int{})
	gtest.Equal(gg.Take([]int{10, 20, 30}, 1), []int{10})
	gtest.Equal(gg.Take([]int{10, 20, 30}, 2), []int{10, 20})

	gtest.Equal(gg.Take([]int{10, 20, 30, 40}, -2), []int{})
	gtest.Equal(gg.Take([]int{10, 20, 30, 40}, -1), []int{})
	gtest.Equal(gg.Take([]int{10, 20, 30, 40}, 0), []int{})
	gtest.Equal(gg.Take([]int{10, 20, 30, 40}, 1), []int{10})
	gtest.Equal(gg.Take([]int{10, 20, 30, 40}, 2), []int{10, 20})
}

func TestDrop(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.Drop([]int(nil), -2), []int(nil))
	gtest.Equal(gg.Drop([]int(nil), -1), []int(nil))
	gtest.Equal(gg.Drop([]int(nil), 0), []int(nil))
	gtest.Equal(gg.Drop([]int(nil), 1), []int(nil))
	gtest.Equal(gg.Drop([]int(nil), 2), []int(nil))

	gtest.Equal(gg.Drop([]int{}, -2), []int{})
	gtest.Equal(gg.Drop([]int{}, -1), []int{})
	gtest.Equal(gg.Drop([]int{}, 0), []int{})
	gtest.Equal(gg.Drop([]int{}, 1), []int{})
	gtest.Equal(gg.Drop([]int{}, 2), []int{})

	gtest.Equal(gg.Drop([]int{10}, -2), []int{10})
	gtest.Equal(gg.Drop([]int{10}, -1), []int{10})
	gtest.Equal(gg.Drop([]int{10}, 0), []int{10})
	gtest.Equal(gg.Drop([]int{10}, 1), []int{})
	gtest.Equal(gg.Drop([]int{10}, 2), []int{})

	gtest.Equal(gg.Drop([]int{10, 20}, -2), []int{10, 20})
	gtest.Equal(gg.Drop([]int{10, 20}, -1), []int{10, 20})
	gtest.Equal(gg.Drop([]int{10, 20}, 0), []int{10, 20})
	gtest.Equal(gg.Drop([]int{10, 20}, 1), []int{20})
	gtest.Equal(gg.Drop([]int{10, 20}, 2), []int{})

	gtest.Equal(gg.Drop([]int{10, 20, 30}, -2), []int{10, 20, 30})
	gtest.Equal(gg.Drop([]int{10, 20, 30}, -1), []int{10, 20, 30})
	gtest.Equal(gg.Drop([]int{10, 20, 30}, 0), []int{10, 20, 30})
	gtest.Equal(gg.Drop([]int{10, 20, 30}, 1), []int{20, 30})
	gtest.Equal(gg.Drop([]int{10, 20, 30}, 2), []int{30})

	gtest.Equal(gg.Drop([]int{10, 20, 30, 40}, -2), []int{10, 20, 30, 40})
	gtest.Equal(gg.Drop([]int{10, 20, 30, 40}, -1), []int{10, 20, 30, 40})
	gtest.Equal(gg.Drop([]int{10, 20, 30, 40}, 0), []int{10, 20, 30, 40})
	gtest.Equal(gg.Drop([]int{10, 20, 30, 40}, 1), []int{20, 30, 40})
	gtest.Equal(gg.Drop([]int{10, 20, 30, 40}, 2), []int{30, 40})
}

func TestTakeWhile(t *testing.T) {
	defer gtest.Catch(t)

	fun := gg.IsPos[int]

	test := func(src, exp []int) {
		gtest.Equal(gg.TakeWhile(src, fun), exp)
		gtest.Equal(gg.TakeWhile(src, nil), src[:0])
	}

	same := func(src ...int) { test(src, src) }

	test(nil, nil)

	test([]int{0}, []int{})
	test([]int{-20, -10, 0}, []int{})
	test([]int{0, -10, -20}, []int{})

	test([]int{20, 10, 0}, []int{20, 10})
	test([]int{0, 10, 20}, []int{})

	test([]int{-20, -10, 0, 10, 20}, []int{})
	test([]int{20, 10, 0, -10, -20}, []int{20, 10})

	test([]int{20, 10, 0, 30, 40}, []int{20, 10})

	same(20, 10, 30, 40)
}

func TestTakeLastWhile(t *testing.T) {
	defer gtest.Catch(t)

	fun := gg.IsPos[int]

	test := func(src, exp []int) {
		gtest.Equal(gg.TakeLastWhile(src, fun), exp)
		gtest.Empty(gg.TakeLastWhile(src, nil))
		gtest.SliceIs(gg.TakeLastWhile(src, nil), src[len(src):])
	}

	same := func(src ...int) { test(src, src) }

	test(nil, nil)

	test([]int{0}, []int{})
	test([]int{-20, -10, 0}, []int{})
	test([]int{0, -10, -20}, []int{})

	test([]int{20, 10, 0}, []int{})
	test([]int{0, 10, 20}, []int{10, 20})

	test([]int{-20, -10, 0, 10, 20}, []int{10, 20})
	test([]int{20, 10, 0, -10, -20}, []int{})

	test([]int{20, 10, 0, 30, 40}, []int{30, 40})

	same(20, 10, 30, 40)
}

func TestDropWhile(t *testing.T) {
	defer gtest.Catch(t)

	fun := gg.IsPos[int]

	test := func(src, exp []int) {
		gtest.Equal(gg.DropWhile(src, fun), exp)
		gtest.Equal(gg.DropWhile(src, nil), src)
	}

	test(nil, nil)

	test([]int{10}, []int{})

	test([]int{10, 20}, []int{})
	test([]int{20, 10}, []int{})

	test([]int{20, 10, 0}, []int{0})
	test([]int{0, 10, 20}, []int{0, 10, 20})

	test([]int{20, 10, -10, -20}, []int{-10, -20})
	test([]int{-20, -10, 10, 20}, []int{-20, -10, 10, 20})

	test([]int{20, 10, 0, -10, -20}, []int{0, -10, -20})
	test([]int{-20, -10, 0, 10, 20}, []int{-20, -10, 0, 10, 20})
}

func TestDropLastWhile(t *testing.T) {
	defer gtest.Catch(t)

	fun := gg.IsPos[int]

	test := func(src, exp []int) {
		gtest.Equal(gg.DropLastWhile(src, fun), exp)
		gtest.SliceIs(gg.DropLastWhile(src, nil), src)
	}

	test(nil, nil)

	test([]int{10}, []int{})

	test([]int{10, 20}, []int{})
	test([]int{20, 10}, []int{})

	test([]int{20, 10, 0}, []int{20, 10, 0})
	test([]int{0, 10, 20}, []int{0})

	test([]int{20, 10, -10, -20}, []int{20, 10, -10, -20})
	test([]int{-20, -10, 10, 20}, []int{-20, -10})

	test([]int{20, 10, 0, -10, -20}, []int{20, 10, 0, -10, -20})
	test([]int{-20, -10, 0, 10, 20}, []int{-20, -10, 0})
}

func TestMap(t *testing.T) {
	defer gtest.Catch(t)

	fun := gg.Map[int, string]
	elemFun := strconv.Itoa

	testMapCommon(fun)

	gtest.Equal(fun([]int{10, 10}, elemFun), []string{`10`, `10`})
	gtest.Equal(fun([]int{10, 20, 10}, elemFun), []string{`10`, `20`, `10`})
	gtest.Equal(fun([]int{10, 20, 10, 20}, elemFun), []string{`10`, `20`, `10`, `20`})
	gtest.Equal(fun([]int{10, 20, 20, 10}, elemFun), []string{`10`, `20`, `20`, `10`})
}

func testMapCommon(mapFun func([]int, func(int) string) []string) {
	elemFun := strconv.Itoa

	gtest.Equal(mapFun([]int(nil), elemFun), []string(nil))
	gtest.Equal(mapFun([]int{}, elemFun), []string{})
	gtest.Equal(mapFun([]int{10}, elemFun), []string{`10`})
	gtest.Equal(mapFun([]int{10, 20}, elemFun), []string{`10`, `20`})
	gtest.Equal(mapFun([]int{10, 20, 30}, elemFun), []string{`10`, `20`, `30`})
}

func BenchmarkMap(b *testing.B) {
	val := gg.Span(32)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Map(val, gg.Inc[int]))
	}
}

func TestMapMut(t *testing.T) {
	defer gtest.Catch(t)

	src := []int{10, 20, 30}
	gtest.SliceIs(gg.MapMut(src, nil), src)
	gtest.Equal(src, []int{10, 20, 30})

	gtest.SliceIs(gg.MapMut(src, gg.Inc[int]), src)
	gtest.Equal(src, []int{11, 21, 31})

	gtest.SliceIs(gg.MapMut(src, gg.Dec[int]), src)
	gtest.Equal(src, []int{10, 20, 30})

	gtest.SliceIs(gg.MapMut(src, gg.Dec[int]), src)
	gtest.Equal(src, []int{9, 19, 29})
}

func TestMapPtr(t *testing.T) {
	defer gtest.Catch(t)

	src := []int{10, 20, 30}
	gtest.Equal(gg.MapPtr(src, DoublePtrStr), []string{`20`, `40`, `60`})
	gtest.Equal(src, []int{20, 40, 60})
}

func DoublePtrStr(ptr *int) string {
	*ptr *= 2
	return gg.String(*ptr)
}

func TestMap2(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(
		gg.Map2([]int(nil), []int(nil), gg.Plus2[int]),
		[]int(nil),
	)

	gtest.Equal(
		gg.Map2([]int{}, []int(nil), gg.Plus2[int]),
		[]int(nil),
	)

	gtest.Equal(
		gg.Map2([]int(nil), []int{}, gg.Plus2[int]),
		[]int(nil),
	)

	gtest.Equal(
		gg.Map2([]int{}, []int{}, gg.Plus2[int]),
		[]int{},
	)

	gtest.PanicStr(`length mismatch`, func() {
		gg.Map2([]int{}, []int{10}, gg.Plus2[int])
	})

	gtest.PanicStr(`length mismatch`, func() {
		gg.Map2([]int{10}, []int{}, gg.Plus2[int])
	})

	gtest.Equal(
		gg.Map2([]int{10, 20, 30}, []int{40, 50, 60}, gg.Plus2[int]),
		[]int{50, 70, 90},
	)
}

func TestMapFlat(t *testing.T) {
	defer gtest.Catch(t)

	fun := gg.MapFlat[[]string, int, string]
	elemFun := itoaPair

	testMapFlatCommon(fun)

	gtest.Equal(fun([]int{10}, elemFun), []string{`10`, `10`})
	gtest.Equal(fun([]int{10, 10}, elemFun), []string{`10`, `10`, `10`, `10`})
	gtest.Equal(fun([]int{10, 20}, elemFun), []string{`10`, `10`, `20`, `20`})
}

func itoaPair(src int) []string {
	return []string{strconv.Itoa(src), strconv.Itoa(src)}
}

func testMapFlatCommon(fun func([]int, func(int) []string) []string) {
	gtest.Equal(
		fun([]int(nil), intStrPair),
		[]string(nil),
	)

	gtest.Equal(
		fun([]int{}, intStrPair),
		[]string(nil),
	)

	gtest.Equal(
		fun([]int{10}, intStrPair),
		[]string{`9`, `11`},
	)

	gtest.Equal(
		fun([]int{10, 20}, intStrPair),
		[]string{`9`, `11`, `19`, `21`},
	)

	gtest.Equal(
		fun([]int{10, 20, 30}, intStrPair),
		[]string{`9`, `11`, `19`, `21`, `29`, `31`},
	)
}

func BenchmarkMapFlat(b *testing.B) {
	val := gg.Span(32)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.MapFlat(val, intPair))
	}
}

func TestMapUniq(t *testing.T) {
	defer gtest.Catch(t)

	fun := gg.MapUniq[int, string]
	elemFun := strconv.Itoa

	testMapCommon(fun)

	gtest.Equal(fun([]int{10, 10}, elemFun), []string{`10`})
	gtest.Equal(fun([]int{10, 20, 10}, elemFun), []string{`10`, `20`})
	gtest.Equal(fun([]int{10, 20, 10, 20}, elemFun), []string{`10`, `20`})
	gtest.Equal(fun([]int{10, 20, 20, 10}, elemFun), []string{`10`, `20`})
}

func TestMapFlatUniq(t *testing.T) {
	defer gtest.Catch(t)

	fun := gg.MapFlatUniq[[]string, int, string]
	elemFun := itoaPair

	testMapFlatCommon(fun)

	gtest.Equal(fun([]int{10}, elemFun), []string{`10`})
	gtest.Equal(fun([]int{10, 10}, elemFun), []string{`10`})
	gtest.Equal(fun([]int{10, 20}, elemFun), []string{`10`, `20`})
	gtest.Equal(fun([]int{10, 20, 20, 10}, elemFun), []string{`10`, `20`})
}

func TestIndex(t *testing.T) {
	defer gtest.Catch(t)

	type Slice = []int
	type Map = map[int]int

	gtest.Zero(gg.Index[Slice, int, int](Slice(nil), nil))
	gtest.Zero(gg.Index[Slice, int, int](Slice{10, 20}, nil))
	gtest.Equal(gg.Index(Slice(nil), gg.Inc[int]), Map{})

	gtest.Equal(
		gg.Index(Slice{10, 20}, gg.Inc[int]),
		Map{11: 10, 21: 20},
	)
}

func TestIndexInto(t *testing.T) {
	defer gtest.Catch(t)

	type Map = map[int]int
	tar := Map{}

	gg.IndexInto(tar, nil, nil)
	gtest.Equal(tar, Map{})

	gg.IndexInto(tar, []int{10, 20}, nil)
	gtest.Equal(tar, Map{})

	gg.IndexInto(tar, []int{10, 20}, gg.Inc[int])
	gtest.Equal(tar, Map{11: 10, 21: 20})
}

func TestIndexPair(t *testing.T) {
	defer gtest.Catch(t)

	type Slice = []int
	type Map = map[int]int

	gtest.Zero(gg.IndexPair[Slice, int, int, int](Slice(nil), nil))
	gtest.Zero(gg.IndexPair[Slice, int, int, int](Slice{10, 20}, nil))
	gtest.Zero(gg.IndexPair(Slice(nil), ToPair[int]))

	gtest.Equal(
		gg.IndexPair(Slice{10, 20}, ToPair[int]),
		Map{9: 11, 19: 21},
	)
}

func TestIndexPairInto(t *testing.T) {
	defer gtest.Catch(t)

	type Map = map[int]int
	tar := Map{}

	gg.IndexPairInto[int](tar, nil, nil)
	gtest.Equal(tar, Map{})

	gg.IndexPairInto(tar, []int{10, 20}, nil)
	gtest.Equal(tar, Map{})

	gg.IndexPairInto(tar, []int{10, 20}, ToPair[int])
	gtest.Equal(tar, Map{9: 11, 19: 21})
}

func TestTimes(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.Times[string](-1, nil))
	gtest.Zero(gg.Times[string](0, nil))
	gtest.Zero(gg.Times[string](1, nil))
	gtest.Zero(gg.Times(-1, gg.String[int]))
	gtest.Zero(gg.Times(0, gg.String[int]))

	gtest.Equal(gg.Times(1, gg.String[int]), []string{`0`})
	gtest.Equal(gg.Times(2, gg.String[int]), []string{`0`, `1`})
	gtest.Equal(gg.Times(3, gg.String[int]), []string{`0`, `1`, `2`})
}

func TestTimesAppend(t *testing.T) {
	defer gtest.Catch(t)

	gtest.NotPanic(func() {
		gg.TimesAppend((*[]string)(nil), -1, nil)
		gg.TimesAppend((*[]string)(nil), 0, nil)
		gg.TimesAppend((*[]string)(nil), 1, nil)

		gg.TimesAppend((*[]string)(nil), -1, gg.String[int])
		gg.TimesAppend((*[]string)(nil), 0, gg.String[int])
		gg.TimesAppend((*[]string)(nil), 1, gg.String[int])
	})

	var tar []string

	gg.TimesAppend(&tar, -1, gg.String[int])
	gtest.Zero(tar)

	gg.TimesAppend(&tar, 0, gg.String[int])
	gtest.Zero(tar)

	gg.TimesAppend(&tar, 1, gg.String[int])
	gtest.Equal(tar, []string{`0`})

	gg.TimesAppend(&tar, 2, gg.String[int])
	gtest.Equal(tar, []string{`0`, `0`, `1`})

	gg.TimesAppend(&tar, 3, gg.String[int])
	gtest.Equal(tar, []string{`0`, `0`, `1`, `0`, `1`, `2`})
}

func TestCount(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.Count([]int(nil), nil))
	gtest.Zero(gg.Count([]int{}, nil))
	gtest.Zero(gg.Count([]int{10}, nil))
	gtest.Zero(gg.Count([]int{10, 20}, nil))

	gtest.Zero(gg.Count([]int(nil), False1[int]))
	gtest.Zero(gg.Count([]int{}, False1[int]))
	gtest.Zero(gg.Count([]int{10}, False1[int]))
	gtest.Zero(gg.Count([]int{10, 20}, False1[int]))

	gtest.Zero(gg.Count([]int(nil), True1[int]))
	gtest.Zero(gg.Count([]int{}, True1[int]))

	gtest.Eq(gg.Count([]int{10}, True1[int]), 1)
	gtest.Eq(gg.Count([]int{10, 20}, True1[int]), 2)
	gtest.Eq(gg.Count([]int{10, 20, 30}, True1[int]), 3)

	gtest.Eq(gg.Count([]int{-10, 10, -20, 20, -30}, gg.IsNeg[int]), 3)
	gtest.Eq(gg.Count([]int{-10, 10, -20, 20, -30}, gg.IsPos[int]), 2)
}

func TestFold(t *testing.T) {
	defer gtest.Catch(t)

	const acc = 10
	gtest.Eq(gg.Fold([]int(nil), acc, nil), acc)
	gtest.Eq(gg.Fold([]int{}, acc, nil), acc)
	gtest.Eq(gg.Fold([]int{20}, acc, nil), acc)
	gtest.Eq(gg.Fold([]int{20, 30}, acc, nil), acc)

	gtest.Eq(gg.Fold([]int{20}, 10, gg.SubUncheck[int]), 10-20)
	gtest.Eq(gg.Fold([]int{20}, 10, gg.Plus2[int]), 10+20)

	gtest.Eq(gg.Fold([]int{20, 30}, 10, gg.SubUncheck[int]), 10-20-30)
	gtest.Eq(gg.Fold([]int{20, 30}, 10, gg.Plus2[int]), 10+20+30)

	gtest.Eq(gg.Fold([]int{20, 30, 40}, 10, gg.SubUncheck[int]), 10-20-30-40)
	gtest.Eq(gg.Fold([]int{20, 30, 40}, 10, gg.Plus2[int]), 10+20+30+40)
}

func TestFoldz(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.Foldz[int]([]int(nil), nil))
	gtest.Zero(gg.Foldz[int]([]int{}, nil))
	gtest.Zero(gg.Foldz[int]([]int{10}, nil))
	gtest.Zero(gg.Foldz[int]([]int{10, 20}, nil))

	gtest.Eq(gg.Foldz([]int{10}, gg.SubUncheck[int]), 0-10)
	gtest.Eq(gg.Foldz([]int{10}, gg.Plus2[int]), 0+10)

	gtest.Eq(gg.Foldz([]int{10, 20}, gg.SubUncheck[int]), 0-10-20)
	gtest.Eq(gg.Foldz([]int{10, 20}, gg.Plus2[int]), 0+10+20)

	gtest.Eq(gg.Foldz([]int{10, 20, 30}, gg.SubUncheck[int]), 0-10-20-30)
	gtest.Eq(gg.Foldz([]int{10, 20, 30}, gg.Plus2[int]), 0+10+20+30)
}

func TestFold1(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.Fold1[int]([]int(nil), nil))
	gtest.Zero(gg.Fold1[int]([]int{}, nil))

	gtest.Eq(gg.Fold1([]int{10}, nil), 10)
	gtest.Eq(gg.Fold1([]int{10}, gg.SubUncheck[int]), 10)
	gtest.Eq(gg.Fold1([]int{10}, gg.Plus2[int]), 10)

	gtest.Eq(gg.Fold1([]int{10, 20}, nil), 10)
	gtest.Eq(gg.Fold1([]int{10, 20}, gg.SubUncheck[int]), 10-20)
	gtest.Eq(gg.Fold1([]int{10, 20}, gg.Plus2[int]), 10+20)

	gtest.Eq(gg.Fold1([]int{10, 20, 30}, nil), 10)
	gtest.Eq(gg.Fold1([]int{10, 20, 30}, gg.SubUncheck[int]), 10-20-30)
	gtest.Eq(gg.Fold1([]int{10, 20, 30}, gg.Plus2[int]), 10+20+30)
}

func TestFilter(t *testing.T) {
	defer gtest.Catch(t)

	type Src = []int

	gtest.Zero(gg.Filter(Src(nil), nil))
	gtest.Zero(gg.Filter(Src{}, nil))
	gtest.Zero(gg.Filter(Src{10}, nil))
	gtest.Zero(gg.Filter(Src{10, 20}, nil))
	gtest.Zero(gg.Filter(Src{10, 20, 30}, nil))
	gtest.Zero(gg.Filter(Src{10, 20, 30}, False1[int]))

	gtest.Equal(gg.Filter(Src{10}, True1[int]), Src{10})
	gtest.Equal(gg.Filter(Src{10, 20}, True1[int]), Src{10, 20})
	gtest.Equal(gg.Filter(Src{10, 20, 30}, True1[int]), Src{10, 20, 30})

	gtest.Equal(
		gg.Filter(Src{-10, 10, -20, 20, -30}, gg.IsNeg[int]),
		Src{-10, -20, -30},
	)

	gtest.Equal(
		gg.Filter(Src{-10, 10, -20, 20, -30}, gg.IsPos[int]),
		Src{10, 20},
	)
}

func TestFilterAppend(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.NotPanic(func() {
		gg.FilterAppend((*Type)(nil), nil, nil)
		gg.FilterAppend((*Type)(nil), nil, True1[int])
		gg.FilterAppend((*Type)(nil), Type{}, True1[int])
		gg.FilterAppend((*Type)(nil), Type{10}, True1[int])
		gg.FilterAppend((*Type)(nil), Type{10, 20}, True1[int])
	})

	var tar Type

	gg.FilterAppend(&tar, nil, True1[int])
	gtest.Zero(tar)

	gg.FilterAppend(&tar, Type{}, True1[int])
	gtest.Zero(tar)

	gg.FilterAppend(&tar, Type{10}, True1[int])
	gtest.Equal(tar, Type{10})

	gg.FilterAppend(&tar, Type{20, 30}, True1[int])
	gtest.Equal(tar, Type{10, 20, 30})

	gg.FilterAppend(&tar, Type{40, 50}, False1[int])
	gtest.Equal(tar, Type{10, 20, 30})

	gg.FilterAppend(&tar, Type{-10, 10, -20, 20, -30}, gg.IsNeg[int])
	gtest.Equal(tar, Type{10, 20, 30, -10, -20, -30})
}

func TestReject(t *testing.T) {
	defer gtest.Catch(t)

	type Src = []int

	gtest.Zero(gg.Reject(Src(nil), nil))
	gtest.Zero(gg.Reject(Src{}, nil))
	gtest.Zero(gg.Reject(Src{10}, nil))
	gtest.Zero(gg.Reject(Src{10, 20}, nil))
	gtest.Zero(gg.Reject(Src{10, 20, 30}, nil))
	gtest.Zero(gg.Reject(Src{10, 20, 30}, True1[int]))

	gtest.Equal(gg.Reject(Src{10}, False1[int]), Src{10})
	gtest.Equal(gg.Reject(Src{10, 20}, False1[int]), Src{10, 20})
	gtest.Equal(gg.Reject(Src{10, 20, 30}, False1[int]), Src{10, 20, 30})

	gtest.Equal(
		gg.Reject(Src{-10, 10, -20, 20, -30}, gg.IsNeg[int]),
		Src{10, 20},
	)

	gtest.Equal(
		gg.Reject(Src{-10, 10, -20, 20, -30}, gg.IsPos[int]),
		Src{-10, -20, -30},
	)
}

func TestRejectAppend(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.NotPanic(func() {
		gg.RejectAppend((*Type)(nil), nil, nil)
		gg.RejectAppend((*Type)(nil), nil, False1[int])
		gg.RejectAppend((*Type)(nil), Type{}, False1[int])
		gg.RejectAppend((*Type)(nil), Type{10}, False1[int])
		gg.RejectAppend((*Type)(nil), Type{10, 20}, False1[int])
	})

	var tar Type

	gg.RejectAppend(&tar, nil, False1[int])
	gtest.Zero(tar)

	gg.RejectAppend(&tar, Type{}, False1[int])
	gtest.Zero(tar)

	gg.RejectAppend(&tar, Type{10}, False1[int])
	gtest.Equal(tar, Type{10})

	gg.RejectAppend(&tar, Type{20, 30}, False1[int])
	gtest.Equal(tar, Type{10, 20, 30})

	gg.RejectAppend(&tar, Type{40, 50}, True1[int])
	gtest.Equal(tar, Type{10, 20, 30})

	gg.RejectAppend(&tar, Type{-10, 10, -20, 20, -30}, gg.IsPos[int])
	gtest.Equal(tar, Type{10, 20, 30, -10, -20, -30})
}

func TestFilterIndex(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.FilterIndex(Type(nil), nil))
	gtest.Zero(gg.FilterIndex(Type{}, nil))
	gtest.Zero(gg.FilterIndex(Type{10, 20, 30}, nil))
	gtest.Zero(gg.FilterIndex(Type{10, 20, 30}, False1[int]))

	gtest.Equal(gg.FilterIndex(Type{10}, True1[int]), []int{0})
	gtest.Equal(gg.FilterIndex(Type{10, 20}, True1[int]), []int{0, 1})
	gtest.Equal(gg.FilterIndex(Type{10, 20, 30}, True1[int]), []int{0, 1, 2})

	gtest.Equal(gg.FilterIndex(Type{-10, 10, -20, 20, -30}, gg.IsNeg[int]), []int{0, 2, 4})
	gtest.Equal(gg.FilterIndex(Type{-10, 10, -20, 20, -30}, gg.IsPos[int]), []int{1, 3})
}

func TestZeroIndex(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.ZeroIndex(Type(nil)))
	gtest.Zero(gg.ZeroIndex(Type{}))
	gtest.Zero(gg.ZeroIndex(Type{10}))
	gtest.Zero(gg.ZeroIndex(Type{10, 20}))
	gtest.Zero(gg.ZeroIndex(Type{10, 20, 30}))

	gtest.Equal(gg.ZeroIndex(Type{0}), []int{0})
	gtest.Equal(gg.ZeroIndex(Type{0, 0}), []int{0, 1})
	gtest.Equal(gg.ZeroIndex(Type{0, 0, 0}), []int{0, 1, 2})

	gtest.Equal(gg.ZeroIndex(Type{0, 10}), []int{0})
	gtest.Equal(gg.ZeroIndex(Type{0, 10, 0}), []int{0, 2})
	gtest.Equal(gg.ZeroIndex(Type{0, 10, 0, 20}), []int{0, 2})
	gtest.Equal(gg.ZeroIndex(Type{0, 10, 0, 20, 0}), []int{0, 2, 4})

	gtest.Equal(gg.ZeroIndex(Type{10, 0, 20, 0, 30}), []int{1, 3})
}

func TestNotZeroIndex(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.NotZeroIndex(Type(nil)))
	gtest.Zero(gg.NotZeroIndex(Type{}))
	gtest.Zero(gg.NotZeroIndex(Type{0}))
	gtest.Zero(gg.NotZeroIndex(Type{0, 0}))
	gtest.Zero(gg.NotZeroIndex(Type{0, 0, 0}))

	gtest.Equal(gg.NotZeroIndex(Type{10}), []int{0})
	gtest.Equal(gg.NotZeroIndex(Type{10, 20}), []int{0, 1})
	gtest.Equal(gg.NotZeroIndex(Type{10, 20, 30}), []int{0, 1, 2})

	gtest.Equal(gg.NotZeroIndex(Type{10, 0}), []int{0})
	gtest.Equal(gg.NotZeroIndex(Type{10, 0, 20}), []int{0, 2})
	gtest.Equal(gg.NotZeroIndex(Type{10, 0, 20, 0}), []int{0, 2})
	gtest.Equal(gg.NotZeroIndex(Type{10, 0, 20, 0, 30}), []int{0, 2, 4})

	gtest.Equal(gg.NotZeroIndex(Type{0, 10, 0, 20, 0}), []int{1, 3})
}

func TestCompact(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.Compact(Type(nil)))
	gtest.Zero(gg.Compact(Type{}))
	gtest.Zero(gg.Compact(Type{0}))
	gtest.Zero(gg.Compact(Type{0, 0}))
	gtest.Zero(gg.Compact(Type{0, 0, 0}))

	gtest.Equal(gg.Compact(Type{10}), Type{10})
	gtest.Equal(gg.Compact(Type{10, 20}), Type{10, 20})
	gtest.Equal(gg.Compact(Type{10, 20, 30}), Type{10, 20, 30})

	gtest.Equal(gg.Compact(Type{10, 0, 20, 0, 30}), Type{10, 20, 30})
	gtest.Equal(gg.Compact(Type{0, 10, 0, 20, 0}), Type{10, 20})
}

func TestFindIndex(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Eq(gg.FindIndex(Type(nil), nil), -1)
	gtest.Eq(gg.FindIndex(Type{}, nil), -1)
	gtest.Eq(gg.FindIndex(Type{0}, nil), -1)
	gtest.Eq(gg.FindIndex(Type{10}, nil), -1)
	gtest.Eq(gg.FindIndex(Type{10, 20}, nil), -1)
	gtest.Eq(gg.FindIndex(Type{10, 20, 30}, nil), -1)

	gtest.Eq(gg.FindIndex(Type{10}, False1[int]), -1)
	gtest.Eq(gg.FindIndex(Type{10, 20}, False1[int]), -1)
	gtest.Eq(gg.FindIndex(Type{10, 20, 30}, False1[int]), -1)

	gtest.Eq(gg.FindIndex(Type{10}, True1[int]), 0)
	gtest.Eq(gg.FindIndex(Type{10, 20}, True1[int]), 0)
	gtest.Eq(gg.FindIndex(Type{10, 20, 30}, True1[int]), 0)

	gtest.Eq(gg.FindIndex(Type{10}, gg.IsNeg[int]), -1)
	gtest.Eq(gg.FindIndex(Type{-10}, gg.IsNeg[int]), 0)
	gtest.Eq(gg.FindIndex(Type{-10, 10}, gg.IsNeg[int]), 0)
	gtest.Eq(gg.FindIndex(Type{10, -10}, gg.IsNeg[int]), 1)
	gtest.Eq(gg.FindIndex(Type{10, -10, 20}, gg.IsNeg[int]), 1)
	gtest.Eq(gg.FindIndex(Type{10, -10, 20, -20}, gg.IsNeg[int]), 1)
}

func TestFindLastIndex(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Eq(gg.FindLastIndex(Type(nil), nil), -1)
	gtest.Eq(gg.FindLastIndex(Type{}, nil), -1)
	gtest.Eq(gg.FindLastIndex(Type{0}, nil), -1)
	gtest.Eq(gg.FindLastIndex(Type{10}, nil), -1)
	gtest.Eq(gg.FindLastIndex(Type{10, 20}, nil), -1)
	gtest.Eq(gg.FindLastIndex(Type{10, 20, 30}, nil), -1)

	gtest.Eq(gg.FindLastIndex(Type{10}, False1[int]), -1)
	gtest.Eq(gg.FindLastIndex(Type{10, 20}, False1[int]), -1)
	gtest.Eq(gg.FindLastIndex(Type{10, 20, 30}, False1[int]), -1)

	gtest.Eq(gg.FindLastIndex(Type{10}, True1[int]), 0)
	gtest.Eq(gg.FindLastIndex(Type{10, 20}, True1[int]), 1)
	gtest.Eq(gg.FindLastIndex(Type{10, 20, 30}, True1[int]), 2)

	gtest.Eq(gg.FindLastIndex(Type{10}, gg.IsNeg[int]), -1)
	gtest.Eq(gg.FindLastIndex(Type{-10}, gg.IsNeg[int]), 0)
	gtest.Eq(gg.FindLastIndex(Type{-10, 10}, gg.IsNeg[int]), 0)
	gtest.Eq(gg.FindLastIndex(Type{10, -10}, gg.IsNeg[int]), 1)
	gtest.Eq(gg.FindLastIndex(Type{10, -10, 20}, gg.IsNeg[int]), 1)
	gtest.Eq(gg.FindLastIndex(Type{10, -10, 20, -20}, gg.IsNeg[int]), 3)
}

func TestFound(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.Tuple2(gg.Found(Type(nil), nil)))
	gtest.Zero(gg.Tuple2(gg.Found(Type{}, nil)))
	gtest.Zero(gg.Tuple2(gg.Found(Type{10}, nil)))
	gtest.Zero(gg.Tuple2(gg.Found(Type{10, 20}, nil)))
	gtest.Zero(gg.Tuple2(gg.Found(Type{10, 20, 30}, nil)))

	gtest.Zero(gg.Tuple2(gg.Found(Type(nil), False1[int])))
	gtest.Zero(gg.Tuple2(gg.Found(Type{}, False1[int])))
	gtest.Zero(gg.Tuple2(gg.Found(Type{10}, False1[int])))
	gtest.Zero(gg.Tuple2(gg.Found(Type{10, 20}, False1[int])))
	gtest.Zero(gg.Tuple2(gg.Found(Type{10, 20, 30}, False1[int])))

	gtest.Eq(
		gg.Tuple2(gg.Found(Type{10}, True1[int])),
		gg.Tuple2(10, true),
	)

	gtest.Eq(
		gg.Tuple2(gg.Found(Type{10, 20}, True1[int])),
		gg.Tuple2(10, true),
	)

	gtest.Eq(
		gg.Tuple2(gg.Found(Type{-10, 10, -20, 20}, gg.IsNeg[int])),
		gg.Tuple2(-10, true),
	)

	gtest.Eq(
		gg.Tuple2(gg.Found(Type{-10, 10, -20, 20}, gg.IsPos[int])),
		gg.Tuple2(10, true),
	)
}

func TestFoundLast(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.Tuple2(gg.FoundLast(Type(nil), nil)))
	gtest.Zero(gg.Tuple2(gg.FoundLast(Type{}, nil)))
	gtest.Zero(gg.Tuple2(gg.FoundLast(Type{10}, nil)))
	gtest.Zero(gg.Tuple2(gg.FoundLast(Type{10, 20}, nil)))
	gtest.Zero(gg.Tuple2(gg.FoundLast(Type{10, 20, 30}, nil)))

	gtest.Zero(gg.Tuple2(gg.FoundLast(Type(nil), False1[int])))
	gtest.Zero(gg.Tuple2(gg.FoundLast(Type{}, False1[int])))
	gtest.Zero(gg.Tuple2(gg.FoundLast(Type{10}, False1[int])))
	gtest.Zero(gg.Tuple2(gg.FoundLast(Type{10, 20}, False1[int])))
	gtest.Zero(gg.Tuple2(gg.FoundLast(Type{10, 20, 30}, False1[int])))

	gtest.Eq(
		gg.Tuple2(gg.FoundLast(Type{10}, True1[int])),
		gg.Tuple2(10, true),
	)

	gtest.Eq(
		gg.Tuple2(gg.FoundLast(Type{10, 20}, True1[int])),
		gg.Tuple2(20, true),
	)

	gtest.Eq(
		gg.Tuple2(gg.FoundLast(Type{-10, 10, -20, 20}, gg.IsNeg[int])),
		gg.Tuple2(-20, true),
	)

	gtest.Eq(
		gg.Tuple2(gg.FoundLast(Type{-10, 10, -20, 20}, gg.IsPos[int])),
		gg.Tuple2(20, true),
	)
}

func TestFind(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.Find(Type(nil), nil))
	gtest.Zero(gg.Find(Type{}, nil))
	gtest.Zero(gg.Find(Type{10}, nil))
	gtest.Zero(gg.Find(Type{10, 20}, nil))
	gtest.Zero(gg.Find(Type{10, 20, 30}, nil))

	gtest.Zero(gg.Find(Type(nil), False1[int]))
	gtest.Zero(gg.Find(Type{}, False1[int]))
	gtest.Zero(gg.Find(Type{10}, False1[int]))
	gtest.Zero(gg.Find(Type{10, 20}, False1[int]))
	gtest.Zero(gg.Find(Type{10, 20, 30}, False1[int]))

	gtest.Eq(gg.Find(Type{10}, True1[int]), 10)
	gtest.Eq(gg.Find(Type{10, 20}, True1[int]), 10)
	gtest.Eq(gg.Find(Type{-10, 10, -20, 20}, gg.IsNeg[int]), -10)
	gtest.Eq(gg.Find(Type{-10, 10, -20, 20}, gg.IsPos[int]), 10)
}

func TestFindLast(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.FindLast(Type(nil), nil))
	gtest.Zero(gg.FindLast(Type{}, nil))
	gtest.Zero(gg.FindLast(Type{10}, nil))
	gtest.Zero(gg.FindLast(Type{10, 20}, nil))
	gtest.Zero(gg.FindLast(Type{10, 20, 30}, nil))

	gtest.Zero(gg.FindLast(Type(nil), False1[int]))
	gtest.Zero(gg.FindLast(Type{}, False1[int]))
	gtest.Zero(gg.FindLast(Type{10}, False1[int]))
	gtest.Zero(gg.FindLast(Type{10, 20}, False1[int]))
	gtest.Zero(gg.FindLast(Type{10, 20, 30}, False1[int]))

	gtest.Eq(gg.FindLast(Type{10}, True1[int]), 10)
	gtest.Eq(gg.FindLast(Type{10, 20}, True1[int]), 20)
	gtest.Eq(gg.FindLast(Type{-10, 10, -20, 20}, gg.IsNeg[int]), -20)
	gtest.Eq(gg.FindLast(Type{-10, 10, -20, 20}, gg.IsPos[int]), 20)
}

func TestProcured(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.Tuple2(gg.Procured[int](Type(nil), nil)))

	gtest.Zero(gg.Tuple2(gg.Procured(Type{10, 20, 30}, Id1False[int])))

	gtest.Eq(
		gg.Tuple2(gg.Procured(Type{0, 10}, Id1True[int])),
		gg.Tuple2(0, true),
	)

	gtest.Eq(
		gg.Tuple2(gg.Procured(Type{10, 20, 30}, Id1True[int])),
		gg.Tuple2(10, true),
	)
}

func TestProcure(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.Procure[int](Type(nil), nil))
	gtest.Zero(gg.Procure(Type{-1}, gg.Inc[int]))
	gtest.Zero(gg.Procure(Type{1}, gg.Dec[int]))

	gtest.Eq(gg.Procure(Type{10}, gg.Inc[int]), 11)
	gtest.Eq(gg.Procure(Type{-1, 10}, gg.Inc[int]), 11)
	gtest.Eq(gg.Procure(Type{-1, -1, 10}, gg.Inc[int]), 11)
	gtest.Eq(gg.Procure(Type{-1, -1, 10, 20}, gg.Inc[int]), 11)

	gtest.Eq(gg.Procure(Type{10}, gg.Dec[int]), 9)
	gtest.Eq(gg.Procure(Type{1, 10}, gg.Dec[int]), 9)
	gtest.Eq(gg.Procure(Type{1, 1, 10}, gg.Dec[int]), 9)
	gtest.Eq(gg.Procure(Type{1, 1, 10, 20}, gg.Dec[int]), 9)
}

func TestAdjoin(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Equal(gg.Adjoin(Type(nil), 0), Type{0})
	gtest.Equal(gg.Adjoin(Type(nil), 10), Type{10})

	gtest.Equal(gg.Adjoin(Type{10, 20, 30}, 10), Type{10, 20, 30})
	gtest.Equal(gg.Adjoin(Type{10, 20, 30}, 20), Type{10, 20, 30})
	gtest.Equal(gg.Adjoin(Type{10, 20, 30}, 30), Type{10, 20, 30})
	gtest.Equal(gg.Adjoin(Type{10, 20, 30}, 0), Type{10, 20, 30, 0})
	gtest.Equal(gg.Adjoin(Type{10, 20, 30}, 40), Type{10, 20, 30, 40})
}

func TestExclude(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.Exclude(Type(nil), 0))
	gtest.Zero(gg.Exclude(Type(nil), 10))
	gtest.Zero(gg.Exclude(Type{}, 0))
	gtest.Zero(gg.Exclude(Type{}, 10))
	gtest.Zero(gg.Exclude(Type{0}, 0))
	gtest.Zero(gg.Exclude(Type{0, 0}, 0))
	gtest.Zero(gg.Exclude(Type{10}, 10))
	gtest.Zero(gg.Exclude(Type{10, 10}, 10))

	gtest.Equal(gg.Exclude(Type{10, 20, 30}, 40), Type{10, 20, 30})
	gtest.Equal(gg.Exclude(Type{10, 20, 30}, 10), Type{20, 30})
	gtest.Equal(gg.Exclude(Type{10, 20, 30}, 20), Type{10, 30})
	gtest.Equal(gg.Exclude(Type{10, 20, 30}, 30), Type{10, 20})
}

func TestExcludeFrom(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.ExcludeFrom(Type(nil)))
	gtest.Zero(gg.ExcludeFrom(Type{}))

	gtest.Zero(gg.ExcludeFrom(Type{10}, Type{10}))

	gtest.Zero(gg.ExcludeFrom(Type{10, 20}, Type{10, 20}))
	gtest.Zero(gg.ExcludeFrom(Type{10, 20}, Type{10, 20, 30}))
	gtest.Zero(gg.ExcludeFrom(Type{10, 20}, Type{10, 20}, Type{30}))
	gtest.Zero(gg.ExcludeFrom(Type{10, 20}, Type{10}, Type{20, 30}))
	gtest.Zero(gg.ExcludeFrom(Type{10, 20}, Type{10}, Type{20}))
	gtest.Zero(gg.ExcludeFrom(Type{10, 20}, Type{10}, Type{20}, Type{30}))

	gtest.Equal(
		gg.ExcludeFrom(Type{10, 20, 30}, Type{10}),
		Type{20, 30},
	)

	gtest.Equal(
		gg.ExcludeFrom(Type{10, 20, 30}, Type{20}),
		Type{10, 30},
	)

	gtest.Equal(
		gg.ExcludeFrom(Type{10, 20, 30}, Type{20}, Type{10}),
		Type{30},
	)
}

func BenchmarkSubtract(b *testing.B) {
	base := []int{10, 20, 30, 40, 50, 60}
	sub := [][]int{{10, 20}, {50}}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ExcludeFrom(base, sub...))
	}
}

func TestIntersect(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.Intersect(Type(nil), Type(nil)))
	gtest.Zero(gg.Intersect(Type{}, Type(nil)))
	gtest.Zero(gg.Intersect(Type(nil), Type{}))
	gtest.Zero(gg.Intersect(Type{}, Type{}))
	gtest.Zero(gg.Intersect(Type{10, 20, 30}, Type(nil)))
	gtest.Zero(gg.Intersect(Type(nil), Type{10, 20, 30}))
	gtest.Zero(gg.Intersect(Type{10}, Type{20}))
	gtest.Zero(gg.Intersect(Type{10, 20}, Type{30, 40}))

	gtest.Equal(gg.Intersect(Type{10, 20, 30}, Type{10}), Type{10})
	gtest.Equal(gg.Intersect(Type{10, 20, 30}, Type{10, 20}), Type{10, 20})
	gtest.Equal(gg.Intersect(Type{10, 20, 30}, Type{10, 20, 30}), Type{10, 20, 30})

	gtest.Equal(gg.Intersect(Type{10}, Type{10, 20, 30}), Type{10})
	gtest.Equal(gg.Intersect(Type{10, 20}, Type{10, 20, 30}), Type{10, 20})
	gtest.Equal(gg.Intersect(Type{10, 20, 30}, Type{10, 20, 30}), Type{10, 20, 30})

	gtest.Equal(gg.Intersect(Type{10, 20}, Type{-10, 20, 30}), Type{20})
	gtest.Equal(gg.Intersect(Type{10, 20, 30}, Type{-10, 20, 30, 40}), Type{20, 30})
}

func TestUnion(t *testing.T) {
	defer gtest.Catch(t)

	type Elem = int
	type Slice = []Elem

	gtest.Zero(gg.Union[Slice]())
	gtest.Zero(gg.Union[Slice](nil))
	gtest.Zero(gg.Union[Slice](nil, nil))
	gtest.Zero(gg.Union[Slice](nil, nil, Slice{}, nil, Slice{}))

	gtest.Equal(gg.Union(Slice{10}), Slice{10})
	gtest.Equal(gg.Union(Slice{10, 10}), Slice{10})
	gtest.Equal(gg.Union(nil, Slice{10, 10}), Slice{10})
	gtest.Equal(gg.Union(Slice{10, 10}, nil), Slice{10})
	gtest.Equal(gg.Union(nil, Slice{10, 10}, nil), Slice{10})

	gtest.Equal(gg.Union(Slice{10}, Slice{10}), Slice{10})
	gtest.Equal(gg.Union(Slice{10, 20}, Slice{10}), Slice{10, 20})
	gtest.Equal(gg.Union(Slice{10, 20}, Slice{10, 20}), Slice{10, 20})
	gtest.Equal(gg.Union(Slice{10, 20}, Slice{20, 10}), Slice{10, 20})
	gtest.Equal(gg.Union(Slice{10, 20}, Slice{20, 10, 30}), Slice{10, 20, 30})
	gtest.Equal(gg.Union(Slice{10, 20}, Slice{20, 10}, Slice{30, 20, 10}), Slice{10, 20, 30})
	gtest.Equal(gg.Union(Slice{}, Slice{20, 10}, Slice{30, 20, 10}), Slice{20, 10, 30})
}

func BenchmarkUnion(b *testing.B) {
	src := [][]int{{10, 20}, {30, 40}, {50, 60}}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Union(src...))
	}
}

func TestUniq(t *testing.T) {
	defer gtest.Catch(t)

	testUniqCommon(gg.Uniq[[]int])
}

func testUniqCommon(fun func([]int) []int) {
	type Slice = []int

	gtest.Zero(fun(Slice(nil)))
	gtest.Zero(fun(Slice{}))

	gtest.Equal(fun(Slice{10}), Slice{10})
	gtest.Equal(fun(Slice{10, 10}), Slice{10})
	gtest.Equal(fun(Slice{10, 10, 10}), Slice{10})
	gtest.Equal(fun(Slice{10, 10, 10, 20}), Slice{10, 20})
	gtest.Equal(fun(Slice{10, 10, 10, 20, 20}), Slice{10, 20})
	gtest.Equal(fun(Slice{10, 20, 20, 10}), Slice{10, 20})
	gtest.Equal(fun(Slice{30, 10, 20, 20, 10, 30}), Slice{30, 10, 20})
}

func TestUniqBy(t *testing.T) {
	defer gtest.Catch(t)

	type Elem = int
	type Slice = []Elem

	testUniqCommon(func(src Slice) Slice {
		return gg.UniqBy(src, gg.Id1[Elem])
	})

	testUniqCommon(func(src Slice) Slice {
		return gg.UniqBy(src, gg.String[Elem])
	})

	gtest.Zero(gg.UniqBy[Slice, Elem, string](Slice(nil), nil))
	gtest.Zero(gg.UniqBy[Slice, Elem, string](Slice{}, nil))
	gtest.Zero(gg.UniqBy[Slice, Elem, string](Slice{10}, nil))
	gtest.Zero(gg.UniqBy[Slice, Elem, string](Slice{10, 20}, nil))
	gtest.Zero(gg.UniqBy[Slice, Elem, string](Slice{10, 20, 30}, nil))

	gtest.Equal(
		gg.UniqBy(
			Slice{10, 20, 10, 30, 10},
			func(Elem) struct{} { return struct{}{} },
		),
		Slice{10},
	)

	gtest.Equal(
		gg.UniqBy(
			Slice{20, 10, 20, 30, 30},
			func(Elem) struct{} { return struct{}{} },
		),
		Slice{20},
	)
}

func TestHas(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.False(gg.Has(Type(nil), 0))
	gtest.False(gg.Has(Type(nil), 10))
	gtest.False(gg.Has(Type{}, 0))
	gtest.False(gg.Has(Type{}, 10))
	gtest.False(gg.Has(Type{10, 20, 30}, 40))
	gtest.False(gg.Has(Type{10, 20, 30}, 0))
	gtest.False(gg.Has(Type{0, 10, 0, 20, 30}, 40))

	gtest.True(gg.Has(Type{10, 20, 30}, 10))
	gtest.True(gg.Has(Type{10, 20, 30}, 20))
	gtest.True(gg.Has(Type{10, 20, 30}, 30))
	gtest.True(gg.Has(Type{0, 10, 0, 20}, 0))
}

func BenchmarkHas(b *testing.B) {
	defer gtest.Catch(b)

	tar := gg.Times(128, gg.Inc[int])
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Has(tar, 256))
	}
}

func TestHasEvery(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.True(gg.HasEvery(Type(nil), nil))
	gtest.True(gg.HasEvery(Type(nil), Type{}))
	gtest.True(gg.HasEvery(Type{}, nil))
	gtest.True(gg.HasEvery(Type{}, Type{}))

	gtest.True(gg.HasEvery(Type{10, 20, 30, 40}, Type{10}))
	gtest.True(gg.HasEvery(Type{10, 20, 30, 40}, Type{20}))
	gtest.True(gg.HasEvery(Type{10, 20, 30, 40}, Type{30}))
	gtest.True(gg.HasEvery(Type{10, 20, 30, 40}, Type{10, 20}))
	gtest.True(gg.HasEvery(Type{10, 20, 30, 40}, Type{20, 30}))
	gtest.True(gg.HasEvery(Type{10, 20, 30, 40}, Type{10, 30}))
	gtest.True(gg.HasEvery(Type{10, 20, 30, 40}, Type{10, 20, 30}))

	gtest.False(gg.HasEvery(Type(nil), Type{10}))
	gtest.False(gg.HasEvery(Type{}, Type{10}))
	gtest.False(gg.HasEvery(Type{10}, Type{20}))
	gtest.False(gg.HasEvery(Type{10, 20}, Type{20, 30}))
}

func TestHasSome(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.False(gg.HasSome(Type(nil), nil))
	gtest.False(gg.HasSome(Type(nil), Type{}))
	gtest.False(gg.HasSome(Type{}, nil))
	gtest.False(gg.HasSome(Type{}, Type{}))

	gtest.True(gg.HasSome(Type{10, 20, 30, 40}, Type{10}))
	gtest.True(gg.HasSome(Type{10, 20, 30, 40}, Type{20}))
	gtest.True(gg.HasSome(Type{10, 20, 30, 40}, Type{30}))
	gtest.True(gg.HasSome(Type{10, 20, 30, 40}, Type{10, 20}))
	gtest.True(gg.HasSome(Type{10, 20, 30, 40}, Type{20, 30}))
	gtest.True(gg.HasSome(Type{10, 20, 30, 40}, Type{10, 30}))
	gtest.True(gg.HasSome(Type{10, 20, 30, 40}, Type{10, 20, 30}))

	gtest.False(gg.HasSome(Type(nil), Type{10}))
	gtest.False(gg.HasSome(Type{}, Type{10}))
	gtest.False(gg.HasSome(Type{10}, Type{20}))
	gtest.False(gg.HasSome(Type{10, 20}, Type{}))
	gtest.False(gg.HasSome(Type{10, 20}, Type{30}))
}

func TestHasNone(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.True(gg.HasNone(Type(nil), nil))
	gtest.True(gg.HasNone(Type(nil), Type{}))
	gtest.True(gg.HasNone(Type{}, nil))
	gtest.True(gg.HasNone(Type{}, Type{}))

	gtest.False(gg.HasNone(Type{10, 20, 30, 40}, Type{10}))
	gtest.False(gg.HasNone(Type{10, 20, 30, 40}, Type{20}))
	gtest.False(gg.HasNone(Type{10, 20, 30, 40}, Type{30}))
	gtest.False(gg.HasNone(Type{10, 20, 30, 40}, Type{10, 20}))
	gtest.False(gg.HasNone(Type{10, 20, 30, 40}, Type{20, 30}))
	gtest.False(gg.HasNone(Type{10, 20, 30, 40}, Type{10, 30}))
	gtest.False(gg.HasNone(Type{10, 20, 30, 40}, Type{10, 20, 30}))

	gtest.True(gg.HasNone(Type(nil), Type{10}))
	gtest.True(gg.HasNone(Type{}, Type{10}))
	gtest.True(gg.HasNone(Type{10}, Type{20}))
	gtest.True(gg.HasNone(Type{10, 20}, Type{}))
	gtest.True(gg.HasNone(Type{10, 20}, Type{30}))
}

func TestSome(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.False(gg.Some(Type(nil), False1[int]))
	gtest.False(gg.Some(Type{}, False1[int]))

	gtest.False(gg.Some(Type(nil), True1[int]))
	gtest.False(gg.Some(Type{}, True1[int]))

	gtest.False(gg.Some(Type{10}, False1[int]))
	gtest.False(gg.Some(Type{10, 20}, False1[int]))
	gtest.False(gg.Some(Type{10, 20, 30}, False1[int]))

	gtest.True(gg.Some(Type{10}, True1[int]))
	gtest.True(gg.Some(Type{10, 20}, True1[int]))
	gtest.True(gg.Some(Type{10, 20, 30}, True1[int]))

	gtest.False(gg.Some(Type{10, 20, 30}, gg.IsNeg[int]))
	gtest.True(gg.Some(Type{10, 20, 30}, gg.IsPos[int]))

	gtest.True(gg.Some(Type{-10, 10, -20, 20}, gg.IsNeg[int]))
	gtest.True(gg.Some(Type{-10, 10, -20, 20}, gg.IsPos[int]))
}

func TestNone(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.True(gg.None(Type(nil), nil))
	gtest.True(gg.None(Type{10, 20}, nil))

	gtest.True(gg.None(Type(nil), False1[int]))
	gtest.True(gg.None(Type{}, False1[int]))

	gtest.True(gg.None(Type(nil), True1[int]))
	gtest.True(gg.None(Type{}, True1[int]))

	gtest.True(gg.None(Type{10}, False1[int]))
	gtest.True(gg.None(Type{10, 20}, False1[int]))
	gtest.True(gg.None(Type{10, 20, 30}, False1[int]))

	gtest.False(gg.None(Type{10}, True1[int]))
	gtest.False(gg.None(Type{10, 20}, True1[int]))
	gtest.False(gg.None(Type{10, 20, 30}, True1[int]))

	gtest.True(gg.None(Type{10, 20, 30}, gg.IsNeg[int]))
	gtest.False(gg.None(Type{10, 20, 30}, gg.IsPos[int]))

	gtest.False(gg.None(Type{-10, 10, -20, 20}, gg.IsNeg[int]))
	gtest.False(gg.None(Type{-10, 10, -20, 20}, gg.IsPos[int]))
}

func TestEvery(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.True(gg.Every(Type(nil), False1[int]))
	gtest.True(gg.Every(Type{}, False1[int]))

	gtest.True(gg.Every(Type(nil), True1[int]))
	gtest.True(gg.Every(Type{}, True1[int]))

	gtest.False(gg.Every(Type{10}, False1[int]))
	gtest.False(gg.Every(Type{10, 20}, False1[int]))
	gtest.False(gg.Every(Type{10, 20, 30}, False1[int]))

	gtest.True(gg.Every(Type{10}, True1[int]))
	gtest.True(gg.Every(Type{10, 20}, True1[int]))
	gtest.True(gg.Every(Type{10, 20, 30}, True1[int]))

	gtest.False(gg.Every(Type{10, 20, 30}, gg.IsNeg[int]))
	gtest.True(gg.Every(Type{10, 20, 30}, gg.IsPos[int]))

	gtest.False(gg.Every(Type{-10, 10, -20, 20}, gg.IsNeg[int]))
	gtest.False(gg.Every(Type{-10, 10, -20, 20}, gg.IsPos[int]))
}

func TestConcat(t *testing.T) {
	defer gtest.Catch(t)

	type Type = []int

	gtest.Zero(gg.Concat[Type]())
	gtest.Zero(gg.Concat[Type](nil))
	gtest.Zero(gg.Concat[Type](nil, nil))
	gtest.Zero(gg.Concat[Type](nil, nil, nil))

	gtest.Equal(gg.Concat(Type{}), Type{})
	gtest.Equal(gg.Concat(Type{}, nil), Type{})
	gtest.Equal(gg.Concat(nil, Type{}), Type{})
	gtest.Equal(gg.Concat(Type{}, Type{}), Type{})

	gtest.Equal(gg.Concat(Type{10}, Type{}), Type{10})
	gtest.Equal(gg.Concat(Type{10, 20}, Type{}), Type{10, 20})
	gtest.Equal(gg.Concat(Type{10}, Type{20}), Type{10, 20})
	gtest.Equal(gg.Concat(Type{10}, Type{20}, Type{30}), Type{10, 20, 30})

	gtest.Equal(
		gg.Concat(Type{10, 20}, Type{20, 30}, Type{30, 40}),
		Type{10, 20, 20, 30, 30, 40},
	)
}

func BenchmarkConcat(b *testing.B) {
	src := [][]int{{10, 20}, {30, 40}, {50, 60}}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Concat(src...))
	}
}

func TestPrimSorted(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(
		gg.SortedPrim(gg.ToSlice(20, 30, 10, 40)),
		gg.ToSlice(10, 20, 30, 40),
	)
}

func BenchmarkPrimSorted(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.SortedPrim(gg.ToSlice(20, 30, 10, 40)))
	}
}

func TestReversed(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.Reversed([]int(nil)), []int(nil))
	gtest.Equal(gg.Reversed([]int{}), []int{})
	gtest.Equal(gg.Reversed([]int{10}), []int{10})

	gtest.Equal(gg.Reversed([]int{10, 20}), []int{20, 10})
	gtest.Equal(gg.Reversed([]int{20, 10}), []int{10, 20})

	gtest.Equal(gg.Reversed([]int{10, 20, 30}), []int{30, 20, 10})
	gtest.Equal(gg.Reversed([]int{30, 20, 10}), []int{10, 20, 30})

	gtest.Equal(gg.Reversed([]int{10, 20, 30, 40}), []int{40, 30, 20, 10})
	gtest.Equal(gg.Reversed([]int{40, 30, 20, 10}), []int{10, 20, 30, 40})
}

func BenchmarkReversed(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Reversed([]int{20, 30, 10, 40}))
	}
}

func BenchmarkTakeWhile(b *testing.B) {
	val := []int{-30, -20, -10, 0, 10, 20, 30}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.TakeWhile(val, gg.IsNeg[int]))
	}
}

func TestMinPrim(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.MinPrim[int](), 0)
	gtest.Eq(gg.MinPrim[float64](), 0.0)
	gtest.Eq(gg.MinPrim[string](), ``)

	gtest.Eq(gg.MinPrim(-10), -10)
	gtest.Eq(gg.MinPrim(0), 0)
	gtest.Eq(gg.MinPrim(10), 10)

	gtest.Eq(gg.MinPrim(-10.5), -10.5)
	gtest.Eq(gg.MinPrim(0.0), 0.0)
	gtest.Eq(gg.MinPrim(10.5), 10.5)

	gtest.Eq(gg.MinPrim(`str`), `str`)

	gtest.Eq(gg.MinPrim(-10, 0), -10)
	gtest.Eq(gg.MinPrim(0, 10), 0)
	gtest.Eq(gg.MinPrim(-10, 10), -10)

	gtest.Eq(gg.MinPrim(``, `10`), ``)
	gtest.Eq(gg.MinPrim(`10`, ``), ``)
	gtest.Eq(gg.MinPrim(`10`, `20`), `10`)
	gtest.Eq(gg.MinPrim(`20`, `10`), `10`)

	gtest.Eq(gg.MinPrim(-20, -10, 0), -20)
	gtest.Eq(gg.MinPrim(0, 10, 20), 0)
	gtest.Eq(gg.MinPrim(-10, 0, 10), -10)

	gtest.Eq(gg.MinPrim(``, `20`, `10`), ``)
	gtest.Eq(gg.MinPrim(``, `10`, `20`), ``)
	gtest.Eq(gg.MinPrim(`10`, ``, `20`), ``)
	gtest.Eq(gg.MinPrim(`10`, `20`, ``), ``)

	gtest.Eq(gg.MinPrim(`10`, `20`, `30`), `10`)
	gtest.Eq(gg.MinPrim(`20`, `10`, `30`), `10`)
	gtest.Eq(gg.MinPrim(`30`, `20`, `10`), `10`)
}

func BenchmarkMinPrim(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.MinPrim(ind-1, ind, ind+1))
	}
}

func TestMaxPrim(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.MaxPrim[int](), 0)
	gtest.Eq(gg.MaxPrim[float64](), 0.0)
	gtest.Eq(gg.MaxPrim[string](), ``)

	gtest.Eq(gg.MaxPrim(-10), -10)
	gtest.Eq(gg.MaxPrim(0), 0)
	gtest.Eq(gg.MaxPrim(10), 10)

	gtest.Eq(gg.MaxPrim(-10.5), -10.5)
	gtest.Eq(gg.MaxPrim(0.0), 0.0)
	gtest.Eq(gg.MaxPrim(10.5), 10.5)

	gtest.Eq(gg.MaxPrim(`str`), `str`)

	gtest.Eq(gg.MaxPrim(-10, 0), 0)
	gtest.Eq(gg.MaxPrim(0, 10), 10)
	gtest.Eq(gg.MaxPrim(-10, 10), 10)

	gtest.Eq(gg.MaxPrim(``, `10`), `10`)
	gtest.Eq(gg.MaxPrim(`10`, ``), `10`)
	gtest.Eq(gg.MaxPrim(`10`, `20`), `20`)
	gtest.Eq(gg.MaxPrim(`20`, `10`), `20`)

	gtest.Eq(gg.MaxPrim(-20, -10, 0), 0)
	gtest.Eq(gg.MaxPrim(0, 10, 20), 20)
	gtest.Eq(gg.MaxPrim(-10, 0, 10), 10)

	gtest.Eq(gg.MaxPrim(``, `20`, `10`), `20`)
	gtest.Eq(gg.MaxPrim(``, `10`, `20`), `20`)
	gtest.Eq(gg.MaxPrim(`10`, ``, `20`), `20`)
	gtest.Eq(gg.MaxPrim(`10`, `20`, ``), `20`)

	gtest.Eq(gg.MaxPrim(`10`, `20`, `30`), `30`)
	gtest.Eq(gg.MaxPrim(`20`, `10`, `30`), `30`)
	gtest.Eq(gg.MaxPrim(`30`, `20`, `10`), `30`)
}

func BenchmarkMaxPrim(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.MaxPrim(ind-1, ind, ind+1))
	}
}

func TestMin(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Min[Comparer[int]](), ComparerOf(0))
	gtest.Eq(gg.Min[Comparer[float64]](), ComparerOf(0.0))
	gtest.Eq(gg.Min[Comparer[string]](), ComparerOf(``))

	gtest.Eq(gg.Min(ComparerOf(-10)), ComparerOf(-10))
	gtest.Eq(gg.Min(ComparerOf(0)), ComparerOf(0))
	gtest.Eq(gg.Min(ComparerOf(10)), ComparerOf(10))

	gtest.Eq(gg.Min(ComparerOf(-10.5)), ComparerOf(-10.5))
	gtest.Eq(gg.Min(ComparerOf(0.0)), ComparerOf(0.0))
	gtest.Eq(gg.Min(ComparerOf(10.5)), ComparerOf(10.5))

	gtest.Eq(gg.Min(ComparerOf(`str`)), ComparerOf(`str`))

	gtest.Eq(gg.Min(ComparerOf(-10), ComparerOf(0)), ComparerOf(-10))
	gtest.Eq(gg.Min(ComparerOf(0), ComparerOf(10)), ComparerOf(0))
	gtest.Eq(gg.Min(ComparerOf(-10), ComparerOf(10)), ComparerOf(-10))

	gtest.Eq(gg.Min(ComparerOf(``), ComparerOf(`10`)), ComparerOf(``))
	gtest.Eq(gg.Min(ComparerOf(`10`), ComparerOf(``)), ComparerOf(``))
	gtest.Eq(gg.Min(ComparerOf(`10`), ComparerOf(`20`)), ComparerOf(`10`))
	gtest.Eq(gg.Min(ComparerOf(`20`), ComparerOf(`10`)), ComparerOf(`10`))

	gtest.Eq(gg.Min(ComparerOf(-20), ComparerOf(-10), ComparerOf(0)), ComparerOf(-20))
	gtest.Eq(gg.Min(ComparerOf(0), ComparerOf(10), ComparerOf(20)), ComparerOf(0))
	gtest.Eq(gg.Min(ComparerOf(-10), ComparerOf(0), ComparerOf(10)), ComparerOf(-10))

	gtest.Eq(gg.Min(ComparerOf(``), ComparerOf(`20`), ComparerOf(`10`)), ComparerOf(``))
	gtest.Eq(gg.Min(ComparerOf(``), ComparerOf(`10`), ComparerOf(`20`)), ComparerOf(``))
	gtest.Eq(gg.Min(ComparerOf(`10`), ComparerOf(``), ComparerOf(`20`)), ComparerOf(``))
	gtest.Eq(gg.Min(ComparerOf(`10`), ComparerOf(`20`), ComparerOf(``)), ComparerOf(``))

	gtest.Eq(gg.Min(ComparerOf(`10`), ComparerOf(`20`), ComparerOf(`30`)), ComparerOf(`10`))
	gtest.Eq(gg.Min(ComparerOf(`20`), ComparerOf(`10`), ComparerOf(`30`)), ComparerOf(`10`))
	gtest.Eq(gg.Min(ComparerOf(`30`), ComparerOf(`20`), ComparerOf(`10`)), ComparerOf(`10`))
}

func TestMax(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Max[Comparer[int]](), ComparerOf(0))
	gtest.Eq(gg.Max[Comparer[float64]](), ComparerOf(0.0))
	gtest.Eq(gg.Max[Comparer[string]](), ComparerOf(``))

	gtest.Eq(gg.Max(ComparerOf(-10)), ComparerOf(-10))
	gtest.Eq(gg.Max(ComparerOf(0)), ComparerOf(0))
	gtest.Eq(gg.Max(ComparerOf(10)), ComparerOf(10))

	gtest.Eq(gg.Max(ComparerOf(-10.5)), ComparerOf(-10.5))
	gtest.Eq(gg.Max(ComparerOf(0.0)), ComparerOf(0.0))
	gtest.Eq(gg.Max(ComparerOf(10.5)), ComparerOf(10.5))

	gtest.Eq(gg.Max(ComparerOf(`str`)), ComparerOf(`str`))

	gtest.Eq(gg.Max(ComparerOf(-10), ComparerOf(0)), ComparerOf(0))
	gtest.Eq(gg.Max(ComparerOf(0), ComparerOf(10)), ComparerOf(10))
	gtest.Eq(gg.Max(ComparerOf(-10), ComparerOf(10)), ComparerOf(10))

	gtest.Eq(gg.Max(ComparerOf(``), ComparerOf(`10`)), ComparerOf(`10`))
	gtest.Eq(gg.Max(ComparerOf(`10`), ComparerOf(``)), ComparerOf(`10`))
	gtest.Eq(gg.Max(ComparerOf(`10`), ComparerOf(`20`)), ComparerOf(`20`))
	gtest.Eq(gg.Max(ComparerOf(`20`), ComparerOf(`10`)), ComparerOf(`20`))

	gtest.Eq(gg.Max(ComparerOf(-20), ComparerOf(-10), ComparerOf(0)), ComparerOf(0))
	gtest.Eq(gg.Max(ComparerOf(0), ComparerOf(10), ComparerOf(20)), ComparerOf(20))
	gtest.Eq(gg.Max(ComparerOf(-10), ComparerOf(0), ComparerOf(10)), ComparerOf(10))

	gtest.Eq(gg.Max(ComparerOf(``), ComparerOf(`20`), ComparerOf(`10`)), ComparerOf(`20`))
	gtest.Eq(gg.Max(ComparerOf(``), ComparerOf(`10`), ComparerOf(`20`)), ComparerOf(`20`))
	gtest.Eq(gg.Max(ComparerOf(`10`), ComparerOf(``), ComparerOf(`20`)), ComparerOf(`20`))
	gtest.Eq(gg.Max(ComparerOf(`10`), ComparerOf(`20`), ComparerOf(``)), ComparerOf(`20`))

	gtest.Eq(gg.Max(ComparerOf(`10`), ComparerOf(`20`), ComparerOf(`30`)), ComparerOf(`30`))
	gtest.Eq(gg.Max(ComparerOf(`20`), ComparerOf(`10`), ComparerOf(`30`)), ComparerOf(`30`))
	gtest.Eq(gg.Max(ComparerOf(`30`), ComparerOf(`20`), ComparerOf(`10`)), ComparerOf(`30`))
}

func TestMinPrimBy(t *testing.T) {
	defer gtest.Catch(t)

	type Out = int
	type Src = Comparer[Out]
	var fun = Src.Get

	gtest.Zero(gg.MinPrimBy[Src, Out](nil, nil))
	gtest.Zero(gg.MinPrimBy[Src, Out]([]Src{}, nil))
	gtest.Zero(gg.MinPrimBy[Src, Out]([]Src{{10}}, nil))
	gtest.Zero(gg.MinPrimBy[Src, Out]([]Src{{10}, {20}}, nil))

	gtest.Zero(gg.MinPrimBy[Src, Out](nil, fun))
	gtest.Zero(gg.MinPrimBy([]Src{}, fun))

	gtest.Eq(gg.MinPrimBy([]Src{{-10}}, fun), -10)
	gtest.Eq(gg.MinPrimBy([]Src{{0}}, fun), 0)
	gtest.Eq(gg.MinPrimBy([]Src{{10}}, fun), 10)

	gtest.Eq(gg.MinPrimBy([]Src{{0}, {-10}}, fun), -10)
	gtest.Eq(gg.MinPrimBy([]Src{{-10}, {0}}, fun), -10)
	gtest.Eq(gg.MinPrimBy([]Src{{0}, {10}}, fun), 0)
	gtest.Eq(gg.MinPrimBy([]Src{{10}, {0}}, fun), 0)
	gtest.Eq(gg.MinPrimBy([]Src{{-10}, {10}}, fun), -10)
	gtest.Eq(gg.MinPrimBy([]Src{{10}, {-10}}, fun), -10)

	gtest.Eq(gg.MinPrimBy([]Src{{-10}, {-20}, {0}}, fun), -20)
	gtest.Eq(gg.MinPrimBy([]Src{{0}, {-10}, {10}}, fun), -10)
	gtest.Eq(gg.MinPrimBy([]Src{{10}, {0}, {20}}, fun), 0)
}

func TestMaxPrimBy(t *testing.T) {
	defer gtest.Catch(t)

	type Out = int
	type Src = Comparer[Out]
	var fun = Src.Get

	gtest.Zero(gg.MaxPrimBy[Src, Out](nil, nil))
	gtest.Zero(gg.MaxPrimBy[Src, Out]([]Src{}, nil))
	gtest.Zero(gg.MaxPrimBy[Src, Out]([]Src{{10}}, nil))
	gtest.Zero(gg.MaxPrimBy[Src, Out]([]Src{{10}, {20}}, nil))

	gtest.Zero(gg.MaxPrimBy[Src, Out](nil, fun))
	gtest.Zero(gg.MaxPrimBy([]Src{}, fun))

	gtest.Eq(gg.MaxPrimBy([]Src{{-10}}, fun), -10)
	gtest.Eq(gg.MaxPrimBy([]Src{{0}}, fun), 0)
	gtest.Eq(gg.MaxPrimBy([]Src{{10}}, fun), 10)

	gtest.Eq(gg.MaxPrimBy([]Src{{0}, {-10}}, fun), 0)
	gtest.Eq(gg.MaxPrimBy([]Src{{-10}, {0}}, fun), 0)
	gtest.Eq(gg.MaxPrimBy([]Src{{0}, {10}}, fun), 10)
	gtest.Eq(gg.MaxPrimBy([]Src{{10}, {0}}, fun), 10)
	gtest.Eq(gg.MaxPrimBy([]Src{{-10}, {10}}, fun), 10)
	gtest.Eq(gg.MaxPrimBy([]Src{{10}, {-10}}, fun), 10)

	gtest.Eq(gg.MaxPrimBy([]Src{{-10}, {0}, {-20}}, fun), 0)
	gtest.Eq(gg.MaxPrimBy([]Src{{0}, {10}, {-10}}, fun), 10)
	gtest.Eq(gg.MaxPrimBy([]Src{{10}, {20}, {0}}, fun), 20)
}

func TestMinBy(t *testing.T) {
	defer gtest.Catch(t)

	type Src = int
	type Out = Comparer[Src]
	var fun = ComparerOf[Src]

	gtest.Zero(gg.MinBy[Src, Out](nil, nil))
	gtest.Zero(gg.MinBy[Src, Out]([]Src{}, nil))
	gtest.Zero(gg.MinBy[Src, Out]([]Src{10}, nil))
	gtest.Zero(gg.MinBy[Src, Out]([]Src{10, 20}, nil))

	gtest.Zero(gg.MinBy[Src, Out](nil, fun))
	gtest.Zero(gg.MinBy([]Src{}, fun))

	gtest.Eq(gg.MinBy([]Src{-10}, fun), Out{-10})
	gtest.Eq(gg.MinBy([]Src{0}, fun), Out{0})
	gtest.Eq(gg.MinBy([]Src{10}, fun), Out{10})

	gtest.Eq(gg.MinBy([]Src{0, -10}, fun), Out{-10})
	gtest.Eq(gg.MinBy([]Src{-10, 0}, fun), Out{-10})
	gtest.Eq(gg.MinBy([]Src{0, 10}, fun), Out{0})
	gtest.Eq(gg.MinBy([]Src{10, 0}, fun), Out{0})
	gtest.Eq(gg.MinBy([]Src{-10, 10}, fun), Out{-10})
	gtest.Eq(gg.MinBy([]Src{10, -10}, fun), Out{-10})

	gtest.Eq(gg.MinBy([]Src{-10, -20, 0}, fun), Out{-20})
	gtest.Eq(gg.MinBy([]Src{0, -10, 10}, fun), Out{-10})
	gtest.Eq(gg.MinBy([]Src{10, 0, 20}, fun), Out{0})
}

func TestMaxBy(t *testing.T) {
	defer gtest.Catch(t)

	type Src = int
	type Out = Comparer[Src]
	var fun = ComparerOf[Src]

	gtest.Zero(gg.MaxBy[Src, Out](nil, nil))
	gtest.Zero(gg.MaxBy[Src, Out]([]Src{}, nil))
	gtest.Zero(gg.MaxBy[Src, Out]([]Src{10}, nil))
	gtest.Zero(gg.MaxBy[Src, Out]([]Src{10, 20}, nil))

	gtest.Zero(gg.MaxBy[Src, Out](nil, fun))
	gtest.Zero(gg.MaxBy([]Src{}, fun))

	gtest.Eq(gg.MaxBy([]Src{-10}, fun), Out{-10})
	gtest.Eq(gg.MaxBy([]Src{0}, fun), Out{0})
	gtest.Eq(gg.MaxBy([]Src{10}, fun), Out{10})

	gtest.Eq(gg.MaxBy([]Src{0, -10}, fun), Out{0})
	gtest.Eq(gg.MaxBy([]Src{-10, 0}, fun), Out{0})
	gtest.Eq(gg.MaxBy([]Src{0, 10}, fun), Out{10})
	gtest.Eq(gg.MaxBy([]Src{10, 0}, fun), Out{10})
	gtest.Eq(gg.MaxBy([]Src{-10, 10}, fun), Out{10})
	gtest.Eq(gg.MaxBy([]Src{10, -10}, fun), Out{10})

	gtest.Eq(gg.MaxBy([]Src{-10, 0, -20}, fun), Out{0})
	gtest.Eq(gg.MaxBy([]Src{0, 10, -10}, fun), Out{10})
	gtest.Eq(gg.MaxBy([]Src{10, 20, 0}, fun), Out{20})
}

func TestSum(t *testing.T) {
	defer gtest.Catch(t)

	type Out = int
	type Src = Comparer[Out]
	var fun = Src.Get

	gtest.Zero(gg.Sum[Src, Out](nil, nil))
	gtest.Zero(gg.Sum[Src, Out]([]Src{}, nil))
	gtest.Zero(gg.Sum[Src, Out]([]Src{{10}}, nil))
	gtest.Zero(gg.Sum[Src, Out]([]Src{{10}, {20}}, nil))

	gtest.Zero(gg.Sum[Src, Out](nil, fun))
	gtest.Zero(gg.Sum([]Src{}, fun))

	gtest.Eq(gg.Sum([]Src{{-10}}, fun), -10)
	gtest.Eq(gg.Sum([]Src{{0}}, fun), 0)
	gtest.Eq(gg.Sum([]Src{{10}}, fun), 10)

	gtest.Eq(gg.Sum([]Src{{0}, {-10}}, fun), -10)
	gtest.Eq(gg.Sum([]Src{{-10}, {0}}, fun), -10)
	gtest.Eq(gg.Sum([]Src{{0}, {10}}, fun), 10)
	gtest.Eq(gg.Sum([]Src{{10}, {0}}, fun), 10)
	gtest.Eq(gg.Sum([]Src{{-10}, {20}}, fun), 10)
	gtest.Eq(gg.Sum([]Src{{20}, {-10}}, fun), 10)
	gtest.Eq(gg.Sum([]Src{{10}, {20}}, fun), 30)
	gtest.Eq(gg.Sum([]Src{{20}, {10}}, fun), 30)

	gtest.Eq(gg.Sum([]Src{{-10}, {-20}, {0}}, fun), -30)
	gtest.Eq(gg.Sum([]Src{{0}, {-10}, {20}}, fun), 10)
	gtest.Eq(gg.Sum([]Src{{10}, {0}, {20}}, fun), 30)
}

func TestCounts(t *testing.T) {
	defer gtest.Catch(t)

	type Key = string
	type Val = byte
	type Src = []Val
	type Tar = map[Key]int

	gtest.Zero(gg.Counts[Src, Key](Src(nil), nil))
	gtest.Zero(gg.Counts[Src, Key](Src{10}, nil))
	gtest.Zero(gg.Counts[Src, Key](Src{10, 20}, nil))

	alwaysSame := func(Val) Key { return `key` }

	gtest.Equal(gg.Counts(Src{}, alwaysSame), Tar{})
	gtest.Equal(gg.Counts(Src{10}, alwaysSame), Tar{`key`: 1})
	gtest.Equal(gg.Counts(Src{10, 20}, alwaysSame), Tar{`key`: 2})
	gtest.Equal(gg.Counts(Src{10, 20, 30}, alwaysSame), Tar{`key`: 3})

	gtest.Equal(gg.Counts(Src{}, gg.String[Val]), Tar{})
	gtest.Equal(gg.Counts(Src{10}, gg.String[Val]), Tar{`10`: 1})
	gtest.Equal(gg.Counts(Src{10, 20}, gg.String[Val]), Tar{`10`: 1, `20`: 1})
	gtest.Equal(gg.Counts(Src{10, 20, 30}, gg.String[Val]), Tar{`10`: 1, `20`: 1, `30`: 1})
}
