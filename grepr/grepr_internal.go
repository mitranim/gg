package grepr

import (
	"fmt"
	"math"
	r "reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mitranim/gg"
)

var (
	typeRtype      = gg.Type[r.Type]()
	typeBool       = gg.Type[bool]()
	typeInt        = gg.Type[int]()
	typeString     = gg.Type[string]()
	typFloat64     = gg.Type[float64]()
	typeComplex64  = gg.Type[complex64]()
	typeComplex128 = gg.Type[complex128]()
	typeGoStringer = gg.Type[fmt.GoStringer]()
)

func (self *Fmt) fmtAny(typ r.Type, src r.Value) {
	if self.fmtedNil(typ, src) || self.fmtedGoString(src) {
		return
	}

	if typ == typeRtype {
		self.fmtReflectType(src)
		return
	}

	switch src.Kind() {
	case r.Invalid:
		self.fmtNil(typ, src)
	case r.Bool:
		self.fmtBool(typ, src)
	case r.Int8:
		self.fmtInt64(typ, src)
	case r.Int16:
		self.fmtInt64(typ, src)
	case r.Int32:
		self.fmtInt64(typ, src)
	case r.Int64:
		self.fmtInt64(typ, src)
	case r.Int:
		self.fmtInt64(typ, src)
	case r.Uint8:
		self.fmtByteHex(typ, src)
	case r.Uint16:
		self.fmtUint64(typ, src)
	case r.Uint32:
		self.fmtUint64(typ, src)
	case r.Uint64:
		self.fmtUint64(typ, src)
	case r.Uint:
		self.fmtUint64(typ, src)
	case r.Uintptr:
		self.fmtUintHex(typ, src)
	case r.Float32:
		self.fmtFloat64(typ, src)
	case r.Float64:
		self.fmtFloat64(typ, src)
	case r.Complex64:
		self.fmtComplex(typ, src)
	case r.Complex128:
		self.fmtComplex(typ, src)
	case r.Array:
		self.fmtArray(src)
	case r.Slice:
		self.fmtSlice(typ, src)
	case r.Chan:
		self.fmtChan(typ, src)
	case r.Func:
		self.fmtFunc(typ, src)
	case r.Interface:
		self.fmtIface(typ, src)
	case r.Map:
		self.fmtMap(typ, src)
	case r.Pointer:
		self.fmtPointer(typ, src)
	case r.UnsafePointer:
		self.fmtUnsafePointer(typ, src)
	case r.String:
		self.fmtString(typ, src)
	case r.Struct:
		self.fmtStruct(src)
	default:
		panic(gg.Errf(`unrecognized reflect kind %q`, src.Kind()))
	}
}

func (self *Fmt) fmtedPointerVisited(typ r.Type, src r.Value) bool {
	if src.IsNil() {
		self.fmtNil(typ, src)
		return true
	}

	ptr := src.UnsafePointer()
	_, ok := self.Visited[ptr]
	if ok {
		self.fmtPointerVisited(src)
		return true
	}

	self.Visited.Init().Add(ptr)
	return false
}

func (self *Fmt) fmtPointerVisited(src r.Value) {
	self.AppendString(`/* visited */ `)
	self.fmtTypeName(src.Type())
	self.AppendByte('(')
	self.fmtUint64Hex(uint64(src.Pointer()))
	self.AppendByte(')')
}

func (self *Fmt) fmtedNil(typ r.Type, src r.Value) bool {
	if !src.IsValid() || gg.IsValueNil(src) {
		self.fmtNil(typ, src)
		return true
	}
	return false
}

func (self *Fmt) fmtNil(typ r.Type, src r.Value) {
	srcTyp := gg.ValueType(src)

	// TODO simplify.
	if !(isTypeInterface(typ) && isValueNilInterface(src) ||
		typ == srcTyp && gg.IsTypeNilable(typ)) {
		defer self.fmtConvOpen(derefIface(src).Type()).fmtConvClose()
	}

	self.AppendString(`nil`)
}

