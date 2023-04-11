package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestCollOf(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.CollOf[SomeKey, SomeModel]())

	testCollOf(func(src []SomeModel) gg.Coll[SomeKey, SomeModel] {
		return gg.CollOf[SomeKey, SomeModel](src...)
	})
}

func testCollOf(fun func([]SomeModel) gg.Coll[SomeKey, SomeModel]) {
	test := func(slice []SomeModel, index map[SomeKey]int) {
		tar := gg.CollOf[SomeKey, SomeModel](slice...)

		gtest.Equal(
			tar,
			gg.Coll[SomeKey, SomeModel]{Slice: slice, Index: index},
		)

		gtest.Is(tar.Slice, slice)
	}

	test(
		[]SomeModel{SomeModel{`10`, `one`}},
		map[SomeKey]int{`10`: 0},
	)

	test(
		[]SomeModel{
			SomeModel{`10`, `one`},
			SomeModel{`20`, `two`},
		},
		map[SomeKey]int{
			`10`: 0,
			`20`: 1,
		},
	)
}

func TestCollFrom(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.CollFrom[SomeKey, SomeModel, []SomeModel]())

	testCollOf(func(src []SomeModel) gg.Coll[SomeKey, SomeModel] {
		return gg.CollFrom[SomeKey, SomeModel](src)
	})

	gtest.Equal(
		gg.CollFrom[SomeKey, SomeModel](
			[]SomeModel{
				SomeModel{`10`, `one`},
				SomeModel{`20`, `two`},
			},
			[]SomeModel{
				SomeModel{`30`, `three`},
				SomeModel{`40`, `four`},
			},
		),
		gg.Coll[SomeKey, SomeModel]{
			Slice: []SomeModel{
				SomeModel{`10`, `one`},
				SomeModel{`20`, `two`},
				SomeModel{`30`, `three`},
				SomeModel{`40`, `four`},
			},
			Index: map[SomeKey]int{
				`10`: 0,
				`20`: 1,
				`30`: 2,
				`40`: 3,
			},
		},
	)
}

func TestColl(t *testing.T) {
	defer gtest.Catch(t)

	exp := SomeColl{
		Slice: []SomeModel{
			SomeModel{`ee24ca`, `Mira`},
			SomeModel{`a19b43`, `Kara`},
		},
		Index: map[SomeKey]int{`ee24ca`: 0, `a19b43`: 1},
	}

	t.Run(`Add`, func(t *testing.T) {
		defer gtest.Catch(t)

		var coll SomeColl
		coll.Add(SomeModel{`ee24ca`, `Mira`})
		coll.Add(SomeModel{`a19b43`, `Kara`})

		gtest.Equal(coll, exp)

		// Must ensure uniqueness.
		coll.Add(SomeModel{`a19b43`, `Kara`})
		gtest.Equal(coll, exp)

		// Must replace existing entry.
		coll.Add(SomeModel{`a19b43`, `Kara_1`})
		gtest.Equal(coll, SomeColl{
			Slice: []SomeModel{
				SomeModel{`ee24ca`, `Mira`},
				SomeModel{`a19b43`, `Kara_1`},
			},
			Index: map[SomeKey]int{`ee24ca`: 0, `a19b43`: 1},
		})
	})

	t.Run(`Clear`, func(t *testing.T) {
		defer gtest.Catch(t)

		tar := gg.CloneDeep(exp)
		gtest.NotZero(tar)

		tar.Clear()
		gtest.Zero(tar)
	})

	t.Run(`MarshalJSON`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.JsonString(SomeColl{}), `null`)

		gtest.Eq(
			gg.JsonString(exp),
			`[{"id":"ee24ca","name":"Mira"},{"id":"a19b43","name":"Kara"}]`,
		)
	})

	t.Run(`UnmarshalJSON`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(gg.JsonParseTo[SomeColl](`[
			{"id": "ee24ca", "name": "Mira"},
			{"id": "a19b43", "name": "Kara"}
		]`), exp)
	})
}
