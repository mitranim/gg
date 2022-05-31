package gg

import u "unsafe"

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
