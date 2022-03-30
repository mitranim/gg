package gg

import (
	"database/sql/driver"
	"fmt"
)

// Short for "signed integer".
type Sint interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

// Short for "unsigned integer".
type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint
}

// Describes all built-in integer types and their typedefs.
type Int interface{ Sint | Uint }

// Describes all built-in float types and their typedefs.
type Float interface{ ~float32 | ~float64 }

/*
Describes all built-in numeric types and their typedefs, excluding complex
numbers.
*/
type Num interface{ Int | Float }

/*
Describes all built-in signed numeric types and their typedefs, excluding
complex numbers.
*/
type Signed interface{ Sint | Float }

// Describes all types that support the "+" operator.
type Plusable interface{ Num | ~string }

// Set of "primitive" types which may be constant.
type Prim interface{ ~bool | ~string | Num }

/*
Describes built-in or well-known types which don't implement text encoding and
decoding intrinsically, but whose text encoding and decoding is supported
across the Go library ecosystem extrinsically.
*/
type Textable interface{ Prim | ~[]byte }

// Describes text types: strings and byte slices.
type Text interface{ ~string | ~[]byte }

/*
Describes all primitive types that support the "<" operator. Counterpart to
`Lesser` which describes types that support comparison via the `.Less` method.
*/
type LesserPrim interface {
	Num | Float | ~uintptr | ~string
}

/*
Describes arbitrary types that support comparison via `.Less`, similar to "<".
Used by various sorting/ordering utilities.
*/
type Lesser[A any] interface{ Less(A) bool }

/*
Short for "primary keyed". See type `Coll` which acts as an ordered map where
each value is indexed on its primary key. Keys must be non-zero. A zero value
is considered an invalid key.
*/
type Pked[A comparable] interface{ Pk() A }

/*
Implemented by various utility types where zero value is considered null in
encoding/decoding contexts such as JSON and SQL.
*/
type Nullable interface{ IsNull() bool }

/*
Implemented by various utility types. Enables compatibility with 3rd party
libraries such as `pgx`.
*/
type Getter interface{ Get() any }

// Implemented by utility types that wrap arbitrary types, such as `Opt`.
type ValGetter[A any] interface{ GetVal() A }

// Implemented by utility types that wrap arbitrary types, such as `Opt`.
type ValSetter[A any] interface{ SetVal(A) }

/*
Implemented by utility types that wrap arbitrary types, such as `Opt`. The
returned pointer must reference the memory of the wrapper, instead of referring
to new memory. Its mutation must affect the wrapper.
*/
type PtrGetter[A any] interface{ GetPtr() *A }

/*
Must clear the receiver. In collection types backed by slices and maps, this
should reduce length to 0, but is allowed to keep capacity.
*/
type Clearer interface{ Clear() }

/*
Interface for types that support parsing from a string. Counterpart to
`encoding.TextUnmarshaler`. Implemented by some utility types.
*/
type Parser interface{ Parse(string) error }

// Copy of `sql.Scanner`. Copied here to avoid a huge import.
type Scanner interface{ Scan(any) error }

// Used by some utility functions.
type ClearerPtrGetter[A any] interface {
	Clearer
	PtrGetter[A]
}

// Used by some utility functions.
type NullableValGetter[A any] interface {
	Nullable
	ValGetter[A]
}

// Used by some utilities.
type Runner interface{ Run() }

/*
Appends a text representation to the given buffer, returning the modified
buffer. Counterpart to `fmt.Stringer`. All types that implement this interface
should also implement `fmt.Stringer`, and in most cases this should be
semantically equivalent to appending the output of `.String`. However, this
interface allows significantly more efficient text encoding.
*/
type Appender interface{ Append([]byte) []byte }

/*
Combination of interfaces related to text encoding implemented by some types in
this package.
*/
type Encoder interface {
	fmt.Stringer
	Appender
	Nullable
	Getter
	driver.Valuer
}

/*
Combination of interfaces related to text decoding implemented by some types in
this package.
*/
type Decoder interface {
	Clearer
	Parser
	Scanner
}

/*
Implemented by the `Err` type. Used by `ErrTrace` to retrieve stack traces from
arbitrary error types.
*/
type StackTraced interface{ StackTrace() []uintptr }

// Used by `Cache`.
type Initer1[A, B any] interface {
	*A
	Init(B)
}
