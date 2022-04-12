package gsql

import "github.com/mitranim/gg"

/*
Must be deferred. Commit if there was no panic, rollback if there was a
panic. Usage:

	defer DbTxDone(conn)
*/
func DbTxDone[A DbTx](val A) {
	DbTxDoneWith(val, gg.AnyErrTraced(recover()))
}

/*
Commit if there was no error, rollback if there was an error.
Used internally by `DbTxDone`.
*/
func DbTxDoneWith[A DbTx](val A, err error) {
	if err != nil {
		_ = val.Rollback()
		panic(err)
	}

	defer gg.Detailf(`failed to commit DB transaction`)
	gg.Try(val.Commit())
}
