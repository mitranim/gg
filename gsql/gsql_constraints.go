package gsql

import (
	"context"
	"database/sql"
	"io"

	"github.com/mitranim/gg"
)

// Implemented by stdlib types such as `sql.DB`.
type Db interface {
	DbConn
	DbTxer
}

// Implemented by stdlib types such as `sql.DB`.
type DbTxer interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

// Implemented by stdlib types such as `sql.Conn` and `sql.Tx`.
type DbConn interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

// Implemented by stdlib types such as `sql.Tx`.
type DbTx interface {
	DbConn
	Commit() error
	Rollback() error
}

// Interface of `sql.Rows`. Used by various scanning tools.
type Rows interface {
	io.Closer
	gg.Errer
	gg.Nexter
	ColumnerScanner
}

// Sub-interface of `Rows` used by `ScanNext`.
type ColumnerScanner interface {
	Columns() ([]string, error)
	Scan(...any) error
}
