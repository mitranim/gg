package gsql

import (
	"database/sql"
	r "reflect"

	"github.com/mitranim/gg"
)

// Returned by `ScanVal` and `ScanAny` when there are too many rows.
const ErrMultipleRows gg.ErrStr = `expected one row, got multiple`

/*
Takes `Rows` and decodes them into a slice of the given type, using `ScanNext`
for each row. Output type must be either scalar or struct. Always closes the
rows.
*/
func ScanVals[Row any, Src Rows](src Src) (out []Row) {
	defer gg.Close(src)
	for src.Next() {
		out = append(out, ScanNext[Row](src))
	}
	gg.ErrOk(src)
	return
}

/*
Takes `Rows` and decodes the first row into a value of the given type, using
`ScanNext` once. The rows must consist of EXACTLY one row, otherwise this
panics. Output type must be either scalar or struct. Always closes the rows.
*/
func ScanVal[Row any, Src Rows](src Src) Row {
	defer gg.Close(src)

	if !src.Next() {
		panic(gg.AnyErrTraced(sql.ErrNoRows))
	}

	out := ScanNext[Row](src)
	gg.ErrOk(src)

	if src.Next() {
		panic(gg.AnyErrTraced(ErrMultipleRows))
	}
	return out
}

/*
Takes `Rows` and decodes the next row into a value of the given type. Output
type must be either scalar or struct. Panics on errors. Must be called only
after `Rows.Next`.
*/
func ScanNext[Row any, Src ColumnerScanner](src Src) Row {
	if isScalar[Row]() {
		return scanNextScalar[Row](src)
	}
	return scanNextStruct[Row](src)
}

/*
Decodes `Rows` into the given dynamically typed output. Counterpart to
`ScanVals` and `ScanVal` which are statically typed. The output must be
a non-nil pointer, any amount of levels deep, to one of the following:

  - Slice of scalars.
  - Slice of structs.
  - Single scalar.
  - Single struct.
  - Interface value hosting a concrete type.

Always closes the rows. If the output is not a slice, verifies that there is
EXACTLY one row in total, otherwise panics.
*/
func ScanAny[Src Rows](src Src, out any) {
	ScanReflect(src, r.ValueOf(out))
}

// Variant of `ScanAny` that takes `reflect.Value` rather than `any`.
func ScanReflect[Src Rows](src Src, out r.Value) {
	tar, iface := derefAlloc(out)

	if tar.Kind() == r.Slice {
		scanValsReflect(src, tar)
	} else {
		scanValReflect(src, tar, true)
	}

	if iface.CanSet() {
		iface.Set(tar.Convert(iface.Type()))
	}
}

/*
Similar to `ScanAny`, but when scanning into a single value (not a slice),
doesn't panic if there are zero rows, leaving the destination unchanged.
When scanning into a slice, behaves exactly like `ScanAny`.
*/
func ScanAnyOpt[Src Rows](src Src, out any) {
	ScanReflectOpt(src, r.ValueOf(out))
}

// Variant of `ScanAnyOpt` that takes `reflect.Value` rather than `any`.
func ScanReflectOpt[Src Rows](src Src, out r.Value) {
	tar, iface := derefAlloc(out)

	if tar.Kind() == r.Slice {
		scanValsReflect(src, tar)
	} else {
		scanValReflect(src, tar, false)
	}

	if iface.CanSet() {
		iface.Set(tar.Convert(iface.Type()))
	}
}
