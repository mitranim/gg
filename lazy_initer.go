package gg

import "sync"

/*
Shortcut for type inference. The following is equivalent:

	NewLazyIniter[Val]()
	new(LazyIniter[Val, Ptr])
*/
func NewLazyIniter[Val any, Ptr IniterPtr[Val]]() *LazyIniter[Val, Ptr] {
	return new(LazyIniter[Val, Ptr])
}

/*
Encapsulates a lazily-initializable value. The first call to `.Get` or `.Ptr`
initializes the underlying value by calling its `.Init` method. Initialization
is performed exactly once. Access is synchronized. All methods are
concurrency-safe. Designed to be embeddable. A zero value is ready to use. When
using this as a struct field, you don't need to explicitly initialize the
field. Contains a mutex and must not be copied.
*/
type LazyIniter[Val any, Ptr IniterPtr[Val]] struct {
	val  Opt[Val]
	lock sync.RWMutex
}

// Returns the underlying value, lazily initializing it on the first call.
func (self *LazyIniter[Val, _]) Get() Val { return *self.Ptr() }

/*
Returns a pointer to the underlying value, lazily initializing it on the first
call.
*/
func (self *LazyIniter[_, Ptr]) Ptr() Ptr {
	if self.inited() {
		return self.ptr()
	}
	return self.init()
}

func (self *LazyIniter[_, Ptr]) ptr() Ptr { return &self.val.Val }

func (self *LazyIniter[_, _]) inited() bool {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.val.IsNotNull()
}

func (self *LazyIniter[_, Ptr]) init() Ptr {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.val.IsNotNull() {
		return self.ptr()
	}

	Ptr(&self.val.Val).Init()
	self.val.Ok = true
	return self.ptr()
}
