package gsql

import (
	"database/sql/driver"
	"encoding/json"
	r "reflect"

	"github.com/mitranim/gg"
)

/*
Shortcut for converting an arbitrary rune-like to `Rune`.
Useful for higher-order functions such as `gg.Map`.
*/
func RuneFrom[A ~rune](val A) Rune { return Rune(val) }

/*
Variant of Go `rune` compatible with text, JSON, and SQL (with caveats). In text
and JSON, behaves like `string`. In Go and SQL, behaves like `rune`/`int32`. As
a special case, zero value is considered empty in text, and null in JSON and
SQL. When parsing, input must be empty or single char.

Some databases, or their Go drivers, may not support representing chars as
int32. For example, Postgres doesn't have an analog of Go `rune`. Its "char"
type is a variable-sized string. This type is not compatible with such
databases.
*/
type Rune rune

// Implement `gg.Nullable`. True if zero value.
func (self Rune) IsNull() bool { return self == 0 }

// Inverse of `.IsNull`.
func (self Rune) IsNotNull() bool { return !self.IsNull() }

// Implement `Clearer`. Zeroes the receiver.
func (self *Rune) Clear() { *self = 0 }

/*
Implement `fmt.Stringer`. If zero, returns an empty string. Otherwise returns
a string containing exactly one character.
*/
func (self Rune) String() string {
	if self.IsNull() {
		return ``
	}
	return string(self)
}

// Implement `AppenderTo`, appending the same representation as `.String`.
func (self Rune) AppendTo(buf []byte) []byte {
	if self.IsNull() {
		return buf
	}
	return append(buf, self.String()...)
}

/*
Implement `Parser`. If the input is empty, clears the receiver via `.Clear`. If
the input has more than one character, returns an error. Otherwise uses the
first and only character from the input.
*/
func (self *Rune) Parse(src string) error {
	chars := []rune(src)
	if len(chars) == 0 {
		self.Clear()
		return nil
	}

	if len(chars) > 1 {
		return gg.Errf(`unable to parse %q as char: too many chars`, src)
	}

	*self = Rune(chars[0])
	return nil
}

// Implement `encoding.TextMarshaler`, returning the same representation as `.String`.
func (self Rune) MarshalText() ([]byte, error) {
	return gg.ToBytes(self.String()), nil
}

// Implement `encoding.TextUnmarshaler`, using the same logic as `.Parse`.
func (self *Rune) UnmarshalText(src []byte) error {
	return self.Parse(gg.ToString(src))
}

/*
Implement `json.Marshaler`. If `.IsNull`, returns a representation of JSON null.
Otherwise uses an equivalent of `json.Marshal(self.String())`.
*/
func (self Rune) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return gg.ToBytes(`null`), nil
	}

	if self == '"' {
		return gg.ToBytes(`"\""`), nil
	}
	return gg.ToBytes(`"` + string(rune(self)) + `"`), nil
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON null,
clears the receiver via `.Clear`. Otherwise requires the input to be a JSON
string and decodes it via `.Parse`.
*/
func (self *Rune) UnmarshalJSON(src []byte) error {
	if gg.IsJsonEmpty(src) {
		self.Clear()
		return nil
	}

	// Inefficient, TODO tune.
	var tar string
	err := json.Unmarshal(src, gg.AnyNoEscUnsafe(&tar))
	if err != nil {
		return err
	}

	return self.Parse(tar)
}

/*
Implement SQL `driver.Valuer`. If `.IsNull`, returns nil. Otherwise returns
rune.
*/
func (self Rune) Value() (driver.Value, error) {
	if self.IsNull() {
		return nil, nil
	}
	return rune(self), nil
}

/*
Implement SQL `Scanner`, decoding arbitrary input, which must be one of:

	* Nil  -> use `.Clear`.
	* Text -> use `.Parse`.
	* Rune -> assign as-is.
*/
func (self *Rune) Scan(src any) error {
	if src == nil {
		self.Clear()
		return nil
	}

	str, ok := gg.AnyToText[string](src)
	if ok {
		return self.Parse(str)
	}

	val := r.ValueOf(gg.AnyNoEscUnsafe(src))
	if val.Kind() == r.Int32 {
		*self = Rune(val.Int())
		return nil
	}

	return gg.ErrConv(src, gg.Type[Rune]())
}
