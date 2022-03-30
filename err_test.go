package gg_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

// Limited sanity check, TODO full test.
func TestErr(t *testing.T) {
	defer gtest.Catch(t)

	err := gg.Err{}.Msgf(`unable to perform %v`, `some operation`).Caused(io.EOF).Traced(0)

	gtest.Eq(err.Error(), `unable to perform some operation: EOF`)

	gtest.TextHas(
		err.Stack(),
		strings.TrimSpace(`
unable to perform some operation: EOF
trace:
    TestErr err_test.go
`),
	)
}

func TestErrTrace(t *testing.T) {
	defer gtest.Catch(t)

	inner := gg.Errf(`inner`)
	outer := fmt.Errorf(`outer: %w`, inner)

	gtest.Eq(outer.Error(), `outer: inner`)
	gtest.True(errors.Is(outer, inner))
	gtest.Equal(gg.ErrTrace(inner), gg.Deref(inner.Trace))
	gtest.Equal(gg.ErrTrace(outer), gg.Deref(inner.Trace))
}

func TestToErrAny(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.ToErrAny(nil))
	gtest.EqAny(gg.ToErrAny(`str`), error(gg.ErrStr(`str`)))
	gtest.EqAny(gg.ToErrAny(gg.ErrStr(`str`)), error(gg.ErrStr(`str`)))
	gtest.EqAny(gg.ToErrAny(gg.ErrAny{`str`}), error(gg.ErrAny{`str`}))
}

func TestErrs_Error(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.Errs(nil).Error())
	gtest.Zero(gg.Errs{}.Error())

	gtest.Eq(
		gg.Errs{nil, testErr0, nil}.Error(),
		`test err 0`,
	)

	gtest.Eq(
		gg.Errs{nil, testErr0, nil, testErr1, nil}.Error(),
		`multiple errors; test err 0; test err 1`,
	)
}

func TestErrs_ErrOpt(t *testing.T) {
	defer gtest.Catch(t)

	testEmpty := func(src gg.Errs) {
		t.Helper()
		gtest.Zero(src.ErrOpt())
	}

	testEmpty(gg.Errs(nil))
	testEmpty(gg.Errs{})
	testEmpty(gg.Errs{nil, nil, nil})

	testOne := func(exp error) {
		t.Helper()
		gtest.Equal(gg.Errs{nil, exp, nil}.ErrOpt(), exp)
		gtest.Equal(gg.Errs{exp, nil}.ErrOpt(), exp)
		gtest.Equal(gg.Errs{nil, exp}.ErrOpt(), exp)
	}

	testOne(testErr0)
	testOne(testErr1)

	errs := gg.Errs{nil, testErr0, nil, testErr1, nil}
	gtest.Equal(errs.ErrOpt(), error(errs))
}

func TestErrs_Unwrap(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.Errs(nil).Unwrap())
	gtest.Zero(gg.Errs{}.Unwrap())
	gtest.Equal(gg.Errs{nil, testErr0}.Unwrap(), testErr0)
	gtest.Equal(gg.Errs{nil, testErr0, testErr1}.Unwrap(), testErr0)
}

func TestErrs_Is(t *testing.T) {
	defer gtest.Catch(t)

	gtest.False(errors.Is(gg.Errs(nil), io.EOF))
	gtest.False(errors.Is(gg.Errs(nil), testErr0))

	gtest.False(errors.Is(gg.Errs{nil, testErr0, nil, testErr1, nil}, io.EOF))
	gtest.False(errors.Is(gg.Errs{nil, testErr0, nil, testErr1, nil}, testErr2))

	gtest.True(errors.Is(gg.Errs{nil, testErr0, nil, testErr1, nil}, testErr0))
	gtest.True(errors.Is(gg.Errs{nil, testErr0, nil, testErr1, nil}, testErr1))

	gtest.True(errors.Is(gg.Errs{nil, gg.Wrapf(testErr0, ``), nil, testErr1, nil}, testErr0))
	gtest.True(errors.Is(gg.Errs{nil, gg.Wrapf(io.EOF, ``), nil, testErr1, nil}, io.EOF))

	gtest.True(errors.Is(gg.Errs{nil, testErr0, nil, gg.Wrapf(testErr1, ``), nil}, testErr1))
	gtest.True(errors.Is(gg.Errs{nil, testErr0, nil, gg.Wrapf(io.EOF, ``), nil}, io.EOF))
}

