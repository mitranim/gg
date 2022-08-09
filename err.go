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
Superior alternative to standard library errors. Supports stack traces and error
wrapping. Provides a convenient builder API.
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

// Implement `Appender`, appending the same representation as `.Error`.
func (self Err) Append(inout []byte) []byte {
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
	buf = errAppendInner(buf, self.Cause)
	return buf
}

/*
Returns a text representation of the full error message with the stack trace, if
any.
*/
func (self Err) Stack() string { return ToString(self.AppendStack(nil)) }

/*
Appends a text representation of the full error message with the stack trace, if
any. The representation is the same as in `.Stack`.
*/
func (self Err) AppendStack(inout []byte) []byte {
	buf := Buf(inout)
	cause := self.Cause
	causeTraced := IsErrTraced(cause)

	if self.Msg == `` {
		if cause == nil {
			return PtrGet(self.Trace).AppendIndent(buf, 0)
		}

		if !causeTraced {
			buf.AppendString(cause.Error())
			buf = errAppendTraceIndent(buf, PtrGet(self.Trace))
			return buf
		}

		buf.Fprintf(`%+v`, cause)
		return buf
	}

	if !causeTraced {
		buf.AppendString(self.Msg)
		buf = errAppendInner(buf, cause)
		buf = errAppendTraceIndent(buf, PtrGet(self.Trace))
		return buf
	}

	buf.AppendString(self.Msg)

	if PtrGet(self.Trace).HasLen() {
		buf = errAppendTraceIndent(buf, PtrGet(self.Trace))
		if cause != nil {
			buf.AppendNewline()
			buf.AppendNewline()
			buf.AppendString(`cause:`)
			buf.AppendNewline()
		}
	} else if cause != nil {
		buf.AppendString(`: `)
	}

	{
		val, _ := cause.(interface{ AppendStack([]byte) []byte })
		if val != nil {
			return val.AppendStack(buf)
		}
	}

	buf.Fprintf(`%+v`, cause)
	return buf
}

// Implement `fmt.Formatter`.
func (self Err) Format(out fmt.State, verb rune) {
	if out.Flag('+') {
		out.Write(self.AppendStack(nil))
		return
	}

	if out.Flag('#') {
		type Error Err
		fmt.Fprintf(out, `%#v`, Error(self))
		return
	}

	_, _ = io.WriteString(out, self.Error())
}

/*
Implement `StackTraced`, which allows to retrieve stack traces from nested
errors.
*/
func (self Err) StackTrace() []uintptr { return PtrGet(self.Trace).Prim() }

// Returns a modified version where `.Msg` is set to the input.
func (self Err) Msgd(val string) Err {
	self.Msg = val
	return self
}

