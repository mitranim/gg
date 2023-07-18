package gtest

import "github.com/mitranim/gg"

// Internal shortcut for generating parts of an error message.
func Msg(msg, det string) string { return gg.JoinLinesOpt(msg, reindent(det)) }

// Internal shortcut for generating parts of an error message.
func MsgOpt(msg, det string) string {
	if det == `` {
		return ``
	}
	return Msg(msg, det)
}

// Internal shortcut for generating parts of an error message.
func MsgExtra(src ...any) string {
	return MsgOpt(`extra:`, gg.SpacedOpt(src...))
}

// Internal shortcut for generating parts of an error message.
func MsgExp[A any](val A) string { return Msg(`expected:`, gg.StringAny(val)) }

// Internal shortcut for generating parts of an error message.
func MsgSingle[A any](val A) string {
	if isSimple(val) {
		return Msg(`value:`, gg.StringAny(val))
	}

	return gg.JoinLinesOpt(
		Msg(`detailed:`, goStringIndent(val)),
		Msg(`simple:`, gg.StringAny(val)),
	)
}

// Internal shortcut for generating parts of an error message.
func MsgEq(act, exp any) string {
	return gg.JoinLinesOpt(`unexpected difference`, MsgEqInner(act, exp))
}

// Internal shortcut for generating parts of an error message.
func MsgEqInner(act, exp any) string {
	if isSimple(act) && isSimple(exp) {
		return gg.JoinLinesOpt(
			Msg(`actual:`, gg.StringAny(act)),
			Msg(`expected:`, gg.StringAny(exp)),
		)
	}

	return gg.JoinLinesOpt(
		MsgEqDetailed(act, exp),
		MsgEqSimple(act, exp),
	)
}

// Internal shortcut for generating parts of an error message.
func MsgEqDetailed(act, exp any) string {
	return gg.JoinLinesOpt(
		Msg(`actual detailed:`, goStringIndent(act)),
		Msg(`expected detailed:`, goStringIndent(exp)),
	)
}

// Internal shortcut for generating parts of an error message.
func MsgEqSimple(act, exp any) string {
	return gg.JoinLinesOpt(
		Msg(`actual simple:`, gg.StringAny(act)),
		Msg(`expected simple:`, gg.StringAny(exp)),
	)
}

// Internal shortcut for generating parts of an error message.
func MsgNotEq[A any](act A) string {
	return gg.JoinLinesOpt(`unexpected equality`, MsgSingle(act))
}

// Internal shortcut for generating parts of an error message.
func MsgErr(err error) string {
	return gg.JoinLinesOpt(
		Msg(`error trace:`, errTrace(err)),
		Msg(`error string:`, gg.StringAny(err)),
	)
}

// Internal shortcut for generating parts of an error message.
func MsgErrNone(test func(error) bool) string {
	return gg.JoinLinesOpt(`unexpected lack of error`, MsgErrTest(test))
}

// Internal shortcut for generating parts of an error message.
func MsgErrActual(err error) string {
	return gg.JoinLinesOpt(
		Msg(`actual error trace:`, errTrace(err)),
		Msg(`actual error string:`, gg.StringAny(err)),
	)
}

// Internal shortcut for generating parts of an error message.
func MsgErrMismatch(fun func(), test func(error) bool, err error) string {
	return gg.JoinLinesOpt(
		`unexpected error mismatch`,
		MsgErrFunTest(fun, test),
		MsgErrActual(err),
	)
}

// Internal shortcut for generating parts of an error message.
func MsgErrMsgMismatch(fun func(), exp, act string) string {
	return gg.JoinLinesOpt(
		`unexpected error message mismatch`,
		MsgFun(fun),
		Msg(`actual error message:`, act),
		Msg(`expected error message substring:`, exp),
	)
}

// Internal shortcut for generating parts of an error message.
func MsgErrIsMismatch(err, exp error) string {
	return gg.JoinLinesOpt(
		`unexpected error mismatch`,
		MsgErrActual(err),
		Msg(`expected error via errors.Is:`, gg.StringAny(exp)),
	)
}

// Internal shortcut for generating parts of an error message.
func MsgErrTest(val func(error) bool) string {
	if val == nil {
		return ``
	}
	return Msg(`error test:`, gg.FuncName(val))
}

// Internal shortcut for generating parts of an error message.
func MsgErrFunTest(fun func(), test func(error) bool) string {
	return gg.JoinLinesOpt(MsgFun(fun), MsgErrTest(test))
}

// Internal shortcut for generating parts of an error message.
func MsgFunErr(fun func(), err error) string {
	return gg.JoinLinesOpt(MsgFun(fun), MsgErr(err))
}

// Internal shortcut for generating parts of an error message.
func MsgFun(val func()) string {
	if val == nil {
		return ``
	}
	return Msg(`function:`, gg.FuncName(val))
}

// Internal shortcut for generating parts of an error message.
func MsgNotPanic() string { return `unexpected lack of panic` }

// Internal shortcut for generating parts of an error message.
func MsgPanicNoneWithTest(fun func(), test func(error) bool) string {
	return gg.JoinLinesOpt(MsgNotPanic(), MsgErrFunTest(fun, test))
}

// Internal shortcut for generating parts of an error message.
func MsgPanicNoneWithStr(fun func(), exp string) string {
	return gg.JoinLinesOpt(MsgNotPanic(), MsgFun(fun), MsgExp(exp))
}

// Internal shortcut for generating parts of an error message.
func MsgPanicNoneWithErr(fun func(), exp error) string {
	return gg.JoinLinesOpt(MsgNotPanic(), MsgFun(fun), MsgExp(exp))
}

// Internal shortcut for generating parts of an error message.
func MsgSliceElemMissing[A ~[]B, B any](src A, val B) string {
	return gg.JoinLinesOpt(`missing element in slice`, MsgSliceElem(src, val))
}

// Internal shortcut for generating parts of an error message.
func MsgSliceElemUnexpected[A ~[]B, B any](src A, val B) string {
	return gg.JoinLinesOpt(`unexpected element in slice`, MsgSliceElem(src, val))
}

// Internal shortcut for generating parts of an error message.
func MsgSliceElem[A ~[]B, B any](src A, val B) string {
	// TODO avoid detailed view when it's unnecessary.
	return gg.JoinLinesOpt(
		Msg(`slice detailed:`, goStringIndent(src)),
		Msg(`element detailed:`, goStringIndent(val)),
		Msg(`slice simple:`, gg.StringAny(src)),
		Msg(`element simple:`, gg.StringAny(val)),
	)
}

// Internal shortcut for generating parts of an error message.
func MsgLess[A any](one, two A) string {
	return gg.JoinLinesOpt(`expected A < B`, MsgAB(one, two))
}

// Internal shortcut for generating parts of an error message.
func MsgLessEq[A any](one, two A) string {
	return gg.JoinLinesOpt(`expected A <= B`, MsgAB(one, two))
}

// Internal shortcut for generating parts of an error message.
func MsgAB[A any](one, two A) string {
	return gg.JoinLinesOpt(
		// TODO avoid detailed view when it's unnecessary.
		Msg(`A detailed:`, goStringIndent(one)),
		Msg(`B detailed:`, goStringIndent(two)),
		Msg(`A simple:`, gg.StringAny(one)),
		Msg(`B simple:`, gg.StringAny(two)),
	)
}
