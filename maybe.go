package gg

import "encoding/json"

// Shortcut for creating a `Maybe` with the given value.
func MaybeVal[A any](val A) Maybe[A] { return Maybe[A]{val, nil} }

// Shortcut for creating a `Maybe` with the given error.
func MaybeErr[A any](err error) Maybe[A] { return Maybe[A]{Zero[A](), err} }

/*
Contains a value or an error. The JSON tags "value" and "error" are chosen due
to their existing popularity in HTTP API.
*/
type Maybe[A any] struct {
	Val A     `json:"value,omitempty"`
	Err error `json:"error,omitempty"`
}

/*
Asserts that the error is nil, returning the resulting value. If the error is
non-nil, panics via `Try`, idempotently adding a stack trace to the error.
*/
func (self Maybe[A]) Ok() A {
	Try(self.Err)
	return self.Val
}

// Implement `Getter`, returning the underlying value as-is.
func (self Maybe[A]) Get() A { return self.Val }

// Implement `Setter`. Sets the underlying value and clears the error.
func (self *Maybe[A]) Set(val A) {
	self.Val = val
	self.Err = nil
}

// Returns the underlying error as-is.
func (self Maybe[_]) GetErr() error { return self.Err }

// Sets the error. If the error is non-nil, clears the value.
func (self *Maybe[A]) SetErr(err error) {
	if err != nil {
		self.Val = Zero[A]()
	}
	self.Err = err
}

// True if error is non-nil.
func (self Maybe[_]) HasErr() bool { return self.Err != nil }

/*
Implement `json.Marshaler`. If the underlying error is non-nil, returns that
error. Otherwise uses `json.Marshal` to encode the underlying value.
*/
func (self Maybe[_]) MarshalJSON() ([]byte, error) {
	if self.Err != nil {
		return nil, self.Err
	}
	return json.Marshal(self.Val)
}

// Implement `json.Unmarshaler`, decoding into the underlying value.
func (self *Maybe[_]) UnmarshalJSON(src []byte) error {
	self.Err = nil
	return json.Unmarshal(src, &self.Val)
}
