package gg

var glss glss_t

/*
Short for "goroutine-local storage storage". Keys are goroutine IDs.

See `BenchmarkDynVar_with_minor_concurrency` for some perf notes.
*/
type glss_t struct{ val SyncMap[uint64, gls] }

func (self *glss_t) get(gid uint64) gls {
	out, _ := self.val.Load(gid)
	return out
}

/*
The apparent race condition between load / check / make / store is benign,
because each goroutine creates and accesses only its own GLS.
*/
func (self *glss_t) getOrMake(gid uint64) gls {
	out := self.get(gid)
	if out != nil {
		return out
	}

	out = make(gls, 1)
	self.val.Store(gid, out)
	return out
}

func (self *glss_t) set(gid uint64, val gls) { self.val.Store(gid, val) }

func (self *glss_t) del(gid uint64) { self.val.Delete(gid) }

/*
Short for "goroutine-local storage".

We create, read, and mutate each `gls` only on its own goroutine.

At the time of writing, our GLS API makes it possible for user code to create
more than one `gls` for the same goroutine, for example by repeatedly calling
[GlsSet] with different sets of values. As a result, most code should access
`gls` from `glss` without caching it; [Gls] is the only exception.
*/
type gls = map[glsKey]any

// GLS keys are [DynVar] pointers.
type glsKey uintptr

/*
We assume that each dyn var is declared statically and never freed. If dynamic
allocation and deallocation of dyn vars was a real use case, we could track
liveness via [runtime.AddCleanup]. Something like:

	type glsKey struct {
		ptr  uintptr
		live atomic.Bool
	}

	func (self *DynVar[_]) getKey() *glsKey {
		// The following must only be done once.
		// Synchronization is elided for example's sake.

		key := &glsKey{ptr: uintptr(unsafe.Pointer(self))}
		runtime.AddCleanup(self, glsOnVarCleanup, key)
		return key
	}

	func glsOnVarCleanup(key *glsKey) {
		key.live.Store(false)

		// And here we would have to delete entries keyed by `key.ptr`
		// from every GLS in `glss`.
	}

When using a [GlsVal], we would check liveness and skip / ignore "dead" vals,
which is easy enough. But deleting the entries from every `gls` would require
additional synchronization at the level of each `gls`, since the finalizer runs
in its own goroutine, and slow everything down for the sake of a use case which
is probably not even real.
*/

/*
Like [GlsSnap] but creates a new map directly. Slightly more efficient than
combining [GlsSnap] with [GlsSet]. Private because the difference is not large
enough to justify bloating the library's interface, and because this is easy to
misuse by passing the resulting [Gls] to multiple child goroutines.
*/
func glsCopy(overrides ...GlsVal) Gls {
	src := glss.get(getGid())
	out := make(gls, len(src)+len(overrides))

	for key, val := range src {
		out[key] = val
	}
	for _, val := range overrides {
		glsValUse(out, val)
	}
	return Gls{val: out}
}

// Must be passed into deferred `glsValsUse`.
func glsValsSwap(vals []GlsVal) (_ uint64, _ gls, _ []GlsVal) {
	// SYNC[gls_vals_swap_nop].
	if len(vals) <= 0 {
		return
	}

	gid := getGid()
	gls := glss.getOrMake(gid)
	out := make([]GlsVal, 0, len(vals))

	for _, val := range vals {
		key := val.key
		prev, had := gls[key]
		out = append(out, GlsVal{key: key, val: prev, del: !had})
		glsValUse(gls, val)
	}

	return gid, gls, out
}

func glsValsUse(gid uint64, gls gls, vals []GlsVal) {
	// SYNC[gls_vals_swap_nop].
	if len(vals) <= 0 {
		return
	}
	for _, val := range vals {
		glsValUse(gls, val)
	}
	if len(gls) <= 0 {
		glss.del(gid)
	}
}

/*
Caution: in most cases, the caller must also check the remaining GLS size
and delete an empty GLS from the GLSS.

SYNC[gls_val_use].
*/
func glsValUse(gls gls, val GlsVal) {
	if val.del {
		delete(gls, val.key)
	} else {
		gls[val.key] = val.val
	}
}

func withGls(gls Gls, run func()) {
	gls.Use()
	defer GlsClear()
	run()
}

func withGls1[A any](gls Gls, run func(A), val A) {
	gls.Use()
	defer GlsClear()
	run(val)
}
