package gg

import (
	"database/sql/driver"
	"fmt"
	"time"
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

// Describes all built-in complex number types and their typedefs.
type Complex interface{ ~complex64 | ~complex128 }

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

/*
Describes text types: strings and byte slices. All types compatible with this
interface can be freely cast to `[]byte` via `ToBytes` and to `string` via
`ToString`, subject to safety gotchas described in those functions' comments.
*/
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
type Pker[A comparable] interface{ Pk() A }

/*
Implemented by various utility types where zero value is considered null in
encoding/decoding contexts such as JSON and SQL.
*/
type Nullable interface{ IsNull() bool }

/*
Used by some 3rd party libraries. Implemented by some of our types for
compatibility.
*/
type AnyGetter Getter[any]

// Implemented by utility types that wrap arbitrary types, such as `Opt`.
type Getter[A any] interface{ Get() A }

// Implemented by utility types that wrap arbitrary types, such as `Opt`.
type Setter[A any] interface{ Set(A) }

/*
Similar to `Getter`, but for types that may perform work in `.Get`, this must
avoid that work and be very cheap to call. Used by `Mem`.
*/
type Peeker[A any] interface{ Peek() A }

/*
Implemented by utility types that wrap arbitrary types, such as `Opt`. The
returned pointer must reference the memory of the wrapper, instead of referring
to new memory. Its mutation must affect the wrapper. If the wrapper is nil,
this should return nil.
*/
type Ptrer[A any] interface{ Ptr() *A }

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

type ParserPtr[A any] interface {
	*A
	Parser
}

// Copy of `sql.Scanner`. Copied here to avoid a huge import.
type Scanner interface{ Scan(any) error }

// Used by some utility functions.
type ClearerPtrGetter[A any] interface {
	Clearer
	Ptrer[A]
}

// Used by some utility functions.
type NullableValGetter[A any] interface {
	Nullable
	Getter[A]
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
type AppenderTo interface{ AppendTo([]byte) []byte }

/*
Appends a text representation of a stack trace to the given buffer, returning
the modified buffer. Implemented by `Err` and used internally in error
formatting.
*/
type StackAppenderTo interface{ AppendStackTo([]byte) []byte }

/*
Combination of interfaces related to text encoding implemented by some types in
this package.
*/
type Encoder interface {
	fmt.Stringer
	AppenderTo
	Nullable
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
arbitrary error types. This interface is also implemented by trace-enabled
errors in "github.com/pkg/errors".
*/
type StackTraced interface{ StackTrace() []uintptr }

/*
The method `.Init` must modify the receiver, initializing any components that
need initialization, for example using `make` to create inner maps or chans.
The receiver must be mutable, usually a pointer. See `IniterPtr` for a more
precise type constraint. Also see `Initer1` which is more commonly used in this
library.
*/
type Initer interface{ Init() }

// Pointer version of `Initer`.
type IniterPtr[A any] interface {
	*A
	Initer
}

/*
The method `.Init` must modify the receiver, initializing any components that
need initialization, for example using `make` to create inner maps or chans.
The receiver must be mutable, usually a pointer. See `Initer1Ptr` for a more
precise type constraint. Also see nullary `Initer`.
*/
type Initer1[A any] interface{ Init(A) }

// Pointer version of `Initer1`.
type Initer1Ptr[A, B any] interface {
	*A
	Initer1[B]
}

/*
The method `.Default` must modify the receiver, applying default values to its
components, usually to struct fields. The receiver must be mutable, usually a
pointer. See `DefaulterPtr` for a more precise type constraint.
*/
type Defaulter interface{ Default() }

// Pointer version of `Defaulter`.
type DefaulterPtr[A any] interface {
	*A
	Defaulter
}

/*
Interface for values which are convertible to `time.Duration` or can specify
a lifetime for other values. Used by `Mem`.
*/
type Durationer interface{ Duration() time.Duration }

/*
Implemented by various types such as `context.Context`, `sql.Rows`, and our own
`Errs`.
*/
type Errer interface{ Err() error }

// Used by various "iterator" types such as `sql.Rows`.
type Nexter interface{ Next() bool }

/*
Implemented by some standard library types such as `time.Time` and
`reflect.IsZero`. Our generic function `IsZero` automatically invokes this
method on inputs that implement it.
*/
type Zeroable interface{ IsZero() bool }

/*
Implemented by some types such as `time.Time`, and invoked automatically by our
function `Equal`.
*/
type Equaler[A any] interface{ Equal(A) bool }

/*
Allows to customize/override `ErrFind`, which prioritizes this interface over
the default behavior.
*/
type ErrFinder interface{ Find(func(error) bool) error }

/*
Must return the length of a collection, such as a slice, map, text, etc.
Implemented by various collection types in this package.
*/
type Lener interface{ Len() int }