func TestErrs_As(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src gg.Errs, ok bool, exp gg.ErrStr) {
		t.Helper()

		var tar gg.ErrStr
		gtest.Eq(errors.As(src, &tar), ok)
		gtest.Eq(tar, exp)
	}

	test(gg.Errs(nil), false, ``)
	test(gg.Errs{}, false, ``)
	test(gg.Errs{testErr0}, false, ``)
	test(gg.Errs{testErr0, testErr1}, false, ``)
	test(gg.Errs{testErrA, testErr0, testErr1}, true, testErrA)
	test(gg.Errs{testErr0, testErrA, testErr1}, true, testErrA)
	test(gg.Errs{testErr0, testErr1, testErrA}, true, testErrA)
	test(gg.Errs{nil, testErr0, nil, testErr1, nil, testErrA, nil}, true, testErrA)
	test(gg.Errs{nil, testErrA, nil, testErrB, nil}, true, testErrA)
}

func TestErrStr(t *testing.T) {
	defer gtest.Catch(t)

	const err0 gg.ErrStr = `err0`
	const err1 gg.ErrStr = `err1`

	wrap0 := gg.Wrapf(err0, `wrap0`)
	wrap1 := gg.Wrapf(err1, `wrap1`)

	gtest.False(errors.Is(err0, err1))
	gtest.False(errors.Is(err1, err0))
	gtest.False(errors.Is(wrap0, err1))
	gtest.False(errors.Is(wrap1, err0))

	gtest.True(errors.Is(err0, err0))
	gtest.True(errors.Is(err1, err1))
	gtest.True(errors.Is(wrap0, err0))
	gtest.True(errors.Is(wrap1, err1))
}

func TestWrapf(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`cause_without_stack`, func(t *testing.T) {
		defer gtest.Catch(t)

		err := gg.Wrapf(gg.ErrAny{`some cause`}, `unable to %v`, `do stuff`).(gg.Err)

		gtest.Equal(err.Msg, `unable to do stuff`)
		gtest.Equal(err.Cause, error(gg.ErrAny{`some cause`}))
		gtest.True(gg.Deref(err.Trace).HasLen())
	})

	t.Run(`cause_with_stack`, func(t *testing.T) {
		defer gtest.Catch(t)

		err := gg.Wrapf(gg.Errf(`some cause`), `unable to %v`, `do stuff`).(gg.Err)

		gtest.Equal(err.Msg, `unable to do stuff`)
		gtest.False(gg.Deref(err.Trace).HasLen())

		cause := err.Cause.(gg.Err)
		gtest.Equal(cause.Msg, `some cause`)
		gtest.True(gg.Deref(cause.Trace).HasLen())
	})
}

func BenchmarkIsErrTraced_error_without_trace(b *testing.B) {
	const err = gg.ErrStr(`some error`)
	gtest.False(gg.IsErrTraced(err))

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.IsErrTraced(err))
	}
}

func BenchmarkIsErrTraced_error_with_trace(b *testing.B) {
	err := error(gg.Errf(`some error`))
	gtest.True(gg.IsErrTraced(err))

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.IsErrTraced(err))
	}
}

func BenchmarkErrTraced_error_without_trace(b *testing.B) {
	const err = gg.ErrStr(`some error`)
	gtest.False(gg.IsErrTraced(err))

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.ErrTraced(err, 1))
	}
}

func BenchmarkErrTraced_error_with_trace(b *testing.B) {
	err := error(gg.Errf(`some error`))
	gtest.True(gg.IsErrTraced(err))

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.ErrTraced(err, 1))
	}
}
