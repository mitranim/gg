package gg_test

import (
	"fmt"
	r "reflect"
	"testing"
	u "unsafe"

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
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(r.TypeOf(SomeModel{}))
	}
}

func BenchmarkTypeOf(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.TypeOf(SomeModel{}))
	}
}

func BenchmarkType(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
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
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.KindOf(SomeModel{}))
	}
}

func BenchmarkAnyToString_miss(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop2(gg.AnyToString(SomeModel{}))
	}
}

func BenchmarkAnyToString_hit(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop2(gg.AnyToString([]byte(`hello world`)))
	}
}

func BenchmarkStructFields_Init(b *testing.B) {
	key := gg.Type[SomeModel]()

	for ind := 0; ind < b.N; ind++ {
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

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.StructFieldCache.Get(key))
	}
}

func TestIsIndirect(t *testing.T) {
	defer gtest.Catch(t)

	gtest.False(gg.IsIndirect(gg.Type[bool]()))
	gtest.False(gg.IsIndirect(gg.Type[int]()))
	gtest.False(gg.IsIndirect(gg.Type[string]()))
	gtest.False(gg.IsIndirect(gg.Type[[0]bool]()))
	gtest.False(gg.IsIndirect(gg.Type[[1]bool]()))
	gtest.False(gg.IsIndirect(gg.Type[[0]*string]()))
	gtest.False(gg.IsIndirect(gg.Type[StructDirect]()))
	gtest.False(gg.IsIndirect(gg.Type[func()]()))
	gtest.False(gg.IsIndirect(gg.Type[func() bool]()))
	gtest.False(gg.IsIndirect(gg.Type[func() *string]()))
	gtest.False(gg.IsIndirect(gg.Type[chan bool]()))
	gtest.False(gg.IsIndirect(gg.Type[chan *string]()))

	gtest.True(gg.IsIndirect(gg.Type[any]()))
	gtest.True(gg.IsIndirect(gg.Type[fmt.Stringer]()))
	gtest.True(gg.IsIndirect(gg.Type[[1]*string]()))
	gtest.True(gg.IsIndirect(gg.Type[[]byte]()))
	gtest.True(gg.IsIndirect(gg.Type[[]string]()))
	gtest.True(gg.IsIndirect(gg.Type[[]*string]()))
	gtest.True(gg.IsIndirect(gg.Type[*bool]()))
	gtest.True(gg.IsIndirect(gg.Type[*StructDirect]()))
	gtest.True(gg.IsIndirect(gg.Type[StructIndirect]()))
	gtest.True(gg.IsIndirect(gg.Type[*StructIndirect]()))
	gtest.True(gg.IsIndirect(gg.Type[map[bool]bool]()))
}

