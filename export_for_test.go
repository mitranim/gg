package gg

import "unsafe"

/*
This file is only included in testing mode (because of `_test.go`)
and allows us to test some internals directly.
*/

func GidSlow() uint64 { return gidSlow() }

func GidWithOverride() uint64 { return getGid() }

func Glss() (out map[uint64]gls) {
	glss.val.Range(func(key uint64, val gls) bool {
		MapInit(&out)[key] = val
		return true
	})
	return
}

func GlsCopy(vals ...GlsVal) Gls { return glsCopy(vals...) }

func GlsKey[A any](dyn *DynVar[A]) glsKey { return glsKey(unsafe.Pointer(dyn)) }

type GlsInternal = gls

func DynVarDef[A any](src *DynVar[A]) func() A {
	defer Lock(&src.lock).Unlock()
	return src.def
}

func DynVarVal[A any](src *DynVar[A]) A {
	defer Lock(&src.lock).Unlock()
	return src.val
}
