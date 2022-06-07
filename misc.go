package gg

import (
	"context"
	"database/sql/driver"
	r "reflect"
)

var (
	Indent  = `    `
	Space   = ` `
	Newline = "\n"
)

/*
Returns true if the input is the zero value of its type. Optionally falls back
on `Zeroable.IsZero` if the input implements this interface. Support for
`Zeroable` is useful for types such as `time.Time`, where "zero" is determined
only by the timestamp, ignoring the timezone.
*/
func IsZero[A any](val A) bool {
	box := AnyNoEscUnsafe(val)
	if box == nil {
		return true
	}

	// Prioritize `==` over reflection if possible. This is measurably faster than
	// the reflect-based version, especially for large value types such as fat
	// structs or arrays,
	if Type[A]().Comparable() {
		// True zero value must always be considered zero, even if the type
		// implements `Zeroable`. More importantly, this safeguards us against
		// unusual cases such as `reflect.Value.IsZero`, which panics when called
		// on the zero value of `reflect.Value`.
		if box == AnyNoEscUnsafe(Zero[A]()) {
			return true
		}

		impl, _ := box.(Zeroable)
		return impl != nil && impl.IsZero()
	}

	impl, _ := box.(Zeroable)
	if impl != nil {
		return impl.IsZero()
	}
	return r.ValueOf(box).IsZero()
}

// Inverse of `IsZero`.
func IsNonZero[A any](val A) bool { return !IsZero(val) }

// Returns a zero value of the given type.
func Zero[A any]() (_ A) { return }

// Generic variant of `Nullable.IsNull`.
func IsNull[A Nullable](val A) bool { return val.IsNull() }

// Inverse of `IsNull`.
func IsNonNull[A Nullable](val A) bool { return !val.IsNull() }

/*
Zeroes the memory referenced by the given pointer. If the pointer is nil, this
is a nop.
*/
func Clear[A any](val *A) {
	if val != nil {
		*val = Zero[A]()
	}
}

/*
Takes an arbitrary value and returns a non-nil pointer to a new memory region
containing a shallow copy of that value.
*/
func Ptr[A any](val A) *A { return &val }

/*
If the pointer is non-nil, dereferences it. Otherwise returns zero value.
TODO consider renaming to `PtrGet`.
*/
func Deref[A any](val *A) A {
	if val != nil {
		return *val
	}
	return Zero[A]()
}

// If the pointer is nil, does nothing. If non-nil, set the given value.
func PtrSet[A any](tar *A, val A) {
	if tar != nil {
		*tar = val
	}
}

/*
Takes two pointers and copies the value from source to target if both pointers
are non-nil. If either is nil, does nothing.
*/
func PtrSetOpt[A any](tar, src *A) {
	if tar != nil && src != nil {
		*tar = *src
	}
}

/*
If the pointer is nil, uses `new` to allocate a new value of the given type,
returning the resulting pointer. Otherwise returns the input as-is.
*/
func PtrInited[A any](val *A) *A {
	if val != nil {
		return val
	}
	return new(A)
}

/*
If the outer pointer is nil, returns nil. If the inner pointer is nil, uses
`new` to allocate a new value, sets and returns the resulting new pointer.
Otherwise returns the inner pointer as-is.
*/
func PtrInit[A any](val **A) *A {
	if val == nil {
		return nil
	}
	if *val == nil {
		*val = new(A)
	}
	return *val
}

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
Returns the lesser of the two inputs, which must be comparable primitives. For
non-primitives, see `Min2`.
*/
func MinPrim2[A LesserPrim](one, two A) A {
	if one < two {
		return one
	}
	return two
}

/*
Returns the lesser of the two inputs. For primitive types that don't implement
`Lesser`, see `MinPrim2`.
*/
func Min2[A Lesser[A]](one, two A) A {
	if one.Less(two) {
		return one
	}
	return two
}

/*
Returns the larger of the two inputs, which must be comparable primitives. For
non-primitives, see `Max2`.
*/
func MaxPrim2[A LesserPrim](one, two A) A {
	if one < two {
		return two
	}
	return one
}

