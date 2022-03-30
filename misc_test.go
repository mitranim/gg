package gg_test

import (
	"context"
	r "reflect"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestIsZero(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(gg.IsZero(0))
	gtest.True(!gg.IsZero(1))
	gtest.True(!gg.IsZero(-1))

	gtest.True(gg.IsZero(``))
	gtest.True(!gg.IsZero(` `))

	gtest.True(gg.IsZero([]string(nil)))
	gtest.True(!gg.IsZero([]string{}))
}

func Benchmark_reflect_IsZero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(r.ValueOf(SomeModel{}).IsZero())
	}
}

func BenchmarkIsZero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.IsZero(SomeModel{}))
	}
}

func TestZero(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.Zero[int](), 0)
	gtest.Equal(gg.Zero[string](), ``)
	gtest.Equal(gg.Zero[struct{}](), struct{}{})
	gtest.Equal(gg.Zero[[]string](), nil)
	gtest.Equal(gg.Zero[func()](), nil)
}

func Benchmark_reflect_Zero(b *testing.B) {
	typ := r.TypeOf(SomeModel{})

	for i := 0; i < b.N; i++ {
		gg.Nop1(r.Zero(typ))
	}
}

func BenchmarkZero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Zero[SomeModel]())
	}
}

func TestClear(t *testing.T) {
	defer gtest.Catch(t)

	gtest.NoPanic(func() {
		gg.Clear((*string)(nil))
	})

	val := `str`
	gg.Clear(&val)
	gtest.Equal(val, ``)
}

func BenchmarkClear(b *testing.B) {
	var val string

	for i := 0; i < b.N; i++ {
		gg.Clear(&val)
		val = `str`
	}
}

func TestDeref(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Deref((*string)(nil)), ``)
	gtest.Eq(gg.Deref(new(string)), ``)
	gtest.Eq(gg.Deref(gg.Ptr(`str`)), `str`)

	gtest.Eq(gg.Deref((*int)(nil)), 0)
	gtest.Eq(gg.Deref(new(int)), 0)
	gtest.Eq(gg.Deref(gg.Ptr(10)), 10)
}

func BenchmarkDeref_miss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Deref((*[]string)(nil)))
	}
}

func BenchmarkDeref_hit(b *testing.B) {
	ptr := gg.Ptr([]string{`one`, `two`})

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Deref(ptr))
	}
}

func TestMinPrim2(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.MinPrim2(10, 20), 10)
	gtest.Eq(gg.MinPrim2(20, 10), 10)
	gtest.Eq(gg.MinPrim2(-10, 10), -10)
	gtest.Eq(gg.MinPrim2(10, -10), -10)
	gtest.Eq(gg.MinPrim2(0, 10), 0)
	gtest.Eq(gg.MinPrim2(10, 0), 0)
	gtest.Eq(gg.MinPrim2(-10, 0), -10)
	gtest.Eq(gg.MinPrim2(0, -10), -10)
}

func TestMaxPrim2(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.MaxPrim2(10, 20), 20)
	gtest.Eq(gg.MaxPrim2(20, 10), 20)
	gtest.Eq(gg.MaxPrim2(-10, 10), 10)
	gtest.Eq(gg.MaxPrim2(10, -10), 10)
	gtest.Eq(gg.MaxPrim2(0, 10), 10)
	gtest.Eq(gg.MaxPrim2(10, 0), 10)
	gtest.Eq(gg.MaxPrim2(-10, 0), 0)
	gtest.Eq(gg.MaxPrim2(0, -10), 0)
}

func TestCtxSet(t *testing.T) {
	defer gtest.Catch(t)

	ctx := context.Background()

	gtest.Zero(ctx.Value((*string)(nil)))

	gtest.Equal(
		gg.CtxSet(ctx, `str`).Value((*string)(nil)),
		any(`str`),
	)
}

func TestCtxGet(t *testing.T) {
	defer gtest.Catch(t)

	ctx := context.Background()

	gtest.Zero(gg.CtxGet[string](ctx))

	gtest.Eq(
		gg.CtxGet[string](gg.CtxSet(ctx, `str`)),
		`str`,
	)

	gtest.Eq(
		gg.CtxGet[string](context.WithValue(gg.CtxSet(ctx, `str`), (*int)(nil), 10)),
		`str`,
	)
}

func TestRange(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(`out of range`, func() {
		gtest.Equal(gg.Range(3, 2), []int{})
		gtest.Equal(gg.Range(-2, -3), []int{})
	})

	gtest.Equal(gg.Range(-2, -2), []int{})
	gtest.Equal(gg.Range(-1, -1), []int{})
	gtest.Equal(gg.Range(0, 0), []int{})
	gtest.Equal(gg.Range(1, 1), []int{})
	gtest.Equal(gg.Range(2, 2), []int{})

	gtest.Equal(gg.Range(-3, -2), []int{-3})
	gtest.Equal(gg.Range(-3, -1), []int{-3, -2})
	gtest.Equal(gg.Range(-3, 0), []int{-3, -2, -1})
	gtest.Equal(gg.Range(-3, 1), []int{-3, -2, -1, 0})
	gtest.Equal(gg.Range(-3, 2), []int{-3, -2, -1, 0, 1})
	gtest.Equal(gg.Range(-3, 3), []int{-3, -2, -1, 0, 1, 2})

	gtest.Equal(gg.Range(0, 1), []int{0})
	gtest.Equal(gg.Range(0, 2), []int{0, 1})
	gtest.Equal(gg.Range(0, 3), []int{0, 1, 2})

	gtest.Equal(gg.Range(3, 4), []int{3})
	gtest.Equal(gg.Range(3, 5), []int{3, 4})
	gtest.Equal(gg.Range(3, 6), []int{3, 4, 5})
}

func TestSpan(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.Span(0), []int{})
	gtest.Equal(gg.Span(1), []int{0})
	gtest.Equal(gg.Span(2), []int{0, 1})
	gtest.Equal(gg.Span(3), []int{0, 1, 2})
}

func Benchmark_eq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(i == i*2)
	}
}

func BenchmarkEq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Eq(i, i*2))
	}
}

func Benchmark_reflect_DeepEqual(b *testing.B) {
	one := []byte(`one`)
	two := []byte(`two`)

	for i := 0; i < b.N; i++ {
		gg.Nop1(r.DeepEqual(one, two))
	}
}

func BenchmarkEqual(b *testing.B) {
	one := []byte(`one`)
	two := []byte(`two`)

	for i := 0; i < b.N; i++ {
		gg.Nop1(gg.Equal(one, two))
	}
}
