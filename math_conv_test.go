package gg_test

import (
	"math"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func Test_safe_int_floats(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.MinSafeIntFloat32, -16_777_215)
	gtest.Eq(gg.MaxSafeIntFloat32, 16_777_215)
	gtest.Eq(gg.MinSafeIntFloat64, -9_007_199_254_740_991)
	gtest.Eq(gg.MaxSafeIntFloat64, 9_007_199_254_740_991)
}

func TestNumConv_width_decrease_within_bounds(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.NumConv[uint8](int16(0)), 0)
	gtest.Eq(gg.NumConv[uint8](int16(128)), 128)
	gtest.Eq(gg.NumConv[uint8](int16(255)), 255)

	gtest.Eq(gg.NumConv[uint8](float32(0)), 0)
	gtest.Eq(gg.NumConv[uint8](float32(128)), 128)
	gtest.Eq(gg.NumConv[uint8](float32(255)), 255)

	gtest.Eq(gg.NumConv[uint8](float64(0)), 0)
	gtest.Eq(gg.NumConv[uint8](float64(128)), 128)
	gtest.Eq(gg.NumConv[uint8](float64(255)), 255)

	gtest.Eq(gg.NumConv[int8](int16(0)), 0)
	gtest.Eq(gg.NumConv[int8](int16(127)), 127)
	gtest.Eq(gg.NumConv[int8](int16(-128)), -128)

	gtest.Eq(gg.NumConv[int8](float32(0)), 0)
	gtest.Eq(gg.NumConv[int8](float32(127)), 127)
	gtest.Eq(gg.NumConv[int8](float32(-128)), -128)

	gtest.Eq(gg.NumConv[int8](float64(0)), 0)
	gtest.Eq(gg.NumConv[int8](float64(127)), 127)
	gtest.Eq(gg.NumConv[int8](float64(-128)), -128)

	gtest.Eq(gg.NumConv[float32](float64(0)), 0)
	gtest.Eq(gg.NumConv[float32](float64(math.MaxFloat32)), math.MaxFloat32)
	gtest.Eq(gg.NumConv[float32](float64(-math.MaxFloat32)), -math.MaxFloat32)
}

func TestNumConv_width_decrease_sign_mismatch(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to safely convert int16 -1 to uint8`,
		func() { gg.NumConv[uint8](int16(-1)) },
	)

	gtest.PanicStr(
		`unable to safely convert int16 -128 to uint8`,
		func() { gg.NumConv[uint8](int16(-128)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -128 to uint8`,
		func() { gg.NumConv[uint8](float32(-128)) },
	)
}

