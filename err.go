package gg

import (
	"errors"
	"fmt"
	"io"
	r "reflect"
)

const (
	ErrInvalidInput ErrStr = `invalid input`
	ErrNyi          ErrStr = `not yet implemented`
)

/*
More powerful alternative to standard library errors. Supports stack traces and
error wrapping. Provides a convenient builder API.
*/
type Err struct {
	Msg   string
	Cause error
	Trace *Trace // By pointer to allow `==` without panics.
}

// Implement `error`.
func (self Err) Error() string { return self.String() }

// Implement a hidden interface for compatibility with `"errors".Unwrap`.
func (self Err) Unwrap() error { return self.Cause }

// Implement a hidden interface for compatibility with `"errors".Is`.
func (self Err) Is(err error) bool {
	val, ok := err.(Err)
	if ok {
		return self.Msg == val.Msg && errors.Is(self.Cause, val.Cause)
	}
	return errors.Is(self.Cause, err)
}

/*
Implement `Errer`. If the receiver is a zero value, returns nil. Otherwise casts
the receiver to an error.
*/
func (self Err) Err() error {
	if IsZero(self) {
		return nil
	}
	return self
}

// Implement `fmt.Stringer`.
func (self Err) String() string {
	if self.Cause == nil {
		return self.Msg
	}
	if self.Msg == `` {
		return self.Cause.Error()
	}
	return AppenderString(self)
}

// Implement `AppenderTo`, appending the same representation as `.Error`.
func (self Err) AppendTo(inout []byte) []byte {
	buf := Buf(inout)

	if self.Cause == nil {
		buf.AppendString(self.Msg)
		return buf
	}

	if self.Msg == `` {
		buf.AppendError(self.Cause)
		return buf
	}

	buf.AppendString(self.Msg)
	return errAppendInner(buf, self.Cause)
}

/*
Returns a text representation of the full error message with the stack trace,
if any. The method's name is chosen for consistency with the getter
`Error.prototype.stack` in JS, which behaves exactly like this method.
*/
func (self Err) Stack() string { return ToString(self.AppendStackTo(nil)) }

/*
Appends a text representation of the full error message with the stack trace, if
any. The representation is the same as in `.Stack`.
*/
func (self Err) AppendStackTo(buf []byte) []byte {
	if self.Cause == nil {
		return self.appendStackWithoutCause(buf)
	}
	return self.appendStackWithCause(buf)
}

func (self Err) appendStackWithoutCause(buf Buf) Buf {
	trace := self.OwnTrace()
	if self.Msg == `` {
		return errAppendTraceIndent(buf, trace)
	}

	buf.AppendString(self.Msg)
	return errAppendTraceIndentWithNewline(buf, trace)
}

func (self Err) appendStackWithCause(buf Buf) Buf {
	if self.Msg == `` {
		return self.appendStackWithCauseWithoutMsg(buf)
	}
	return self.appendStackWithCauseWithMsg(buf)
}

func (self Err) appendStackWithCauseWithoutMsg(buf Buf) Buf {
	cause := self.Cause
	trace := self.OwnTrace()

	if trace.IsEmpty() {
		buf.AppendErrorStack(cause)
		return buf
	}

	if !IsErrTraced(cause) {
		// Outer doesn't have a message and inner doesn't have a trace, so we treat
		// them as complements to each other.
		buf.AppendError(cause)
		return errAppendTraceIndentWithNewline(buf, trace)
	}

	buf = errAppendTraceIndent(buf, trace)
	buf.AppendNewline()
	buf.AppendErrorStack(cause)
	return buf
}

