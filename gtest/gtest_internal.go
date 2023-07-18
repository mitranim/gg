package gtest

import (
	"fmt"
	r "reflect"
	"strings"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
)

// Suboptimal, TODO revise.
func reindent(src string) string {
	return gg.JoinLines(gg.Map(gg.SplitLines(src), indent)...)
}

func indent(src string) string {
	if src == `` {
		return src
	}
	return gg.Indent + src
}

// TODO rename and make public.
func goStringIndent[A any](val A) string { return grepr.StringIndent(val, 1) }

func errTrace(err error) string {
	return strings.TrimSpace(gg.ErrTrace(err).StringIndent(1))
}

/*
Should return `true` when stringifying the given value via `fmt.Sprint` produces
basically the same representation as pretty-printing it via `grepr`, with no
significant difference in information. We "discount" the string quotes in this
case. TODO rename and move to `grepr`. This test for `fmt.Stringer` but ignores
other text-encoding interfaces such as `gg.Appender` or `encoding.TextMarshaler`
because `gtest` produces the "simple" representation by calling `fmt.Sprint`,
which does not support any of those additional interfaces.
*/
func isSimple(src any) bool {
	return src == nil || (!gg.AnyIs[fmt.Stringer](src) &&
		!gg.AnyIs[fmt.GoStringer](src) &&
		isPrim(src))
}

// TODO should probably move to `gg` and make public.
func isPrim(src any) bool {
	val := r.ValueOf(src)

	switch val.Kind() {
	case r.Bool,
		r.Int8, r.Int16, r.Int32, r.Int64, r.Int,
		r.Uint8, r.Uint16, r.Uint32, r.Uint64, r.Uint, r.Uintptr,
		r.Float32, r.Float64,
		r.Complex64, r.Complex128,
		r.String:
		return true
	default:
		return false
	}
}
