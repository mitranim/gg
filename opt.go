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
func (self Opt[A]) IsNull() bool { return !self.Ok }

// Inverse of `.IsNull`.
func (self Opt[A]) IsNonNull() bool { return self.Ok }

// Implement `Clearer`. Zeroes the receiver.
//go:noinline
func (self *Opt[A]) Clear() { Clear(self) }

/*
Implement `Getter` for compatibility with 3rd party libraries such as `pgx`.
If `.IsNull`, returns nil. Otherwise returns the underlying value, invoking
its own `Getter` if possible.
*/
//go:noinline
func (self Opt[A]) Get() any { return GetNull[A](self) }

// Implement `ValGetter`, returning the underlying value as-is.
func (self Opt[A]) GetVal() A { return self.Val }

/*
Implement `ValSetter`. Modifies the underlying value and sets `.Ok = true`.
The result is considered non-null even if the value is "zero".
*/
func (self *Opt[A]) SetVal(val A) {
	self.Val = val
	self.Ok = true
}

// Implement `PtrGetter`, returning a pointer to the underlying value.
func (self *Opt[A]) GetPtr() *A {
	if self == nil {
		return nil
	}
	return &self.Val
}

/*
Implement `fmt.Stringer`. If `.IsNull`, returns an empty string. Otherwise uses
the `String` function to encode the inner value.
*/
//go:noinline
func (self Opt[A]) String() string { return StringNull[A](self) }

/*
Implement `Parser`. If the input is empty, clears the receiver via `.Clear`.
Otherwise uses the `ParseCatch` function, decoding into the underlying value.
*/
//go:noinline
func (self *Opt[A]) Parse(src string) error {
	return self.with(ParseClearCatch[A](src, self))
}

// Implement `Appender`, appending the same representation as `.String`.
//go:noinline
func (self Opt[A]) Append(buf []byte) []byte { return AppendNull[A](buf, self) }

// Implement `encoding.TextMarshaler`, returning the same representation as `.String`.
//go:noinline
func (self Opt[A]) MarshalText() ([]byte, error) { return MarshalNullCatch[A](self) }

// Implement `encoding.TextUnmarshaler`, using the same logic as `.Parse`.
//go:noinline
func (self *Opt[A]) UnmarshalText(src []byte) error {
	return self.with(ParseClearCatch[A](src, self))
}

/*
Implement `json.Marshaler`. If `.IsNull`, returns a representation of JSON null.
Otherwise uses `json.Marshal` to encode the underlying value.
*/
//go:noinline
func (self Opt[A]) MarshalJSON() ([]byte, error) {
	return JsonBytesNullCatch[A](self)
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON null,
clears the receiver via `.Clear`. Otherwise uses `json.Unmarshaler` to decode
into the underlying value.
*/
//go:noinline
func (self *Opt[A]) UnmarshalJSON(src []byte) error {
	return self.with(JsonParseClearCatch[A](src, self))
}

/*
Implement SQL `driver.Valuer`. If `.IsNull`, returns nil. If the underlying
value implements `driver.Valuer`, delegates to its method. Otherwise returns
the underlying value, invoking its own `Getter` if possible.
*/
//go:noinline
func (self Opt[A]) Value() (driver.Value, error) { return ValueNull[A](self) }

/*
Implement SQL `Scanner`, decoding an arbitrary input into the underlying value.
If the underlying type implements `Scanner`, its own implementation is used.
Otherwise input must be nil or text-like (see `Text`). Text decoding uses the
same logic as `.Parse`.
*/
//go:noinline
func (self *Opt[A]) Scan(src any) error {
	self.Clear()
	return self.with(ScanCatch[A](GetOpt(src), self))
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
	if src.IsNonNull() && fun != nil {
		out.SetVal(fun(src.Val))
	}
	return
}
