package gg

/*
Short for "zero optional value". Workaround for the lack of type inference in
type literals.
*/
func ZopVal[A any](val A) Zop[A] { return Zop[A]{val} }

/*
Short for "zero optional". The zero value is considered empty/null in JSON. Note
that "encoding/json" doesn't support ",omitempty" for struct values. This
wrapper allows empty structs to become "null". This type doesn't implement any
other encoding or decoding methods, and is intended only for non-scalar values
such as "models" / "data classes".
*/
type Zop[A any] struct {
	/**
	The `role:"ref"` annotation, where "ref" is short for "reference", indicates
	that this field is a reference/pointer to the inner type/value.
	Reflection-based code may use this to treat `Zop` like a pointer.
	*/
	Val A `role:"ref"`
}

// Implement `Nullable`. True if zero value of its type.
func (self Zop[_]) IsNull() bool { return IsZero(self.Val) }

// Inverse of `.IsNull`.
func (self Zop[_]) IsNonNull() bool { return !IsZero(self.Val) }

// Implement `Clearer`. Zeroes the receiver.
func (self *Zop[_]) Clear() { Clear(&self.Val) }

/*
Implement `Getter` for compatibility with 3rd party libraries such as `pgx`.
If `.IsNull`, returns nil. Otherwise returns the underlying value, invoking
its own `Getter` if possible.
*/
//go:noinline
func (self Zop[A]) Get() any { return GetNull[A](self) }

// Implement `ValGetter`, returning the underlying value as-is.
func (self Zop[A]) GetVal() A { return self.Val }

// Implement `ValSetter`, modifying the underlying value.
func (self *Zop[A]) SetVal(val A) { self.Val = val }

// Implement `PtrGetter`, returning a pointer to the underlying value.
func (self *Zop[A]) GetPtr() *A {
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
func (self Zop[A]) MarshalJSON() ([]byte, error) { return JsonBytesNullCatch[A](self) }

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON null,
clears the receiver via `.Clear`. Otherwise uses `json.Unmarshaler` to decode
into the underlying value.
*/
//go:noinline
func (self *Zop[A]) UnmarshalJSON(src []byte) error { return JsonParseClearCatch[A](src, self) }

/*
FP-style "mapping". If the original value is zero, or if the function is nil,
the output is zero. Otherwise the output is the result of calling the function
with the previous value.
*/
func ZopMap[A, B any](src Zop[A], fun func(A) B) (out Zop[B]) {
	if src.IsNonNull() && fun != nil {
		out.Val = fun(src.Val)
	}
	return
}
