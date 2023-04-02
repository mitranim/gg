package gg

import (
	"context"
	"database/sql/driver"
	r "reflect"
	u "unsafe"
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

	/**
	When possible, we want to "natively" compare to a zero value before invoking
	`Zeroable`. Depending on the Go version, this may slightly speed us up, or
	slightly slow us down. More importantly, this improves correctness,
	safeguarding us against bizarre implementations of `Zeroable` such as
	`reflect.Value.IsZero`, which panics when called on the zero value of
	`reflect.Value`.

	It would be ideal to compare concrete values rather than `any`, but requiring
	`comparable` would make this function much less usable, and comparing `any`
	seems fast enough.
	*/
	if Type[A]().Comparable() {
		if box == AnyNoEscUnsafe(Zero[A]()) {
			return true
		}

		/**
		Terminate here to avoid falling down to the reflection-based clause.
		We already know that our value is not considered zero by Go.
		*/
		impl, _ := box.(Zeroable)
		return impl != nil && impl.IsZero()
	}

	impl, _ := box.(Zeroable)
	if impl != nil {
		return impl.IsZero()
	}

	/**
	Implementation note. For some types, some non-zero byte patterns are
	considered a zero value. The most notable example is strings. For example,
	the following expression makes a string with a non-zero data pointer, which
	is nevertheless considered a zero value by Go:

		`some_text`[:0]

	`reflect.Value.IsZero` correctly handles such cases, and doesn't seem to be
	outrageously slow.
	*/
	return r.ValueOf(box).IsZero()
}

// Inverse of `IsZero`.
func IsNotZero[A any](val A) bool { return !IsZero(val) }

/*
True if every byte in the given value is zero. Not equivalent to `IsZero`.
Most of the time, you should prefer `IsZero`, which is more performant and
more correct.
*/
func IsTrueZero[A any](val A) bool {
	size := u.Sizeof(val)
	for off := uintptr(0); off < size; off++ {
		if *(*byte)(u.Pointer(uintptr(u.Pointer(&val)) + off)) != 0 {
			return false
		}
	}
	return true
}

// Returns a zero value of the given type.
func Zero[A any]() (_ A) { return }

// Generic variant of `Nullable.IsNull`.
func IsNull[A Nullable](val A) bool { return val.IsNull() }

// Inverse of `IsNull`.
func IsNotNull[A Nullable](val A) bool { return !val.IsNull() }

/*
Zeroes the memory referenced by the given pointer. If the pointer is nil, does
nothing. Also see the interface `Clearer` and method `.Clear` implemented by
various types.
*/
func PtrClear[A any](val *A) {
	if val != nil {
		*val = Zero[A]()
	}
}

// Calls `PtrClear` and returns the same pointer.
func PtrCleared[A any](val *A) *A {
	PtrClear(val)
	return val
}

/*
Takes an arbitrary value and returns a non-nil pointer to a new memory region
containing a shallow copy of that value.
*/
func Ptr[A any](val A) *A { return &val }

