package gg

import (
	"bytes"
	r "reflect"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
	u "unsafe"
)

/*
Same as `len`. Sometimes useful for higher-order functions. Note that this does
NOT count Unicode characters. For that, use `CharCount`.
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
Implementation note. We could write `TextDat` as following, but it would not
be an improvement, because it still makes assumptions about the underlying
structure of the data, specifically it assumes that strings and byte slices
have a different width. At the time of writing, Go doesn't seem to provide a
safe and free way to detect if we have `~string` or `~[]byte`. A type switch
on `any(src)` works only for core types such as `string`, but not for typedefs
conforming to `~string` and `~[]byte`. Alternatives involve overheads such as
calling interface methods of `reflect.Type`, which would stop this function
from being a free cast.

	func TextDat[A Text](src A) *byte {
		if u.Sizeof(src) == SizeofString {
			return u.StringData(string(src))
		}
		if u.Sizeof(src) == SizeofSlice {
			return u.SliceData([]byte(src))
		}
		panic(`unreachable`)
	}
*/

/*
Allocation-free conversion between two types conforming to the `Text`
constraint, typically variants of `string` and/or `[]byte`.
*/
func ToText[Out, Src Text](src Src) Out {
	out := CastUnsafe[Out](src)

	/**
	Implementation note. We could also write the condition as shown below:

		Kind[Src]() == r.String && Kind[Out]() == r.Slice

	But the above would be measurably slower than the unsafe trick.
	In addition, sizeof lets us ensure that the target can be cast into
	`SliceHeader` without affecting other memory.
	*/
	if u.Sizeof(src) == SizeofString && u.Sizeof(out) == SizeofSliceHeader {
		CastUnsafe[*SliceHeader](&out).Cap = len(out)
	}

	return out
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
character, appends a newline at the end.
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
func Spaced(src ...any) string { return JoinAny(src, ` `) }

/*
Converts arguments to strings and joins the results with a single space,
ignoring empty strings. See `StringCatch` for the encoding rules. Also see
`JoinSpacedOpt` for a more limited but more efficient version that doesn't
involve `any`.
*/
func SpacedOpt(src ...any) string { return JoinAnyOpt(src, ` `) }

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
func JoinSpaced[A Text](val ...A) string { return Join(val, ` `) }

// Joins non-empty strings with a space.
func JoinSpacedOpt[A Text](val ...A) string { return JoinOpt(val, ` `) }

// Joins the given strings with newlines.
func JoinLines[A Text](val ...A) string { return Join(val, "\n") }

// Joins non-empty strings with newlines.
func JoinLinesOpt[A Text](val ...A) string { return JoinOpt(val, "\n") }

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
		buf.GrowCap(SumBy(src, TextLen[A]) + (len(sep) * (len(src) - 1)))

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
				size += wid + len(sep)
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
				found = true
				buf = append(buf, src...)
			}
		}
		return buf.String()
	}
}

/*
Similar to `strings.Split` and `bytes.Split`. Differences:

  - Supports all text types.
  - Returns nil for empty input.
*/
func Split[A Text](src, sep A) []A {
	if len(src) <= 0 {
		return nil
	}
	if Kind[A]() == r.String {
		return CastSlice[A](strings.Split(ToString(src), ToString(sep)))
	}
	return CastSlice[A](bytes.Split(ToBytes(src), ToBytes(sep)))
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
newline characters. If the input is empty, the output is empty. Avoids
information loss: preserves empty lines, allowing the caller to transform and
join the lines without losing blanks. The following sequences are considered
newlines: "\r\n", "\r", "\n".
*/
func SplitLines[A Text](src A) []A {
	/**
	In our benchmark in Go 1.20.2, this runs about 20-30 times faster than the
	equivalent regexp-based implementation.

	It would be much simpler to use `strings.FieldsFunc` and `bytes.FieldsFunc`,
	but they would elide empty lines, losing information and making this
	non-reversible. They would also be about 2 times slower.

	TODO simpler implementation.
	*/

	var out []A
	var prev int
	var next int
	max := len(src)

	/**
	Iterating bytes is significantly faster than runes, and in valid UTF-8 it's
	not possible to encounter '\r' or '\n' in multi-byte characters, making this
	safe for valid text.
	*/
	for next < max {
		char := src[next]

		if char == '\r' && next < len(src)-1 && src[next+1] == '\n' {
			out = append(out, src[prev:next])
			next = next + 2
			prev = next
			continue
		}

		if char == '\n' || char == '\r' {
			out = append(out, src[prev:next])
			next++
			prev = next
			continue
		}

		next++
	}

	if next > 0 {
		out = append(out, src[prev:next])
	}
	return out
}

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
	return isByteNewline(TextLastByte(src))
}

/*
If the given text is non-empty and does not end with a newline character,
appends a newline and returns the result. Otherwise returns the text unchanged.
If the input type is a typedef of `[]byte` and has enough capacity, it's
mutated. In other cases, the text is reallocated. Also see
`Buf.AppendNewlineOpt` and `Strln`.
*/
func AppendNewlineOpt[A Text](val A) A {
	if len(val) > 0 && !HasNewlineSuffix(val) {
		return ToText[A](append([]byte(val), '\n'))
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
	return regexp.MustCompile(`\p{Lu}+[\p{Ll}\d]*|[\p{Ll}\d]+`)
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

// Combines the words via "".
func (self Words) Dense() string { return self.Join(``) }

// Combines the words via " ".
func (self Words) Spaced() string { return self.Join(` `) }

// Combines the words via "_".
func (self Words) Snake() string { return self.Join(`_`) }

// Combines the words via "-".
func (self Words) Kebab() string { return self.Join(`-`) }

// Combines the words via ",".
func (self Words) Comma() string { return self.Join(`,`) }

// Combines the words via "|".
func (self Words) Piped() string { return self.Join(`|`) }

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
Truncates text to the given count of Unicode characters (not bytes). The limit
can't exceed `math.MaxInt`. Also see `TextTruncWith` which is more general.
*/
func TextTrunc[A Text](src A, limit uint) (_ A) {
	return TextTruncWith(src, Zero[A](), limit)
}

/*
Shortcut for `TextTruncWith(src, "…")`. Truncates the given text to the given total
count of Unicode characters with an ellipsis.
*/
func TextEllipsis[A Text](src A, limit uint) A {
	return TextTruncWith(src, ToText[A](`…`), limit)
}

/*
Truncates the given text to the given total count of Unicode characters (not
bytes) with a suffix. If the text is under the limit, it's returned unchanged,
otherwise it's truncated and the given suffix is appended. The total count
includes the character count of the given suffix string. The limit can't exceed
`math.MaxInt`. Also see shortcut `TextEllipsis` which uses this with the
ellipsis character '…'.
*/
func TextTruncWith[A Text](src, suf A, limit uint) A {
	if limit == 0 {
		return Zero[A]()
	}

	lim := safeUintToInt(limit)
	sufCharLen := CharCount(suf)
	str := ToString(src)
	prevInd := 0
	nextInd := 0
	charInd := 0

	for nextInd = range str {
		if charInd+sufCharLen > lim {
			return ToText[A](str[:prevInd] + ToString(suf))
		}
		prevInd = nextInd
		charInd++
	}
	return src
}
