//go:build (386 || amd64 || arm || arm64 || riscv64 || s390x) && !(go1.23 || go1.24 || go1.25)

package gg

func gid() uint64 { return gidSlow() }
