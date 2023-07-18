package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestZop_MarshalJSON(t *testing.T) {
	defer gtest.Catch(t)

	type Type = gg.Zop[int]

	gtest.Eq(gg.JsonString(gg.Zero[Type]()), `null`)
	gtest.Eq(gg.JsonString(Type{123}), `123`)
	gtest.Eq(gg.JsonString(gg.ZopVal(123)), `123`)
}

func TestZop_UnmarshalJSON(t *testing.T) {
	defer gtest.Catch(t)

	type Type = gg.Zop[int]

	gtest.Zero(gg.JsonDecodeTo[Type](`null`))

	gtest.Equal(
		gg.JsonDecodeTo[Type](`123`),
		gg.ZopVal(123),
	)
}
