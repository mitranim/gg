package gg

/*
If the error is nil, returns void. If the error is non-nil, panics with that
error, idempotently adding a stack trace.
*/
func Try(err error) {
	if err != nil {
		panic(ErrTracedAt(err, 1))
	}
}

/*
If the error is nil, returns void. If the error is non-nil, panics with that
error, idempotently adding a stack trace and skipping the given number of stack
frames.
*/
func TryAt(err error, skip int) {
	if err != nil {
		panic(ErrTracedAt(err, skip+1))
	}
}

/*
If the error is nil, returns void. If the error is non-nil, panics with that
exact error. Unlike `Try` and other similar functions, this function does NOT
generate a stack trace or modify the error in any way. Most of the time, you
should prefer `Try` and other similar functions. This function is intended
for "internal" library code, for example to avoid the stack trace check when
the trace is guaranteed, or to avoid generating the trace when the error is
meant to be caught before being returned to the caller.
*/
func TryErr(err error) {
	if err != nil {
		panic(err)
	}
}

/*
Like `Try`, but for multiple errors. Uses `Errs.Err` to combine the errors.
If the resulting error is nil, returns void. If the resulting error is non-nil,
panics with that error, idempotently adding a stack trace.
*/
func TryMul(src ...error) {
	err := ErrMul(src...)
	if err != nil {
		panic(ErrTracedAt(err, 1))
	}
}

/*
If the error is nil, returns the given value. If the error is non-nil, panics
with that error, idempotently adding a stack trace.
*/
func Try1[A any](val A, err error) A {
	if err != nil {
		panic(ErrTracedAt(err, 1))
	}
	return val
}

/*
If the error is nil, returns the given values. If the error is non-nil, panics
with that error, idempotently adding a stack trace.
*/
func Try2[A, B any](one A, two B, err error) (A, B) {
	if err != nil {
		panic(ErrTracedAt(err, 1))
	}
	return one, two
}

/*
If the error is nil, returns the given values. If the error is non-nil, panics
with that error, idempotently adding a stack trace.
*/
func Try3[A, B, C any](one A, two B, three C, err error) (A, B, C) {
	if err != nil {
		panic(ErrTracedAt(err, 1))
	}
	return one, two, three
}

/*
Must be deferred. Recovers from panics, writing the resulting error, if any, to
the given pointer. Should be used together with "try"-style functions.
Idempotently adds a stack trace.
*/
func Rec(out *error) {
	if out == nil {
		return
	}

	err := AnyErrTracedAt(recover(), 1)
	if err != nil {
		*out = err
	}
}

/*
Must be deferred. Recovers from panics, writing the resulting error, if any, to
the given pointer. Should be used together with "try"-style functions. Unlike
`Rec` and other similar functions, this function does NOT generate a stack
trace or modify the error in any way. Most of the time, you should prefer `Rec`
and other similar functions. This function is intended for "internal" library
code, for example to avoid the stack trace check when the trace is guaranteed,
or to avoid generating the trace when the error is meant to be caught before
being returned to the caller.
*/
func RecErr(out *error) {
	if out == nil {
		return
	}

	err := AnyErr(recover())
	if err != nil {
		*out = err
	}
}

/*
Must be deferred. Same as `Rec`, but skips the given amount of stack frames when
capturing a trace.
*/
func RecN(out *error, skip int) {
	if out == nil {
		return
	}

	err := AnyErrTracedAt(recover(), skip+1)
	if err != nil {
		*out = err
	}
}

/*
Must be deferred. Filtered version of `Rec`. Recovers from panics that
satisfy the provided test. Re-panics on non-nil errors that don't satisfy the
test. Does NOT check errors that are returned normally, without a panic.
Idempotently adds a stack trace.
*/
func RecOnly(ptr *error, test func(error) bool) {
	err := AnyErrTraced(recover())
	if err == nil {
		return
	}

	*ptr = err
	if test != nil && test(err) {
		return
	}

	panic(err)
}

/*
Must be deferred. Recovery for background goroutines which are not allowed to
crash. Calls the provided function ONLY if the error is non-nil.
*/
func RecWith(fun func(error)) {
	err := AnyErrTraced(recover())
	if err != nil && fun != nil {
		fun(err)
	}
}

/*
Runs the given function, converting a panic to an error.
Idempotently adds a stack trace.
*/
func Catch(fun func()) (err error) {
	defer RecN(&err, 1)
	if fun != nil {
		fun()
	}
	return
}

/*
Runs the given function with the given input, converting a panic to an error.
Idempotently adds a stack trace. Compare `Catch01` and `Catch11`.
*/
func Catch10[A any](fun func(A), val A) (err error) {
	defer RecN(&err, 1)
	if fun != nil {
		fun(val)
	}
	return
}

/*
Runs the given function, returning the function's result along with its panic
converted to an error. Idempotently adds a stack trace. Compare `Catch10` and
`Catch11`.
*/
func Catch01[A any](fun func() A) (val A, err error) {
	defer RecN(&err, 1)
	if fun != nil {
		val = fun()
	}
	return
}

/*
Runs the given function with the given input, returning the function's result
along with its panic converted to an error. Idempotently adds a stack trace.
Compare `Catch10` and `Catch01`.
*/
func Catch11[A, B any](fun func(A) B, val0 A) (val1 B, err error) {
	defer RecN(&err, 1)
	if fun != nil {
		val1 = fun(val0)
	}
	return
}

