package gg_test

import (
	"testing"
	"time"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestTimeMicro(t *testing.T) {
	t.Run(`UTC`, func(t *testing.T) {
		defer gtest.Catch(t)
		defer gg.SnapSwap(&time.Local, nil)
		testTimeMicro(t)
	})

	t.Run(`non_UTC`, func(t *testing.T) {
		defer gtest.Catch(t)
		defer gg.SnapSwap(&time.Local, time.FixedZone(``, 60*60*3))
		testTimeMicro(t)
	})
}

func testTimeMicro(t *testing.T) {
	t.Run(`IsNull`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.True(gg.TimeMicro(0).IsNull())
		gtest.False(gg.TimeMicro(-1).IsNull())
		gtest.False(gg.TimeMicro(1).IsNull())
	})

	t.Run(`Clear`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.NotPanic((*gg.TimeMicro)(nil).Clear)

		tar := gg.TimeMicro(10)

		tar.Clear()
		gtest.Zero(tar)

		tar.Clear()
		gtest.Zero(tar)
	})

	t.Run(`Time`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.TimeMicro(0).Time(), time.UnixMicro(0))
		gtest.Eq(gg.TimeMicro(1).Time(), time.UnixMicro(1))
		gtest.Eq(gg.TimeMicro(123).Time(), time.UnixMicro(123))
		gtest.Eq(gg.TimeMicro(123456).Time(), time.UnixMicro(123456))
	})

	t.Run(`Get`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(gg.TimeMicro(0).Get(), any(nil))
		gtest.Equal(gg.TimeMicro(1).Get(), any(time.UnixMicro(1)))
		gtest.Equal(gg.TimeMicro(123).Get(), any(time.UnixMicro(123)))
		gtest.Equal(gg.TimeMicro(123456).Get(), any(time.UnixMicro(123456)))
	})

	t.Run(`SetInt64`, func(t *testing.T) {
		defer gtest.Catch(t)

		var tar gg.TimeMicro

		tar.SetInt64(123)
		gtest.Eq(tar, 123)

		tar.SetInt64(0)
		gtest.Zero(tar)

		tar.SetInt64(-123)
		gtest.Eq(tar, -123)
	})

	t.Run(`SetTime`, func(t *testing.T) {
		defer gtest.Catch(t)

		var tar gg.TimeMicro

		tar.SetTime(time.UnixMicro(-123))
		gtest.Eq(tar, -123)

		tar.SetTime(time.UnixMicro(0))
		gtest.Zero(tar)

		tar.SetTime(time.UnixMicro(123))
		gtest.Eq(tar, 123)
	})

	t.Run(`Parse`, func(t *testing.T) {
		testTimeMicroParse(t, (*gg.TimeMicro).Parse)
	})

	t.Run(`String`, func(t *testing.T) {
		defer gtest.Catch(t)
		testTimeMicroString(gg.TimeMicro.String)
	})

	t.Run(`AppenderTo`, func(t *testing.T) {
		defer gtest.Catch(t)

		testTimeMicroString(gg.AppenderString[gg.TimeMicro])

		test := func(src string, tar gg.TimeMicro, exp string) {
			gtest.Eq(
				gg.ToString(tar.AppendTo(gg.ToBytes(src))),
				exp,
			)
		}

		test(`<prefix>`, gg.TimeMicro(0), `<prefix>`)
		test(`<prefix>`, gg.TimeMicro(1), `<prefix>1`)
		test(`<prefix>`, gg.TimeMicro(-1), `<prefix>-1`)
		test(`<prefix>`, gg.TimeMicro(123), `<prefix>123`)
		test(`<prefix>`, gg.TimeMicro(-123), `<prefix>-123`)
	})

	t.Run(`MarshalText`, func(t *testing.T) {
		defer gtest.Catch(t)
		testTimeMicroString(timeMicroStringViaMarshalText)
	})

	t.Run(`UnmarshalText`, func(t *testing.T) {
		defer gtest.Catch(t)
		testTimeMicroParse(t, timeMicroParseViaUnmarshalText)
	})

	t.Run(`MarshalJSON`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src gg.TimeMicro, exp string) {
			gtest.Eq(gg.ToString(gg.Try1(src.MarshalJSON())), exp)
		}

		test(0, `null`)
		test(123, `123`)
		test(-123, `-123`)
	})

	t.Run(`UnmarshalJSON`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src string, exp gg.TimeMicro) {
			var tar gg.TimeMicro
			gtest.NoError(tar.UnmarshalJSON(gg.ToBytes(src)))
			gtest.Eq(tar, exp)
		}

		test(`null`, 0)

		test(`123`, 123)
		test(`-123`, -123)

		test(`"0001-01-01T00:00:00Z"`, gg.TimeMicro(timeZeroToUnixMicro()))
		test(`"1234-05-06T07:08:09.123456789Z"`, -23215049510876544)
		test(`"9234-05-06T07:08:09.123456789Z"`, 229240566489123456)
	})

	t.Run(`Scan`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src any, exp gg.TimeMicro) {
			var tar gg.TimeMicro

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

		test(time.UnixMicro(0), 0)
		test(time.UnixMicro(123), 123)
		test(time.UnixMicro(-123), -123)

		test(int64(0), 0)
		test(int64(123), 123)
		test(int64(-123), -123)

		test(gg.TimeMicro(0), 0)
		test(gg.TimeMicro(123), 123)
		test(gg.TimeMicro(-123), -123)

		test(`0`, 0)
		test(`123`, 123)
		test(`-123`, -123)

		test(`0001-01-01T00:00:00Z`, gg.TimeMicro(timeZeroToUnixMicro()))
		test(`1234-05-06T07:08:09.123456789Z`, -23215049510876544)
		test(`9234-05-06T07:08:09.123456789Z`, 229240566489123456)
	})
}