func (self Err) appendStackWithCauseWithMsg(buf Buf) Buf {
	cause := self.Cause
	trace := self.OwnTrace()

	if trace.IsEmpty() {
		buf.AppendString(self.Msg)
		buf.AppendString(`: `)
		buf.AppendErrorStack(cause)
		return buf
	}

	if !IsErrTraced(cause) {
		buf.AppendString(self.Msg)
		buf = errAppendInner(buf, cause)
		return errAppendTraceIndentWithNewline(buf, trace)
	}

	buf.AppendString(self.Msg)
	buf = errAppendTraceIndentWithNewline(buf, trace)
	buf.AppendNewline()
	buf.AppendString(`cause: `)
	buf.AppendErrorStack(cause)
	return buf
}

// Implement `fmt.Formatter`.
func (self Err) Format(out fmt.State, verb rune) {
	if out.Flag('+') {
		_, _ = out.Write(self.AppendStackTo(nil))
		return
	}

	if out.Flag('#') {
		type Error Err
		_, _ = fmt.Fprintf(out, `%#v`, Error(self))
		return
	}

	_, _ = io.WriteString(out, self.Error())
}

// Safely dereferences `.Trace`, returning nil if the pointer is nil.
func (self Err) OwnTrace() Trace { return PtrGet(self.Trace) }

/*
Implement `StackTraced`, which allows to retrieve stack traces from nested
errors.
*/
func (self Err) StackTrace() []uintptr { return self.OwnTrace() }

// Returns a modified version where `.Msg` is set to the input.
func (self Err) Msgd(val string) Err {
	self.Msg = val
	return self
}

// Returns a modified version where `.Msg` is set from `fmt.Sprintf`.
func (self Err) Msgf(pat string, arg ...any) Err {
	self.Msg = fmt.Sprintf(pat, NoEscUnsafe(arg)...)
	return self
}

/*
Returns a modified version where `.Msg` is set to a concatenation of strings
generated from the arguments, via `Str`. See `StringCatch` for the encoding
rules.
*/
func (self Err) Msgv(src ...any) Err {
	self.Msg = Str(src...)
	return self
}

// Returns a modified version with the given `.Cause`.
func (self Err) Caused(val error) Err {
	self.Cause = val
	return self
}

/*
Returns a modified version where `.Trace` is initialized idempotently if
`.Trace` was nil. Skips the given amount of stack frames when capturing the
trace, where 1 corresponds to the caller's frame.
*/
func (self Err) TracedAt(skip int) Err {
	if self.Trace == nil {
		self.Trace = Ptr(CaptureTrace(skip + 1))
	}
	return self
}

/*
Returns a modified version where `.Trace` is initialized idempotently if neither
the error nor `.Cause` had a trace. Skips the given amount of stack frames when
capturing the trace, where 1 corresponds to the caller's frame.
*/
func (self Err) TracedOptAt(skip int) Err {
	if self.IsTraced() {
		return self
	}
	return self.TracedAt(skip + 1)
}

// True if either the error or its cause has a non-empty stack trace.
func (self Err) IsTraced() bool {
	return self.OwnTrace().IsNotEmpty() || IsErrTraced(self.Cause)
}

/*
Shortcut for combining multiple errors via `Errs.Err`. Does NOT generate a stack
trace or modify the errors in any way.
*/
func ErrMul(src ...error) error { return Errs(src).Err() }

/*
Combines multiple errors. Used by `Conc`. Caution: although this implements the
`error` interface, avoid converting this to `error`. Even when the slice is nil,
the resulting interface value would be non-nil, which is incorrect. Instead,
call the method `Errs.Err`, which will correctly return a nil interface value
when all errors are nil.
*/
type Errs []error

// Implement `error`.
func (self Errs) Error() string { return self.String() }

// Implement a hidden interface for compatibility with `"errors".Unwrap`.
func (self Errs) Unwrap() error { return self.First() }

// Implement a hidden interface for compatibility with `"errors".Is`.
func (self Errs) Is(err error) bool {
	return Some(self, func(val error) bool {
		return val != nil && errors.Is(val, err)
	})
}

