//go:build 386

package gg

import "math"

func isInt32SafeForUint(val int32) bool {
	return val >= 0 && val <= math.MaxInt32
}

func isUint32SafeForInt(val uint32) bool { return val <= math.MaxInt32 }

func isInt64SafeForInt(val int64) bool {
	return val >= math.MinInt32 && val <= math.MaxInt32
}

func isInt64SafeForUint(val int64) bool {
	return val >= 0 && val <= math.MaxUint32
}

func isUint64SafeForInt(val uint64) bool {
	return val <= math.MaxInt32
}

func isUint64SafeForUint(val uint64) bool {
	return val <= math.MaxUint32
}

func isIntSafeForFloat64(int) bool { return true }

func isUintSafeForInt64(uint) bool { return true }

func isUintSafeForFloat64(uint) bool { return true }

func isFloat64SafeForInt(val float64) bool {
	return !IsFrac(val) && val >= math.MinInt32 && val <= math.MaxInt32
}

func isFloat64SafeForUint(val float64) bool {
	return !IsFrac(val) && val >= 0 && val <= math.MaxUint32
}
