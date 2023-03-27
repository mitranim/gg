package grepr

import (
	"fmt"
	r "reflect"
	"strconv"
	"unicode/utf8"

	"github.com/mitranim/gg"
)

func fmtAny(buf *Fmt, src r.Value) {
	if fmtedNil(buf, src) || fmtedGoString(buf, src) {
		return
	}

	if src.Type() == gg.Type[r.Type]() {
		fmtType(buf, src)
		return
	}

	switch src.Kind() {
	case r.Invalid:
		fmtNil(buf)
	case r.Bool:
		buf.AppendBool(src.Bool())
	case r.Int8:
		buf.AppendInt64(src.Int())
	case r.Int16:
		buf.AppendInt64(src.Int())
	case r.Int32:
		buf.AppendInt64(src.Int())
	case r.Int64:
		buf.AppendInt64(src.Int())
	case r.Int:
		buf.AppendInt64(src.Int())
	case r.Uint8:
		fmtUintHex(buf, src.Uint())
	case r.Uint16:
		buf.AppendUint64(src.Uint())
	case r.Uint32:
		buf.AppendUint64(src.Uint())
	case r.Uint64:
		buf.AppendUint64(src.Uint())
	case r.Uint:
		buf.AppendUint64(src.Uint())
	case r.Uintptr:
		fmtUintHex(buf, uintptr(src.Uint()))
	case r.Float32:
		buf.AppendFloat64(src.Float())
	case r.Float64:
		buf.AppendFloat64(src.Float())
	case r.Complex64:
		fmtComplex(buf, src.Complex())
	case r.Complex128:
		fmtComplex(buf, src.Complex())
	case r.Array:
		fmtArray(buf, src)
	case r.Slice:
		fmtSlice(buf, src)
	case r.Chan:
		fmtChan(buf, src)
	case r.Func:
		fmtFunc(buf, src)
	case r.Interface:
		fmtIface(buf, src)
	case r.Map:
		fmtMap(buf, src)
	case r.Pointer:
		fmtPointer(buf, src)
	case r.UnsafePointer:
		fmtUintHex(buf, src.Pointer())
	case r.String:
		fmtString(buf, src.String())
	case r.Struct:
		fmtStruct(buf, src)
	default:
		panic(gg.Errf(`unrecognized reflect kind %q`, src.Kind()))
	}
}

func fmtedVisited(buf *Fmt, src r.Value) bool {
	if src.IsNil() {
		fmtNil(buf)
		return true
	}

	ptr := src.UnsafePointer()
	_, ok := buf.Visited[ptr]
	if ok {
		fmtVisited(buf, src)
		return true
	}

	buf.Visited.Init().Add(ptr)
	return false
}

func fmtVisited(buf *Fmt, src r.Value) {
	buf.AppendString(`/* visited */ (`)
	fmtTypeName(buf, src.Type())
	buf.AppendString(`)(`)
	fmtUintHex(buf, src.Pointer())
	buf.AppendString(`)`)
}

func fmtedNil(buf *Fmt, src r.Value) bool {
	if !src.IsValid() || isValueNil(src) {
		fmtNil(buf)
		return true
	}
	return false
}

func fmtNil(buf *Fmt) { buf.AppendString(`nil`) }

/*
TODO: consider custom interface such as `.AppendGoString`, possibly with
indentation support.

TODO: if, rather than implementing `.GoString` directly, the input inherits the
method from an embedded type, we should do nothing and return false.
*/
func fmtedGoString(buf *Fmt, src r.Value) bool {
	if src.CanConvert(gg.Type[fmt.GoStringer]()) {
		buf.AppendString(src.Interface().(fmt.GoStringer).GoString())
		return true
	}
	return false
}

func fmtComplex(buf *Fmt, src complex128) {
	buf.AppendByte('(')
	buf.AppendFloat64(real(src))
	img := imag(src)
	if !(img < 0) {
		buf.AppendByte('+')
	}
	buf.AppendFloat64(img)
	buf.AppendByte('i')
	buf.AppendByte(')')
}

func fmtUintHex[A uint64 | uintptr](buf *Fmt, src A) {
	buf.AppendString(`0x`)
	buf.Buf = strconv.AppendUint(buf.Buf, uint64(src), 16)
}

/*
TODO: if the type is not exactly `string` and the value is not used in a
strongly typed context, wrap the literal in a cast.

TODO: same for other literals: bools, ints, floats, bytes, runes, complex.
*/
func fmtString(buf *Fmt, src string) {
	if buf.IsMulti() && CanBackquote(src) {
		buf.AppendByte('`')
		buf.AppendString(src)
		buf.AppendByte('`')
	} else {
		buf.Buf = strconv.AppendQuote(buf.Buf, src)
	}
}

func fmtSlice(buf *Fmt, src r.Value) {
	if fmtedNil(buf, src) {
		return
	}

	if gg.IsValueBytes(src) {
		// TODO: elide type name when elidable AND when not printing as a string.
		fmtTypeName(buf, src.Type())
		fmtBytesInner(buf, src.Bytes())
		return
	}

	fmtArray(buf, src)
}

