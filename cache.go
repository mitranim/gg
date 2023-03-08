package gg

import (
	r "reflect"
	"sync"
	"time"
)

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

// Type-inferring shortcut for creating a `Cache` of the given type.
func CacheOf[
	Key comparable,
	Val any,
	Ptr Initer1Ptr[Val, Key],
]() *Cache[Key, Val, Ptr] {
	return new(Cache[Key, Val, Ptr])
}

type Cache[
	Key comparable,
	Val any,
	Ptr Initer1Ptr[Val, Key],
] struct {
	Lock sync.RWMutex
	Map  map[Key]Ptr
}

func (self *Cache[Key, Val, Ptr]) Get(key Key) Val { return *self.Ptr(key) }

func (self *Cache[Key, Val, Ptr]) Ptr(key Key) Ptr {
	ptr := self.get(key)
	if ptr != nil {
		return ptr
	}

	defer Lock(&self.Lock).Unlock()

	ptr = self.Map[key]
	if ptr != nil {
		return ptr
	}

	ptr = new(Val)
	ptr.Init(key)
	MapInit(&self.Map)[key] = ptr
	return ptr
}

func (self *Cache[Key, _, Ptr]) get(key Key) Ptr {
	defer Lock(self.Lock.RLocker()).Unlock()
	return self.Map[key]
}

func (self *Cache[Key, _, _]) Del(key Key) {
	defer Lock(&self.Lock).Unlock()
	delete(self.Map, key)
}

// Type-inferring shortcut for creating a `TypeCache` of the given type.
func TypeCacheOf[Val any, Ptr Initer1Ptr[Val, r.Type]]() *TypeCache[Val, Ptr] {
	return new(TypeCache[Val, Ptr])
}

type TypeCache[Val any, Ptr Initer1Ptr[Val, r.Type]] struct {
	Lock sync.RWMutex
	Map  map[r.Type]Ptr
}

func (self *TypeCache[Val, Ptr]) Get(key r.Type) Val { return *self.Ptr(key) }

func (self *TypeCache[Val, Ptr]) Ptr(key r.Type) Ptr {
	ptr := self.get(key)
	if ptr != nil {
		return ptr
	}

	defer Lock(&self.Lock).Unlock()

	ptr = self.Map[key]
	if ptr != nil {
		return ptr
	}

	ptr = new(Val)
	ptr.Init(key)

	if self.Map == nil {
		self.Map = map[r.Type]Ptr{}
	}
	self.Map[key] = ptr

	return ptr
}

func (self *TypeCache[Val, Ptr]) get(key r.Type) Ptr {
	defer Lock(self.Lock.RLocker()).Unlock()
	return self.Map[key]
}

/*
Tool for deduplicating and caching expensive work.
All methods are safe for concurrent use.
*/
type Mem[A any] struct {
	sync.RWMutex
	Timed[A]
}

// Clears the inner value and timestamp.
func (self *Mem[A]) Clear() {
	defer Lock(self).Unlock()
	Clear(&self.Val)
	Clear(&self.Inst)
}

// Returns the inner `Timed`.
func (self *Mem[A]) GetTimed() Timed[A] {
	defer Lock(self.RLocker()).Unlock()
	return self.Timed
}

// Returns the inner value, if any.
func (self *Mem[A]) Get() A { return self.GetTimed().Val }

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
	val := self.GetTimed()
	if !val.IsExpired(life) {
		return val.Val
	}

	defer Lock(self).Unlock()

	if fun != nil && self.Timed.IsExpired(life) {
		self.Timed.Set(fun())
	}
	return self.Timed.Val
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
