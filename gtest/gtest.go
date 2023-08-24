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
)

/*
Used internally by assertion utils. Error wrapper whose default stringing
uses "%+v" formatting on the inner error, causing it to be ALWAYS formatted
with a stack trace, which is useful when panics are not caught.
*/
type VerbErr struct{ gg.Err }

/*
Implement `error` by using full formatting on the inner error: multiline with a
stack trace.
*/
func (self VerbErr) Error() string { return self.String() }

/*
Implement `fmt.Stringer` by using full formatting on the inner error: multiline
with a stack trace.
*/
func (self VerbErr) String() string {
	var buf gg.Buf
	buf = self.Err.AppendStackTo(buf)
	buf.AppendString(gg.Newline)
	return buf.String()
}

/*
Shortcut for generating a test error (of type `VerbErr` provided by this
package) with the given message, skipping the given amount of stack frames.
*/
func ErrAt(skip int, msg ...any) VerbErr {
	return VerbErr{gg.Err{}.Msgv(msg...).TracedAt(skip + 1)}
}

/*
Shortcut for generating an error where the given messages are combined as
lines.
*/
func ErrLines(msg ...any) VerbErr {
	// Suboptimal but not anyone's bottleneck.
	return ErrAt(1, gg.JoinLinesOpt(gg.Map(msg, gg.String[any])...))
}

/*
Must be deferred. Usage:

	func TestSomething(t *testing.T) {
		// Catches panics and uses `t.Fatalf`.
		defer gtest.Catch(t)

		// Test assertion. Panics and gets caught above.
		gtest.Eq(10, 20)
	}
*/
func Catch(t testing.TB) {
	t.Helper()
	val := gg.AnyErrTracedAt(recover(), 1)
	if val != nil {
		t.Fatalf(`%+v`, val)
	}
}

/*
Asserts that the input is `true`, or fails the test, printing the optional
additional messages and the stack trace.
*/
func True(val bool, opt ...any) {
	if !val {
		panic(ErrAt(1, gg.JoinLinesOpt(`expected true, got false`, MsgExtra(opt...))))
	}
}

/*
Asserts that the input is `false`, or fails the test, printing the optional
additional messages and the stack trace.
*/
func False(val bool, opt ...any) {
	if val {
		panic(ErrAt(1, gg.JoinLinesOpt(`expected false, got true`, MsgExtra(opt...))))
	}
}

/*
Asserts that the inputs are byte-for-byte identical, via `gg.Is`. Otherwise
fails the test, printing the optional additional messages and the stack trace.
Intended for interface values, maps, chans, funcs. For slices, use `SliceIs`.
*/
func Is[A any](act, exp A, opt ...any) {
	if gg.Is(act, exp) {
		return
	}

	if gg.Equal(act, exp) {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`inputs are equal but not identical`,
			MsgEqInner(act, exp),
			MsgExtra(opt...),
		)))
	}

	panic(ErrAt(1, gg.JoinLinesOpt(MsgEq(act, exp), MsgExtra(opt...))))
}

/*
Asserts that the inputs are NOT byte-for-byte identical, via `gg.Is`. Otherwise
fails the test, printing the optional additional messages and the stack trace.
Intended for interface values, maps, chans, funcs. For slices, use `NotSliceIs`.
*/
func NotIs[A any](act, exp A, opt ...any) {
	if gg.Is(act, exp) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgNotEq(act), MsgExtra(opt...))))
	}
}

/*
Asserts that the inputs are equal via `==`, or fails the test, printing the
optional additional messages and the stack trace.
*/
func Eq[A comparable](act, exp A, opt ...any) {
	if act != exp {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgEq(act, exp), MsgExtra(opt...))))
	}
}

/*
Asserts that the inputs are equal via `==`, or fails the test, printing the
optional additional messages and the stack trace. Doesn't statically require
the inputs to be comparable, but may panic if they aren't.
*/
func AnyEq[A any](act, exp A, opt ...any) {
	if any(act) != any(exp) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgEq(act, exp), MsgExtra(opt...))))
	}
}

