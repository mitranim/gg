package gg

import (
	"regexp"
	"strings"
)

/*
Same as `len`. Limited to `Text` types but can be passed to higher-order
functions.
*/
func StrLen[A Text](val A) int { return len(val) }

// Returns the first byte or 0.
func StrHead[A Text](val A) byte {
	if len(val) > 0 {
		return val[0]
	}
	return 0
}

// Returns the last byte or 0.
func StrLast[A Text](val A) byte {
	if len(val) > 0 {
		return val[len(val)-1]
	}
	return 0
}

// True if len > 0.
func IsStrNonEmpty[A Text](val A) bool { return len(val) > 0 }

// True if len == 0.
func IsStrEmpty[A Text](val A) bool { return len(val) == 0 }

/*
Returns the underlying data pointer of the given string or byte slice.
Mutations may trigger segfaults or cause undefined behavior.
*/
func StrDat[A Text](src A) *byte { return CastUnsafe[*byte](src) }

/*
Allocation-free conversion. Reinterprets an arbitrary string-like or bytes-like
value as a regular string.
*/
func ToString[A Text](val A) string { return CastUnsafe[string](val) }

/*
Allocation-free conversion. Reinterprets an arbitrary string-like or bytes-like
value as bytes.
*/
func ToBytes[A Text](val A) []byte {
	slice := CastUnsafe[SliceHeader](val)
	slice.Cap = slice.Len
	return CastUnsafe[[]byte](slice)
}

// Joins the given strings with newlines.
func JoinLines(val ...string) string { return strings.Join(val, Newline) }

// Joins non-empty strings with newlines.
func JoinLinesOpt[A Text](val ...A) string { return JoinOpt(val, Newline) }

// Joins the given strings with a space.
func Spaced(val ...string) string { return strings.Join(val, Space) }

// Joins non-empty strings with a space.
func SpacedOpt[A Text](val ...A) string { return JoinOpt(val, Space) }

// Similar to `strings.Join` but ignores empty strings.
func JoinOpt[A Text](val []A, sep string) string {
	if len(val) == 0 {
		return ``
	}
	if len(val) == 1 {
		return ToString(val[0])
	}

	var size int
	for _, val := range val {
		wid := len(val)
		if wid > 0 {
			size = size + wid + len(sep)
		}
	}

	var buf Buf
	buf.GrowCap(size)

	var found bool
	for _, val := range val {
		if len(val) > 0 {
			if found {
				buf.AppendString(sep)
			}
			buf = append(buf, val...)
			found = true
		}
	}
	return buf.String()
}

/*
Searches for the given separator and returns the part of the string before the
separator, removing that prefix from the original string referenced by the
pointer. The separator is excluded from both chunks. As a special case, if
the separator is empty, pops the entire given string.
*/
func StrPop[A, B ~string](ptr *A, sep B) A {
	if ptr == nil {
		return ``
	}

	src := *ptr

	if len(sep) == 0 {
		*ptr = ``
		return src
	}

	ind := strings.Index(string(src), string(sep))
	if !(ind >= 0 && ind < len(src)) {
		*ptr = ``
		return src
	}

	*ptr = src[ind+len(sep):]
	return src[:ind]
}

// True if the string ends with a line feed or carriage return.
func HasNewlineSuffix[A Text](val A) bool {
	return StrLast(val) == '\n' || StrLast(val) == '\r'
}

/*
Regexp for splitting arbitrary text into words, Unicode-aware. Used by
`ToWords`.
*/
var ReWord = Lazy1(
	regexp.MustCompile,
	`\p{Lu}[\p{Ll}\d]*|\p{Lu}[\p{Lu}\d]*|\p{Ll}[\p{Ll}\d]*`,
)

/*
Splits arbitrary text into words, Unicode-aware. Suitable for conversion between
typographic cases such as `camelCase` and `snake_case`.
*/
func ToWords[A Text](val A) Words {
	return ReWord().FindAllString(ToString(val), -1)
}

/*
Tool for converting between typographic cases such as `camelCase` and
`snake_case`.
*/
type Words []string

func (self Words) Spaced() string { return self.Join(` `) }
func (self Words) Snake() string  { return self.Join(`_`) }
func (self Words) Kebab() string  { return self.Join(`-`) }
func (self Words) Solid() string  { return self.Join(``) }
func (self Words) Comma() string  { return self.Join(`,`) }

func (self Words) Lower() Words    { return self.MapMut(strings.ToLower) }
func (self Words) Upper() Words    { return self.MapMut(strings.ToUpper) }
func (self Words) Title() Words    { return self.MapMut(strings.Title) }
func (self Words) Sentence() Words { return self.Lower().MapHead(strings.Title) }
func (self Words) Camel() Words    { return self.Title().MapHead(strings.ToLower) }

// Same as `strings.Join`.
func (self Words) Join(val string) string { return strings.Join(self, val) }

/*
Mutates the receiver by replacing each element with the result of calling the
given function on that element.
*/
func (self Words) MapMut(fun func(string) string) Words {
	if fun != nil {
		for ind := range self {
			self[ind] = fun(self[ind])
		}
	}
	return self
}

/*
Mutates the receiver by replacing the first element with the result of calling
the given function on that element. If the receiver is empty, this is a nop.
*/
func (self Words) MapHead(fun func(string) string) Words {
	if fun != nil && len(self) > 0 {
		self[0] = fun(self[0])
	}
	return self
}
