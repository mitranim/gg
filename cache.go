package gg

import (
	r "reflect"
	"sync"
)

// Type-inferring shortcut for creating a `Cache` of the given type.
func CacheOf[
	Key comparable,
	Val any,
	Ptr Initer1Ptr[Val, Key],
]() *Cache[Key, Val, Ptr] {
	return new(Cache[Key, Val, Ptr])
}

// Concurrency-safe cache. See the method reference.
type Cache[
	Key comparable,
	Val any,
	Ptr Initer1Ptr[Val, Key],
] struct {
	Map  map[Key]Ptr
	Lock sync.RWMutex
}

/*
Shortcut for using `.Ptr` and dereferencing the result. May be invalid if the
resulting value is non-copyable, for example when it contains a mutex.
*/
func (self *Cache[Key, Val, Ptr]) Get(key Key) Val { return *self.Ptr(key) }

/*
Returns the cached value for the given key. If the value did not previously
exist, idempotently initializes it by calling `.Init` (by pointer) and caches
the result. For any given key, the value is initialized exactly once, even if
multiple goroutines are trying to access it simultaneously.
*/
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

// Deletes the value for the given key.
func (self *Cache[Key, _, _]) Del(key Key) {
	defer Lock(&self.Lock).Unlock()
	delete(self.Map, key)
}

// Type-inferring shortcut for creating a `TypeCache` of the given type.
func TypeCacheOf[Val any, Ptr Initer1Ptr[Val, r.Type]]() *TypeCache[Val, Ptr] {
	return new(TypeCache[Val, Ptr])
}

/*
Tool for storing information derived from `reflect.Type` that can be generated
once and then cached. Used internally. All methods are concurrency-safe.
*/
type TypeCache[Val any, Ptr Initer1Ptr[Val, r.Type]] struct {
	Map  map[r.Type]Ptr
	Lock sync.RWMutex
}

/*
Shortcut for using `.Ptr` and dereferencing the result. May be invalid if the
resulting value is non-copyable, for example when it contains a mutex.
*/
func (self *TypeCache[Val, Ptr]) Get(key r.Type) Val { return *self.Ptr(key) }

/*
Returns the cached value for the given key. If the value did not previously
exist, idempotently initializes it by calling `.Init` (by pointer) and caches
the result. For any given key, the value is initialized exactly once, even if
multiple goroutines are trying to access it simultaneously.
*/
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
