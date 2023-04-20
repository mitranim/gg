package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestLazyColl(t *testing.T) {
	defer gtest.Catch(t)
	testColl[*gg.LazyColl[SomeKey, SomeModel]]()
}

func TestLazyCollOf(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.LazyCollOf[SomeKey, SomeModel]())

	testLazyCollMake(func(src ...SomeModel) SomeLazyColl {
		return gg.LazyCollOf[SomeKey, SomeModel](src...)
	})
}

func testLazyCollMake[Coll AnyColl](fun func(...SomeModel) Coll) {
	test := func(slice []SomeModel, index map[SomeKey]int) {
		tar := fun(slice...)
		testCollEqual(tar, SomeColl{Slice: slice, Index: index})
		gtest.Is(getCollSlice(tar), slice)
	}

	test(
		[]SomeModel{SomeModel{10, `one`}},
		nil,
	)

	test(
		[]SomeModel{
			SomeModel{10, `one`},
			SomeModel{20, `two`},
		},
		nil,
	)
}

func TestLazyCollFrom(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.LazyCollFrom[SomeKey, SomeModel, []SomeModel]())

	testLazyCollMake(func(src ...SomeModel) SomeLazyColl {
		return gg.LazyCollFrom[SomeKey, SomeModel](src)
	})
}
