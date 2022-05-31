/*
Missing feature of the standard library: terse, expressive test assertions.
*/
package gtest

import (
	e "errors"
	"fmt"
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
)

/*
Used internally by assertion utils. Error wrapper whose default stringing
uses "%+v" formatting on the inner error, causing it to be ALWAYS formatted
with a stack trace, which is useful when panics are not caught.
*/
type Err struct{ gg.Err }

/*
Implement `error` by using full formatting on the inner error: multiline with a
stack trace.
*/
func (self Err) Error() string { return self.String() }

/*
Implement `fmt.Stringer` by using full formatting on the inner error: multiline
with a stack trace.
*/
func (self Err) String() string {
	var buf gg.Buf
	buf = self.Err.AppendStack(buf)
	buf.AppendString(gg.Newline)
	return buf.String()
}

/*
Must be deferred. Usage:

	func TestSomething(t *testing.T) {
		// Catches panics and uses `t.Fatalf`.
		defer gtest.Catch(t)

		// Test assertion that panics.
		gtest.Eq(10, 20)
	}
*/
func Catch(t testing.TB) {
	t.Helper()
	val := gg.AnyErrTraced(recover())
	if val != nil {
		t.Fatalf(`%+v`, val)
	}
}

func ToErr(skip int, msg string) Err {
	return Err{gg.Err{}.Msgd(msg).Traced(skip + 1)}
}

/*
Asserts that the input is `true`, or fails the test, printing the optional
additional messages and the stack trace.
*/
func True(val bool, opt ...any) {
	if !val {
		panic(ToErr(1, msgOpt(opt, `expected true, got false`)))
	}
}

/*
Asserts that the input is `false`, or fails the test, printing the optional
additional messages and the stack trace.
*/
func False(val bool, opt ...any) {
	if val {
		panic(ToErr(1, msgOpt(opt, `expected false, got true`)))
	}
}

/*
Asserts that the inputs are equal via "==", or fails the test, printing the
optional additional messages and the stack trace.
*/
func Eq[A comparable](act, exp A, opt ...any) {
	if act != exp {
		panic(ToErr(1, msgOpt(opt, msgEq(act, exp))))
	}
}

/*
Asserts that the inputs are equal via "==", or fails the test, printing the
optional additional messages and the stack trace. Doesn't statically require
the inputs to be comparable, but may panic if they aren't.
*/
func EqAny[A any](act, exp A, opt ...any) {
	if any(act) != any(exp) {
		panic(ToErr(1, msgOpt(opt, msgEq(act, exp))))
	}
}

func msgEq(act, exp any) string {
	return gg.JoinLinesOpt(`unexpected difference`, msgEqInner(act, exp))
}

func msgEqInner(act, exp any) string {
	return gg.JoinLinesOpt(
		msgDet(`actual detailed:`, goStringIndent(act)),
		msgDet(`expected detailed:`, goStringIndent(exp)),
		msgDet(`actual simple:`, gg.StringAny(act)),
		msgDet(`expected simple:`, gg.StringAny(exp)),
	)
}

func goStringIndent[A any](val A) string { return grepr.StringIndent(val, 1) }

func msgDet(msg, det string) string {
	if det == `` {
		return ``
	}
	return gg.JoinLinesOpt(msg, gg.Indent+det)
}

/*
Asserts that the inputs are not equal via "!=", or fails the test, printing the
optional additional messages and the stack trace.
*/
func NotEq[A any](act, nom A, opt ...any) {
	if any(act) == any(nom) {
		panic(ToErr(1, msgOpt(opt, msgNotEq(act, nom))))
	}
}

func msgNotEq[A any](act, nom A) string {
	return gg.JoinLinesOpt(`unexpected equality`, msgNotEqInner(act, nom))
}

func msgNotEqInner[A any](act, nom A) string {
	return gg.JoinLinesOpt(
		msgDet(`actual detailed:`, goStringIndent(act)),
		msgDet(`nominal detailed:`, goStringIndent(nom)),
		msgDet(`actual simple:`, gg.StringAny(act)),
		msgDet(`nominal simple:`, gg.StringAny(nom)),
	)
}

/*
Asserts that the inputs are deeply equal, or fails the test, printing the
optional additional messages and the stack trace.
*/
func Equal[A any](act, exp A, opt ...any) {
	if !gg.Equal(act, exp) {
		panic(ToErr(1, msgOpt(opt, msgEq(act, exp))))
	}
}

