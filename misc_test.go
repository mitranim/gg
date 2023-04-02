package gg_test

import (
	"context"
	"fmt"
	"math"
	r "reflect"
	"testing"
	"time"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

var testArr32_0 = [4]byte{0x00, 0x2e, 0x7a, 0xb0}
var testArr32_1 = [4]byte{0x4b, 0x38, 0xa9, 0x65}
var testArr64_0 = [8]byte{0x00, 0x2e, 0x7a, 0xb0, 0xef, 0x3f, 0x44, 0x88}
var testArr64_1 = [8]byte{0x4b, 0x38, 0xa9, 0x65, 0x13, 0x0d, 0x46, 0x29}
var testArr128_0 = [16]byte{0x00, 0x2e, 0x7a, 0xb0, 0xef, 0x3f, 0x44, 0x88, 0x95, 0x88, 0xc1, 0xf1, 0x10, 0xeb, 0xc2, 0x08}
var testArr128_1 = [16]byte{0x4b, 0x38, 0xa9, 0x65, 0x13, 0x0d, 0x46, 0x29, 0xb7, 0x98, 0xd8, 0x69, 0x6f, 0xdf, 0xc7, 0xf2}
var testArr192_0 = [32]byte{0x00, 0x2e, 0x7a, 0xb0, 0xef, 0x3f, 0x44, 0x88, 0x95, 0x88, 0xc1, 0xf1, 0x10, 0xeb, 0xc2, 0x08, 0xe9, 0x68, 0x33, 0x30, 0xb3, 0xdb, 0x4b, 0x82}
var testArr192_1 = [32]byte{0x4b, 0x38, 0xa9, 0x65, 0x13, 0x0d, 0x46, 0x29, 0xb7, 0x98, 0xd8, 0x69, 0x6f, 0xdf, 0xc7, 0xf2, 0x41, 0x50, 0xd4, 0xc4, 0x4f, 0x45, 0x45, 0x13}
var testArr256_0 = [32]byte{0x00, 0x2e, 0x7a, 0xb0, 0xef, 0x3f, 0x44, 0x88, 0x95, 0x88, 0xc1, 0xf1, 0x10, 0xeb, 0xc2, 0x08, 0xe9, 0x68, 0x33, 0x30, 0xb3, 0xdb, 0x4b, 0x82, 0x8e, 0x1d, 0xb5, 0xe5, 0x1a, 0x90, 0xe4, 0xa2}
var testArr256_1 = [32]byte{0x4b, 0x38, 0xa9, 0x65, 0x13, 0x0d, 0x46, 0x29, 0xb7, 0x98, 0xd8, 0x69, 0x6f, 0xdf, 0xc7, 0xf2, 0x41, 0x50, 0xd4, 0xc4, 0x4f, 0x45, 0x45, 0x13, 0x81, 0xce, 0x33, 0xcb, 0x28, 0x13, 0x17, 0x32}

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

	t.Run(`struct_inner_field_fake_zero`, func(t *testing.T) {
		defer gtest.Catch(t)

		{
			var tar FatStruct
			tar.Name = `str`[:0]
			gtest.True(gg.IsZero(tar))
		}

		{
			var tar FatStructNonComparable
			tar.Name = `str`[:0]
			gtest.True(gg.IsZero(tar))
		}
	})
}

// This is a control. Our version should be significantly more efficient.
func Benchmark_is_zero_reflect_fat_struct_zero(b *testing.B) {
	var val FatStruct

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(r.ValueOf(val).IsZero())
	}
}

// This is a control. Our version should be significantly more efficient.
func Benchmark_is_zero_reflect_fat_struct_non_zero(b *testing.B) {
	val := FatStruct{Id: 10}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(r.ValueOf(val).IsZero())
	}
}

func Benchmark_is_zero_IsZero_fat_struct_zero(b *testing.B) {
	var val FatStruct

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(val))
	}
}

func Benchmark_is_zero_IsZero_fat_struct_non_zero(b *testing.B) {
	val := FatStruct{Id: 10}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(val))
	}
}

func Benchmark_is_zero_IsZero_fat_struct_non_comparable_zero(b *testing.B) {
	var val FatStructNonComparable

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(val))
	}
}

