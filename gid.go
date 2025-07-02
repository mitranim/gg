package gg

import (
	"runtime"
	"strings"
)

/*
Allows user code to plug in an arbitrary function for getting the goroutine ID.
All code in this package tries this override before falling back on [Gid].
On platforms where we don't support the fast mode of [Gid], this allows user
code to plug-in an alternative.
*/
var GidFunc func() uint64

/*
Returns the current goroutine ID. On major CPU architectures, in specific
versions of Go, uses a very fast implementation which hacks into Go internals:

  - Supported architectures: `386 || amd64 || arm || arm64 || riscv64 || s390x`.
  - Supported Go versions: `go1.23 || go1.24 || go1.25`.
  - Only the official Go compiler (build flag `gc`) is supported.

In "unsupported" cases, we fall back on parsing the output of [runtime.Stack],
which is extremely slow.

We restrict the fast path to the Go versions whose internals are known to this
library at the time of writing, because the memory layout of the structures
we're accessing changes between releases, and we could end up mis-identifying
arbitrary non-GID memory as a GID.
*/
func Gid() uint64 { return gid() }

/*
Fallback version of [Gid], used on unsupported architectures and in unsupported
Go versions. Parses the current goroutine id from the output of [runtime.Stack].
Panics if the output format is unrecognized. Very slow.

Placed in this file and not in `gid_internal_slow.go` because we need this
available in testing.
*/
func gidSlow() uint64 {
	const errPre = `unable to determine goroutine id`

	/**
	`30` is the max text length of "goroutine N" where N is the byte count of
	`uint64` encoded in base 10 in ASCII / UTF-8.
	*/
	var arr [30]byte
	buf := NoEscUnsafe(arr[:]) // Saves an allocation.
	buf = buf[:runtime.Stack(buf, false)]
	str := ToString(buf)

	const pre = `goroutine `
	if !strings.HasPrefix(str, pre) {
		panic(Errv(errPre + `: unrecognized stack format`))
	}

	str = str[len(pre):]
	var ind uint64
	for len(str) > 0 && isDigit(str[ind]) {
		ind++
	}

	if ind <= 0 {
		panic(Errv(errPre + `: missing id in stack`))
	}

	out := ParseTo[uint64](str[:ind])
	runtime.KeepAlive(arr)
	return out
}

func getGid() uint64 {
	fun := GidFunc
	if fun != nil {
		return fun()
	}
	return Gid()
}
