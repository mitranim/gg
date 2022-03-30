package gg

import (
	"path"
	r "reflect"
	"runtime"
	rt "runtime"
	"strings"
)

/*
Returns `reflect.Type` of the given type. Differences from `reflect.TypeOf`:

	* Avoids spurious heap escape and copying.

	* Output is always non-nil.

	* When the given type is an interface, including the empty interface `any`,
	  the output is a non-nil `reflect.Type` describing the given interface.
*/
func Type[A any]() r.Type { return r.TypeOf((*A)(nil)).Elem() }

/*
Similar to `reflect.TypeOf`, with the following differences:

	* Avoids spurious heap escape and copying.

	* Output is always non-nil.

	* When the given type is an interface, including the empty interface `any`,
	  the output is a non-nil `reflect.Type` describing the given interface.
*/
func TypeOf[A any](A) r.Type { return r.TypeOf((*A)(nil)).Elem() }

/*
Nil-safe version of `reflect.Type.Kind`. If the input is nil, returns
`reflect.Invalid`.
*/
func TypeKind(val r.Type) r.Kind {
	if val == nil {
		return r.Invalid
	}
	return val.Kind()
}

/*
Returns `reflect.Kind` of the given `any`. Compare our generic functions `Kind`
and `KindOf` which take a concrete type.
*/
func KindOfAny(val any) r.Kind {
	return TypeKind(r.TypeOf(AnyNoEscUnsafe(val)))
}

/*
Returns `reflect.Kind` of the given type. Never returns `reflect.Invalid`. If
the type parameter is an interface, the output is `reflect.Interface`.
*/
func Kind[A any]() r.Kind {
	return r.TypeOf((*A)(nil)).Elem().Kind()
}

/*
Returns `reflect.Kind` of the given type. Never returns `reflect.Invalid`. If
the type parameter is an interface, the output is `reflect.Interface`.
*/
func KindOf[A any](A) r.Kind {
	return r.TypeOf((*A)(nil)).Elem().Kind()
}

// Uses `reflect.Zero` to create a zero value of the given type.
func ZeroValue[A any]() r.Value { return r.Zero(Type[A]()) }

// Takes an arbitrary function and returns its name.
func FuncName(val any) string { return FuncNameBase(RuntimeFunc(val)) }

// Takes an arbitrary function and returns its `runtime.Func`.
func RuntimeFunc(val any) *rt.Func {
	return runtime.FuncForPC(r.ValueOf(val).Pointer())
}

// Returns the given function's name without the package path prefix.
func FuncNameBase(fun *rt.Func) string {
	if fun == nil {
		return ``
	}
	return path.Base(fun.Name())
}

/*
Returns the name of the given function stripped of various namespaces: package
path prefix, package name, type name.
*/
func FuncNameShort(name string) string {
	// TODO cleanup.

	name = path.Base(name)

	for len(name) > 0 {
		ind := strings.IndexByte(name, '.')
		if ind >= 0 &&
			len(name) > (ind+1) &&
			!(name[ind+1] == '.' || name[ind+1] == ']') &&
			!isFuncNameAnon(name[:ind]) {
			name = name[ind+1:]
			continue
		}
		break
	}

	return name
}

// True if the value's underlying type is convertible to `[]byte`.
func IsValueBytes(val r.Value) bool {
	return val.IsValid() && IsTypeBytes(val.Type())
}

// True if the type is convertible to `[]byte`.
func IsTypeBytes(typ r.Type) bool {
	return (typ != nil) &&
		(typ.Kind() == r.Slice || typ.Kind() == r.Array) &&
		(typ.Elem().Kind() == r.Uint8)
}

/*
If the underlying type is compatible with `Text`, unwraps and converts it to a
string. Otherwise returns zero value. Boolean indicates success.
*/
func AnyToString(src any) (string, bool) {
	switch src := AnyNoEscUnsafe(src).(type) {
	case string:
		return src, true
	case []byte:
		return ToString(src), true
	}

	return ValueToString(r.ValueOf(AnyNoEscUnsafe(src)))
}

// Reflection-based component of `AnyToString`. For internal use.
func ValueToString(val r.Value) (string, bool) {
	if !val.IsValid() {
		return ``, true
	}

	if val.Kind() == r.String {
		return val.String(), true
	}

	if IsValueBytes(val) {
		return ToString(val.Bytes()), true
	}

	return ``, false
}

/*
If the underlying type is compatible with `Text`, unwraps and converts it to the
given text type. Otherwise returns zero value. Boolean indicates success. If the
given value is backed by `string` but the output type is backed by `[]byte`,
or vice versa, this performs a copy. Otherwise this doesn't allocate.
*/
func AnyToText[A Text](src any) (A, bool) {
	return ValueToText[A](r.ValueOf(AnyNoEscUnsafe(src)))
}

// Reflection-based component of `AnyToText`. For internal use.
func ValueToText[A Text](val r.Value) (A, bool) {
	if !val.IsValid() {
		return Zero[A](), true
	}

	if val.Kind() == r.String {
		return A(val.String()), true
	}

	if IsValueBytes(val) {
		return A(val.Bytes()), true
	}

	return Zero[A](), false
}

/*
Same as `ValueToString` but instead of boolean true/false, returns a nil/non-nil
error. The error describes the failure to convert the input to a string.
*/
func ValueToStringCatch(val r.Value) (string, error) {
	out, ok := ValueToString(val)
	if ok {
		return out, nil
	}
	return out, ErrConv(val, Type[string]())
}

var StructFieldCache = TypeCacheOf[StructFields]()

type StructFields []r.StructField

func (self *StructFields) Init(src r.Type) {
	*self = Times(src.NumField(), src.Field)
}

/*
Takes a struct field tag and returns its identifier part, following the
"encoding/json" conventions. Ident "-" is converted to "". Usage:

	ident := TagIdent(someField.Tag.Get(`json`))
	ident := TagIdent(someField.Tag.Get(`db`))

Rules:

	json:"ident"         -> "ident"
	json:"ident,<extra>" -> "ident"
	json:"-"             -> ""
	json:"-,<extra>"     -> ""
*/
func TagIdent(val string) string {
	ind := strings.IndexRune(string(val), ',')
	if ind >= 0 {
		val = val[:ind]
	}
	if val == `-` {
		return ``
	}
	return val
}

/*
Returns the field's DB/SQL column name from the "db" tag, following the same
conventions as the `encoding/json` package.
*/
func FieldDbName(val r.StructField) string {
	return TagIdent(val.Tag.Get(`db`))
}

/*
Returns the field's JSON column name from the "json" tag, following the same
conventions as the `encoding/json` package.
*/
func FieldJsonName(val r.StructField) string {
	return TagIdent(val.Tag.Get(`json`))
}

/*
Returns the element type of the provided type, automatically dereferencing
pointer types. If the input is nil, returns nil.
*/
func TypeDeref(val r.Type) r.Type {
	for val != nil && val.Kind() == r.Pointer {
		val = val.Elem()
	}
	return val
}

/*
Dereferences the provided value until it's no longer a pointer. If the input is
a nil pointer or a pointer to a nil pointer (recursively), returns an
empty/invalid value.
*/
func ValueDeref(val r.Value) r.Value {
	for val.Kind() == r.Pointer {
		if val.IsNil() {
			return r.Value{}
		}
		val = val.Elem()
	}
	return val
}
