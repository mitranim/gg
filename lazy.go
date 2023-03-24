package gg

import "sync"

/*
Creates `Lazy` with the given function. See the type's description for details.
*/
func NewLazy[A any](fun func() A) *Lazy[A] { return &Lazy[A]{fun: fun} }

/*
Similar to `sync.Once`, but specialized for creating and caching one value,
instead of relying on nullary functions and side effects. Created via `NewLazy`.
Calling `.Get` on the resulting object will idempotently call the given function
and cache the result, and discard the function. Uses `sync.Once` internally for
synchronization.
*/
type Lazy[A any] struct {
	once sync.Once
	val  A
	fun  func() A
}

/*
Returns the inner value after idempotently creating it.
This method is synchronized and safe for concurrent use.
*/
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
