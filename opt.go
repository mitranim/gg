package gg

import "database/sql/driver"

/*
Short for "optional value". Instantiates an optional with the given val.
The result is considered non-null.
*/
func OptVal[A any](val A) Opt[A] { return Opt[A]{val, true} }

/*
Shortcut for creating an optional from a given value and boolean indicating
validity. If the boolean is false, the output is considered "null" even if the
value is not "zero".
*/
func OptFrom[A any](val A, ok bool) Opt[A] { return Opt[A]{val, ok} }

/*
Short for "optional". Wraps an arbitrary type. When `.Ok` is false, the value is
considered empty/null in various contexts such as text encoding, JSON encoding,
SQL encoding, even if the value is not "zero".
*/
type Opt[A any] struct {
	Val A
	Ok  bool
}

// Implement `Nullable`. True if not `.Ok`.
func (self Opt[_]) IsNull() bool { return !self.Ok }

// Inverse of `.IsNull`.
func (self Opt[_]) IsNotNull() bool { return self.Ok }

// Implement `Clearer`. Zeroes the receiver.
func (self *Opt[_]) Clear() { PtrClear(self) }

// Implement `Getter`, returning the underlying value as-is.
func (self Opt[A]) Get() A { return self.Val }

/*
Implement `Setter`. Modifies the underlying value and sets `.Ok = true`.
The resulting state is considered non-null even if the value is "zero".
*/
func (self *Opt[A]) Set(val A) {
	self.Val = val
	self.Ok = true
}

// Implement `Ptrer`, returning a pointer to the underlying value.
func (self *Opt[A]) Ptr() *A {
	if self == nil {
		return nil
	}
	return &self.Val
}

/*
Implement `fmt.Stringer`. If `.IsNull`, returns an empty string. Otherwise uses
the `String` function to encode the inner value.
*/
func (self Opt[A]) String() string { return StringNull[A](self) }

/*
Implement `Parser`. If the input is empty, clears the receiver via `.Clear`.
Otherwise uses the `ParseCatch` function, decoding into the underlying value.
*/
func (self *Opt[A]) Parse(src string) error {
	return self.with(ParseClearCatch[A](src, self))
}

// Implement `AppenderTo`, appending the same representation as `.String`.
func (self Opt[A]) AppendTo(buf []byte) []byte { return AppendNull[A](buf, self) }

// Implement `encoding.TextMarshaler`, returning the same representation as `.String`.
func (self Opt[A]) MarshalText() ([]byte, error) { return MarshalNullCatch[A](self) }

// Implement `encoding.TextUnmarshaler`, using the same logic as `.Parse`.
func (self *Opt[A]) UnmarshalText(src []byte) error {
	if len(src) <= 0 {
		self.Clear()
		return nil
	}
	return self.with(ParseClearCatch[A](src, self))
}

/*
Implement `json.Marshaler`. If `.IsNull`, returns a representation of JSON null.
Otherwise uses `json.Marshal` to encode the underlying value.
*/
func (self Opt[A]) MarshalJSON() ([]byte, error) {
	return JsonBytesNullCatch[A](self)
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON null,
clears the receiver via `.Clear`. Otherwise uses `JsonParseCatch` to decode
into the underlying value.
*/
func (self *Opt[A]) UnmarshalJSON(src []byte) error {
	if IsJsonEmpty(src) {
		self.Clear()
		return nil
	}
	return self.with(JsonParseCatch(src, &self.Val))
}

/*
Implement SQL `driver.Valuer`. If `.IsNull`, returns nil. If the underlying
value implements `driver.Valuer`, delegates to its method. Otherwise returns
the underlying value as-is.
*/
func (self Opt[A]) Value() (driver.Value, error) { return ValueNull[A](self) }

/*
Implement SQL `Scanner`, decoding an arbitrary input into the underlying value.
If the underlying type implements `Scanner`, delegates to that implementation.
Otherwise input must be nil or text-like (see `Text`). Text decoding uses the
same logic as `.Parse`.
*/
func (self *Opt[A]) Scan(src any) error {
	if src == nil {
		self.Clear()
		return nil
	}

	val, ok := src.(A)
	if ok {
		self.Set(val)
		return nil
	}

	return self.with(ScanCatch[A](src, self))
}

func (self *Opt[_]) with(err error) error {
	self.Ok = err == nil
	return err
}

/*
FP-style "mapping". If the original value is considered "null", or if the
function is nil, the output is "zero" and "null". Otherwise the output is the
result of calling the function with the previous value, and is considered
non-"null" even if the value is zero.
*/
func OptMap[A, B any](src Opt[A], fun func(A) B) (out Opt[B]) {
	if src.IsNotNull() && fun != nil {
		out.Set(fun(src.Val))
	}
	return
}