/*
Asserts that the input slices have the same set of elements, or fails the test,
printing the optional additional messages and the stack trace.
*/
func EqualSet[A ~[]B, B comparable](act, exp A, opt ...any) {
	if !gg.Equal(gg.SetFrom(act), gg.SetFrom(exp)) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`unexpected difference in element sets`,
			msgEqInner(act, exp)),
		)))
	}
}

/*
Asserts that the inputs are not deeply equal, or fails the test, printing the
optional additional messages and the stack trace.
*/
func NotEqual[A any](act, nom A, opt ...any) {
	if gg.Equal(act, nom) {
		panic(ToErr(1, msgOpt(opt, msgNotEq(act, nom))))
	}
}

/*
Asserts that the given slice headers (not their elements) are equal via
`gg.SliceIs`. This means they have the same data pointer, length, capacity.
Does NOT compare individual elements, unlike `Equal`. Otherwise fails the test,
printing the optional additional messages and the stack trace.
*/
func SliceIs[A ~[]B, B any](act, exp A, opt ...any) {
	if !gg.SliceIs(act, exp) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`unexpected difference in slice headers`,
			msgDet(`actual header:`, goStringIndent(gg.CastUnsafe[gg.SliceHeader](act))),
			msgDet(`expected header:`, goStringIndent(gg.CastUnsafe[gg.SliceHeader](exp))),
		))))
	}
}

/*
Asserts that the given slice headers (not their elements) are distinct. This
means at least one of the following fields is different: data pointer, length,
capacity. Does NOT compare individual elements, unlike `NotEqual`. Otherwise
fails the test, printing the optional additional messages and the stack trace.
*/
func NotSliceIs[A ~[]B, B any](act, nom A, opt ...any) {
	if gg.SliceIs(act, nom) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`expected given slice headers to be distinct, but they were equal`,
			msgDet(`actual header:`, goStringIndent(gg.CastUnsafe[gg.SliceHeader](act))),
			msgDet(`nominal header:`, goStringIndent(gg.CastUnsafe[gg.SliceHeader](nom))),
		))))
	}
}

/*
Asserts that the input is zero via `gg.IsZero`, or fails the test, printing the
optional additional messages and the stack trace.
*/
func Zero[A any](val A, opt ...any) {
	if !gg.IsZero(val) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`unexpected non-zero value`,
			msgSingle(val),
		))))
	}
}

/*
Asserts that the input is zero via `gg.IsZero`, or fails the test, printing the
optional additional messages and the stack trace.
*/
func NotZero[A any](val A, opt ...any) {
	if gg.IsZero(val) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`unexpected zero value`,
			msgSingle(val),
		))))
	}
}

func msgSingle[A any](val A) string {
	return gg.JoinLinesOpt(
		msgDet(`detailed:`, goStringIndent(val)),
		msgDet(`simple:`, gg.StringAny(val)),
	)
}

func msgOpt(opt []any, src string) string {
	return gg.JoinLinesOpt(
		src,
		msgDet(`extra:`, gg.SpacedOpt(gg.Map(opt, gg.StringAny[any])...)),
	)
}

/*
Asserts that the given function panics AND that the resulting error satisfies
the given error-testing function. Otherwise fails the test, printing the
optional additional messages and the stack trace.
*/
func Panic(test func(error) bool, fun func(), opt ...any) {
	err := gg.Catch(fun)

	if err == nil {
		panic(ToErr(1, msgOpt(opt, msgPanicNone(fun, test))))
	}

	if !test(err) {
		panic(ToErr(1, msgOpt(opt, msgErrMismatch(fun, test, err))))
	}
}

/*
Asserts that the given function panics with an error whose message contains the
given substring, or fails the test, printing the optional additional messages
and the stack trace.
*/
func PanicStr(exp string, fun func(), opt ...any) {
	err := gg.Catch(fun)
	if err == nil {
		panic(ToErr(1, msgOpt(opt, msgPanicNone(fun, nil))))
	}

	msg := err.Error()
	if !strings.Contains(msg, exp) {
		panic(ToErr(1, msgOpt(opt, msgErrMsgMismatch(fun, exp, msg))))
	}
}

/*
Asserts that the given function panics and the panic result matches the given
error via `errors.Is`, or fails the test, printing the optional additional
messages and the stack trace.
*/
func PanicErrIs(exp error, fun func(), opt ...any) {
	if exp == nil {
		panic(ToErr(1, msgOpt(opt, `expected error must be non-nil`)))
	}

	err := gg.Catch(fun)
	if err == nil {
		panic(ToErr(1, msgOpt(opt, msgPanicNone(fun, nil))))
	}

	if !e.Is(err, exp) {
		panic(ToErr(1, msgOpt(opt, msgErrIsMismatch(err, exp))))
	}
}

