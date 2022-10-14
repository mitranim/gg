package gg

import (
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

// Compares two text chunks via `==`.
func StrEq[A Text](one, two A) bool { return ToString(one) == ToString(two) }

/*
Returns the underlying data pointer of the given string or byte slice.
Mutations may trigger segfaults or cause undefined behavior.
*/
func StrDat[A Text](src A) *byte { return CastUnsafe[*byte](src) }

/*
Allocation-free conversion. Reinterprets arbitrary text as a string. If the
string is used with an API that relies on string immutability, for example as a
map key, the source memory must not be mutated afterwards.
*/
func ToString[A Text](val A) string { return CastUnsafe[string](val) }

/*
Allocation-free conversion. Reinterprets arbitrary text as bytes. If the source
was a string, the output must NOT be mutated. Mutating memory that belongs to a
string may produce segfaults or undefined behavior.
*/
func ToBytes[A Text](val A) []byte {
	return u.Slice(CastUnsafe[*byte](val), len(val))
}

/*
Converts arguments to strings and concatenates the results. See `StringCatch`
for the encoding rules. Also see `JoinDense` for a more limited but more
efficient version that doesn't involve `any`.
*/
func Str(src ...any) string { return JoinAny(src, ``) }

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
		buf.GrowCap(Sum(src, StrLen[A]) + (len(sep) * (len(src) - 1)))

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

// Self-explanatory. Splits the given text into lines.
func SplitLines[A Text](src A) []string {
	if len(src) == 0 {
		return nil
	}
	return RE_NEWLINE().Split(ToString(src), -1)
}

// Matches any newline. Supports Windows, Unix, and old MacOS styles.
var RE_NEWLINE = Lazy1(regexp.MustCompile, `(?:\r\n|\r|\n)`)

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

// Missing/private half of `strings.TrimSpace`. Trims only the prefix.
func TrimSpacePrefix[A Text](src A) A {
	return CastUnsafe[A](strings.TrimLeftFunc(ToString(src), unicode.IsSpace))
}

// Missing/private half of `strings.TrimSpace`. Trims only the suffix.
func TrimSpaceSuffix[A Text](src A) A {
	return CastUnsafe[A](strings.TrimRightFunc(ToString(src), unicode.IsSpace))
}

/*
Regexp for splitting arbitrary text into words, Unicode-aware. Used by
`ToWords`.
*/
var ReWord = Lazy1(
	regexp.MustCompile,
	`\p{Lu}\p{Ll}+\d*|\p{Lu}+\d*|\p{Ll}+\d*`,
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
func (self Words) Title() Words { return MapMut(self, strings.Title) }

/*
Converts the first word to Titlecase and each other word to lowercase. Mutates
and returns the receiver.
*/
func (self Words) Sentence() Words { return self.Lower().MapHead(strings.Title) }

/*
Converts the first word to lowercase and each other word to Titlecase. Mutates
and returns the receiver.
*/
func (self Words) Camel() Words { return self.Title().MapHead(strings.ToLower) }

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

// Uses `utf8.RuneCountInString` to count chars in arbitrary text.
func CharCount[A Text](val A) int {
	return utf8.RuneCountInString(ToString(val))
}
