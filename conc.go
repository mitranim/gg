package gg

import (
	"context"
	"sync"
)

/*
Tiny shortcut for gradually building a list of funcs which are later to be
executed concurrently. This type's methods invoke global funcs such as `Conc`.
Compare `ConcRaceSlice`.
*/
type ConcSlice []func()

// If the given func is non-nil, adds it to the slice for later execution.
func (self *ConcSlice) Add(fun func()) {
	if fun != nil {
		*self = append(*self, fun)
	}
}

// Same as calling `Conc` with the given slice of funcs.
func (self ConcSlice) Run() { Conc(self...) }

// Same as calling `ConcCatch` with the given slice of funcs.
func (self ConcSlice) RunCatch() []error { return ConcCatch(self...) }

/*
Runs the given funcs sequentially rather than concurrently. Provided for
performance debugging.
*/
func (self ConcSlice) RunSeq() {
	for _, fun := range self {
		if fun != nil {
			fun()
		}
	}
}

/*
Shortcut for concurrent execution. Runs the given functions via `ConcCatch`.
If there is at least one error, panics with the combined error, adding a stack
trace pointing to the call site of `Conc`.
*/
func Conc(val ...func()) {
	TryAt(ErrMul(ConcCatch(val...)...), 1)
}

/*
Concurrently runs the given functions. Catches their panics, converts them to
`error`, and returns the resulting errors. The error slice always has the same
length as the number of given functions.

Ensures that the resulting errors have stack traces, both on the current
goroutine, and on any sub-goroutines used for concurrent execution, wrapping
errors as necessary.

The error slice can be converted to `error` via `ErrMul` or `Errs.Err`. The
slice is returned as `[]error` rather than `Errs` to avoid accidental incorrect
conversion of empty `Errs` to non-nil `error`.
*/
func ConcCatch(val ...func()) []error {
	switch len(val) {
	case 0:
		return nil

	case 1:
		return []error{Catch(val[0])}

	default:
		tar := concCatch(val)
		Errs(tar).WrapTracedAt(1)
		return tar
	}
}

func concCatch(src []func()) []error {
	tar := make([]error, len(src))
	var gro sync.WaitGroup

	for ind, fun := range src {
		if fun == nil {
			continue
		}
		gro.Add(1)
		go concCatchRun(&gro, &tar[ind], fun)
	}

	gro.Wait()
	return tar
}

func concCatchRun(gro *sync.WaitGroup, errPtr *error, fun func()) {
	defer gro.Add(-1)
	defer Rec(errPtr)
	fun()
}

/*
Concurrently calls the given function on each element of the given slice. If the
function is nil, does nothing. Also see `Conc`.
*/
func ConcEach[A any](src []A, fun func(A)) {
	TryAt(ErrMul(ConcEachCatch(src, fun)...), 1)
}

/*
Concurrently calls the given function on each element of the given slice,
returning the resulting panics if any. If the function is nil, does nothing and
returns nil. Also see `ConcCatch`

Ensures that the resulting errors have stack traces, both on the current
goroutine, and on any sub-goroutines used for concurrent execution, wrapping
errors as necessary.

The error slice can be converted to `error` via `ErrMul` or `Errs.Err`. The
slice is returned as `[]error` rather than `Errs` to avoid accidental incorrect
conversion of empty `Errs` to non-nil `error`.
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
		tar := concEachCatch(src, fun)
		Errs(tar).WrapTracedAt(1)
		return tar
	}
}

func concEachCatch[A any](src []A, fun func(A)) []error {
	tar := make([]error, len(src))
	var gro sync.WaitGroup

	for ind, val := range src {
		gro.Add(1)
		go concCatchEachRun(&gro, &tar[ind], fun, val)
	}

	gro.Wait()
	return tar
}

func concCatchEachRun[A any](gro *sync.WaitGroup, errPtr *error, fun func(A), val A) {
	defer gro.Add(-1)
	defer Rec(errPtr)
	fun(val)
}

// Like `Map` but concurrent. Also see `Conc`.
func ConcMap[A, B any](src []A, fun func(A) B) []B {
	vals, errs := ConcMapCatch(src, fun)
	TryMul(errs...)
	return vals
}

/*
Like `Map` but concurrent. Returns the resulting values along with the caught
panics, if any. Also see `ConcCatch`.

Ensures that the resulting errors have stack traces, both on the current
goroutine, and on any sub-goroutines used for concurrent execution, wrapping
errors as necessary.

The error slice can be converted to `error` via `ErrMul` or `Errs.Err`. The
slice is returned as `[]error` rather than `Errs` to avoid accidental incorrect
conversion of empty `Errs` to non-nil `error`.
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

// Partial application / thunk of `ConcMap`, suitable for `Conc`.
func ConcMapFunc[A, B any](tar *[]B, src []A, fun func(A) B) func() {
	if IsEmpty(src) || fun == nil {
		return nil
	}
	return func() { *tar = ConcMap(src, fun) }
}

/*
Shortcut for constructing `ConcRaceSlice` in a variadic call with parens rather
than braces.
*/
func ConcRace(src ...func(context.Context)) ConcRaceSlice {
	return ConcRaceSlice(src)
}

/*
Tool for concurrent execution. Similar to `ConcSlice`, but with support for
context and cancelation. See `ConcRaceSlice.RunCatch` for details.
*/
type ConcRaceSlice []func(context.Context)

// If the given func is non-nil, adds it to the slice for later execution.
func (self *ConcRaceSlice) Add(fun func(context.Context)) {
	if fun != nil {
		*self = append(*self, fun)
	}
}

/*
Shortcut. Runs the functions via `ConcRaceSlice.RunCatch`. If the resulting
error is non-nil, panics with that error, idempotently adding a stack trace.
*/
func (self ConcRaceSlice) Run(ctx context.Context) {
	TryAt(self.RunCatch(ctx), 1)
}

/*
Runs the functions concurrently. Blocks until all functions complete
successfully, returning nil. If one of the functions panics, cancels the
context passed to each function, and immediately returns the resulting error,
without waiting for the other functions to terminate. In this case, the panics
in other functions, if any, are caught and ignored.

Ensures that the resulting error has a stack trace, both on the current
goroutine, and on the sub-goroutine used for concurrent execution, wrapping
errors as necessary.
*/
func (self ConcRaceSlice) RunCatch(ctx context.Context) (err error) {
	switch len(self) {
	case 0:
		return nil

	case 1:
		fun := self[0]
		if fun == nil {
			return nil
		}

		defer Rec(&err)
		fun(ctx)
		return

	default:
		return WrapTracedAt(self.run(ctx), 1)
	}
}

func (self ConcRaceSlice) run(ctx context.Context) error {
	var gro sync.WaitGroup
	var errChan chan error

	for _, fun := range self {
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

	// Happens when every element is nil and len >= 2.
	if errChan == nil {
		return nil
	}

	go closeConcCtx(&gro, errChan)
	return <-errChan
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

func recSend(errChan chan error) {
	err := AnyErrTracedAt(recover(), 1)
	if err != nil {
		SendOpt(errChan, err)
	}
}
