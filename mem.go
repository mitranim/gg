package gg

import (
	"sync"
	"time"
)

// Should be embedded in other types. See `Mem.DedupFrom` for an example.
type MemSecond[A any] struct {
	DurSecond
	Mem[A]
}

// Should be embedded in other types. See `Mem.DedupFrom` for an example.
type MemMinute[A any] struct {
	DurMinute
	Mem[A]
}

// Should be embedded in other types. See `Mem.DedupFrom` for an example.
type MemHour[A any] struct {
	DurHour
	Mem[A]
}

/*
Tool for deduplicating and caching expensive work. All methods are safe for
concurrent use.
*/
type Mem[A any] struct {
	timed Timed[A]
	lock  sync.RWMutex
}

// Implement `Getter`. Returns the inner value, if any.
func (self *Mem[A]) Get() A { return self.Timed().Val }

/*
Returns `Timed` that contains the current value and the timestamp at which the
value was last set.
*/
func (self *Mem[A]) Timed() Timed[A] {
	defer Lock(self.lock.RLocker()).Unlock()
	return self.timed
}

// Clears the inner value and timestamp.
func (self *Mem[A]) Clear() {
	defer Lock(&self.lock).Unlock()
	self.timed.Clear()
}

/*
Shortcut for types such as `MemHour` that embed `Mem` and implement `Dur`.
Usage:

	type MemExample struct { MemHour[string] }

	func (self *MemExample) Get() string {return self.DedupFrom(self)}
	func (*MemExample) Make() string {return `some_value`}
*/
func (self *Mem[A]) DedupFrom(val DurMaker[A]) A {
	return self.Dedup(val.Duration(), val.Make)
}

/*
Either reuses the existing value, or calls the given function to regenerate it.
The given duration is the allotted lifetime of the previous value, if any.

In addition to reusing previous values, this method deduplicates concurrent
work. When called concurrently by multiple goroutines, only 1 goroutine
performs work, while the others simply wait for it.

Usage:

	type MemExample struct { Mem[string] }

	func (self *MemExample) Get() string {
		return self.Dedup(time.Hour, func() string {return `some_value`})
	}
*/
func (self *Mem[A]) Dedup(life time.Duration, fun func() A) A {
	val := self.Timed()
	if !val.IsExpired(life) {
		return val.Val
	}

	defer Lock(&self.lock).Unlock()

	if fun != nil && self.timed.IsExpired(life) {
		self.timed.Set(fun())
	}
	return self.timed.Val
}

/*
Implement `json.Marshaler` by proxying to `Timed.MarshalJSON` on the inner
instance of `.Timed`. Like other methods, this is safe for concurrent use.
*/
func (self *Mem[_]) MarshalJSON() ([]byte, error) {
	defer Lock(&self.lock).Unlock()
	return self.timed.MarshalJSON()
}

/*
Implement `json.Unmarshaler` by proxying to `Timed.UnmarshalJSON` on the inner
instance of `.Timed`. Like other methods, this is safe for concurrent use.
*/
func (self *Mem[_]) UnmarshalJSON(src []byte) error {
	defer Lock(&self.lock).Unlock()
	return self.timed.UnmarshalJSON(src)
}

/*
Implements `Dur` by returning `time.Second`. This type is zero-sized, and can be
embedded in other types, like a mixin, at no additional cost.
*/
type DurSecond struct{}

// Implement `Dur` by returning `time.Second`.
func (DurSecond) Duration() time.Duration { return time.Second }

/*
Implements `Dur` by returning `time.Minute`. This type is zero-sized, and can be
embedded in other types, like a mixin, at no additional cost.
*/
type DurMinute struct{}

// Implement `Dur` by returning `time.Minute`.
func (DurMinute) Duration() time.Duration { return time.Minute }

/*
Implements `Dur` by returning `time.Hour`. This type is zero-sized, and can be
embedded in other types, like a mixin, at no additional cost.
*/
type DurHour struct{}

// Implement `Dur` by returning `time.Hour`.
func (DurHour) Duration() time.Duration { return time.Hour }