func fmtArray(buf *Fmt, src r.Value) {
	prev := setElideType(buf, isNotInterface(src.Type().Elem()))
	defer prev.Done()

	if !prev.Val {
		fmtTypeName(buf, src.Type())
	}

	if src.Len() == 0 {
		buf.AppendString(`{}`)
		return
	}

	if buf.IsSingle() {
		fmtArraySingle(buf, src)
		return
	}

	fmtArrayMulti(buf, src)
}

func fmtArraySingle(buf *Fmt, src r.Value) {
	buf.AppendByte('{')
	for ind := range gg.Iter(src.Len()) {
		if ind > 0 {
			buf.AppendString(`, `)
		}
		fmtAny(buf, src.Index(ind))
	}
	buf.AppendByte('}')
}

func fmtArrayMulti(buf *Fmt, src r.Value) {
	buf.AppendByte('{')
	buf.AppendNewline()
	snap := incLvl(buf)

	for ind := range gg.Iter(src.Len()) {
		fmtIndent(buf)
		fmtAny(buf, src.Index(ind))
		buf.AppendByte(',')
		buf.AppendNewline()
	}

	snap.Done()
	fmtIndent(buf)
	buf.AppendByte('}')
}

/*
TODO:

	* Looks like text -> append like string.
	* Otherwise -> append like bytes.
*/
func fmtBytesInner(buf *Fmt, src []byte) {
	buf.AppendByte('(')
	fmtString(buf, gg.ToString(src))
	buf.AppendByte(')')
}

func fmtChan(buf *Fmt, src r.Value) {
	fmtUnfmtable(buf, src)
}

func fmtFunc(buf *Fmt, src r.Value) {
	fmtUnfmtable(buf, src)
}

func fmtIface(buf *Fmt, src r.Value) {
	if fmtedNil(buf, src) {
		return
	}
	fmtAny(buf, src.Elem())
}

func fmtMap(buf *Fmt, src r.Value) {
	if fmtedNil(buf, src) {
		return
	}

	prev := setElideType(buf, isNotInterface(src.Type().Elem()))
	defer prev.Done()

	if !prev.Val {
		fmtTypeName(buf, src.Type())
	}

	if src.Len() == 0 {
		buf.AppendString(`{}`)
		return
	}

	if buf.IsSingle() {
		fmtMapSingle(buf, src)
		return
	}

	fmtMapMulti(buf, src)
}

func fmtMapSingle(buf *Fmt, src r.Value) {
	buf.AppendByte('{')

	iter := src.MapRange()
	var found bool

	for iter.Next() {
		if found {
			buf.AppendString(`, `)
		}
		found = true

		fmtAny(buf, iter.Key())
		buf.AppendString(`: `)
		fmtAny(buf, iter.Value())
	}

	buf.AppendByte('}')
}

func fmtMapMulti(buf *Fmt, src r.Value) {
	buf.AppendByte('{')
	buf.AppendNewline()

	iter := src.MapRange()
	snap := incLvl(buf)

	for iter.Next() {
		fmtIndent(buf)
		fmtAny(buf, iter.Key())
		buf.AppendString(`: `)
		fmtAny(buf, iter.Value())
		buf.AppendByte(',')
		buf.AppendNewline()
	}

	snap.Done()
	fmtIndent(buf)
	buf.AppendByte('}')
}

func fmtPointer(buf *Fmt, src r.Value) {
	if fmtedVisited(buf, src) {
		return
	}

	defer setElideType(buf, false).Done()
	src = src.Elem()

	if canAmpersand(src.Kind()) {
		buf.AppendByte('&')
		fmtAny(buf, src)
		return
	}

	buf.AppendString(`gg.Ptr`)
	fmtTypeArg(buf, src.Type())
	buf.AppendByte('(')
	fmtAny(buf, src)
	buf.AppendByte(')')
}

func fmtStruct(buf *Fmt, src r.Value) {
	prev := setElideType(buf, false)
	defer prev.Done()

	if !prev.Val {
		fmtTypeName(buf, src.Type())
	}

	if src.NumField() == 0 {
		buf.AppendString(`{}`)
		return
	}

	if buf.IsSingle() {
		fmtStructSingle(buf, src)
		return
	}

	fmtStructMulti(buf, src)
}

func fmtStructField(buf *Fmt, src r.Value, field r.StructField) {
	buf.AppendString(field.Name)
	buf.AppendString(`: `)
	fmtAny(buf, src)
}

func fmtStructSingle(buf *Fmt, src r.Value) {
	if isStructUnit(src) {
		fmtStructSingleAnon(buf, src)
		return
	}
	fmtStructSingleNamed(buf, src)
}

func fmtStructSingleAnon(buf *Fmt, src r.Value) {
	src = src.Field(0)

	buf.AppendByte('{')

	if !skipField(buf, src) {
		fmtAny(buf, src)
	}

	buf.AppendByte('}')
}