/*
Asserts that the given function panics, or fails the test, printing the optional
additional messages and the stack trace.
*/
func PanicAny(fun func(), opt ...any) {
	err := gg.Catch(fun)

	if err == nil {
		panic(ToErr(1, msgOpt(opt, msgPanicNone(fun, nil))))
	}
}

func msgPanicNone(fun func(), test func(error) bool) string {
	return gg.JoinLinesOpt(`unexpected lack of panic`, msgErrFunTest(fun, test))
}

func msgErrFunTest(fun func(), test func(error) bool) string {
	return gg.JoinLinesOpt(msgFun(fun), msgErrTest(test))
}

func msgFun(val func()) string {
	if val == nil {
		return ``
	}
	return msgDet(`function:`, gg.FuncName(val))
}

func msgErrTest(val func(error) bool) string {
	if val == nil {
		return ``
	}
	return msgDet(`error test:`, gg.FuncName(val))
}

func msgErrMismatch(fun func(), test func(error) bool, err error) string {
	return gg.JoinLinesOpt(
		`unexpected error mismatch`,
		msgErrFunTest(fun, test),
		msgDet(`actual error:`, gg.StringAny(err)),
	)
}

func msgErrMsgMismatch(fun func(), exp, act string) string {
	return gg.JoinLinesOpt(
		`unexpected error message mismatch`,
		msgFun(fun),
		msgDet(`actual error message:`, act),
		msgDet(`expected error message substring:`, exp),
	)
}

func msgErrIsMismatch(err, exp error) string {
	return gg.JoinLinesOpt(
		`unexpected error mismatch`,
		msgDet(`actual error:`, gg.StringAny(err)),
		msgDet(`expected error via errors.Is:`, gg.StringAny(exp)),
	)
}

/*
Asserts that the given function doesn't panic, or fails the test, printing the
error's trace if possible, the optional additional messages, and the stack
trace.
*/
func NoPanic(fun func(), opt ...any) {
	err := gg.Catch(fun)
	if err != nil {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`unexpected panic`,
			msgFunErr(fun, err),
		))))
	}
}

/*
Asserts that the given error is non-nil AND satisfies the given error-testing
function. Otherwise fails the test, printing the optional additional messages
and the stack trace.
*/
func Error(test func(error) bool, err error, opt ...any) {
	if err == nil {
		panic(ToErr(1, msgOpt(opt, msgErrorNone(test))))
	}

	if !test(err) {
		panic(ToErr(1, msgOpt(opt, msgErrMismatch(nil, test, err))))
	}
}

/*
Asserts that the given error is non-nil and its message contains the given
substring, or fails the test, printing the optional additional messages and the
stack trace.
*/
func ErrorStr(exp string, err error, opt ...any) {
	if err == nil {
		panic(ToErr(1, msgOpt(opt, msgErrorNone(nil))))
	}

	msg := err.Error()

	if !strings.Contains(msg, exp) {
		panic(ToErr(1, msgOpt(opt, msgErrMsgMismatch(nil, exp, msg))))
	}
}

/*
Asserts that the given error is non-nil, or fails the test, printing the
optional additional messages and the stack trace.
*/
func ErrorAny(err error, opt ...any) {
	if err == nil {
		panic(ToErr(1, msgOpt(opt, msgErrorNone(nil))))
	}
}

func msgErrorNone(test func(error) bool) string {
	return gg.JoinLinesOpt(`unexpected lack of error`, msgErrFunTest(nil, test))
}

/*
Asserts that the given error is nil, or fails the test, printing the error's
trace if possible, the optional additional messages, and the stack trace.
*/
func NoError(err error, opt ...any) {
	if err != nil {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`unexpected error`,
			msgErr(err),
		))))
	}
}

func msgFunErr(fun func(), err error) string {
	return gg.JoinLinesOpt(msgFun(fun), msgErr(err))
}

func msgErr(err error) string {
	return gg.JoinLinesOpt(
		msgDet(`error trace:`, strings.TrimSpace(gg.ErrTrace(err).StringIndent(1))),
		msgDet(`error string:`, gg.StringAny(err)),
	)
}

// Shortcut for error testing.
type ErrMsgTest string

// Tests that the given error has the given message.
func (self ErrMsgTest) Is(err error) bool {
	return err != nil && strings.Contains(err.Error(), string(self))
}