/*
Runs a given function, converting a panic to an error IF the error satisfies
the provided test. Idempotently adds a stack trace.
*/
func CatchOnly(test func(error) bool, fun func()) (err error) {
	defer RecOnly(&err, test)
	if fun != nil {
		fun()
	}
	return
}

/*
Shortcut for `Catch() != nil`. Useful when you want to handle all errors while
ignoring their content.
*/
func Caught(fun func()) bool {
	return Catch(fun) != nil
}

/*
Shortcut for `CatchOnly() != nil`. Useful when you want to handle a specific
error while ignoring its content.
*/
func CaughtOnly(test func(error) bool, fun func()) bool {
	return CatchOnly(test, fun) != nil
}

/*
Must be deferred. Catches panics; ignores errors that satisfy the provided
test; re-panics on other non-nil errors. Idempotently adds a stack trace.
*/
func SkipOnly(test func(error) bool) {
	err := AnyErrTraced(recover())
	if err != nil && test != nil && test(err) {
		return
	}
	Try(err)
}

// Runs a function, catching and ignoring ALL panics.
func Skipping(fun func()) {
	defer Skip()
	if fun != nil {
		fun()
	}
}

/*
Runs a function, catching and ignoring only the panics that satisfy the provided
test. Idempotently adds a stack trace.
*/
func SkippingOnly(test func(error) bool, fun func()) {
	defer SkipOnly(test)
	if fun != nil {
		fun()
	}
}

// Must be deferred. Catches and ignores ALL panics.
func Skip() { _ = recover() }

/*
Must be deferred. Tool for adding a stack trace to an arbitrary panic. Unlike
the "rec" functions, this does NOT prevent the panic from propagating. It
simply ensures that there's a stack trace, then re-panics.

Caution: due to idiosyncrasies of `recover()`, this works ONLY when deferred
directly. Anything other than `defer gg.Traced()` will NOT work.
*/
func Traced() { TryErr(AnyErrTraced(recover())) }

/*
Must be deferred. Version of `Traced` that skips the given number of stack
frames when generating a stack trace.
*/
func TracedAt(skip int) { Try(AnyErrTracedAt(recover(), skip+1)) }

/*
Must be deferred. Runs the function only if there's no panic. Idempotently adds
a stack trace.
*/
func Ok(fun func()) {
	Try(AnyErrTraced(recover()))
	if fun != nil {
		fun()
	}
}

/*
Must be deferred. Runs the function ONLY if there's an ongoing panic, and then
re-panics. Idempotently adds a stack trace.
*/
func Fail(fun func(error)) {
	err := AnyErrTraced(recover())
	if err != nil && fun != nil {
		fun(err)
	}
	Try(err)
}

/*
Must be deferred. Always runs the given function, passing either the current
panic or nil. If the error is non-nil, re-panics.
*/
func Finally(fun func(error)) {
	err := AnyErrTraced(recover())
	if fun != nil {
		fun(err)
	}
	Try(err)
}

/*
Must be deferred. Short for "transmute" or "transform". Catches an ongoing
panic, transforms the error by calling the provided function, and then
re-panics via `Try`. Idempotently adds a stack trace.
*/
func Trans(fun func(error) error) {
	err := AnyErrTraced(recover())
	if err != nil && fun != nil {
		err = fun(err)
	}
	Try(err)
}

/*
Runs a function, "transmuting" or "transforming" the resulting panic by calling
the provided transformer. See `Trans`.
*/
func Transing(trans func(error) error, fun func()) {
	defer Trans(trans)
	if fun != nil {
		fun()
	}
}

/*
Must be deferred. Similar to `Trans`, but transforms only non-nil errors that
satisfy the given predicate. Idempotently adds a stack trace.
*/
func TransOnly(test func(error) bool, trans func(error) error) {
	err := AnyErrTraced(recover())
	if err != nil && test != nil && trans != nil && test(err) {
		err = trans(err)
	}
	Try(err)
}

/*
Must be deferred. Wraps non-nil panics, prepending the error message and
idempotently adding a stack trace. The message is converted to a string via
`Str(msg...)`. Usage:

	defer gg.Detail(`unable to do X`)
	defer gg.Detail(`unable to do A with B `, someEntity.Id)
*/
func Detail(msg ...any) {
	Try(Wrap(AnyErr(recover()), msg...))
}

/*
Must be deferred. Wraps non-nil panics, prepending the error message and
idempotently adding a stack trace. Usage:

	defer gg.Detailf(`unable to %v`, `do X`)

The first argument must be a hardcoded pattern string compatible with
`fmt.Sprintf` and other similar functions. If the first argument is an
expression rather than a hardcoded string, use `Detail` instead.
*/
func Detailf(pat string, val ...any) {
	Try(Wrapf(AnyErr(recover()), pat, val...))
}

/*
Must be deferred. Wraps non-nil panics, prepending the error message, ONLY if
they satisfy the provided test. Idempotently adds a stack trace.
*/
func DetailOnlyf(test func(error) bool, pat string, val ...any) {
	err := AnyErrTraced(recover())
	if err != nil && test != nil && test(err) {
		err = Wrapf(err, pat, val...)
	}
	Try(err)
}
