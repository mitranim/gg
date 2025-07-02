//go:build !(386 || arm)

package gg

import "math"

func isInt32SafeForUint(val int32) bool { return val >= 0 }

func isUint32SafeForInt(uint32) bool { return true }

func isInt64SafeForInt(int64) bool { return true }

func isInt64SafeForUint(val int64) bool { return val >= 0 }

func isUint64SafeForInt(val uint64) bool {
	return val <= math.MaxInt64
}

func isUint64SafeForUint(uint64) bool { return true }

func isIntSafeForFloat64(val int) bool {
	return val >= MinSafeIntFloat64 && val <= MaxSafeIntFloat64
}

func isUintSafeForInt64(val uint) bool {
	return val <= math.MaxInt
}

func isUintSafeForFloat64(val uint) bool { return val <= MaxSafeIntFloat64 }

func isFloat64SafeForInt(val float64) bool {
	return !IsFrac(val) && val >= MinSafeIntFloat64 && val <= MaxSafeIntFloat64
}

func isFloat64SafeForUint(val float64) bool {
	return !IsFrac(val) && val >= 0 && val <= MaxSafeIntFloat64
}
