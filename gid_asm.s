//go:build (386 || amd64 || arm || arm64 || riscv64 || s390x) && (go1.23 || go1.24 || go1.25)

#include "textflag.h"

/*
Relevant note in `go/src/runtime/HACKING.md`:

> `getg()` alone returns the current `g`, but when executing on the system
> or signal stacks, this will return the current M's "g0" or "gsignal",
> respectively. This is usually not what you want.

...which shouldn't affect our hack because it's only called from user code,
which is supposed to only run on regular goroutines.
*/

// func getg() *g
TEXT ·getg(SB), NOSPLIT | NOFRAME, $0-0

/*
Relevant note in Go sources:

  https://github.com/golang/go/blob/d19e377f6ea3b84e94d309894419f2995e7b56bd/src/cmd/internal/obj/x86/obj6.go#L73-L111

Also see: https://go.dev/doc/asm#x86.
*/
#ifdef GOARCH_386
  MOVL 0(TLS), AX
  MOVL AX, ret+0(FP)
#endif

/*
Relevant note in Go sources:

  https://github.com/golang/go/blob/d19e377f6ea3b84e94d309894419f2995e7b56bd/src/cmd/internal/obj/x86/obj6.go#L73-L101

Also see: https://go.dev/doc/asm#amd64.
*/
#ifdef GOARCH_amd64
  MOVQ 0(TLS), AX
  MOVQ AX, ret+0(FP)
#endif

/*
Examples can be found in `go/src/runtime/asm_arm.s`.
Also see: https://go.dev/doc/asm#arm.
*/
#ifdef GOARCH_arm
  MOVW g, ret+0(FP)
#endif

/*
Examples can be found in `go/src/runtime/asm_arm64.s`.
Also see: https://go.dev/doc/asm#arm64.
*/
#ifdef GOARCH_arm64
  MOVD g, ret+0(FP)
#endif

// Examples can be found in `go/src/runtime/asm_riscv64.s`.
#ifdef GOARCH_riscv64
  MOV g, ret+0(FP)
#endif

// Examples can be found in `go/src/runtime/asm_s390x.s`.
#ifdef GOARCH_s390x
  MOVD g, ret+0(FP)
#endif

  RET
