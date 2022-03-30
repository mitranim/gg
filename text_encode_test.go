package gg_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

// TODO test invalid.
func TestString(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.String(any(nil)), ``)

	gtest.Eq(gg.String(false), `false`)
	gtest.Eq(gg.String(true), `true`)

	gtest.Eq(gg.String(0), `0`)
	gtest.Eq(gg.String(10), `10`)
	gtest.Eq(gg.String(-10), `-10`)

	gtest.Eq(gg.String(0.0), `0`)
	gtest.Eq(gg.String(10.23), `10.23`)
	gtest.Eq(gg.String(-10.23), `-10.23`)

	gtest.Eq(gg.String(``), ``)
	gtest.Eq(gg.String(`str`), `str`)

	gtest.Eq(gg.String([]byte(nil)), ``)
	gtest.Eq(gg.String([]byte{}), ``)
	gtest.Eq(gg.String([]byte(`str`)), `str`)

	gtest.Eq(gg.String(&url.URL{Path: `/one`}), `/one`)
}

func BenchmarkString(b *testing.B) {
	val := `str`

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.String(val))
	}
}

func TestAppend(t *testing.T) {
	defer gtest.Catch(t)

	type Bui []byte

	gtest.Equal(gg.Append(Bui(nil), any(nil)), Bui(nil))
	gtest.Equal(gg.Append(Bui(``), any(nil)), Bui(``))
	gtest.Equal(gg.Append(Bui(`str`), any(nil)), Bui(`str`))
	gtest.Equal(gg.Append(Bui(nil), 10), Bui(`10`))
	gtest.Equal(gg.Append(Bui(`str`), 10), Bui(`str10`))
}

func Benchmark_fmt_Sprint(b *testing.B) {
	var val SomeModel

	for i := 0; i < b.N; i++ {
		gg.Nop1(fmt.Sprint(val))
	}
}

func BenchmarkStringAny(b *testing.B) {
	var val SomeModel

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.StringAny(val))
	}
}
