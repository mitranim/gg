package gg_test

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

var testTime = time.Date(1234, time.February, 23, 0, 0, 0, 0, time.UTC)

// TODO test invalid.
func TestString(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(`unable to convert value { } of type gg_test.SomeModel to type string`, func() {
		gg.String(SomeModel{})
	})

	gtest.Eq(gg.String(any(nil)), ``)

	gtest.Eq(gg.String(false), `false`)
	gtest.Eq(gg.String(true), `true`)

	gtest.Eq(gg.String(0), `0`)
	gtest.Eq(gg.String(10), `10`)
	gtest.Eq(gg.String(-10), `-10`)

	gtest.Eq(gg.String(0.0), `0`)
	gtest.Eq(gg.String(10.0), `10`)
	gtest.Eq(gg.String(-10.0), `-10`)
	gtest.Eq(gg.String(10.23), `10.23`)
	gtest.Eq(gg.String(-10.23), `-10.23`)

	gtest.Eq(gg.String(``), ``)
	gtest.Eq(gg.String(`str`), `str`)

	gtest.Eq(gg.String([]byte(nil)), ``)
	gtest.Eq(gg.String([]byte{}), ``)
	gtest.Eq(gg.String([]byte(`str`)), `str`)

	gtest.Eq(gg.String(&url.URL{Path: `/one`}), `/one`)

	t.Run(`time.Time`, func(t *testing.T) {
		defer gtest.Catch(t)

		// Unfortunate default which we choose to override/replace.
		gtest.Eq(testTime.String(), `1234-02-23 00:00:00 +0000 UTC`)

		gtest.Eq(gg.String(testTime), `1234-02-23T00:00:00Z`)
	})
}

func BenchmarkString_string(b *testing.B) {
	val := `str`

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.String(val))
	}
}

func BenchmarkString_bool(b *testing.B) {
	val := true

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.String(val))
	}
}

func BenchmarkString_int(b *testing.B) {
	val := 123

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.String(val))
	}
}

func BenchmarkString_stringer(b *testing.B) {
	val := &url.URL{Path: `/one`}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.String(val))
	}
}

func TestAppend(t *testing.T) {
	defer gtest.Catch(t)

	type Bui []byte

	gtest.Equal(gg.AppendTo(Bui(nil), any(nil)), Bui(nil))
	gtest.Equal(gg.AppendTo(Bui(``), any(nil)), Bui(``))
	gtest.Equal(gg.AppendTo(Bui(`pre_`), any(nil)), Bui(`pre_`))
	gtest.Equal(gg.AppendTo(Bui(nil), 10), Bui(`10`))
	gtest.Equal(gg.AppendTo(Bui(`pre_`), 10), Bui(`pre_10`))

	t.Run(`time.Time`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(
			gg.AppendTo(Bui(nil), testTime),
			Bui(`1234-02-23T00:00:00Z`),
		)

		gtest.Equal(
			gg.AppendTo(Bui(`pre_`), testTime),
			Bui(`pre_1234-02-23T00:00:00Z`),
		)
	})
}

func Benchmark_string_any_fmt_Sprint(b *testing.B) {
	var val SomeModel

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(fmt.Sprint(val))
	}
}

func Benchmark_string_any_StringAny(b *testing.B) {
	var val SomeModel

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.StringAny(val))
	}
}
