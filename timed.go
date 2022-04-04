package gg

import "time"

/*
Shortcut for creating a `Timed` with the given value, using the current
timestamp.
*/
func TimedVal[A any](val A) (out Timed[A]) {
	out.Set(val)
	return
}

/*
Describes an arbitrary value with a timestamp. The timestamp indicates when the
value was obtained. In JSON encoding and decoding, acts as a transparent
proxy/reference/pointer to the inner value.
*/
type Timed[A any] struct {
	Val  A `role:"ref"`
	Inst time.Time
}

// True if timestamp is unset.
func (self Timed[_]) IsNull() bool { return self.Inst.IsZero() }

// Inverse of `.IsNull`.
func (self Timed[_]) IsNonNull() bool { return !self.Inst.IsZero() }

// Implement `Clearer`. Zeroes the receiver.
//go:noinline
func (self *Timed[A]) Clear() { Clear(self) }

// Implement `Getter`, returning the underlying value as-is.
func (self Timed[A]) Get() A { return self.Val }

/*
Implement `Setter`. Modifies the underlying value and sets the current
timestamp. The resulting state is considered non-null even if the value
is "zero".
*/
func (self *Timed[A]) Set(val A) {
	self.Val = val
	self.Inst = time.Now()
}

// Implement `Ptrer`, returning a pointer to the underlying value.
func (self *Timed[A]) Ptr() *A {
	if self == nil {
		return nil
	}
	return &self.Val
}

/*
Implement `json.Marshaler`. If `.IsNull`, returns a representation of JSON null.
Otherwise uses `json.Marshal` to encode the underlying value.
*/
//go:noinline
func (self Timed[A]) MarshalJSON() ([]byte, error) {
	return JsonBytesNullCatch[A](self)
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON null,
clears the receiver via `.Clear`. Otherwise uses `json.Unmarshaler` to decode
into the underlying value, and sets the current timestamp on success.
*/
//go:noinline
func (self *Timed[A]) UnmarshalJSON(src []byte) error {
	return self.with(JsonParseClearCatch[A](src, self))
}

// True if the timestamp is unset, or if timestamp + duration > now.
func (self Timed[_]) IsExpired(dur time.Duration) bool {
	return self.Inst.IsZero() || self.Inst.Add(dur).Before(time.Now())
}

// Inverse of `.IsExpired`.
func (self Timed[_]) IsLive(dur time.Duration) bool {
	return !self.IsExpired(dur)
}

func (self *Timed[_]) with(err error) error {
	if err != nil {
		Clear(&self.Inst)
	} else {
		self.Inst = time.Now()
	}
	return err
}
