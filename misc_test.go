package gg_test

import (
	"context"
	"fmt"
	r "reflect"
	"testing"
	"time"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestIsZero(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(gg.IsZero(0))
	gtest.False(gg.IsZero(1))
	gtest.False(gg.IsZero(-1))

	gtest.True(gg.IsZero(``))
	gtest.False(gg.IsZero(` `))

	gtest.True(gg.IsZero([]string(nil)))
	gtest.False(gg.IsZero([]string{}))

	t.Run(`method`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.True(gg.IsZero[IsZeroAlwaysTrue](``))
		gtest.True(gg.IsZero[IsZeroAlwaysFalse](``))
		gtest.True(gg.IsZero[IsZeroAlwaysTrue](`str`))
		gtest.False(gg.IsZero[IsZeroAlwaysFalse](`str`))

		gtest.True(gg.IsZero(r.ValueOf(nil)))
		gtest.True(gg.IsZero(r.ValueOf(``)))
		gtest.False(gg.IsZero(r.ValueOf(`str`)))
	})

	t.Run(`time`, func(t *testing.T) {
		defer gtest.Catch(t)

		const minSec = 60
		const hourMin = 60
		const offsetHour = 1

		local := time.FixedZone(`local`, minSec*hourMin*offsetHour)

		testZero := func(src time.Time) {
			gtest.True(src.IsZero())
			gtest.True(gg.IsZero(src))
			gtest.Eq(gg.IsZero(src), src.IsZero())
		}

		testZero(time.Time{})
		testZero(time.Time{}.In(time.UTC))
		testZero(time.Time{}.In(local))

		gtest.Eq(time.Time{}, time.Time{}.In(time.UTC))
		gtest.NotEq(time.Time{}, time.Time{}.In(local))

		gtest.Eq(
			time.Date(1, 1, 1, offsetHour, 0, 0, 0, local),
			time.Time{}.In(local),
		)

		gtest.True(gg.IsZero(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)))
		gtest.True(gg.IsZero(time.Date(1, 1, 1, offsetHour, 0, 0, 0, local)))

		gtest.False(gg.IsZero(time.Date(1, 1, 1, 0, 0, 0, 1, time.UTC)))
		gtest.False(gg.IsZero(time.Date(1, 1, 0, 0, 0, 0, 1, local)))
	})
}

// This is a control. Our version should be significantly more efficient.
func Benchmark_is_zero_reflect_struct_zero(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(r.ValueOf(FatStruct{}).IsZero())
	}
}

// This is a control. Our version should be significantly more efficient.
func Benchmark_is_zero_reflect_struct_non_zero(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(r.ValueOf(FatStruct{Id: 10}).IsZero())
	}
}

func Benchmark_is_zero_IsZero_struct_zero(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(FatStruct{}))
	}
}

func Benchmark_is_zero_IsZero_struct_non_zero(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(FatStruct{Id: 10}))
	}
}

func Benchmark_is_zero_IsZero_time_Time_zero(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(time.Time{}))
	}
}

func Benchmark_is_zero_IsZero_time_Time_non_zero(b *testing.B) {
	inst := time.Now()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(inst))
	}
}

func Benchmark_is_zero_method_time_Time(b *testing.B) {
	inst := time.Now()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(inst.IsZero())
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

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(r.Zero(typ))
	}
}

func BenchmarkZero(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
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

	for ind := 0; ind < b.N; ind++ {
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
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Deref((*[]string)(nil)))
	}
}

func BenchmarkDeref_hit(b *testing.B) {
	ptr := gg.Ptr([]string{`one`, `two`})

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Deref(ptr))
	}
}

