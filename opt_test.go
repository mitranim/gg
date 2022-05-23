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

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(val.String())
	}
}

func TestOpt_Scan(t *testing.T) {
	defer gtest.Catch(t)

	type Type = gg.Opt[float64]

	var tar Type
	gtest.NoError(tar.Scan(float64(9.999999682655225e-18)))
	gtest.Eq(tar.Val, 9.999999682655225e-18)

	tar.Clear()
	gtest.Zero(tar)
	gtest.NoError(tar.Scan(`9.999999682655225e-18`))
	gtest.Eq(tar.Val, 9.999999682655225e-18)
}
