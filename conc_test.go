package gg_test

import (
	"context"
	e "errors"
	"sync"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/mitranim/gg/gtest"
)

const (
	testErrUntracedA = gg.ErrStr(`test err untraced A`)
	testErrUntracedB = gg.ErrStr(`test err untraced B`)
	testErrUntracedC = gg.ErrStr(`test err untraced C`)
)

var (
	testErrTraced0 = error(gg.Errf(`test err traced 0`))
	testErrTraced1 = error(gg.Errf(`test err traced 1`))
	testErrTraced2 = error(gg.Errf(`test err traced 2`))
)

func testPanicUntracedA() { panic(testErrUntracedA) }
func testPanicUntracedB() { panic(testErrUntracedB) }

func testPanicTraced0() { panic(testErrTraced0) }
func testPanicTraced1() { panic(testErrTraced1) }

func testNopCtx(context.Context) {}

func testPanicCtxUntracedA(context.Context) { panic(testErrUntracedA) }

func testPanicCtxTraced0(context.Context) { panic(testErrTraced0) }
func testPanicCtxTraced1(context.Context) { panic(testErrTraced1) }
func testPanicCtxTraced2(context.Context) { panic(testErrTraced2) }

func TestConc(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`no_panic`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Zero(gg.ConcCatch())
		gtest.Equal(gg.ConcCatch(nil, nil, nil), []error{nil, nil, nil})

		gtest.Equal(
			gg.ConcCatch(gg.Nop),
			[]error{nil},
		)

		gtest.Equal(
			gg.ConcCatch(gg.Nop, gg.Nop),
			[]error{nil, nil},
		)

		gtest.Equal(
			gg.ConcCatch(gg.Nop, nil, gg.Nop),
			[]error{nil, nil, nil},
		)

		gtest.Equal(
			gg.ConcCatch(nil, gg.Nop, nil, gg.Nop, nil),
			[]error{nil, nil, nil, nil, nil},
		)
	})

	t.Run(`only_panic`, func(t *testing.T) {
		defer gtest.Catch(t)

		testWrappedErrs(
			gg.ConcCatch(testPanicUntracedA),
			[]error{testErrUntracedA},
		)

		gtest.Equal(
			gg.ConcCatch(testPanicTraced0),
			[]error{testErrTraced0},
			`when only one function is provided and it panics with a traced error, that error must be preserved as-is without redundant wrapping`,
		)

		testWrappedErrs(
			gg.ConcCatch(testPanicUntracedA, testPanicUntracedB),
			[]error{testErrUntracedA, testErrUntracedB},
		)

		testWrappedErrs(
			gg.ConcCatch(testPanicTraced0, testPanicTraced1),
			[]error{testErrTraced0, testErrTraced1},
		)
	})

	t.Run(`mixed`, func(t *testing.T) {
		defer gtest.Catch(t)

		testWrappedErrs(
			gg.ConcCatch(gg.Nop, testPanicUntracedA, gg.Nop, testPanicTraced1, gg.Nop),
			gg.Errs{nil, testErrUntracedA, nil, testErrTraced1, nil},
		)
	})
}

func testWrappedErr(act, exp error, opt ...any) {
	gtest.ErrIs(act, exp, opt...)
	gtest.NotIs(act, exp, opt...)

	if act != nil && !gg.IsErrTraced(act) {
		panic(gg.Errv(gg.JoinLines(
			`unexpected lack of stack trace in error:`,
			grepr.StringIndent(act, 1),
			gtest.MsgExtra(opt...),
		)))
	}
}

func testWrappedErrs(act, exp []error, opt ...any) {
	gtest.Len(act, len(exp))

	msgAll := gg.JoinLinesOpt(
		`all errors:`,
		grepr.StringIndent(act, 1),
		gtest.MsgExtra(opt...),
	)

	for ind, val := range act {
		if val != nil && !gg.IsErrTraced(val) {
			panic(gg.Errv(gg.JoinLinesOpt(
				gg.Str(`expected every error to be traced, found untraced at index `, ind),
				msgAll,
			)))
		}
	}

	for ind, valAct := range act {
		valExp := exp[ind]

		if valExp == nil {
			if valAct != nil {
				panic(gg.JoinLinesOpt(
					gg.Str(`unexpected non-nil error at index `, ind),
					msgAll,
				))
			}
			continue
		}

		if !e.Is(valAct, valExp) {
			panic(gg.JoinLinesOpt(
				gtest.MsgErrIsMismatch(valAct, valExp),
				msgAll,
			))
		}

		if gg.Is(valAct, valExp) {
			panic(gg.JoinLinesOpt(
				gg.Str(`unexpected unwrapped error at index `, ind),
				msgAll,
			))
		}
	}
}

