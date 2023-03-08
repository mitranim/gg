package gsql

import (
	"database/sql/driver"
	"strings"

	"github.com/mitranim/gg"
)

/*
Variant of `string` intended as an operand for SQL "like" and "ilike" operators.
When generating an SQL argument via `.Value`, the string is wrapped in `%` to
ensure partial match, escaping any pre-existing `%` and `_`. As a special case,
an empty string is used as-is, and doesn't match anything when used with
`like` or `ilike`.
*/
type Like string

// Implement `fmt.Stringer`. Returns the underlying string unchanged.
func (self Like) String() string { return string(self) }

// Implement `driver.Valuer`, returning the escaped string from `.Esc`.
func (self Like) Value() (driver.Value, error) { return self.Esc(), nil }

// Implement `sql.Scanner`.
func (self *Like) Scan(src any) error {
	str, ok := gg.AnyToText[string](src)
	if ok {
		*self = Like(str)
		return nil
	}
	return gg.ErrConv(src, gg.Type[Like]())
}

/*
Returns an escaped string suitable as an operand for SQL "like" or "ilike".
As a special case, an empty string is returned as-is.
*/
func (self Like) Esc() string {
	if self == `` {
		return ``
	}
	return `%` + replaceLike.Replace(string(self)) + `%`
}

var replaceLike = strings.NewReplacer(`%`, `\%`, `_`, `\_`)