/*
Returns the larger of the two inputs. For primitive types that don't implement
`Lesser`, see `MaxPrim2`.
*/
func Max2[A Lesser[A]](one, two A) A {
	if one.Less(two) {
		return two
	}
	return one
}

// True if the given number is > 0.
func IsPos[A Signed](val A) bool { return val > 0 }

// True if the given number is < 0.
func IsNeg[A Signed](val A) bool { return val < 0 }

// Same as `==`. Sometimes useful with higher-order functions.
func Eq[A comparable](one, two A) bool { return one == two }

/*
Short for "equal". For types that implement `Equaler`, this simply calls their
equality method. Otherwise falls back on `reflect.DeepEqual`. Compared to
`reflect.DeepEqual`, this has better type safety and performance, even when
calling it in fallback mode.
*/
func Equal[A any](one, two A) bool {
	impl, _ := AnyNoEscUnsafe(one).(Equaler[A])
	if impl != nil {
		return impl.Equal(two)
	}
	return r.DeepEqual(AnyNoEscUnsafe(one), AnyNoEscUnsafe(two))
}

/*
True if both inputs are not zero values of their type, and are equal to each
other via `==`.
*/
func EqNonZero[A comparable](one, two A) bool {
	return one != Zero[A]() && one == two
}

/*
True if the given slices have the same data pointer, length, capacity.
Does not compare individual elements.
*/
func SliceIs[A any](one, two []A) bool {
	return CastUnsafe[r.SliceHeader](one) == CastUnsafe[r.SliceHeader](two)
}

// Returns the first non-zero value from among the inputs.
func Or[A any](val ...A) A { return Find(val, IsNonZero[A]) }

/*
Variant of `Or` compatible with `Nullable`. Returns the first non-"null" value
from among the inputs.
*/
func NullOr[A Nullable](val ...A) A { return Find(val, IsNonNull[A]) }

/*
Returns true if the given `any` can be usefully converted into a value of the
given type. If the result is true, `src.(A)` doesn't panic. If the output is
false, `src.(A)` panics.
*/
func AnyIs[A any](src any) bool {
	_, ok := AnyNoEscUnsafe(src).(A)
	return ok
}

/*
Non-asserting interface conversion. Safely converts the given `any` into the
given type, returning zero value on failure.
*/
func AnyAs[A any](src any) A {
	val, _ := AnyNoEscUnsafe(src).(A)
	return val
}

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
source → "runtime/malloc.go") and does not allocate. The example loops should
compile to approximately the same instructions as "normal" counted loops.
*/
func Iter(size int) []struct{} { return make([]struct{}, size) }

/*
Returns a slice of numbers from "min" to "max". The range is inclusive at the
start but exclusive at the end: "[min,max)".
*/
func Range[A Int](min, max A) []A {
	buf := make([]A, max-min)
	for ind := range buf {
		buf[ind] = A(ind) + min
	}
	return buf
}

// Shortcut for creating a range from 0 to N.
func Span[A Int](val A) []A { return Range(0, val) }

// Combines two inputs via "+". Also see variadic `Plus`.
func Plus2[A Plusable](one, two A) A { return one + two }

// Same as `one - two`.
func Minus2[A Num](one, two A) A { return one - two }

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

	defer Swap(&somePtr, someVal).Done()
*/
func Swap[A any](ptr *A, next A) Snapshot[A] {
	prev := *ptr
	*ptr = next
	return Snapshot[A]{ptr, prev}
}

// Short for "snapshot". Used by `Swap`.
type Snapshot[A any] struct {
	Ptr *A
	Val A
}

// If the pointer is non-nil, writes the value to it. See `Swap`.
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
func SnapSlice[Slice ~[]Elem, Elem any](ptr *Slice) SliceSnapshot[Elem] {
	return SliceSnapshot[Elem]{CastUnsafe[*[]Elem](ptr), PtrLen(ptr)}
}

/*
Analogous to `Snapshot`, but instead of storing a value, stores a length.
When done, reverts the referenced slice to the given length.
*/
type SliceSnapshot[A any] struct {
	Ptr *[]A
	Len int
}

/*
Analogous to `Snapshot.Done`. Reverts the referenced slice to `self.Len` while
keeping the capacity.
*/
func (self SliceSnapshot[_]) Done() {
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