/*
Asserts that the given slice contains the given value, or fails the test,
printing the optional additional messages and the stack trace.
*/
func Has[A ~[]B, B comparable](src A, val B, opt ...any) {
	if !gg.Has(src, val) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`missing element in slice`,
			msgSliceElem(src, val),
		))))
	}
}

/*
Asserts that the given slice does not contain the given value, or fails the
test, printing the optional additional messages and the stack trace.
*/
func NotHas[A ~[]B, B comparable](src A, val B, opt ...any) {
	if gg.Has(src, val) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`unexpected element in slice`,
			msgSliceElem(src, val),
		))))
	}
}

func msgSliceElem[A ~[]B, B comparable](src A, val B) string {
	return gg.JoinLinesOpt(
		msgDet(`slice detailed:`, goStringIndent(src)),
		msgDet(`element detailed:`, goStringIndent(val)),
		msgDet(`slice simple:`, gg.StringAny(src)),
		msgDet(`element simple:`, gg.StringAny(val)),
	)
}

/*
Asserts that the first slice contains all elements from the second slice. In
other words, asserts that the first slice is a strict superset of the second.
Otherwise fails the test, printing the optional additional messages and the
stack trace.
*/
func HasEvery[A ~[]B, B comparable](src, exp A, opt ...any) {
	missing := gg.Subtract(exp, src)

	if len(missing) > 0 {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`expected outer slice to contain all elements from inner slice`,
			msgDet(`outer detailed:`, goStringIndent(src)),
			msgDet(`inner detailed:`, goStringIndent(exp)),
			msgDet(`missing detailed:`, goStringIndent(missing)),
			msgDet(`outer simple:`, gg.StringAny(src)),
			msgDet(`inner simple:`, gg.StringAny(exp)),
			msgDet(`missing simple:`, gg.StringAny(missing)),
		))))
	}
}

/*
Asserts that the first slice contains some elements from the second slice. In
other words, asserts that the element sets have an intersection. Otherwise
fails the test, printing the optional additional messages and the stack trace.
*/
func HasSome[A ~[]B, B comparable](src, exp A, opt ...any) {
	if !gg.HasSome(src, exp) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`unexpected lack of shared elements in two slices`,
			msgDet(`left detailed:`, goStringIndent(src)),
			msgDet(`right detailed:`, goStringIndent(exp)),
			msgDet(`left simple:`, gg.StringAny(src)),
			msgDet(`right simple:`, gg.StringAny(exp)),
		))))
	}
}

/*
Asserts that the first slice does not contain any from the second slice. In
other words, asserts that the element sets are disjoint. Otherwise fails the
test, printing the optional additional messages and the stack trace.
*/
func HasNone[A ~[]B, B comparable](src, exp A, opt ...any) {
	inter := gg.Intersect(src, exp)

	if len(inter) > 0 {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`expected left slice to contain no elements from right slice`,
			msgDet(`left detailed:`, goStringIndent(src)),
			msgDet(`right detailed:`, goStringIndent(exp)),
			msgDet(`intersection detailed:`, goStringIndent(inter)),
			msgDet(`left simple:`, gg.StringAny(src)),
			msgDet(`right simple:`, gg.StringAny(exp)),
			msgDet(`intersection simple:`, gg.StringAny(inter)),
		))))
	}
}

/*
Asserts that the given slice contains no duplicates, or fails the test, printing
the optional additional messages and the stack trace.
*/
func Uniq[A ~[]B, B comparable](src A, opt ...any) {
	dup, ok := foundDup(src)
	if ok {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`unexpected duplicate in slice`,
			msgSingle(dup),
		))))
	}
}

func foundDup[A comparable](src []A) (A, bool) {
	for ind, val := range src {
		for _, more := range src[ind+1:] {
			if val == more {
				return val, true
			}
		}
	}
	return gg.Zero[A](), false
}

/*
Asserts that the given chunk of text contains the given substring, or fails the
test, printing the optional additional messages and the stack trace.
*/
func TextHas[A, B gg.Text](src A, exp B, opt ...any) {
	if !strings.Contains(gg.ToString(src), gg.ToString(exp)) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`text does not contain substring`,
			msgDet(`full text:`, goStringIndent(gg.ToString(src))),
			msgDet(`substring:`, goStringIndent(gg.ToString(exp))),
		))))
	}
}

/*
Asserts that the given slice is empty, or fails the test, printing the optional
additional messages and the stack trace.
*/
func Empty[A ~[]B, B any](src A, opt ...any) {
	if len(src) != 0 {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`unexpected non-empty slice`,
			msgDet(`detailed:`, goStringIndent(src)),
			msgDet(`simple:`, gg.StringAny(src)),
		))))
	}
}

