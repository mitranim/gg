package gsql

import (
	"database/sql"
	r "reflect"
	"strings"
	"time"

	"github.com/mitranim/gg"
)

func popSqlArrSegment(ptr *string) string {
	var lvl int
	src := *ptr

	for ind, char := range src {
		if char == '{' {
			lvl++
			continue
		}

		if char == '}' {
			lvl--
			if lvl < 0 {
				panic(gg.ErrInvalidInput)
			}
			continue
		}

		if char == ',' && lvl == 0 {
			*ptr = src[ind+1:]
			return src[:ind]
		}
	}

	*ptr = ``
	return src
}

func typeReferenceField(typ r.Type) (_ r.StructField, _ bool) {
	if typ.Kind() == r.Struct && typ.NumField() > 0 {
		field := gg.StructFieldCache.Get(typ)[0]
		if field.Tag.Get(`role`) == `ref` {
			return field, true
		}
	}
	return
}

func isTypeNonScannableStruct(typ r.Type) bool {
	return typ.Kind() == r.Struct &&
		!typ.ConvertibleTo(gg.Type[time.Time]()) &&
		!typ.Implements(gg.Type[gg.Scanner]()) &&
		!r.PointerTo(typ).Implements(gg.Type[gg.Scanner]())
}

var typeMetaCache = gg.TypeCacheOf[typeMeta]()

type typeMeta struct {
	typ  r.Type
	dict map[string][]int
}

func (self typeMeta) Get(key string) []int {
	val, ok := self.dict[key]
	if !ok {
		panic(gg.Errf(`unknown column %q in type %v`, key, self.typ))
	}
	return val
}

func (self typeMeta) IsScalar() bool { return self.dict == nil }

// Called by `TypeCache`.
func (self *typeMeta) Init(typ r.Type) {
	self.typ = typ
	self.addAny(nil, nil, typ)
}

func (self *typeMeta) addAny(index []int, cols []string, typ r.Type) {
	field, ok := typeReferenceField(typ)
	if ok {
		self.addAny(gg.Concat(index, field.Index), cols, field.Type)
	}

	if isTypeNonScannableStruct(typ) {
		self.addStruct(index, cols, typ)
		return
	}

	if len(cols) > 0 {
		self.initMap()[strings.Join(cols, `.`)] = index
	}
}

func (self *typeMeta) addStruct(index []int, cols []string, typ r.Type) {
	self.initMap()
	for _, field := range gg.StructPublicFieldCache.Get(typ) {
		self.addField(index, cols, field)
	}
}

func (self *typeMeta) addField(index []int, cols []string, field r.StructField) {
	col := gg.FieldDbName(field)
	typ := field.Type

	if col != `` {
		self.addAny(gg.Concat(index, field.Index), gg.CloneAppend(cols, col), typ)
		return
	}

	if !field.Anonymous {
		return
	}

	if typ.Kind() == r.Struct {
		self.addStruct(gg.Concat(index, field.Index), cols, typ)
		return
	}

	panic(gg.Errf(
		`unsupported embedded type %q; embedded fields must be structs`,
		typ,
	))
}

func (self *typeMeta) initMap() map[string][]int {
	return gg.MapInit(&self.dict)
}

func scanNextScalar[Row any, Src ColumnerScanner](src Src) (out Row) {
	gg.Try(src.Scan(gg.AnyNoEscUnsafe(&out)))
	return
}

func scanNextStruct[Row any, Src ColumnerScanner](src Src) (out Row) {
	scanStructReflect(src, r.ValueOf(gg.AnyNoEscUnsafe(&out)).Elem())
	return
}

/*
TODO needs performance tuning.

Would be nice to use an implementation similar to this:

	gg.Try(src.Scan(gg.Map(RowsCols(src), func(key string) any {
		return tar.FieldByIndex(meta.Get(key)).Addr().Interface()
	})...))

...But the SQL driver doesn't allow to decode SQL "null" into non-nullable
destinations such as `string` fields. This behavior is inconsistent with
JSON, and unfortunate for our purposes.
*/
func scanStructReflect[Src ColumnerScanner](src Src, tar r.Value) {
	typ := tar.Type()
	meta := typeMetaCache.Get(typ)
	cols := gg.Try1(src.Columns())
	indir := gg.Map(cols, func(key string) r.Value {
		return r.New(r.PointerTo(typ.FieldByIndex(meta.Get(key)).Type))
	})

	gg.Try(src.Scan(gg.Map(indir, r.Value.Interface)...))

	gg.Each2(cols, indir, func(key string, val r.Value) {
		val = val.Elem()
		if !val.IsNil() {
			tar.FieldByIndex(meta.Get(key)).Set(val.Elem())
		}
	})
}

func scanValsReflect[Src Rows](src Src, tar r.Value) {
	defer gg.Close(src)

	for src.Next() {
		const off = 1

		// Increase length by one, effectively appending a zero value to the slice.
		// Similar to `r.Append(r.New(typ).Elem())`, but should be marginally more
		// efficient.
		ind := tar.Len()
		tar.Grow(off)
		tar.SetLen(ind + off)

		// Settable, addressable reference to the newly appended zero value.
		out := tar.Index(ind)

		// Hide new value from consumer code until scan is successful.
		tar.SetLen(ind)

		scanReflect(src, out)

		// After successful scan, we can reveal new element to consumer code.
		tar.SetLen(ind + off)
	}

	gg.ErrOk(src)
}

func scanValReflect[Src Rows](src Src, tar r.Value) {
	defer gg.Close(src)

	if !src.Next() {
		panic(gg.AnyErrTraced(sql.ErrNoRows))
	}

	scanReflect(src, tar)
	gg.ErrOk(src)

	if src.Next() {
		panic(gg.AnyErrTraced(ErrMultipleRows))
	}
}

func scanReflect[Src ColumnerScanner](src Src, tar r.Value) {
	if isValueScalar(tar) {
		scanScalarReflect(src, tar)
		return
	}
	scanStructReflect(src, tar)
}

func scanScalarReflect[Src ColumnerScanner](src Src, tar r.Value) {
	gg.Try(src.Scan(tar.Addr().Interface()))
}

func isScalar[A any]() bool {
	return typeMetaCache.Get(gg.Type[A]()).IsScalar()
}

func isValueScalar(val r.Value) bool {
	return typeMetaCache.Get(val.Type()).IsScalar()
}
