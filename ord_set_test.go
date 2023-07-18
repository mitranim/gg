package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestOrdSetOf(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.OrdSetOf[int]())

	gtest.Equal(
		gg.OrdSetOf(10, 30, 20, 20, 10),
		gg.OrdSet[int]{
			Slice: []int{10, 30, 20},
			Index: gg.SetOf(10, 20, 30),
		},
	)
}

func TestOrdSetFrom(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.OrdSetFrom[[]int, int]())

	gtest.Equal(
		gg.OrdSetFrom([]int{10, 30}, []int{20, 20, 10}),
		gg.OrdSet[int]{
			Slice: []int{10, 30, 20},
			Index: gg.SetOf(10, 20, 30),
		},
	)
}

func TestOrdSet(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`Has`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.False(gg.OrdSetOf[int]().Has(0))
		gtest.False(gg.OrdSetOf[int]().Has(10))
		gtest.False(gg.OrdSetOf[int](10).Has(20))
		gtest.True(gg.OrdSetOf[int](10).Has(10))
		gtest.True(gg.OrdSetOf[int](10, 20).Has(10))
		gtest.True(gg.OrdSetOf[int](10, 20).Has(20))
	})

	t.Run(`Add`, func(t *testing.T) {
		defer gtest.Catch(t)

		var tar gg.OrdSet[int]

		tar.Add(10)
		gtest.Equal(tar, gg.OrdSetOf(10))

		tar.Add(20)
		gtest.Equal(tar, gg.OrdSetOf(10, 20))

		tar.Add(10, 10)
		gtest.Equal(tar, gg.OrdSetOf(10, 20))

		tar.Add(20, 20)
		gtest.Equal(tar, gg.OrdSetOf(10, 20))

		tar.Add(40, 30)
		gtest.Equal(tar, gg.OrdSetOf(10, 20, 40, 30))
	})

	t.Run(`Clear`, func(t *testing.T) {
		defer gtest.Catch(t)

		tar := gg.OrdSetOf(10, 20, 30)
		gtest.NotZero(tar)

		tar.Clear()
		gtest.Zero(tar)
	})

	t.Run(`MarshalJSON`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.JsonString(gg.OrdSet[int]{}), `null`)

		gtest.Eq(
			gg.JsonString(gg.OrdSetOf(20, 10, 30)),
			`[20,10,30]`,
		)
	})

	t.Run(`UnmarshalJSON`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(
			gg.JsonDecodeTo[gg.OrdSet[int]](`[20, 10, 30]`),
			gg.OrdSetOf(20, 10, 30),
		)
	})
}
