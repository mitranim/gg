package gg

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

/*
Constructs a [DynVar]: a dynamically-scoped variable.

When [DynVar.Get] is called on goroutines where the variable's value is not set,
the default function provided here, if any, is used to create the default
value, which is stored in the dynamic variable and reused in other similar
cases.

If the default function panics, and the process doesn't crash, future calls to
[DynVar.Get] will call it again, until it succeeds. After the first success,
the function is removed.

Default function calls within one [DynVar] are serialized and never overlap.
*/
func NewDynVar[A any](def func() A) *DynVar[A] { return &DynVar[A]{def: def} }

/*
Represents a dynamically-scoped variable. Building block for goroutine-local
storage (GLS), similar to thread-local storage (TLS) but adapted for Go.

Uses [GidFunc], if set, to determine the GID, falling back on [Gid].

The authors of Go have always objected to exposing goroutine IDs or providing
any form of GLS, but Go itself uses TLS internally, and for testing they import
an external library for getting goroutine IDs. Public access was requested
many, many times, and there are other libraries which provide GLS support.
Clearly, use cases exist.

Worth noting that many major languages support various types of contextual
storage which, unlike regular TLS, is automatically inherited, which is what
you ACTUALLY want. Tends to be tied to async/await and called something like
"async context", but sometimes TLS is supported too. Examples:

  - JS: `node:async_hooks` and TC39 proposal for async context;
    only async/await, and doesn't propagate to workers.
  - Python: `contextvars` (only async/await).
  - Rust: `tokio::task_local!` (only async/await).
  - Java: `InheritableThreadLocal` and `ScopedValue`.
  - C#: `AsyncLocal` (threads _and_ async/await).

The major gotcha of naive GLS / TLS is the lack of implicit inheritance.
Non-inherited storage is usually NOT what you want. To inherit GLS, spawn
goroutines with [GlsGo] or [GlsGo1] rather than the `go` keyword. When running
code in goroutines whose spawning you don't control, copy the GLS from the
parent to the child goroutine via [GlsSnap] and [GlsSet].

Performance: on architectures and in Go versions supported by the "fast path"
of [Gid], [DynVar] operations are fairly cheap. The overhead of [DynVar.Get]
is mostly one [sync.Map.Load]; in Go 1.24 on M3 Pro, it clocks at around 10ns
in a single-goroutine benchmark, though costs may vary under contention.
[DynVar.Set] is similar, but also involves a conversion of the input to `any`
which makes a heap copy of any value wider than a machine word; very large
objects should be passed by pointer, just like everywhere else in Go.

A zero value is ready to use. [DynVar] contains a synchronization primitive
and must not be copied after first use.

Our GLS implementation assumes that every [DynVar] is declared statically in
module root. If a [DynVar] is deallocated by GC, every GLS which had that
variable's value will keep it until that entire GLS is cleaned up. However,
temporary dynamic variables are valid and safe when created and used on a
single goroutine which doesn't leak them to other goroutines, and does not
propagate its GLS to child goroutines.

See [GlsGo] for usage examples.
*/
type DynVar[A any] struct {
	lock sync.Mutex
	has  atomic.Bool
	def  func() A
	val  A
}

/*
If the value of this dynamic variable is set on the current goroutine,
returns that value. Otherwise:

  - If the default function was provided via [NewDynVar], returns the result of
    calling that function; this result is also stored and reused in later calls
    to [DynVar.Get] on goroutines where the value is not set.
  - Otherwise: returns the zero value.
*/
func (self *DynVar[A]) Get() A {
	gid := getGid()
	val, ok := self.got(gid)
	if ok {
		return val
	}
	return self.getDef(gid)
}

/*
Returns the current value of this dynamic variable, if set on the current
goroutine, or the zero value, and a boolean indicating if the value was set.
Does _not_ fall back on the default value, even if the default function was
provided in [NewDynVar].
*/
func (self *DynVar[A]) Got() (A, bool) { return self.got(getGid()) }

