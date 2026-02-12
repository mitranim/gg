//go:build (386 || amd64 || arm || arm64 || riscv64 || s390x) && go1.23 && !go1.27

#include "textflag.h"

/*
Go keeps the current `*g` in a dedicated location which varies by architecture.
On 386 and amd64 this involves the pseudo-register `TLS` which maps to a real
register which has a pointer to thread-local storage which starts with `*g`.
On other architectures, Go defines the pseudo-register `g` which has `*g`;
it maps to one of the real registers, varying by architecture.

See also:

  - go/src/cmd/compile/abi-internal.md
  - go/src/cmd/asm/internal/arch/arch.go
  - go/src/runtime/go_tls.h
  - https://github.com/golang/go/blob/d19e377f6ea3b84e94d309894419f2995e7b56bd/src/cmd/internal/obj/x86/obj6.go#L73-L111

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
See also:

  - go/src/runtime/asm_386.s
  - https://go.dev/doc/asm#x86
*/
#ifdef GOARCH_386
  MOVL 0(TLS), AX
  MOVL AX, ret+0(FP)
#endif

/*
See also:

  - go/src/runtime/asm_amd64.s
  - https://go.dev/doc/asm#amd64
*/
#ifdef GOARCH_amd64
  MOVQ 0(TLS), AX
  MOVQ AX, ret+0(FP)
#endif

/*
See also:

  - go/src/runtime/asm_arm.s
  - https://go.dev/doc/asm#arm
*/
#ifdef GOARCH_arm
  MOVW g, ret+0(FP)
#endif

/*
See also:

  - go/src/runtime/asm_arm64.s
  - https://go.dev/doc/asm#arm64
*/
#ifdef GOARCH_arm64
  MOVD g, ret+0(FP)
#endif

// See also: `go/src/runtime/asm_riscv64.s`.
#ifdef GOARCH_riscv64
  MOV g, ret+0(FP)
#endif

// See also: `go/src/runtime/asm_s390x.s`.
#ifdef GOARCH_s390x
  MOVD g, ret+0(FP)
#endif

  RET
