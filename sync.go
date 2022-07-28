package gg

import (
	"sync"
	"sync/atomic"
)

/*
Shortcut for mutexes. Usage:

	defer Lock(someLock).Unlock()
*/
func Lock(val sync.Locker) sync.Locker {
	val.Lock()
	return val
}

/*
Shortcut for dereferencing a pointer under a lock. Uses `PtrGet`, returning the
zero value of the given type if the pointer is nil.
*/
func LockGet[A any](lock sync.Locker, ptr *A) A {
	if ptr == nil {
		return Zero[A]()
	}
	defer Lock(lock).Unlock()
	return *ptr
}

// Shortcut for writing to a pointer under a lock.
func LockSet[A any](lock sync.Locker, ptr *A, val A) {
	if ptr == nil {
		return
	}
	defer Lock(lock).Unlock()
	*ptr = val
}

/*
Typed version of `atomic.Value`. Currently implemented as a typedef of
`atomic.Value` where the value is internally stored as `any`, which may cause
the value to be automatically copied when stored. Thus, large values should be
stored by pointer to minimize copying. This may change in the future.
*/
type Atom[A any] atomic.Value

/*
Like `.Load` but returns true if anything was previously stored, and false if
nothing was previously stored.
*/
func (self *Atom[A]) Loaded() (A, bool) {
	val := (*atomic.Value)(self).Load()
	return AnyAs[A](val), val != nil
}

// Typed version of `atomic.Value.Load`.
func (self *Atom[A]) Load() A {
	return AnyAs[A]((*atomic.Value)(self).Load())
}

// Typed version of `atomic.Value.Store`.
func (self *Atom[A]) Store(val A) {
	(*atomic.Value)(self).Store(val)
}

// Typed version of `atomic.Value.Swap`.
func (self *Atom[A]) Swap(val A) A {
	return AnyAs[A]((*atomic.Value)(self).Swap(val))
}

// Typed version of `atomic.Value.CompareAndSwap`.
func (self *Atom[A]) CompareAndSwap(prev, next A) bool {
	return (*atomic.Value)(self).CompareAndSwap(prev, next)
}

/*
Typed version of `sync.Map`. Currently implemented as a typedef of `sync.Map`
where both keys and valus are internally stored as `any`, which may cause them
to be automatically copied when stored. Thus, large values should be stored by
pointer to minimize copying. This may change in the future.
*/
type SyncMap[Key comparable, Val any] sync.Map

// Typed version of `sync.Map.Load`.
func (self *SyncMap[Key, Val]) Load(key Key) (Val, bool) {
	iface, ok := (*sync.Map)(self).Load(key)
	return AnyAs[Val](iface), ok
}

// Typed version of `sync.Map.LoadOrStore`.
func (self *SyncMap[Key, Val]) LoadOrStore(key Key, val Val) (Val, bool) {
	iface, ok := (*sync.Map)(self).LoadOrStore(key, val)
	return AnyAs[Val](iface), ok
}

// Typed version of `sync.Map.LoadAndDelete`.
func (self *SyncMap[Key, Val]) LoadAndDelete(key Key) (Val, bool) {
	iface, ok := (*sync.Map)(self).LoadAndDelete(key)
	return AnyAs[Val](iface), ok
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

// Alias of `chan` with additional convenience methods.
type Chan[A any] chan A

// Closes the channel unless it's nil.
func (self Chan[_]) Close() {
	if self != nil {
		close(self)
	}
}

// Same as global `ChanInit`.
//go:noinline
func (self *Chan[A]) Init() Chan[A] { return ChanInit(self) }

/*
Idempotently initializes the channel. If the pointer is non-nil and the channel
is nil, creates a new unbuffered channel and assigns it to the pointer. Returns
the resulting channel.
*/
func ChanInit[Tar ~chan Val, Val any](ptr *Tar) Tar {
	if ptr == nil {
		return nil
	}
	if *ptr == nil {
		*ptr = make(Tar)
	}
	return *ptr
}

// Same as global `ChanInitCap`.
//go:noinline
func (self *Chan[A]) InitCap(cap int) Chan[A] { return ChanInitCap(self, cap) }

/*
Idempotently initializes the channel. If the pointer is non-nil and the channel
is nil, creates a new buffered channel with the given capacity and assigns it
to the pointer. Returns the resulting channel.
*/
func ChanInitCap[Tar ~chan Val, Val any](ptr *Tar, cap int) Tar {
	if ptr == nil {
		return nil
	}
	if *ptr == nil {
		*ptr = make(Tar, cap)
	}
	return *ptr
}

// Same as global `SendOpt`.
//go:noinline
func (self Chan[A]) SendOpt(val A) { SendOpt(self, val) }

// Shortcut for sending a value over a channel in a non-blocking fashion.
func SendOpt[Tar ~chan Val, Val any](tar Tar, val Val) {
	select {
	case tar <- val:
	default:
	}
}

// Same as global `SendZeroOpt`.
//go:noinline
func (self Chan[A]) SendZeroOpt() { SendZeroOpt(self) }

// Shortcut for sending a zero value over a channel in a non-blocking fashion.
func SendZeroOpt[Tar ~chan Val, Val any](tar Tar) {
	select {
	case tar <- Zero[Val]():
	default:
	}
}
