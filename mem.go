package gg

import (
	"sync"
	"time"
)

/*
Tool for deduplicating and caching expensive work. All methods are safe for
concurrent use. The first type parameter is used to determine expiration
duration, and should be a zero-sized stateless type, such as `DurSecond`,
`DurMinute`, `DurHour`, and `DurForever` provided by this package. The given
type `Tar` must implement `Initer` on its pointer type: `(*Tar).Init`. The init
method is used to populate data whenever it's missing or expired. See methods
`Mem.Get` and `Mem.Peek`. A zero value of `Mem` is ready for use. Contains a
synchronization primitive and must not be copied.

Usage example:

	type Model struct { Id uint64 }

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
	val  Tar
	ok   bool
	inst time.Time
	lock sync.RWMutex
}

/*
Returns the inner value after ensuring it's initialized and not expired. If the
data is missing or expired, it's initialized by calling `(*Tar).Init`. Otherwise
the data is returned as-is.

This method avoids redundant concurrent work. When called concurrently by
multiple goroutines, only 1 goroutine performs work, while the others simply
wait for it.

Method `(*Tar).Init` is always called on a new pointer to a zero value, for
multiple reasons. If `(*Tar).Init` appends data to the target instead of
replacing it, this avoids accumulating redundant data and leaking memory.
Additionally, this avoids accidental concurrent modification between `Mem` and
its callers that could lead to observing an inconsistent state of the data.

Expiration is determined by consulting the `Durationer` type provided to `Mem`
as its first type parameter, calling `.Duration` on a zero value. As a special
case, 0 duration is considered indefinite, making the `Mem` never expire, and
thus functionally equivalent to `LazyIniter`. Negative durations cause the `Mem`
to expire immediately, making it pointless.

Compare `Mem.Peek` which does not perform initialization.
*/
func (self *Mem[Dur, Tar, Ptr]) Get() Tar {
	defer Lock(&self.lock).Unlock()

	val, ok, inst := self.val, self.ok, self.inst
	if ok && !isExpired(inst, Zero[Dur]().Duration()) {
		return val
	}

	var tar Tar
	Ptr(&tar).Init()
	self.set(tar)
	return tar
}

/*
Similar to `Mem.Get` but returns the inner value as-is, without checking
expiration. If the value was never initialized, it's zero.
*/
func (self *Mem[_, Tar, _]) Peek() Tar {
	defer Lock(self.lock.RLocker()).Unlock()
	return self.val
}

// Clears the inner value and timestamp.
func (self *Mem[_, Tar, _]) Clear() {
	defer Lock(&self.lock).Unlock()
	self.clear()
}

/*
Implement `json.Marshaler`. If the value is not initialized, returns a
representation of JSON null. Otherwise uses `json.Marshal` to encode the
underlying value, even if expired. Like other methods, this is safe for
concurrent use.
*/
func (self *Mem[_, _, _]) MarshalJSON() ([]byte, error) {
	defer Lock(self.lock.RLocker()).Unlock()
	if !self.ok {
		return ToBytes(`null`), nil
	}
	return JsonBytesCatch(self.val)
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON null,
clears the inner value and the timestamp. Otherwise uses `json.Unmarshal` to
decode into the inner value, setting the current timestamp on success. Like
`Mem.Get`, this uses an intermediary zero value, avoiding corruption of the
existing inner value in cases of partially failed decoding. Like other methods,
this is safe for concurrent use.
*/
func (self *Mem[_, Tar, _]) UnmarshalJSON(src []byte) error {
	defer Lock(&self.lock).Unlock()

	if IsJsonEmpty(src) {
		self.clear()
		return nil
	}

	var tar Tar
	err := JsonDecodeCatch(src, &self.val)
	if err != nil {
		return err
	}

	self.set(tar)
	return nil
}

func (self *Mem[_, Tar, _]) clear() {
	var val Tar
	self.val = val
	self.ok = false
	self.inst = time.Time{}
}

func (self *Mem[_, Tar, _]) set(val Tar) {
	self.val = val
	self.ok = true
	self.inst = time.Now()
}

func isExpired(inst time.Time, dur time.Duration) bool {
	return dur != 0 && inst.Add(dur).Before(time.Now())
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

/*
Implements `Durationer` by returning 0, which is understood by `Mem` as
indefinite, making it never expire. This type is zero-sized, and can be
embedded in other types, like a mixin, at no additional cost.
*/
type DurForever struct{}

// Implement `Durationer` by returning 0.
func (DurForever) Duration() time.Duration { return 0 }