func TestPtrSet(t *testing.T) {
	defer gtest.Catch(t)

	gtest.NoPanic(func() {
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

	gtest.NoPanic(func() {
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

func TestPtrInited(t *testing.T) {
	defer gtest.Catch(t)

	gtest.NotZero(gg.PtrInited((*string)(nil)))

	src := new(string)
	gtest.Eq(gg.PtrInited(src), src)
}

func PtrInit(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.PtrInit((**string)(nil)))

	var tar *string
	gtest.Eq(gg.PtrInit(&tar), tar)
	gtest.NotZero(tar)
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

func TestMin2(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Min2(ComparerOf(10), ComparerOf(20)), ComparerOf(10))
	gtest.Eq(gg.Min2(ComparerOf(20), ComparerOf(10)), ComparerOf(10))
	gtest.Eq(gg.Min2(ComparerOf(-10), ComparerOf(10)), ComparerOf(-10))
	gtest.Eq(gg.Min2(ComparerOf(10), ComparerOf(-10)), ComparerOf(-10))
	gtest.Eq(gg.Min2(ComparerOf(0), ComparerOf(10)), ComparerOf(0))
	gtest.Eq(gg.Min2(ComparerOf(10), ComparerOf(0)), ComparerOf(0))
	gtest.Eq(gg.Min2(ComparerOf(-10), ComparerOf(0)), ComparerOf(-10))
	gtest.Eq(gg.Min2(ComparerOf(0), ComparerOf(-10)), ComparerOf(-10))
}

func TestMax2(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Max2(ComparerOf(10), ComparerOf(20)), ComparerOf(20))
	gtest.Eq(gg.Max2(ComparerOf(20), ComparerOf(10)), ComparerOf(20))
	gtest.Eq(gg.Max2(ComparerOf(-10), ComparerOf(10)), ComparerOf(10))
	gtest.Eq(gg.Max2(ComparerOf(10), ComparerOf(-10)), ComparerOf(10))
	gtest.Eq(gg.Max2(ComparerOf(0), ComparerOf(10)), ComparerOf(10))
	gtest.Eq(gg.Max2(ComparerOf(10), ComparerOf(0)), ComparerOf(10))
	gtest.Eq(gg.Max2(ComparerOf(-10), ComparerOf(0)), ComparerOf(0))
	gtest.Eq(gg.Max2(ComparerOf(0), ComparerOf(-10)), ComparerOf(0))
}

func BenchmarkMaxPrim2(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.MaxPrim2(10, 20))
	}
}

func BenchmarkMax2(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Max2(ComparerOf(10), ComparerOf(20)))
	}
}

func TestAnyIs(t *testing.T) {
	t.Run(`nil`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.False(gg.AnyIs[any](nil))
		gtest.False(gg.AnyIs[int](nil))
		gtest.False(gg.AnyIs[string](nil))
		gtest.False(gg.AnyIs[*string](nil))
		gtest.False(gg.AnyIs[fmt.Stringer](nil))
	})

	t.Run(`mismatch`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.False(gg.AnyIs[int](`str`))
		gtest.False(gg.AnyIs[string](10))
		gtest.False(gg.AnyIs[*string](`str`))
		gtest.False(gg.AnyIs[fmt.Stringer](`str`))
	})

	t.Run(`match`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.True(gg.AnyIs[any](`str`))
		gtest.True(gg.AnyIs[int](10))
		gtest.True(gg.AnyIs[string](`str`))
		gtest.True(gg.AnyIs[*string]((*string)(nil)))
		gtest.True(gg.AnyIs[*string](gg.Ptr(`str`)))
		gtest.True(gg.AnyIs[fmt.Stringer](gg.ErrStr(`str`)))
	})
}

