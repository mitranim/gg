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
	* `.Pkg`, if set, indicates the package name to strip from type names.
*/
type Conf struct {
	Indent     string
	ZeroFields bool
	Pkg        string
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

/*
Similar to `fmt.Sprintf("%#v")` or `gg.GoString`, but more advanced.
Formats the input as Go code, using the config `ConfDefault`, returning
the resulting string.
*/
func String[A any](src A) string { return StringIndent(src, 0) }

/*
Similar to `fmt.Sprintf("%#v")` or `gg.GoString`, but more advanced.
Formats the input as Go code, using the config `ConfDefault`, returning
the resulting bytes.
*/
func Bytes[A any](src A) []byte { return BytesIndent(src, 0) }

/*
Formats the input as Go code, using the given config, returning the
resulting string.
*/
func StringC[A any](conf Conf, src A) string { return StringIndentC(conf, src, 0) }

/*
Formats the input as Go code, using the given config, returning the
resulting bytes.
*/
func BytesC[A any](conf Conf, src A) []byte { return BytesIndentC(conf, src, 0) }

/*
Formats the input as Go code, using the default config with the given
indentation level, returning the resulting string.
*/
func StringIndent[A any](src A, lvl int) string {
	return StringIndentC(ConfDefault, src, lvl)
}

/*
Formats the input as Go code, using the default config with the given
indentation level, returning the resulting bytes.
*/
func BytesIndent[A any](src A, lvl int) []byte {
	return BytesIndentC(ConfDefault, src, lvl)
}

/*
Formats the input as Go code, using the given config with the given indentation
level, returning the resulting string.
*/
func StringIndentC[A any](conf Conf, src A, lvl int) string {
	return gg.ToString(BytesIndentC(conf, src, lvl))
}

/*
Formats the input as Go code, using the given config with the given indentation
level, returning the resulting bytes.
*/
func BytesIndentC[A any](conf Conf, src A, lvl int) []byte {
	buf := conf.Fmt()
	buf.Lvl += lvl
	buf.fmtAny(gg.Type[A](), r.ValueOf(gg.AnyNoEscUnsafe(src)))
	return buf.Buf
}

/*
Shortcut for printing the input as Go code, prefixed with the given description,
using the default config. Handy for debug-printing.
*/
func Prn[A any](desc string, src A) { fmt.Println(desc, String(src)) }

// Shortcut for printing the input as Go code, using the default config.
func Println[A any](src A) { fmt.Println(String(src)) }

// Shortcut for printing the input as Go code, using the given config.
func PrintlnC[A any](conf Conf, src A) { fmt.Println(StringC(conf, src)) }
