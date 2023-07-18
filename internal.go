package gg

import (
	"io/fs"
	"math"
	"path/filepath"
	r "reflect"
	"regexp"
	"strings"
	u "unsafe"
)

func typeBitSize(typ r.Type) int { return int(typ.Size() * 8) }

// Borrowed from the standard library. Requires caution.
func noescape(src u.Pointer) u.Pointer {
	out := uintptr(src)
	//nolint:staticcheck
	return u.Pointer(out ^ 0)
}

func errAppendInner(buf Buf, err error) Buf {
	if err != nil {
		buf.AppendString(`: `)
		buf.AppendError(err)
	}
	return buf
}

func errAppendTraceIndentWithNewline(buf Buf, trace Trace) Buf {
	if trace.IsNotEmpty() {
		buf.AppendNewline()
		return errAppendTraceIndent(buf, trace)
	}
	return buf
}

func errAppendTraceIndent(buf Buf, trace Trace) Buf {
	if trace.IsNotEmpty() {
		buf.AppendString(`trace:`)
		buf = trace.AppendIndentTo(buf, 1)
	}
	return buf
}

func isFuncNameAnon(val string) bool {
	const pre = `func`
	return strings.HasPrefix(val, pre) && hasPrefixDigit(val[len(pre):])
}

func hasPrefixDigit(val string) bool { return isDigit(TextHeadByte(val)) }

func isDigit(val byte) bool { return val >= '0' && val <= '9' }

func validateLenMatch(one, two int) {
	if one != two {
		panic(Errf(
			`unable to iterate pairwise: length mismatch: %v and %v`,
			one, two,
		))
	}
}

// Note: `strconv.ParseBool` is too permissive for our taste.
func parseBool(src string, out r.Value) error {
	switch src {
	case `true`:
		out.SetBool(true)
		return nil

	case `false`:
		out.SetBool(false)
		return nil

	default:
		return ErrParse(ErrInvalidInput, src, Type[bool]())
	}
}

/*
Somewhat similar to `filepath.Rel`, but doesn't support `..`, performs
significantly better, and returns the path as-is when it doesn't start
with the given base.
*/
func relOpt(base, src string) string {
	if strings.HasPrefix(src, base) {
		rem := src[len(base):]
		if len(rem) > 0 && rem[0] == filepath.Separator {
			return rem[1:]
		}
	}
	return src
}

func isIntString(val string) bool {
	if len(val) <= 0 {
		return false
	}

	if len(val) > 0 && (val[0] == '+' || val[0] == '-') {
		val = val[1:]
	}

	if len(val) <= 0 {
		return false
	}

	// Note: here we iterate bytes rather than UTF-8 characters because digits
	// are always single byte and we abort on the first mismatch. This may be
	// slightly more efficient than iterating characters.
	for ind := 0; ind < len(val); ind++ {
		if !isDigit(val[ind]) {
			return false
		}
	}
	return true
}

func isCliFlag(val string) bool { return TextHeadByte(val) == '-' }

func isCliFlagValid(val string) bool { return reCliFlag.Get().MatchString(val) }

/*
Must begin with `-` and consist of alphanumeric characters, optionally
containing `-` between those characters.

TODO test.
*/
var reCliFlag = NewLazy(func() *regexp.Regexp {
	return regexp.MustCompile(`^-+[\p{L}\d]+(?:[\p{L}\d-]*[\p{L}\d])?$`)
})

func cliFlagSplit(src string) (_ string, _ string, _ bool) {
	if !isCliFlag(src) {
		return
	}

	ind := strings.IndexRune(src, '=')
	if ind >= 0 {
		return src[:ind], src[ind+1:], true
	}

	return src, ``, false
}

/*
Represents nodes in a linked list. Normally in Go, linked lists tend to be an
anti-pattern; slices perform better in most scenarios, and don't require an
additional abstraction. However, there is one valid scenario for linked lists:
when nodes are pointers to local variables, when those local variables don't
escape, and when they represent addresses to actual memory regions in stack
frames. In this case, this may provide us with a resizable data structure
allocated entirely on the stack, which is useful for book-keeping in recursive
tree-walking or graph-walking algorithms. We currently do not verify if the
trick has the expected efficiency, as the overheads are minimal.
*/
type node[A comparable] struct {
	tail *node[A]
	val  A
}

func (self node[A]) has(val A) bool {
	return self.val == val || (self.tail != nil && self.tail.has(val))
}

func (self *node[A]) cons(val A) (out node[A]) {
	out.tail = self
	out.val = val
	return
}

/*
Suboptimal: doesn't preallocate capacity. We only call this in case of errors,
so the overhead should be negligible.
*/
func (self node[A]) vals() (out []A) {
	out = append(out, self.val)
	node := self.tail
	for node != nil {
		out = append(out, node.val)
		node = node.tail
	}
	return
}

func safeUintToInt(src uint) int {
	if src > math.MaxInt {
		return math.MaxInt
	}
	return int(src)
}

func isByteNewline(val byte) bool { return val == '\n' || val == '\r' }

func errCollMissing[Val, Key any](key Key) Err {
	return Errf(`missing value of type %v for key %v`, Type[Val](), key)
}

func dirEntryToFileName(src fs.DirEntry) (_ string) {
	if src == nil || src.IsDir() {
		return
	}
	return src.Name()
}
