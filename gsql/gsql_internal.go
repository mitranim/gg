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

type typeMeta map[string][]int

func (self typeMeta) Get(key string) []int {
	val, ok := self[key]
	if !ok {
		panic(gg.Errf(`unknown column %q`, key))
	}
	return val
}

func (self typeMeta) IsScalar() bool { return self == nil }

// Called by `TypeCache`.
func (self *typeMeta) Init(typ r.Type) { self.addAny(nil, nil, typ) }

//go:noinline
func (self *typeMeta) initMap() typeMeta { return gg.MapPtrInit(self) }

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
	for _, field := range gg.StructFieldCache.Get(typ) {
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

func scanNextScalar[Row any, Src ColumnerScanner](src Src) (out Row) {
	gg.Try(src.Scan(gg.AnyNoEscUnsafe(&out)))
	return
}

func scanNextStruct[Row any, Src ColumnerScanner](src Src, meta typeMeta) (out Row) {
	tar := r.ValueOf(gg.AnyNoEscUnsafe(&out)).Elem()
	scanStruct(src, meta, tar)
	return
}

/*
TODO needs performance tuning.

Would be nice to use an implementation similar to this:

	gg.Try(src.Scan(gg.Map(RowsCols(src), func(key string) any {
		return tar.FieldByIndex(meta.Get(key)).Addr().Interface()
	})...))

...But the SQL driver doesn't allow to decode "null" into non-nullable
destinations such as `string` fields. This behavior is inconsistent with
JSON, and unfortunate for our purposes.
*/
func scanStruct[Src ColumnerScanner](src Src, meta typeMeta, tar r.Value) {
	typ := tar.Type()
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

func scanValsAny[Src Rows](src Src, tar r.Value) {
	defer gg.Close(src)

	elem := tar.Type().Elem()
	meta := typeMetaCache.Get(elem)

	for src.Next() {
		tar.Set(r.Append(tar, scanNextAny(src, meta, elem)))
	}

	gg.ErrOk(src)
}

func scanValAny[Src Rows](src Src, tar r.Value) {
	defer gg.Close(src)

	if !src.Next() {
		panic(gg.ToErrTraced(sql.ErrNoRows, 1))
	}

	typ := tar.Type()
	tar.Set(scanNextAny(src, typeMetaCache.Get(typ), typ))
	gg.ErrOk(src)

	if src.Next() {
		panic(gg.ToErrTraced(ErrMultipleRows, 1))
	}
}

func scanNextAny[Src ColumnerScanner](src Src, meta typeMeta, typ r.Type) r.Value {
	if meta.IsScalar() {
		return scanNextScalarAny(src, typ)
	}
	return scanNextStructAny(src, meta, typ)
}

func scanNextScalarAny[Src ColumnerScanner](src Src, typ r.Type) r.Value {
	tar := r.New(typ)
	gg.Try(src.Scan(tar.Interface()))
	return tar.Elem()
}

func scanNextStructAny[Src ColumnerScanner](src Src, meta typeMeta, typ r.Type) r.Value {
	tar := gg.NewElem(typ)
	scanStruct(src, meta, tar)
	return tar
}