func testTimeMicroString(fun func(gg.TimeMicro) string) {
	gtest.Eq(fun(gg.TimeMicro(0)), ``)
	gtest.Eq(fun(gg.TimeMicro(1)), `1`)
	gtest.Eq(fun(gg.TimeMicro(-1)), `-1`)
	gtest.Eq(fun(gg.TimeMicro(123)), `123`)
	gtest.Eq(fun(gg.TimeMicro(-123)), `-123`)
}

func timeMicroStringViaMarshalText(src gg.TimeMicro) string {
	return gg.ToString(gg.Try1(src.MarshalText()))
}

func testTimeMicroParse(
	t *testing.T,
	fun func(*gg.TimeMicro, string) error,
) {
	t.Run(`invalid`, func(t *testing.T) {
		defer gtest.Catch(t)

		var tar gg.TimeMicro

		gtest.ErrorStr(
			`parsing time "wtf" as "2006-01-02T15:04:05Z07:00"`,
			fun(&tar, `wtf`),
		)
		gtest.Zero(tar)
	})

	t.Run(`empty`, func(t *testing.T) {
		defer gtest.Catch(t)

		tar := gg.TimeMicro(123)
		gtest.NoError(fun(&tar, ``))
		gtest.Zero(tar)
	})

	t.Run(`integer`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src string, exp gg.TimeMicro) {
			var tar gg.TimeMicro
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

		test := func(src string, exp gg.TimeMicro) {
			var tar gg.TimeMicro
			gtest.NoError(fun(&tar, src))
			gtest.Eq(tar, exp)

			inst := timeParse(src)

			gtest.Eq(
				tar.Time(),
				time.UnixMicro(inst.UnixMicro()).In(inst.Location()),
			)
		}

		test(`0001-01-01T00:00:00Z`, gg.TimeMicro(timeZeroToUnixMicro()))
		test(`1234-05-06T07:08:09.123456789Z`, -23215049510876544)
		test(`9234-05-06T07:08:09.123456789Z`, 229240566489123456)
	})
}

func timeMicroParseViaUnmarshalText(tar *gg.TimeMicro, src string) error {
	return tar.UnmarshalText(gg.ToBytes(src))
}

func timeZeroToUnixMicro() int64 { return time.Time{}.UnixMicro() }

func timeParse(src string) time.Time {
	return gg.Try1(time.Parse(time.RFC3339, src)).In(gg.Or(time.Local, time.UTC))
}

func BenchmarkTimeMicro_Parse_integer(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.TimeMicroParse(`-1234567890123456`)
	}
}

func BenchmarkTimeMicro_Parse_RFC3339(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.TimeMicroParse(`1234-05-06T07:08:09Z`)
	}
}
