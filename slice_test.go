package gg_test

import (
	"strconv"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func BenchmarkSliceDat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.SliceDat([]byte(`hello world`)))
	}
}

func TestLens(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Lens[[]int](), 0)
	gtest.Eq(gg.Lens([]int{}, []int{10}), 1)
	gtest.Eq(gg.Lens([]int{}, []int{10}, []int{20, 30}), 3)
}

func BenchmarkLens(b *testing.B) {
	val := [][]int{{}, {10}, {20, 30}, {40, 50, 60}}

	for i := 0; i < b.N; i++ {
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
			t.Helper()

			tar := gg.GrowLen(cur, size)
			gtest.Equal(src, expSrc)
			gtest.Equal(tar, expTar)
			gtest.Eq(cap(tar), cap(src))
			gtest.Eq(gg.SliceDat(src), gg.SliceDat(tar))
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
		gtest.NotEq(gg.SliceDat(src), gg.SliceDat(tar))
	})
}

func BenchmarkGrowLen(b *testing.B) {
	buf := make([]byte, 0, 1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
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

		gtest.Eq(gg.SliceDat(src), gg.SliceDat(tar))
		gtest.Eq(len(tar), len(src))
		gtest.Eq(cap(tar), cap(src))
	})
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

func TestMap(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.Map([]int(nil), strconv.Itoa), []string(nil))
	gtest.Equal(gg.Map([]int{}, strconv.Itoa), []string{})
	gtest.Equal(gg.Map([]int{10}, strconv.Itoa), []string{`10`})
	gtest.Equal(gg.Map([]int{10, 20}, strconv.Itoa), []string{`10`, `20`})
	gtest.Equal(gg.Map([]int{10, 20, 30}, strconv.Itoa), []string{`10`, `20`, `30`})
}

func BenchmarkMap(b *testing.B) {
	val := gg.Span(32)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Map(val, gg.Inc[int]))
	}
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

func TestPrimSorted(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(
		gg.SortedPrim(gg.SliceOf(20, 30, 10, 40)),
		gg.SliceOf(10, 20, 30, 40),
	)
}

func BenchmarkPrimSorted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.SortedPrim(gg.SliceOf(20, 30, 10, 40)))
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
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Reversed([]int{20, 30, 10, 40}))
	}
}

func TestSubtract(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.Subtract([]int(nil)), []int(nil))
	gtest.Equal(gg.Subtract([]int{}), []int(nil))
	gtest.Equal(gg.Subtract([]int{10}), []int{10})
	gtest.Equal(gg.Subtract([]int{10}, []int{20}), []int{10})
	gtest.Equal(gg.Subtract([]int{10}, []int{10, 20}), []int(nil))
	gtest.Equal(gg.Subtract([]int{10, 20, 30}, []int{10, 20}), []int{30})
}

func BenchmarkSubtract(b *testing.B) {
	base := []int{10, 20, 30, 40, 50, 60}
	sub := [][]int{{10, 20}, {50}}

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Subtract(base, sub...))
	}
}

func BenchmarkTakeWhile(b *testing.B) {
	val := []int{-30, -20, -10, 0, 10, 20, 30}

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.TakeWhile(val, gg.IsNeg[int]))
	}
}

func TestPlus(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`Num`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.Plus[int](), 0)
		gtest.Eq(gg.Plus(0), 0)
		gtest.Eq(gg.Plus(10), 10)
		gtest.Eq(gg.Plus(10, 20), 30)
		gtest.Eq(gg.Plus(-10, 0), -10)
		gtest.Eq(gg.Plus(-10, 0, 10), 0)
		gtest.Eq(gg.Plus(-10, 0, 10, 20), 20)
		gtest.Eq(gg.Plus(-10, 0, 10, 20, 30), 50)
	})

	t.Run(`string`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.Plus[string](), ``)
		gtest.Eq(gg.Plus(``), ``)
		gtest.Eq(gg.Plus(`one`), `one`)
		gtest.Eq(gg.Plus(`one`, ``), `one`)
		gtest.Eq(gg.Plus(``, `two`), `two`)
		gtest.Eq(gg.Plus(`one`, `two`), `onetwo`)
		gtest.Eq(gg.Plus(`one`, `two`, `three`), `onetwothree`)
	})
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Plus(10, 20, 30, 40, 50, 60, 70, 80, 90))
	}
}

func BenchmarkMinPrim(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.MinPrim(i-1, i, i+1))
	}
}

func BenchmarkMaxPrim(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.MaxPrim(i-1, i, i+1))
	}
}
