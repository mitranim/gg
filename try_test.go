package gg_test

import (
	"errors"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestCatch(t *testing.T) {
	defer gtest.Catch(t)

	caught := gg.Catch(func() { panic(`some cause`) })

	var err gg.Err
	gtest.True(errors.As(caught, &err))

	gtest.Equal(err.Msg, ``)
	gtest.Equal(err.Cause, error(gg.ErrStr(`some cause`)))
	gtest.True(err.Trace.HasLen())
}

func TestDetailf(t *testing.T) {
	defer gtest.Catch(t)

	err := gg.Catch(func() {
		defer gg.Detailf(`unable to %v`, `do stuff`)
		panic(`some cause`)
	}).(gg.Err)

	gtest.Equal(err.Msg, `unable to do stuff`)
	gtest.Equal(err.Cause, error(gg.ErrStr(`some cause`)))
	gtest.True(err.Trace.HasLen())
}