func BenchmarkConcCatch_one(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		_ = gg.ConcCatch(testPanicTraced0)
	}
}

func BenchmarkConcCatch_multi(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		_ = gg.ConcCatch(gg.Nop, testPanicTraced0, gg.Nop, testPanicTraced1, gg.Nop)
	}
}

// Needs more test cases.
func TestConcMapCatch(t *testing.T) {
	defer gtest.Catch(t)

	src := []int{10, 20, 30}
	vals, errs := gg.ConcMapCatch(src, testConcMapFunc)

	gtest.Len(vals, len(src))
	gtest.Len(errs, len(src))

	gtest.Equal(vals, []string{`10`, ``, `30`})
	testWrappedErrs(errs, []error{nil, testErrTraced1, nil})
}

func testConcMapFunc(src int) string {
	if src == 20 {
		panic(testErrTraced1)
	}
	return gg.String(src)
}

// Needs more test cases.
func TestConcRaceCatch(t *testing.T) {
	defer gtest.Catch(t)

	//nolint:staticcheck
	gtest.Zero(gg.ConcRaceCatch(nil))
	//nolint:staticcheck
	gtest.Zero(gg.ConcRaceCatch(nil, nil))
	//nolint:staticcheck
	gtest.Zero(gg.ConcRaceCatch(nil, testNopCtx))

	ctx := context.Background()
	gtest.Zero(gg.ConcRaceCatch(ctx))
	gtest.Zero(gg.ConcRaceCatch(ctx, nil))
	gtest.Zero(gg.ConcRaceCatch(ctx, nil, nil))
	gtest.Zero(gg.ConcRaceCatch(ctx, nil, nil, nil))

	gtest.Zero(gg.ConcRaceCatch(ctx))
	gtest.Zero(gg.ConcRaceCatch(ctx, testNopCtx))
	gtest.Zero(gg.ConcRaceCatch(ctx, testNopCtx, testNopCtx))
	gtest.Zero(gg.ConcRaceCatch(ctx, testNopCtx, testNopCtx, testNopCtx))

	gtest.Is(
		//nolint:staticcheck
		gg.ConcRaceCatch(nil, testPanicCtxTraced0),
		testErrTraced0,
	)

	testWrappedErr(
		gg.ConcRaceCatch(ctx, testPanicCtxUntracedA),
		testErrUntracedA,
		`when only one function is provided and it panics with an untraced error, that error must be wrapped with a trace`,
	)

	gtest.Is(
		gg.ConcRaceCatch(ctx, testPanicCtxTraced0),
		testErrTraced0,
		`when only one function is provided and it panics with a traced error, that error must be returned as-is`,
	)

	testWrappedErr(
		gg.ConcRaceCatch(ctx, testPanicCtxUntracedA, testNopCtx),
		testErrUntracedA,
		`when multiple functions are provided and one panics with an untraced error, that error must be wrapped with a trace`,
	)

	testWrappedErr(
		gg.ConcRaceCatch(ctx, testPanicCtxTraced0, testNopCtx),
		testErrTraced0,
		`when multiple functions are provided and one panics with a traced error, that error must be wrapped with an additional trace`,
	)

	testWrappedErr(
		gg.ConcRaceCatch(ctx, testNopCtx, testPanicCtxTraced0),
		testErrTraced0,
	)

	testWrappedErr(
		gg.ConcRaceCatch(ctx, testNopCtx, testPanicCtxTraced0, testNopCtx),
		testErrTraced0,
	)

	testWrappedErr(
		gg.ConcRaceCatch(ctx, testPanicCtxTraced0, testPanicCtxTraced0),
		testErrTraced0,
	)

	testWrappedErr(
		gg.ConcRaceCatch(ctx, testPanicCtxTraced0, testPanicCtxTraced0, testPanicCtxTraced0),
		testErrTraced0,
	)

	// TODO: would be ideal to also verify wrapping via `testWrappedErr`.
	gtest.Has(
		[]string{testErrTraced0.Error(), testErrTraced1.Error(), testErrTraced2.Error()},
		gg.ConcRaceCatch(ctx, testPanicCtxTraced0, testPanicCtxTraced1, testPanicCtxTraced2).Error(),
	)

	/**
	Every function must receive the same cancelable context, and the context
	must be canceled after completion, regardless if we have full success,
	partial success, or full failure.
	*/
	{
		ctxRoot := context.Background()
		test := func(funs ...func(context.Context)) {
			ctxs := make([]context.Context, len(funs))

			// This test requires additional syncing because `.Run` or `.RunCatch`
			// terminate on the first panic, without waiting for the termination
			// of the remaining functions. This is by design, but in this test,
			// we must wait for their termination to ensure that the slice of
			// contexts is fully mutated.
			var gro sync.WaitGroup

			for ind, fun := range funs {
				ind, fun := ind, fun
				gro.Add(1)
				funs[ind] = func(ctxChild context.Context) {
					defer gro.Add(-1)
					ctxs[ind] = ctxChild
					fun(ctxChild)
				}
			}

			gg.Nop1(gg.ConcRaceCatch(ctxRoot, funs...))
			gro.Wait()

			testIsContextConsistent(ctxs...)

			if len(ctxs) == 1 {
				gtest.Is(ctxs[0], ctxRoot)
			} else {
				testIsCtxCanceled(ctxs...)
			}
		}

		test(testNopCtx)
		test(testNopCtx, testNopCtx)
		test(testNopCtx, testNopCtx, testNopCtx)
		test(testPanicCtxTraced0, testNopCtx, testNopCtx)
		test(testNopCtx, testPanicCtxTraced0, testNopCtx)
		test(testNopCtx, testNopCtx, testPanicCtxTraced0)
		test(testPanicCtxTraced0, testNopCtx, testPanicCtxTraced0)
		test(testPanicCtxTraced0, testPanicCtxTraced0, testPanicCtxTraced0)
	}

	/**
	On the first panic, we must immediately cancel the context before returning
	the caught error. Some of the concurrently launched functions may still
	continue running in the background. They're expected to respect context
	cancelation and terminate as soon as reasonably possible, but that's up
	to the user of the library. Our responsibility is to terminate and cancel
	as soon as the first panic is found.
	*/
	{
		var gro0 sync.WaitGroup
		var gro1 sync.WaitGroup
		var state0 CtxState
		var state1 CtxState
		var state2 CtxState

		gro0.Add(1)
		gro1.Add(2)

		testWrappedErr(
			/**
			This must terminate and return the error even though some inner functions
			are still blocked on the wait group, which is unblocked AFTER this test
			phase. This ensures that the concurrent run terminates on the first
			panic without waiting for all functions. Otherwise, the test would
			deadlock and eventually time out.
			*/
			gg.ConcRaceCatch(
				ctx,
				func(ctx context.Context) {
					defer gro1.Add(-1)
					gro0.Wait()
					state0 = ToCtxState(ctx)
				},
				func(ctx context.Context) {
					state1 = ToCtxState(ctx)
					panic(testErrTraced0)
				},
				func(ctx context.Context) {
					defer gro1.Add(-1)
					gro0.Wait()
					state2 = ToCtxState(ctx)
				},
			),
			testErrTraced0,
		)

		// This should unblock the inner functions, whose context must now be
		// canceled.
		gro0.Add(-1)
		gro1.Wait()

		gtest.Equal(state0, CtxState{true, context.Canceled})
		gtest.Equal(state1, CtxState{false, nil})
		gtest.Equal(state2, CtxState{true, context.Canceled})
	}
}

