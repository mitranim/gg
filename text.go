package gg

import (
	"math"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
	u "unsafe"
)

/*
Same as `len`. Limited to `Text` types but can be passed to higher-order
functions.
*/
func TextLen[A Text](val A) int { return len(val) }

// True if len <= 0. Inverse of `IsTextNotEmpty`.
func IsTextEmpty[A Text](val A) bool { return len(val) <= 0 }

// True if len > 0. Inverse of `IsTextEmpty`.
func IsTextNotEmpty[A Text](val A) bool { return len(val) > 0 }

// Returns the first byte or 0.
func TextHeadByte[A Text](val A) byte {
	if len(val) > 0 {
		return val[0]
	}
	return 0
}

// Returns the last byte or 0.
func TextLastByte[A Text](val A) byte {
	if len(val) > 0 {
		return val[len(val)-1]
	}
	return 0
}

/*
Like `utf8.DecodeRuneInString`, but faster at the time of writing, and without
`utf8.RuneError`. On decoding error, the result is `(0, 0)`.
*/
func TextHeadChar[A Text](src A) (char rune, size int) {
	for ind, val := range ToText[string](src) {
		if ind == 0 {
			char = val
			size = len(src)
		} else {
			size = ind
			break
		}
	}
	return
}

/*
True if the inputs would be `==` if compared as strings. When used on typedefs
of `[]byte`, this is the same as `bytes.Equal`.
*/
func TextEq[A Text](one, two A) bool { return ToString(one) == ToString(two) }

/*
Similar to `unsafe.StringData`, but takes arbitrary text as input. Returns the
pointer to the first byte of the underlying data array for the given string or
byte slice. Use caution. Mutating the underlying data may trigger segfaults or
cause undefined behavior.
*/
func TextDat[A Text](val A) *byte { return CastUnsafe[*byte](val) }

/*
TODO: in `TextDat`, consider using `unsafe.StringData` for strings and
`unsafe.SliceData` for byte slices. Something like the following pseudocode.
May run slower.

func TextDat[A Text](val A) *byte {
	if string {
		return u.StringData(val)
	}
	if bytes {
		return u.SliceData(val)
	}
	panic(`unreachable`)
}
*/

/*
Allocation-free conversion between two types conforming to the `Text`
constraint, typically variants of `string` and/or `[]byte`.
*/
func ToText[Out, Src Text](val Src) Out {
	tar := CastUnsafe[Out](val)

	/**
	Implementation note. We could also write the condition as shown below, but
	this would be significantly slower than the unsafe trick:

		Kind[Src]() == r.String && Kind[Out]() == r.Slice

	In addition, sizeof lets us ensure that the target can be cast into
	`SliceHeader` without affecting other memory.
	*/
	if u.Sizeof(Zero[Src]()) == SizeofString && u.Sizeof(Zero[Out]()) == SizeofSliceHeader {
		CastUnsafe[*SliceHeader](&tar).Cap = len(tar)
	}

	return tar
}

/*
Allocation-free conversion. Reinterprets arbitrary text as a string. If the
string is used with an API that relies on string immutability, for example as a
map key, the source memory must not be mutated afterwards.
*/
func ToString[A Text](val A) string { return CastUnsafe[string](val) }

/*
Implementation note. `ToString` could be written as shown below. This passes our
test, but runs marginally slower than our current implementation, and does not
improve correctness, because `TextDat` also makes assumptions about the
underlying structure of the string header.

	func ToString[A Text](val A) string { return u.String(TextDat(val), len(val)) }
*/

/*
Allocation-free conversion. Reinterprets arbitrary text as bytes. If the source
was a string, the output must NOT be mutated. Mutating memory that belongs to a
string may produce segfaults or undefined behavior.
*/
func ToBytes[A Text](val A) []byte { return u.Slice(TextDat(val), len(val)) }

/*
Converts arguments to strings and concatenates the results. See `StringCatch`
for the encoding rules. Also see `JoinDense` for a simpler version that doesn't
involve `any`.
*/
func Str(src ...any) string { return JoinAny(src, ``) }

/*
Similar to `Str`. Concatenates string representations of the input values.
Additionally, if the output is non-empty and doesn't end with a newline
character, appends `Newline` at the end.
*/
func Strln(src ...any) string {
	switch len(src) {
	case 0:
		return ``

	case 1:
		return AppendNewlineOpt(String(src[0]))

	default:
		var buf Buf
		buf.AppendAnysln(src...)
		return buf.String()
	}
}

/*
Converts arguments to strings and joins the results with a single space. See
`StringCatch` for encoding rules. Also see `JoinSpaced` for a more limited but
more efficient version that doesn't involve `any`.
*/
func Spaced(src ...any) string { return JoinAny(src, Space) }

/*
Converts arguments to strings and joins the results with a single space,
ignoring empty strings. See `StringCatch` for the encoding rules. Also see
`JoinSpacedOpt` for a more limited but more efficient version that doesn't
involve `any`.
*/
func SpacedOpt(src ...any) string { return JoinAnyOpt(src, Space) }