// Implement a hidden interface for compatibility with `"errors".As`.
func (self Errs) As(out any) bool {
	return Some(self, func(val error) bool {
		return errors.As(val, out)
	})
}

/*
Returns the first error that satisfies the given test function, by calling
`ErrFind` on each element. Order is depth-first rather than breadth-first.
*/
func (self Errs) Find(fun func(error) bool) error {
	if fun != nil {
		for _, val := range self {
			val = ErrFind(val, fun)
			if val != nil {
				return val
			}
		}
	}
	return nil
}

/*
Shortcut for `.Find(fun) != nil`. Returns true if at least one error satisfies
the given predicate function, using `ErrFind` to unwrap.
*/
func (self Errs) Some(fun func(error) bool) bool { return self.Find(fun) != nil }

// If there are any non-nil errors, panics with a stack trace.
func (self Errs) Try() { Try(self.Err()) }

/*
Implement `Errer`. If there are any non-nil errors, returns a non-nil error,
unwrapping if possible. Otherwise returns nil. Does NOT generate a stack trace
or modify the errors in any way.
*/
func (self Errs) Err() error {
	err, count := self.find()

	switch count {
	case 0:
		return nil

	case 1:
		return err

	default:
		return self
	}
}

/*
Counts successes (represented with nil errors) and failures (represented with
non-nil errors). The order of return parameters matches the general convention
for returning values followed by errors, in that order.
*/
func (self Errs) Count() (nils, notNils int) {
	for _, val := range self {
		if val == nil {
			nils++
		} else {
			notNils++
		}
	}
	return
}

// True if there is at least one non-nil error.
func (self Errs) HasErr() bool {
	_, notNils := self.Count()
	return notNils > 0
}

// First non-nil error.
func (self Errs) First() error { return Find(self, IsErrNotNil) }

// Appends non-nil errors, ignoring nils.
func (self *Errs) Add(src ...error) {
	for _, src := range src {
		if src != nil {
			*self = append(*self, src)
		}
	}
}

// Returns an error message. Same as `.Error`.
func (self Errs) String() string {
	err, count := self.find()

	/**
	TODO also implement `fmt.Formatter` with support for %+v, show stacktraces for
	inner errors.
	*/
	switch count {
	case 0:
		return ``

	case 1:
		return err.Error()

	default:
		return ToString(self.AppendTo(nil))
	}
}

/*
Appends a text representation of the error or errors. The text is the same as
returned by `.Error`.
*/
func (self Errs) AppendTo(buf []byte) []byte {
	err, count := self.find()

	switch count {
	case 0:
		return buf

	case 1:
		buf := Buf(buf)
		buf.AppendError(err)
		return buf

	default:
		return self.append(buf)
	}
}

func (self Errs) append(buf Buf) Buf {
	buf.AppendString(`multiple errors`)

	for _, val := range self {
		if val == nil {
			continue
		}
		buf.AppendString(`; `)
		buf.AppendError(val)
	}

	return buf
}

// Implement `fmt.Formatter`.
func (self Errs) Format(out fmt.State, verb rune) {
	if out.Flag('+') {
		_, _ = out.Write(self.AppendStackTo(nil))
		return
	}

	if out.Flag('#') {
		type Errs []error
		_, _ = fmt.Fprintf(out, `%#v`, Errs(self))
		return
	}

	_, _ = io.WriteString(out, self.Error())
}

// Appends a text representation of the errors with stack traces, if any.
func (self Errs) AppendStackTo(buf []byte) []byte {
	err, count := self.find()

	switch count {
	case 0:
		return buf

	case 1:
		buf := Buf(buf)
		buf.AppendErrorStack(err)
		return buf

	default:
		buf := Buf(buf)
		buf.AppendString(`multiple errors:`)

		for _, val := range self {
			if val == nil {
				continue
			}
			buf.AppendNewlines(2)
			buf.AppendErrorStack(val)
		}

		return buf
	}
}