/*
TODO: consider custom interface such as `.AppendGoString`, possibly with
indentation support.

TODO: if, rather than implementing `.GoString` directly, the input inherits the
method from an embedded type, we should do nothing and return false.
*/
func (self *Fmt) fmtedGoString(src r.Value) bool {
	if !src.Type().Implements(typeGoStringer) {
		return false
	}
	self.AppendString(src.Interface().(fmt.GoStringer).GoString())
	return true
}

func (self *Fmt) fmtBool(typ r.Type, src r.Value) {
	defer self.fmtConvOpt(typ, src.Type(), typeBool).fmtConvClose()
	self.AppendBool(src.Bool())
}

/*
Adapted from `strconv.FormatComplex` with minor changes. We don't bother
printing the real part when it's zero, and we avoid the scientific notation
when formatting floats.
*/
func (self *Fmt) fmtComplex(typ r.Type, src r.Value) {
	done := self.fmtConvOpt(typ, src.Type(), typeComplex128)
	if done != nil {
		defer done.fmtConvClose()
	} else {
		self.AppendByte('(')
		defer self.AppendByte(')')
	}

	val := src.Complex()
	realPart := real(val)
	imagPart := imag(val)

	if realPart == 0 {
		self.AppendFloat64(imagPart)
		self.AppendByte('i')
		return
	}

	self.AppendFloat64(realPart)
	if !(imagPart < 0) {
		self.AppendByte('+')
	}

	self.AppendFloat64(imagPart)
	self.AppendByte('i')
}

func (self *Fmt) fmtInt64(typ r.Type, src r.Value) {
	defer self.fmtConvOpt(typ, src.Type(), typeInt).fmtConvClose()
	self.AppendInt64(src.Int())
}

func (self *Fmt) fmtUint64(typ r.Type, src r.Value) {
	defer self.fmtConvOpt(typ, src.Type()).fmtConvClose()
	self.AppendUint64(src.Uint())
}

func (self *Fmt) fmtByteHex(typ r.Type, src r.Value) {
	defer self.fmtConvOpt(typ, src.Type()).fmtConvClose()
	self.AppendString(`0x`)
	self.AppendByteHex(byte(src.Uint()))
}

func (self *Fmt) fmtUintHex(typ r.Type, src r.Value) {
	defer self.fmtConvOpt(typ, src.Type()).fmtConvClose()
	self.fmtUint64Hex(src.Uint())
}

func (self *Fmt) fmtFloat64(typ r.Type, src r.Value) {
	val := src.Float()
	srcTyp := src.Type()
	if isTypeInterface(typ) && (srcTyp != typFloat64 || !(math.Remainder(val, 1) > 0)) {
		defer self.fmtConvOpen(srcTyp).fmtConvClose()
	}

	if val < 0 {
		if math.IsInf(val, -1) {
			self.AppendString(`math.Inf(-1)`)
			return
		}
		self.AppendFloat64(val)
		return
	}

	if val >= 0 {
		if math.IsInf(val, 1) {
			self.AppendString(`math.Inf(0)`)
			return
		}
		self.AppendFloat64(val)
		return
	}

	self.AppendString(`math.NaN()`)
}

func (self *Fmt) fmtString(typ r.Type, src r.Value) {
	defer self.fmtConvOpt(typ, src.Type(), typeString).fmtConvClose()
	text := src.String()
	self.fmtStringInner(text, textPrintability(text))
}

