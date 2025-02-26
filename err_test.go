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

	err := gg.Err{}.Msgf(`unable to perform %v`, `some operation`).Caused(io.EOF).TracedAt(0)

	gtest.Eq(err.Error(), `unable to perform some operation: EOF`)

	gtest.Eq(err.Stack(), strings.TrimSpace(`
unable to perform some operation: EOF
trace:
    gg_test.TestErr err_test.go:18
`))
}

/*
Placed early in the file because this hardcodes line numbers, which could change
when adding or modifying tests.
*/
func TestErrStack(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.ErrStack(nil))

	t.Run(`without_trace`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.ErrStack(gg.ErrStr(`str`)), `str`)
		gtest.Eq(gg.ErrStack(gg.Err{Msg: `str`}), `str`)
		gtest.Eq(gg.ErrStack(fmt.Errorf(`str`)), `str`)
	})

	t.Run(`Err_outer_traced`, func(t *testing.T) {
		defer gtest.Catch(t)

		err := gg.Err{}.TracedAt(0)

		gtest.Eq(gg.ErrStack(err), strings.TrimSpace(`
trace:
    gg_test.TestErrStack.func2 err_test.go:49
`))
	})

	t.Run(`Err_inner_traced_outer_blank`, func(t *testing.T) {
		defer gtest.Catch(t)

		inner := gg.Err{Msg: `inner`}.TracedAt(0)
		outer := gg.Err{Cause: inner}

		gtest.Eq(gg.ErrStack(outer), inner.Stack())

		gtest.Eq(gg.ErrStack(outer), strings.TrimSpace(`
inner
trace:
    gg_test.TestErrStack.func3 err_test.go:60
`))
	})

	t.Run(`Err_inner_messaged_traced_outer_messaged`, func(t *testing.T) {
		defer gtest.Catch(t)

		inner := gg.Err{Msg: `inner`}.TracedAt(0)
		outer := gg.Err{Msg: `outer`, Cause: inner}

		gtest.Eq(gg.ErrStack(outer), strings.TrimSpace(`
outer: inner
trace:
    gg_test.TestErrStack.func4 err_test.go:75
`))
	})

	t.Run(`Err_inner_messaged_traced_outer_traced`, func(t *testing.T) {
		defer gtest.Catch(t)

		inner := gg.Err{Msg: `inner`}.TracedAt(0)
		outer := gg.Err{Cause: inner}.TracedAt(0)

		gtest.Eq(gg.ErrStack(outer), strings.TrimSpace(`
trace:
    gg_test.TestErrStack.func5 err_test.go:89
inner
trace:
    gg_test.TestErrStack.func5 err_test.go:88
`))
	})

	t.Run(`Err_inner_messaged_traced_outer_messaged_traced`, func(t *testing.T) {
		defer gtest.Catch(t)

		inner := gg.Err{Msg: `inner`}.TracedAt(0)
		outer := gg.Err{Msg: `outer`, Cause: inner}.TracedAt(0)

		gtest.Eq(gg.ErrStack(outer), strings.TrimSpace(`
outer
trace:
    gg_test.TestErrStack.func6 err_test.go:104
cause: inner
trace:
    gg_test.TestErrStack.func6 err_test.go:103
`))
	})
}

func TestErrTrace(t *testing.T) {
	defer gtest.Catch(t)

	inner := gg.Errf(`inner`)
	outer := fmt.Errorf(`outer: %w`, inner)

	gtest.Eq(outer.Error(), `outer: inner`)
	gtest.True(errors.Is(outer, inner))
	gtest.Equal(gg.ErrTrace(inner), gg.PtrGet(inner.Trace))
	gtest.Equal(gg.ErrTrace(outer), gg.PtrGet(inner.Trace))
}

func TestAnyErr(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.AnyErr(nil))
	gtest.AnyEq(gg.AnyErr(`str`), error(gg.ErrStr(`str`)))
	gtest.AnyEq(gg.AnyErr(gg.ErrStr(`str`)), error(gg.ErrStr(`str`)))
	gtest.AnyEq(gg.AnyErr(gg.ErrAny{`str`}), error(gg.ErrAny{`str`}))
}

func TestErrs_Error(t *testing.T) {
	defer gtest.Catch(t)
	testErrsError(gg.Errs.Error)
}

