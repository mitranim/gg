package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

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
