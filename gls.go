package gg

/*
This file implements various tools for goroutine-local storage (GLS)
accessed via dynamic variables ([DynVar]).
*/

/*
Wraps the keyword `go` and ensures that the child goroutine inherits the GLS
(goroutine-local storage) of the parent goroutine, and that the GLS is cleaned
up when done, avoiding a memory leak. Accepts any amount of GLS overrides,
which take effect on the child goroutine. GLS is accessed via [DynVar].

The GLS is snapshotted on the parent goroutine before starting the child.
The child receives a copy of the parent's GLS. Modification of the parent's
storage, or termination of the parent, has no effect on the child.

Example:

	// Define this statically in module root.
	var CTX = gg.NewDynVar(context.Background)

	// Any code can access the value.
	func subFunc() { _ = CTX.Get() }

	// In some goroutine, set the value locally.
	// Always defer cleanup to avoid leaking memory.
	func someFunc() {
		defer CTX.Set(someCtx).Use()

		// Spawn sub-goroutines inheriting the GLS:
		gg.GlsGo(subFunc)

		// Spawn with overrides:
		gg.GlsGo(subFunc, CTX.With(someOtherCtx))
	}
*/
func GlsGo(run func(), overrides ...GlsVal) {
	if run == nil {
		return
	}
	go withGls(glsCopy(overrides...), run)
}

/*
Variant of [GlsGo] which passes one given argument to the function it runs.
Sometimes useful for avoiding an anonymous closure when wrapping the given
function in another. Example:

	func Go(run func()) { gg.GlsGo1(runWithRecovery, run) }

	func runWithRecovery(run func()) {
		defer recoverAndLogError()
		run()
	}
*/
func GlsGo1[A any](run func(A), val A, overrides ...GlsVal) {
	if run == nil {
		return
	}
	go withGls1(glsCopy(overrides...), run, val)
}

/*
Like [GlsGo] but runs the given function immediately on the current goroutine.
The provided GLS entry overrides are in effect for the duration of the call;
previous values are restored before returning. Non-overridden GLS entries
behave normally.
*/
func GlsRun(run func(), overrides ...GlsVal) {
	if run == nil {
		return
	}
	if len(overrides) > 0 {
		defer glsValsUse(glsValsSwap(overrides))
	}
	run()
}

// Like [GlsRun], but takes and passes an additional argument like [GlsGo1].
func GlsRun1[A any](run func(A), val A, overrides ...GlsVal) {
	if run == nil {
		return
	}
	if len(overrides) > 0 {
		defer glsValsUse(glsValsSwap(overrides))
	}
	run(val)
}

/*
Replaces the current goroutine's GLS (goroutine-local storage), if any, with
the given values. The values are usually created via [GlsSnap] on a parent
goroutine, and installed on a child goroutine by calling this function.

Manually propagating GLS via [GlsSnap] and [GlsSet] is necessary for goroutines
which are started by external code you don't control, but invoke your own
callback. Normally, you should start goroutines via [GlsGo] or [GlsGo1], which
automatically propagate GLS from parents to children, and ensure cleanup.

Invoking this with no arguments is equivalent to [GlsClear]. The latter is
provided because it can be passed to functions which only take `func()`.

Returns the previous GLS; if GLS was empty, returns a zero value.

Should be chained into deferred [Gls.Use] to restore the previous state later.
If the previous state was empty, [Gls.Use] will delete the GLS from the central
registry. Without this, the new GLS might never be deleted, and you will leak
memory. Canonical usage:

	vals := gg.GlsSnap()

	someGoroutineCreatingFunc(func() {
		defer gg.GlsSet(vals...).Use()
	})

At the beginning of a goroutine, a valid alternative is to use a deferred
[GlsClear]:

	vals := gg.GlsSnap()

	someGoroutineCreatingFunc(func() {
		defer gg.GlsClear()
		gg.GlsSet(vals...)
	})
*/
func GlsSet(vals ...GlsVal) Gls {
	gid := getGid()
	prev := glss.get(gid)
	size := len(vals)

	if size <= 0 {
		glss.del(gid)
	} else {
		gls := make(gls, size)
		for _, val := range vals {
			if !val.del {
				gls[val.key] = val.val
			}
		}
		glss.set(gid, gls)
	}
	return Gls{val: prev}
}

