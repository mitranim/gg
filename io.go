package gg

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	r "reflect"
	"strings"
)

/*
Creates a read-closer able to read from the given string or byte slice. Similar
to the following, but shorter and avoids allocation in case of bytes-to-string
or string-to-bytes conversion:

	io.NopCloser(strings.NewReader(string(val)))
	io.NopCloser(bytes.NewReader([]byte(val)))
*/
func NewReadCloser[A Text](val A) io.ReadCloser {
	if Kind[A]() == r.String {
		return new(StringReadCloser).Reset(ToString(val))
	}
	return new(BytesReadCloser).Reset(ToBytes(val))
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
Shortcut for using `os.Stat` to check if there is an existing file or directory
at the given path.
*/
func PathExists(path string) bool {
	info := fileInfo(path)
	return info != nil
}

/*
Shortcut for using `os.Stat` to check if there is an existing directory at the
given path.
*/
func DirExists(path string) bool {
	info := fileInfo(path)
	return info != nil && info.IsDir()
}

/*
Shortcut for using `os.Stat` to check if the file at the given path exists,
and is not a directory.
*/
func FileExists(path string) bool {
	info := fileInfo(path)
	return info != nil && !info.IsDir()
}

func fileInfo(path string) os.FileInfo {
	if path == `` {
		return nil
	}
	info, _ := os.Stat(path)
	return info
}

// Shortcut for `os.ReadDir`. Panics on error.
func ReadDir(path string) []fs.DirEntry { return Try1(os.ReadDir(path)) }

/*
Shortcut for `os.ReadFile`. Panics on error. Converts the content to the
requested text type without an additional allocation.
*/
func ReadFile[A Text](path string) A {
	return ToText[A](Try1(os.ReadFile(path)))
}

/*
Fully reads the given stream via `io.ReadAll` and returns two "forks". If
reading fails, panics. If the input is nil, both outputs are nil.
*/
func ForkReader(src io.Reader) (_, _ io.Reader) {
	if src == nil {
		return nil, nil
	}

	defer Detailf(`failed to read for forking`)
	text := ReadAll(src)
	return NewReadCloser(text), NewReadCloser(text)
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
	text := ReadCloseAll(src)
	return NewReadCloser(text), NewReadCloser(text)
}

// Shortcut for `os.Getwd` that panics on error.
func Cwd() string { return Try1(os.Getwd()) }

// If the given closer is non-nil, closes it. Panics on error.
func Close(val io.Closer) {
	if val != nil {
		Try(val.Close())
	}
}

// Shortcut for `os.MkdirAll` with `os.ModePerm`.
func MkdirAll(path string) {
	Try(os.MkdirAll(path, os.ModePerm))
}

// Shortcut for `os.Stat` that panics on error.
func Stat(path string) fs.FileInfo { return Try1(os.Stat(path)) }

/*
Shortcut for writing the given text to the given `io.Writer`.
Automatically converts text to bytes and panics on errors.
*/
func Write[Out io.Writer, Src Text](out Out, src Src) {
	Try1(out.Write(ToBytes(src)))
}

const (
	// Standard terminal escape sequence.
	TermEsc = "\x1b"

	/**
	Escape sequence recognized by many terminals. When printed, should cause
	the terminal to scroll down as much as needed to create an appearance of
	clearing the window. Scrolling up should reveal previous content.
	*/
	TermEscClearSoft = TermEsc + `c`

	/**
	Escape sequence recognized by many terminals. When printed, should cause
	the terminal to clear the scrollback buffer without clearing the currently
	visible content.
	*/
	TermEscClearScrollback = TermEsc + `[3J`

	/**
	Escape sequence recognized by many terminals. When printed, should cause
	the terminal to clear both the scrollback buffer and the currently visible
	content.
	*/
	TermEscClearHard = TermEscClearSoft + TermEscClearScrollback
)

/*
Prints `TermEscClearScrollback` to `os.Stdout`, causing the current TTY to clear
the scrollback buffer.
*/
func TermClearScrollback() {
	_, _ = io.WriteString(os.Stdout, TermEscClearScrollback)
}

/*
Prints `TermEscClearScrollback` to `os.Stdout`, causing the current TTY to push
existing content out of view.
*/
func TermClearSoft() {
	_, _ = io.WriteString(os.Stdout, TermEscClearSoft)
}

// Prints `TermEscClearScrollback` to `os.Stdout`, clearing the current TTY.
func TermClearHard() {
	_, _ = io.WriteString(os.Stdout, TermEscClearHard)
}