func testErrsError(fun func(gg.Errs) string) {
	gtest.Zero(fun(gg.Errs(nil)))
	gtest.Zero(fun(gg.Errs{}))

	gtest.Eq(
		fun(gg.Errs{nil, testErrTraced0, nil}),
		`test err traced 0`,
	)

	gtest.Eq(
		fun(gg.Errs{nil, testErrTraced0, nil, testErrTraced1, nil}),
		`multiple errors; test err traced 0; test err traced 1`,
	)
}

func TestErrs_Format_basic(t *testing.T) {
	defer gtest.Catch(t)
	testErrsError(func(val gg.Errs) string {
		return fmt.Sprintf(`%v`, val)
	})
}

func TestErrs_Format(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`empty`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(fmt.Sprintf(`%v`, gg.Errs(nil)), ``)
		gtest.Eq(fmt.Sprintf(`%+v`, gg.Errs(nil)), ``)
		gtest.Eq(fmt.Sprintf(`%#v`, gg.Errs(nil)), `gg.Errs(nil)`)

		gtest.Eq(fmt.Sprintf(`%v`, gg.Errs{}), ``)
		gtest.Eq(fmt.Sprintf(`%+v`, gg.Errs{}), ``)
		gtest.Eq(fmt.Sprintf(`%#v`, gg.Errs{}), `gg.Errs{}`)
	})

	const errStr = gg.ErrStr(`err_simple`)
	errTraced := gg.Errf(`err_traced`)
	errWrapped := gg.Wrapf(io.EOF, `err_wrapped`)

	t.Run(`single_untraced`, func(t *testing.T) {
		defer gtest.Catch(t)

		err := errStr
		src := gg.Errs{err}

		gtest.Eq(fmt.Sprintf(`%v`, src), err.Error())
		gtest.Eq(fmt.Sprintf(`%+v`, src), err.Error())
		gtest.Eq(fmt.Sprintf(`%#v`, src), `gg.Errs{"err_simple"}`)
	})

	t.Run(`single_traced`, func(t *testing.T) {
		defer gtest.Catch(t)

		err := errTraced
		src := gg.Errs{err}

		gtest.Eq(fmt.Sprintf(`%v`, src), err.Error())
		gtest.Eq(fmt.Sprintf(`%+v`, src), err.Stack())
		gtest.Eq(fmt.Sprintf(`%#v`, src), `gg.Errs{`+fmt.Sprintf(`%#v`, err)+`}`)
	})

	t.Run(`multiple_mixed`, func(t *testing.T) {
		defer gtest.Catch(t)

		src := gg.Errs{nil, errStr, nil, errTraced, nil, errWrapped, nil}

		gtest.Eq(
			fmt.Sprintf(`%v`, src),
			`multiple errors; err_simple; err_traced; err_wrapped: EOF`,
		)

		var buf gg.Buf
		buf.AppendString(`multiple errors:`)

		buf.AppendNewlines(2)
		buf.AppendString(errStr.Error())

		buf.AppendNewlines(2)
		buf.Fprintf(`%+v`, errTraced)

		buf.AppendNewlines(2)
		buf.Fprintf(`%+v`, errWrapped)

		gtest.Eq(fmt.Sprintf(`%+v`, src), buf.String())

		gtest.Eq(
			fmt.Sprintf(`%#v`, src),
			gg.Str(
				`gg.Errs{`,
				`error(nil)`,
				`, `,
				`"err_simple"`,
				`, `,
				`error(nil)`,
				`, `,
				fmt.Sprintf(`%#v`, errTraced),
				`, `,
				`error(nil)`,
				`, `,
				fmt.Sprintf(`%#v`, errWrapped),
				`, `,
				`error(nil)`,
				`}`,
			),
		)
	})
}

func TestErrs_Err(t *testing.T) {
	defer gtest.Catch(t)

	testEmpty := func(src gg.Errs) {
		gtest.Zero(src.Err())
	}

	testEmpty(gg.Errs(nil))
	testEmpty(gg.Errs{})
	testEmpty(gg.Errs{nil, nil, nil})

	testOne := func(exp error) {
		gtest.Equal(gg.Errs{nil, exp, nil}.Err(), exp)
		gtest.Equal(gg.Errs{exp, nil}.Err(), exp)
		gtest.Equal(gg.Errs{nil, exp}.Err(), exp)
	}

	testOne(testErrTraced0)
	testOne(testErrTraced1)

	errs := gg.Errs{nil, testErrTraced0, nil, testErrTraced1, nil}
	gtest.Equal(errs.Err(), error(errs))
}

