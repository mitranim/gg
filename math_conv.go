package gg

import (
	"math"
	r "reflect"
)

/*
Checked numeric conversion. Same as a built-in Go conversion, but panics in case
of overflow, underflow, or imprecision when converting between integers and
floats. Converting NaN and infinities between different floating point types
is allowed. Performance overhead is measurable but small.
*/
func NumConv[Out, Src Num](src Src) Out {
	out := Out(src)
	outKind := r.TypeOf(out).Kind()

	switch r.TypeOf(src).Kind() {
	case r.Int8:
		switch outKind {
		case r.Uint8, r.Uint16, r.Uint32, r.Uint64, r.Uint:
			if src >= 0 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int8, r.Int16, r.Int32, r.Int64, r.Int, r.Float32, r.Float64:
			return out
		}

	case r.Uint8:
		switch outKind {
		case r.Int8:
			if src <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8, r.Int16, r.Uint16, r.Int32, r.Uint32, r.Int64, r.Uint64, r.Int, r.Uint, r.Float32, r.Float64:
			return out
		}

	case r.Int16:
		switch outKind {
		case r.Int8:
			val := int16(src)
			if val >= math.MinInt8 && val <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8:
			if src >= 0 && int16(src) <= math.MaxUint8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint16, r.Uint32, r.Uint64, r.Uint:
			if src >= 0 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int16, r.Int32, r.Int64, r.Int, r.Float32, r.Float64:
			return out
		}

	case r.Uint16:
		switch outKind {
		case r.Int8:
			if src <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8:
			if uint16(src) <= math.MaxUint8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int16:
			if uint16(src) <= math.MaxInt16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int32, r.Uint32, r.Int64, r.Uint64, r.Int, r.Uint, r.Float32, r.Float64:
			return out
		}

	case r.Int32:
		switch outKind {
		case r.Int8:
			if int32(src) >= math.MinInt8 && src <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8:
			if src >= 0 && int32(src) <= math.MaxUint8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int16:
			val := int32(src)
			if val >= math.MinInt16 && val <= math.MaxInt16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint16:
			if src >= 0 && int32(src) <= math.MaxUint16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint32, r.Uint64:
			if int32(src) >= 0 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int:
			return out
		case r.Uint:
			if isInt32SafeForUint(int32(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float32:
			val := int32(src)
			if val >= MinSafeIntFloat32 && val <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int32, r.Int64, r.Float64:
			return out
		}

	case r.Uint32:
		switch outKind {
		case r.Int8:
			if src <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8:
			if uint32(src) <= math.MaxUint8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int16:
			if uint32(src) <= math.MaxInt16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint16:
			if uint32(src) <= math.MaxUint16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int32:
			if uint32(src) <= math.MaxInt32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int:
			if isUint32SafeForInt(uint32(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float32:
			if uint32(src) <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint32, r.Int64, r.Uint64, r.Uint, r.Float64:
			return out
		}

	case r.Int64:
		switch outKind {
		case r.Int8:
			if int64(src) >= math.MinInt8 && src <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8:
			if src >= 0 && int64(src) <= math.MaxUint8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int16:
			val := int64(src)
			if val >= math.MinInt16 && val <= math.MaxInt16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint16:
			if src >= 0 && int64(src) <= math.MaxUint16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int32:
			val := int64(src)
			if val >= math.MinInt32 && val <= math.MaxInt32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint32:
			if src >= 0 && int64(src) <= math.MaxUint32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int64:
			return out
		case r.Uint64:
			if src >= 0 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int:
			if isInt64SafeForInt(int64(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint:
			if isInt64SafeForUint(int64(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float32:
			val := int64(src)
			if val >= MinSafeIntFloat32 && val <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float64:
			val := int64(src)
			if val >= MinSafeIntFloat64 && val <= MaxSafeIntFloat64 {
				return out
			}
			panic(errNumConv[Out](src))
		}

	case r.Uint64:
		switch outKind {
		case r.Int8:
			if src <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8:
			if uint64(src) <= math.MaxUint8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int16:
			if uint64(src) <= math.MaxInt16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint16:
			if uint64(src) <= math.MaxUint16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int32:
			if uint64(src) <= math.MaxInt32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint32:
			if uint64(src) <= math.MaxUint32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int64:
			if uint64(src) <= math.MaxInt64 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint64:
			return out
		case r.Int:
			if isUint64SafeForInt(uint64(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint:
			if isUint64SafeForUint(uint64(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float32:
			if uint64(src) <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float64:
			if uint64(src) <= MaxSafeIntFloat64 {
				return out
			}
			panic(errNumConv[Out](src))
		}

	case r.Int:
		switch outKind {
		case r.Int8:
			if int(src) >= math.MinInt8 && src <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8:
			if src >= 0 && int(src) <= math.MaxUint8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int16:
			val := int(src)
			if val >= math.MinInt16 && val <= math.MaxInt16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint16:
			if src >= 0 && int(src) <= math.MaxUint16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int32:
			val := int(src)
			if val >= math.MinInt32 && val <= math.MaxInt32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint32:
			if src >= 0 && int(src) <= math.MaxInt32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int64:
			return out
		case r.Uint64:
			if src >= 0 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int:
			return out
		case r.Uint:
			if src >= 0 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float32:
			val := int(src)
			if val >= MinSafeIntFloat32 && val <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float64:
			if isIntSafeForFloat64(int(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		}

	case r.Uint:
		switch outKind {
		case r.Int8:
			if src <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8:
			if uint(src) <= math.MaxUint8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int16:
			if uint(src) <= math.MaxInt16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint16:
			if uint(src) <= math.MaxUint16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int32:
			if uint(src) <= math.MaxInt32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint32:
			if uint(src) <= math.MaxUint32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int64:
			if isUintSafeForInt64(uint(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint64, r.Uint:
			return out
		case r.Int:
			if uint(src) <= math.MaxInt {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float32:
			if uint(src) <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float64:
			if isUintSafeForFloat64(uint(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		}

	case r.Float32:
		switch outKind {
		case r.Int8:
			val := float32(src)
			if !IsFrac(val) && val >= math.MinInt8 && val <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8:
			val := float32(src)
			if !IsFrac(val) && val >= 0 && val <= math.MaxUint8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int16:
			val := float32(src)
			if !IsFrac(val) && val >= math.MinInt16 && val <= math.MaxInt16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint16:
			val := float32(src)
			if !IsFrac(val) && val >= 0 && val <= math.MaxUint16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int32:
			val := float32(src)
			if !IsFrac(val) && val >= MinSafeIntFloat32 && val <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint32:
			val := float32(src)
			if !IsFrac(val) && val >= 0 && val <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int64:
			val := float32(src)
			if !IsFrac(val) && val >= MinSafeIntFloat32 && val <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint64:
			val := float32(src)
			if !IsFrac(val) && val >= 0 && val <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int:
			val := float32(src)
			if !IsFrac(val) && val >= MinSafeIntFloat32 && val <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint:
			val := float32(src)
			if !IsFrac(val) && val >= 0 && val <= MaxSafeIntFloat32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float32, r.Float64:
			return out
		}

	case r.Float64:
		switch outKind {
		case r.Int8:
			val := float64(src)
			if !IsFrac(val) && val >= math.MinInt8 && val <= math.MaxInt8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint8:
			val := float64(src)
			if !IsFrac(val) && val >= 0 && val <= math.MaxUint8 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int16:
			val := float64(src)
			if !IsFrac(val) && val >= math.MinInt16 && val <= math.MaxInt16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint16:
			val := float64(src)
			if !IsFrac(val) && val >= 0 && val <= math.MaxUint16 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int32:
			val := float64(src)
			if !IsFrac(val) && val >= math.MinInt32 && val <= math.MaxInt32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint32:
			val := float64(src)
			if !IsFrac(val) && val >= 0 && val <= math.MaxUint32 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int64:
			val := float64(src)
			if !IsFrac(val) && val >= MinSafeIntFloat64 && val <= MaxSafeIntFloat64 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint64:
			val := float64(src)
			if !IsFrac(val) && val >= 0 && val <= MaxSafeIntFloat64 {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Int:
			if isFloat64SafeForInt(float64(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Uint:
			if isFloat64SafeForUint(float64(src)) {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float32:
			val := float64(src)
			if !IsFin(val) || (val >= -math.MaxFloat32 && val <= math.MaxFloat32) {
				return out
			}
			panic(errNumConv[Out](src))
		case r.Float64:
			return out
		}
	}

	/**
	Technical note. An older version of this function used the following check:

		if src == Src(out) &&
			!(src < 0 && out >= 0) &&
			!(src >= 0 && out < 0) {
			return out
		}

	...But we can't rely on this. For some Go versions and CPU architectures,
	some conversions between ints and floats produce an incorrect result in one
	direction, but are reversible and thus not detectable with this approach.
	*/

	panic(errNumConv[Out](src))
}

// Uses `String` to avoid the scientific notation for floats.
func errNumConv[Out, Src Num](src Src) Err {
	return Errf(
		`unable to safely convert %v %v to %v`,
		Type[Src](), String(src), Type[Out](),
	)
}
