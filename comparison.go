package gg

import (
	r "reflect"
	u "unsafe"
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

// Generic variant of `Nullable.IsNull`.
func IsNull[A Nullable](val A) bool { return val.IsNull() }

// Inverse of `IsNull`.
func IsNotNull[A Nullable](val A) bool { return !val.IsNull() }

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
True if the given slice headers are byte-for-byte identical. In other words,
true if the given slices have the same data pointer, length, capacity. Does not
compare individual elements.
*/
func SliceIs[A any](one, two []A) bool {
	return CastUnsafe[r.SliceHeader](one) == CastUnsafe[r.SliceHeader](two)
}

// Returns the first non-zero value from among the inputs.
func Or[A any](val ...A) A { return Find(val, IsNotZero[A]) }

/*
Variant of `Or` compatible with `Nullable`. Returns the first non-"null" value
from among the inputs.
*/
func NullOr[A Nullable](val ...A) A { return Find(val, IsNotNull[A]) }

// Version of `<` for non-primitives that implement `Lesser`.
func Less2[A Lesser[A]](one, two A) bool { return one.Less(two) }

/*
Variadic version of `<` for non-primitives that implement `Lesser`.
Shortcut for `IsSorted` with `Less2`.
*/
func Less[A Lesser[A]](src ...A) bool { return IsSorted(src, Less2[A]) }

// Same as Go's `<` operator, expressed as a generic function.
func LessPrim2[A LesserPrim](one, two A) bool { return one < two }

/*
Variadic version of `<` for non-primitives that implement `Lesser`.
Shortcut for `IsSorted` with `LessPrim2`.
*/
func LessPrim[A LesserPrim](src ...A) bool { return IsSorted(src, LessPrim2[A]) }

// Version of `<=` for non-primitives that implement `Lesser`.
func LessEq2[A Lesser[A]](one, two A) bool { return one.Less(two) || Equal(one, two) }

/*
Variadic version of `<=` for non-primitives that implement `Lesser`.
Shortcut for `IsSorted` with `LessEq2`.
*/
func LessEq[A Lesser[A]](src ...A) bool { return IsSorted(src, LessEq2[A]) }

// Same as Go's `<=` operator, expressed as a generic function.
func LessEqPrim2[A LesserPrim](one, two A) bool { return one <= two }

/*
Variadic version of Go's `<=` operator.
Shortcut for `IsSorted` with `LessEqPrim2`.
*/
func LessEqPrim[A LesserPrim](src ...A) bool { return IsSorted(src, LessEqPrim2[A]) }

/*
Returns the lesser of the two inputs. For primitive types that don't implement
`Lesser`, see `MinPrim2`. For a variadic variant, see `Min`.
*/
func Min2[A Lesser[A]](one, two A) A {
	if one.Less(two) {
		return one
	}
	return two
}

/*
Returns the lesser of the two inputs, which must be comparable primitives. For
non-primitives, see `Min2`. For a variadic variant, see `MinPrim`.
*/
func MinPrim2[A LesserPrim](one, two A) A {
	if one < two {
		return one
	}
	return two
}

/*
Returns the larger of the two inputs. For primitive types that don't implement
`Lesser`, see `MaxPrim2`. For a variadic variant, see `Max`.
*/
func Max2[A Lesser[A]](one, two A) A {
	if one.Less(two) {
		return two
	}
	return one
}

/*
Returns the larger of the two inputs, which must be comparable primitives. For
non-primitives, see `Max2`. For a variadic variant, see `MaxPrim`.
*/
func MaxPrim2[A LesserPrim](one, two A) A {
	if one < two {
		return two
	}
	return one
}