func TestErrs_Unwrap(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.Errs(nil).Unwrap())
	gtest.Zero(gg.Errs{}.Unwrap())
	gtest.Equal(gg.Errs{nil, testErrTraced0}.Unwrap(), testErrTraced0)
	gtest.Equal(gg.Errs{nil, testErrTraced0, testErrTraced1}.Unwrap(), testErrTraced0)
}

func TestErrs_Is(t *testing.T) {
	defer gtest.Catch(t)

	gtest.False(errors.Is(gg.Errs(nil), io.EOF))
	gtest.False(errors.Is(gg.Errs(nil), testErrTraced0))

	gtest.False(errors.Is(gg.Errs{nil, testErrTraced0, nil, testErrTraced1, nil}, io.EOF))
	gtest.False(errors.Is(gg.Errs{nil, testErrTraced0, nil, testErrTraced1, nil}, testErrTraced2))

	gtest.True(errors.Is(gg.Errs{nil, testErrTraced0, nil, testErrTraced1, nil}, testErrTraced0))
	gtest.True(errors.Is(gg.Errs{nil, testErrTraced0, nil, testErrTraced1, nil}, testErrTraced1))

	gtest.True(errors.Is(gg.Errs{nil, gg.Wrapf(testErrTraced0, ``), nil, testErrTraced1, nil}, testErrTraced0))
	gtest.True(errors.Is(gg.Errs{nil, gg.Wrapf(io.EOF, ``), nil, testErrTraced1, nil}, io.EOF))

	gtest.True(errors.Is(gg.Errs{nil, testErrTraced0, nil, gg.Wrapf(testErrTraced1, ``), nil}, testErrTraced1))
	gtest.True(errors.Is(gg.Errs{nil, testErrTraced0, nil, gg.Wrapf(io.EOF, ``), nil}, io.EOF))
}

func TestErrs_As(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src gg.Errs, ok bool, exp gg.ErrStr) {
		var tar gg.ErrStr
		gtest.Eq(errors.As(src, &tar), ok)
		gtest.Eq(tar, exp)
	}

	test(gg.Errs(nil), false, ``)
	test(gg.Errs{}, false, ``)
	test(gg.Errs{testErrTraced0}, false, ``)
	test(gg.Errs{testErrTraced0, testErrTraced1}, false, ``)
	test(gg.Errs{testErrUntracedA, testErrTraced0, testErrTraced1}, true, testErrUntracedA)
	test(gg.Errs{testErrTraced0, testErrUntracedA, testErrTraced1}, true, testErrUntracedA)
	test(gg.Errs{testErrTraced0, testErrTraced1, testErrUntracedA}, true, testErrUntracedA)
	test(gg.Errs{nil, testErrTraced0, nil, testErrTraced1, nil, testErrUntracedA, nil}, true, testErrUntracedA)
	test(gg.Errs{nil, testErrUntracedA, nil, testErrUntracedB, nil}, true, testErrUntracedA)
}

func TestErrs_Find(t *testing.T) {
	defer gtest.Catch(t)

	for _, src := range []gg.Errs{nil, {}, {nil}, {nil, nil}} {
		for _, fun := range []func(error) bool{nil, True1[error], False1[error]} {
			gtest.Zero(src.Find(fun))
		}
	}

	match := func(err error) bool { return err == testErrTraced0 }

	for _, src := range []gg.Errs{
		{testErrTraced0, testErrTraced1, testErrUntracedA, testErrUntracedB},
		{testErrTraced1, testErrTraced0, testErrUntracedA, testErrUntracedB},
		{testErrTraced1, testErrUntracedA, testErrTraced0, testErrUntracedB},
		{testErrTraced1, testErrUntracedA, testErrUntracedB, testErrTraced0},
		{testErrTraced1, nil, nil, testErrTraced0},
		{testErrTraced1, gg.Wrapf(fmt.Errorf(`one: %w`, testErrTraced0), `two`)},
	} {
		gtest.Equal(src.Find(match), testErrTraced0)
	}
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
		gtest.True(gg.PtrGet(err.Trace).IsNotEmpty())
	})

	t.Run(`cause_with_stack`, func(t *testing.T) {
		defer gtest.Catch(t)

		err := gg.Wrapf(gg.Errf(`some cause`), `unable to %v`, `do stuff`).(gg.Err)

		gtest.Equal(err.Msg, `unable to do stuff`)
		gtest.False(gg.PtrGet(err.Trace).IsNotEmpty())

		cause := err.Cause.(gg.Err)
		gtest.Equal(cause.Msg, `some cause`)
		gtest.True(gg.PtrGet(cause.Trace).IsNotEmpty())
	})
}