/*
Asserts that the inputs are equal via `gg.TextEq`, or fails the test, printing
the optional additional messages and the stack trace.
*/
func TextEq[A gg.Text](act, exp A, opt ...any) {
	if !gg.TextEq(act, exp) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgEq(act, exp), MsgExtra(opt...))))
	}
}

/*
Asserts that the inputs are not equal via `!=`, or fails the test, printing the
optional additional messages and the stack trace.
*/
func NotEq[A comparable](act, nom A, opt ...any) {
	if act == nom {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgNotEq(act), MsgExtra(opt...))))
	}
}

/*
Asserts that the inputs are not equal via `gg.TextEq`, or fails the test,
printing the optional additional messages and the stack trace.
*/
func NotTextEq[A gg.Text](act, nom A, opt ...any) {
	if gg.TextEq(act, nom) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgNotEq(act), MsgExtra(opt...))))
	}
}

/*
Asserts that the inputs are deeply equal, or fails the test, printing the
optional additional messages and the stack trace.
*/
func Equal[A any](act, exp A, opt ...any) {
	if !gg.Equal(act, exp) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgEq(act, exp), MsgExtra(opt...))))
	}
}

/*
Asserts that the input slices have the same set of elements, or fails the test,
printing the optional additional messages and the stack trace.
*/
func EqualSet[A ~[]B, B comparable](act, exp A, opt ...any) {
	missingAct := gg.Exclude(exp, act...)
	var msgMissingAct string

	missingExp := gg.Exclude(act, exp...)
	var msgMissingExp string

	if len(missingAct) > 0 {
		msgMissingAct = Msg(`missing from actual:`, goStringIndent(missingAct))
	}

	if len(missingExp) > 0 {
		msgMissingExp = Msg(`missing from expected:`, goStringIndent(missingExp))
	}

	if !gg.Equal(gg.SetFrom(act), gg.SetFrom(exp)) {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`unexpected difference in element sets`,
			Msg(`actual:`, goStringIndent(act)),
			Msg(`expected:`, goStringIndent(exp)),
			msgMissingAct,
			msgMissingExp,
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the inputs are not deeply equal, or fails the test, printing the
optional additional messages and the stack trace.
*/
func NotEqual[A any](act, nom A, opt ...any) {
	if gg.Equal(act, nom) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgNotEq(act), MsgExtra(opt...))))
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
		panic(ErrAt(1, gg.JoinLinesOpt(
			`expected given slice headers to be identical, but they were distinct`,
			Msg(`actual header:`, goStringIndent(gg.SliceHeaderOf(act))),
			Msg(`expected header:`, goStringIndent(gg.SliceHeaderOf(exp))),
			MsgExtra(opt...),
		)))
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
		panic(ErrAt(1, gg.JoinLinesOpt(
			`expected given slice headers to be distinct, but they were identical`,
			Msg(`actual header:`, goStringIndent(gg.SliceHeaderOf(act))),
			Msg(`nominal header:`, goStringIndent(gg.SliceHeaderOf(nom))),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the input is zero via `gg.IsZero`, or fails the test, printing the
optional additional messages and the stack trace.
*/
func Zero[A any](val A, opt ...any) {
	if !gg.IsZero(val) {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`unexpected non-zero value`,
			MsgSingle(val),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the input is zero via `gg.IsZero`, or fails the test, printing the
optional additional messages and the stack trace.
*/
func NotZero[A any](val A, opt ...any) {
	if gg.IsZero(val) {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`unexpected zero value`,
			MsgSingle(val),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the given function panics AND that the resulting error satisfies
the given error-testing function. Otherwise fails the test, printing the
optional additional messages and the stack trace.
*/
func Panic(test func(error) bool, fun func(), opt ...any) {
	err := gg.Catch(fun)

	if err == nil {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgPanicNoneWithTest(fun, test), MsgExtra(opt...))))
	}

	if !test(err) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgErrMismatch(fun, test, err), MsgExtra(opt...))))
	}
}

