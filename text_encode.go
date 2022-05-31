package gg

import (
	"encoding"
	"fmt"
	r "reflect"
	"strconv"
)

/*
Shortcut for implementing string encoding of `Nullable` types.
Mostly for internal use.
*/
func StringNull[A any, B NullableValGetter[A]](val B) string {
	if val.IsNull() {
		return ``
	}
	return String(val.Get())
}

/*
Alias for `fmt.Sprint` defined as a generic function for compatibility with
higher-order functions like `Map`. Slightly more efficient than `fmt.Sprint`:
avoids spurious heap escape and copying.

The output of this function is intended only for debug purposes. For machine
consumption or user display, use `String`, which is more restrictive.
*/
func StringAny[A any](val A) string { return fmt.Sprint(AnyNoEscUnsafe(val)) }

// Stringifies an arbitrary value via `StringCatch`. Panics on errors.
func String[A any](val A) string { return Try1(StringCatch(val)) }

/*
Missing feature of the standard library. Converts an arbitrary value to a
string, allowing only INTENTIONALLY stringable values. Rules:

	* Nil is considered "".

	* A string is returned as-is.

	* A byte slice is cast to a string.

	* Any other primitive value (see constraint `Prim`) is encoded via `strconv`.

	* Types that support `fmt.Stringer`, `Appender` or `encoding.TextMarshaler`
	  are encoded by using the corresponding method.

	* Any other type causes an error.
*/
func StringCatch[A any](val A) (string, error) {
	box := AnyNoEscUnsafe(val)

	stringer, _ := box.(fmt.Stringer)
	if stringer != nil {
		return stringer.String(), nil
	}

	appender, _ := box.(Appender)
	if appender != nil {
		return ToString(appender.Append(nil)), nil
	}

	marshaler, _ := box.(encoding.TextMarshaler)
	if marshaler != nil {
		val, err := marshaler.MarshalText()
		return ToString(val), err
	}

	switch val := box.(type) {
	case nil:
		return ``, nil

	case string:
		return val, nil

	case []byte:
		return ToString(val), nil

	case bool:
		return strconv.FormatBool(val), nil

	case int8:
		return strconv.FormatInt(int64(val), 10), nil
	case int16:
		return strconv.FormatInt(int64(val), 10), nil
	case int32:
		return strconv.FormatInt(int64(val), 10), nil
	case int64:
		return strconv.FormatInt(int64(val), 10), nil
	case int:
		return strconv.FormatInt(int64(val), 10), nil

	case uint8:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint:
		return strconv.FormatUint(uint64(val), 10), nil

	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(float64(val), 'f', -1, 64), nil

	default:
		return StringReflectCatch(r.ValueOf(box))
	}
}

/*
Reflection-based component of `StringCatch`.
Mostly for internal use.
*/
func StringReflectCatch(val r.Value) (string, error) {
	if !val.IsValid() {
		return ``, nil
	}

	typ := val.Type()

	switch typ.Kind() {
	case r.String:
		return val.String(), nil

	case r.Bool:
		return strconv.FormatBool(val.Bool()), nil

	case r.Int8, r.Int16, r.Int32, r.Int64, r.Int:
		return strconv.FormatInt(val.Int(), 10), nil

	case r.Uint8, r.Uint16, r.Uint32, r.Uint64, r.Uint:
		return strconv.FormatUint(val.Uint(), 10), nil

	case r.Float32, r.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil

	default:
		return ValueToStringCatch(val)
	}
}

/*
Shortcut for implementing string encoding of `Nullable` types.
Mostly for internal use.
*/
func AppendNull[A any, B NullableValGetter[A]](buf []byte, src B) []byte {
	if src.IsNull() {
		return buf
	}
	return Append(buf, src.Get())
}

/*
Appends text representation of the input to the given buffer,
using `AppendCatch`. Panics on errors.
*/
func Append[A ~[]byte, B any](buf A, src B) A {
	return Try1(AppendCatch(buf, src))
}

