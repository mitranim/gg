package gg

/*
Shortcut for converting an arbitrary map to `Dict`. Workaround for the
limitations of type inference in Go generics.
*/
func ToDict[Src ~map[Key]Val, Key comparable, Val any](val Src) Dict[Key, Val] {
	return Dict[Key, Val](val)
}

/*
Typedef of an arbitrary map with various methods that duplicate global map
functions. Useful as a shortcut for creating bound methods that can be passed
to higher-order functions.
*/
type Dict[Key comparable, Val any] map[Key]Val

/*
Idempotently initializes the map at the given pointer via `make`, returning the
result. Pointer must be non-nil. If map was non-nil, it's unchanged. Output is
always non-nil.
*/
func MapInit[Map ~map[Key]Val, Key comparable, Val any](val *Map) Map {
	if val == nil {
		return nil
	}
	if *val == nil {
		*val = make(map[Key]Val)
	}
	return *val
}

// Self as global `MapInit`.
//go:noinline
func (self *Dict[Key, Val]) Init() Dict[Key, Val] { return MapInit(self) }

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

// Self as global `MapClone`.
//go:noinline
func (self Dict[Key, Val]) Clone() Dict[Key, Val] { return MapClone(self) }

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

// Self as global `MapKeys`.
//go:noinline
func (self Dict[Key, _]) Keys() []Key { return MapKeys(self) }

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

// Self as global `MapVals`.
//go:noinline
func (self Dict[_, Val]) Vals() []Val { return MapVals(self) }

// Same as `_, ok := tar[key]`, expressed as a generic function.
func MapHas[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) bool {
	_, ok := tar[key]
	return ok
}

// Self as global `MapHas`.
//go:noinline
func (self Dict[Key, _]) Has(key Key) bool { return MapHas(self, key) }

// Same as `val, ok := tar[key]`, expressed as a generic function.
func MapGot[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) (Val, bool) {
	val, ok := tar[key]
	return val, ok
}

// Self as global `MapGot`.
//go:noinline
func (self Dict[Key, Val]) Got(key Key) (Val, bool) { return MapGot(self, key) }

// Same as `val := tar[key]`, expressed as a generic function.
func MapGet[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) Val {
	return tar[key]
}

// Self as global `MapGet`.
//go:noinline
func (self Dict[Key, Val]) Get(key Key) Val { return MapGet(self, key) }

// Same as `tar[key] = val`, expressed as a generic function.
func MapSet[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key, val Val) {
	tar[key] = val
}

// Self as global `MapSet`.
//go:noinline
func (self Dict[Key, Val]) Set(key Key, val Val) { MapSet(self, key, val) }

/*
Same as `MapSet`, but key and value should be be non-zero.
If either is zero, this ignores the inputs and does nothing.
*/
func MapSetOpt[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key, val Val) {
	if IsNonZero(key) && IsNonZero(val) {
		MapSet(tar, key, val)
	}
}

// Self as global `MapSetOpt`.
//go:noinline
func (self Dict[Key, Val]) SetOpt(key Key, val Val) { MapSetOpt(self, key, val) }

// Same as `delete(tar, key)`, expressed as a generic function.
func MapDel[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) {
	delete(tar, key)
}

// Self as global `MapDel`.
//go:noinline
func (self Dict[Key, _]) Del(key Key) { delete(self, key) }

// Deletes all entries, returning the resulting map. Passing nil is safe.
func MapClear[Map ~map[Key]Val, Key comparable, Val any](tar Map) {
	for key := range tar {
		delete(tar, key)
	}
}

// Self as global `MapClear`.
//go:noinline
func (self Dict[_, _]) Clear() { MapClear(self) }

// Needs a better name.
func MapDict[Key comparable, A, B any](src map[Key]A, fun func(A) B) map[Key]B {
	if src == nil || fun == nil {
		return nil
	}

	out := make(map[Key]B, len(src))
	for key, val := range src {
		out[key] = fun(val)
	}
	return out
}
