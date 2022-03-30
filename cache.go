package gg

import (
	r "reflect"
	"sync"
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

func (self *Cache[Key, Val, Ptr]) Get(key Key) Val { return *self.GetPtr(key) }

func (self *Cache[Key, Val, Ptr]) GetPtr(key Key) Ptr {
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

func (self *TypeCache[Val, Ptr]) Get(key r.Type) Val { return *self.GetPtr(key) }

func (self *TypeCache[Val, Ptr]) GetPtr(key r.Type) Ptr {
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
