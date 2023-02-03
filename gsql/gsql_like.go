package gsql

import (
	"database/sql/driver"
	"strings"
)

/*
Variant of `string` intended as an operand for SQL "like" and "ilike" operators.
When generating an SQL argument via `.Value`, the string is wrapped in `%` to
ensure partial match, escaping any pre-existing `%` and `_`.
*/
type Like string

// Implement `gg.Nullable`. True if the string is empty.
func (self Like) IsNull() bool { return self == `` }

// Implement `fmt.Stringer`. Returns the underlying string unchanged.
func (self Like) String() string { return string(self) }

// Implement `driver.Valuer`, returning the escaped string from `.Esc`.
func (self Like) Value() (driver.Value, error) {
	if self.IsNull() {
		return nil, nil
	}
	return self.Esc(), nil
}

// Returns an escaped string suitable as an operand for SQL "like" or "ilike".
func (self Like) Esc() string {
	if self.IsNull() {
		return ``
	}
	return `%` + replaceLike.Replace(string(self)) + `%`
}

var replaceLike = strings.NewReplacer(`%`, `\%`, `_`, `\_`)
