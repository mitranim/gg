package gg_test

import (
	"context"
	"sync"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

var (
	testErr0 = error(gg.Errf(`test err 0`))
	testErr1 = error(gg.Errf(`test err 1`))
	testErr2 = error(gg.Errf(`test err 2`))
)

const (
	testErrA = gg.ErrStr(`test err A`)
	testErrB = gg.ErrStr(`test err B`)
)

func testPanic0() { panic(testErr0) }
func testPanic1() { panic(testErr1) }

func testNopCtx(context.Context)    {}
func testPanicCtx0(context.Context) { panic(testErr0) }
func testPanicCtx1(context.Context) { panic(testErr1) }
func testPanicCtx2(context.Context) { panic(testErr2) }

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

		gtest.Equal(
			gg.ConcCatch(testPanic0),
			[]error{testErr0},
		)

		gtest.Equal(
			gg.ConcCatch(testPanic0, testPanic1),
			[]error{testErr0, testErr1},
		)
	})

	t.Run(`mixed`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(
			gg.ConcCatch(gg.Nop, testPanic0, gg.Nop, testPanic1, gg.Nop),
			gg.Errs{nil, testErr0, nil, testErr1, nil},
		)
	})
}

func BenchmarkConcCatch_one(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		_ = gg.ConcCatch(testPanic0)
	}
}

func BenchmarkConcCatch_multi(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		_ = gg.ConcCatch(gg.Nop, testPanic0, gg.Nop, testPanic1, gg.Nop)
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
	gtest.Equal(errs, []error{nil, testErr1, nil})
}

func testConcMapFunc(src int) string {
	if src == 20 {
		panic(testErr1)
	}
	return gg.String(src)
}

// Needs more test cases.
func TestConcRace(t *testing.T) {
	defer gtest.Catch(t)

	//nolint:staticcheck
	gtest.Zero(gg.ConcRace().RunCatch(nil))
	//nolint:staticcheck
	gtest.Zero(gg.ConcRace(nil).RunCatch(nil))
	//nolint:staticcheck
	gtest.Zero(gg.ConcRace(testNopCtx).RunCatch(nil))

	ctx := context.Background()
	gtest.Zero(gg.ConcRace().RunCatch(ctx))
	gtest.Zero(gg.ConcRace(nil).RunCatch(ctx))
	gtest.Zero(gg.ConcRace(nil, nil).RunCatch(ctx))
	gtest.Zero(gg.ConcRace(nil, nil, nil).RunCatch(ctx))

	gtest.Zero(gg.ConcRace(testNopCtx).RunCatch(ctx))
	gtest.Zero(gg.ConcRace(testNopCtx, testNopCtx).RunCatch(ctx))
	gtest.Zero(gg.ConcRace(testNopCtx, testNopCtx, testNopCtx).RunCatch(ctx))

	gtest.Is(
		//nolint:staticcheck
		gg.ConcRace(testPanicCtx0).RunCatch(nil),
		testErr0,
	)

	gtest.Is(
		gg.ConcRace(testPanicCtx0, testNopCtx).RunCatch(ctx),
		testErr0,
	)

	gtest.Is(
		gg.ConcRace(testNopCtx, testPanicCtx0).RunCatch(ctx),
		testErr0,
	)

	gtest.Is(
		gg.ConcRace(testNopCtx, testPanicCtx0, testNopCtx).RunCatch(ctx),
		testErr0,
	)

	gtest.Is(
		gg.ConcRace(testPanicCtx0, testPanicCtx0).RunCatch(ctx),
		testErr0,
	)

	gtest.Is(
		gg.ConcRace(testPanicCtx0, testPanicCtx0, testPanicCtx0).RunCatch(ctx),
		testErr0,
	)

	gtest.HasEqual(
		[]error{testErr0, testErr1, testErr2},
		gg.ConcRace(testPanicCtx0, testPanicCtx1, testPanicCtx2).RunCatch(ctx),
	)

	/**
	Every function must receive the same cancelable context, and the context
	must be canceled after completion, regardless if we have full success,
	partial success, or full failure.
	*/
	{
		test := func(funs ...func(context.Context)) {
			conc := make(gg.ConcRaceSlice, len(funs))
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

				conc.Add(func(ctx context.Context) {
					defer gro.Add(-1)
					ctxs[ind] = ctx
					fun(ctx)
				})
			}

			gg.Nop1(conc.RunCatch(ctx))
			gro.Wait()

			testIsContextConsistent(ctxs...)
			testIsCtxCanceled(ctxs[0])
		}

		test(testNopCtx)
		test(testNopCtx, testNopCtx)
		test(testNopCtx, testNopCtx, testNopCtx)
		test(testPanicCtx0, testNopCtx, testNopCtx)
		test(testNopCtx, testPanicCtx0, testNopCtx)
		test(testNopCtx, testNopCtx, testPanicCtx0)
		test(testPanicCtx0, testNopCtx, testPanicCtx0)
		test(testPanicCtx0, testPanicCtx0, testPanicCtx0)
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

		gtest.Is(
			/**
			This must terminate and return the error even though some inner functions
			are still blocked on the wait group, which is unblocked AFTER this test
			phase. This ensures that the concurrent run terminates on the first
			panic without waiting for all functions. Otherwise, the test would
			deadlock and eventually time out.
			*/
			gg.ConcRace(
				func(ctx context.Context) {
					defer gro1.Add(-1)
					gro0.Wait()
					state0 = ToCtxState(ctx)
				},
				func(ctx context.Context) {
					state1 = ToCtxState(ctx)
					panic(testErr0)
				},
				func(ctx context.Context) {
					defer gro1.Add(-1)
					gro0.Wait()
					state2 = ToCtxState(ctx)
				},
			).RunCatch(ctx),
			testErr0,
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
	if !(len(vals) > 1) {
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

func testIsCtxCanceled(ctx context.Context) {
	if !isCtxDone(ctx) {
		panic(`expected context to be done`)
	}
	gtest.ErrorIs(ctx.Err(), context.Canceled)
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
