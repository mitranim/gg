package gg

import (
	"database/sql/driver"
	"encoding/json"
	r "reflect"
	"strconv"
	"time"
)

// Calls `time.Now` and converts to `TimeMicro`, truncating precision.
func TimeMicroNow() TimeMicro { return TimeMicro(time.Now().UnixMicro()) }

// Shortcut for parsing text into `TimeMicro`. Panics on error.
func TimeMicroParse[A Text](src A) TimeMicro {
	var out TimeMicro
	Try(out.Parse(ToString(src)))
	return out
}

/*
Represents a Unix timestamp in microseconds. In text and JSON, this type
supports parsing numeric timestamps and RFC3339 timestamps, but always encodes
as a number. In SQL, this type is represented in the RFC3339 format. This type
is "zero-optional" or "zero-nullable". The zero value is considered empty in
text and null in JSON/SQL. Conversion to `time.Time` doesn't specify a
timezone, which means it uses `time.Local` by default. If you prefer UTC,
enforce it across the app by updating `time.Local`.

Caution: corresponding DB columns MUST be restricted to microsecond precision.
Without this restriction, encoding and decoding is not reversible. After losing
precision to an encoding-decoding roundtrip, you might be unable to find a
corresponding value in a database, if timestamp precision is higher than a
microsecond.

Also see `TimeMilli`, which uses milliseconds.
*/
type TimeMicro int64

// Implement `Nullable`. True if zero.
func (self TimeMicro) IsNull() bool { return self == 0 }

// Implement `Clearer`, zeroing the receiver.
func (self *TimeMicro) Clear() {
	if self != nil {
		*self = 0
	}
}

/*
Convert to `time.Time` by calling `time.UnixMicro`. The resulting timestamp has
the timezone `time.Local`. To enforce UTC, modify `time.Local` at app startup,
or call `.In(time.UTC)`.
*/
func (self TimeMicro) Time() time.Time { return time.UnixMicro(int64(self)) }

/*
Implement `AnyGetter` for compatibility with some 3rd party libraries. If zero,
returns `nil`, otherwise creates `time.Time` by calling `TimeMicro.Time`.
*/
func (self TimeMicro) Get() any {
	if self.IsNull() {
		return nil
	}
	return self.Time()
}

// Sets the receiver to the given input.
func (self *TimeMicro) SetInt64(val int64) { *self = TimeMicro(val) }

// Sets the receiver to the result of `time.Time.UnixMicro`.
func (self *TimeMicro) SetTime(val time.Time) { self.SetInt64(val.UnixMicro()) }

/*
Implement `Parser`. The input must be either an integer in base 10, representing
a Unix millisecond timestamp, or an RFC3339 timestamp. RFC3339 is the default
time encoding/decoding format in Go and some other languages.
*/
func (self *TimeMicro) Parse(src string) error {
	if len(src) == 0 {
		self.Clear()
		return nil
	}

	if isIntString(src) {
		num, err := strconv.ParseInt(src, 10, 64)
		if err != nil {
			return err
		}
		self.SetInt64(num)
		return nil
	}

	inst, err := time.Parse(time.RFC3339, src)
	if err != nil {
		return err
	}
	self.SetTime(inst)
	return nil
}

/*
Implement `fmt.Stringer`. If zero, returns an empty string. Otherwise returns
the base 10 representation of the underlying number.
*/
func (self TimeMicro) String() string {
	if self.IsNull() {
		return ``
	}
	return strconv.FormatInt(int64(self), 10)
}

// Implement `Appender`, using the same representation as `.String`.
func (self TimeMicro) Append(buf []byte) []byte {
	if self.IsNull() {
		return buf
	}
	return strconv.AppendInt(buf, int64(self), 10)
}

/*
Implement `encoding.TextMarhaler`. If zero, returns nil. Otherwise returns the
same representation as `.String`.
*/
func (self TimeMicro) MarshalText() ([]byte, error) {
	if self.IsNull() {
		return nil, nil
	}
	return ToBytes(self.String()), nil
}

// Implement `encoding.TextUnmarshaler`, using the same algorithm as `.Parse`.
func (self *TimeMicro) UnmarshalText(src []byte) error {
	return self.Parse(ToString(src))
}

/*
Implement `json.Marshaler`. If zero, returns bytes representing `null`.
Otherwise encodes as a JSON number.
*/
func (self TimeMicro) MarshalJSON() ([]byte, error) {
	if self.IsNull() {
		return ToBytes(`null`), nil
	}
	return json.Marshal(int64(self))
}

/*
Implement `json.Unmarshaler`. If the input is empty or represents JSON `null`,
zeroes the receiver. If the input is a JSON number, parses it in accordance
with `.Parse`. Otherwise uses the default `json.Unmarshal` behavior for
`*time.Time` and stores the resulting timestamp in milliseconds.
*/
func (self *TimeMicro) UnmarshalJSON(src []byte) error {
	if ToString(src) == `null` {
		self.Clear()
		return nil
	}

	if isIntString(ToString(src)) {
		num, err := strconv.ParseInt(ToString(src), 10, 64)
		if err != nil {
			return err
		}
		self.SetInt64(num)
		return nil
	}

	var inst time.Time
	err := json.Unmarshal(src, &inst)
	if err != nil {
		return err
	}
	self.SetTime(inst)
	return nil
}

// Implement `driver.Valuer`, using `.Get`.
func (self TimeMicro) Value() (driver.Value, error) {
	return self.Get(), nil
}

/*
Implement `sql.Scanner`, converting an arbitrary input to `TimeMicro` and
modifying the receiver. Acceptable inputs:

	* `nil`         -> use `.Clear`
	* integer       -> assign, assuming milliseconds
	* text          -> use `.Parse`
	* `time.Time`   -> use `.SetTime`
	* `*time.Time`  -> use `.Clear` or `.SetTime`
	* `AnyGetter`   -> scan underlying value
*/
func (self *TimeMicro) Scan(src any) error {
	str, ok := AnyToText[string](src)
	if ok {
		return self.Parse(str)
	}

	switch src := src.(type) {
	case nil:
		self.Clear()
		return nil

	case time.Time:
		self.SetTime(src)
		return nil

	case *time.Time:
		if src == nil {
			self.Clear()
		} else {
			self.SetTime(*src)
		}
		return nil

	case int64:
		self.SetInt64(src)
		return nil

	case TimeMicro:
		*self = src
		return nil

	default:
		val := r.ValueOf(src)
		if val.CanInt() {
			self.SetInt64(val.Int())
			return nil
		}
		return ErrConv(src, Type[TimeMicro]())
	}
}
