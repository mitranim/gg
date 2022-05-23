package gg_test

import (
	"testing"
	"time"
	u "unsafe"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestParseTo(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(`unsupported kind interface`, func() { gg.ParseTo[any](``) })
	gtest.PanicStr(`unsupported kind struct`, func() { gg.ParseTo[SomeModel](``) })
	gtest.PanicStr(`unsupported kind slice`, func() { gg.ParseTo[[]string](``) })
	gtest.PanicStr(`unsupported kind chan`, func() { gg.ParseTo[chan struct{}](``) })
	gtest.PanicStr(`unsupported kind func`, func() { gg.ParseTo[func()](``) })
	gtest.PanicStr(`unsupported kind uintptr`, func() { gg.ParseTo[uintptr](``) })
	gtest.PanicStr(`unsupported kind unsafe.Pointer`, func() { gg.ParseTo[u.Pointer](``) })

	gtest.PanicStr(`invalid syntax`, func() { gg.ParseTo[int](``) })
	gtest.PanicStr(`invalid syntax`, func() { gg.ParseTo[*int](``) })

	gtest.Equal(gg.ParseTo[string](``), ``)
	gtest.Equal(gg.ParseTo[string](`str`), `str`)

	gtest.Equal(gg.ParseTo[*string](``), gg.Ptr(``))
	gtest.Equal(gg.ParseTo[*string](`str`), gg.Ptr(`str`))

	gtest.Equal(gg.ParseTo[int](`0`), 0)
	gtest.Equal(gg.ParseTo[int](`123`), 123)

	gtest.Equal(gg.ParseTo[*int](`0`), gg.Ptr(0))
	gtest.Equal(gg.ParseTo[*int](`123`), gg.Ptr(123))

	gtest.Equal(
		gg.ParseTo[time.Time](`1234-05-23T12:34:56Z`),
		time.Date(1234, 5, 23, 12, 34, 56, 0, time.UTC),
	)

	gtest.Equal(
		gg.ParseTo[*time.Time](`1234-05-23T12:34:56Z`),
		gg.Ptr(time.Date(1234, 5, 23, 12, 34, 56, 0, time.UTC)),
	)
}

func BenchmarkParseTo_int(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ParseTo[int](`123`))
	}
}

func BenchmarkParseTo_int_ptr(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ParseTo[*int](`123`))
	}
}

func BenchmarkParseTo_Parser(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ParseTo[ParserStr](`863872f79b1d4cc9a45e8027a6ad66ad`))
	}
}

func BenchmarkParseTo_Parser_ptr(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ParseTo[*ParserStr](`863872f79b1d4cc9a45e8027a6ad66ad`))
	}
}

func BenchmarkParseTo_Unmarshaler(b *testing.B) {
	src := []byte(`863872f79b1d4cc9a45e8027a6ad66ad`)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ParseTo[UnmarshalerBytes](src))
	}
}

func BenchmarkParseTo_Unmarshaler_ptr(b *testing.B) {
	src := []byte(`863872f79b1d4cc9a45e8027a6ad66ad`)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ParseTo[*UnmarshalerBytes](src))
	}
}

func BenchmarkParseTo_time_Time(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ParseTo[time.Time](`1234-05-23T12:34:56Z`))
	}
}

func BenchmarkParseTo_time_Time_ptr(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ParseTo[*time.Time](`1234-05-23T12:34:56Z`))
	}
}

func BenchmarkParse(b *testing.B) {
	var val int

	for ind := 0; ind < b.N; ind++ {
		gg.Parse(`123`, &val)
	}
}
