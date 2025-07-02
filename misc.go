package gg

import (
	"context"
	"database/sql/driver"
)

var Indent = `    `

// Does nothing.
func Nop() {}

// Does nothing.
func Nop1[A any](A) {}

// Does nothing.
func Nop2[A, B any](A, B) {}

// Does nothing.
func Nop3[A, B, C any](A, B, C) {}

// Identity function. Returns input as-is.
func Id1[A any](val A) A { return val }

// Identity function. Returns input as-is.
func Id2[A, B any](val0 A, val1 B) (A, B) { return val0, val1 }

// Identity function. Returns input as-is.
func Id3[A, B, C any](val0 A, val1 B, val2 C) (A, B, C) { return val0, val1, val2 }

// Returns a zero value of the given type.
func Zero[A any]() (_ A) { return }

/*
Same as Go's `+` operator, expressed as a generic function. Input type may be
numeric or ~string. When the input type is numeric, this is unchecked and may
overflow. For integers, prefer [Add] whenever possible, which has overflow
checks.
*/
func Plus2[A Plusable](one, two A) A { return one + two }

/*
Variadic version of Go's `+` operator. Input type may be numeric or ~string.
If the input is empty, returns a zero value. Use caution: this has no overflow
checks for numbers. Prefer [Add] for integers. See [Sum] for a non-variadic
equivalent.
*/
func Plus[A Plusable](val ...A) A { return Foldz(val, Plus2[A]) }

/*
Shortcut for implementing `driver.Valuer` on `Nullable` types that wrap other
types, such as `Opt`. Mostly for internal use.
*/
func ValueNull[A any, B NullableValGetter[A]](src B) (driver.Value, error) {
	if src.IsNull() {
		return nil, nil
	}

	val := src.Get()

	impl, _ := AnyNoEscUnsafe(val).(driver.Valuer)
	if impl != nil {
		return impl.Value()
	}

	return val, nil
}

/*
Equivalent to the following, but more convenient:

	_, ok := src.(A)
	_ = ok
*/
func AnyIs[A any](src any) bool {
	_, ok := AnyNoEscUnsafe(src).(A)
	return ok
}

/*
Equivalent to the following, but more convenient. When conversion can't be
performed, the output is a zero value of `A`.

	val, _ := src.(A)
	_ = val
*/
func AnyAs[A any](src any) A {
	val, _ := AnyNoEscUnsafe(src).(A)
	return val
}

/*
Converts the argument to `any` and returns it. Sometimes useful in higher-order
functions.
*/
func ToAny[A any](val A) any { return val }

/*
Uses `context.WithValue` to create a context with the given value, using the
type's nil pointer "(*A)(nil)" as the key.
*/
func CtxSet[A any](ctx context.Context, val A) context.Context {
	return context.WithValue(ctx, (*A)(nil), val)
}

/*
Uses `ctx.Value` to get the value of the given type, using the type's nil pointer
"(*A)(nil)" as the key. If the context is nil or doesn't contain the value,
returns zero value and false.
*/
func CtxGot[A any](ctx context.Context) (A, bool) {
	if ctx == nil {
		return Zero[A](), false
	}
	val, ok := ctx.Value((*A)(nil)).(A)
	return val, ok
}

// Same as `CtxGot` but returns only the boolean.
func CtxHas[A any](ctx context.Context) bool {
	_, ok := CtxGot[A](ctx)
	return ok
}

/*
Same as `CtxGot` but returns only the resulting value. If value was not found,
output is zero.
*/
func CtxGet[A any](ctx context.Context) A {
	val, _ := CtxGot[A](ctx)
	return val
}

/*
Short for "iterator". Returns a slice of the given length that can be iterated
by using a `range` loop. Usage:

	for range Iter(size) { ... }
	for ind := range Iter(size) { ... }

Because `struct{}` is zero-sized, `[]struct{}` is backed by "zerobase" (see Go
source → "runtime/malloc.go") and does not allocate. Loops using this should
compile to approximately the same instructions as "normal" counted loops.

This function is unnecessary in Go 1.22 and higher, where the `for` loop
has built-in support for ranging over an integer.
*/
func Iter(size int) []struct{} { return make([]struct{}, size) }

/*
Returns a slice of numbers from `min` to `max`. The range is inclusive at the
start but exclusive at the end: `[min,max)`. If `!(max > min)`, returns nil.
Values must be within the range of the Go type `int`.
*/
func Range[A Int](min, max A) []A {
	// We must check this before calling `max-1` to avoid underflow.
	if !(max > min) {
		return nil
	}
	return RangeIncl(min, max-1)
}

