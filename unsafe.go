package gg

import u "unsafe"

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
func Sizeof[A any]() uintptr { return u.Sizeof(Zero[A]()) }

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
Reinterprets existing memory as a byte slice. The resulting byte slice is backed
by the given pointer. Mutations of the resulting slice are reflected in the
source memory. Length and capacity are equal to the size of the referenced
memory. If the pointer is nil, the output is nil.
*/
func AsBytes[A any](tar *A) []byte {
	if tar == nil {
		return nil
	}
	return u.Slice(CastUnsafe[*byte](tar), Sizeof[A]())
}