/*
Creates a snapshot of the current goroutine's GLS (goroutine-local storage)
suitable for installation on other goroutines via [GlsSet]. The snapshot is
considered immutable, is not tied to any particular goroutine, and can be used
on any amount of goroutines.

Accepts any amount of overrides for specific dynamic variables.
Overrides are created with [DynVar.With] and [DynVar.WithClear].

Unnecessary when using [GlsGo] and [GlsGo1], which automatically
copy GLS from parent to child goroutines.

The order of entries is undefined and may vary between calls.
*/
func GlsSnap(overrides ...GlsVal) []GlsVal {
	gls := glss.get(getGid())
	size := len(overrides)
	keys := make(Set[glsKey], size)
	out := make([]GlsVal, 0, size+len(gls))

	for _, val := range overrides {
		out = append(out, val)
		keys.Add(val.key)
	}
	for key, val := range gls {
		if !keys.Has(key) {
			out = append(out, GlsVal{key: key, val: val})
		}
	}
	return out
}

/*
Clears the current goroutine's GLS (the value of every [DynVar]). Has no effect
on the GLS of other goroutines. See [Gls.Use] for use cases and an example.
Unnecessary when using [GlsGo] or [GlsGo1].

When running an HTTP server, make sure to include this at the start of your
request handler:

	func handleRequest(http.ResponseWriter, *http.Request) {
		defer gg.ClsClear()
		// ... respond
	}

More generally, make sure to defer this at the start of every function which may
be repeatedly invoked on a new goroutine, such as HTTP handlers, TCP handlers,
websocket handlers, and so on. Otherwise, your app may gradually leak memory by
accumulating leftover GLS of goroutines which are no longer running.

In testing, each test, sub-test, and benchmark runs on a new goroutine.
GLS cleanup in testing is optional, but can be done like this:

	t.Cleanup(gg.GlsClear)
	b.Cleanup(gg.GlsClear)

This function is equivalent to calling [GlsSet] with no arguments. The latter
is more flexible. The former is provided because it can be passed to functions
which only take `func()`, like in the example above.
*/
func GlsClear() { glss.del(getGid()) }

/*
An opaque structure representing GLS (goroutine-local storage) of a single
goroutine. Returned by [GlsSet]. May contain mutable state. Permanently moving
this to another goroutine is valid, but sharing this between goroutines is
invalid and may lead to silent logical errors and occasional panics.

GLS is accessed via dynamic variables: [DynVar].
*/
type Gls struct{ val gls }

/*
Replaces the current goroutine's GLS with this snapshot. If the snapshot is
empty, deletes the GLS from the central registry. Returns the previous state of
the current goroutine's GLS, which is no longer stored in the central registry
and can be used to restore that state. See [GlsSet] for usage examples.
*/
func (self Gls) Use() Gls {
	gid := getGid()
	prev := glss.get(gid)
	gls := self.val
	if len(gls) > 0 {
		glss.set(gid, gls)
	} else {
		glss.del(gid)
	}
	return Gls{val: prev}
}

/*
An opaque structure representing a single entry in goroutine-local storage,
for a specific [DynVar]. May represent either existence of a specific value,
or non-existence of the variable's value.

For any [GlsVal] created with [DynVar.Set] or [DynVar.Clear], its [GlsVal.Use]
should be immediately deferred to restore the previous state when done.

A [GlsVal] created via [DynVar.With] or [DynVar.WithClear] can be passed as an
override when calling [GlsGo], [GlsGo1], [GlsRun], [GlsRun1].

A [GlsVal] is considered immutable, is not tied to any particular goroutine,
and can be used on any amount of goroutines.

A [GlsVal] is tied to its [DynVar] and its lifetime. A [GlsVal] not created
by a [DynVar] is invalid. A [GlsVal] whose [DynVar] is no longer reachable
is invalid. Prefer to declare dynamic variables statically in module root.
*/
type GlsVal struct {
	key glsKey
	val any
	del bool
}

/*
Applies this entry to the current GLS, returning a snapshot of the previous
state of that entry in the current GLS, which can be chained into deferred
[GlsVal.Use] to restore it later. See [DynVar.Set].
*/
func (self GlsVal) Use() GlsVal {
	// SYNC[gls_val_use].

	gid := getGid()
	key := self.key

	if !self.del {
		gls := glss.getOrMake(gid)
		prev, had := gls[key]
		gls[key] = self.val
		return GlsVal{key: key, val: prev, del: !had}
	}

	gls := glss.get(gid)
	prev, had := gls[key]
	delete(gls, key)
	if len(gls) <= 0 {
		glss.del(gid)
	}
	return GlsVal{key: key, val: prev, del: !had}
}
