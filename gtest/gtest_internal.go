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

func msgSingle[A any](val A) string {
	if isSimple(val) {
		return Msg(`value:`, gg.StringAny(val))
	}

	return gg.JoinLinesOpt(
		Msg(`detailed:`, goStringIndent(val)),
		Msg(`simple:`, gg.StringAny(val)),
	)
}

func msgOpt(opt []any, src string) string {
	return gg.JoinLinesOpt(
		src,
		MsgOpt(`extra:`, gg.SpacedOpt(opt...)),
	)
}

func msgPanicNoneWithTest(fun func(), test func(error) bool) string {
	return gg.JoinLinesOpt(msgNotPanic(), msgErrFunTest(fun, test))
}

func msgPanicNoneWithStr(fun func(), exp string) string {
	return gg.JoinLinesOpt(msgNotPanic(), msgFun(fun), msgExp(exp))
}

func msgPanicNoneWithErr(fun func(), exp error) string {
	return gg.JoinLinesOpt(msgNotPanic(), msgFun(fun), msgExp(exp))
}

func msgErrFunTest(fun func(), test func(error) bool) string {
	return gg.JoinLinesOpt(msgFun(fun), msgErrTest(test))
}

func msgNotPanic() string { return `unexpected lack of panic` }

func msgFun(val func()) string {
	if val == nil {
		return ``
	}
	return Msg(`function:`, gg.FuncName(val))
}

func msgErrTest(val func(error) bool) string {
	if val == nil {
		return ``
	}
	return Msg(`error test:`, gg.FuncName(val))
}

func msgErrMismatch(fun func(), test func(error) bool, err error) string {
	return gg.JoinLinesOpt(
		`unexpected error mismatch`,
		msgErrFunTest(fun, test),
		msgErrActual(err),
	)
}

func msgErrMsgMismatch(fun func(), exp, act string) string {
	return gg.JoinLinesOpt(
		`unexpected error message mismatch`,
		msgFun(fun),
		Msg(`actual error message:`, act),
		Msg(`expected error message substring:`, exp),
	)
}

func msgErrIsMismatch(err, exp error) string {
	return gg.JoinLinesOpt(
		`unexpected error mismatch`,
		msgErrActual(err),
		Msg(`expected error via errors.Is:`, gg.StringAny(exp)),
	)
}

func msgExp[A any](val A) string { return Msg(`expected:`, gg.StringAny(val)) }

func msgErrorNone(test func(error) bool) string {
	return gg.JoinLinesOpt(`unexpected lack of error`, msgErrTest(test))
}

func msgFunErr(fun func(), err error) string {
	return gg.JoinLinesOpt(msgFun(fun), msgErr(err))
}

func msgErr(err error) string {
	return gg.JoinLinesOpt(
		Msg(`error trace:`, errTrace(err)),
		Msg(`error string:`, gg.StringAny(err)),
	)
}

func msgErrActual(err error) string {
	return gg.JoinLinesOpt(
		Msg(`actual error trace:`, errTrace(err)),
		Msg(`actual error string:`, gg.StringAny(err)),
	)
}

func errTrace(err error) string {
	return strings.TrimSpace(gg.ErrTrace(err).StringIndent(1))
}

func msgSliceElemMissing[A ~[]B, B any](src A, val B) string {
	return gg.JoinLinesOpt(`missing element in slice`, msgSliceElem(src, val))
}

func msgSliceElemUnexpected[A ~[]B, B any](src A, val B) string {
	return gg.JoinLinesOpt(`unexpected element in slice`, msgSliceElem(src, val))
}

func msgSliceElem[A ~[]B, B any](src A, val B) string {
	// TODO avoid detailed view when it's unnecessary.
	return gg.JoinLinesOpt(
		Msg(`slice detailed:`, goStringIndent(src)),
		Msg(`element detailed:`, goStringIndent(val)),
		Msg(`slice simple:`, gg.StringAny(src)),
		Msg(`element simple:`, gg.StringAny(val)),
	)
}

func msgLess[A any](one, two A) string {
	return gg.JoinLinesOpt(`expected A < B`, msgAB(one, two))
}

func msgLessEq[A any](one, two A) string {
	return gg.JoinLinesOpt(`expected A <= B`, msgAB(one, two))
}

func msgAB[A any](one, two A) string {
	return gg.JoinLinesOpt(
		// TODO avoid detailed view when it's unnecessary.
		Msg(`A detailed:`, goStringIndent(one)),
		Msg(`B detailed:`, goStringIndent(two)),
		Msg(`A simple:`, gg.StringAny(one)),
		Msg(`B simple:`, gg.StringAny(two)),
	)
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