/*
Similar to `strings.Join` but takes `[]any`, converting elements to strings. See
`StringCatch` for the encoding rules. Also see `Join`, `JoinOpt`,
`JoinAnyOpt`.
*/
func JoinAny(src []any, sep string) string {
	switch len(src) {
	case 0:
		return ``

	case 1:
		return String(src[0])

	default:
		var buf Buf
		for ind, src := range src {
			if ind > 0 {
				buf.AppendString(sep)
			}
			buf.AppendAny(src)
		}
		return buf.String()
	}
}

// Like `JoinAny` but ignores empty strings.
func JoinAnyOpt(src []any, sep string) string {
	switch len(src) {
	case 0:
		return ``

	case 1:
		return String(src[0])

	default:
		var buf Buf

		for ind, src := range src {
			len0 := buf.Len()
			if ind > 0 {
				buf.AppendString(sep)
			}
			len1 := buf.Len()

			buf.AppendAny(src)

			if ind > 0 && buf.Len() == len1 {
				buf.TruncLen(len0)
			}
		}

		return buf.String()
	}
}

// Concatenates the given text without any separators.
func JoinDense[A Text](val ...A) string { return Join(val, ``) }

// Joins the given strings with a space.
func JoinSpaced[A Text](val ...A) string { return Join(val, Space) }

// Joins non-empty strings with a space.
func JoinSpacedOpt[A Text](val ...A) string { return JoinOpt(val, Space) }

// Joins the given strings with newlines.
func JoinLines[A Text](val ...A) string { return Join(val, Newline) }

// Joins non-empty strings with newlines.
func JoinLinesOpt[A Text](val ...A) string { return JoinOpt(val, Newline) }

/*
Similar to `strings.Join` but works on any input compatible with the `Text`
interface. Also see `JoinOpt`, `JoinAny`, `JoinAnyOpt`.
*/
func Join[A Text](src []A, sep string) string {
	switch len(src) {
	case 0:
		return ``

	case 1:
		return ToString(src[0])

	default:
		var buf Buf
		buf.GrowCap(Sum(src, TextLen[A]) + (len(sep) * (len(src) - 1)))

		buf.AppendString(ToString(src[0]))
		for _, src := range src[1:] {
			buf.AppendString(sep)
			buf.AppendString(ToString(src))
		}
		return buf.String()
	}
}

/*
Similar to `strings.Join` but works for any input compatible with the `Text`
interface and ignores empty strings.
*/
func JoinOpt[A Text](src []A, sep string) string {
	switch len(src) {
	case 0:
		return ``

	case 1:
		return ToString(src[0])

	default:
		var size int
		for _, src := range src {
			wid := len(src)
			if wid > 0 {
				size = size + wid + len(sep)
			}
		}

		var buf Buf
		buf.GrowCap(size)

		var found bool
		for _, src := range src {
			if len(src) > 0 {
				if found {
					buf.AppendString(sep)
				}
				buf = append(buf, src...)
				found = true
			}
		}
		return buf.String()
	}
}

/*
Similar to `strings.SplitN` for N = 1. More efficient: returns a tuple instead
of allocating a slice. Safer: returns zero values if split doesn't succeed.
*/
func Split2[A Text](src A, sep string) (A, A) {
	ind := strings.Index(ToString(src), sep)
	if ind >= 0 {
		return src[:ind], src[ind+len(sep):]
	}
	return src, Zero[A]()
}

/*
Splits the given text into lines. The resulting strings do not contain any
newline characters. If the input is empty, the output is nil. The following
sequences are considered newlines: "\r\n", "\r", "\n".
*/
func SplitLines[A Text](src A) []string {
	if len(src) <= 0 {
		return nil
	}

	// Probably vastly suboptimal. Needs tuning.
	return ReNewline.Get().Split(ToString(src), -1)
}

// Matches any newline. Supports Windows, Unix, and old MacOS styles.
var ReNewline = NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`(?:\r\n|\r|\n)`)
})

/*
Similar to `SplitLines`, but splits only on the first newline occurrence,
returning the first line and the remainder, plus the number of bytes in the
elided line separator. The following sequences are considered newlines:
"\r\n", "\r", "\n".
*/
func SplitLines2[A Text](src A) (A, A, int) {
	size := len(src)
	limit := size - 1

	for ind, char := range ToString(src) {
		if char == '\r' {
			if ind < limit && src[ind+1] == '\n' {
				return src[:ind], src[ind+2:], 2
			}
			return src[:ind], src[ind+1:], 1
		}
		if char == '\n' {
			return src[:ind], src[ind+1:], 1
		}
	}
	return src, Zero[A](), 0
}

/*
Searches for the given separator and returns the part of the text before the
separator, removing that prefix from the original text referenced by the
pointer. The separator is excluded from both chunks. As a special case, if the
separator is empty, pops the entire source text.
*/
func TextPop[Src, Sep Text](ptr *Src, sep Sep) Src {
	if ptr == nil {
		return Zero[Src]()
	}

	src := *ptr

	if len(sep) == 0 {
		PtrClear(ptr)
		return src
	}

	ind := strings.Index(ToString(src), ToString(sep))
	if !(ind >= 0 && ind < len(src)) {
		PtrClear(ptr)
		return src
	}

	*ptr = src[ind+len(sep):]
	return src[:ind]
}

