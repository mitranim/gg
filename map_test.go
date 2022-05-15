package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestMapInit(t *testing.T) {
	defer gtest.Catch(t)

	var tar IntMap
	gtest.Equal(gg.MapInit(&tar), tar)
	gtest.NotZero(tar)

	tar[10] = 20
	gtest.Equal(gg.MapInit(&tar), tar)
	gtest.Equal(tar, IntMap{10: 20})
}

func TestMapClone(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.MapClone(IntMap(nil)), IntMap(nil))

	src := IntMap{10: 20, 30: 40}
	out := gg.MapClone(src)
	gtest.Equal(out, src)

	src[10] = 50
	gtest.Equal(src, IntMap{10: 50, 30: 40})
	gtest.Equal(out, IntMap{10: 20, 30: 40})
}

func TestMapKeys(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src IntMap, exp []int) {
		gtest.Equal(gg.SortedPrim(gg.MapKeys(src)), exp)
	}

	test(IntMap(nil), []int(nil))
	test(IntMap{}, []int{})
	test(IntMap{10: 20}, []int{10})
	test(IntMap{10: 20, 30: 40}, []int{10, 30})
}

func TestMapVals(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src IntMap, exp []int) {
		gtest.Equal(gg.SortedPrim(gg.MapVals(src)), exp)
	}

	test(IntMap(nil), []int(nil))
	test(IntMap{}, []int{})
	test(IntMap{10: 20}, []int{20})
	test(IntMap{10: 20, 30: 40}, []int{20, 40})
}

func TestMapHas(t *testing.T) {
	defer gtest.Catch(t)

	gtest.False(gg.MapHas(IntMap(nil), 10))
	gtest.False(gg.MapHas(IntMap{10: 20}, 20))
	gtest.True(gg.MapHas(IntMap{10: 20}, 10))
}

func TestMapGot(t *testing.T) {
	defer gtest.Catch(t)

	{
		val, ok := gg.MapGot(IntMap(nil), 10)
		gtest.Zero(val)
		gtest.False(ok)
	}

	{
		val, ok := gg.MapGot(IntMap{10: 20}, 20)
		gtest.Zero(val)
		gtest.False(ok)
	}

	{
		val, ok := gg.MapGot(IntMap{10: 20}, 10)
		gtest.Eq(val, 20)
		gtest.True(ok)
	}
}

func TestMapGet(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.MapGet(IntMap(nil), 10))
	gtest.Zero(gg.MapGet(IntMap{10: 20}, 20))
	gtest.Eq(gg.MapGet(IntMap{10: 20}, 10), 20)
}

func TestMapSet(t *testing.T) {
	defer gtest.Catch(t)

	tar := IntMap{}

	gg.MapSet(tar, 10, 20)
	gtest.Equal(tar, IntMap{10: 20})

	gg.MapSet(tar, 10, 30)
	gtest.Equal(tar, IntMap{10: 30})
}

func TestMapSetOpt(t *testing.T) {
	defer gtest.Catch(t)

	tar := IntMap{}

	gg.MapSetOpt(tar, 0, 20)
	gtest.Equal(tar, IntMap{})

	gg.MapSetOpt(tar, 10, 0)
	gtest.Equal(tar, IntMap{})

	gg.MapSetOpt(tar, 10, 20)
	gtest.Equal(tar, IntMap{10: 20})

	gg.MapSetOpt(tar, 10, 30)
	gtest.Equal(tar, IntMap{10: 30})
}

func TestMapClear(t *testing.T) {
	defer gtest.Catch(t)

	var tar IntMap
	gg.MapClear(tar)
	gtest.Equal(tar, IntMap(nil))

	tar = IntMap{10: 20, 30: 40}
	gg.MapClear(tar)
	gtest.Equal(tar, IntMap{})
}