func TestCloneDeep(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`direct`, func(t *testing.T) {
		defer gtest.Catch(t)

		testCloneDeepSame(true)
		testCloneDeepSame(10)
		testCloneDeepSame(`str`)
		testCloneDeepSame([0]string{})
		testCloneDeepSame([2]string{`one`, `two`})
		testCloneDeepSame([0]*string{})

		// Private fields are ignored.
		testCloneDeepSame(StructDirect{
			Public0: 10,
			Public1: `one`,
			private: gg.Ptr(`two`),
		})
	})

	t.Run(`pointer`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.CloneDeep((*string)(nil)), (*string)(nil))

		{
			src := gg.Ptr(`one`)
			out := gg.CloneDeep(src)

			gtest.NotEq(out, src)

			*src = `two`
			gtest.Equal(src, gg.Ptr(`two`))
			gtest.Equal(out, gg.Ptr(`one`))
		}

		{
			src := gg.Ptr(gg.Ptr(`one`))
			out := gg.CloneDeep(src)

			gtest.NotEq(out, src)
			gtest.NotEq(*out, *src)

			**src = `two`
			gtest.Equal(src, gg.Ptr(gg.Ptr(`two`)))
			gtest.Equal(out, gg.Ptr(gg.Ptr(`one`)))
		}
	})

	t.Run(`slice`, func(t *testing.T) {
		defer gtest.Catch(t)

		testCloneDeepSameSlice([]string(nil))
		testCloneDeepSameSlice([]string{})

		// Slices with zero length but non-zero capacity must still be cloned.
		testCloneDeepDifferentSlice([]string{`one`, `two`}[:0])
		testCloneDeepDifferentSlice([]*string{gg.Ptr(`one`), gg.Ptr(`two`)}[:0])

		{
			src := []string{`one`, `two`}
			out := gg.CloneDeep(src)

			testCloneDeepDifferentSlice(src)

			src[0] = `three`
			gtest.Equal(src, []string{`three`, `two`})
			gtest.Equal(out, []string{`one`, `two`})
		}

		{
			src := []*string{gg.Ptr(`one`), gg.Ptr(`two`)}
			out := gg.CloneDeep(src)

			testCloneDeepDifferentSlice(src)

			*src[0] = `three`
			gtest.Equal(src, []*string{gg.Ptr(`three`), gg.Ptr(`two`)})
			gtest.Equal(out, []*string{gg.Ptr(`one`), gg.Ptr(`two`)})
		}
	})

	t.Run(`slice_of_struct_pointers`, func(t *testing.T) {
		defer gtest.Catch(t)

		one := SomeModel{Id: `10`}
		two := SomeModel{Id: `20`}
		src := []*SomeModel{&one, &two}
		out := gg.CloneDeep(src)

		gtest.Equal(out, src)

		one.Id = `30`
		two.Id = `40`
		src = append(src, &SomeModel{Id: `50`})

		gtest.Equal(
			src,
			[]*SomeModel{
				&SomeModel{Id: `30`},
				&SomeModel{Id: `40`},
				&SomeModel{Id: `50`},
			},
		)

		gtest.Equal(
			out,
			[]*SomeModel{
				&SomeModel{Id: `10`},
				&SomeModel{Id: `20`},
			},
		)
	})

	t.Run(`inner_interface`, func(t *testing.T) {
		defer gtest.Catch(t)

		type Type struct{ Val fmt.Stringer }

		srcInner := gg.ErrStr(`one`)
		src := Type{&srcInner}
		out := gg.CloneDeep(src)

		gtest.Equal(src, Type{gg.Ptr(gg.ErrStr(`one`))})
		gtest.Equal(out, Type{gg.Ptr(gg.ErrStr(`one`))})

		srcInner = `two`

		gtest.Equal(src, Type{gg.Ptr(gg.ErrStr(`two`))})
		gtest.Equal(out, Type{gg.Ptr(gg.ErrStr(`one`))})
	})
}

func testCloneDeepSame[A comparable](src A) {
	gtest.Eq(gg.CloneDeep(src), src)
}

func testCloneDeepSameSlice[A any](src []A) {
	gtest.Equal(gg.CloneDeep(src), src)

	gtest.Eq(gg.SliceDat(gg.Clone(src)), gg.SliceDat(src))
	gtest.Eq(gg.SliceDat(gg.CloneDeep(src)), gg.SliceDat(src))

	gtest.SliceIs(gg.Clone(src), src)
	gtest.SliceIs(gg.CloneDeep(src), src)
}

/*
Note: this doesn't verify the deep cloning of slice elements, which must be
checked separately.
*/
func testCloneDeepDifferentSlice[A any](src []A) {
	gtest.Equal(gg.CloneDeep(src), src)
	gtest.Equal(gg.CloneDeep(src), gg.Clone(src))

	gtest.NotEq(gg.SliceDat(gg.Clone(src)), gg.SliceDat(src))
	gtest.NotEq(gg.SliceDat(gg.CloneDeep(src)), gg.SliceDat(src))

	gtest.NotSliceIs(gg.Clone(src), src)
	gtest.NotSliceIs(gg.CloneDeep(src), src)
}

