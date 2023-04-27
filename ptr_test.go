package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestPtrInited(t *testing.T) {
	defer gtest.Catch(t)

	gtest.NotZero(gg.PtrInited((*string)(nil)))

	src := new(string)
	gtest.Eq(gg.PtrInited(src), src)
}

func TestPtrInit(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.PtrInit((**string)(nil)))

	var tar *string
	gtest.Eq(gg.PtrInit(&tar), tar)
	gtest.NotZero(tar)
}

func TestPtrClear(t *testing.T) {
	defer gtest.Catch(t)

	gtest.NotPanic(func() {
		gg.PtrClear((*string)(nil))
	})

	val := `str`
	gg.PtrClear(&val)
	gtest.Equal(val, ``)
}

func BenchmarkPtrClear(b *testing.B) {
	var val string

	for ind := 0; ind < b.N; ind++ {
		gg.PtrClear(&val)
		val = `str`
	}
}

func TestPtrGet(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.PtrGet((*string)(nil)), ``)
	gtest.Eq(gg.PtrGet(new(string)), ``)
	gtest.Eq(gg.PtrGet(gg.Ptr(`str`)), `str`)

	gtest.Eq(gg.PtrGet((*int)(nil)), 0)
	gtest.Eq(gg.PtrGet(new(int)), 0)
	gtest.Eq(gg.PtrGet(gg.Ptr(10)), 10)
}

func BenchmarkPtrGet_miss(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.PtrGet((*[]string)(nil)))
	}
}

func BenchmarkPtrGet_hit(b *testing.B) {
	ptr := gg.Ptr([]string{`one`, `two`})

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.PtrGet(ptr))
	}
}

func TestPtrSet(t *testing.T) {
	defer gtest.Catch(t)

	gtest.NotPanic(func() {
		gg.PtrSet((*string)(nil), ``)
		gg.PtrSet((*string)(nil), `str`)
	})

	var tar string

	gg.PtrSet(&tar, `one`)
	gtest.Eq(tar, `one`)

	gg.PtrSet(&tar, `two`)
	gtest.Eq(tar, `two`)
}

func TestPtrSetOpt(t *testing.T) {
	defer gtest.Catch(t)

	gtest.NotPanic(func() {
		gg.PtrSetOpt((*string)(nil), (*string)(nil))
		gg.PtrSetOpt(new(string), (*string)(nil))
		gg.PtrSetOpt((*string)(nil), new(string))
	})

	var tar string
	gg.PtrSetOpt(&tar, gg.Ptr(`one`))
	gtest.Eq(tar, `one`)

	gg.PtrSetOpt(&tar, gg.Ptr(`two`))
	gtest.Eq(tar, `two`)
}

func TestPtrPop(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src *string, exp string) {
		gtest.Eq(gg.PtrPop(src), exp)
	}

	test(nil, ``)
	test(gg.Ptr(``), ``)
	test(gg.Ptr(`val`), `val`)
}
