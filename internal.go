package gg

import (
	"path/filepath"
	r "reflect"
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

func isJsonEmpty[A Text](val A) bool {
	return len(val) == 0 || strings.TrimSpace(ToString(val)) == `null`
}

func errAppendInner(buf Buf, err error) Buf {
	if err != nil {
		buf.AppendString(`: `)
		buf.AppendError(err)
	}
	return buf
}

func errAppendTraceIndent(buf Buf, trace Trace) Buf {
	if trace.HasLen() {
		buf.AppendNewline()
		buf.AppendString(`trace:`)
		buf = trace.AppendIndent(buf, 1)
	}
	return buf
}

func isFuncNameAnon(val string) bool {
	const pre = `func`
	return strings.HasPrefix(val, pre) && hasPrefixDigit(val[len(pre):])
}

func hasPrefixDigit(val string) bool { return isDigit(StrHead(val)) }

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
	if len(val) == 0 {
		return false
	}

	if len(val) > 0 && (val[0] == '+' || val[0] == '-') {
		val = val[1:]
	}

	if len(val) == 0 {
		return false
	}
	for ind := 0; ind < len(val); ind++ {
		if !isDigit(val[ind]) {
			return false
		}
	}
	return true
}