/*
Adds a stack trace to each non-nil error via `WrapTracedAt`. Useful when writing
tools that internally use goroutines.
*/
func (self Errs) WrapTracedAt(skip int) {
	for ind := range self {
		self[ind] = WrapTracedAt(self[ind], skip)
	}
}

func (self Errs) find() (err error, count int) {
	for _, val := range self {
		if val == nil {
			continue
		}
		if err == nil {
			err = val
		}
		count++
	}
	return
}

/*
Implementation of `error` that wraps an arbitrary value. Useful in panic
recovery. Used internally by `AnyErr` and some other error-related functions.
*/
type ErrAny struct{ Val any }

// Implement `error`.
func (self ErrAny) Error() string { return fmt.Sprint(self.Val) }

// Implement a hidden interface in "errors".
func (self ErrAny) Unwrap() error { return AnyAs[error](self.Val) }

/*
String typedef that implements `error`. Errors of this type can be defined as
constants.
*/
type ErrStr string

// Implement `error`.
func (self ErrStr) Error() string { return string(self) }

// Implement `fmt.Stringer`.
func (self ErrStr) String() string { return string(self) }

// Implement a hidden interface for compatibility with `"errors".Is`.
func (self ErrStr) Is(err error) bool {
	tar, ok := err.(ErrStr)
	return ok && self == tar
}

// Self-explanatory.
func IsErrNil(val error) bool { return val == nil }

// Self-explanatory.
func IsErrNotNil(val error) bool { return val != nil }

/*
True if the error has a stack trace. Shortcut for `ErrTrace(val).IsNotEmpty()`.
*/
func IsErrTraced(val error) bool { return ErrTrace(val).IsNotEmpty() }

/*
Creates an error where the message is generated by passing the arguments to
`fmt.Sprintf`, with a stack trace. Also see `Errv`.
*/
func Errf(pat string, arg ...any) Err { return Err{}.Msgf(pat, arg...).TracedAt(1) }

/*
Creates an error where the message is generated by passing the arguments to
`Str`, with a stack trace. Suffix "v" is short for "vector", alluding to how
all arguments are treated equally, as opposed to "f" ("format") where the first
argument is a formatting pattern.
*/
func Errv(val ...any) Err { return Err{}.Msgv(val...).TracedAt(1) }

/*
If the error is nil, returns nil. Otherwise wraps the error, prepending the
given message and idempotently adding a stack trace. The message is converted
to a string via `Str(msg...)`.
*/
func Wrap(err error, msg ...any) error {
	if err == nil {
		return nil
	}
	return Err{}.Caused(err).Msgv(msg...).TracedOptAt(1)
}

/*
If the error is nil, returns nil. Otherwise wraps the error, prepending the
given message and idempotently adding a stack trace. The pattern argument must
be a hardcoded pattern string compatible with `fmt.Sprintf` and other similar
functions. If the pattern argument is an expression rather than a hardcoded
string, use `Wrap` instead.
*/
func Wrapf(err error, pat string, arg ...any) error {
	if err == nil {
		return nil
	}
	return Err{}.Caused(err).Msgf(pat, arg...).TracedOptAt(1)
}

/*
If the error is nil, returns nil. Otherwise wraps the error into an instance of
`Err` without a message with a new stack trace, skipping the given number of
frames. Unlike other "traced" functions, this one is NOT idempotent: it doesn't
check if the existing errors already have traces, and always adds a trace. This
is useful when writing tools that internally use goroutines, in order to
"connect" the traces between the goroutine and the caller.
*/
func WrapTracedAt(err error, skip int) error {
	if err == nil {
		return nil
	}
	return Err{Cause: err, Trace: Ptr(CaptureTrace(skip + 1))}
}

