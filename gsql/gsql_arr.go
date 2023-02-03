package gsql

import (
	"database/sql/driver"

	"github.com/mitranim/gg"
)

/*
Shortcut for casting into `Arr`. Workaround for the lack of type inference in
type literals and casts.
*/
func ToArr[A any](val []A) Arr[A] { return val }

// Shortcut for creating `Arr` from the arguments.
func ArrOf[A any](val ...A) Arr[A] { return val }

/*
Short for "array". A slice type that supports SQL array encoding and decoding,
using the `{}` format. Examples:

	Arr[int]{10, 20}                  <-> '{10,20}'
	Arr[Arr[int]]{{10, 20}, {30, 40}} <-> '{{10,20},{30,40}}'
*/
type Arr[A any] []A

var (
	_ = gg.Encoder(gg.Zero[Arr[any]]())
	_ = gg.Decoder(gg.Zero[*Arr[any]]())
)

// Implement `gg.Nullable`. True if the slice is nil.
func (self Arr[A]) IsNull() bool { return self == nil }

// Implement `fmt.Stringer`. Returns an SQL encoding of the array.
func (self Arr[A]) String() string { return gg.AppenderString(self) }

/*
Implement `Appender`, appending the array's SQL encoding to the buffer.
If the slice is nil, appends nothing.
*/
func (self Arr[A]) Append(buf []byte) []byte {
	if self != nil {
		buf = append(buf, '{')
		buf = self.AppendInner(buf)
		buf = append(buf, '}')
	}
	return buf
}

// Same as `.Append` but without the enclosing `{}`.
func (self Arr[A]) AppendInner(buf []byte) []byte {
	var found bool
	for _, val := range self {
		if found {
			buf = append(buf, ',')
		}
		found = true
		buf = gg.Append(buf, val)
	}
	return buf
}

// Decodes from an SQL array literal string. Supports nested arrays.
func (self *Arr[A]) Parse(src string) (err error) {
	defer gg.Rec(&err)
	defer gg.Detailf(`unable to decode %q into %T`, src, self)

	self.Clear()

	if len(src) == 0 {
		return nil
	}

	if src == `{}` {
		if *self == nil {
			*self = Arr[A]{}
		}
		return nil
	}

	if !(gg.StrHead(src) == '{' && gg.StrLast(src) == '}') {
		panic(gg.ErrInvalidInput)
	}
	src = src[1 : len(src)-1]

	for len(src) > 0 {
		gg.AppendVals(self, gg.ParseTo[A](popSqlArrSegment(&src)))
	}
	return nil
}

// Truncates the length, keeping the capacity.
func (self *Arr[A]) Clear() { gg.SliceTrunc(self) }

// Implement `driver.Valuer`.
func (self Arr[A]) Value() (driver.Value, error) {
	if self.IsNull() {
		return nil, nil
	}
	return self.String(), nil
}

// Implement `sql.Scanner`.
func (self *Arr[A]) Scan(src any) error {
	str, ok := gg.AnyToText[string](src)
	if ok {
		return self.Parse(str)
	}

	switch src := src.(type) {
	case Arr[A]:
		*self = src
		return nil

	default:
		return gg.ErrConv(src, gg.Type[Arr[A]]())
	}
}
