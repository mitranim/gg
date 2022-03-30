package gg_test

import (
	"fmt"
	r "reflect"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestType(t *testing.T) {
	defer gtest.Catch(t)

	testType[int]()
	testType[*int]()
	testType[**int]()

	testType[string]()
	testType[string]()
	testType[*string]()
	testType[**string]()

	testType[SomeModel]()
	testType[*SomeModel]()
	testType[**SomeModel]()

	testType[func()]()

	testTypeIface[any](r.TypeOf((*any)(nil)).Elem())
	testTypeIface[fmt.Stringer](r.TypeOf((*fmt.Stringer)(nil)).Elem())
}

func testType[A any]() {
	gtest.EqAny(gg.Type[A](), r.TypeOf(gg.Zero[A]()))
}

func testTypeIface[A any](exp r.Type) {
	gtest.EqAny(gg.Type[A](), exp)
}

func TestTypeOf(t *testing.T) {
	defer gtest.Catch(t)

	testTypeOf(int(0))
	testTypeOf(int(10))
	testTypeOf((*int)(nil))
	testTypeOf((**int)(nil))

	testTypeOf(string(``))
	testTypeOf(string(`str`))
	testTypeOf((*string)(nil))
	testTypeOf((**string)(nil))

	testTypeOf(SomeModel{})
	testTypeOf((*SomeModel)(nil))
	testTypeOf((**SomeModel)(nil))

	testTypeOf((func())(nil))

	testTypeOfIface(any(nil), r.TypeOf((*any)(nil)).Elem())
	testTypeOfIface(fmt.Stringer(nil), r.TypeOf((*fmt.Stringer)(nil)).Elem())
}

func testTypeOf[A any](src A) {
	gtest.Equal(gg.TypeOf(src), r.TypeOf(src))
}

func testTypeOfIface[A any](src A, exp r.Type) {
	gtest.Equal(gg.TypeOf(src), exp)
}

/*
This benchmark is defective. It fails to reproduce spurious escapes commonly
observed in code using this function.
*/
func Benchmark_reflect_TypeOf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(r.TypeOf(SomeModel{}))
	}
}

func BenchmarkTypeOf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.TypeOf(SomeModel{}))
	}
}

func BenchmarkType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Type[SomeModel]())
	}
}

func TestKindOfAny(t *testing.T) {
	defer gtest.Catch(t)

	// Difference from `KindOf`.
	testKindOfAny(any(nil), r.Invalid)
	testKindOfAny(fmt.Stringer(nil), r.Invalid)

	testKindOfAny(``, r.String)
	testKindOfAny((*string)(nil), r.Pointer)
	testKindOfAny(SomeModel{}, r.Struct)
	testKindOfAny((*SomeModel)(nil), r.Pointer)
	testKindOfAny([]string(nil), r.Slice)
	testKindOfAny((*[]string)(nil), r.Pointer)
	testKindOfAny((func())(nil), r.Func)
}

func testKindOfAny(src any, exp r.Kind) {
	gtest.Eq(gg.KindOfAny(src), exp)
}

func TestKindOf(t *testing.T) {
	defer gtest.Catch(t)

	// Difference from `KindOfAny`.
	testKindOf(any(nil), r.Interface)
	testKindOf(fmt.Stringer(nil), r.Interface)

	testKindOf(``, r.String)
	testKindOf((*string)(nil), r.Pointer)
	testKindOf(SomeModel{}, r.Struct)
	testKindOf((*SomeModel)(nil), r.Pointer)
	testKindOf([]string(nil), r.Slice)
	testKindOf((*[]string)(nil), r.Pointer)
	testKindOf((func())(nil), r.Func)
}

func testKindOf[A any](src A, exp r.Kind) {
	gtest.Eq(gg.KindOf(src), exp)
}

func BenchmarkKindOf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.KindOf(SomeModel{}))
	}
}

func BenchmarkAnyToString_miss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop2(gg.AnyToString(SomeModel{}))
	}
}

func BenchmarkAnyToString_hit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop2(gg.AnyToString([]byte(`hello world`)))
	}
}

func BenchmarkStructFields_Init(b *testing.B) {
	key := gg.Type[SomeModel]()

	for i := 0; i < b.N; i++ {
		var tar gg.StructFields
		tar.Init(key)
	}
}

func TestStructFieldCache(t *testing.T) {
	typ := gg.Type[SomeModel]()

	gtest.NotEmpty(gg.StructFieldCache.Get(typ))

	gtest.Equal(
		gg.StructFieldCache.Get(typ),
		gg.Times(typ.NumField(), typ.Field),
	)
}

func BenchmarkStructFieldCache(b *testing.B) {
	key := gg.Type[SomeModel]()

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.StructFieldCache.Get(key))
	}
}