/*
Caution: this operation is prone to race conditions, and may produce
"corrupted" states, such as `{Done: false, Err: context.Canceled}`,
depending on the execution timing. Our tests ensure that we always
see very specific results, and anything else is considered a test
failure. Avoid this pattern in actual code.
*/
func ToCtxState(ctx context.Context) CtxState {
	return CtxState{isCtxDone(ctx), ctx.Err()}
}

type CtxState struct {
	Done bool
	Err  error
}

func testIsContextConsistent(vals ...context.Context) {
	if len(vals) <= 1 {
		return
	}

	exp := vals[0]
	gtest.NotZero(exp)

	for _, val := range vals {
		if exp != val {
			panic(gtest.ErrLines(
				`unexpected difference between context values`,
				gtest.MsgEqDetailed(val, exp),
			))
		}
	}
}

func testIsCtxCanceled(vals ...context.Context) {
	for ind, val := range vals {
		if !isCtxDone(val) {
			panic(gg.Errf(`expected context %v to be done`, ind))
		}
		gtest.ErrIs(val.Err(), context.Canceled)
	}
}

/*
Warning: the output is correct only when it's `true`. When the output is
`false`, then sometimes it's incorrect and the actual result is UNKNOWN
because the channel may be concurrently closed before your next line of code.

In other words, the return type of this function isn't exactly a boolean.
It's a union of "true" and "unknowable".
*/
func isCtxDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
