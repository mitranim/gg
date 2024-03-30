package gg_test

import (
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

type NullString string

func (self NullString) IsZero() bool { return self == `` }

type Weekday struct{ NullString }

func Benchmark_is_zero_struct_wrapping_string_typedef_IsZero(b *testing.B) {
	val := Weekday{`wednesday`}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.IsZero(val))
	}
}

func Benchmark_is_zero_struct_wrapping_string_typedef_method_call(b *testing.B) {
	val := Weekday{`wednesday`}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(val.IsZero())
	}
}

func Benchmark_is_zero_struct_wrapping_string_typedef_inline(b *testing.B) {
	val := Weekday{`wednesday`}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(val == (Weekday{}))
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
	defer gtest.Catch(t)

	const msgNotIs = `expected given slice headers to be distinct, but they were identical`
	const msgIs = `expected given slice headers to be identical, but they were distinct`

	gtest.True(gg.SliceIs([]byte(nil), []byte(nil)))
	gtest.True(gg.SliceIs([]string(nil), []string(nil)))

	gtest.SliceIs([]byte(nil), []byte(nil))
	gtest.SliceIs([]string(nil), []string(nil))

	gtest.PanicStr(msgNotIs, func() { gtest.NotSliceIs([]byte(nil), []byte(nil)) })
	gtest.PanicStr(msgNotIs, func() { gtest.NotSliceIs([]string(nil), []string(nil)) })

	gtest.True(gg.SliceIs([]byte{}, []byte{}), `zerobase`)
	gtest.True(gg.SliceIs([]string{}, []string{}), `zerobase`)

	gtest.SliceIs([]byte{}, []byte{}, `zerobase`)
	gtest.SliceIs([]string{}, []string{}, `zerobase`)

	gtest.PanicStr(msgNotIs, func() { gtest.NotSliceIs([]byte{}, []byte{}) })
	gtest.PanicStr(msgNotIs, func() { gtest.NotSliceIs([]string{}, []string{}) })

	gtest.False(gg.SliceIs([]byte{0}, []byte{0}))
	gtest.False(gg.SliceIs([]string{``}, []string{``}))

	gtest.NotSliceIs([]byte{0}, []byte{0})
	gtest.NotSliceIs([]string{``}, []string{``})

	gtest.PanicStr(msgIs, func() { gtest.SliceIs([]byte{0}, []byte{0}) })
	gtest.PanicStr(msgIs, func() { gtest.SliceIs([]string{``}, []string{``}) })
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

func TestLess(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(gg.Less[Comparer[int]]())

	gtest.True(gg.Less(ComparerOf(0)))
	gtest.False(gg.Less(ComparerOf(0), ComparerOf(0)))

	gtest.True(gg.Less(ComparerOf(0), ComparerOf(10)))
	gtest.False(gg.Less(ComparerOf(0), ComparerOf(10), ComparerOf(10)))
	gtest.False(gg.Less(ComparerOf(0), ComparerOf(0), ComparerOf(10)))

	gtest.True(gg.Less(ComparerOf(0), ComparerOf(10), ComparerOf(20)))
	gtest.False(gg.Less(ComparerOf(10), ComparerOf(0), ComparerOf(20)))
	gtest.False(gg.Less(ComparerOf(20), ComparerOf(0), ComparerOf(10)))
}

func TestLessPrim(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(gg.LessPrim[int]())

	gtest.True(gg.LessPrim(0))
	gtest.False(gg.LessPrim(0, 0))

	gtest.True(gg.LessPrim(0, 10))
	gtest.False(gg.LessPrim(0, 10, 10))
	gtest.False(gg.LessPrim(0, 0, 10))

	gtest.True(gg.LessPrim(0, 10, 20))
	gtest.False(gg.LessPrim(10, 0, 20))
	gtest.False(gg.LessPrim(20, 0, 10))
}

func TestLessEq(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(gg.LessEq[Comparer[int]]())

	gtest.True(gg.LessEq(ComparerOf(0)))
	gtest.True(gg.LessEq(ComparerOf(0), ComparerOf(0)))

	gtest.True(gg.LessEq(ComparerOf(0), ComparerOf(10)))
	gtest.True(gg.LessEq(ComparerOf(0), ComparerOf(10), ComparerOf(10)))
	gtest.True(gg.LessEq(ComparerOf(0), ComparerOf(0), ComparerOf(10)))

	gtest.True(gg.LessEq(ComparerOf(0), ComparerOf(10), ComparerOf(20)))
	gtest.False(gg.LessEq(ComparerOf(10), ComparerOf(0), ComparerOf(20)))
	gtest.False(gg.LessEq(ComparerOf(20), ComparerOf(0), ComparerOf(10)))
}

func TestLessEqPrim(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(gg.LessEqPrim[int]())

	gtest.True(gg.LessEqPrim(0))
	gtest.True(gg.LessEqPrim(0, 0))

	gtest.True(gg.LessEqPrim(0, 10))
	gtest.True(gg.LessEqPrim(0, 10, 10))
	gtest.True(gg.LessEqPrim(0, 0, 10))

	gtest.True(gg.LessEqPrim(0, 10, 20))
	gtest.False(gg.LessEqPrim(10, 0, 20))
	gtest.False(gg.LessEqPrim(20, 0, 10))
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
