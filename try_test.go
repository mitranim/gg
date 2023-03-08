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
	gtest.True(err.Trace.HasLen())

	err = gg.Catch(func() { panic(gg.ErrStr(`string_error`)) }).(gg.Err)
	gtest.Zero(err.Msg)
	gtest.Equal(err.Cause, error(gg.ErrStr(`string_error`)))
	gtest.True(err.Trace.HasLen())
}

func TestDetailf(t *testing.T) {
	defer gtest.Catch(t)

	err := gg.Catch(func() {
		defer gg.Detailf(`unable to %v`, `do stuff`)
		panic(`string_panic`)
	}).(gg.Err)

	gtest.Equal(err.Msg, `unable to do stuff`)
	gtest.Equal(err.Cause, error(gg.ErrStr(`string_panic`)))
	gtest.True(err.Trace.HasLen())
}
