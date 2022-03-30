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
func (self *Buf) AppendString(val string) *Buf {
	*self = append(*self, val...)
	return self
}

// Appends the given string N times. Mutates and returns the receiver.
func (self *Buf) AppendStringN(val string, count int) *Buf {
	if len(val) > 0 {
		for count > 0 {
			count--
			self.AppendString(val)
		}
	}
	return self
}

// Appends `Indent`. Mutates and returns the receiver.
func (self *Buf) AppendIndent() *Buf {
	return self.AppendString(Indent)
}

// Appends `Indent` N times. Mutates and returns the receiver.
func (self *Buf) AppendIndents(lvl int) *Buf {
	return self.AppendStringN(Indent, lvl)
}

// Appends the given bytes. Mutates and returns the receiver.
func (self *Buf) AppendBytes(val []byte) *Buf {
	*self = append(*self, val...)
	return self
}

// Appends the given byte. Mutates and returns the receiver.
func (self *Buf) AppendByte(val byte) *Buf {
	*self = append(*self, val)
	return self
}

// Appends the given rune. Mutates and returns the receiver.
func (self *Buf) AppendRune(val rune) *Buf {
	*self = utf8.AppendRune(*self, val)
	return self
}

// Appends a single space. Mutates and returns the receiver.
func (self *Buf) AppendSpace() *Buf { return self.AppendByte(' ') }

// Appends a space N times. Mutates and returns the receiver.
func (self *Buf) AppendSpaces(count int) *Buf {
	return self.AppendByteN(' ', count)
}

// Appends the given byte N times. Mutates and returns the receiver.
func (self *Buf) AppendByteN(val byte, count int) *Buf {
	for count > 0 {
		count--
		self.AppendByte(val)
	}
	return self
}

// Appends `Newline`. Mutates and returns the receiver.
func (self *Buf) AppendNewline() *Buf {
	return self.AppendString(Newline)
}

// Appends `Newline` N times. Mutates and returns the receiver.
func (self *Buf) AppendNewlines(count int) *Buf {
	return self.AppendStringN(Newline, count)
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendInt(val int) *Buf {
	*self = strconv.AppendInt(*self, int64(val), 10)
	return self
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendInt64(val int64) *Buf {
	*self = strconv.AppendInt(*self, val, 10)
	return self
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendFloat32(val float32) *Buf {
	*self = strconv.AppendFloat(*self, float64(val), 'f', -1, 32)
	return self
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendFloat64(val float64) *Buf {
	*self = strconv.AppendFloat(*self, val, 'f', -1, 64)
	return self
}

/*
Appends text representation of the input, using "strconv". Mutates and returns
the receiver.
*/
func (self *Buf) AppendBool(val bool) *Buf {
	*self = strconv.AppendBool(*self, val)
	return self
}

/*
Appends the string representation of the given error. If the input is nil, this
is a nop. Mutates and returns the receiver.
*/
func (self *Buf) AppendError(val error) *Buf {
	if val == nil {
		return self
	}

	impl, _ := val.(Appender)
	if impl != nil {
		*self = impl.Append(*self)
		return self
	}

	return self.AppendString(val.Error())
}

/*
Appends the text representation of the input, using the `Append` function.
Mutates and returns the receiver.
*/
func (self *Buf) AppendAny(val any) *Buf {
	*self = Append(*self, val)
	return self
}

/*
Appends the text representation of the input, using the `AppendGoString`
function. Mutates and returns the receiver.
*/
func (self *Buf) AppendGoString(val any) *Buf {
	*self = AppendGoString(*self, val)
	return self
}

// Shortcut for appending a formatted string.
func (self *Buf) Fprintf(pat string, val ...any) *Buf {
	_, _ = fmt.Fprintf(self, pat, NoEscUnsafe(val)...)
	return self
}

// Shortcut for appending a formatted string with an idempotent trailing newline.
func (self *Buf) Fprintlnf(pat string, val ...any) *Buf {
	str := fmt.Sprintf(pat, NoEscUnsafe(val)...)
	self.AppendString(str)
	if !HasNewlineSuffix(str) {
		self.AppendNewline()
	}
	return self
}

// Same as `len(buf)`.
func (self Buf) Len() int { return len(self) }

/*
Increases the buffer's length by N zero values.
Mutates and returns the receiver.
*/
func (self *Buf) GrowLen(size int) *Buf {
	*self = GrowLen(*self, size)
	return self
}

/*
Increases the buffer's capacity sufficiently to accommodate N additional
elements. Mutates and returns the receiver.
*/
func (self *Buf) GrowCap(size int) *Buf {
	*self = GrowCap(*self, size)
	return self
}

/*
Truncates the buffer's length, preserving the capacity.
Does not modify the content. Mutates and returns the receiver.
*/
func (self *Buf) Clear() *Buf {
	if self != nil && *self != nil {
		*self = (*self)[:0]
	}
	return self
}
