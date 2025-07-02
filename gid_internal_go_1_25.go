//go:build (386 || amd64 || arm || arm64 || riscv64 || s390x) && go1.25 && !go1.26

package gg

import "sync/atomic"

func gid() uint64 { return getg().goid }

func getg() *g

// https://github.com/golang/go/blob/6e676ab2b809d46623acb5988248d95d1eb7939c/src/runtime/runtime2.go#L394
type g struct {
	stack struct {
		lo uintptr
		hi uintptr
	}
	stackguard0 uintptr
	stackguard1 uintptr
	_panic      uintptr
	_defer      uintptr
	m           uintptr
	sched       struct {
		sp   uintptr
		pc   uintptr
		g    uintptr
		ctxt uintptr
		lr   uintptr
		bp   uintptr
	}
	syscallsp    uintptr
	syscallpc    uintptr
	syscallbp    uintptr
	stktopsp     uintptr
	param        uintptr
	atomicstatus atomic.Uint32
	stackLock    uint32
	goid         uint64 // Goroutine ID.
}
