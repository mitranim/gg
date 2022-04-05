package gg

/*
Copies the given map. If the input is nil, the output is nil. Otherwise the
output is a shallow copy.
*/
func MapClone[Map ~map[Key]Val, Key comparable, Val any](src Map) Map {
	if src == nil {
		return nil
	}

	out := make(Map, len(src))
	for key, val := range src {
		out[key] = val
	}
	return out
}

// Returns the maps's keys as a slice. Order is random.
func MapKeys[Key comparable, Val any](src map[Key]Val) []Key {
	if src == nil {
		return nil
	}

	out := make([]Key, 0, len(src))
	for key := range src {
		out = append(out, key)
	}
	return out
}

// Returns the maps's values as a slice. Order is random.
func MapVals[Key comparable, Val any](src map[Key]Val) []Val {
	if src == nil {
		return nil
	}

	out := make([]Val, 0, len(src))
	for _, val := range src {
		out = append(out, val)
	}
	return out
}

/*
Idempotently initializes the map at the given pointer via `make`, returning the
result. Pointer must be non-nil. If map was non-nil, it's unchanged. Output is
always non-nil.
*/
func MapPtrInit[Map ~map[Key]Val, Key comparable, Val any](val *Map) Map {
	*val = MapInit(*val)
	return *val
}

/*
Idempotently initializes the map via `make`. If the input is already non-nil,
it's returned as-is.
*/
func MapInit[Map ~map[Key]Val, Key comparable, Val any](val Map) Map {
	if val == nil {
		val = make(map[Key]Val)
	}
	return val
}

// Same as `_, ok := tar[key]`, expressed as a generic function.
func MapHas[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) bool {
	_, ok := tar[key]
	return ok
}

// Same as `val, ok := tar[key]`, expressed as a generic function.
func MapGot[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) (Val, bool) {
	val, ok := tar[key]
	return val, ok
}

// Same as `val := tar[key]`, expressed as a generic function.
func MapGet[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) Val {
	return tar[key]
}

// Same as `tar[key] = val`, expressed as a generic function.
func MapSet[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key, val Val) map[Key]Val {
	tar[key] = val
	return tar
}

// Same as `delete(tar, key)`, expressed as a generic function.
func MapDel[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) map[Key]Val {
	delete(tar, key)
	return tar
}

// Deletes all entries, returning the resulting map. Passing nil is safe.
func MapClear[Map ~map[Key]Val, Key comparable, Val any](tar Map) map[Key]Val {
	for key := range tar {
		delete(tar, key)
	}
	return tar
}