/*
Asserts that the given slice is not empty, or fails the test, printing the
optional additional messages and the stack trace.
*/
func NotEmpty[A ~[]B, B any](src A, opt ...any) {
	if len(src) == 0 {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(`unexpected empty slice`, msgSingle(src)))))
	}
}

/*
Asserts that the given slice is not empty, or fails the test, printing the
optional additional messages and the stack trace.
*/
func MapNotEmpty[Src ~map[Key]Val, Key comparable, Val any](src Src, opt ...any) {
	if len(src) == 0 {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(`unexpected empty map`, msgSingle(src)))))
	}
}

/*
Asserts that the given slice has exactly the given length, or fails the test,
printing the optional additional messages and the stack trace.
*/
func Len[A ~[]B, B any](src A, exp int, opt ...any) {
	if len(src) != exp {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			fmt.Sprintf(`got slice length %v, expected %v`, len(src), exp),
			msgSingle(src),
		))))
	}
}

/*
Asserts that the given text has exactly the given length, or fails the test,
printing the optional additional messages and the stack trace.
*/
func TextLen[A gg.Text](src A, exp int, opt ...any) {
	if len(src) != exp {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			fmt.Sprintf(`got text length %v, expected %v`, len(src), exp),
			msgSingle(src),
		))))
	}
}

/*
Asserts that `.String` of the input matches the expected string, or fails the
test, printing the optional additional messages and the stack trace.
*/
func Str[A fmt.Stringer](src A, exp string, opt ...any) {
	Eq(src.String(), exp, opt...)
}

/*
Asserts `one < two`, or fails the test, printing the optional additional
messages and the stack trace. For non-primitives that implement `gg.Lesser`,
see `Less`. Also see `LessEqPrim`.
*/
func LessPrim[A gg.LesserPrim](one, two A, opt ...any) {
	if !(one < two) {
		panic(ToErr(1, msgOpt(opt, msgLess(one, two))))
	}
}

/*
Asserts `one < two`, or fails the test, printing the optional additional
messages and the stack trace. For primitives, see `LessPrim`.
*/
func Less[A gg.Lesser[A]](one, two A, opt ...any) {
	if !one.Less(two) {
		panic(ToErr(1, msgOpt(opt, msgLess(one, two))))
	}
}

func msgLess[A any](one, two A) string {
	return gg.JoinLinesOpt(`expected A < B`, msgAB(one, two))
}

/*
Asserts `one <= two`, or fails the test, printing the optional additional
messages and the stack trace. For non-primitives that implement `gg.Lesser`,
see `LessEq`. Also see `LessPrim`.
*/
func LessEqPrim[A gg.LesserPrim](one, two A, opt ...any) {
	if !(one <= two) {
		panic(ToErr(1, msgOpt(opt, msgLessEq(one, two))))
	}
}

/*
Asserts `one <= two`, or fails the test, printing the optional additional
messages and the stack trace. For primitives, see `LessEqPrim`. Also see
`Less`.
*/
func LessEq[A interface {
	gg.Lesser[A]
	comparable
}](one, two A, opt ...any) {
	if !(one == two || one.Less(two)) {
		panic(ToErr(1, msgOpt(opt, msgLessEq(one, two))))
	}
}

func msgLessEq[A any](one, two A) string {
	return gg.JoinLinesOpt(`expected A <= B`, msgAB(one, two))
}

func msgAB[A any](one, two A) string {
	return gg.JoinLinesOpt(
		msgDet(`A detailed:`, goStringIndent(one)),
		msgDet(`B detailed:`, goStringIndent(two)),
		msgDet(`A simple:`, gg.StringAny(one)),
		msgDet(`B simple:`, gg.StringAny(two)),
	)
}

/*
Asserts that the given number is > 0, or fails the test, printing the optional
additional messages and the stack trace.
*/
func Pos[A gg.Signed](src A, opt ...any) {
	if !gg.IsPos(src) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`expected > 0, got value out of range`,
			msgSingle(src),
		))))
	}
}

/*
Asserts that the given number is < 0, or fails the test, printing the optional
additional messages and the stack trace.
*/
func Neg[A gg.Signed](src A, opt ...any) {
	if !gg.IsNeg(src) {
		panic(ToErr(1, msgOpt(opt, gg.JoinLinesOpt(
			`expected < 0, got value out of range`,
			msgSingle(src),
		))))
	}
}
