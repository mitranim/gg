package gg

import (
	"sync"
	"time"
)

/*
Tool for deduplicating and caching expensive work. All methods are safe for
concurrent use. The first type parameter is used to determine expiration
duration, and should typically be a zero-sized stateless type, such as
`DurSecond`, `DurMinute`, `DurHour` provided by this package. The given type
`Tar` must implement `Initer` on its pointer type: `(*Tar).Init`. The init
method is used to populate data whenever it's missing or expired. See methods
`Mem.Get` and `Mem.Peek`. A zero value of `Mem` is ready for use.
Contains a synchronization primitive and must not be copied.

Usage example:

	type Model struct { Id string }

	type DatModels []Model

	// Pretend that this is expensive.
	func (self *DatModels) Init() { *self = DatModels{{10}, {20}, {30}} }

	type Dat struct {
		Models gg.Mem[gg.DurHour, DatModels, *DatModels]
		// ... other caches
	}

	var dat Dat

	func init() { fmt.Println(dat.Models.Get()) }
*/
type Mem[Dur Durationer, Tar any, Ptr IniterPtr[Tar]] struct {
	timed Timed[Tar]
	lock  sync.RWMutex
}

/*
Returns the inner value after ensuring it's initialized and not expired. If the
data is missing or expired, it's initialized by calling `(*Tar).Init`. Otherwise
the data is returned as-is.

This method avoids redundant concurrent work. When called concurrently by
multiple goroutines, only 1 goroutine performs work, while the others simply
wait for it.

Expiration is determined by consulting the `Dur` type provided to `Mem` as its
first type parameter, calling `Dur.Duration` on a zero value.

Method `(*Tar).Init` is always called on a new pointer to a zero value, for
multiple reasons. If `(*Tar).Init` appends data to the target instead of
replacing it, this avoids accumulating redundant data and leaking memory.
Additionally, this avoids the possibility of concurrent modification
(between `Mem` and its callers) that could lead to observing an inconsistent
state of the data.

Compare `Mem.Peek` which does not perform initialization.
*/
func (self *Mem[_, Tar, _]) Get() Tar { return self.Timed().Get() }

/*
Same as `Mem.Get` but returns `Timed` which contains both the inner value and
the timestamp at which it was generated or set.
*/
func (self *Mem[_, Tar, Ptr]) Timed() Timed[Tar] {
	dur := self.Duration()
	val := self.PeekTimed()
	if !val.IsExpired(dur) {
		return val
	}

	defer Lock(&self.lock).Unlock()
	timed := &self.timed

	if timed.IsExpired(dur) {
		// See comment on `Mem.Get` why we avoid calling this on `&self.timed.Val`.
		var tar Tar
		Ptr(&tar).Init()
		timed.Set(tar)
	}
	return *timed
}

/*
Similar to `Mem.Get` but returns the inner value as-is, without performing
initialization.
*/
func (self *Mem[_, Tar, _]) Peek() Tar { return self.PeekTimed().Get() }

/*
Similar to `Mem.Timed` but returns inner `Timed` as-is, without performing
initialization. The result contains the current state of the data, and the
timestamp at which the value was last set. If the data has not been initialized
or set, the timestamp is zero.
*/
func (self *Mem[_, Tar, _]) PeekTimed() Timed[Tar] {
	defer Lock(self.lock.RLocker()).Unlock()
	return self.timed
}

// Clears the inner value and timestamp.
func (self *Mem[_, _, _]) Clear() {
	defer Lock(&self.lock).Unlock()
	self.timed.Clear()
}

/*
Implement `Durationer`. Shortcut for `Zero[Dur]().Duration()` using the first
type parameter provided to this type. For internal use.
*/
func (*Mem[Dur, _, _]) Duration() time.Duration { return Zero[Dur]().Duration() }

/*
Implement `json.Marshaler` by proxying to `Timed.MarshalJSON` on the inner
instance of `.Timed`. Like other methods, this is safe for concurrent use.
*/
func (self *Mem[_, _, _]) MarshalJSON() ([]byte, error) {
	defer Lock(&self.lock).Unlock()
	return self.timed.MarshalJSON()
}

/*
Implement `json.Unmarshaler` by proxying to `Timed.UnmarshalJSON` on the inner
instance of `.Timed`. Like other methods, this is safe for concurrent use.
*/
func (self *Mem[_, _, _]) UnmarshalJSON(src []byte) error {
	defer Lock(&self.lock).Unlock()
	return self.timed.UnmarshalJSON(src)
}

/*
Implements `Durationer` by returning `time.Second`. This type is zero-sized, and
can be embedded in other types, like a mixin, at no additional cost.
*/
type DurSecond struct{}

// Implement `Durationer` by returning `time.Second`.
func (DurSecond) Duration() time.Duration { return time.Second }

/*
Implements `Durationer` by returning `time.Minute`. This type is zero-sized, and
can be embedded in other types, like a mixin, at no additional cost.
*/
type DurMinute struct{}

// Implement `Durationer` by returning `time.Minute`.
func (DurMinute) Duration() time.Duration { return time.Minute }

/*
Implements `Durationer` by returning `time.Hour`. This type is zero-sized, and
can be embedded in other types, like a mixin, at no additional cost.
*/
type DurHour struct{}

// Implement `Durationer` by returning `time.Hour`.
func (DurHour) Duration() time.Duration { return time.Hour }
