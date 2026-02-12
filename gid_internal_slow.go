//go:build (386 || amd64 || arm || arm64 || riscv64 || s390x) && (!go1.23 || go1.27)

package gg

func gid() uint64 { return gidSlow() }