func Benchmark_clone_direct_CloneDeep(b *testing.B) {
	src := [8]SomeModel{{Id: `10`}, {Id: `20`}, {Id: `30`}}
	gtest.Equal(gg.CloneDeep(src), src)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.CloneDeep(src))
	}
}

func Benchmark_clone_direct_native(b *testing.B) {
	src := [8]SomeModel{{Id: `10`}, {Id: `20`}, {Id: `30`}}

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(esc(src))
	}
}

func Benchmark_clone_slice_CloneDeep(b *testing.B) {
	src := []SomeModel{{Id: `10`}, {Id: `20`}, {Id: `30`}}
	gtest.Equal(gg.CloneDeep(src), src)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.CloneDeep(src))
	}
}

func Benchmark_clone_slice_Clone(b *testing.B) {
	src := []SomeModel{{Id: `10`}, {Id: `20`}, {Id: `30`}}
	gtest.Equal(gg.Clone(src), src)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Clone(src))
	}
}

func Benchmark_clone_map_CloneDeep(b *testing.B) {
	src := gg.Index(
		[]SomeModel{{Id: `10`}, {Id: `20`}, {Id: `30`}},
		gg.ValidPk[SomeKey, SomeModel],
	)
	gtest.Equal(gg.CloneDeep(src), src)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.CloneDeep(src))
	}
}

func Benchmark_clone_map_MapClone(b *testing.B) {
	src := gg.Index(
		[]SomeModel{{Id: `10`}, {Id: `20`}, {Id: `30`}},
		gg.ValidPk[SomeKey, SomeModel],
	)
	gtest.Equal(gg.MapClone(src), src)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.MapClone(src))
	}
}

func TestStructDeepPublicFieldCache(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.StructDeepPublicFieldCache.Get(gg.Type[struct{}]()))

	gtest.Equal(
		gg.StructDeepPublicFieldCache.Get(gg.Type[StructDirect]()),
		gg.StructDeepPublicFields{
			{
				Name:   `Public0`,
				Type:   gg.Type[int](),
				Offset: 0,
				Index:  []int{0},
			},
			{
				Name:   `Public1`,
				Type:   gg.Type[string](),
				Offset: u.Offsetof(StructDirect{}.Public1),
				Index:  []int{1},
			},
		},
	)

	gtest.Equal(
		gg.StructDeepPublicFieldCache.Get(gg.Type[Outer]()),
		gg.StructDeepPublicFields{
			{
				Name:  `OuterId`,
				Type:  gg.Type[int](),
				Index: []int{0}},
			{
				Name:   `OuterName`,
				Type:   gg.Type[string](),
				Offset: u.Offsetof(Outer{}.OuterName),
				Index:  []int{1}},
			{
				Name:   `EmbedId`,
				Type:   gg.Type[int](),
				Offset: u.Offsetof(Outer{}.Embed) + u.Offsetof(Embed{}.EmbedId),
				Index:  []int{2, 0}},
			{
				Name:   `EmbedName`,
				Type:   gg.Type[string](),
				Offset: u.Offsetof(Outer{}.Embed) + u.Offsetof(Embed{}.EmbedName),
				Index:  []int{2, 1}},
			{
				Name:   `Inner`,
				Type:   gg.Type[*Inner](),
				Offset: u.Offsetof(Outer{}.Inner),
				Index:  []int{3}},
		},
	)
}

func TestJsonNameToDbNameCache(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.JsonNameToDbNameCache.Get(gg.Type[struct{}]()))

	gtest.Equal(
		gg.JsonNameToDbNameCache.Get(gg.Type[SomeJsonDbMapper]()),
		gg.JsonNameToDbName{
			`someName`:  `some_name`,
			`someValue`: `some_value`,
		},
	)
}

func TestDbNameToJsonNameCache(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.DbNameToJsonNameCache.Get(gg.Type[struct{}]()))

	gtest.Equal(
		gg.DbNameToJsonNameCache.Get(gg.Type[SomeJsonDbMapper]()),
		gg.DbNameToJsonName{
			`some_name`:  `someName`,
			`some_value`: `someValue`,
		},
	)
}