func (self *Fmt) fmtStringInner(src string, prn printability) {
	/**
	For the most part we're more restrictive than `strconv.CanBackquote`, but
	unlike `strconv.CanBackquote` we allow '\n' in backquoted strings. We want
	to avoid loss of information when a multiline string is displayed in a
	terminal and copied to an editor. `strconv.CanBackquote` allows too many
	characters which may fail to display properly. This includes the tabulation
	character '\t', which is usually converted to spaces, and a variety of
	Unicode code points without corresponding graphical symbols. On the other
	hand, we want to support multiline strings in the common case, without edge
	case breakage. We assume that terminals, and other means of displaying the
	output of a program that may be using `grepr`, do not convert printed '\n'
	to '\r' or "\r\n", but may convert printed '\r' or "\r\n" to '\n', as Unix
	line endings tend to be preferred by any tooling used by Go developers. This
	allows us to support displaying strings as multiline in the common case,
	while avoiding information loss in the case of strings with '\r'.
	*/
	if !prn.errors && !prn.unprintables && !prn.backquotes && !prn.carriageReturns && (!prn.lineFeeds || self.IsMulti()) {
		self.AppendByte('`')
		self.AppendString(src)
		self.AppendByte('`')
	} else {
		self.Buf = strconv.AppendQuote(self.Buf, src)
	}
}

func (self *Fmt) fmtSlice(typ r.Type, src r.Value) {
	if self.fmtedNil(typ, src) {
		return
	}

	if gg.IsValueBytes(src) {
		self.fmtBytes(typ, src)
		return
	}

	self.fmtArray(src)
}

func (self *Fmt) fmtArray(src r.Value) {
	typ := src.Type()

	self.fmtTypeNameOpt(typ)
	defer self.setElideType(!isTypeInterface(typ.Elem())).Done()

	if src.Len() == 0 || src.IsZero() {
		self.AppendString(`{}`)
		return
	}

	if self.IsSingle() {
		self.fmtArraySingle(typ, src)
		return
	}

	self.fmtArrayMulti(typ, src)
}

func (self *Fmt) fmtArraySingle(typ r.Type, src r.Value) {
	typElem := typ.Elem()

	self.AppendByte('{')
	for ind := range gg.Iter(src.Len()) {
		if ind > 0 {
			self.AppendString(`, `)
		}
		self.fmtAny(typElem, src.Index(ind))
	}
	self.AppendByte('}')
}

func (self *Fmt) fmtArrayMulti(typ r.Type, src r.Value) {
	typElem := typ.Elem()

	self.AppendByte('{')
	self.AppendNewline()
	snap := self.lvlInc()

	for ind := range gg.Iter(src.Len()) {
		self.fmtIndent()
		self.fmtAny(typElem, src.Index(ind))
		self.AppendByte(',')
		self.AppendNewline()
	}

	snap.Done()
	self.fmtIndent()
	self.AppendByte('}')
}

func (self *Fmt) fmtBytes(typ r.Type, src r.Value) {
	if self.fmtedNil(typ, src) {
		return
	}

	text := src.Bytes()

	if len(text) > 0 {
		prn := textPrintability(text)

		if !prn.errors && !prn.unprintables {
			self.fmtTypeName(src.Type())
			self.AppendByte('(')
			self.fmtStringInner(gg.ToString(text), prn)
			self.AppendByte(')')
			return
		}
	}

	self.fmtBytesHex(src.Type(), text)
}

/*
Similar to `.fmtArray`, but much faster and always single-line. TODO consider
supporting column width in `Conf`, which would allow us to print bytes in
rows.
*/
func (self *Fmt) fmtBytesHex(typ r.Type, src []byte) {
	self.fmtTypeNameOpt(typ)
	self.AppendByte('{')
	for ind, val := range src {
		if ind > 0 {
			self.AppendString(`, `)
		}
		self.AppendString(`0x`)
		self.AppendByteHex(val)
	}
	self.AppendByte('}')
}

func (self *Fmt) fmtChan(typ r.Type, src r.Value) {
	self.fmtUnfmtable(typ, src)
}

func (self *Fmt) fmtFunc(typ r.Type, src r.Value) {
	self.fmtUnfmtable(typ, src)
}

func (self *Fmt) fmtIface(typ r.Type, src r.Value) {
	if self.fmtedNil(typ, src) {
		return
	}
	self.fmtAny(typ, src.Elem())
}

