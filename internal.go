package gg

import (
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

func hasPrefixDigit(val string) bool {
	char := StrHead(val)
	return char >= '0' && char <= '9'
}

func validateLenMatch(one, two int) {
	if one != two {
		panic(Errf(
			`unable to iterate pairwise: length mismatch: %v and %v`,
			one, two,
		))
	}
}
