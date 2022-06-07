package gg

/*
Short for "zero optional value". Workaround for the lack of type inference in
struct literals.
*/
func ZopVal[A any](val A) Zop[A] { return Zop[A]{val} }

/*
Short for "zero optional". The zero value is considered empty/null in JSON. Note
that "encoding/json" doesn't support ",omitempty" for structs. This wrapper
allows empty structs to become "null". This type doesn't implement any other
encoding or decoding methods, and is intended only for non-scalar values such
as "models" / "data classes".
*/
type Zop[A any] struct {
	/**
	Annotation `role:"ref"` indicates that this field is a reference/pointer to
	the inner type/value. Reflection-based code may use this to treat this type
	like a pointer.
	*/
	Val A `role:"ref"`
}

// Implement `Nullable`. True if zero value of its type.
func (self Zop[_]) IsNull() bool { return IsZero(self.Val) }

// Inverse of `.IsNull`.
func (self Zop[_]) IsNonNull() bool { return !IsZero(self.Val) }

// Implement `Clearer`. Zeroes the receiver.
func (self *Zop[_]) Clear() { Clear(&self.Val) }

// Implement `Getter`, returning the underlying value as-is.
func (self Zop[A]) Get() A { return self.Val }

// Implement `Setter`, modifying the underlying value.
func (self *Zop[A]) Set(val A) { self.Val = val }

// Implement `Ptrer`, returning a pointer to the underlying value.
func (self *Zop[A]) Ptr() *A {
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
func (self Zop[A]) MarshalJSON() ([]byte, error) {
	return JsonBytesNullCatch[A](self)
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON null,
clears the receiver via `.Clear`. Otherwise uses `JsonParseCatch` to decode
into the underlying value.
*/
//go:noinline
func (self *Zop[A]) UnmarshalJSON(src []byte) error {
	if IsJsonEmpty(src) {
		self.Clear()
		return nil
	}
	return JsonParseCatch(src, &self.Val)
}

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
