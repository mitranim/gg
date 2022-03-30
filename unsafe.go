package gg

import u "unsafe"

// Memory representation of an arbitrary Go slice.
type SliceHeader struct {
	Dat u.Pointer
	Len int
	Cap int
}

/*
Dangerous tool for performance fine-tuning. Converts the given pointer to
`unsafe.Pointer` and tricks the compiler, preventing escape analysis of the
resulting pointer from moving the underlying memory to the heap. Can negate
failures of Go escape analysis, but can also introduce tricky bugs. The caller
MUST ensure that the original is not freed while the resulting pointer is still
in use.
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
