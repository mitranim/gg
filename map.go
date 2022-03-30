package gg

func MapClone[Coll ~map[Key]Val, Key comparable, Val any](src Coll) Coll {
	if src == nil {
		return nil
	}

	out := make(Coll, len(src))
	for key, val := range src {
		out[key] = val
	}
	return out
}

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

func MapPtrInit[Map ~map[Key]Val, Key comparable, Val any](val *Map) Map {
	*val = MapInit(*val)
	return *val
}

func MapInit[Map ~map[Key]Val, Key comparable, Val any](val Map) Map {
	if val == nil {
		val = make(map[Key]Val)
	}
	return val
}

func MapHas[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) bool {
	_, ok := tar[key]
	return ok
}

func MapGot[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) (Val, bool) {
	val, ok := tar[key]
	return val, ok
}

func MapGet[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) Val {
	return tar[key]
}

func MapSet[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key, val Val) map[Key]Val {
	tar[key] = val
	return tar
}

func MapDel[Map ~map[Key]Val, Key comparable, Val any](tar Map, key Key) map[Key]Val {
	delete(tar, key)
	return tar
}

func MapClear[Map ~map[Key]Val, Key comparable, Val any](tar Map) map[Key]Val {
	for key := range tar {
		delete(tar, key)
	}
	return tar
}
