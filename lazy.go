package gg

import (
	"sync"
	"sync/atomic"
)

/*
Very similar to [sync.OnceValue], with the difference in panic handling.
[sync.OnceValue] stores a panic and re-panics on subsequent calls, while
our [Once] retries on subsequent calls. Calls do not overlap.
*/
func Once[A any](fun func() A) func() A {
	if fun == nil {
		return Zero[A]
	}

	/**
	We're not using [sync.Once] because its panic handling is also incompatible
	with ours, but in a different way from [sync.OnceValue]. [sync.Once] always
	considers the first call successful. If the first call panics and does not
	crash the process due to the caller recovering, the function never reruns.
	In our case, that would result in always returning a zero value.

	Minor perf notes.

	[sync.OnceValue] groups the closured state into a struct to ensure a single
	heap allocation, but in Go 1.24, our microbenchmark shows a single allocation
	and marginally lower memory usage with regular inline variables.

	[sync.OnceValue] shows a marginally lower cost per call than [Once] in our
	microbenchmark, but when we exactly copy its implementation, and benchmark
	it side by side with the original, the copy performs worse. Special compiler
	treatment for stdlib? In any case, the difference is tiny and the overheads
	are absolutely negligible.
	*/
	var done atomic.Bool
	var lock sync.Mutex
	var val A

	create := func() {
		defer Lock(&lock).Unlock()
		if done.Load() {
			return
		}
		val = fun()
		fun = nil
		done.Store(true)
	}

	return func() A {
		if !done.Load() {
			create()
		}
		return val
	}
}

/*
Creates [Lazy] with the given function. Deprecated in favor of [Once],
which is easier to use.
*/
func NewLazy[A any](fun func() A) *Lazy[A] { return &Lazy[A]{fun: fun} }

// Deprecated in favor of [Once], which is easier to use.
type Lazy[A any] struct {
	val  A
	fun  func() A
	once sync.Once
}

// Returns the inner value after idempotently creating it.
func (self *Lazy[A]) Get() A {
	self.once.Do(self.init)
	return self.val
}

func (self *Lazy[_]) init() {
	fun := self.fun
	self.fun = nil
	if fun != nil {
		self.val = fun()
	}
}