// Returns a modified version where `.Msg` is set from `fmt.Sprintf`.
func (self Err) Msgf(pat string, val ...any) Err {
	self.Msg = fmt.Sprintf(pat, NoEscUnsafe(val)...)
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
func (self Err) Traced(skip int) Err {
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
func (self Err) TracedOpt(skip int) Err {
	if self.IsTraced() {
		return self
	}
	return self.Traced(skip + 1)
}

// True if either the error or its cause has a non-empty stack trace.
func (self Err) IsTraced() bool {
	return PtrGet(self.Trace).HasLen() || IsErrTraced(self.Cause)
}

/*
Combines multiple errors. Used by `Conc`. Avoid casting this to `error`. Instead
call the method `Errs.Err`, which will correctly return a nil interface when
all errors are nil.
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
unwrapping if possible. Otherwise returns nil.
*/
func (self Errs) Err() error {
	switch self.LenNonNil() {
	case 0:
		return nil

	case 1:
		return self.First()

	default:
		return self
	}
}

// Counts non-nil errors.
func (self Errs) LenNonNil() int { return Count(self, IsErrNonNil) }

// Counts nil errors.
func (self Errs) LenNil() int { return Count(self, IsErrNil) }

// True if there are any non-nil errors.
func (self Errs) HasLen() bool { return self.LenNonNil() > 0 }

// First non-nil error.
func (self Errs) First() error { return Find(self, IsErrNonNil) }

// Returns an error message. Same as `.Error`.
func (self Errs) String() string {
	switch self.LenNonNil() {
	case 0:
		return ``

	case 1:
		return self.First().Error()

	default:
		return ToString(self.Append(nil))
	}
}

/*
Appends a text representation of the error or errors. The text is the same as
returned by `.Error`.
*/
func (self Errs) Append(buf []byte) []byte {
	switch self.LenNonNil() {
	case 0:
		return buf

	case 1:
		buf := Buf(buf)
		buf.AppendError(self.First())
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

/*
Implementation of `error` that wraps an arbitrary value. Useful in panic
recovery. Used internally by `Rec`.
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

func IsErrNil(val error) bool { return val == nil }

func IsErrNonNil(val error) bool { return val != nil }

/*
True if the error has a stack trace. Relies on a hidden interface
implemented by `Err`.
*/
func IsErrTraced(val error) bool { return ErrTrace(val).HasLen() }

/*
Creates an error with a stack trace and a message formatted via `fmt.Sprintf`.
*/
func Errf(pat string, val ...any) Err {
	return Err{}.Msgf(pat, val...).Traced(1)
}

/*
Wraps the given error, prepending the given message and idempotently adding a
stack trace.
*/
func Wrapf(err error, pat string, val ...any) error {
	if err == nil {
		return nil
	}
	return Err{}.Caused(err).Msgf(pat, val...).TracedOpt(1)
}

/*
Returns the stack trace of the given error, unwrapping it as much as necessary.
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
If the input implements `error`, tries to find its stack trace via `ErrTrace`.
If no trace is found, generates a new trace, skipping the given amount of
frames. Suitable for use on `any` values returned by `recover`.
*/
func AnyTrace(val any, skip int) Trace {
	out := ErrTrace(AnyAs[error](val))
	if out != nil {
		return out
	}
	return CaptureTrace(skip + 1)
}

// Idempotently adds a stack trace, skipping the given number of frames.
func ErrTraced(err error, skip int) error {
	if err == nil {
		return nil
	}
	if IsErrTraced(err) {
		return err
	}
	return errTraced(err, skip+1)
}

// Outlined to avoid weird deoptimization of `ErrTraced`.
func errTraced(err error, skip int) Err {
	val, ok := err.(Err)
	if ok {
		return val.Traced(skip + 1)
	}
	return Err{}.Caused(err).Traced(skip + 1)
}

/*
Idempotently converts the input to an error. If the input is nil, the output is
nil. If the input implements `error`, it's returned as-is. If the input does
not implement `error`, it's converted to `ErrStr` or wrapped with `ErrAny`.
*/
func ToErrAny(val any) error {
	if val == nil {
		return nil
	}

	err, _ := val.(error)
	if err != nil {
		return err
	}

	str, ok := val.(string)
	if ok {
		return ErrStr(str)
	}

	return ErrAny{val}
}

/*
Converts an arbitrary value to an error. Idempotently adds a stack trace. If the
input is a non-nil non-error, it's wrapped into `ErrAny`.
*/
func ToErrTraced(val any, skip int) error {
	return ErrTraced(ToErrAny(val), skip+1)
}

// Same as `ToErrTraced(val, 1)`.
func AnyErrTraced(val any) error { return ToErrTraced(val, 2) }

// If the error is nil, returns ``. Otherwise uses `.Error`.
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
	return Wrapf(err, `unable to decode %q into type %v`, src, typ)
}

/*
Shortcut for flushing errors out of error containers such as `context.Context`
or `sql.Rows`. If the inner error is non-nil, panics, idempotently adding a
stack trace. Otherwise does nothing.
*/
func ErrOk[A Errer](val A) { Try(ErrTraced(val.Err(), 1)) }

/*
Similar to `errors.As`. Differences:

	* Instead of taking a pointer and returning a boolean, this returns the
	  unwrapped error value. On success, output is non-zero. On failure, output
	  is zero.

	* Automatically tries non-pointer and pointer versions of the given type. The
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
error that satisfies the predicate, or nil. If the error
*/
func ErrFind(err error, fun func(error) bool) error {
	for err != nil && fun != nil {
		impl, _ := err.(ErrFinder)
		if impl != nil {
			return impl.Find(fun)
		}

		if fun(err) {
			return err
		}

		next := errors.Unwrap(err)
		if next != nil && next != err {
			err = next
			continue
		}
		return nil
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
