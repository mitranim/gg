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
`ScanVals` and `ScanVal` which are statically typed. Output must be a non-nil
pointer to one of the following:

	* Slice of scalars.
	* Slice of structs.
	* Single scalar.
	* Single struct.

Always closes the rows. If output is not a slice, verifies that there is EXACTLY
one row in total, otherwise panics.
*/
func ScanAny[Src Rows](src Src, out any) {
	ScanReflect(src, r.ValueOf(out))
}

// Variant of `ScanAny` that takes a reflect value rather than `any`.
func ScanReflect[Src Rows](src Src, out r.Value) {
	if out.Kind() != r.Pointer {
		panic(gg.Errf(`scan destination must be a pointer, got %q`, out.Type()))
	}
	if out.IsNil() {
		panic(gg.Errf(`scan destination must be non-nil, got nil %q`, out.Type()))
	}
	out = gg.ValueDerefAlloc(out)

	if out.Kind() == r.Slice {
		scanValsReflect(src, out)
	} else {
		scanValReflect(src, out)
	}
}
