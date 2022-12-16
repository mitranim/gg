package gg

import "sync"

/**
Tiny shortcut for gradually building a list of funcs which are later to be
executed concurrently. This type's methods invoke global funcs such as `Conc`.
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

/**
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

// Concurrently runs the given functions.
func Conc(val ...func()) { Errs(ConcCatch(val...)).Try() }

// Concurrently runs the given functions, returning their panics.
func ConcCatch(val ...func()) []error {
	switch len(val) {
	case 0:
		return nil

	case 1:
		return []error{Catch(val[0])}

	default:
		return concCatch(val)
	}
}

func concCatch(src []func()) []error {
	out := make([]error, 0, len(src))
	var gro sync.WaitGroup

	for _, fun := range src {
		if fun == nil {
			AppendVals(&out, nil)
		} else {
			gro.Add(1)
			go concCatchRun(&gro, fun, AppendPtrZero(&out))
		}
	}

	gro.Wait()
	return out
}

func concCatchRun(gro *sync.WaitGroup, fun func(), err *error) {
	defer gro.Add(-1)
	defer Rec(err)
	fun()
}

// Concurrently calls the given function on each element of the given slice.
func ConcEach[A any](src []A, fun func(A)) { Errs(ConcEachCatch(src, fun)).Try() }

/*
Concurrently calls the given function on each element of the given slice,
returning the resulting panics if any.
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
		return concEachCatch(src, fun)
	}
}

func concEachCatch[A any](src []A, fun func(A)) []error {
	out := make([]error, 0, len(src))
	var gro sync.WaitGroup

	for _, val := range src {
		gro.Add(1)
		go concCatchEachRun(&gro, fun, val, AppendPtrZero(&out))
	}

	gro.Wait()
	return out
}

func concCatchEachRun[A any](gro *sync.WaitGroup, fun func(A), val A, err *error) {
	defer gro.Add(-1)
	defer Rec(err)
	fun(val)
}

// Like `Map` but concurrent.
func ConcMap[A, B any](src []A, fun func(A) B) []B {
	vals, errs := ConcMapCatch(src, fun)
	Errs(errs).Try()
	return vals
}

/*
Like `Map` but concurrent. Returns the resulting values along with the caught
panics, if any.
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
		return concMapCatch(src, fun)
	}
}

func concMapCatch[A, B any](src []A, fun func(A) B) ([]B, []error) {
	vals := make([]B, 0, len(src))
	errs := make([]error, 0, len(src))
	var gro sync.WaitGroup

	for _, val := range src {
		gro.Add(1)
		go concCatchMapRun(&gro, fun, val, AppendPtrZero(&vals), AppendPtrZero(&errs))
	}

	gro.Wait()
	return vals, errs
}

func concCatchMapRun[A, B any](gro *sync.WaitGroup, fun func(A) B, val A, out *B, err *error) {
	defer gro.Add(-1)
	defer Rec(err)
	*out = fun(val)
}

// Partial application / thunk of `ConcMap`, suitable for `Conc`.
func ConcMapFunc[A, B any](out *[]B, src []A, fun func(A) B) func() {
	if len(src) == 0 || fun == nil {
		return nil
	}
	return func() { *out = ConcMap(src, fun) }
}