func Benchmark_is_zero_IsZero_fat_struct_non_comparable_non_zero(b *testing.B) {
	var val FatStructNonComparable
	val.Id = 10

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(val))
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

func Benchmark_is_zero_string_non_zero(b *testing.B) {
	val := `some_string`

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(val))
	}
}

func Benchmark_is_zero_int_non_zero(b *testing.B) {
	val := 123

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(val))
	}
}

func TestIsTrueZero(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(gg.IsTrueZero(false))
	gtest.False(gg.IsTrueZero(true))

	gtest.True(gg.IsTrueZero(byte(0)))
	gtest.False(gg.IsTrueZero(byte(1)))

	gtest.True(gg.IsTrueZero(0))
	gtest.False(gg.IsTrueZero(1))

	gtest.True(gg.IsTrueZero(``))
	gtest.False(gg.IsTrueZero(`str`[:0]))

	gtest.True(gg.IsTrueZero(FatStruct{}))
	gtest.False(gg.IsTrueZero(FatStruct{Id: 1}))
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
	var tar []string

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.AnyAs[fmt.Stringer](tar))
	}
}

func BenchmarkAnyAs_concrete_hit(b *testing.B) {
	var tar gg.Err

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.AnyAs[gg.Err](tar))
	}
}

func BenchmarkAnyAs_iface_hit(b *testing.B) {
	var tar gg.Err

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.AnyAs[error](tar))
	}
}

func BenchmarkAnyAs_analog_native_miss(b *testing.B) {
	var tar []string

	for ind := 0; ind < b.N; ind++ {
		impl, ok := any(tar).(fmt.Stringer)
		gg.Nop2(impl, ok)
	}
}

func BenchmarkAnyAs_analog_native_hit(b *testing.B) {
	var tar gg.Err

	for ind := 0; ind < b.N; ind++ {
		impl, ok := any(tar).(fmt.Stringer)
		gg.Nop2(impl, ok)
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

	rangeU := gg.Range[uint8]
	rangeS := gg.Range[int8]

	gtest.Equal(
		rangeU(0, math.MaxUint8),
		[]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254},
	)

	gtest.Equal(
		rangeS(math.MinInt8, math.MaxInt8),
		[]int8{-128, -127, -126, -125, -124, -123, -122, -121, -120, -119, -118, -117, -116, -115, -114, -113, -112, -111, -110, -109, -108, -107, -106, -105, -104, -103, -102, -101, -100, -99, -98, -97, -96, -95, -94, -93, -92, -91, -90, -89, -88, -87, -86, -85, -84, -83, -82, -81, -80, -79, -78, -77, -76, -75, -74, -73, -72, -71, -70, -69, -68, -67, -66, -65, -64, -63, -62, -61, -60, -59, -58, -57, -56, -55, -54, -53, -52, -51, -50, -49, -48, -47, -46, -45, -44, -43, -42, -41, -40, -39, -38, -37, -36, -35, -34, -33, -32, -31, -30, -29, -28, -27, -26, -25, -24, -23, -22, -21, -20, -19, -18, -17, -16, -15, -14, -13, -12, -11, -10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126},
	)

	gtest.Equal(rangeU(math.MaxUint8, 0), nil)
	gtest.Equal(rangeU(math.MaxUint8, 1), nil)
	gtest.Equal(rangeU(math.MaxUint8, math.MaxUint8), nil)
	gtest.Equal(rangeU(math.MaxUint8/2, 1), nil)
	gtest.Equal(rangeU(math.MaxUint8/2, math.MaxUint8/4), nil)
	gtest.Equal(rangeU(0, 0), nil)
	gtest.Equal(rangeU(1, 1), nil)
	gtest.Equal(rangeU(2, 2), nil)

	gtest.Equal(rangeS(math.MaxInt8, 0), nil)
	gtest.Equal(rangeS(math.MaxInt8, 1), nil)
	gtest.Equal(rangeS(math.MaxInt8, math.MaxInt8), nil)
	gtest.Equal(rangeS(0, math.MinInt8), nil)
	gtest.Equal(rangeS(1, math.MinInt8), nil)
	gtest.Equal(rangeS(math.MinInt8, math.MinInt8), nil)
	gtest.Equal(rangeS(math.MaxInt8, math.MinInt8), nil)
	gtest.Equal(rangeS(math.MaxInt8/2, math.MinInt8/2), nil)
	gtest.Equal(rangeS(-2, -2), nil)
	gtest.Equal(rangeS(-1, -1), nil)
	gtest.Equal(rangeS(0, 0), nil)
	gtest.Equal(rangeS(1, 1), nil)
	gtest.Equal(rangeS(2, 2), nil)

	gtest.Equal(rangeS(-3, -2), []int8{-3})
	gtest.Equal(rangeS(-3, -1), []int8{-3, -2})
	gtest.Equal(rangeS(-3, 0), []int8{-3, -2, -1})
	gtest.Equal(rangeS(-3, 1), []int8{-3, -2, -1, 0})
	gtest.Equal(rangeS(-3, 2), []int8{-3, -2, -1, 0, 1})
	gtest.Equal(rangeS(-3, 3), []int8{-3, -2, -1, 0, 1, 2})

	gtest.Equal(rangeS(0, 1), []int8{0})
	gtest.Equal(rangeS(0, 2), []int8{0, 1})
	gtest.Equal(rangeS(0, 3), []int8{0, 1, 2})

	gtest.Equal(rangeS(3, 4), []int8{3})
	gtest.Equal(rangeS(3, 5), []int8{3, 4})
	gtest.Equal(rangeS(3, 6), []int8{3, 4, 5})

	gtest.Equal(gg.Range(math.MinInt, math.MinInt+1), []int{math.MinInt})
	gtest.Equal(gg.Range(math.MaxInt-1, math.MaxInt), []int{math.MaxInt - 1})

	gtest.PanicStr(`unable to safely convert uint 18446744073709551614 to int -2`, func() {
		gg.Range[uint](math.MaxUint-1, math.MaxUint)
	})
}