// If the pointer is non-nil, dereferences it. Otherwise returns zero value.
func PtrGet[A any](val *A) A {
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
If the pointer is non-nil, returns its value while zeroing the destination.
Otherwise returns zero value.
*/
func PtrPop[A any](src *A) (out A) {
	if src != nil {
		out, *src = *src, out
	}
	return
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

/*
True if the inputs are byte-for-byte identical. This function is not meant for
common use. Nearly always, you should use `Eq` or `Equal` instead. This one is
sometimes useful for testing purposes, such as asserting that two interface
values refer to the same underlying data. This may lead to brittle code that is
not portable between different Go implementations. Performance is similar to
`==` for small value sizes (up to 2-4 machine words) but is significantly worse
for large value sizes.
*/
func Is[A any](one, two A) bool {
	/**
	Note. The "ideal" implementation looks like this:

		const size = u.Sizeof(one)
		return CastUnsafe[[size]byte](one) == CastUnsafe[[size]byte](two)

	But at the time of writing, in Go 1.19, `unsafe.Sizeof` on a type parameter
	is considered non-constant. If this changes in the future, we'll switch to
	the implementation above.
	*/

	size := u.Sizeof(one)

	switch size {
	case 0:
		return true

	case SizeofWord:
		return *(*uint)(u.Pointer(&one)) == *(*uint)(u.Pointer(&two))

	// Common case: comparing interfaces or strings.
	case SizeofWord * 2:
		return (*(*uint)(u.Pointer(&one)) == *(*uint)(u.Pointer(&two))) &&
			(*(*uint)(u.Pointer(uintptr(u.Pointer(&one)) + SizeofWord)) ==
				*(*uint)(u.Pointer(uintptr(u.Pointer(&two)) + SizeofWord)))

	// Common case: comparing slices.
	case SizeofWord * 3:
		return (*(*uint)(u.Pointer(&one)) == *(*uint)(u.Pointer(&two))) &&
			(*(*uint)(u.Pointer(uintptr(u.Pointer(&one)) + SizeofWord)) ==
				*(*uint)(u.Pointer(uintptr(u.Pointer(&two)) + SizeofWord))) &&
			(*(*uint)(u.Pointer(uintptr(u.Pointer(&one)) + SizeofWord*2)) ==
				*(*uint)(u.Pointer(uintptr(u.Pointer(&two)) + SizeofWord*2)))

	default:
		/**
		Implementation note. We could also walk word-by-word by using padded structs
		to ensure sufficient empty memory. It would improve the performance
		slightly, but not enough to bother. The resulting performance is still
		much worse than `==` on whole values.
		*/
		for off := uintptr(0); off < size; off++ {
			oneChunk := *(*byte)(u.Pointer(uintptr(u.Pointer(&one)) + off))
			twoChunk := *(*byte)(u.Pointer(uintptr(u.Pointer(&two)) + off))
			if oneChunk != twoChunk {
				return false
			}
		}
		return true
	}
}

// Same as `==`. Sometimes useful with higher-order functions.
func Eq[A comparable](one, two A) bool { return one == two }

/*
Short for "equal". For types that implement `Equaler`, this simply calls their
equality method. Otherwise falls back on `reflect.DeepEqual`. Compared to
`reflect.DeepEqual`, this has better type safety, and in many cases this has
better performance, even when calling `reflect.DeepEqual` in fallback mode.
*/
func Equal[A any](one, two A) bool {
	impl, _ := AnyNoEscUnsafe(one).(Equaler[A])
	if impl != nil {
		return impl.Equal(two)
	}
	return r.DeepEqual(AnyNoEscUnsafe(one), AnyNoEscUnsafe(two))
}

/*
True if the inputs are equal via `==`, and neither is a zero value of its type.
For non-equality, use `NotEqNotZero`.
*/
func EqNotZero[A comparable](one, two A) bool {
	return one == two && one != Zero[A]()
}

/*
True if the inputs are non-equal via `!=`, and at least one is not a zero value
of its type. For equality, use `EqNotZero`.
*/
func NotEqNotZero[A comparable](one, two A) bool {
	return one != two && one != Zero[A]()
}

/*
True if the given slices have the same data pointer, length, capacity.
Does not compare individual elements.
*/
func SliceIs[A any](one, two []A) bool {
	return u.SliceData(one) == u.SliceData(two) &&
		len(one) == len(two) &&
		cap(one) == cap(two)
}

// Returns the first non-zero value from among the inputs.
func Or[A any](val ...A) A { return Find(val, IsNotZero[A]) }

/*
Variant of `Or` compatible with `Nullable`. Returns the first non-"null" value
from among the inputs.
*/
func NullOr[A Nullable](val ...A) A { return Find(val, IsNotNull[A]) }

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
Non-asserting interface conversion. Converts the given `any` into the given
type, returning zero value on failure.
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

	defer PtrSwap(&somePtr, someVal).Done()
*/
func PtrSwap[A any](ptr *A, next A) Snapshot[A] {
	prev := *ptr
	*ptr = next
	return Snapshot[A]{ptr, prev}
}

// Short for "snapshot". Used by `PtrSwap`.
type Snapshot[A any] struct {
	Ptr *A
	Val A
}

// If the pointer is non-nil, writes the value to it. See `PtrSwap`.
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