func (self *Fmt) fmtMap(typ r.Type, src r.Value) {
	if self.fmtedNil(typ, src) {
		return
	}

	srcTyp := src.Type()

	self.fmtTypeNameOpt(srcTyp)
	defer self.setElideType(!isTypeInterface(srcTyp.Elem())).Done()

	if src.Len() == 0 {
		self.AppendString(`{}`)
		return
	}

	if self.IsSingle() {
		self.fmtMapSingle(src)
		return
	}

	self.fmtMapMulti(src)
}

func (self *Fmt) fmtMapSingle(src r.Value) {
	typ := src.Type()
	typKey := typ.Key()
	typVal := typ.Elem()

	self.AppendByte('{')

	iter := src.MapRange()
	var found bool

	for iter.Next() {
		if found {
			self.AppendString(`, `)
		}
		found = true

		self.fmtAny(typKey, iter.Key())
		self.AppendString(`: `)
		self.fmtAny(typVal, iter.Value())
	}

	self.AppendByte('}')
}

func (self *Fmt) fmtMapMulti(src r.Value) {
	typ := src.Type()
	typKey := typ.Key()
	typVal := typ.Elem()

	self.AppendByte('{')
	self.AppendNewline()

	iter := src.MapRange()
	snap := self.lvlInc()

	for iter.Next() {
		self.fmtIndent()
		self.fmtAny(typKey, iter.Key())
		self.AppendString(`: `)
		self.fmtAny(typVal, iter.Value())
		self.AppendByte(',')
		self.AppendNewline()
	}

	snap.Done()
	self.fmtIndent()
	self.AppendByte('}')
}

func (self *Fmt) fmtPointer(typ r.Type, src r.Value) {
	if self.fmtedNil(typ, src) || self.fmtedPointerVisited(typ, src) {
		return
	}

	defer self.setElideType(false).Done()
	src = src.Elem()

	if canAmpersand(src.Kind()) {
		self.AppendByte('&')
		self.fmtAny(typ, src)
		return
	}

	self.fmtIdent(`gg`, `Ptr`)
	self.fmtTypeArg(src.Type())
	self.AppendByte('(')
	self.fmtAny(typ, src)
	self.AppendByte(')')
}

func (self *Fmt) fmtUnsafePointer(typ r.Type, src r.Value) {
	defer self.fmtConvOpt(typ, src.Type()).fmtConvClose()
	self.fmtUint64Hex(uint64(src.Pointer()))
}

func (self *Fmt) fmtStruct(src r.Value) {
	self.fmtTypeNameOpt(src.Type())
	defer self.setElideType(false).Done()

	if src.NumField() == 0 {
		self.AppendString(`{}`)
		return
	}

	if self.IsSingle() {
		self.fmtStructSingle(src)
		return
	}

	self.fmtStructMulti(src)
}

func (self *Fmt) fmtStructField(src r.Value, field r.StructField) {
	self.AppendString(field.Name)
	self.AppendString(`: `)
	self.fmtAny(field.Type, src)
}

func (self *Fmt) fmtStructSingle(src r.Value) {
	if isStructUnit(src.Type()) {
		self.fmtStructSingleAnon(src)
		return
	}
	self.fmtStructSingleNamed(src)
}

func (self *Fmt) fmtStructSingleAnon(src r.Value) {
	head := src.Field(0)

	self.AppendByte('{')
	if !self.skipField(head) {
		self.fmtAny(structHeadType(src.Type()), head)
	}
	self.AppendByte('}')
}

func (self *Fmt) fmtStructSingleNamed(src r.Value) {
	self.AppendByte('{')

	var found bool

	for _, field := range gg.StructPublicFieldCache.Get(src.Type()) {
		src := src.FieldByIndex(field.Index)
		if self.skipField(src) {
			continue
		}

		if found {
			self.AppendString(`, `)
		}
		found = true

		self.fmtStructField(src, field)
	}

	self.AppendByte('}')
}

func (self *Fmt) fmtStructMulti(src r.Value) {
	if isStructUnit(src.Type()) {
		self.fmtStructMultiAnon(src)
		return
	}
	self.fmtStructMultiNamed(src)
}

