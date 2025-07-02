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

	gtest.Zero(gg.JsonDecodeTo[Type](`null`))

	gtest.Equal(
		gg.JsonDecodeTo[Type](`123`),
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
	gtest.NoErr(tar.Scan(float64(9.999999682655225e-18)))
	gtest.Eq(tar.Val, 9.999999682655225e-18)

	tar.Clear()
	gtest.Zero(tar)
	gtest.NoErr(tar.Scan(`9.999999682655225e-18`))
	gtest.Eq(tar.Val, 9.999999682655225e-18)
}

func TestOpt_Value(t *testing.T) {
	defer gtest.Catch(t)

	type Type = gg.Opt[float64]
	var tar Type
	tar.Val = 123.456

	gtest.Zero(gg.Try1(tar.Value()))

	tar.Ok = true
	gtest.Eq(gg.Try1(tar.Value()), 123.456)
}
