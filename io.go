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
Creates a read-closer able to read from the given string or byte slice, where
the "close" operation does nothing. Similar to combining stdlib functions, but
shorter and avoids allocation in case of bytes-to-string or string-to-bytes
conversion:

	// Longer and marginally less efficient:
	io.NopCloser(bytes.NewReader([]byte(`some_data`)))
	io.NopCloser(strings.NewReader(string(`some_data`)))

	// Equivalent, shorter, marginally more efficient:
	gg.NewReadCloser([]byte(`some_data`))
	gg.NewReadCloser(string(`some_data`))
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
Shortcut for using `os.ReadDir` to return a list of file names in the given
directory. Panics on error.
*/
func ReadDirFileNames(path string) []string {
	return MapCompact(ReadDir(path), dirEntryToFileName)
}

/*
Shortcut for `os.ReadFile`. Panics on error. Converts the content to the
requested text type without an additional allocation.
*/
func ReadFile[A Text](path string) A {
	return ToText[A](Try1(os.ReadFile(path)))
}

/*
Shortcut for `os.WriteFile` with default permissions `os.ModePerm`. Panics on
error. Takes an arbitrary text type conforming to `Text` and converts it to
bytes without an additional allocation.
*/
func WriteFile[A Text](path string, body A) {
	Try(os.WriteFile(path, ToBytes(body), os.ModePerm))
}

/*
Fully reads the given stream via `io.ReadAll` and returns two "forks". If
reading fails, panics. If the input is nil, both outputs are nil.
*/
func ForkReader(src io.Reader) (_, _ io.Reader) {
	if src == nil {
		return nil, nil
	}

	defer Detail(`failed to read for forking`)
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

	defer Detail(`failed to read for forking`)
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

// Shortcut for `os.MkdirAll` with `os.ModePerm`. Panics on error.
func MkdirAll(path string) { Try(os.MkdirAll(path, os.ModePerm)) }

// Shortcut for `os.Stat` that panics on error.
func Stat(path string) fs.FileInfo { return Try1(os.Stat(path)) }

/*
Shortcut for writing the given text to the given `io.Writer`.
Automatically converts text to bytes and panics on errors.
*/
func Write[Out io.Writer, Src Text](out Out, src Src) {
	Try1(out.Write(ToBytes(src)))
}

// https://en.wikipedia.org/wiki/ANSI_escape_code
const (
	// Standard terminal escape sequence. Same as "\x1b" or "\033".
	TermEsc = string(rune(27))

	// Control Sequence Introducer. Used for other codes.
	TermEscCsi = TermEsc + `[`

	// Update cursor position to first row, first column.
	TermEscCup = TermEscCsi + `1;1H`

	// Supposed to clear the screen without clearing the scrollback, aka soft
	// clear. Seems insufficient on its own, at least in some terminals.
	TermEscErase2 = TermEscCsi + `2J`

	// Supposed to clear the screen and the scrollback, aka hard clear. Seems
	// insufficient on its own, at least in some terminals.
	TermEscErase3 = TermEscCsi + `3J`

	// Supposed to reset the terminal to initial state, aka super hard clear.
	// Seems insufficient on its own, at least in some terminals.
	TermEscReset = TermEsc + `c`

	// Clear screen without clearing scrollback.
	TermEscClearSoft = TermEscCup + TermEscErase2

	// Clear screen AND scrollback.
	TermEscClearHard = TermEscCup + TermEscReset + TermEscErase3
)

/*
Prints `TermEscClearSoft` to `os.Stdout`, causing the current TTY to clear the
screen but not the scrollback, pushing existing content out of view.
*/
func TermClearSoft() { _, _ = io.WriteString(os.Stdout, TermEscClearSoft) }

/*
Prints `TermEscClearHard` to `os.Stdout`, clearing the current TTY completely
(both screen and scrollback).
*/
func TermClearHard() { _, _ = io.WriteString(os.Stdout, TermEscClearHard) }
