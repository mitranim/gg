//go:build (386 || amd64 || arm || arm64 || riscv64 || s390x) && (go1.23 || go1.24) && !go1.25

package gg

import "sync/atomic"

func gid() uint64 { return getg().goid }

func getg() *g

/*
Go 1.23:

	https://github.com/golang/go/blob/6885bad7dd86880be6929c02085e5c7a67ff2887/src/runtime/runtime2.go#L422

Go 1.24:

	https://github.com/golang/go/blob/3901409b5d0fb7c85a3e6730a59943cc93b2835c/src/runtime/runtime2.go#L396
*/
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
		ret  uintptr
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
