package gg

import (
	"runtime"
	"sync"
	"sync/atomic"
)

/*
This file implements stack IDs, which are somewhat similar to goroutine IDs,
but must be stored on the stack manually, by calling either [WithGivenSid]
or [WithSid]; the latter ensures uniqueness of concurrent SIDs.

This implementation is vaguely inspired by `github.com/jtolio/gls` and uses
basically the same mechanism. The code is original; the differences are mainly
cosmetic. Our version seems to perform slightly better in Go 1.24.0, but not by
much. The portability between different Go versions has not been investigated.
Our version does not have special support for JS/WASM ports.
*/

/*
Short for "stack ID". Retrieves the nearest ID, if any, encoded into the stack
via [WithSid] or [WithGivenSid]. Also see [Gid] for goroutine ID, and [DynVar]
for dynamic variables.

Valid SIDs start at 1; 0 should be considered unset / missing.

This was implemented out of horrified amazement at the idea,
and probably has no real use case.
*/
func Sid() (out uint64) {
	/**
	In our testing, this is always stack-allocated. We only pass this to
	`runtime.Callers` where it doesn't escape. May vary by Go version.
	Last tested in Go 1.24.0.
	*/
	buf := make([]uintptr, 64)

	// Skip this PC and `runtime.Callers` itself.
	off := 2

	for {
		size := runtime.Callers(off, buf)
		if size <= 0 {
			break
		}

		// Linter is bugged.
		// nolint ineffassign
		off += size

		for _, addr := range buf[:size] {
			if addr == 0 || addr == sidEndPc {
				return
			}

			dig, ok := sidPcsToDigits[addr]
			if !ok {
				continue
			}

			out <<= sidDigitBits
			out += uint64(dig)
		}

		break
	}

	return
}

/*
Encodes a unique stack ID (SID) into the stack and runs the given function
inside. Calling [Sid] retrieves the SID, unless shadowed by an inner call to
[WithSid] or [WithGivenSid]. The minimum SID value is 1.

This is vaguely similar to goroutine ID (GID). Requires a bit of extra setup,
but doesn't depend on Go internals, and allows to provide arbitrary IDs via
[WithGivenSid].

When invoked concurrently, or when nested / overlapping, each invocation of
[WithSid] provides a unique SID. When the call ends, the SID is reclaimed into
the pool, to be reused in future calls. This is done to avoid running out of
unused SIDs in the edge case of calling this function [math.MaxUint64] times.
*/
func WithSid(fun func()) {
	if fun == nil {
		return
	}
	sid := sidGet()
	defer sidFree(sid)
	sid_end(sid, fun)
}

/*
Encodes the given stack ID (SID) into the stack and runs the given function
inside. Calling [Sid] retrieves the given SID, unless shadowed by an inner
call to [WithSid] or [WithGivenSid].
*/
func WithGivenSid(sid uint64, fun func()) {
	if fun == nil {
		return
	}
	sid_end(sid, fun)
}

/* Internal */

/*
Known minor issue: we never truly "free" previously-used SIDs. This means that
if a program spawns a million billion concurrent goroutines which all use
[WithSid], and they all exit, and the program never does that again, their
identifiers stay in memory until the program exits. However, that only happens
for _concurrent_ goroutines, whose memory consumption should gigantically dwarf
our buffer's size, making this cleanup a low priority.

`sync.Pool` seems unsuitable here, because it's allowed to "lose" values,
thus preventing us from recycling them.
*/
var sidPoolMutex sync.Mutex
var sidPool []uint64
var sidHighest atomic.Uint64

func sidGet() uint64 {
	sid, ok := sidFromPool()
	if ok {
		return sid
	}

	/**
	No overflow check for now. The overflow only affects code which has at least
	`math.MaxUint64` concurrent or otherwise overlapping calls to this function,
	not necessarily on different goroutines, which is unreal for real programs.
	*/
	return sidHighest.Add(1)
}

func sidFree(sid uint64) {
	defer Lock(&sidPoolMutex).Unlock()
	sidPool = append(sidPool, sid)
}

func sidFromPool() (_ uint64, _ bool) {
	defer Lock(&sidPoolMutex).Unlock()

	ind := len(sidPool) - 1
	if ind < 0 {
		return
	}

	sid := sidPool[ind]
	sidPool = sidPool[:ind]
	return sid, true
}

