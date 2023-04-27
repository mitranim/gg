package gg

/*
Takes an arbitrary value and returns a non-nil pointer to a new memory region
containing a shallow copy of that value.
*/
func Ptr[A any](val A) *A { return &val }

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
