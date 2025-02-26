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
func (self *LazyIniter[Val, _]) Get() Val {
	out, ok := self.got()
	if ok {
		return out
	}
	return self.get()
}

/*
Clears the underlying value. After this call, the next call to `LazyIniter.Get`
or `LazyIniter.Ptr` will reinitialize by invoking the `.Init` method of the
underlying value.
*/
func (self *LazyIniter[_, _]) Clear() {
	defer Lock(&self.lock).Unlock()
	self.val.Clear()
}

/*
Resets the underlying value to the given input. After this call, the underlying
value is considered to be initialized. Further calls to `LazyIniter.Get` or
`LazyIniter.Ptr` will NOT reinitialize until `.Clear` is called.
*/
func (self *LazyIniter[Val, _]) Reset(src Val) {
	defer Lock(&self.lock).Unlock()
	self.val.Set(src)
}

func (self *LazyIniter[Val, _]) got() (_ Val, _ bool) {
	defer Lock(self.lock.RLocker()).Unlock()
	if self.val.IsNotNull() {
		return self.val.Val, true
	}
	return
}

func (self *LazyIniter[Val, Ptr]) get() Val {
	defer Lock(&self.lock).Unlock()

	if self.val.IsNull() {
		Ptr(&self.val.Val).Init()
		self.val.Ok = true
	}

	return self.val.Val
}