/*
Returns the current value of this dynamic variable, if set on the current
goroutine. Otherwise uses the given function to create a new value, and
sets it on the current goroutine. Subsequent calls to this method on the
same goroutine return the created value instead of calling the given function,
unless the value is unset at a later point via [DynVar.Clear] or [GlsClear].

When the given function is nil, this is equivalent to [DynVar.Get].

See [DynVar.Set] for notes on avoiding memory leakage, and follow its
recommendations for proper cleanup.
*/
func (self *DynVar[A]) GetOr(fun func() A) A {
	if fun == nil {
		return self.Get()
	}

	gid := getGid()
	gls := glss.getOrMake(gid)
	key := self.key()
	val, ok := gls[key].(A)
	if ok {
		return val
	}

	val = fun()
	gls[key] = val
	return val
}

/*
Sets the given value on the current goroutine. All subsequent calls to
[DynVar.Get] and [DynVar.Got] on the same goroutine will return this value,
unless the GLS is modified again.

Returns a [GlsVal] which snapshots the previous state of this variable on this
goroutine, and should be chained into deferred [GlsVal.Use] to restore it on
completion:

	defer someVar.Set(someVal).Use()

To avoid leaking memory, user code must ensure cleanup of GLS on goroutine
termination, by meeting at least one of the following conditions:

  - Option 0: chain [DynVar.Set] into deferred [GlsVal.Use].
  - Option 1: defer [GlsClear] in the top-level function of the current
    goroutine, before any other GLS modification.
  - Option 2: the current goroutine is spawned by [GlsGo] or [GlsGo1],
    which automatically use Option 1.

It's often useful to put variable modifications in small utility functions,
which are unable to defer cleanup. Deferring [GlsVal.Use] is also impossible
when using [DynVar.GetOr]. In such cases, user code must ensure either Option 1
or Option 2.
*/
func (self *DynVar[A]) Set(val A) GlsVal {
	return GlsVal{key: self.key(), val: val}.Use()
}

/*
Deletes the value of this dynamic variable from this goroutine's storage (GLS).
Has no effect on other goroutines.

Returns a [GlsVal] which snapshots the previous state of this variable on this
goroutine, and should be chained into deferred [GlsVal.Use] to restore the
entry on completion. See the comment on [DynVar.Set] for additional info.
*/
func (self *DynVar[A]) Clear() GlsVal {
	return GlsVal{key: self.key(), del: true}.Use()
}

/*
Returns a [GlsVal] whose [GlsVal.Use] sets this dynamic variable to the given
value on the goroutine where it's invoked. Eligible for any goroutine. Can be
passed as an override when calling [GlsGo], [GlsGo1], [GlsRun], [GlsRun1].
*/
func (self *DynVar[A]) With(val A) GlsVal {
	return GlsVal{key: self.key(), val: val}
}

/*
Returns a [GlsVal] whose [GlsVal.Use] clears this dynamic variable on the
goroutine where it's invoked. Eligible for any goroutine. Can be passed
as an override when calling [GlsGo], [GlsGo1], [GlsRun], [GlsRun1].
*/
func (self *DynVar[A]) WithClear() GlsVal {
	return GlsVal{key: self.key(), del: true}
}

/* Internal */

func (self *DynVar[A]) key() glsKey { return glsKey(unsafe.Pointer(self)) }

func (self *DynVar[A]) got(gid uint64) (A, bool) {
	src := glss.get(gid)[self.key()]
	val, ok := src.(A)
	return val, ok
}

func (self *DynVar[A]) getDef(gid uint64) (_ A) {
	if self.has.Load() {
		return self.val
	}
	return self.initDef(gid)
}

func (self *DynVar[A]) initDef(gid uint64) A {
	defer Lock(&self.lock).Unlock()

	if self.has.Load() {
		return self.val
	}

	fun := self.def

	/**
	Because the default function can only be set once, via the constructor,
	we can henceforth always route into the fast path.
	*/
	if fun == nil {
		self.has.Store(true)
		return self.val
	}

	val := fun()
	self.val = val
	self.def = nil // Nullify only on success.
	self.has.Store(true)
	return val
}
