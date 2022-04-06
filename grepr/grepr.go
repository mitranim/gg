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
var Default = Conf{Indent: gg.Indent}

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
	var buf Fmt
	buf.Conf = Default
	buf.Lvl += lvl
	buf.Any(src)
	return buf.String()
}

/*
Shortcut for printing a given value as Go code, prefixed with a description.
Very handy for debug-printing.
*/
func Prn(desc string, src any) {
	fmt.Println(desc, String(src))
}

// Corrected version of `strconv.CanBackquote` that allows newlines.
func CanBackquote[A gg.Text](src A) bool {
	for _, char := range gg.ToString(src) {
		if isNonBackquotable(char) {
			return false
		}
	}
	return true
}
