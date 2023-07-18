package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

const testPanicStrEmptyKey = `unexpected empty key gg_test.SomeKey in gg_test.SomeModel`

type AnyColl interface{ SomeColl | SomeLazyColl }

/*
This is a workaround. Ideally we'd just use `AnyColl`, but Go wouldn't allow us
to access ANY of its methods or properties.
*/
type CollPtr[Coll AnyColl] interface {
	*Coll
	Len() int
	IsEmpty() bool
	IsNotEmpty() bool
	Has(SomeKey) bool
	Get(SomeKey) SomeModel
	GetReq(SomeKey) SomeModel
	Got(SomeKey) (SomeModel, bool)
	Ptr(SomeKey) *SomeModel
	PtrReq(SomeKey) *SomeModel
	Add(...SomeModel) *Coll
	Reset(...SomeModel) *Coll
	Clear() *Coll
	Reindex() *Coll
	Swap(int, int)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

func TestColl(t *testing.T) {
	defer gtest.Catch(t)
	testColl[*SomeColl]()
}

func testColl[Ptr CollPtr[Coll], Coll AnyColl]() {
	testCollReset[Ptr]()
	testCollClear[Ptr]()
	testCollHas[Ptr]()
	testCollGet[Ptr]()
	testCollGetReq[Ptr]()
	testCollGot[Ptr]()
	testCollPtr[Ptr]()
	testCollPtrReq[Ptr]()
	testColl_Len_IsEmpty_IsNotEmpty[Ptr]()
	testCollAdd[Ptr]()
	testCollReindex[Ptr]()
	testCollSwap[Ptr]()
	testCollMarshalJSON[Ptr]()
	testCollUnmarshalJSON[Ptr]()
}

func testCollReset[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	testEmpty := func() {
		gtest.Eq(ptr, ptr.Reset())
		gtest.Zero(tar)
	}

	testEmpty()

	test := func(src ...SomeModel) {
		gtest.Eq(ptr, ptr.Reset(gg.Clone(src)...))
		testCollSlice(tar, src)
	}

	test()
	test(SomeModel{10, `one`})
	test(SomeModel{10, `one`}, SomeModel{20, `two`})

	// Artifact of the implementation. Divergence from `Coll.Add` which would
	// detect and deduplicate entries with the same primary key, although
	// `LazyColl.Add` would not. This behavior exists in both `Coll` and
	// `LazyColl` because `.Reset` is meant to store the slice as-is.
	test(SomeModel{10, `one`}, SomeModel{10, `two`}, SomeModel{30, `three`})

	if gg.EqType[Coll, SomeColl]() {
		// We're forced to exclude redundant elements from the index.
		testCollIndex(tar, map[SomeKey]int{10: 1, 30: 2})
	} else if gg.EqType[Coll, SomeLazyColl]() {
		testCollIndex(tar, nil)
	}

	testEmpty()
}

func testCollEqual[Coll AnyColl](src Coll, exp SomeColl) {
	gtest.Equal(gg.Cast[SomeColl](src), exp)
}

/*
Workaround because Go doesn't let us access properties on values of type
parameters constrained by `IColl`.
*/
func testCollSlice[Coll AnyColl](src Coll, exp []SomeModel) {
	gtest.Equal(getCollSlice(src), exp)
}

func getCollSlice[Coll AnyColl](src Coll) []SomeModel {
	return gg.Cast[SomeColl](src).Slice
}

/*
Workaround because Go doesn't let us access properties on values of type
parameters constrained by `IColl`.
*/
func testCollIndex[Coll AnyColl](src Coll, exp map[SomeKey]int) {
	gtest.Equal(getCollIndex(src), exp)
}

func getCollIndex[Coll SomeColl | SomeLazyColl](src Coll) map[SomeKey]int {
	return gg.Cast[SomeColl](src).Index
}

func testCollClear[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	ptr.Reset(SomeModel{10, `one`}, SomeModel{20, `two`})
	gtest.NotZero(tar)

	gtest.Eq(ptr, ptr.Clear())
	gtest.Zero(tar)
}

func testCollHas[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	gtest.False(ptr.Has(0))
	gtest.False(ptr.Has(10))
	gtest.False(ptr.Has(20))
	gtest.False(ptr.Has(30))

	ptr.Reset(SomeModel{10, `one`}, SomeModel{20, `two`})

	gtest.False(ptr.Has(0))
	gtest.True(ptr.Has(10))
	gtest.True(ptr.Has(20))
	gtest.False(ptr.Has(30))

	testIdempotentIndexing(func(tar *SomeLazyColl) { tar.Has(0) })
}

func testIdempotentIndexing(fun func(*SomeLazyColl)) {
	var tar SomeLazyColl
	tar.Reset(SomeModel{10, `one`}, SomeModel{20, `two`})
	testCollIndex(tar, nil)

	fun(&tar)

	testCollIndex(tar, map[SomeKey]int{10: 0, 20: 1})
	index := getCollIndex(tar)

	fun(&tar)

	testCollIndex(tar, map[SomeKey]int{10: 0, 20: 1})

	// Technically, this doesn't guarantee that the collection doesn't rebuild the
	// index by deleting and re-adding entries. We rely on this check because we
	// know that we only clear the index by setting it to nil, and finding a
	// different reference would indicate that it's been rebuilt.
	gtest.Is(getCollIndex(tar), index, `must preserve existing index as-is`)
}

func testCollGet[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	gtest.Zero(ptr.Get(0))
	gtest.Zero(ptr.Get(10))
	gtest.Zero(ptr.Get(20))
	gtest.Zero(ptr.Get(30))

	ptr.Reset(SomeModel{10, `one`}, SomeModel{20, `two`})

	gtest.Zero(ptr.Get(0))
	gtest.Eq(ptr.Get(10), SomeModel{10, `one`})
	gtest.Eq(ptr.Get(20), SomeModel{20, `two`})
	gtest.Zero(ptr.Get(30))

	testIdempotentIndexing(func(tar *SomeLazyColl) { tar.Get(0) })
}

func testCollGetReq[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 0`, func() { ptr.GetReq(0) })
	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 10`, func() { ptr.GetReq(10) })
	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 20`, func() { ptr.GetReq(20) })
	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 30`, func() { ptr.GetReq(30) })

	ptr.Reset(SomeModel{10, `one`}, SomeModel{20, `two`})

	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 0`, func() { ptr.GetReq(0) })
	gtest.Eq(ptr.GetReq(10), SomeModel{10, `one`})
	gtest.Eq(ptr.GetReq(20), SomeModel{20, `two`})
	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 30`, func() { ptr.GetReq(30) })

	testIdempotentIndexing(func(tar *SomeLazyColl) { tar.GetReq(10) })
}

