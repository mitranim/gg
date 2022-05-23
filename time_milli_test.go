package gg_test

import (
	"testing"
	"time"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestTimeMilli(t *testing.T) {
	t.Run(`UTC`, func(t *testing.T) {
		defer gtest.Catch(t)
		defer gg.Swap(&time.Local, nil)
		testTimeMilli(t)
	})

	t.Run(`non_UTC`, func(t *testing.T) {
		defer gtest.Catch(t)
		defer gg.Swap(&time.Local, time.FixedZone(``, 60*60*3))
		testTimeMilli(t)
	})
}

func testTimeMilli(t *testing.T) {
	t.Run(`IsNull`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.True(gg.TimeMilli(0).IsNull())
		gtest.False(gg.TimeMilli(-1).IsNull())
		gtest.False(gg.TimeMilli(1).IsNull())
	})

	t.Run(`Clear`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.NoPanic((*gg.TimeMilli)(nil).Clear)

		tar := gg.TimeMilli(10)

		tar.Clear()
		gtest.Zero(tar)

		tar.Clear()
		gtest.Zero(tar)
	})

	t.Run(`Time`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.TimeMilli(0).Time(), time.UnixMilli(0))
		gtest.Eq(gg.TimeMilli(1).Time(), time.UnixMilli(1))
		gtest.Eq(gg.TimeMilli(123).Time(), time.UnixMilli(123))
		gtest.Eq(gg.TimeMilli(123456).Time(), time.UnixMilli(123456))
	})

	t.Run(`Get`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(gg.TimeMilli(0).Get(), any(nil))
		gtest.Equal(gg.TimeMilli(1).Get(), any(time.UnixMilli(1)))
		gtest.Equal(gg.TimeMilli(123).Get(), any(time.UnixMilli(123)))
		gtest.Equal(gg.TimeMilli(123456).Get(), any(time.UnixMilli(123456)))
	})

	t.Run(`SetInt64`, func(t *testing.T) {
		defer gtest.Catch(t)

		var tar gg.TimeMilli

		tar.SetInt64(123)
		gtest.Eq(tar, 123)

		tar.SetInt64(0)
		gtest.Zero(tar)

		tar.SetInt64(-123)
		gtest.Eq(tar, -123)
	})

	t.Run(`SetTime`, func(t *testing.T) {
		defer gtest.Catch(t)

		var tar gg.TimeMilli

		tar.SetTime(time.UnixMilli(-123))
		gtest.Eq(tar, -123)

		tar.SetTime(time.UnixMilli(0))
		gtest.Zero(tar)

		tar.SetTime(time.UnixMilli(123))
		gtest.Eq(tar, 123)
	})

	t.Run(`Parse`, func(t *testing.T) {
		testTimeMilliParse(t, (*gg.TimeMilli).Parse)
	})

	t.Run(`String`, func(t *testing.T) {
		defer gtest.Catch(t)
		testTimeMilliString(gg.TimeMilli.String)
	})

	t.Run(`Append`, func(t *testing.T) {
		defer gtest.Catch(t)

		testTimeMilliString(gg.AppenderString[gg.TimeMilli])

		test := func(src string, tar gg.TimeMilli, exp string) {
			gtest.Eq(
				gg.ToString(tar.Append(gg.ToBytes(src))),
				exp,
			)
		}

		test(`<prefix>`, gg.TimeMilli(0), `<prefix>`)
		test(`<prefix>`, gg.TimeMilli(1), `<prefix>1`)
		test(`<prefix>`, gg.TimeMilli(-1), `<prefix>-1`)
		test(`<prefix>`, gg.TimeMilli(123), `<prefix>123`)
		test(`<prefix>`, gg.TimeMilli(-123), `<prefix>-123`)
	})

	t.Run(`MarshalText`, func(t *testing.T) {
		defer gtest.Catch(t)
		testTimeMilliString(timeMilliStringViaMarshalText)
	})

	t.Run(`UnmarshalText`, func(t *testing.T) {
		defer gtest.Catch(t)
		testTimeMilliParse(t, timeMilliParseViaUnmarshalText)
	})

	t.Run(`MarshalJSON`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src gg.TimeMilli, exp string) {
			gtest.Eq(gg.ToString(gg.Try1(src.MarshalJSON())), exp)
		}

		test(0, `null`)
		test(123, `123`)
		test(-123, `-123`)
	})

	t.Run(`UnmarshalJSON`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src string, exp gg.TimeMilli) {
			var tar gg.TimeMilli
			gtest.NoError(tar.UnmarshalJSON(gg.ToBytes(src)))
			gtest.Eq(tar, exp)
		}

		test(`null`, 0)

		test(`123`, 123)
		test(`-123`, -123)

		test(`"0001-01-01T00:00:00Z"`, gg.TimeMilli(timeZeroToUnixMilli()))
		test(`"1234-05-06T07:08:09.123456789Z"`, -23215049510877)
		test(`"9234-05-06T07:08:09.123456789Z"`, 229240566489123)
	})

	t.Run(`Scan`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src any, exp gg.TimeMilli) {
			var tar gg.TimeMilli

			tar = 0
			gtest.NoError(tar.Scan(src))
			gtest.Eq(tar, exp)

			tar = 123
			gtest.NoError(tar.Scan(src))
			gtest.Eq(tar, exp)
		}

		test(nil, 0)
		test((*time.Time)(nil), 0)
		test(``, 0)
		test([]byte(nil), 0)
		test([]byte{}, 0)

		test(time.UnixMilli(0), 0)
		test(time.UnixMilli(123), 123)
		test(time.UnixMilli(-123), -123)

		test(int64(0), 0)
		test(int64(123), 123)
		test(int64(-123), -123)

		test(gg.TimeMilli(0), 0)
		test(gg.TimeMilli(123), 123)
		test(gg.TimeMilli(-123), -123)

		test(`0`, 0)
		test(`123`, 123)
		test(`-123`, -123)

		test(`0001-01-01T00:00:00Z`, gg.TimeMilli(timeZeroToUnixMilli()))
		test(`1234-05-06T07:08:09.123456789Z`, -23215049510877)
		test(`9234-05-06T07:08:09.123456789Z`, 229240566489123)
	})
}