func TestNumConv_width_decrease_out_of_bounds(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to safely convert int16 256 to uint8`,
		func() { gg.NumConv[uint8](int16(256)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 256 to uint8`,
		func() { gg.NumConv[uint8](float32(256)) },
	)

	gtest.PanicStr(
		`unable to safely convert int16 128 to int8`,
		func() { gg.NumConv[int8](int16(128)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -170141173319264430000000000000000000000 to int16`,
		func() { gg.NumConv[int16](float32(-math.MaxFloat32 / 2)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 170141173319264430000000000000000000000 to int16`,
		func() { gg.NumConv[int16](float32(math.MaxFloat32 / 2)) },
	)

	gtest.PanicStr(
		`unable to safely convert float64 680564693277057700000000000000000000000 to float32`,
		func() { gg.NumConv[float32](float64(math.MaxFloat32 * 2)) },
	)

	gtest.PanicStr(
		`unable to safely convert float64 -680564693277057700000000000000000000000 to float32`,
		func() { gg.NumConv[float32](float64(-math.MaxFloat32 * 2)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 128 to int8`,
		func() { gg.NumConv[int8](float32(128)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 NaN to int16`,
		func() { gg.NumConv[int16](float32(math.NaN())) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 +Inf to int16`,
		func() { gg.NumConv[int16](float32(math.Inf(1))) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -Inf to int16`,
		func() { gg.NumConv[int16](float32(math.Inf(-1))) },
	)
}

func TestNumConv_width_decrease_imprecision(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to safely convert float32 10.5 to int16`,
		func() { gg.NumConv[int16](float32(10.5)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -10.5 to int16`,
		func() { gg.NumConv[int16](float32(-10.5)) },
	)
}

func TestNumConv_width_match_within_bounds(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.NumConv[uint8](int8(0)), 0)
	gtest.Eq(gg.NumConv[int8](uint8(0)), 0)

	gtest.Eq(gg.NumConv[uint8](int8(127)), 127)
	gtest.Eq(gg.NumConv[int8](uint8(127)), 127)

	gtest.Eq(gg.NumConv[int](uint(math.MaxInt)), math.MaxInt)
	gtest.Eq(gg.NumConv[uint](int(math.MaxInt)), math.MaxInt)

	gtest.Eq(gg.NumConv[int32](float32(0)), 0)
	gtest.Eq(gg.NumConv[int32](float32(0)), 0)

	gtest.Eq(gg.NumConv[int32](float32(gg.MinSafeIntFloat32)), gg.MinSafeIntFloat32)
	gtest.Eq(gg.NumConv[int32](float32(gg.MaxSafeIntFloat32)), gg.MaxSafeIntFloat32)

	gtest.Eq(gg.NumConv[float32](int32(gg.MinSafeIntFloat32)), gg.MinSafeIntFloat32)
	gtest.Eq(gg.NumConv[float32](int32(gg.MaxSafeIntFloat32)), gg.MaxSafeIntFloat32)
}

func TestNumConv_width_match_sign_mismatch(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to safely convert int8 -1 to uint8`,
		func() { gg.NumConv[uint8](int8(-1)) },
	)

	gtest.PanicStr(
		`unable to safely convert int8 -128 to uint8`,
		func() { gg.NumConv[uint8](int8(-128)) },
	)

	gtest.PanicStr(
		`unable to safely convert uint8 128 to int8`,
		func() { gg.NumConv[int8](uint8(128)) },
	)

	gtest.PanicStr(
		`unable to safely convert uint8 255 to int8`,
		func() { gg.NumConv[int8](uint8(255)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -1 to uint32`,
		func() { gg.NumConv[uint32](float32(-1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float64 -1 to uint64`,
		func() { gg.NumConv[uint64](float64(-1)) },
	)
}

func TestNumConv_width_match_out_of_bounds(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to safely convert uint 9223372036854775808 to int`,
		func() { gg.NumConv[int](uint(math.MaxInt + 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -170141173319264430000000000000000000000 to int32`,
		func() { gg.NumConv[int32](float32(-math.MaxFloat32 / 2)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 170141173319264430000000000000000000000 to int32`,
		func() { gg.NumConv[int32](float32(math.MaxFloat32 / 2)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 NaN to int32`,
		func() { gg.NumConv[int32](float32(math.NaN())) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 +Inf to int32`,
		func() { gg.NumConv[int32](float32(math.Inf(1))) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -Inf to int32`,
		func() { gg.NumConv[int32](float32(math.Inf(-1))) },
	)

	gtest.PanicStr(
		`unable to safely convert int32 -2147483648 to float32`,
		func() { gg.NumConv[float32](int32(math.MinInt32)) },
	)

	gtest.PanicStr(
		`unable to safely convert int32 2147483647 to float32`,
		func() { gg.NumConv[float32](int32(math.MaxInt32)) },
	)

	gtest.PanicStr(
		`unable to safely convert int64 -9223372036854775808 to float64`,
		func() { gg.NumConv[float64](int64(math.MinInt64)) },
	)

	gtest.PanicStr(
		`unable to safely convert int64 9223372036854775807 to float64`,
		func() { gg.NumConv[float64](int64(math.MaxInt64)) },
	)
}

func TestNumConv_width_match_imprecision(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to safely convert int32 16777216 to float32`,
		func() { gg.NumConv[float32](int32(gg.MaxSafeIntFloat32 + 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert int32 -16777216 to float32`,
		func() { gg.NumConv[float32](int32(gg.MinSafeIntFloat32 - 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert int32 2147483647 to float32`,
		func() { gg.NumConv[float32](int32(math.MaxInt32)) },
	)

	gtest.PanicStr(
		`unable to safely convert int64 9007199254740992 to float64`,
		func() { gg.NumConv[float64](int64(gg.MaxSafeIntFloat64 + 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert int64 -9007199254740992 to float64`,
		func() { gg.NumConv[float64](int64(gg.MinSafeIntFloat64 - 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert int64 9223372036854775807 to float64`,
		func() { gg.NumConv[float64](int64(math.MaxInt64)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 10.5 to int32`,
		func() { gg.NumConv[int32](float32(10.5)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -10.5 to int32`,
		func() { gg.NumConv[int32](float32(-10.5)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 16777216 to int32`,
		func() { gg.NumConv[int32](float32(gg.MaxSafeIntFloat32 + 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -16777216 to int32`,
		func() { gg.NumConv[int32](float32(gg.MinSafeIntFloat32 - 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 2147483648 to int32`,
		func() { gg.NumConv[int32](float32(math.MaxInt32)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -2147483648 to int32`,
		func() { gg.NumConv[int32](float32(math.MinInt32)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 16777216 to uint32`,
		func() { gg.NumConv[uint32](float32(gg.MaxSafeIntFloat32 + 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 4294967296 to uint32`,
		func() { gg.NumConv[uint32](float32(math.MaxUint32)) },
	)

	gtest.PanicStr(
		`unable to safely convert float64 9007199254740992 to int64`,
		func() { gg.NumConv[int64](float64(gg.MaxSafeIntFloat64 + 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float64 -9007199254740992 to int64`,
		func() { gg.NumConv[int64](float64(gg.MinSafeIntFloat64 - 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float64 9223372036854776000 to int64`,
		func() { gg.NumConv[int64](float64(math.MaxInt64)) },
	)

	gtest.PanicStr(
		`unable to safely convert float64 -9223372036854776000 to int64`,
		func() { gg.NumConv[int64](float64(math.MinInt64)) },
	)

	gtest.PanicStr(
		`unable to safely convert float64 9007199254740992 to uint64`,
		func() { gg.NumConv[uint64](float64(gg.MaxSafeIntFloat64 + 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float64 18446744073709552000 to uint64`,
		func() { gg.NumConv[uint64](float64(math.MaxUint64)) },
	)
}

func TestNumConv_width_increase_within_bounds(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Eq(gg.NumConv[uint16](uint8(0)), 0)
	gtest.Eq(gg.NumConv[uint16](uint8(128)), 128)
	gtest.Eq(gg.NumConv[uint16](uint8(255)), 255)

	gtest.Eq(gg.NumConv[int16](uint8(0)), 0)
	gtest.Eq(gg.NumConv[int16](uint8(128)), 128)
	gtest.Eq(gg.NumConv[int16](uint8(255)), 255)

	gtest.Eq(gg.NumConv[int16](int8(0)), 0)
	gtest.Eq(gg.NumConv[int16](int8(127)), 127)
	gtest.Eq(gg.NumConv[int16](int8(-128)), -128)

	gtest.Eq(gg.NumConv[float64](int32(math.MaxInt32)), math.MaxInt32)
	gtest.Eq(gg.NumConv[int64](float32(gg.MaxSafeIntFloat32)), gg.MaxSafeIntFloat32)

	gtest.Eq(gg.NumConv[float64](int32(math.MinInt32)), math.MinInt32)
	gtest.Eq(gg.NumConv[int64](float32(gg.MinSafeIntFloat32)), gg.MinSafeIntFloat32)

	gtest.Eq(gg.NumConv[float64](float32(0)), 0)
	gtest.Eq(gg.NumConv[float64](float32(math.MaxFloat32)), math.MaxFloat32)
	gtest.Eq(gg.NumConv[float64](float32(-math.MaxFloat32)), -math.MaxFloat32)
}

func TestNumConv_width_increase_sign_mismatch(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to safely convert int8 -1 to uint16`,
		func() { gg.NumConv[uint16](int8(-1)) },
	)

	gtest.PanicStr(
		`unable to safely convert int8 -128 to uint16`,
		func() { gg.NumConv[uint16](int8(-128)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -1 to uint64`,
		func() { gg.NumConv[uint64](float32(-1)) },
	)
}

func TestNumConv_width_increase_out_of_bounds(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to safely convert float32 -170141173319264430000000000000000000000 to int64`,
		func() { gg.NumConv[int64](float32(-math.MaxFloat32 / 2)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 170141173319264430000000000000000000000 to int64`,
		func() { gg.NumConv[int64](float32(math.MaxFloat32 / 2)) },
	)
}

func TestNumConv_width_increase_imprecision(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(
		`unable to safely convert float32 10.5 to int64`,
		func() { gg.NumConv[int64](float32(10.5)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -10.5 to int64`,
		func() { gg.NumConv[int64](float32(-10.5)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 NaN to int64`,
		func() { gg.NumConv[int64](float32(math.NaN())) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 +Inf to int64`,
		func() { gg.NumConv[int64](float32(math.Inf(1))) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -Inf to int64`,
		func() { gg.NumConv[int64](float32(math.Inf(-1))) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 16777216 to int64`,
		func() { gg.NumConv[int64](float32(gg.MaxSafeIntFloat32 + 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -16777216 to int64`,
		func() { gg.NumConv[int64](float32(gg.MinSafeIntFloat32 - 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 2147483648 to int64`,
		func() { gg.NumConv[int64](float32(math.MaxInt32)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -2147483648 to int64`,
		func() { gg.NumConv[int64](float32(math.MinInt32)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 9223372036854776000 to int64`,
		func() { gg.NumConv[int64](float32(math.MaxInt64)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 -9223372036854776000 to int64`,
		func() { gg.NumConv[int64](float32(math.MinInt64)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 16777216 to uint64`,
		func() { gg.NumConv[uint64](float32(gg.MaxSafeIntFloat32 + 1)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 4294967296 to uint64`,
		func() { gg.NumConv[uint64](float32(math.MaxUint32)) },
	)

	gtest.PanicStr(
		`unable to safely convert float32 18446744073709552000 to uint64`,
		func() { gg.NumConv[uint64](float32(math.MaxUint64)) },
	)
}

func TestNumConv_NaN(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(math.IsNaN(float64(gg.NumConv[float32](float32(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[Float32](float32(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[float64](float32(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[Float64](float32(math.NaN())))))

	gtest.True(math.IsNaN(float64(gg.NumConv[float32](float64(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[Float32](float64(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[float64](float64(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[Float64](float64(math.NaN())))))

	gtest.True(math.IsNaN(float64(gg.NumConv[float32](Float32(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[Float32](Float32(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[float64](Float32(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[Float64](Float32(math.NaN())))))

	gtest.True(math.IsNaN(float64(gg.NumConv[float32](Float64(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[Float32](Float64(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[float64](Float64(math.NaN())))))
	gtest.True(math.IsNaN(float64(gg.NumConv[Float64](Float64(math.NaN())))))
}

func TestNumConv_Inf(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(math.IsInf(float64(gg.NumConv[float32](float32(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float32](float32(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[float64](float32(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float64](float32(math.Inf(-1)))), -1))

	gtest.True(math.IsInf(float64(gg.NumConv[float32](float32(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float32](float32(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[float64](float32(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float64](float32(math.Inf(+1)))), +1))

	gtest.True(math.IsInf(float64(gg.NumConv[float32](float64(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float32](float64(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[float64](float64(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float64](float64(math.Inf(-1)))), -1))

	gtest.True(math.IsInf(float64(gg.NumConv[float32](float64(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float32](float64(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[float64](float64(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float64](float64(math.Inf(+1)))), +1))

	gtest.True(math.IsInf(float64(gg.NumConv[float32](Float32(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float32](Float32(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[float64](Float32(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float64](Float32(math.Inf(-1)))), -1))

	gtest.True(math.IsInf(float64(gg.NumConv[float32](Float32(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float32](Float32(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[float64](Float32(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float64](Float32(math.Inf(+1)))), +1))

	gtest.True(math.IsInf(float64(gg.NumConv[float32](Float64(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float32](Float64(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[float64](Float64(math.Inf(-1)))), -1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float64](Float64(math.Inf(-1)))), -1))

	gtest.True(math.IsInf(float64(gg.NumConv[float32](Float64(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float32](Float64(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[float64](Float64(math.Inf(+1)))), +1))
	gtest.True(math.IsInf(float64(gg.NumConv[Float64](Float64(math.Inf(+1)))), +1))
}

//go:noinline
func makeInt32() int32 { return gg.MaxSafeIntFloat32 }

//go:noinline
func makeFloat32() float32 { return gg.MaxSafeIntFloat32 }

//go:noinline
func numConvNativeIntToFloat(src int32) float64 { return float64(src) }

//go:noinline
func numConvNativeFloatToInt(src float32) int64 { return int64(src) }

//go:noinline
func numConvOursIntToFloat(src int32) float64 { return gg.NumConv[float64](src) }

//go:noinline
func numConvOursFloatToInt(src float32) int64 { return gg.NumConv[int64](src) }

//go:noinline
func numConvOursIntToInt(src int32) Int32 { return gg.NumConv[Int32](src) }

//go:noinline
func numConvOursFloatToFloat(src float32) Float32 { return gg.NumConv[Float32](src) }

func Benchmark_NumConv_int_to_float_native(b *testing.B) {
	defer gtest.Catch(b)
	src := makeInt32()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(numConvNativeIntToFloat(src))
	}
}

func Benchmark_NumConv_int_to_float_ours(b *testing.B) {
	defer gtest.Catch(b)
	src := makeInt32()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(numConvOursIntToFloat(src))
	}
}

func Benchmark_NumConv_float_to_int_native(b *testing.B) {
	defer gtest.Catch(b)
	src := makeFloat32()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(numConvNativeFloatToInt(src))
	}
}

func Benchmark_NumConv_float_to_int_ours(b *testing.B) {
	defer gtest.Catch(b)
	src := makeFloat32()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(numConvOursFloatToInt(src))
	}
}

func Benchmark_NumConv_equivalent_int_ours(b *testing.B) {
	defer gtest.Catch(b)
	src := makeInt32()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(numConvOursIntToInt(src))
	}
}

func Benchmark_NumConv_equivalent_float_ours(b *testing.B) {
	defer gtest.Catch(b)
	src := makeFloat32()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(numConvOursFloatToFloat(src))
	}
}
