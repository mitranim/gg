package gg

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
Creates a read-closer able to read from the given string or byte slice.
Equivalent to `io.NopCloser(strings.NewReader(string(val)))`
but marginally more efficient.
*/
func NewReadCloser[A Text](val A) *StringReadCloser {
	return new(StringReadCloser).Reset(ToString(val))
}

// Variant of `strings.Reader` that also implements nop `io.Closer`.
type StringReadCloser struct{ strings.Reader }

// Calls `(*strings.Reader).Reset`.
func (self *StringReadCloser) Reset(src string) *StringReadCloser {
	self.Reader.Reset(src)
	return self
}

// Implement `io.Closer`. This is a nop. The error is always nil.
func (*StringReadCloser) Close() error { return nil }

// Variant of `bytes.Reader` that also implements nop `io.Closer`.
type BytesReadCloser struct{ bytes.Reader }

// Calls `(*bytes.Reader).Reset`.
func (self *BytesReadCloser) Reset(src []byte) *BytesReadCloser {
	self.Reader.Reset(src)
	return self
}

// Implement `io.Closer`. This is a nop. The error is always nil.
func (*BytesReadCloser) Close() error { return nil }

/*
Same as `io.ReadAll` but with different error handling.
If reader is nil, returns nil. Panics on errors.
*/
func ReadAll(src io.Reader) []byte {
	if src == nil {
		return nil
	}
	return Try1(io.ReadAll(src))
}

/*
Variant of `ReadAll` that closes the provided reader when done.
If reader is nil, returns nil. Panics on errors.
*/
func ReadCloseAll(src io.ReadCloser) []byte {
	if src == nil {
		return nil
	}
	defer src.Close()
	return Try1(io.ReadAll(src))
}

/*
Shortcut for `os.ReadFile`. Panics on error. Converts the content to the
requested text type without an additional allocation.
*/
func ReadFile[A Text](path string) A {
	return CastUnsafe[A](Try1(os.ReadFile(path)))
}

/*
Fully reads the given stream via `io.ReadAll`, closing it at the end, and
returns two "forks". Used internally by `(*gh.Req).CloneBody` and
`(*gh.Res).CloneBody`. If reading fails, panics. If the input is nil, both
outputs are nil.
*/
func ForkReadCloser(src io.ReadCloser) (_, _ io.ReadCloser) {
	if src == nil {
		return nil, nil
	}

	defer Detailf(`failed to read for forking`)
	val := ReadCloseAll(src)

	return NewReadCloser(val), NewReadCloser(val)
}

// Shortcut for `os.Getwd` that panics on error.
func Cwd() string { return Try1(os.Getwd()) }

/*
"Optional" variant of `filepath.Rel`. If the given path can't be made relative,
it's returned as-is.
*/
func RelOpt(base, src string) string {
	out, err := filepath.Rel(base, src)
	if err == nil && len(out) < len(src) {
		return out
	}
	return src
}

// If the given closer is non-nil, closes it, ignoring the error.
func Close(val io.Closer) {
	if val != nil {
		_ = val.Close()
	}
}
