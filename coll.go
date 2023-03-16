package gg

import "encoding/json"

/*
Short for "valid primary key". Returns the primary key generated by the given
input, asserts that the key is non-zero, and returns the resulting key.
Used internally by `Coll`.
*/
func ValidPk[
	Key comparable,
	Val Pked[Key],
](val Val) Key {
	key := val.Pk()
	if IsZero(key) {
		panic(Errf(`unexpected empty key of type %T for %#v`, key, val))
	}
	return key
}

/*
Syntactic shortcut for making a `Coll` of the given arguments, with type
inference.
*/
func CollOf[Key comparable, Val Pked[Key]](src ...Val) Coll[Key, Val] {
	var tar Coll[Key, Val]
	tar.Add(src...)
	return tar
}

/*
Syntactic shortcut for making a `Coll` from any number of source slices, with
type inference.
*/
func CollFrom[Slice ~[]Val, Key comparable, Val Pked[Key]](src ...Slice) Coll[Key, Val] {
	var tar Coll[Key, Val]
	for _, src := range src {
		tar.Add(src...)
	}
	return tar
}

/*
Short for "collection". Represents an ordered map where keys are automatically
derived from values. Keys must be non-zero. Similarly to a map, this ensures
value uniqueness by primary key, and allows efficient access by key. Unlike a
map, values in this type are ordered and can be iterated cheaply, because they
are stored in a publicly-accessible slice. However, as a tradeoff, this type
does not support deletion.
*/
type Coll[
	Key comparable,
	Val Pked[Key],
] struct {
	Slice []Val `role:"ref"`
	Index map[Key]int
}

// Same as `len(self.Slice)`.
func (self Coll[_, _]) Len() int { return len(self.Slice) }

// True if `.Len` > 0.
func (self Coll[_, _]) HasLen() bool { return self.Len() > 0 }

// Inverse of `.HasLen`.
func (self Coll[_, _]) IsEmpty() bool { return !self.HasLen() }

/*
True if the index has the given key. Doesn't check if the index is within the
bounds of the inner slice.
*/
func (self Coll[Key, _]) Has(key Key) bool {
	_, ok := self.Index[key]
	return ok
}

// Returns the value indexed on the given key, or the zero value of that type.
func (self Coll[Key, Val]) Get(key Key) Val {
	return PtrGet(self.Ptr(key))
}

/*
Returns the value indexed on the given key and a boolean indicating if the value
was actually present.
*/
func (self Coll[Key, Val]) Got(key Key) (Val, bool) {
	ptr := self.Ptr(key)
	return PtrGet(ptr), ptr != nil
}

/*
Find the value indexed on the given key and returns the pointer to its position
in the slice. If the value is not found, returns nil. Caution: sorting the inner
slice invalidates such pointers.
*/
func (self Coll[Key, Val]) Ptr(key Key) *Val {
	ind, ok := MapGot(self.Index, key)
	if !ok {
		return nil
	}
	return GetPtr(self.Slice, ind)
}

/*
Idempotently adds each given value to both the inner slice and the inner index.
Every value whose key already exists in the index is replaced at the existing
position in the slice.
*/
func (self *Coll[Key, Val]) Add(src ...Val) *Coll[Key, Val] {
	index := MapInit(&self.Index)

	for _, val := range src {
		key := ValidPk[Key](val)
		ind, ok := index[key]
		if ok {
			self.Slice[ind] = val
			continue
		}
		index[key] = AppendIndex(&self.Slice, val)
	}

	return self
}

// Nullifies both the slice and the index. Does not preserve their capacity.
func (self *Coll[Key, Val]) Clear() *Coll[Key, Val] {
	if self != nil {
		self.Slice = nil
		self.Index = nil
	}
	return self
}

/*
Rebuilds the inner index from the inner slice, without checking the validity of
the existing index. Can be useful for external code that directly modifies the
inner `.Slice`, for example by sorting it. This is NOT used when adding items
via `.Add`, which modifies the index incrementally rather than all-at-once.
*/
func (self *Coll[Key, Val]) Reindex() {
	src := self.Slice
	self.Clear()
	self.Slice = src[:0]
	self.Add(src...)
}

// Implement `json.Marshaler`. Encodes the inner slice, ignoring the index.
func (self Coll[_, _]) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Slice)
}

// Unmarshals the input into the inner slice and rebuilds the index.
func (self *Coll[_, _]) UnmarshalJSON(src []byte) error {
	err := json.Unmarshal(src, &self.Slice)
	self.Reindex()
	return err
}
