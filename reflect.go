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
func TypeOf[A any](A) r.Type { return Type[A]() }

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
func Kind[A any]() r.Kind { return Type[A]().Kind() }

/*
Returns `reflect.Kind` of the given type. Never returns `reflect.Invalid`. If
the type parameter is an interface, the output is `reflect.Interface`.
*/
func KindOf[A any](A) r.Kind { return Type[A]().Kind() }

// Returns `reflect.Type.Size` of the given type.
func Size[A any]() uintptr { return Type[A]().Size() }

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
	return out, ErrConv(val.Interface(), Type[string]())
}

var StructFieldCache = TypeCacheOf[StructFields]()

type StructFields []r.StructField

func (self *StructFields) Init(src r.Type) {
	TimesAppend(self, src.NumField(), src.Field)
}

var StructPublicFieldCache = TypeCacheOf[StructPublicFields]()

type StructPublicFields []r.StructField

func (self *StructPublicFields) Init(src r.Type) {
	FilterAppend(self, StructFieldCache.Get(src), IsFieldPublic)
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
Self-explanatory. For some reason this is not provided in usable form by
the "reflect" package.
*/
func IsFieldPublic(val r.StructField) bool { return val.PkgPath == `` }

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

/*
True if the given type may contain any indirections (pointers). For any "direct"
type, assigning a value to another variable via `A := B` makes a complete copy.
For any "indirect" type, reassignment is insufficient to make a copy.

Special exceptions:

	* Strings are considered to be non-indirect, despite containing a pointer.
	  Generally in Go, strings are considered to be immutable and reassignment is
	  considered to be a copy.

	* Chans are considered to be non-indirect.

	* Funcs are considered to be non-indirect.

	* For structs, only public fields are checked.
*/
func IsIndirect(typ r.Type) bool {
	switch TypeKind(typ) {
	case r.Array:
		return typ.Len() > 0 && IsIndirect(typ.Elem())
	case r.Slice:
		return true
	case r.Interface:
		return true
	case r.Map:
		return true
	case r.Pointer:
		return true
	case r.Struct:
		return Some(StructPublicFieldCache.Get(typ), IsFieldIndirect)
	default:
		return false
	}
}

// Shortcut for testing if the field's type is `IsIndirect`.
func IsFieldIndirect(val r.StructField) bool { return IsIndirect(val.Type) }

/*
Returns a deep clone of the given value. Doesn't clone chans and funcs,
preserving them as-is. If the given value is "direct" (see `IsIndirect`), this
function doesn't allocate and simply returns the input as-is.
*/
func CloneDeep[A any](src A) A {
	ValueClone(r.ValueOf(AnyNoEscUnsafe(&src)).Elem())
	return src
}

/*
Mutates the input in-place, replacing the underlying data with a deep clone if
necessary. The given value must be settable.
*/
func ValueClone(src r.Value) {
	switch src.Kind() {
	case r.Array:
		cloneArray(src)
	case r.Slice:
		cloneSlice(src)
	case r.Interface:
		cloneInterface(src)
	case r.Map:
		cloneMap(src)
	case r.Pointer:
		clonePointer(src)
	case r.Struct:
		cloneStruct(src)
	}
}

// Similar to `CloneDeep` but takes and returns `reflect.Value`.
func ValueCloned(src r.Value) r.Value {
	switch src.Kind() {
	case r.Array:
		return clonedArray(src)
	case r.Slice:
		return clonedSlice(src)
	case r.Interface:
		return clonedInterface(src)
	case r.Map:
		return clonedMap(src)
	case r.Pointer:
		return clonedPointer(src)
	case r.Struct:
		return clonedStruct(src)
	default:
		return src
	}
}

// Idempotent set. Calls `reflect.Value.Set` only if the inputs are distinct.
func ValueSet(tar, src r.Value) {
	if tar != src {
		tar.Set(src)
	}
}

// Shortcut for `reflect.New(typ).Elem()`.
func NewElem(typ r.Type) r.Value { return r.New(typ).Elem() }