func testTimeMilliString(fun func(gg.TimeMilli) string) {
	gtest.Eq(fun(gg.TimeMilli(0)), ``)
	gtest.Eq(fun(gg.TimeMilli(1)), `1`)
	gtest.Eq(fun(gg.TimeMilli(-1)), `-1`)
	gtest.Eq(fun(gg.TimeMilli(123)), `123`)
	gtest.Eq(fun(gg.TimeMilli(-123)), `-123`)
}

func timeMilliStringViaMarshalText(src gg.TimeMilli) string {
	return gg.ToString(gg.Try1(src.MarshalText()))
}

func testTimeMilliParse(
	t *testing.T,
	fun func(*gg.TimeMilli, string) error,
) {
	t.Run(`invalid`, func(t *testing.T) {
		defer gtest.Catch(t)

		var tar gg.TimeMilli

		gtest.ErrorStr(
			`parsing time "wtf" as "2006-01-02T15:04:05Z07:00"`,
			fun(&tar, `wtf`),
		)
		gtest.Zero(tar)
	})

	t.Run(`empty`, func(t *testing.T) {
		defer gtest.Catch(t)

		tar := gg.TimeMilli(123)
		gtest.NoError(fun(&tar, ``))
		gtest.Zero(tar)
	})

	t.Run(`integer`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src string, exp gg.TimeMilli) {
			var tar gg.TimeMilli
			gtest.NoError(fun(&tar, src))
			gtest.Eq(tar, exp)
		}

		test(`0`, 0)
		test(`-0`, 0)
		test(`+0`, 0)

		test(`1`, 1)
		test(`-1`, -1)
		test(`+1`, +1)

		test(`12`, 12)
		test(`-12`, -12)
		test(`+12`, +12)

		test(`123`, 123)
		test(`-123`, -123)
		test(`+123`, +123)
	})

	t.Run(`RFC3339`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src string, exp gg.TimeMilli) {
			var tar gg.TimeMilli
			gtest.NoError(fun(&tar, src))
			gtest.Eq(tar, exp)

			inst := timeParse(src)

			gtest.Eq(
				tar.Time(),
				time.UnixMilli(inst.UnixMilli()).In(inst.Location()),
			)
		}

		test(`0001-01-01T00:00:00Z`, gg.TimeMilli(timeZeroToUnixMilli()))
		test(`1234-05-06T07:08:09.123456789Z`, -23215049510877)
		test(`9234-05-06T07:08:09.123456789Z`, 229240566489123)
	})
}

func timeMilliParseViaUnmarshalText(tar *gg.TimeMilli, src string) error {
	return tar.UnmarshalText(gg.ToBytes(src))
}

func timeZeroToUnixMilli() int64 { return time.Time{}.UnixMilli() }

func BenchmarkTimeMilli_Parse_integer(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.TimeMilliParse(`-1234567890123`)
	}
}

func BenchmarkTimeMilli_Parse_RFC3339(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.TimeMilliParse(`1234-05-06T07:08:09Z`)
	}
}
