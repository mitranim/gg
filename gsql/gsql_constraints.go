package gsql

import (
	"context"
	"database/sql"
	"io"
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
	Errer
	Nexter
	ColumnerScanner
}

type (
	Columner interface{ Columns() ([]string, error) }
	Errer    interface{ Err() error }
	Scanner  interface{ Scan(...any) error }
	Nexter   interface{ Next() bool }
)

type ColumnerScanner interface {
	Columner
	Scanner
}
