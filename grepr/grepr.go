/*
Missing feature of the standard library: printing arbitrary inputs as Go code,
with proper spacing and support for multi-line output with indentation. The
name "repr" stands for "representation" and alludes to the Python function with
the same name.
*/
package grepr

import (
	"fmt"
	r "reflect"
	u "unsafe"

	"github.com/mitranim/gg"
)

// Default config used by top-level formatting functions in this package.
var ConfDefault = Conf{Indent: gg.Indent}

// Config that allows formatting of struct zero fields.
var ConfFull = Conf{Indent: gg.Indent, ZeroFields: true}

/*
Formatting config.

	* `.Indent` controls indentation. If empty, output is single line.

	* `.ZeroFields`, if set, forces printing of zero fields in structs.
		By default zero fields are skipped.
*/
type Conf struct {
	Indent     string
	ZeroFields bool
}

/*
Short for "is single line". If `.Indent` is empty, this is true, and output
is single-line. Otherwise output is multi-line.
*/
func (self Conf) IsSingle() bool { return self.Indent == `` }

// Short for "is multi line". Inverse of `.IsSingle`.
func (self Conf) IsMulti() bool { return self.Indent != `` }

// Inverse of `.ZeroFields`.
func (self Conf) SkipZeroFields() bool { return !self.ZeroFields }

/*
Shortcut for printing the input as Go code, using this config, prefixed with the
given description. Handy for debug-printing.
*/
func (self Conf) Prn(desc string, src any) { fmt.Println(desc, self.AnyString(src)) }

// Shortcut for printing the input as Go code, using this config
func (self Conf) Println(src any) { fmt.Println(self.AnyString(src)) }

// Shortcut for using this config to format the input as Go code.
func (self Conf) AnyString(src any) string {
	buf := self.Fmt()
	buf.Any(src)
	return buf.String()
}

// Shortcut for using this config to format the input as Go code.
func (self Conf) ValueString(src r.Value) string {
	buf := self.Fmt()
	buf.Value(src)
	return buf.String()
}

// Shortcut for creating a pretty-formatter with this config.
func (self Conf) Fmt() Fmt {
	var buf Fmt
	buf.Conf = self
	return buf
}

// Short for "formatter".
type Fmt struct {
	Conf
	gg.Buf
	Lvl       int
	ElideType bool
	Visited   gg.Set[u.Pointer]
}

// Formats the input into the inner buffer, using the inner config.
func (self *Fmt) Any(src any) { fmtAny(self, r.ValueOf(src)) }

// Formats the input into the inner buffer, using the inner config.
func (self *Fmt) Value(src r.Value) { fmtAny(self, src) }

/*
Similar to `fmt.Sprintf("%#v")` or `gg.GoString`, but more advanced. Formats
the input as Go code, using the default config.
*/
func String(src any) string { return StringIndent(src, 0) }

/*
Formats the input as Go code, using the default config with the given
indentation level.
*/
func StringIndent(src any, lvl int) string {
	buf := ConfDefault.Fmt()
	buf.Lvl += lvl
	buf.Any(src)
	return buf.String()
}

/*
Shortcut for printing the input as Go code, prefixed with the given description,
using the default config. Handy for debug-printing.
*/
func Prn(desc string, src any) { ConfDefault.Prn(desc, src) }

// Shortcut for printing the input as Go code, using the default config.
func Println(src any) { ConfDefault.Println(src) }

// Corrected version of `strconv.CanBackquote` that allows newlines.
func CanBackquote[A gg.Text](src A) bool {
	for _, char := range gg.ToString(src) {
		if isNonBackquotable(char) {
			return false
		}
	}
	return true
}