func testCollGot[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	test := func(key SomeKey, expVal SomeModel, expOk bool) {
		val, ok := ptr.Got(key)
		gtest.Eq(expVal, val)
		gtest.Eq(expOk, ok)
	}

	test(0, SomeModel{}, false)
	test(10, SomeModel{}, false)
	test(20, SomeModel{}, false)
	test(30, SomeModel{}, false)

	ptr.Reset(SomeModel{10, `one`}, SomeModel{20, `two`})

	test(0, SomeModel{}, false)
	test(10, SomeModel{10, `one`}, true)
	test(20, SomeModel{20, `two`}, true)
	test(30, SomeModel{}, false)

	testIdempotentIndexing(func(tar *SomeLazyColl) { tar.Got(0) })
}

func testCollPtr[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	gtest.Zero(ptr.Ptr(0))
	gtest.Zero(ptr.Ptr(10))
	gtest.Zero(ptr.Ptr(20))
	gtest.Zero(ptr.Ptr(30))

	ptr.Reset(SomeModel{10, `one`}, SomeModel{20, `two`})

	gtest.Zero(ptr.Ptr(0))

	gtest.Equal(ptr.Ptr(10), &SomeModel{10, `one`})
	gtest.Eq(ptr.Ptr(10), ptr.Ptr(10))

	gtest.Equal(ptr.Ptr(20), &SomeModel{20, `two`})
	gtest.Eq(ptr.Ptr(20), ptr.Ptr(20))

	gtest.Zero(ptr.Ptr(30))

	ptr.Ptr(10).Name = `three`
	gtest.Equal(ptr.Ptr(10), &SomeModel{10, `three`})

	ptr.Ptr(10).Name = `four`
	gtest.Equal(ptr.Ptr(10), &SomeModel{10, `four`})

	testIdempotentIndexing(func(tar *SomeLazyColl) { tar.Ptr(0) })
}