/*
Asserts that the given function panics with an error whose message contains the
given substring, or fails the test, printing the optional additional messages
and the stack trace.
*/
func PanicStr(exp string, fun func(), opt ...any) {
	if exp == `` {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`refusing to test for panic with an empty expected error message`,
			MsgFun(fun),
			MsgExtra(opt...),
		)))
	}

	err := gg.Catch(fun)
	if err == nil {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgPanicNoneWithStr(fun, exp), MsgExtra(opt...))))
	}

	msg := err.Error()
	if !strings.Contains(msg, exp) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgErrMsgMismatch(fun, exp, msg), MsgExtra(opt...))))
	}
}

/*
Asserts that the given function panics and the panic result matches the given
error via `errors.Is`, or fails the test, printing the optional additional
messages and the stack trace.
*/
func PanicErrIs(exp error, fun func(), opt ...any) {
	if exp == nil {
		panic(ErrAt(1, gg.JoinLinesOpt(`expected error must be non-nil`, MsgExtra(opt...))))
	}

	err := gg.Catch(fun)
	if err == nil {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgPanicNoneWithErr(fun, exp), MsgExtra(opt...))))
	}

	if !e.Is(err, exp) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgErrIsMismatch(err, exp), MsgExtra(opt...))))
	}
}

/*
Asserts that the given function panics, or fails the test, printing the optional
additional messages and the stack trace.
*/
func PanicAny(fun func(), opt ...any) {
	err := gg.Catch(fun)

	if err == nil {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgPanicNoneWithTest(fun, nil), MsgExtra(opt...))))
	}
}

/*
Asserts that the given function doesn't panic, or fails the test, printing the
error's trace if possible, the optional additional messages, and the stack
trace.
*/
func NotPanic(fun func(), opt ...any) {
	err := gg.Catch(fun)
	if err != nil {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`unexpected panic`,
			MsgFunErr(fun, err),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the given error is non-nil AND satisfies the given error-testing
function. Otherwise fails the test, printing the optional additional messages
and the stack trace.
*/
func Err(test func(error) bool, err error, opt ...any) {
	if err == nil {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgErrNone(test), MsgExtra(opt...))))
	}

	if !test(err) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgErrMismatch(nil, test, err), MsgExtra(opt...))))
	}
}

/*
Asserts that the given error is non-nil and its message contains the given
substring, or fails the test, printing the optional additional messages and the
stack trace.
*/
func ErrStr(exp string, err error, opt ...any) {
	if err == nil {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgErrNone(nil), MsgExtra(opt...))))
	}

	msg := err.Error()

	if !strings.Contains(msg, exp) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgErrMsgMismatch(nil, exp, msg), MsgExtra(opt...))))
	}
}

/*
Asserts that the given error is non-nil and matches the expected error via
`errors.Is`, or fails the test, printing the optional additional messages and
the stack trace.
*/
func ErrIs(act, exp error, opt ...any) {
	if exp == nil {
		panic(ErrAt(1, gg.JoinLinesOpt(`expected error must be non-nil`, MsgExtra(opt...))))
	}

	if !e.Is(act, exp) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgErrIsMismatch(act, exp), MsgExtra(opt...))))
	}
}

/*
Asserts that the given error is non-nil, or fails the test, printing the
optional additional messages and the stack trace.
*/
func ErrAny(err error, opt ...any) {
	if err == nil {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgErrNone(nil), MsgExtra(opt...))))
	}
}

