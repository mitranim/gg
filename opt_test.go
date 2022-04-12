package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestOpt_MarshalJSON(t *testing.T) {
	defer gtest.Catch(t)

	type Type = gg.Opt[int]

	gtest.Eq(gg.JsonString(gg.Zero[Type]()), `null`)
	gtest.Eq(gg.JsonString(Type{Val: 123}), `null`)
	gtest.Eq(gg.JsonString(gg.OptVal(123)), `123`)
}

func TestOpt_UnmarshalJSON(t *testing.T) {
	defer gtest.Catch(t)

	type Type = gg.Opt[int]

	gtest.Zero(gg.JsonParseTo[Type](`null`))

	gtest.Equal(
		gg.JsonParseTo[Type](`123`),
		gg.OptVal(123),
	)
}

func BenchmarkOpt_String(b *testing.B) {
	val := gg.OptVal(`str`)

	for i := 0; i < b.N; i++ {
		gg.Nop1(val.String())
	}
}