/*
Internally `Wrap` is very similar to `Wrapf` which has a fuller test.
Here, we only need to test message generation.
*/
func TestWrap(t *testing.T) {
	defer gtest.Catch(t)

	err := gg.Wrap(gg.ErrAny{`some cause`}, `unable to %v `, `do stuff`).(gg.Err)
	gtest.Equal(err.Msg, `unable to %v do stuff`)
	gtest.Equal(err.Cause, error(gg.ErrAny{`some cause`}))
}

func BenchmarkIsErrTraced_error_without_trace(b *testing.B) {
	const err = gg.ErrStr(`some error`)
	gtest.False(gg.IsErrTraced(err))

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsErrTraced(err))
	}
}

func BenchmarkIsErrTraced_error_with_trace(b *testing.B) {
	err := error(gg.Errf(`some error`))
	gtest.True(gg.IsErrTraced(err))

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsErrTraced(err))
	}
}

func BenchmarkErrTraced_error_without_trace(b *testing.B) {
	const err = gg.ErrStr(`some error`)
	gtest.False(gg.IsErrTraced(err))

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ErrTracedAt(err, 1))
	}
}

func BenchmarkErrTraced_error_with_trace(b *testing.B) {
	err := error(gg.Errf(`some error`))
	gtest.True(gg.IsErrTraced(err))

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.ErrTracedAt(err, 1))
	}
}

func TestErrAs(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.ErrAs[gg.Err](nil))
	gtest.Zero(gg.ErrAs[gg.ErrStr](nil))
	gtest.Zero(gg.ErrAs[PtrErrStr](nil))

	gtest.Zero(gg.ErrAs[gg.Err](fmt.Errorf(``)))
	gtest.Zero(gg.ErrAs[gg.ErrStr](fmt.Errorf(``)))
	gtest.Zero(gg.ErrAs[PtrErrStr](fmt.Errorf(``)))

	gtest.Zero(gg.ErrAs[gg.Err](fmt.Errorf(``)))
	gtest.Zero(gg.ErrAs[gg.ErrStr](fmt.Errorf(``)))
	gtest.Zero(gg.ErrAs[PtrErrStr](fmt.Errorf(``)))

	gtest.Zero(gg.ErrAs[gg.Err](fmt.Errorf(`%w`, fmt.Errorf(``))))
	gtest.Zero(gg.ErrAs[gg.ErrStr](fmt.Errorf(`%w`, fmt.Errorf(``))))
	gtest.Zero(gg.ErrAs[PtrErrStr](fmt.Errorf(`%w`, fmt.Errorf(``))))

	gtest.Eq(
		gg.ErrAs[gg.ErrStr](gg.ErrStr(`one`)),
		gg.ErrStr(`one`),
	)

	gtest.Eq(
		gg.ErrAs[gg.ErrStr](fmt.Errorf(`%w`, gg.ErrStr(`one`))),
		gg.ErrStr(`one`),
	)

	gtest.Eq(
		gg.ErrAs[PtrErrStr](gg.Ptr(PtrErrStr(`one`))),
		PtrErrStr(`one`),
	)

	gtest.Eq(
		gg.ErrAs[PtrErrStr](fmt.Errorf(`%w`, gg.Ptr(PtrErrStr(`one`)))),
		PtrErrStr(`one`),
	)
}

func TestErrFind(t *testing.T) {
	defer gtest.Catch(t)

	inner := error(gg.ErrStr(`one`))
	outer := gg.Wrapf(inner, `two`)

	gtest.Zero(gg.ErrFind(nil, nil))
	gtest.Zero(gg.ErrFind(nil, True1[error]))
	gtest.Zero(gg.ErrFind(nil, False1[error]))
	gtest.Zero(gg.ErrFind(inner, False1[error]))
	gtest.Zero(gg.ErrFind(outer, False1[error]))

	gtest.Equal(gg.ErrFind(inner, True1[error]), inner)
	gtest.Equal(gg.ErrFind(outer, True1[error]), outer)

	test := func(err error) bool { return err == inner }
	gtest.Equal(gg.ErrFind(outer, test), inner)
}