/*
Asserts that the given error is nil, or fails the test, printing the error's
trace if possible, the optional additional messages, and the stack trace.
*/
func NoErr(err error, opt ...any) {
	if err != nil {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`unexpected error`,
			MsgErr(err),
			MsgExtra(opt...),
		)))
	}
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
		panic(ErrAt(1, gg.JoinLinesOpt(MsgSliceElemMissing(src, val), MsgExtra(opt...))))
	}
}

/*
Asserts that the given slice does not contain the given value, or fails the
test, printing the optional additional messages and the stack trace.
*/
func NotHas[A ~[]B, B comparable](src A, val B, opt ...any) {
	if gg.Has(src, val) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgSliceElemUnexpected(src, val), MsgExtra(opt...))))
	}
}

/*
Asserts that the given slice contains the given value, or fails the test,
printing the optional additional messages and the stack trace. Uses `gg.Equal`
to compare values. For values that implement `comparable`, use `Has` which is
simpler and faster.
*/
func HasEqual[A ~[]B, B any](src A, val B, opt ...any) {
	if !gg.HasEqual(src, val) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgSliceElemMissing(src, val), MsgExtra(opt...))))
	}
}

/*
Asserts that the given slice does not contain the given value, or fails the
test, printing the optional additional messages and the stack trace. Uses
`gg.Equal` to compare values. For values that implement `comparable`, use
`HasNot` which is simpler and faster.
*/
func NotHasEqual[A ~[]B, B any](src A, val B, opt ...any) {
	if gg.HasEqual(src, val) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgSliceElemUnexpected(src, val), MsgExtra(opt...))))
	}
}