func (self *Fmt) fmtStructMultiAnon(src r.Value) {
	head := src.Field(0)

	self.AppendByte('{')

	if !self.skipField(head) {
		defer self.lvlInc().Done()
		self.fmtAny(structHeadType(src.Type()), head)
	}

	self.AppendByte('}')
}

func (self *Fmt) fmtStructMultiNamed(src r.Value) {
	fields := gg.StructPublicFieldCache.Get(src.Type())

	if self.SkipZeroFields() {
		test := func(field r.StructField) bool {
			return !src.FieldByIndex(field.Index).IsZero()
		}

		count := gg.Count(fields, test)

		if count == 0 {
			self.AppendString(`{}`)
			return
		}

		if count == 1 {
			field := gg.Find(fields, test)
			self.fmtStructMultiNamedUnit(src.FieldByIndex(field.Index), field)
			return
		}
	}

	self.fmtStructMultiNamedLines(src, fields)
}

func (self *Fmt) fmtStructMultiNamedUnit(src r.Value, field r.StructField) {
	self.AppendByte('{')
	self.fmtStructField(src, field)
	self.AppendByte('}')
}

func (self *Fmt) fmtStructMultiNamedLines(src r.Value, fields []r.StructField) {
	self.AppendByte('{')
	self.AppendNewline()
	snap := self.lvlInc()

	for _, field := range fields {
		src := src.FieldByIndex(field.Index)
		if self.skipField(src) {
			continue
		}

		self.fmtIndent()
		self.fmtStructField(src, field)
		self.AppendByte(',')
		self.AppendNewline()
	}

	snap.Done()
	self.fmtIndent()
	self.AppendByte('}')
}

func (self *Fmt) fmtUnfmtable(typ r.Type, src r.Value) {
	if self.fmtedNil(typ, src) {
		return
	}

	self.fmtTypeName(src.Type())
	self.AppendByte('(')
	self.fmtUint64Hex(uint64(src.Pointer()))
	self.AppendByte(')')
}

func (self *Fmt) fmtTypeArg(typ r.Type) {
	if isTypeDefaultForLiteral(typ) {
		return
	}

	self.AppendByte('[')
	self.fmtTypeName(typ)
	self.AppendByte(']')
}

func (self *Fmt) fmtReflectType(src r.Value) {
	self.fmtIdent(`gg`, `Type`)
	self.AppendByte('[')
	self.fmtTypeName(src.Interface().(r.Type))
	self.AppendString(`]()`)
}

func (self *Fmt) fmtUint64Hex(val uint64) {
	self.AppendString(`0x`)
	self.AppendUint64Hex(val)
}

func (self *Fmt) fmtIndent() { self.AppendStringN(self.Indent, self.Lvl) }

func (self *Fmt) fmtIdent(pkg, name string) {
	if self.Pkg != pkg {
		self.AppendString(pkg)
		self.AppendByte('.')
	}
	self.AppendString(name)
}

func (self *Fmt) fmtTypeNameOpt(typ r.Type) {
	if !self.ElideType {
		self.fmtTypeName(typ)
	}
}

func (self *Fmt) fmtTypeName(typ r.Type) {
	self.AppendString(self.typeName(typ))
}

func (self *Fmt) fmtConvOpen(typ r.Type) *Fmt {
	self.fmtTypeName(typ)
	self.AppendByte('(')
	return self
}

func (self *Fmt) fmtConvOpt(outer, inner r.Type, excl ...r.Type) *Fmt {
	if !isTypeInterface(outer) || gg.Has(excl, inner) {
		return nil
	}
	return self.fmtConvOpen(inner)
}

func (self *Fmt) fmtConvClose() {
	// The nil check is relevant for `defer`. See `.fmtConvOpt`.
	if self != nil {
		self.AppendByte(')')
	}
}

func (self *Fmt) setElideType(val bool) gg.Snapshot[bool] {
	return gg.PtrSwap(&self.ElideType, val)
}

