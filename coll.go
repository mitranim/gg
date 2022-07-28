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
Short for "collection". Represents an ordered map where keys are automatically
derived from values, and must be non-zero.
*/
type Coll[
	Key comparable,
	Val Pked[Key],
] struct {
	Slice []Val `role:"ref"`
	Index map[Key]int
}

/*
Reindexes the collection. Must be invoked after appending elements to the slice
through external means. Note that all built-in methods of this type perform
indexing automatically. This method must be invoked if the collection is
modified by directly accessing `.Slice` and/or `.Index`.
*/
func (self *Coll[Key, Val]) Calc() {
	if !self.isIndexed() {
		index := self.initIndex()
		MapClear(index)

		for ind, val := range self.Slice {
			index[ValidPk[Key](val)] = ind
		}
	}
}

// Clears the collection, keeping capacity.
func (self *Coll[Key, Val]) Clear() *Coll[Key, Val] {
	SliceTrunc(&self.Slice)
	MapClear(self.Index)
	return self
}

// Adds the given elements to both the inner slice and the inner index.
func (self *Coll[Key, Val]) Add(val ...Val) *Coll[Key, Val] {
	for _, val := range val {
		index := self.initIndex()
		key := ValidPk[Key](val)

		AppendVals(&self.Slice, val)
		index[key] = len(self.Slice) - 1
	}
	return self
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
Find the value indexed on the given key and returns the pointer to its position
in the slice. If the value is not found, returns nil.
*/
func (self Coll[Key, Val]) Ptr(key Key) *Val {
	ind, ok := MapGot(self.Index, key)
	if !ok {
		return nil
	}
	return GetPtr(self.Slice, ind)
}

// Same as `len(self.Index)`.
func (self Coll[_, _]) Len() int { return len(self.Index) }

// True if length > 0.
func (self Coll[_, _]) HasLen() bool { return self.Len() > 0 }

// Inverse of `.HasLen`.
func (self Coll[_, _]) IsEmpty() bool { return !self.HasLen() }

// Implement `json.Marshaler`. Encodes the inner slice, ignoring the index.
func (self Coll[_, _]) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Slice)
}

// Unmarshals the input into the inner slice and rebuilds the index.
func (self *Coll[_, _]) UnmarshalJSON(src []byte) error {
	err := json.Unmarshal(src, &self.Slice)
	self.Calc()
	return err
}

func (self *Coll[Key, _]) initIndex() map[Key]int {
	return MapInit(&self.Index)
}

func (self *Coll[_, _]) isIndexed() bool {
	return len(self.Slice) == len(self.Index) && Every(self.Slice, self.isValIndexed)
}

func (self *Coll[_, Val]) isValIndexed(val Val) bool {
	return MapHas(self.Index, val.Pk())
}
