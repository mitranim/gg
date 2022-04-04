package gg

import (
	r "reflect"
	"sync"
	"time"
)

/*
Takes a function that ought to be called no more than once. Returns a function
that caches and reuses the result of the original function. Uses `sync.Once`
internally.
*/
func Lazy[A any](fun func() A) func() A { return (&lazy[A]{fun: fun}).do }

/*
Variant of `Lazy` that takes an additional argument and passes it to the given
function when it's executed, which happens no more than once.
*/
func Lazy1[A, B any](fun func(B) A, val B) func() A {
	return Lazy(func() A { return fun(val) })
}

/*
Variant of `Lazy` that takes additional arguments and passes them to the given
function when it's executed, which happens no more than once.
*/
func Lazy2[A, B, C any](fun func(B, C) A, val0 B, val1 C) func() A {
	return Lazy(func() A { return fun(val0, val1) })
}

/*
Variant of `Lazy` that takes additional arguments and passes them to the given
function when it's executed, which happens no more than once.
*/
func Lazy3[A, B, C, D any](fun func(B, C, D) A, val0 B, val1 C, val2 D) func() A {
	return Lazy(func() A { return fun(val0, val1, val2) })
}

type lazy[A any] struct {
	ref sync.Once
	val A
	fun func() A
}

func (self *lazy[A]) do() A {
	self.ref.Do(self.init)
	return self.val
}

func (self *lazy[A]) init() {
	fun := self.fun
	if fun != nil {
		self.fun = nil
		self.val = fun()
	}
}

func CacheOf[
	Key comparable,
	Val any,
	Ptr Initer1[Val, Key],
]() *Cache[Key, Val, Ptr] {
	return &Cache[Key, Val, Ptr]{Map: map[Key]Ptr{}}
}

type Cache[
	Key comparable,
	Val any,
	Ptr Initer1[Val, Key],
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

	defer Locked(&self.Lock).Unlock()

	ptr = self.Map[key]
	if ptr != nil {
		return ptr
	}

	ptr = new(Val)
	ptr.Init(key)
	self.Map[key] = ptr
	return ptr
}

func (self *Cache[Key, _, Ptr]) get(key Key) Ptr {
	defer Locked(self.Lock.RLocker()).Unlock()
	return self.Map[key]
}

func (self *Cache[Key, _, _]) Del(key Key) {
	defer Locked(&self.Lock).Unlock()
	delete(self.Map, key)
}

func TypeCacheOf[Val any, Ptr Initer1[Val, r.Type]]() *TypeCache[Val, Ptr] {
	return &TypeCache[Val, Ptr]{Map: map[r.Type]Ptr{}}
}

type TypeCache[Val any, Ptr Initer1[Val, r.Type]] struct {
	Lock sync.RWMutex
	Map  map[r.Type]Ptr
}

func (self *TypeCache[Val, Ptr]) Get(key r.Type) Val { return *self.Ptr(key) }

func (self *TypeCache[Val, Ptr]) Ptr(key r.Type) Ptr {
	ptr := self.get(key)
	if ptr != nil {
		return ptr
	}

	defer Locked(&self.Lock).Unlock()

	ptr = self.Map[key]
	if ptr != nil {
		return ptr
	}

	ptr = new(Val)
	ptr.Init(key)
	self.Map[key] = ptr
	return ptr
}

func (self *TypeCache[Val, Ptr]) get(key r.Type) Ptr {
	defer Locked(self.Lock.RLocker()).Unlock()
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
	defer Locked(self).Unlock()
	Clear(&self.Val)
	Clear(&self.Inst)
}

// Returns the inner `Timed`.
func (self *Mem[A]) GetTimed() Timed[A] {
	defer Locked(self.RLocker()).Unlock()
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

	defer Locked(self).Unlock()

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
