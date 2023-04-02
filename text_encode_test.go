package gg_test

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
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

func TestNumBits(t *testing.T) {
	defer gtest.Catch(t)

	{
		testNumBits[uint8](0, `00000000`)
		testNumBits[uint8](1, `00000001`)
		testNumBits[uint8](2, `00000010`)

		testNumBits[uint8](3, `00000011`)
		testNumBits[uint8](4, `00000100`)
		testNumBits[uint8](5, `00000101`)

		testNumBits[uint8](7, `00000111`)
		testNumBits[uint8](8, `00001000`)
		testNumBits[uint8](9, `00001001`)

		testNumBits[uint8](15, `00001111`)
		testNumBits[uint8](16, `00010000`)
		testNumBits[uint8](17, `00010001`)

		testNumBits[uint8](31, `00011111`)
		testNumBits[uint8](32, `00100000`)
		testNumBits[uint8](33, `00100001`)

		testNumBits[uint8](63, `00111111`)
		testNumBits[uint8](64, `01000000`)
		testNumBits[uint8](65, `01000001`)

		testNumBits[uint8](127, `01111111`)
		testNumBits[uint8](128, `10000000`)
		testNumBits[uint8](129, `10000001`)

		testNumBits[uint8](255, `11111111`)
	}

	{
		testNumBits[int8](0, `00000000`)
		testNumBits[int8](1, `00000001`)
		testNumBits[int8](2, `00000010`)

		testNumBits[int8](3, `00000011`)
		testNumBits[int8](4, `00000100`)
		testNumBits[int8](5, `00000101`)

		testNumBits[int8](7, `00000111`)
		testNumBits[int8](8, `00001000`)
		testNumBits[int8](9, `00001001`)

		testNumBits[int8](15, `00001111`)
		testNumBits[int8](16, `00010000`)
		testNumBits[int8](17, `00010001`)

		testNumBits[int8](31, `00011111`)
		testNumBits[int8](32, `00100000`)
		testNumBits[int8](33, `00100001`)

		testNumBits[int8](63, `00111111`)
		testNumBits[int8](64, `01000000`)
		testNumBits[int8](65, `01000001`)

		testNumBits[int8](127, `01111111`)

		testNumBits[int8](-128, `10000000`)
		testNumBits[int8](-127, `10000001`)
		testNumBits[int8](-126, `10000010`)

		testNumBits[int8](-125, `10000011`)
		testNumBits[int8](-124, `10000100`)
		testNumBits[int8](-123, `10000101`)

		testNumBits[int8](-121, `10000111`)
		testNumBits[int8](-120, `10001000`)
		testNumBits[int8](-119, `10001001`)

		testNumBits[int8](-113, `10001111`)
		testNumBits[int8](-112, `10010000`)
		testNumBits[int8](-111, `10010001`)

		testNumBits[int8](-97, `10011111`)
		testNumBits[int8](-96, `10100000`)
		testNumBits[int8](-95, `10100001`)

		testNumBits[int8](-65, `10111111`)
		testNumBits[int8](-64, `11000000`)
		testNumBits[int8](-63, `11000001`)

		testNumBits[int8](-1, `11111111`)
	}

	{
		testNumBits[uint16](0, `0000000000000000`)
		testNumBits[uint16](1, `0000000000000001`)
		testNumBits[uint16](2, `0000000000000010`)

		testNumBits[uint16](3, `0000000000000011`)
		testNumBits[uint16](4, `0000000000000100`)
		testNumBits[uint16](5, `0000000000000101`)

		testNumBits[uint16](7, `0000000000000111`)
		testNumBits[uint16](8, `0000000000001000`)
		testNumBits[uint16](9, `0000000000001001`)

		testNumBits[uint16](15, `0000000000001111`)
		testNumBits[uint16](16, `0000000000010000`)
		testNumBits[uint16](17, `0000000000010001`)

		testNumBits[uint16](31, `0000000000011111`)
		testNumBits[uint16](32, `0000000000100000`)
		testNumBits[uint16](33, `0000000000100001`)

		testNumBits[uint16](63, `0000000000111111`)
		testNumBits[uint16](64, `0000000001000000`)
		testNumBits[uint16](65, `0000000001000001`)

		testNumBits[uint16](127, `0000000001111111`)
		testNumBits[uint16](128, `0000000010000000`)
		testNumBits[uint16](129, `0000000010000001`)

		testNumBits[uint16](255, `0000000011111111`)
		testNumBits[uint16](256, `0000000100000000`)
		testNumBits[uint16](257, `0000000100000001`)

		testNumBits[uint16](511, `0000000111111111`)
		testNumBits[uint16](512, `0000001000000000`)
		testNumBits[uint16](513, `0000001000000001`)

		testNumBits[uint16](1023, `0000001111111111`)
		testNumBits[uint16](1024, `0000010000000000`)
		testNumBits[uint16](1025, `0000010000000001`)

		testNumBits[uint16](2047, `0000011111111111`)
		testNumBits[uint16](2048, `0000100000000000`)
		testNumBits[uint16](2049, `0000100000000001`)

		testNumBits[uint16](4095, `0000111111111111`)
		testNumBits[uint16](4096, `0001000000000000`)
		testNumBits[uint16](4097, `0001000000000001`)

		testNumBits[uint16](8191, `0001111111111111`)
		testNumBits[uint16](8192, `0010000000000000`)
		testNumBits[uint16](8193, `0010000000000001`)

		testNumBits[uint16](16383, `0011111111111111`)
		testNumBits[uint16](16384, `0100000000000000`)
		testNumBits[uint16](16385, `0100000000000001`)

		testNumBits[uint16](32767, `0111111111111111`)
		testNumBits[uint16](32768, `1000000000000000`)
		testNumBits[uint16](32769, `1000000000000001`)

		testNumBits[uint16](65535, `1111111111111111`)
	}

	{
		testNumBits[int16](0, `0000000000000000`)
		testNumBits[int16](1, `0000000000000001`)
		testNumBits[int16](2, `0000000000000010`)

		testNumBits[int16](3, `0000000000000011`)
		testNumBits[int16](4, `0000000000000100`)
		testNumBits[int16](5, `0000000000000101`)

		testNumBits[int16](7, `0000000000000111`)
		testNumBits[int16](8, `0000000000001000`)
		testNumBits[int16](9, `0000000000001001`)

		testNumBits[int16](15, `0000000000001111`)
		testNumBits[int16](16, `0000000000010000`)
		testNumBits[int16](17, `0000000000010001`)

		testNumBits[int16](31, `0000000000011111`)
		testNumBits[int16](32, `0000000000100000`)
		testNumBits[int16](33, `0000000000100001`)

		testNumBits[int16](63, `0000000000111111`)
		testNumBits[int16](64, `0000000001000000`)
		testNumBits[int16](65, `0000000001000001`)

		testNumBits[int16](127, `0000000001111111`)
		testNumBits[int16](128, `0000000010000000`)
		testNumBits[int16](129, `0000000010000001`)

		testNumBits[int16](255, `0000000011111111`)
		testNumBits[int16](256, `0000000100000000`)
		testNumBits[int16](257, `0000000100000001`)

		testNumBits[int16](511, `0000000111111111`)
		testNumBits[int16](512, `0000001000000000`)
		testNumBits[int16](513, `0000001000000001`)

		testNumBits[int16](1023, `0000001111111111`)
		testNumBits[int16](1024, `0000010000000000`)
		testNumBits[int16](1025, `0000010000000001`)

		testNumBits[int16](2047, `0000011111111111`)
		testNumBits[int16](2048, `0000100000000000`)
		testNumBits[int16](2049, `0000100000000001`)

		testNumBits[int16](4095, `0000111111111111`)
		testNumBits[int16](4096, `0001000000000000`)
		testNumBits[int16](4097, `0001000000000001`)

		testNumBits[int16](8191, `0001111111111111`)
		testNumBits[int16](8192, `0010000000000000`)
		testNumBits[int16](8193, `0010000000000001`)

		testNumBits[int16](16383, `0011111111111111`)
		testNumBits[int16](16384, `0100000000000000`)
		testNumBits[int16](16385, `0100000000000001`)

		testNumBits[int16](32767, `0111111111111111`)

		testNumBits[int16](-32768, `1000000000000000`)
		testNumBits[int16](-32767, `1000000000000001`)
		testNumBits[int16](-32766, `1000000000000010`)

		testNumBits[int16](-32765, `1000000000000011`)
		testNumBits[int16](-32764, `1000000000000100`)
		testNumBits[int16](-32763, `1000000000000101`)

		testNumBits[int16](-32761, `1000000000000111`)
		testNumBits[int16](-32760, `1000000000001000`)
		testNumBits[int16](-32759, `1000000000001001`)

		testNumBits[int16](-32753, `1000000000001111`)
		testNumBits[int16](-32752, `1000000000010000`)
		testNumBits[int16](-32751, `1000000000010001`)

		testNumBits[int16](-32737, `1000000000011111`)
		testNumBits[int16](-32736, `1000000000100000`)
		testNumBits[int16](-32735, `1000000000100001`)

		testNumBits[int16](-32705, `1000000000111111`)
		testNumBits[int16](-32704, `1000000001000000`)
		testNumBits[int16](-32703, `1000000001000001`)

		testNumBits[int16](-32641, `1000000001111111`)
		testNumBits[int16](-32640, `1000000010000000`)
		testNumBits[int16](-32639, `1000000010000001`)

		testNumBits[int16](-32513, `1000000011111111`)
		testNumBits[int16](-32512, `1000000100000000`)
		testNumBits[int16](-32511, `1000000100000001`)

		testNumBits[int16](-32257, `1000000111111111`)
		testNumBits[int16](-32256, `1000001000000000`)
		testNumBits[int16](-32255, `1000001000000001`)

		testNumBits[int16](-31745, `1000001111111111`)
		testNumBits[int16](-31744, `1000010000000000`)
		testNumBits[int16](-31743, `1000010000000001`)

		testNumBits[int16](-30721, `1000011111111111`)
		testNumBits[int16](-30720, `1000100000000000`)
		testNumBits[int16](-30719, `1000100000000001`)

		testNumBits[int16](-28673, `1000111111111111`)
		testNumBits[int16](-28672, `1001000000000000`)
		testNumBits[int16](-28671, `1001000000000001`)

		testNumBits[int16](-24577, `1001111111111111`)
		testNumBits[int16](-24576, `1010000000000000`)
		testNumBits[int16](-24575, `1010000000000001`)

		testNumBits[int16](-16385, `1011111111111111`)
		testNumBits[int16](-16384, `1100000000000000`)
		testNumBits[int16](-16383, `1100000000000001`)

		testNumBits[int16](-1, `1111111111111111`)
	}
}

func testNumBits[A gg.Int](src A, exp string) {
	gtest.Eq(gg.NumBits[A](src), exp)
}

func Benchmark_num_bits_strconv(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		strconv.FormatUint(math.MaxUint64/3, 2)
	}
}

func Benchmark_num_bits_ours(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.NumBits[uint64](math.MaxUint64 / 3)
	}
}