func TestAnyAs(t *testing.T) {
	t.Run(`nil`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Zero(gg.AnyAs[any](nil))
		gtest.Zero(gg.AnyAs[int](nil))
		gtest.Zero(gg.AnyAs[string](nil))
		gtest.Zero(gg.AnyAs[*string](nil))
		gtest.Zero(gg.AnyAs[fmt.Stringer](nil))
	})

	t.Run(`mismatch`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Zero(gg.AnyAs[int](`str`))
		gtest.Zero(gg.AnyAs[string](10))
		gtest.Zero(gg.AnyAs[*string](`str`))
		gtest.Zero(gg.AnyAs[fmt.Stringer](`str`))
	})

	t.Run(`match`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(gg.AnyAs[any](`str`), `str`)
		gtest.Eq(gg.AnyAs[int](10), 10)
		gtest.Eq(gg.AnyAs[string](`str`), `str`)
		gtest.Equal(gg.AnyAs[*string](gg.Ptr(`str`)), gg.Ptr(`str`))

		gtest.Equal(
			gg.AnyAs[fmt.Stringer](gg.ErrStr(`str`)),
			fmt.Stringer(gg.ErrStr(`str`)),
		)
	})
}
func BenchmarkAnyAs_concrete_miss(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.AnyAs[gg.ErrStr](0))
	}
}

func BenchmarkAnyAs_iface_miss(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.AnyAs[fmt.Stringer](0))
	}
}

func BenchmarkAnyAs_concrete_hit(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.AnyAs[gg.ErrStr](gg.ErrStr(``)))
	}
}

func BenchmarkAnyAs_iface_hit(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.AnyAs[fmt.Stringer](gg.ErrStr(``)))
	}
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

func TestPlus2(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Plus2(10, 20), 30)
	gtest.Eq(gg.Plus2(`10`, `20`), `1020`)
}

func Benchmark_eq(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(ind == ind*2)
	}
}

func BenchmarkEq(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Eq(ind, ind*2))
	}
}

func Benchmark_reflect_DeepEqual_int(b *testing.B) {
	one := 123
	two := 123

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(r.DeepEqual(one, two))
	}
}

func BenchmarkEqual_int(b *testing.B) {
	one := 123
	two := 123

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Equal(one, two))
	}
}

func Benchmark_reflect_DeepEqual_bytes(b *testing.B) {
	one := []byte(`one`)
	two := []byte(`two`)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(r.DeepEqual(one, two))
	}
}

func BenchmarkEqual_bytes(b *testing.B) {
	one := []byte(`one`)
	two := []byte(`two`)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Equal(one, two))
	}
}

func Benchmark_reflect_DeepEqual_time_Time(b *testing.B) {
	one := time.Date(1234, 5, 23, 12, 34, 56, 0, time.UTC)
	two := time.Date(1234, 5, 23, 12, 34, 56, 0, time.UTC)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(r.DeepEqual(one, two))
	}
}

func BenchmarkEqual_time_Time(b *testing.B) {
	one := time.Date(1234, 5, 23, 12, 34, 56, 0, time.UTC)
	two := time.Date(1234, 5, 23, 12, 34, 56, 0, time.UTC)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Equal(one, two))
	}
}

func TestSliceIs(t *testing.T) {
	t.Run(`nil`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.SliceIs([]byte(nil), []byte(nil))
		gtest.SliceIs([]string(nil), []string(nil))
	})

	t.Run(`zerobase`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.SliceIs([]byte{}, []byte{})
		gtest.SliceIs([]string{}, []string{})
	})
}

func TestNotSliceIs(t *testing.T) {
	defer gtest.Catch(t)

	gtest.NotSliceIs([]byte(`str`), []byte(`str`))
	gtest.NotSliceIs([]string{`str`}, []string{`str`})
}

func TestSnapSlice(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.SnapSlice((*[]int)(nil)))

	tar := []int{10, 20}
	snap := gg.SnapSlice(&tar)

	gtest.Equal(snap, gg.SliceSnapshot[int]{&tar, 2})
	gtest.Eq(cap(tar), 2)

	tar = []int{10, 20, 30, 40}
	gtest.Equal(tar, []int{10, 20, 30, 40})
	gtest.Eq(cap(tar), 4)

	snap.Done()
	gtest.Equal(tar, []int{10, 20})
	gtest.Eq(cap(tar), 4)
}
