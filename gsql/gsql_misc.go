package gsql

import "github.com/mitranim/gg"

/*
Must be deferred. Commit if there was no panic, rollback if there was a
panic. Usage:

	defer DbTxDone(conn)
*/
func DbTxDone[A DbTx](val A) {
	err := gg.ToErrTraced(recover(), 1)

	if err != nil {
		_ = val.Rollback()
		panic(err)
	}

	defer gg.Detailf(`failed to commit DB transaction`)
	gg.Try(val.Commit())
}