/*
Asserts that the first slice contains all elements from the second slice. In
other words, asserts that the first slice is a strict superset of the second.
Otherwise fails the test, printing the optional additional messages and the
stack trace.
*/
func HasEvery[A ~[]B, B comparable](src, exp A, opt ...any) {
	missing := gg.Exclude(exp, src...)

	if len(missing) > 0 {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`expected outer slice to contain all elements from inner slice`,
			// TODO avoid detailed view when it's unnecessary.
			Msg(`outer detailed:`, goStringIndent(src)),
			Msg(`inner detailed:`, goStringIndent(exp)),
			Msg(`missing detailed:`, goStringIndent(missing)),
			Msg(`outer simple:`, gg.StringAny(src)),
			Msg(`inner simple:`, gg.StringAny(exp)),
			Msg(`missing simple:`, gg.StringAny(missing)),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the first slice contains some elements from the second slice. In
other words, asserts that the element sets have an intersection. Otherwise
fails the test, printing the optional additional messages and the stack trace.
*/
func HasSome[A ~[]B, B comparable](src, exp A, opt ...any) {
	if !gg.HasSome(src, exp) {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`unexpected lack of shared elements in two slices`,
			Msg(`left detailed:`, goStringIndent(src)),
			Msg(`right detailed:`, goStringIndent(exp)),
			Msg(`left simple:`, gg.StringAny(src)),
			Msg(`right simple:`, gg.StringAny(exp)),
			MsgExtra(opt...),
		)))
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
		panic(ErrAt(1, gg.JoinLinesOpt(
			`expected left slice to contain no elements from right slice`,
			Msg(`left detailed:`, goStringIndent(src)),
			Msg(`right detailed:`, goStringIndent(exp)),
			Msg(`intersection detailed:`, goStringIndent(inter)),
			Msg(`left simple:`, gg.StringAny(src)),
			Msg(`right simple:`, gg.StringAny(exp)),
			Msg(`intersection simple:`, gg.StringAny(inter)),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that every element of the given slice satisfies the given predicate
function, or fails the test, printing the optional additional messages and the
stack trace.
*/
func Every[A ~[]B, B any](src A, fun func(B) bool, opt ...any) {
	for ind, val := range src {
		if fun == nil || !fun(val) {
			panic(ErrAt(1, gg.JoinLinesOpt(
				gg.Str(
					`expected every element to satisfy predicate `, gg.FuncName(fun),
					`; element at index `, ind, ` did not satisfy`,
				),
				Msg(`slice detailed:`, goStringIndent(src)),
				Msg(`element detailed:`, goStringIndent(val)),
				Msg(`slice simple:`, gg.StringAny(src)),
				Msg(`element simple:`, gg.StringAny(val)),
				MsgExtra(opt...),
			)))
		}
	}
}

/*
Asserts that at least one element of the given slice satisfies the given
predicate function, or fails the test, printing the optional additional
messages and the stack trace.
*/
func Some[A ~[]B, B any](src A, fun func(B) bool, opt ...any) {
	if gg.Some(src, fun) {
		return
	}

	panic(ErrAt(1, gg.JoinLinesOpt(
		gg.Str(
			`expected at least one element to satisfy predicate `, gg.FuncName(fun),
			`; found no such elements`,
		),
		Msg(`slice detailed:`, goStringIndent(src)),
		Msg(`slice simple:`, gg.StringAny(src)),
		MsgExtra(opt...),
	)))
}

/*
Asserts that no elements of the given slice satisfy the given predicate
function, or fails the test, printing the optional additional messages and the
stack trace.
*/
func None[A ~[]B, B any](src A, fun func(B) bool, opt ...any) {
	for ind, val := range src {
		if fun == nil || fun(val) {
			panic(ErrAt(1, gg.JoinLinesOpt(
				gg.Str(
					`expected every element to fail predicate `, gg.FuncName(fun),
					`; element at index `, ind, ` did not fail`,
				),
				Msg(`slice detailed:`, goStringIndent(src)),
				Msg(`element detailed:`, goStringIndent(val)),
				Msg(`slice simple:`, gg.StringAny(src)),
				Msg(`element simple:`, gg.StringAny(val)),
				MsgExtra(opt...),
			)))
		}
	}
}

/*
Asserts that the given slice contains no duplicates, or fails the test, printing
the optional additional messages and the stack trace.
*/
func Uniq[A ~[]B, B comparable](src A, opt ...any) {
	dup, ind0, ind1, ok := foundDup(src)
	if ok {
		panic(ErrAt(1, gg.JoinLinesOpt(
			fmt.Sprintf(`unexpected duplicate at indexes %v and %v`, ind0, ind1),
			MsgSingle(dup),
			MsgExtra(opt...),
		)))
	}
}

func foundDup[A comparable](src []A) (_ A, _ int, _ int, _ bool) {
	size := len(src)
	if !(size > 0) {
		return
	}

	found := make(map[A]int, size)
	for ind1, val := range src {
		ind0, ok := found[val]
		if ok {
			return val, ind0, ind1, true
		}
		found[val] = ind1
	}
	return
}

/*
Asserts that the given chunk of text contains the given substring, or fails the
test, printing the optional additional messages and the stack trace.
*/
func TextHas[A, B gg.Text](src A, exp B, opt ...any) {
	if !strings.Contains(gg.ToString(src), gg.ToString(exp)) {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`text does not contain substring`,
			Msg(`full text:`, goStringIndent(gg.ToString(src))),
			Msg(`substring:`, goStringIndent(gg.ToString(exp))),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the given chunk of text does not contain the given substring, or
fails the test, printing the optional additional messages and the stack trace.
*/
func NotTextHas[A, B gg.Text](src A, exp B, opt ...any) {
	if strings.Contains(gg.ToString(src), gg.ToString(exp)) {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`text contains unexpected substring`,
			Msg(`full text:`, goStringIndent(gg.ToString(src))),
			Msg(`substring:`, goStringIndent(gg.ToString(exp))),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the given slice is empty, or fails the test, printing the optional
additional messages and the stack trace.
*/
func Empty[A ~[]B, B any](src A, opt ...any) {
	if len(src) != 0 {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`unexpected non-empty slice`,
			Msg(`detailed:`, goStringIndent(src)),
			Msg(`simple:`, gg.StringAny(src)),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the given slice is not empty, or fails the test, printing the
optional additional messages and the stack trace.
*/
func NotEmpty[A ~[]B, B any](src A, opt ...any) {
	if len(src) <= 0 {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`unexpected empty slice`,
			MsgSingle(src),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the given slice is not empty, or fails the test, printing the
optional additional messages and the stack trace.
*/
func MapNotEmpty[Src ~map[Key]Val, Key comparable, Val any](src Src, opt ...any) {
	if len(src) <= 0 {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`unexpected empty map`,
			MsgSingle(src),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the given slice has exactly the given length, or fails the test,
printing the optional additional messages and the stack trace.
*/
func Len[A ~[]B, B any](src A, exp int, opt ...any) {
	if len(src) != exp {
		panic(ErrAt(1, gg.JoinLinesOpt(
			fmt.Sprintf(`got slice length %v, expected %v`, len(src), exp),
			MsgSingle(src),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the given slice has exactly the given capacity, or fails the test,
printing the optional additional messages and the stack trace.
*/
func Cap[A ~[]B, B any](src A, exp int, opt ...any) {
	if cap(src) != exp {
		panic(ErrAt(1, gg.JoinLinesOpt(
			fmt.Sprintf(`got slice capacity %v, expected %v`, cap(src), exp),
			MsgSingle(src),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the given text has exactly the given length, or fails the test,
printing the optional additional messages and the stack trace.
*/
func TextLen[A gg.Text](src A, exp int, opt ...any) {
	if len(src) != exp {
		panic(ErrAt(1, gg.JoinLinesOpt(
			fmt.Sprintf(`got text length %v, expected %v`, len(src), exp),
			MsgSingle(src),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that `.String` of the input matches the expected string, or fails the
test, printing the optional additional messages and the stack trace.
*/
func Str[A any](src A, exp string, opt ...any) {
	Eq(gg.String(src), exp, opt...)
}

/*
Asserts `one < two`, or fails the test, printing the optional additional
messages and the stack trace. For non-primitives that implement `gg.Lesser`,
see `Less`. Also see `LessEqPrim`.
*/
func LessPrim[A gg.LesserPrim](one, two A, opt ...any) {
	if !(one < two) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgLess(one, two), MsgExtra(opt...))))
	}
}

/*
Asserts `one < two`, or fails the test, printing the optional additional
messages and the stack trace. For primitives, see `LessPrim`.
*/
func Less[A gg.Lesser[A]](one, two A, opt ...any) {
	if !one.Less(two) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgLess(one, two), MsgExtra(opt...))))
	}
}

/*
Asserts `one <= two`, or fails the test, printing the optional additional
messages and the stack trace. For non-primitives that implement `gg.Lesser`,
see `LessEq`. Also see `LessPrim`.
*/
func LessEqPrim[A gg.LesserPrim](one, two A, opt ...any) {
	if !(one <= two) {
		panic(ErrAt(1, gg.JoinLinesOpt(MsgLessEq(one, two), MsgExtra(opt...))))
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
		panic(ErrAt(1, gg.JoinLinesOpt(MsgLessEq(one, two), MsgExtra(opt...))))
	}
}

/*
Asserts that the given number is > 0, or fails the test, printing the optional
additional messages and the stack trace.
*/
func Pos[A gg.Signed](src A, opt ...any) {
	if !gg.IsPos(src) {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`expected > 0, got value out of range`,
			MsgSingle(src),
			MsgExtra(opt...),
		)))
	}
}

/*
Asserts that the given number is < 0, or fails the test, printing the optional
additional messages and the stack trace.
*/
func Neg[A gg.Signed](src A, opt ...any) {
	if !gg.IsNeg(src) {
		panic(ErrAt(1, gg.JoinLinesOpt(
			`expected < 0, got value out of range`,
			MsgSingle(src),
			MsgExtra(opt...),
		)))
	}
}