/*
Returns a slice of numbers from `min` to `max`. The range is inclusive at the
start and at the end: `[min,max]`. If `!(max >= min)`, returns nil. Values must
be within the range of the Go type `int`.

While the exclusive range `[min,max)` implemented by `Range` is more
traditional, this function allows to create a range that includes the maximum
representable value of any given integer type, such as 255 for `uint8`, which
cannot be done with `Range`.
*/
func RangeIncl[A Int](min, max A) []A {
	if !(max >= min) {
		return nil
	}

	minInt := NumConv[int](min)
	maxInt := NumConv[int](max)
	buf := make([]A, (maxInt-minInt)+1)
	for ind := range buf {
		buf[ind] = A(ind + minInt)
	}
	return buf
}

// Shortcut for creating range `[0,N)`, exclusive at the end.
func Span[A Int](val A) []A { return Range(0, val) }

/*
Takes a pointer and a fallback value which must be non-zero. If the pointer
destination is zero, sets the fallback and returns true. Otherwise returns
false.
*/
func Fellback[A any](tar *A, fallback A) bool {
	if IsZero(fallback) {
		panic(Errf(`invalid non-zero fallback %#v`, fallback))
	}

	if tar == nil {
		return false
	}

	if IsZero(*tar) {
		*tar = fallback
		return true
	}
	return false
}

/*
Snapshots the current value at the given pointer and returns a snapshot
that can restore this value. Usage:

	defer Snap(&somePtr).Done()
	somePtr.SomeField = someValue
*/
func Snap[A any](ptr *A) Snapshot[A] { return Snapshot[A]{ptr, *ptr} }

/*
Snapshots the previous value, sets the next value, and returns a snapshot
that can restore the previous value. Usage:

	defer SnapSwap(&somePtr, someVal).Done()
*/
func SnapSwap[A any](ptr *A, next A) Snapshot[A] {
	prev := *ptr
	*ptr = next
	return Snapshot[A]{ptr, prev}
}

// Short for "snapshot". Used by `SnapSwap`.
type Snapshot[A any] struct {
	Ptr *A
	Val A
}

// If the pointer is non-nil, writes the value to it. See `SnapSwap`.
func (self Snapshot[_]) Done() {
	if self.Ptr != nil {
		*self.Ptr = self.Val
	}
}

/*
Snapshots the length of the given slice and returns a snapshot that can restore
the previous length. Usage:

	defer SnapSlice(&somePtr).Done()
*/
func SnapSlice[Slice ~[]Elem, Elem any](ptr *Slice) SliceSnapshot[Slice, Elem] {
	return SliceSnapshot[Slice, Elem]{ptr, PtrLen(ptr)}
}

/*
Analogous to `Snapshot`, but instead of storing a value, stores a length.
When done, reverts the referenced slice to the given length.
*/
type SliceSnapshot[Slice ~[]Elem, Elem any] struct {
	Ptr *Slice
	Len int
}

/*
Analogous to `Snapshot.Done`. Reverts the referenced slice to `self.Len` while
keeping the capacity.
*/
func (self SliceSnapshot[_, _]) Done() {
	if self.Ptr != nil {
		*self.Ptr = (*self.Ptr)[:self.Len]
	}
}

// Shortcut for making a pseudo-tuple with two elements.
func Tuple2[A, B any](valA A, valB B) Tup2[A, B] {
	return Tup2[A, B]{valA, valB}
}

// Represents a pseudo-tuple with two elements.
type Tup2[A, B any] struct {
	A A
	B B
}

// Converts the pseudo-tuple to a proper Go tuple.
func (self Tup2[A, B]) Get() (A, B) { return self.A, self.B }

// Shortcut for making a pseudo-tuple with three elements.
func Tuple3[A, B, C any](valA A, valB B, valC C) Tup3[A, B, C] {
	return Tup3[A, B, C]{valA, valB, valC}
}

// Represents a pseudo-tuple with three elements.
type Tup3[A, B, C any] struct {
	A A
	B B
	C C
}

// Converts the pseudo-tuple to a proper Go tuple.
func (self Tup3[A, B, C]) Get() (A, B, C) { return self.A, self.B, self.C }

/*
Makes a zero value of the given type, passes it to the given mutator functions
by pointer, and returns the modified value. Nil functions are ignored.
*/
func With[A any](funs ...func(*A)) (out A) {
	for _, fun := range funs {
		if fun != nil {
			fun(&out)
		}
	}
	return
}
