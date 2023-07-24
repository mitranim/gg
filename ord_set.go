package gg

import "encoding/json"

/*
Syntactic shortcut for making an `OrdSet` of the given arguments, with type
inference.
*/
func OrdSetOf[Val comparable](src ...Val) OrdSet[Val] {
	var tar OrdSet[Val]
	tar.Add(src...)
	return tar
}

/*
Syntactic shortcut for making an `OrdSet` from any number of source slices, with
type inference.
*/
func OrdSetFrom[Slice ~[]Val, Val comparable](src ...Slice) OrdSet[Val] {
	var tar OrdSet[Val]
	for _, src := range src {
		tar.Add(src...)
	}
	return tar
}

/*
Represents an ordered set. Compare `OrdMap` which represents an ordered map.
This implementation is specialized for easy and efficient appending, iteration,
and membership testing, but as a tradeoff, it does not support deletion.
For "proper" ordered sets that support deletion, see the library
https://github.com/mitranim/gord.
*/
type OrdSet[Val comparable] struct {
	Slice []Val `role:"ref"`
	Index Set[Val]
}

// Same as `len(self.Slice)`.
func (self OrdSet[_]) Len() int { return len(self.Slice) }

// True if `.Len` <= 0. Inverse of `.IsNotEmpty`.
func (self OrdSet[_]) IsEmpty() bool { return self.Len() <= 0 }

// True if `.Len` > 0. Inverse of `.IsEmpty`.
func (self OrdSet[_]) IsNotEmpty() bool { return self.Len() > 0 }

// True if the index has the given value. Ignores the inner slice.
func (self OrdSet[Val]) Has(val Val) bool { return self.Index.Has(val) }

/*
Idempotently adds each given value to both the inner slice and the inner
index, skipping duplicates.
*/
func (self *OrdSet[Val]) Add(src ...Val) *OrdSet[Val] {
	for _, val := range src {
		if !self.Has(val) {
			Append(&self.Slice, val)
			self.Index.Init().Add(val)
		}
	}
	return self
}

/*
Replaces `.Slice` with the given slice and rebuilds `.Index`. Uses the slice
as-is with no reallocation. Callers must be careful to avoid modifying the
source data, which may invalidate the collection's index.
*/
func (self *OrdSet[Val]) Reset(src ...Val) *OrdSet[Val] {
	self.Slice = src
	self.Reindex()
	return self
}

// Nullifies both the slice and the index. Does not preserve their capacity.
func (self *OrdSet[Val]) Clear() *OrdSet[Val] {
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
func (self *OrdSet[Val]) Reindex() { self.Index = SetOf(self.Slice...) }

// Implement `json.Marshaler`. Encodes the inner slice, ignoring the index.
func (self OrdSet[_]) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Slice)
}

// Unmarshals the input into the inner slice and rebuilds the index.
func (self *OrdSet[_]) UnmarshalJSON(src []byte) error {
	err := json.Unmarshal(src, &self.Slice)
	self.Reindex()
	return err
}