func fmtStructSingleNamed(buf *Fmt, src r.Value) {
	buf.AppendByte('{')

	var found bool

	for _, field := range gg.StructPublicFieldCache.Get(src.Type()) {
		src := src.FieldByIndex(field.Index)
		if skipField(buf, src) {
			continue
		}

		if found {
			buf.AppendString(`, `)
		}
		found = true

		fmtStructField(buf, src, field)
	}

	buf.AppendByte('}')
}

func fmtStructMulti(buf *Fmt, src r.Value) {
	if isStructUnit(src) {
		fmtStructMultiAnon(buf, src)
		return
	}
	fmtStructMultiNamed(buf, src)
}

func fmtStructMultiAnon(buf *Fmt, src r.Value) {
	src = src.Field(0)

	buf.AppendByte('{')

	if !skipField(buf, src) {
		defer incLvl(buf).Done()
		fmtAny(buf, src)
	}

	buf.AppendByte('}')
}

func fmtStructMultiNamed(buf *Fmt, src r.Value) {
	fields := gg.StructPublicFieldCache.Get(src.Type())

	if buf.SkipZeroFields() {
		test := func(field r.StructField) bool {
			return !src.FieldByIndex(field.Index).IsZero()
		}

		count := gg.Count(fields, test)

		if count == 0 {
			buf.AppendString(`{}`)
			return
		}

		if count == 1 {
			field := gg.Find(fields, test)
			fmtStructMultiNamedUnit(buf, src.FieldByIndex(field.Index), field)
			return
		}
	}

	fmtStructMultiNamedLines(buf, src, fields)
}

func fmtStructMultiNamedUnit(buf *Fmt, src r.Value, field r.StructField) {
	buf.AppendByte('{')
	fmtStructField(buf, src, field)
	buf.AppendByte('}')
}

func fmtStructMultiNamedLines(buf *Fmt, src r.Value, fields []r.StructField) {
	buf.AppendByte('{')
	buf.AppendNewline()
	snap := incLvl(buf)

	for _, field := range fields {
		src := src.FieldByIndex(field.Index)
		if skipField(buf, src) {
			continue
		}

		fmtIndent(buf)
		fmtStructField(buf, src, field)
		buf.AppendByte(',')
		buf.AppendNewline()
	}

	snap.Done()
	fmtIndent(buf)
	buf.AppendByte('}')
}

func fmtUnfmtable(buf *Fmt, src r.Value) {
	if fmtedNil(buf, src) {
		return
	}

	fmtTypeName(buf, src.Type())
	buf.AppendByte('(')
	fmtUintHex(buf, src.Pointer())
	buf.AppendByte(')')
}

func fmtTypeName(buf *Fmt, typ r.Type) {
	if typ == gg.Type[[]byte]() {
		buf.AppendString(`[]byte`)
	} else {
		buf.AppendString(typeName(typ))
	}
}

func fmtTypeArg(buf *Fmt, typ r.Type) {
	if isTypeDefaultForLiteral(typ) {
		return
	}

	buf.AppendByte('[')
	buf.AppendString(typeName(typ))
	buf.AppendByte(']')
}

func fmtType(buf *Fmt, src r.Value) {
	buf.AppendString(`gg.Type[`)
	fmtTypeName(buf, src.Interface().(r.Type))
	buf.AppendString(`]()`)
}

func fmtIndent(buf *Fmt) {
	buf.AppendStringN(buf.Indent, buf.Lvl)
}

/*
TODO: consider eliding the name of the "current" package. Is that possible?
We can inspect the stack trace, but we might be unable to define "current".
*/
func typeName(typ r.Type) string { return typ.String() }

func canAmpersand(kind r.Kind) bool {
	return kind == r.Array || kind == r.Slice || kind == r.Struct
}

func isTypeDefaultForLiteral(typ r.Type) bool {
	return typ == gg.Type[bool]() ||
		typ == gg.Type[int]() ||
		typ == gg.Type[string]()
}

func isValueNil(src r.Value) bool {
	return isValueNilable(src) && src.IsNil()
}

func isValueNilable(src r.Value) bool {
	switch src.Kind() {
	case r.Invalid, r.Chan, r.Func, r.Interface, r.Map, r.Pointer, r.Slice:
		return true
	default:
		return false
	}
}

func setElideType(buf *Fmt, val bool) gg.Snapshot[bool] {
	return gg.PtrSwap(&buf.ElideType, val)
}

func incLvl(buf *Fmt) gg.Snapshot[int] {
	return gg.PtrSwap(&buf.Lvl, buf.Lvl+1)
}

func skipField(buf *Fmt, src r.Value) bool {
	return buf.SkipZeroFields() && src.IsZero()
}

func isNotBackquotable(char rune) bool {
	const bom = '\ufeff'

	return char == utf8.RuneError ||
		char == '`' ||
		char == bom ||
		(char < ' ' && !(char == '\t' || char == '\n' || char == '\r'))
}

// TODO take type instead of value.
func isStructUnit(val r.Value) bool { return val.NumField() == 1 }

func isNotInterface(typ r.Type) bool {
	return gg.TypeKind(typ) != r.Interface
}
