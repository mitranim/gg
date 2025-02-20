package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestCatch(t *testing.T) {
	defer gtest.Catch(t)

	err := gg.Catch(func() { panic(`string_panic`) }).(gg.Err)
	gtest.Equal(err.Msg, `string_panic`)
	gtest.True(err.Trace.IsNotEmpty())

	err = gg.Catch(func() { panic(gg.ErrStr(`string_error`)) }).(gg.Err)
	gtest.Zero(err.Msg)
	gtest.Equal(err.Cause, error(gg.ErrStr(`string_error`)))
	gtest.True(err.Trace.IsNotEmpty())
}

func TestDetailf(t *testing.T) {
	defer gtest.Catch(t)

	err := gg.Catch(func() {
		defer gg.Detailf(`unable to %v`, `do stuff`)
		panic(`string_panic`)
	}).(gg.Err)

	gtest.Equal(err.Msg, `unable to do stuff`)
	gtest.Equal(err.Cause, error(gg.ErrStr(`string_panic`)))
	gtest.True(err.Trace.IsNotEmpty())
}

/*
func BenchmarkPanicSkip(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		benchmarkPanicSkip()
	}
}

func BenchmarkPanicSkipTraced(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		benchmarkPanicSkipTraced()
	}
}

func BenchmarkFileExists(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		benchmarkFileExists()
	}
}

func benchmarkPanicSkip() {
	defer gg.Skip()
	panic(`error_message`)
}

func benchmarkPanicSkipTraced() {
	defer gg.Skip()
	panic(gg.Err{}.Msgd(`error_message`).TracedAt(1))
}

func benchmarkFileExists() {
	gtest.True(gg.FileExists(`try_test.go`))
}
*/
