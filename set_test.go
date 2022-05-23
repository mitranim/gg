package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestSetOf(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.SetOf[int](), IntSet{})
	gtest.Equal(gg.SetOf(10), IntSet{10: void})
	gtest.Equal(gg.SetOf(10, 20), IntSet{10: void, 20: void})
	gtest.Equal(gg.SetOf(10, 20, 30), IntSet{10: void, 20: void, 30: void})
	gtest.Equal(gg.SetOf(10, 20, 30, 10, 20), IntSet{10: void, 20: void, 30: void})
}

func BenchmarkSetOf_empty(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.SetOf[int]())
	}
}

func BenchmarkSetOf_non_empty(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.SetOf(10, 20, 30, 40, 50, 60, 70, 80, 90))
	}
}

func TestSet(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`SetOf`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(gg.SetOf[int](), IntSet{})
		gtest.Equal(gg.SetOf[int](10), IntSet{10: void})
		gtest.Equal(gg.SetOf[int](10, 20), IntSet{10: void, 20: void})
	})

	t.Run(`Add`, func(t *testing.T) {
		defer gtest.Catch(t)

		set := IntSet{}
		gtest.Equal(set, IntSet{})

		set.Add(10)
		gtest.Equal(set, IntSet{10: void})

		set.Add(20, 30)
		gtest.Equal(set, IntSet{10: void, 20: void, 30: void})
	})

	t.Run(`Clear`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(IntSet{10: void, 20: void}.Clear(), IntSet{})
	})

	t.Run(`Reset`, func(t *testing.T) {
		defer gtest.Catch(t)

		set := IntSet{}
		gtest.Equal(set, IntSet{})

		set.Add(10)
		gtest.Equal(set, IntSet{10: void})

		set.Reset(20, 30)
		gtest.Equal(set, IntSet{20: void, 30: void})
	})

	// TODO test multiple values (issue: ordering).
	t.Run(`Slice`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Zero(IntSet(nil).Slice())
		gtest.Equal(IntSet{}.Slice(), []int{})
		gtest.Equal(IntSet{10: void}.Slice(), []int{10})
	})

	// TODO test multiple values (issue: ordering).
	t.Run(`Filter`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Zero(IntSet(nil).Filter(gg.IsPos[int]))
		gtest.Zero(IntSet{}.Filter(gg.IsPos[int]))
		gtest.Zero(IntSet{-10: void}.Filter(gg.IsPos[int]))
		gtest.Equal(IntSet{10: void}.Filter(gg.IsPos[int]), []int{10})
	})

	// TODO test multiple values (issue: ordering).
	t.Run(`MarshalJSON`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(set IntSet, exp string) {
			gtest.Equal(gg.JsonString(set), exp)
		}

		test(IntSet(nil), `null`)
		test(IntSet{}, `[]`)
		test(IntSet{}.Add(10), `[10]`)
	})

	// TODO test multiple values (issue: ordering).
	t.Run(`UnmarshalJSON`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src string, exp IntSet) {
			gtest.Equal(gg.JsonParseTo[IntSet](src), exp)
		}

		test(`[]`, IntSet{})
		test(`[10]`, IntSet{}.Add(10))
	})

	// TODO test multiple values (issue: ordering).
	t.Run(`GoString`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.GoString(IntSet(nil)), `gg.Set[int](nil)`)
		gtest.Eq(gg.GoString(IntSet{}), `gg.Set[int]{}`)
		gtest.Eq(gg.GoString(IntSet{}.Add(10)), `gg.Set[int]{}.Add(10)`)
	})
}

func Benchmark_Set_GoString(b *testing.B) {
	val := gg.SetOf(10, 20, 30)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(val.GoString())
	}
}
