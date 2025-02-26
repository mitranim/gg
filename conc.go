package gg

import (
	"context"
	"sync"
)

/*
Runs the given functions concurrently via `ConcCatch`, waiting for all to finish
and collecting their panics as `error` values. After all functions are finished,
if there is at least one non-nil error, panics with the combined error.
*/
func Conc(funs ...func()) {
	TryErr(ErrMul(ConcCatch(funs...)...))
}

/*
Concurrently runs the given functions. Catches their panics, converts them to
`error`, and returns the resulting errors. The error slice always has the same
length as the number of given functions.

Ensures that the resulting errors have stack traces, combining the traces from
goroutines spawned for the given functions (if necessary) with the trace from
the calling goroutine.

If `error` is desired instead of `[]error`, use `ErrMul` to convert correctly.
*/
func ConcCatch(funs ...func()) []error {
	switch len(funs) {
	case 0:
		return nil

	case 1:
		return []error{Catch(funs[0])}

	default:
		errs := make(Errs, len(funs))
		var gro sync.WaitGroup

		for ind, fun := range funs {
			if fun == nil {
				continue
			}
			gro.Add(1)
			go concCatchRun(&gro, &errs[ind], fun)
		}

		gro.Wait()
		errs.WrapTracedAt(1)
		return errs
	}
}

func concCatchRun(gro *sync.WaitGroup, errPtr *error, fun func()) {
	defer gro.Add(-1)
	defer Rec(errPtr)
	fun()
}

/*
Concurrently calls the given function on each element of the given slice. If the
function is nil, does nothing. Also see `ConcEachCatch`, `ConcMap`, `Conc`.
*/
func ConcEach[A any](src []A, fun func(A)) {
	TryErr(ErrMul(ConcEachCatch(src, fun)...))
}

/*
Concurrently calls the given function on each element of the given slice.
Collects panics from each call as error values. Also see `ConcEach`, `ConcMap`,
`Conc`.

Ensures that the resulting errors have stack traces, combining the traces from
goroutines spawned for the given functions (if necessary) with the trace from
the calling goroutine.

If `error` is desired instead of `[]error`, use `ErrMul` to convert correctly.
*/
func ConcEachCatch[A any](src []A, fun func(A)) []error {
	if fun == nil {
		return nil
	}

	switch len(src) {
	case 0:
		return nil

	case 1:
		return []error{Catch10(fun, src[0])}

	default:
		errs := concEachCatch(src, fun)
		Errs(errs).WrapTracedAt(1)
		return errs
	}
}

func concEachCatch[A any](src []A, fun func(A)) []error {
	errs := make([]error, len(src))
	var gro sync.WaitGroup

	for ind, val := range src {
		gro.Add(1)
		go concCatchEachRun(&gro, &errs[ind], fun, val)
	}

	gro.Wait()
	return errs
}

func concCatchEachRun[A any](gro *sync.WaitGroup, errPtr *error, fun func(A), val A) {
	defer gro.Add(-1)
	defer Rec(errPtr)
	fun(val)
}

// Like `Map` but concurrent. Also see `ConcMapCatch`, `ConcEach`, `Conc`.
func ConcMap[A, B any](src []A, fun func(A) B) []B {
	vals, errs := ConcMapCatch(src, fun)
	TryMul(errs...)
	return vals
}

/*
Like `Map` but concurrent. Returns the resulting values along with the caught
panics, if any. Also see `ConcMap`, `ConcCatch`, `ConcEachCatch`.

Ensures that the resulting errors have stack traces, combining the traces from
goroutines spawned for the given functions (if necessary) with the trace from
the calling goroutine.

If `error` is desired instead of `[]error`, use `ErrMul` to convert correctly.
*/
func ConcMapCatch[A, B any](src []A, fun func(A) B) ([]B, []error) {
	if fun == nil {
		return nil, nil
	}

	switch len(src) {
	case 0:
		return nil, nil

	case 1:
		val, err := Catch11(fun, src[0])
		return []B{val}, []error{err}

	default:
		out, errs := concMapCatch(src, fun)
		Errs(errs).WrapTracedAt(1)
		return out, errs
	}
}

func concMapCatch[A, B any](src []A, fun func(A) B) ([]B, []error) {
	vals := make([]B, len(src))
	errs := make([]error, len(src))
	var gro sync.WaitGroup

	for ind, val := range src {
		gro.Add(1)
		go concCatchMapRun(&gro, &vals[ind], &errs[ind], fun, val)
	}

	gro.Wait()
	return vals, errs
}

func concCatchMapRun[A, B any](gro *sync.WaitGroup, tar *B, errPtr *error, fun func(A) B, val A) {
	defer gro.Add(-1)
	defer Rec(errPtr)
	*tar = fun(val)
}

/*
Partial application / thunk of `ConcMap`, suitable for passing to `Conc`.
When the resulting function is fully executed, it replaces the referenced
slice with the accumulated result.
*/
func ConcMapFunc[A, B any](tar *[]B, src []A, fun func(A) B) func() {
	if IsEmpty(src) || fun == nil {
		return nil
	}
	return func() { *tar = ConcMap(src, fun) }
}

/*
Variant of `Conc` with support for context and cancelation. Runs the given
functions via `ConcRaceCatch`. If there was an error, panics with that error.
*/
func ConcRace(ctx context.Context, funs ...func(context.Context)) {
	TryErr(ConcRaceCatch(ctx, funs...))
}

/*
Runs the functions concurrently. Blocks until all functions complete
successfully, returning nil. If one of the functions panics, cancels the
context passed to each function, and immediately returns the resulting error,
without waiting for the other functions to terminate. In this case, the panics
in other functions, if any, are caught and ignored.

Ensures that the resulting error, if any, has a stack trace, combining the trace
from the goroutine spawned for the panicked function (if necessary) with the
trace from the calling goroutine.
*/
func ConcRaceCatch(ctx context.Context, funs ...func(context.Context)) (err error) {
	switch len(funs) {
	case 0:
		return nil

	case 1:
		fun := funs[0]
		if fun == nil {
			return nil
		}

		defer Rec(&err)
		fun(ctx)
		return

	default:
		var gro sync.WaitGroup
		var errChan chan error

		for _, fun := range funs {
			if fun == nil {
				continue
			}

			if errChan == nil {
				errChan = make(chan error, 1)

				var cancel func()
				ctx, cancel = context.WithCancel(ctx)

				// Note: unlike `defer` in some other languages, `defer` in Go is
				// function-scoped, not block-scoped. This will be executed once
				// we're done waiting on the error channel.
				defer cancel()
			}

			gro.Add(1)
			go runConcCtx(&gro, errChan, fun, ctx)
		}

		// Happens when len > 1 but every func is nil.
		if errChan == nil {
			return nil
		}

		go closeConcCtx(&gro, errChan)
		return WrapTracedAt(<-errChan, 1)
	}
}

func runConcCtx(gro *sync.WaitGroup, errChan chan error, fun func(context.Context), ctx context.Context) {
	defer gro.Add(-1)
	defer recSend(errChan)
	fun(ctx)
}

func closeConcCtx(gro *sync.WaitGroup, errChan chan error) {
	defer close(errChan)
	gro.Wait()
}

func recSend(tar chan error) {
	err := AnyErrTracedAt(recover(), 1)
	if err != nil {
		SendOpt(tar, err)
	}
}
