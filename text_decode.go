package gg

import (
	"encoding"
	r "reflect"
	"strconv"
)

// Decodes arbitrary text into value of specific type, using `Parse`.
func ParseTo[Out any, Src Text](src Src) (out Out) {
	Parse(src, &out)
	return
}

/*
Decodes arbitrary text into value of specific type, using `ParseCatch`.
Panics on errors.
*/
func Parse[Out any, Src Text](src Src, out *Out) {
	Try(ParseCatch(src, out))
}

/*
Missing feature of the standard library. Decodes arbitrary text into a value of
an arbitrary type. The output must either implement `Parser`, or implement
`encoding.TextUnmarshaler`, or be a pointer to any of the types described by
the constraint `Textable` defined by this package. If the output is not
decodable, this returns an error.
*/
func ParseCatch[Out any, Src Text](src Src, out *Out) error {
	if out == nil {
		return nil
	}

	box := AnyNoEscUnsafe(out)

	parser, _ := box.(Parser)
	if parser != nil {
		return parser.Parse(string(src))
	}

	unmarshaler, _ := box.(encoding.TextUnmarshaler)
	if unmarshaler != nil {
		return unmarshaler.UnmarshalText(ToBytes(src))
	}

	return ParseValueCatch(src, r.ValueOf(box).Elem())
}

/*
Reflection-based component of `ParseCatch`.
Mostly for internal use.
*/
func ParseValueCatch[A Text](src A, out r.Value) error {
	out = NoEscUnsafe(out)
	typ := out.Type()
	kind := typ.Kind()

	switch kind {
	case r.Int8, r.Int16, r.Int32, r.Int64, r.Int:
		val, err := strconv.ParseInt(ToString(src), 10, typeBitSize(typ))
		out.SetInt(val)
		return ErrParse(err, src, typ)

	case r.Uint8, r.Uint16, r.Uint32, r.Uint64, r.Uint:
		val, err := strconv.ParseUint(ToString(src), 10, typeBitSize(typ))
		out.SetUint(val)
		return ErrParse(err, src, typ)

	case r.Float32, r.Float64:
		val, err := strconv.ParseFloat(ToString(src), typeBitSize(typ))
		out.SetFloat(val)
		return ErrParse(err, src, typ)

	case r.Bool:
		return parseBool(ToString(src), out)

	case r.String:
		out.SetString(string(src))
		return nil

	default:
		if IsTypeBytes(typ) {
			out.SetBytes([]byte(src))
			return nil
		}
		return Errf(`unable to convert string into %v: unsupported kind %v`, typ, kind)
	}
}

// Note: `strconv.ParseBool` is too permissive for our taste.
func parseBool(src string, out r.Value) error {
	switch src {
	case `true`:
		out.SetBool(true)
		return nil

	case `false`:
		out.SetBool(false)
		return nil

	default:
		return ErrParse(ErrInvalidInput, src, Type[bool]())
	}
}

/*
Shortcut for implementing text decoding of types that wrap other types, such as
`Opt`. Mostly for internal use.
*/
func ParseClearCatch[Out any, Tar ClearerPtrGetter[Out], Src Text](src Src, tar Tar) error {
	if len(src) == 0 {
		tar.Clear()
		return nil
	}
	return ParseCatch(src, tar.GetPtr())
}

/*
Shortcut for implementing `sql.Scanner` on types that wrap other types, such as
`Opt`. Mostly for internal use.
*/
func ScanCatch[Inner any, Outer PtrGetter[Inner]](src any, tar Outer) error {
	if src == nil {
		return nil
	}

	ptr := tar.GetPtr()

	impl, _ := AnyNoEscUnsafe(ptr).(Scanner)
	if impl != nil {
		return impl.Scan(src)
	}

	str, ok := AnyToText[string](src)
	if ok {
		return ParseCatch(str, ptr)
	}

	return ErrConv(src, Type[Outer]())
}