type sidFunc = func(uint64, func())

const sidDigitBits = 4

var sidEndPc = sidFuncToPc(sid_end)
var sidDigitsToFuns [16]sidFunc
var sidDigitsToPcs [16]uintptr
var sidPcsToDigits = make(map[uintptr]byte, 16)

/*
We have to use `init` instead of variable initialization expressions
because otherwise the compiler complains about an initialization cycle.
*/
func init() {
	sidDigitsToFuns[0x0] = sid_digit_0x0
	sidDigitsToFuns[0x1] = sid_digit_0x1
	sidDigitsToFuns[0x2] = sid_digit_0x2
	sidDigitsToFuns[0x3] = sid_digit_0x3
	sidDigitsToFuns[0x4] = sid_digit_0x4
	sidDigitsToFuns[0x5] = sid_digit_0x5
	sidDigitsToFuns[0x6] = sid_digit_0x6
	sidDigitsToFuns[0x7] = sid_digit_0x7
	sidDigitsToFuns[0x8] = sid_digit_0x8
	sidDigitsToFuns[0x9] = sid_digit_0x9
	sidDigitsToFuns[0xa] = sid_digit_0xa
	sidDigitsToFuns[0xb] = sid_digit_0xb
	sidDigitsToFuns[0xc] = sid_digit_0xc
	sidDigitsToFuns[0xd] = sid_digit_0xd
	sidDigitsToFuns[0xe] = sid_digit_0xe
	sidDigitsToFuns[0xf] = sid_digit_0xf

	sidDigitsToPcs[0x0] = sidFuncToPc(sid_digit_0x0)
	sidDigitsToPcs[0x1] = sidFuncToPc(sid_digit_0x1)
	sidDigitsToPcs[0x2] = sidFuncToPc(sid_digit_0x2)
	sidDigitsToPcs[0x3] = sidFuncToPc(sid_digit_0x3)
	sidDigitsToPcs[0x4] = sidFuncToPc(sid_digit_0x4)
	sidDigitsToPcs[0x5] = sidFuncToPc(sid_digit_0x5)
	sidDigitsToPcs[0x6] = sidFuncToPc(sid_digit_0x6)
	sidDigitsToPcs[0x7] = sidFuncToPc(sid_digit_0x7)
	sidDigitsToPcs[0x8] = sidFuncToPc(sid_digit_0x8)
	sidDigitsToPcs[0x9] = sidFuncToPc(sid_digit_0x9)
	sidDigitsToPcs[0xa] = sidFuncToPc(sid_digit_0xa)
	sidDigitsToPcs[0xb] = sidFuncToPc(sid_digit_0xb)
	sidDigitsToPcs[0xc] = sidFuncToPc(sid_digit_0xc)
	sidDigitsToPcs[0xd] = sidFuncToPc(sid_digit_0xd)
	sidDigitsToPcs[0xe] = sidFuncToPc(sid_digit_0xe)
	sidDigitsToPcs[0xf] = sidFuncToPc(sid_digit_0xf)

	for key, val := range sidDigitsToPcs {
		sidPcsToDigits[val] = byte(key)
	}
}

func sidFuncToPc(fun sidFunc) (out uintptr) {
	fun(0, func() {
		buf := make([]uintptr, 1)

		/**
		Skip 3 PCs:

		- `runtime.Callers`
		- this closure
		- `sid_start` which is used by `fun`

		This gets us the PC of `fun`.
		*/
		out = buf[:runtime.Callers(3, buf)][0]
	})

	if out == 0 {
		panic(Errf(`internal error: unable to determine program counter for function %v`, fun))
	}
	return
}

func sidFun(rem uint64) sidFunc {
	if !(rem > 0) {
		// Big-endian encoding.
		return sid_start
	}
	return sidDigitsToFuns[rem&0xf]
}

//go:noinline
func sid_start(_ uint64, fun func()) { fun() }

//go:noinline
func sid_end(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0x0(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0x1(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0x2(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0x3(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0x4(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0x5(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0x6(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0x7(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0x8(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0x9(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0xa(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0xb(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0xc(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0xd(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0xe(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }

//go:noinline
func sid_digit_0xf(rem uint64, fun func()) { sidFun(rem)(rem>>sidDigitBits, fun) }