func TestRangeIncl(t *testing.T) {
	defer gtest.Catch(t)

	rangeU := gg.RangeIncl[uint8]
	rangeS := gg.RangeIncl[int8]

	gtest.Equal(
		rangeU(0, math.MaxUint8),
		[]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255},
	)

	gtest.Equal(
		rangeS(math.MinInt8, math.MaxInt8),
		[]int8{-128, -127, -126, -125, -124, -123, -122, -121, -120, -119, -118, -117, -116, -115, -114, -113, -112, -111, -110, -109, -108, -107, -106, -105, -104, -103, -102, -101, -100, -99, -98, -97, -96, -95, -94, -93, -92, -91, -90, -89, -88, -87, -86, -85, -84, -83, -82, -81, -80, -79, -78, -77, -76, -75, -74, -73, -72, -71, -70, -69, -68, -67, -66, -65, -64, -63, -62, -61, -60, -59, -58, -57, -56, -55, -54, -53, -52, -51, -50, -49, -48, -47, -46, -45, -44, -43, -42, -41, -40, -39, -38, -37, -36, -35, -34, -33, -32, -31, -30, -29, -28, -27, -26, -25, -24, -23, -22, -21, -20, -19, -18, -17, -16, -15, -14, -13, -12, -11, -10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127},
	)

	gtest.Equal(rangeU(math.MaxUint8, 0), nil)
	gtest.Equal(rangeU(math.MaxUint8, 1), nil)
	gtest.Equal(rangeU(math.MaxUint8/2, 1), nil)
	gtest.Equal(rangeU(math.MaxUint8/2, math.MaxUint8/4), nil)

	gtest.Equal(rangeS(math.MaxInt8, 0), nil)
	gtest.Equal(rangeS(math.MaxInt8, 1), nil)
	gtest.Equal(rangeS(0, math.MinInt8), nil)
	gtest.Equal(rangeS(1, math.MinInt8), nil)
	gtest.Equal(rangeS(math.MaxInt8, math.MinInt8), nil)
	gtest.Equal(rangeS(math.MaxInt8/2, math.MinInt8/2), nil)

	gtest.Equal(rangeU(math.MaxUint8, math.MaxUint8), []uint8{math.MaxUint8})
	gtest.Equal(rangeS(math.MaxInt8, math.MaxInt8), []int8{math.MaxInt8})
	gtest.Equal(rangeS(math.MinInt8, math.MinInt8), []int8{math.MinInt8})

	gtest.Equal(rangeU(0, 0), []uint8{0})
	gtest.Equal(rangeU(1, 1), []uint8{1})
	gtest.Equal(rangeU(2, 2), []uint8{2})

	gtest.Equal(rangeS(-2, -2), []int8{-2})
	gtest.Equal(rangeS(-1, -1), []int8{-1})
	gtest.Equal(rangeS(0, 0), []int8{0})
	gtest.Equal(rangeS(1, 1), []int8{1})
	gtest.Equal(rangeS(2, 2), []int8{2})

	gtest.Equal(rangeS(-3, -2), []int8{-3, -2})
	gtest.Equal(rangeS(-3, -1), []int8{-3, -2, -1})
	gtest.Equal(rangeS(-3, 0), []int8{-3, -2, -1, 0})
	gtest.Equal(rangeS(-3, 1), []int8{-3, -2, -1, 0, 1})
	gtest.Equal(rangeS(-3, 2), []int8{-3, -2, -1, 0, 1, 2})
	gtest.Equal(rangeS(-3, 3), []int8{-3, -2, -1, 0, 1, 2, 3})

	gtest.Equal(rangeS(0, 1), []int8{0, 1})
	gtest.Equal(rangeS(0, 2), []int8{0, 1, 2})
	gtest.Equal(rangeS(0, 3), []int8{0, 1, 2, 3})

	gtest.Equal(rangeS(3, 4), []int8{3, 4})
	gtest.Equal(rangeS(3, 5), []int8{3, 4, 5})
	gtest.Equal(rangeS(3, 6), []int8{3, 4, 5, 6})

	gtest.Equal(gg.RangeIncl(math.MinInt, math.MinInt+1), []int{math.MinInt, math.MinInt + 1})
	gtest.Equal(gg.RangeIncl(math.MaxInt-1, math.MaxInt), []int{math.MaxInt - 1, math.MaxInt})

	gtest.PanicStr(`unable to safely convert uint 18446744073709551614 to int -2`, func() {
		gg.RangeIncl[uint](math.MaxUint-1, math.MaxUint)
	})
}

