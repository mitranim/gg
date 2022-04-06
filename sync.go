package gg

import (
	"sync"
	"sync/atomic"
)

/*
Shortcut for mutexes. Usage:

	defer Locked(someLock).Unlock()
*/
func Locked[A interface{ Lock() }](val A) A {
	val.Lock()
	return val
}

/*
Typed version of `atomic.Value`. Currently implemented as a typedef of
`atomic.Value` where the value is internally stored as `any`, which may cause
the value to be automatically copied when stored. Thus, large values should be
stored by pointer, unless copying is desirable. This may change in the future.
*/
type Atom[A any] atomic.Value

/*
Like `.Load` but returns true if anything was previously stored, and false if
nothing was previously stored.
*/
func (self *Atom[A]) Loaded() (A, bool) {
	val := (*atomic.Value)(self).Load()
	return AnyTo[A](val), val != nil
}

// Typed version of `atomic.Value.Load`.
func (self *Atom[A]) Load() A {
	return AnyTo[A]((*atomic.Value)(self).Load())
}

// Typed version of `atomic.Value.Store`.
func (self *Atom[A]) Store(val A) {
	(*atomic.Value)(self).Store(val)
}

// Typed version of `atomic.Value.Swap`.
func (self *Atom[A]) Swap(val A) A {
	return AnyTo[A]((*atomic.Value)(self).Swap(val))
}

// Typed version of `atomic.Value.CompareAndSwap`.
func (self *Atom[A]) CompareAndSwap(prev, next A) bool {
	return (*atomic.Value)(self).CompareAndSwap(prev, next)
}

/*
Typed version of `sync.Map`. Currently implemented as a typedef of `sync.Map`
where both keys and valus are internally stored as `any`, which may cause them
to be automatically copied when stored. Thus, large values should be stored by
pointer, unless copying is desirable. This may change in the future.
*/
type SyncMap[Key comparable, Val any] sync.Map

// Typed version of `sync.Map.Load`.
func (self *SyncMap[Key, Val]) Load(key Key) (Val, bool) {
	iface, ok := (*sync.Map)(self).Load(key)
	return AnyTo[Val](iface), ok
}

// Typed version of `sync.Map.LoadOrStore`.
func (self *SyncMap[Key, Val]) LoadOrStore(key Key, val Val) (Val, bool) {
	iface, ok := (*sync.Map)(self).LoadOrStore(key, val)
	return AnyTo[Val](iface), ok
}

// Typed version of `sync.Map.LoadAndDelete`.
func (self *SyncMap[Key, Val]) LoadAndDelete(key Key) (Val, bool) {
	iface, ok := (*sync.Map)(self).LoadAndDelete(key)
	return AnyTo[Val](iface), ok
}

// Typed version of `sync.Map.Store`.
func (self *SyncMap[Key, Val]) Store(key Key, val Val) Val {
	(*sync.Map)(self).Store(key, val)
	return val
}

// Typed version of `sync.Map.Delete`.
func (self *SyncMap[Key, Val]) Delete(key Key) {
	(*sync.Map)(self).Delete(key)
}

// Typed version of `sync.Map.Range`.
func (self *SyncMap[Key, Val]) Range(fun func(Key, Val) bool) {
	if fun == nil {
		return
	}
	(*sync.Map)(self).Range(func(key, val any) bool {
		return fun(key.(Key), val.(Val))
	})
}