func (self *Fmt) lvlInc() gg.Snapshot[int] {
	return gg.PtrSwap(&self.Lvl, self.Lvl+1)
}

func (self *Fmt) skipField(src r.Value) bool {
	return self.SkipZeroFields() && src.IsZero()
}

func (self *Fmt) typeName(typ r.Type) string {
	return elidePkg(typeName(typ), self.Pkg)
}

func typeName(typ r.Type) string { return string(typeNameCache.Get(typ)) }

var typeNameCache = gg.TypeCacheOf[typeNameStr]()

type typeNameStr string

func (self *typeNameStr) Init(typ r.Type) {
	if typ == nil {
		return
	}

	tar := typ.String()

	/**
	Some types must be wrapped in parens because we use the resulting type name
	in expression context, not in type context. Wrapping avoids ambiguity with
	value expression syntax.
	*/
	if typ.Kind() == r.Func && strings.HasPrefix(tar, `func(`) ||
		typ.Kind() == r.Pointer && strings.HasPrefix(tar, `*`) {
		tar = `(` + tar + `)`
	}

	tar = strings.ReplaceAll(tar, `interface {}`, `any`)
	tar = reUint8.Get().ReplaceAllString(tar, `byte`)
	*self = typeNameStr(tar)
}

var reUint8 = gg.NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`\buint8\b`)
})

func elidePkg(src, pkg string) string {
	if pkg == `` {
		return src
	}

	tar := strings.TrimPrefix(src, pkg)
	if len(src) != len(tar) && len(tar) > 0 && tar[0] == '.' {
		return tar[1:]
	}
	return src
}

func canAmpersand(kind r.Kind) bool {
	return kind == r.Array || kind == r.Slice || kind == r.Struct
}

func isTypeDefaultForLiteral(typ r.Type) bool {
	switch typ {
	case nil, typeBool, typeInt, typeString, typeComplex64, typeComplex128:
		return true
	default:
		return false
	}
}

func isValueNilInterface(src r.Value) bool {
	return !src.IsValid() || isValueInterface(src) && src.IsNil()
}

func isValueInterface(src r.Value) bool {
	return src.Kind() == r.Interface
}

func isStructUnit(typ r.Type) bool { return typ.NumField() == 1 }

func isTypeInterface(typ r.Type) bool { return gg.TypeKind(typ) == r.Interface }

func structHeadType(typ r.Type) r.Type {
	return gg.StructPublicFieldCache.Get(typ)[0].Type
}

func derefIface(val r.Value) r.Value {
	for val.Kind() == r.Interface {
		val = val.Elem()
	}
	return val
}

func textPrintability[A gg.Text](src A) (out printability) {
	out.init(gg.ToString(src))
	return
}

type printability struct {
	errors, lineFeeds, carriageReturns, backquotes, escapes, unprintables bool
}

func (self *printability) init(src string) {
	for _, val := range src {
		/**
		`unicode.IsPrint` uses `unicode.S` which includes `utf8.RuneError`.
		As a result, it considers error runes printable, which is wildly
		inappropriate for our purposes. So we have to handle it separately.
		*/
		if val == utf8.RuneError {
			self.errors = true
			return
		}

		if val == '\n' {
			self.lineFeeds = true
			continue
		}

		if val == '\r' {
			self.carriageReturns = true
			continue
		}

		if val == '`' {
			self.backquotes = true
			continue
		}

		if int(val) < len(stringEsc) && stringEsc[byte(val)] {
			self.escapes = true
			continue
		}

		if !unicode.IsPrint(val) {
			self.unprintables = true
			return
		}
	}
}

// https://go.dev/ref/spec#String_literals
var stringEsc = [256]bool{
	byte('\a'): true,
	byte('\b'): true,
	byte('\f'): true,
	byte('\n'): true,
	byte('\r'): true,
	byte('\t'): true,
	byte('\v'): true,
	byte('\\'): true,
	byte('"'):  true,
}