/*
Idempotently converts the input to an error. If the input is nil, the output is
nil. If the input implements `error`, it's returned as-is. If the input does
not implement `error`, it's converted to `ErrStr` or wrapped with `ErrAny`.
Does NOT generate a stack trace or modify an underlying `error` in any way.
See `AnyErrTraced` for that.
*/
func AnyErr(val any) error {
	switch val := val.(type) {
	case nil:
		return nil
	case error:
		return val
	case string:
		return ErrStr(val)
	default:
		return ErrAny{val}
	}
}

// Same as `AnyTraceAt(val, 1)`.
func AnyTrace(val any) Trace {
	/**
	Note for attentive readers: 1 in the comment and 2 here is intentional.
	It's required for the equivalence between `AnyTraceAt(val, 1)` and
	`AnyTrace(val)` at the call site.
	*/
	return AnyTraceAt(val, 2)
}

/*
If the input implements `error`, tries to find its stack trace via `ErrTrace`.
If no trace is found, generates a new trace, skipping the given amount of
frames. Suitable for `any` values returned by `recover`. The given value is
used only as a possible trace carrier, and its other properties are ignored.
Also see `ErrTrace` which is similar but does not capture a new trace.
*/
func AnyTraceAt(val any, skip int) Trace {
	out := ErrTrace(AnyAs[error](val))
	if out != nil {
		return out
	}
	return CaptureTrace(skip + 1)
}

/*
Returns the stack trace of the given error, unwrapping it as much as necessary.
Uses the `StackTraced` interface to detect the trace; the interface is
implemented by the type `Err` provided by this library, and by trace-enabled
errors in "github.com/pkg/errors". Does NOT generate a new trace. Also see
`ErrStack` which returns a string that includes both the error message and the
trace's representation, and `AnyTraceAt` which is suitable for use with
`recover` and idempotently adds a trace if one is missing.
*/
func ErrTrace(val error) Trace {
	for val != nil {
		impl, _ := val.(StackTraced)
		if impl != nil {
			out := impl.StackTrace()
			if out != nil {
				return ToTrace(out)
			}
		}
		val = errors.Unwrap(val)
	}
	return nil
}

/*
Returns a string that includes both the message and the representation of the
trace of the given error, if possible. If the error is nil, the output is zero.
Does not capture a new trace. Also see `ErrTrace` which returns the `Trace` of
the given error, if possible. The name of this function is consistent with the
method `Err.Stack`.
*/
func ErrStack(err error) string { return Err{Cause: err}.Stack() }

// Same as `ErrTracedAt(val, 1)`.
func ErrTraced(err error) error {
	// See `AnyTrace` for notes on 1 vs 2.
	return ErrTracedAt(err, 2)
}

// Idempotently adds a stack trace, skipping the given number of frames.
func ErrTracedAt(err error, skip int) error {
	if err == nil {
		return nil
	}
	if IsErrTraced(err) {
		return err
	}
	return errTracedAt(err, skip+1)
}

// Outlined to avoid deoptimization of `ErrTracedAt` observed in benchmarks.
func errTracedAt(err error, skip int) error {
	if err == nil {
		return nil
	}
	val, ok := err.(Err)
	if ok {
		return val.TracedAt(skip + 1)
	}
	return Err{}.Caused(err).TracedAt(skip + 1)
}

// Same as `AnyErrTracedAt(val, 1)`.
func AnyErrTraced(val any) error {
	// See `AnyTrace` for notes on 1 vs 2.
	return AnyErrTracedAt(val, 2)
}

/*
Converts an arbitrary value to an error. Idempotently adds a stack trace.
If the input is a non-nil non-error, it's wrapped into `ErrAny`.
*/
func AnyErrTracedAt(val any, skip int) error {
	switch val := val.(type) {
	case nil:
		return nil
	case error:
		return ErrTracedAt(val, skip+1)
	case string:
		return Err{Msg: val}.TracedAt(skip + 1)
	default:
		return Err{Cause: ErrAny{val}}.TracedAt(skip + 1)
	}
}

