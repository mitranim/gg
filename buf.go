package gg

import (
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

/*
Short for "buffer". Simpler, cleaner, more usable alternative to
`strings.Builder` and `bytes.Buffer`.
*/
type Buf []byte

var (
	_ = fmt.Stringer(Zero[Buf]())
	_ = Appender(Zero[Buf]())
	_ = io.Writer(Zero[*Buf]())
	_ = io.StringWriter(Zero[*Buf]())
)

/*
Free cast to a string. Mutation of the original buffer affects the resulting
string.
*/
func (self Buf) String() string { return ToString(self) }

/*
Implement `Appender`. Appends its own content to the given buffer.
If the given buffer has no capacity, returns itself.
*/
func (self Buf) Append(val []byte) []byte {
	if cap(val) == 0 {
		return self
	}
	return append(val, self...)
}

/*
Implement `io.StringWriter`, appending the input to the buffer.
The error is always nil and may be ignored.
*/
func (self *Buf) WriteString(val string) (int, error) {
	*self = append(*self, val...)
	return len(val), nil
}

/*
Implement `io.Writer`, appending the input to the buffer.
The error is always nil and may be ignored.
*/
func (self *Buf) Write(val []byte) (int, error) {
	*self = append(*self, val...)
	return len(val), nil
}

// Appends the given string. Mutates and returns the receiver.
func (self *Buf) AppendString(val string) { *self = append(*self, val...) }

// Appends the given string N times. Mutates and returns the receiver.
func (self *Buf) AppendStringN(val string, count int) {
	if len(val) > 0 {
		for count > 0 {
			count--
			self.AppendString(val)
		}
	}
}

// Appends `Indent`. Mutates and returns the receiver.
func (self *Buf) AppendIndent() { self.AppendString(Indent) }

// Appends `Indent` N times. Mutates and returns the receiver.
func (self *Buf) AppendIndents(lvl int) { self.AppendStringN(Indent, lvl) }

// Appends the given bytes. Mutates and returns the receiver.
func (self *Buf) AppendBytes(val []byte) { *self = append(*self, val...) }

// Appends the given byte. Mutates and returns the receiver.
func (self *Buf) AppendByte(val byte) { *self = append(*self, val) }

// Appends the given rune. Mutates and returns the receiver.
func (self *Buf) AppendRune(val rune) { *self = utf8.AppendRune(*self, val) }

// Appends a single space. Mutates and returns the receiver.
func (self *Buf) AppendSpace() { self.AppendByte(' ') }

// Appends a space N times. Mutates and returns the receiver.
func (self *Buf) AppendSpaces(count int) { self.AppendByteN(' ', count) }

// Appends the given byte N times. Mutates and returns the receiver.
func (self *Buf) AppendByteN(val byte, count int) {
	for range Iter(count) {
		self.AppendByte(val)
	}
}

// Appends `Newline`. Mutates and returns the receiver.
func (self *Buf) AppendNewline() { self.AppendString(Newline) }

// Appends `Newline` N times. Mutates and returns the receiver.
func (self *Buf) AppendNewlines(count int) { self.AppendStringN(Newline, count) }

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendUint(val int) {
	*self = strconv.AppendUint(*self, uint64(val), 10)
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendUint64(val uint64) {
	*self = strconv.AppendUint(*self, val, 10)
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendInt(val int) {
	*self = strconv.AppendInt(*self, int64(val), 10)
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendInt64(val int64) {
	*self = strconv.AppendInt(*self, val, 10)
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendFloat32(val float32) {
	*self = strconv.AppendFloat(*self, float64(val), 'f', -1, 32)
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendFloat64(val float64) {
	*self = strconv.AppendFloat(*self, val, 'f', -1, 64)
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendBool(val bool) { *self = strconv.AppendBool(*self, val) }

/*
Appends the string representation of the given error. If the input is nil, this
is a nop. Mutates and returns the receiver.
*/
func (self *Buf) AppendError(val error) {
	if val == nil {
		return
	}

	impl, _ := val.(Appender)
	if impl != nil {
		*self = impl.Append(*self)
		return
	}

	self.AppendString(val.Error())
}

/*
Appends the text representation of the input, using the `Append` function.
Mutates and returns the receiver.
*/
func (self *Buf) AppendAny(val any) { *self = Append(*self, val) }

/*
Appends the text representation of the input, using the `AppendGoString`
function. Mutates and returns the receiver.
*/
func (self *Buf) AppendGoString(val any) { *self = AppendGoString(*self, val) }

// Shortcut for appending a formatted string.
func (self *Buf) Fprintf(pat string, val ...any) {
	_, _ = fmt.Fprintf(self, pat, NoEscUnsafe(val)...)
}

// Shortcut for appending a formatted string with an idempotent trailing newline.
func (self *Buf) Fprintlnf(pat string, val ...any) {
	str := fmt.Sprintf(pat, NoEscUnsafe(val)...)
	self.AppendString(str)
	if !HasNewlineSuffix(str) {
		self.AppendNewline()
	}
}

// Same as `len(buf)`.
func (self Buf) Len() int { return len(self) }

// Replaces the buffer with the given slice.
func (self *Buf) Reset(src []byte) { *self = src }

/*
Increases the buffer's length by N zero values.
Mutates and returns the receiver.
*/
func (self *Buf) GrowLen(size int) { *self = GrowLen(*self, size) }

/*
Increases the buffer's capacity sufficiently to accommodate N additional
elements. Mutates and returns the receiver.
*/
func (self *Buf) GrowCap(size int) { *self = GrowCap(*self, size) }

/*
Truncates the buffer's length, preserving the capacity.
Does not modify the content. Mutates and returns the receiver.
*/
func (self *Buf) Clear() {
	if self != nil && *self != nil {
		*self = (*self)[:0]
	}
}
