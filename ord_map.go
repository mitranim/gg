package gg

/*
Implementation note: this currently lacks its own tests, but is indirectly
tested through `Coll`. TODO dedicated tests.
*/

/*
Represents an ordered map. Compare `OrdSet` which has only values, not key-value
pairs. Compare `Coll` which is an ordered map where each value determines its
own key.

This implementation is specialized for easy and efficient appending, iteration,
and membership testing, but as a tradeoff, it does not support deletion.
For "proper" ordered sets that support deletion, see the library
https://github.com/mitranim/gord.

Known limitations:

	* Lack of support for JSON encoding and decoding. An implementation using Go
	  maps would be easy but incorrect: element positions would be randomized.
*/
type OrdMap[Key comparable, Val any] struct {
	Slice []Val `role:"ref"`
	Index map[Key]int
}

// Same as `len(self.Slice)`.
func (self OrdMap[_, _]) Len() int { return len(self.Slice) }

// True if `.Len` <= 0. Inverse of `.IsNotEmpty`.
func (self OrdMap[_, _]) IsEmpty() bool { return self.Len() <= 0 }

// True if `.Len` > 0. Inverse of `.IsEmpty`.
func (self OrdMap[_, _]) IsNotEmpty() bool { return self.Len() > 0 }

// True if the index has the given key.
func (self OrdMap[Key, _]) Has(key Key) bool {
	return MapHas(self.Index, key)
}

// Returns the value indexed on the given key, or the zero value of that type.
func (self OrdMap[Key, Val]) Get(key Key) Val {
	return PtrGet(self.Ptr(key))
}

/*
Short for "get required". Returns the value indexed on the given key. Panics if
the value is missing.
*/
func (self OrdMap[Key, Val]) GetReq(key Key) Val {
	ptr := self.Ptr(key)
	if ptr != nil {
		return *ptr
	}
	panic(errCollMissing[Val](key))
}

/*
Returns the value indexed on the given key and a boolean indicating if the value
was actually present.
*/
func (self OrdMap[Key, Val]) Got(key Key) (Val, bool) {
	// Note: we must check `ok` because if the entry is missing, `ind` is `0`,
	// which is invalid.
	ind, ok := self.Index[key]
	if ok {
		return Got(self.Slice, ind)
	}
	return Zero[Val](), false
}

/*
Short for "pointer". Returns a pointer to the value indexed on the given key, or
nil if the value is missing. Because this type does not support deletion, the
correspondence of positions in `.Slice` and indexes in `.Index` does not change
when adding or replacing values. The pointer should remain valid for the
lifetime of the ordered map, unless `.Slice` is directly mutated by external
means.
*/
func (self OrdMap[Key, Val]) Ptr(key Key) *Val {
	// Note: we must check `ok` because if the entry is missing, `ind` is `0`,
	// which is invalid.
	ind, ok := self.Index[key]
	if ok {
		return GetPtr(self.Slice, ind)
	}
	return nil
}

/*
Short for "pointer required". Returns a non-nil pointer to the value indexed
on the given key, or panics if the value is missing.
*/
func (self OrdMap[Key, Val]) PtrReq(key Key) *Val {
	ptr := self.Ptr(key)
	if ptr != nil {
		return ptr
	}
	panic(errCollMissing[Val](key))
}

/*
Idempotently adds or replaces the given value, updating both the inner slice and
the inner index. If the key was already registered in the map, the new value
replaces the old value at the same position in the inner slice.
*/
func (self *OrdMap[Key, Val]) Set(key Key, val Val) *OrdMap[Key, Val] {
	index := MapInit(&self.Index)
	ind, ok := index[key]
	if ok {
		self.Slice[ind] = val
		return self
	}

	index[key] = AppendIndex(&self.Slice, val)
	return self
}

// Same as `OrdMap.Set`, but panics if the key was already present in the index.
func (self *OrdMap[Key, Val]) Add(key Key, val Val) *OrdMap[Key, Val] {
	index := MapInit(&self.Index)

	if MapHas(index, key) {
		panic(Errf(
			`unexpected redundant %v with key %v`,
			Type[Val](), key,
		))
	}

	index[key] = AppendIndex(&self.Slice, val)
	return self
}

// Nullifies both the slice and the index. Does not preserve their capacity.
func (self *OrdMap[Key, Val]) Clear() *OrdMap[Key, Val] {
	if self != nil {
		self.Slice = nil
		self.Index = nil
	}
	return self
}