func testCollPtrReq[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 0`, func() { ptr.PtrReq(0) })
	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 10`, func() { ptr.PtrReq(10) })
	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 20`, func() { ptr.PtrReq(20) })
	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 30`, func() { ptr.PtrReq(30) })

	ptr.Reset(SomeModel{10, `one`}, SomeModel{20, `two`})

	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 0`, func() { ptr.PtrReq(0) })
	gtest.Equal(ptr.PtrReq(10), &SomeModel{10, `one`})
	gtest.Equal(ptr.PtrReq(20), &SomeModel{20, `two`})
	gtest.PanicStr(`missing value of type gg_test.SomeModel for key 30`, func() { ptr.PtrReq(30) })

	ptr.PtrReq(10).Name = `three`
	gtest.Equal(ptr.PtrReq(10), &SomeModel{10, `three`})

	ptr.PtrReq(10).Name = `four`
	gtest.Equal(ptr.PtrReq(10), &SomeModel{10, `four`})

	testIdempotentIndexing(func(tar *SomeLazyColl) { tar.PtrReq(10) })
}

func testColl_Len_IsEmpty_IsNotEmpty[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	gtest.Eq(ptr.Len(), 0)
	gtest.True(ptr.IsEmpty())
	gtest.False(ptr.IsNotEmpty())

	ptr.Reset(SomeModel{10, `one`})
	gtest.Eq(ptr.Len(), 1)
	gtest.False(ptr.IsEmpty())
	gtest.True(ptr.IsNotEmpty())

	ptr.Reset(SomeModel{10, `one`}, SomeModel{20, `two`})
	gtest.Eq(ptr.Len(), 2)
	gtest.False(ptr.IsEmpty())
	gtest.True(ptr.IsNotEmpty())

	ptr.Clear()
	gtest.Eq(ptr.Len(), 0)
	gtest.True(ptr.IsEmpty())
	gtest.False(ptr.IsNotEmpty())
}

func testCollAdd[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	gtest.PanicStr(testPanicStrEmptyKey, func() { ptr.Add(SomeModel{}) })
	testCollSlice(tar, nil)

	ptr.Add(SomeModel{10, `one`})
	testCollSlice(tar, []SomeModel{{10, `one`}})

	ptr.Add(SomeModel{20, `two`})
	testCollSlice(tar, []SomeModel{{10, `one`}, {20, `two`}})

	if gg.EqType[Coll, SomeColl]() {
		testCollIndex(tar, map[SomeKey]int{10: 0, 20: 1})
	} else if gg.EqType[Coll, SomeLazyColl]() {
		testCollIndex(tar, nil)
	}

	// Known divergence between `Coll` and `LazyColl`.
	{
		var tar SomeColl
		tar.Add(SomeModel{10, `one`})
		tar.Add(SomeModel{20, `two`})
		tar.Add(SomeModel{10, `three`})
		testCollIndex(tar, map[SomeKey]int{10: 0, 20: 1})
		testCollSlice(tar, []SomeModel{{10, `three`}, {20, `two`}})
	}
	// This happens when `LazyColl` is not indexed.
	{
		var tar SomeLazyColl
		tar.Add(SomeModel{10, `one`})
		tar.Add(SomeModel{20, `two`})
		tar.Add(SomeModel{10, `three`})
		testCollIndex(tar, nil)
		testCollSlice(tar, []SomeModel{{10, `one`}, {20, `two`}, {10, `three`}})
	}
	// Reindexing (for any reason) avoids the problem.
	{
		var tar SomeLazyColl
		tar.Add(SomeModel{10, `one`})
		tar.Add(SomeModel{20, `two`})
		testCollIndex(tar, nil)
		tar.Reindex()
		testCollIndex(tar, map[SomeKey]int{10: 0, 20: 1})
		tar.Add(SomeModel{10, `three`})
		testCollSlice(tar, []SomeModel{{10, `three`}, {20, `two`}})
	}
}

func testCollReindex[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	testCollIndex(tar, nil)
	ptr.Reset(SomeModel{10, `one`})

	ptr.Reindex()
	testCollIndex(tar, map[SomeKey]int{10: 0})
}

func testCollSwap[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	ptr.Reset(SomeModel{10, `one`}, SomeModel{20, `two`}, SomeModel{30, `three`})
	prev := gg.CloneDeep(tar)

	same := func(ind int) {
		ptr.Swap(ind, ind)
		gtest.Equal(prev, tar)
	}
	same(0)
	same(1)
	same(2)
	testCollSlice(tar, []SomeModel{{10, `one`}, {20, `two`}, {30, `three`}})

	testGet := func() {
		gtest.Equal(ptr.GetReq(10), SomeModel{10, `one`})
		gtest.Equal(ptr.GetReq(20), SomeModel{20, `two`})
		gtest.Equal(ptr.GetReq(30), SomeModel{30, `three`})
	}

	ptr.Swap(0, 1)
	testGet()

	if gg.EqType[Coll, Coll]() {
		testCollIndex(tar, map[SomeKey]int{10: 1, 20: 0, 30: 2})
	} else if gg.EqType[Coll, SomeLazyColl]() {
		testCollIndex(tar, nil)
	}

	ptr.Swap(2, 0)
	testGet()

	if gg.EqType[Coll, Coll]() {
		testCollIndex(tar, map[SomeKey]int{10: 1, 20: 2, 30: 0})
	} else if gg.EqType[Coll, SomeLazyColl]() {
		testCollIndex(tar, nil)
	}
}

func testCollMarshalJSON[Ptr CollPtr[Coll], Coll AnyColl]() {
	var tar Coll
	ptr := Ptr(&tar)

	gtest.Eq(gg.JsonString(tar), `null`)

	ptr.Reset(SomeModel{10, `one`}, SomeModel{20, `two`})

	gtest.Eq(
		gg.JsonString(tar),
		`[{"id":10,"name":"one"},{"id":20,"name":"two"}]`,
	)
}

func testCollUnmarshalJSON[Ptr CollPtr[Coll], Coll AnyColl]() {
	tar := gg.JsonDecodeTo[Coll](`[
		{"id": 10, "name": "one"},
		{"id": 20, "name": "two"}
	]`)

	testCollSlice(tar, []SomeModel{{10, `one`}, {20, `two`}})

	if gg.EqType[Coll, SomeColl]() {
		testCollIndex(tar, map[SomeKey]int{10: 0, 20: 1})
	} else if gg.EqType[Coll, SomeLazyColl]() {
		testCollIndex(tar, nil)
	}
}

func TestCollOf(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.CollOf[SomeKey, SomeModel]())

	testCollMake(func(src ...SomeModel) SomeColl {
		return gg.CollOf[SomeKey, SomeModel](src...)
	})
}

func testCollMake[Coll AnyColl](fun func(...SomeModel) Coll) {
	test := func(slice []SomeModel, index map[SomeKey]int) {
		tar := fun(slice...)
		testCollEqual(tar, SomeColl{Slice: slice, Index: index})
		gtest.Is(getCollSlice(tar), slice)
	}

	test(
		[]SomeModel{SomeModel{10, `one`}},
		map[SomeKey]int{10: 0},
	)

	test(
		[]SomeModel{
			SomeModel{10, `one`},
			SomeModel{20, `two`},
		},
		map[SomeKey]int{
			10: 0,
			20: 1,
		},
	)
}

func TestCollFrom(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.CollFrom[SomeKey, SomeModel, []SomeModel]())

	testCollMake(func(src ...SomeModel) SomeColl {
		return gg.CollFrom[SomeKey, SomeModel](src)
	})

	gtest.Equal(
		gg.CollFrom[SomeKey, SomeModel](
			[]SomeModel{
				SomeModel{10, `one`},
				SomeModel{20, `two`},
			},
			[]SomeModel{
				SomeModel{30, `three`},
				SomeModel{40, `four`},
			},
		),
		SomeColl{
			Slice: []SomeModel{
				SomeModel{10, `one`},
				SomeModel{20, `two`},
				SomeModel{30, `three`},
				SomeModel{40, `four`},
			},
			Index: map[SomeKey]int{
				10: 0,
				20: 1,
				30: 2,
				40: 3,
			},
		},
	)
}
