package gsql

import (
	"database/sql"
	"io"
	r "reflect"

	"github.com/mitranim/gg"
)

func ScanVals[Val any, Src Rows](src Src) (out []Val) {
	defer RowsDone(src)
	for src.Next() {
		out = append(out, NextVal[Val](src))
	}
	RowsOk(src)
	return
}

func ScanVal[Val any, Src Rows](src Src) Val {
	defer RowsDone(src)
	if src.Next() {
		return NextVal[Val](src)
	}
	RowsOk(src)
	panic(gg.ToErrTraced(sql.ErrNoRows, 1))
}

func NextVal[Val any, Src ColumnerScanner](src Src) Val {
	meta := typeMetaCache.Get(gg.Type[Val]())
	if meta.IsScalar() {
		return nextScalar[Val](src)
	}
	return nextStruct[Val](src, meta)
}

func ScanAny[Src Rows](src Src, out any) {
	tar := gg.ValueDeref(r.ValueOf(out))

	if tar.Kind() == r.Slice {
		scanValsAny(src, tar)
	} else {
		scanValAny(src, tar)
	}
}

func RowsDone[A io.Closer](val A) { _ = val.Close() }

func RowsErr[A Errer](val A) error {
	return gg.ErrTraced(val.Err(), 1)
}

func RowsOk[A Errer](val A) {
	gg.Try(gg.ErrTraced(val.Err(), 1))
}

func RowsCols[A Columner](val A) []string {
	return gg.Try1(val.Columns())
}

func RowsScan[A Scanner](src A, buf []any) {
	gg.Try(src.Scan(gg.NoEscUnsafe(buf)...))
}