/*
Similar to `AnyErrTracedAt`, but always returns a value of the concrete type
`Err`. If the input is nil, the output is zero. Otherwise the output is always
non-zero. The message is derived from the input. The stack trace is reused from
the input if possible, otherwise it's generated here, skipping the given amount
of stack frames.
*/
func AnyToErrTracedAt(val any, skip int) (_ Err) {
	switch val := val.(type) {
	case nil:
		return
	case Err:
		return val.TracedOptAt(skip + 1)
	case string:
		return Err{Msg: val}.TracedAt(skip + 1)
	case error:
		return Err{}.Caused(val).TracedOptAt(skip + 1)
	default:
		return Err{Cause: ErrAny{val}}.TracedAt(skip + 1)
	}
}

// If the error is nil, returns "". Otherwise uses `.Error`.
func ErrString(val error) string {
	if val != nil {
		return val.Error()
	}
	return ``
}

/*
Returns an error that describes a failure to convert the given input to the
given output type. Used internally in various conversions.
*/
func ErrConv(src any, typ r.Type) error {
	return Errf(
		`unable to convert value %v of type %v to type %v`,
		src, r.TypeOf(src), typ,
	)
}

/*
Returns an error that describes a failure to decode the given string into the
given output type. Used internally in various conversions.
*/
func ErrParse[A Text](err error, src A, typ r.Type) error {
	return Wrapf(err, `unable to decode %q into %v`, src, typ)
}

/*
Shortcut for flushing errors out of error containers such as `context.Context`
or `sql.Rows`. If the inner error is non-nil, panics, idempotently adding a
stack trace. Otherwise does nothing.
*/
func ErrOk[A Errer](val A) { TryErr(ErrTracedAt(val.Err(), 1)) }

/*
Safely compares two error values, avoiding panics due to `==` on incomparable
underlying types. Returns true if both errors are nil, or if the underlying
types are comparable and the errors are `==`, or if the errors are identical
via `Is`.
*/
func ErrEq(err0, err1 error) bool {
	if err0 == nil && err1 == nil {
		return true
	}
	if err0 == nil || err1 == nil {
		return false
	}
	if r.TypeOf(err0).Comparable() && r.TypeOf(err1).Comparable() {
		return err0 == err1
	}
	return Is(err0, err1)
}

/*
Similar to `errors.As`. Differences:

  - Instead of taking a pointer and returning a boolean, this returns the
    unwrapped error value. On success, output is non-zero. On failure, output
    is zero.
  - Automatically tries non-pointer and pointer versions of the given type. The
    caller should always specify a non-pointer type. This provides nil-safety
    for types that implement `error` on the pointer type. The caller doesn't
    have to remember whether to use pointer or non-pointer.
*/
func ErrAs[
	Tar any,
	Ptr interface {
		*Tar
		error
	},
](src error) Tar {
	var tar Tar
	if AnyIs[error](tar) && errors.As(src, &tar) {
		return tar
	}

	var ptr Ptr
	if errors.As(src, &ptr) {
		return PtrGet((*Tar)(ptr))
	}

	return Zero[Tar]()
}

/*
Somewhat analogous to `errors.Is` and `errors.As`, but instead of comparing an
error to another error value or checking its type, uses a predicate function.
Uses `errors.Unwrap` to traverse the error chain and returns the outermost
error that satisfies the predicate, or nil.
*/
func ErrFind(err error, fun func(error) bool) error {
	if fun == nil {
		return nil
	}

	for err != nil {
		impl, _ := err.(ErrFinder)
		if impl != nil {
			return impl.Find(fun)
		}

		if fun(err) {
			return err
		}

		next := errors.Unwrap(err)
		if ErrEq(next, err) {
			break
		}
		err = next
	}

	return nil
}

/*
Shortcut that returns true if `ErrFind` is non-nil for the given error and
predicate function.
*/
func ErrSome(err error, fun func(error) bool) bool {
	return ErrFind(err, fun) != nil
}