func TestSpan(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(gg.Span(0), nil)
	gtest.Equal(gg.Span(1), []int{0})
	gtest.Equal(gg.Span(2), []int{0, 1})
	gtest.Equal(gg.Span(3), []int{0, 1, 2})
}

func TestPlus2(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.Plus2(10, 20), 30)
	gtest.Eq(gg.Plus2(`10`, `20`), `1020`)
}

func Benchmark_eq_operator(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(ind == ind*2)
	}
}

func BenchmarkEq_int(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Eq(ind, ind*2))
	}
}

func BenchmarkEq_array_64(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Eq(testArr64_0, testArr64_1))
	}
}

func BenchmarkEq_array_128(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Eq(testArr128_0, testArr128_1))
	}
}

func BenchmarkEq_array_256(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Eq(testArr256_0, testArr256_1))
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

func BenchmarkSliceIs(b *testing.B) {
	defer gtest.Catch(b)
	src := nonEmptyByteSlice()
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.SliceIs(src, src)
	}
}

//go:noinline
func nonEmptyByteSlice() []byte { return []byte(`one`) }

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

	gtest.Equal(snap, gg.SliceSnapshot[[]int, int]{&tar, 2})
	gtest.Eq(cap(tar), 2)

	tar = []int{10, 20, 30, 40}
	gtest.Equal(tar, []int{10, 20, 30, 40})
	gtest.Eq(cap(tar), 4)

	snap.Done()
	gtest.Equal(tar, []int{10, 20})
	gtest.Eq(cap(tar), 4)
}

