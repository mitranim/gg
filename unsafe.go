package gg

import (
	r "reflect"
	u "unsafe"
)

/*
Amount of bytes in `uintptr`. At the time of writing, in Go 1.20, this is also
the amount of bytes in `int` and `uint`.
*/
const SizeofWord = u.Sizeof(uintptr(0))

/*
Amount of bytes in any value of type `~string`. Note that this is the size of
the string header, not of the underlying data.
*/
const SizeofString = u.Sizeof(``)

/*
Amount of bytes in any slice header, for example of type `[]byte`. Note that the
size of a slice header is constant and does not reflect the size of the
underlying data.
*/
const SizeofSlice = u.Sizeof([]byte(nil))

/*
Amount of bytes in our own `SliceHeader`. In the official Go implementation
(version 1.20 at the time of writing), this is equal to `SizeofSlice`.
In case of mismatch, using `SliceHeader` for anything is invalid.
*/
const SizeofSliceHeader = u.Sizeof(SliceHeader{})

/*
Returns `unsafe.Sizeof` for the given type. Equivalent to `reflect.Type.Size`
for the same type. Due to Go's limitations, the result is not a constant, thus
you should prefer direct use of `unsafe.Sizeof` which returns a constant.
*/
func Size[A any]() uintptr { return u.Sizeof(Zero[A]()) }

/*
Memory representation of an arbitrary Go slice. Same as `reflect.SliceHeader`
but with `unsafe.Pointer` instead of `uintptr`.
*/
type SliceHeader struct {
	Dat u.Pointer
	Len int
	Cap int
}

/*
Takes a regular slice header and converts it to its underlying representation
`SliceHeader`.
*/
func SliceHeaderOf[A any](src []A) SliceHeader {
	return CastUnsafe[SliceHeader](src)
}

/*
Dangerous tool for performance fine-tuning. Converts the given pointer to
`unsafe.Pointer` and tricks the compiler into thinking that the memory
underlying the pointer should not be moved to the heap. Can negate failures of
Go escape analysis, but can also introduce tricky bugs. The caller MUST ensure
that the original is not freed while the resulting pointer is still in use.
*/
func PtrNoEscUnsafe[A any](val *A) u.Pointer { return noescape(u.Pointer(val)) }

// Dangerous tool for performance fine-tuning.
func NoEscUnsafe[A any](val A) A { return *(*A)(PtrNoEscUnsafe(&val)) }

// Dangerous tool for performance fine-tuning.
func AnyNoEscUnsafe(src any) any { return NoEscUnsafe(src) }

/*
Self-explanatory. Slightly cleaner and less error prone than direct use of
unsafe pointers.
*/
func CastUnsafe[Out, Src any](val Src) Out { return *(*Out)(u.Pointer(&val)) }

/*
Same as `CastUnsafe` but with additional validation: `unsafe.Sizeof` must be the
same for both types, otherwise this panics.
*/
func Cast[Out, Src any](src Src) Out {
	out := CastUnsafe[Out](src)
	srcSize := u.Sizeof(src)
	outSize := u.Sizeof(out)
	if srcSize == outSize {
		return out
	}
	panic(errSizeMismatch(Type[Src](), srcSize, Type[Out](), outSize))
}

func errSizeMismatch(src r.Type, srcSize uintptr, out r.Type, outSize uintptr) Err {
	return Errf(
		`size mismatch: %v (size %v) vs %v (size %v)`,
		src, srcSize, out, outSize,
	)
}

/*
Similar to `CastUnsafe` between two slice types but with additional validation:
`unsafe.Sizeof` must be the same for both element types, otherwise this
panics.
*/
func CastSlice[Out, Src any](src []Src) []Out {
	srcSize := Size[Src]()
	outSize := Size[Out]()
	if srcSize == outSize {
		return CastUnsafe[[]Out](src)
	}
	panic(errSizeMismatch(Type[Src](), srcSize, Type[Out](), outSize))
}

/*
Reinterprets existing memory as a byte slice. The resulting byte slice is backed
by the given pointer. Mutations of the resulting slice are reflected in the
source memory. Length and capacity are equal to the size of the referenced
memory. If the pointer is nil, the output is nil.
*/
func AsBytes[A any](tar *A) []byte {
	if tar == nil {
		return nil
	}
	return u.Slice(CastUnsafe[*byte](tar), Size[A]())
}