/*
Same as `StringCatch`, but instead of returning a string, appends the text
representation of the input to the given buffer. See `StringCatch` for the
encoding rules.
*/
func AppendCatch[A ~[]byte, B any](buf A, src B) (A, error) {
	box := AnyNoEscUnsafe(src)

	switch val := box.(type) {
	case nil:
		return buf, nil

	case string:
		return append(buf, val...), nil

	case []byte:
		return append(buf, val...), nil

	case bool:
		return strconv.AppendBool(buf, val), nil

	case int8:
		return strconv.AppendInt(buf, int64(val), 10), nil
	case int16:
		return strconv.AppendInt(buf, int64(val), 10), nil
	case int32:
		return strconv.AppendInt(buf, int64(val), 10), nil
	case int64:
		return strconv.AppendInt(buf, int64(val), 10), nil
	case int:
		return strconv.AppendInt(buf, int64(val), 10), nil

	case uint8:
		return strconv.AppendUint(buf, uint64(val), 10), nil
	case uint16:
		return strconv.AppendUint(buf, uint64(val), 10), nil
	case uint32:
		return strconv.AppendUint(buf, uint64(val), 10), nil
	case uint64:
		return strconv.AppendUint(buf, uint64(val), 10), nil
	case uint:
		return strconv.AppendUint(buf, uint64(val), 10), nil

	case float32:
		return strconv.AppendFloat(buf, float64(val), 'f', -1, 64), nil
	case float64:
		return strconv.AppendFloat(buf, float64(val), 'f', -1, 64), nil

	default:
		appender, _ := val.(Appender)
		if appender != nil {
			return appender.Append(buf), nil
		}

		marshaler, _ := val.(encoding.TextMarshaler)
		if marshaler != nil {
			val, err := marshaler.MarshalText()
			return append(buf, val...), err
		}

		stringer, _ := val.(fmt.Stringer)
		if stringer != nil {
			return append(buf, stringer.String()...), nil
		}

		return AppendReflectCatch(buf, r.ValueOf(box))
	}
}

/*
Reflection-based component of `AppendCatch`.
Mostly for internal use.
*/
func AppendReflectCatch[A ~[]byte](buf A, val r.Value) (A, error) {
	if !val.IsValid() {
		return buf, nil
	}

	typ := val.Type()

	switch typ.Kind() {
	case r.String:
		return append(buf, val.String()...), nil

	case r.Bool:
		return strconv.AppendBool(buf, val.Bool()), nil

	case r.Int8, r.Int16, r.Int32, r.Int64, r.Int:
		return strconv.AppendInt(buf, val.Int(), 10), nil

	case r.Uint8, r.Uint16, r.Uint32, r.Uint64, r.Uint:
		return strconv.AppendUint(buf, val.Uint(), 10), nil

	case r.Float32, r.Float64:
		return strconv.AppendFloat(buf, val.Float(), 'f', -1, 64), nil

	default:
		str, err := ValueToStringCatch(val)
		return append(buf, str...), err
	}
}

/*
Shortcut for implementing `encoding.TextMarshaler` on arbitrary types. Mostly
for internal use. Uses `StringCatch` internally. The resulting bytes may be
backed by constant storage and must not be mutated.
*/
func Marshal(val any) ([]byte, error) {
	str, err := StringCatch(val)
	return ToBytes(str), err
}

/*
Shortcut for implementing `encoding.TextMarshaler` on `Nullable` types. Mostly
for internal use. Uses `StringCatch` internally. The resulting bytes may be
backed by constant storage and must not be mutated.
*/
func MarshalNullCatch[A any, B NullableValGetter[A]](val B) ([]byte, error) {
	if val.IsNull() {
		return nil, nil
	}
	return Marshal(val.Get())
}

/*
Shortcut for stringifying a type that implements `Appender`.
Mostly for internal use.
*/
func AppenderString[A Appender](val A) string {
	return ToString(val.Append(nil))
}

/*
Appends the `fmt.GoStringer` representation of the given input to the given
buffer. Also see the function `GoString`.
*/
func AppendGoString[A any](inout []byte, val A) []byte {
	buf := Buf(inout)
	box := AnyNoEscUnsafe(val)

	impl, _ := box.(fmt.GoStringer)
	if impl != nil {
		buf.AppendString(impl.GoString())
		return buf
	}

	fmt.Fprintf(NoEscUnsafe(&buf), `%#v`, box)
	return buf
}

/*
Returns the `fmt.GoStringer` representation of the given input.
Equivalent to `fmt.Sprintf("%#v", val)` but marginally more efficient.
*/
func GoString[A any](val A) string {
	box := AnyNoEscUnsafe(val)

	impl, _ := box.(fmt.GoStringer)
	if impl != nil {
		return impl.GoString()
	}

	return fmt.Sprintf(`%#v`, box)
}
