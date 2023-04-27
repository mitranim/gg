package gg_test

import (
	"context"
	"fmt"
	"math"
	r "reflect"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

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

func TestPlus(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`Num`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.Plus[int](), 0)
		gtest.Eq(gg.Plus(0), 0)
		gtest.Eq(gg.Plus(10), 10)
		gtest.Eq(gg.Plus(10, 20), 30)
		gtest.Eq(gg.Plus(-10, 0), -10)
		gtest.Eq(gg.Plus(-10, 0, 10), 0)
		gtest.Eq(gg.Plus(-10, 0, 10, 20), 20)
		gtest.Eq(gg.Plus(-10, 0, 10, 20, 30), 50)
	})

	t.Run(`string`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Eq(gg.Plus[string](), ``)
		gtest.Eq(gg.Plus(``), ``)
		gtest.Eq(gg.Plus(`one`), `one`)
		gtest.Eq(gg.Plus(`one`, ``), `one`)
		gtest.Eq(gg.Plus(``, `two`), `two`)
		gtest.Eq(gg.Plus(`one`, `two`), `onetwo`)
		gtest.Eq(gg.Plus(`one`, `two`, `three`), `onetwothree`)
	})
}

func BenchmarkPlus(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Plus(10, 20, 30, 40, 50, 60, 70, 80, 90))
	}
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
