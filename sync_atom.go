package gg

import (
	"sync/atomic"
)

/*
Typedef of `atomic.Pointer` with added methods which operate on values rather
than pointers, which is often more convenient.
*/
type Atom[A any] atomic.Pointer[A]

/*
Returns the value behind the latest stored pointer, if any, or a zero value,
along with a boolean indicating if a pointer was actually stored.
*/
func (self *Atom[A]) LoadedVal() (_ A, _ bool) {
	val := self.LoadPtr()
	if val == nil {
		return
	}
	return *val, true
}

// Returns the value behind the latest stored pointer, if any, or a zero value.
func (self *Atom[A]) LoadVal() A {
	val, _ := self.LoadedVal()
	return val
}

// Allocates a copy of the provided value and stores a pointer to that memory.
func (self *Atom[A]) StoreVal(val A) {
	(*atomic.Pointer[A])(self).Store(&val)
}

// Same as `atomic.Pointer.Load`.
func (self *Atom[A]) LoadPtr() *A { return (*atomic.Pointer[A])(self).Load() }

// Same as `atomic.Pointer.Store`.
func (self *Atom[A]) StorePtr(val *A) {
	(*atomic.Pointer[A])(self).Store(val)
}

// Same as `atomic.Pointer.Swap`.
func (self *Atom[A]) SwapPtr(val *A) *A {
	return (*atomic.Pointer[A])(self).Swap(val)
}

// Same as `atomic.Pointer.CompareAndSwap`.
func (self *Atom[A]) CompareAndSwapPtr(prev, next *A) bool {
	return (*atomic.Pointer[A])(self).CompareAndSwap(prev, next)
}

// Replaces any currently stored pointer with nil.
func (self *Atom[A]) Clear() { (*atomic.Pointer[A])(self).Store(nil) }