// True if the string ends with a line feed or carriage return.
func HasNewlineSuffix[A Text](src A) bool {
	val := TextLastByte(src)
	return val == '\n' || val == '\r'
}

/*
If the given text is non-empty and does not end with a newline character,
appends `Newline` and returns the result. Otherwise returns the text unchanged.
If the input type is a typedef of `[]byte` and has enough capacity, it's
mutated. In other cases, the text is reallocated. Also see
`Buf.AppendNewlineOpt` and `Strln`.
*/
func AppendNewlineOpt[A Text](val A) A {
	if len(val) > 0 && !HasNewlineSuffix(val) {
		return ToText[A](append([]byte(val), Newline...))
	}
	return val
}

// Missing/private half of `strings.TrimSpace`. Trims only the prefix.
func TrimSpacePrefix[A Text](src A) A {
	return ToText[A](strings.TrimLeftFunc(ToString(src), unicode.IsSpace))
}

// Missing/private half of `strings.TrimSpace`. Trims only the suffix.
func TrimSpaceSuffix[A Text](src A) A {
	return ToText[A](strings.TrimRightFunc(ToString(src), unicode.IsSpace))
}

/*
Regexp for splitting arbitrary text into words, Unicode-aware. Used by
`ToWords`.
*/
var ReWord = NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`\p{Lu}\p{Ll}+\d*|\p{Lu}+\d*|\p{Ll}+\d*`)
})

/*
Splits arbitrary text into words, Unicode-aware. Suitable for conversion between
typographic cases such as `camelCase` and `snake_case`.
*/
func ToWords[A Text](val A) Words {
	return ReWord.Get().FindAllString(ToString(val), -1)
}

/*
Tool for converting between typographic cases such as `camelCase` and
`snake_case`.
*/
type Words []string

// Combines the words via " ".
func (self Words) Spaced() string { return self.Join(` `) }

// Combines the words via "_".
func (self Words) Snake() string { return self.Join(`_`) }

// Combines the words via "-".
func (self Words) Kebab() string { return self.Join(`-`) }

// Combines the words via "".
func (self Words) Dense() string { return self.Join(``) }

// Combines the words via ",".
func (self Words) Comma() string { return self.Join(`,`) }

// Converts each word to lowercase. Mutates and returns the receiver.
func (self Words) Lower() Words { return MapMut(self, strings.ToLower) }

// Converts each word to UPPERCASE. Mutates and returns the receiver.
func (self Words) Upper() Words { return MapMut(self, strings.ToUpper) }

// Converts each word to Titlecase. Mutates and returns the receiver.
func (self Words) Title() Words {
	//nolint:staticcheck
	return MapMut(self, strings.Title)
}

/*
Converts the first word to Titlecase and each other word to lowercase. Mutates
and returns the receiver.
*/
func (self Words) Sentence() Words {
	//nolint:staticcheck
	return self.MapHead(strings.Title).MapTail(strings.ToLower)
}

/*
Converts the first word to lowercase and each other word to Titlecase. Mutates
and returns the receiver.
*/
func (self Words) Camel() Words {
	//nolint:staticcheck
	return self.MapHead(strings.ToLower).MapTail(strings.Title)
}

// Same as `strings.Join`.
func (self Words) Join(val string) string { return strings.Join(self, val) }

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

/*
Mutates the receiver by replacing elements, other than the first, with the
results of the given function.
*/
func (self Words) MapTail(fun func(string) string) Words {
	if len(self) > 0 {
		MapMut(self[1:], fun)
	}
	return self
}

// Uses `utf8.RuneCountInString` to count chars in arbitrary text.
func CharCount[A Text](val A) int {
	return utf8.RuneCountInString(ToString(val))
}

/*
Similar to `src[start:end]`, but instead of slicing text at byte positions,
slices text at character positions. Similar to `string([]rune(src)[start:end])`,
but slightly more performant and more permissive.
*/
func TextCut[A Text](src A, start, end int) (_ A) {
	if !(end > start) {
		return
	}

	startInd := 0
	endInd := len(src)
	charInd := 0

	for byteInd := range ToString(src) {
		if charInd == start {
			startInd = byteInd
		}
		if charInd == end {
			endInd = byteInd
			break
		}
		charInd++
	}

	return src[startInd:endInd]
}

/*
Truncates the given text to the given total count of Unicode characters
(not bytes) with an ellipsis, if needed. The total count includes the ellipsis
character '…'. The limit's can't exceed `math.MaxInt`.
*/
func Ellipsis[A Text](src A, limit uint) string {
	if limit == 0 {
		return ``
	}

	var lim int
	if limit > math.MaxInt {
		lim = math.MaxInt
	} else {
		lim = int(limit)
	}

	const suf = `…`
	const sufCharLen = 1

	str := ToString(src)
	prevInd := 0
	nextInd := 0
	charInd := 0

	for nextInd = range str {
		if charInd+sufCharLen > lim {
			return str[:prevInd] + suf
		}
		prevInd = nextInd
		charInd++
	}
	return str
}