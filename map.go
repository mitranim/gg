package gg

/*
Non-idempotent version of `MapInit`. If the target pointer is nil, does nothing
and returns nil. If the target pointer is non-nil, allocates the map via
`make`, stores it at the target pointer, and returns the resulting non-nil
map.
*/
func MapMake[Map ~map[Key]Val, Key comparable, Val any](ptr *Map) Map {
	if ptr == nil {
		return nil
	}
	val := make(map[Key]Val)
	*ptr = val
	return val
}

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

// Same as `len(self)`.
func (self Dict[_, _]) Len() int { return len(self) }

// Same as `len(self) <= 0`. Inverse of `.IsNotEmpty`.
func (self Dict[_, _]) IsEmpty() bool { return len(self) <= 0 }

// Same as `len(self) > 0`. Inverse of `.IsEmpty`.
func (self Dict[_, _]) IsNotEmpty() bool { return len(self) > 0 }

/*
Idempotent map initialization. If the target pointer is nil, does nothing and
returns nil. If the map at the target pointer is non-nil, does nothing and
returns that map. Otherwise allocates the map via `make`, stores it at the
target pointer, and returns the resulting non-nil map.
*/
func MapInit[Map ~map[Key]Val, Key comparable, Val any](ptr *Map) Map {
	if ptr == nil {
		return nil
	}
	val := *ptr
	if val == nil {
		val = make(map[Key]Val)
		*ptr = val
	}
	return val
}

// Self as global `MapInit`.
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
func (self Dict[_, Val]) Vals() []Val { return MapVals(self) }

// Same as `_, ok := tar[key]`, expressed as a generic function.
func MapHas[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) bool {
	_, ok := tar[key]
	return ok
}

// Self as global `MapHas`.
func (self Dict[Key, _]) Has(key Key) bool { return MapHas(self, key) }

// Same as `val, ok := tar[key]`, expressed as a generic function.
func MapGot[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) (Val, bool) {
	val, ok := tar[key]
	return val, ok
}

// Self as global `MapGot`.
func (self Dict[Key, Val]) Got(key Key) (Val, bool) { return MapGot(self, key) }

// Same as `val := tar[key]`, expressed as a generic function.
func MapGet[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) Val {
	return tar[key]
}

// Self as global `MapGet`.
func (self Dict[Key, Val]) Get(key Key) Val { return MapGet(self, key) }

// Same as `tar[key] = val`, expressed as a generic function.
func MapSet[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key, val Val) {
	tar[key] = val
}

// Self as global `MapSet`.
func (self Dict[Key, Val]) Set(key Key, val Val) { MapSet(self, key, val) }

/*
Same as `MapSet`, but key and value should be be non-zero.
If either is zero, this ignores the inputs and does nothing.
*/
func MapSetOpt[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key, val Val) {
	if IsNotZero(key) && IsNotZero(val) {
		MapSet(tar, key, val)
	}
}

// Self as global `MapSetOpt`.
func (self Dict[Key, Val]) SetOpt(key Key, val Val) { MapSetOpt(self, key, val) }

// Same as `delete(tar, key)`, expressed as a generic function.
func MapDel[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) {
	delete(tar, key)
}

// Self as global `MapDel`.
func (self Dict[Key, _]) Del(key Key) { delete(self, key) }

/*
Deletes all entries, returning the resulting map. Passing nil is safe.
Note that this involves iterating the map, which is inefficient in Go.
In many cases, it's more efficient to make a new map.
*/
func MapClear[Map ~map[Key]Val, Key comparable, Val any](tar Map) {
	for key := range tar {
		delete(tar, key)
	}
}

// Self as global `MapClear`.
func (self Dict[_, _]) Clear() { MapClear(self) }