func TestIs(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(gg.Is[int](10, 10))
	gtest.False(gg.Is[int](10, 20))

	gtest.True(gg.Is[any](10, 10))
	gtest.False(gg.Is[any](10, 20))

	var one int = 10
	var two int = 10
	gtest.True(gg.Is(one, one))
	gtest.True(gg.Is(one, two))
	gtest.True(gg.Is(&one, &one))
	gtest.False(gg.Is(&one, &two))

	gtest.True(gg.Is([1]byte{}, [1]byte{}))
	gtest.True(gg.Is([2]byte{}, [2]byte{}))
	gtest.True(gg.Is([4]byte{}, [4]byte{}))
	gtest.True(gg.Is([8]byte{}, [8]byte{}))
	gtest.True(gg.Is([16]byte{}, [16]byte{}))
	gtest.True(gg.Is([32]byte{}, [32]byte{}))
	gtest.True(gg.Is([64]byte{}, [64]byte{}))

	gtest.True(gg.Is(testArr32_0, testArr32_0))
	gtest.False(gg.Is(testArr32_0, testArr32_1))

	gtest.True(gg.Is(testArr64_0, testArr64_0))
	gtest.False(gg.Is(testArr64_0, testArr64_1))

	gtest.True(gg.Is(testArr128_0, testArr128_0))
	gtest.False(gg.Is(testArr128_0, testArr128_1))

	gtest.True(gg.Is(testArr256_0, testArr256_0))
	gtest.False(gg.Is(testArr256_0, testArr256_1))

	// Nil slices must be identical regardless of the element type.
	gtest.True(gg.Is([]struct{}(nil), []struct{}(nil)))
	gtest.True(gg.Is([]int(nil), []int(nil)))
	gtest.True(gg.Is([]string(nil), []string(nil)))

	// Slices of zero-sized types and empty slices of non-zero-sized types are
	// backed by the same "zerobase" pointer, which makes them identical.
	// This may vary between Go implementations and versions.
	gtest.True(gg.Is([]struct{}{}, []struct{}{}))
	gtest.True(gg.Is(make([]struct{}, 128), make([]struct{}, 128)))
	gtest.True(gg.Is([]int{}, []int{}))
	gtest.True(gg.Is([]string{}, []string{}))

	// Non-empty slices of non-zero-sized types must always be distinct.
	gtest.False(gg.Is([]int{0}, []int{0}))
	gtest.False(gg.Is([]string{``}, []string{``}))

	// Even though strings are reference types, string constants may be identical
	// when equal. The same may occur for interface values created from string
	// constants. This may vary between Go implementations and versions.
	gtest.True(gg.Is(``, ``))
	gtest.True(gg.Is(`one`, `one`))
	gtest.True(gg.Is[any](`one`, `one`))
	gtest.False(gg.Is(`one`, `two`))
	gtest.False(gg.Is[any](`one`, `two`))

	// However, slicing string constants may produce strings which are equal but
	// not identical. This may vary between Go implementations and versions.
	gtest.True(`123_one`[3:] == `456_one`[3:])
	gtest.False(gg.Is(`123_one`[3:], `456_one`[3:]))
}

func BenchmarkIs_ifaces(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Is(any(nil), any(nil)))
	}
}

func BenchmarkIs_slices(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Is([]byte(nil), []byte(nil)))
	}
}

func BenchmarkIs_array_32(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Is(testArr32_0, testArr32_1))
	}
}

func BenchmarkIs_array_64(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Is(testArr64_0, testArr64_1))
	}
}

func BenchmarkIs_array_128(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Is(testArr128_0, testArr128_1))
	}
}

func BenchmarkIs_array_192(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Is(testArr192_0, testArr192_1))
	}
}

func BenchmarkIs_array_256(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Is(testArr256_0, testArr256_1))
	}
}

func BenchmarkIs_struct_256(b *testing.B) {
	type Type struct{ _, _, _, _ uint64 }

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Is(Type{}, Type{}))
	}
}

func BenchmarkIs_struct_512(b *testing.B) {
	type Type struct{ _, _, _, _, _, _, _, _ uint64 }

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Is(Type{}, Type{}))
	}
}

func BenchmarkIs_struct_fat(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Is(FatStruct{}, FatStruct{}))
	}
}
