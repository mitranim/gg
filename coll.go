package gg

import "encoding/json"

/*
Short for "valid primary key". Returns the primary key generated by the given
input, asserts that the key is non-zero, and returns the resulting key.
Used internally by `Coll` and `LazyColl`.
*/
func ValidPk[Key comparable, Val Pker[Key]](val Val) Key {
	key := val.Pk()
	if IsZero(key) {
		panic(Errf(`unexpected empty key %v in %v`, Type[Key](), Type[Val]()))
	}
	return key
}

/*
Syntactic shortcut for making a `Coll` of the given arguments. Reuses the given
slice as-is with no reallocation.
*/
func CollOf[Key comparable, Val Pker[Key]](src ...Val) Coll[Key, Val] {
	var tar Coll[Key, Val]
	tar.Reset(src...)
	return tar
}

/*
Syntactic shortcut for making a `Coll` from any number of source slices. When
called with exactly one argument, this reuses the given slice as-is with no
reallocation.
*/
func CollFrom[Key comparable, Val Pker[Key], Slice ~[]Val](src ...Slice) Coll[Key, Val] {
	var tar Coll[Key, Val]

	switch len(src) {
	case 1:
		tar.Reset(src[0]...)
	default:
		for _, src := range src {
			tar.Add(src...)
		}
	}

	return tar
}

/*
Short for "collection". Represents an ordered map where keys are automatically
derived from values. Compare `OrdMap` where keys are provided externally. Keys
must be non-zero. Similarly to a map, this ensures value uniqueness by primary
key, and allows efficient access by key. Unlike a map, values in this type are
ordered and can be iterated cheaply, because they are stored in a
publicly-accessible slice. However, as a tradeoff, this type does not support
deletion.
*/
type Coll[Key comparable, Val Pker[Key]] OrdMap[Key, Val]

// Same as `OrdMap.Len`.
func (self Coll[_, _]) Len() int { return self.OrdMap().Len() }

// Same as `OrdMap.IsEmpty`.
func (self Coll[_, _]) IsEmpty() bool { return self.OrdMap().IsEmpty() }

// Same as `OrdMap.IsNotEmpty`.
func (self Coll[_, _]) IsNotEmpty() bool { return self.OrdMap().IsNotEmpty() }

// Same as `OrdMap.Has`.
func (self Coll[Key, _]) Has(key Key) bool { return self.OrdMap().Has(key) }

// Same as `OrdMap.Get`.
func (self Coll[Key, Val]) Get(key Key) Val { return self.OrdMap().Get(key) }

// Same as `OrdMap.GetReq`.
func (self Coll[Key, Val]) GetReq(key Key) Val { return self.OrdMap().GetReq(key) }

// Same as `OrdMap.Got`.
func (self Coll[Key, Val]) Got(key Key) (Val, bool) { return self.OrdMap().Got(key) }

// Same as `OrdMap.Ptr`.
func (self Coll[Key, Val]) Ptr(key Key) *Val { return self.OrdMap().Ptr(key) }

// Same as `OrdMap.PtrReq`.
func (self Coll[Key, Val]) PtrReq(key Key) *Val { return self.OrdMap().PtrReq(key) }

/*
Idempotently adds each given value to both the inner slice and the inner index.
Every value whose key already exists in the index is replaced at the existing
position in the slice.
*/
func (self *Coll[Key, Val]) Add(src ...Val) *Coll[Key, Val] {
	for _, src := range src {
		self.OrdMap().Set(ValidPk[Key](src), src)
	}
	return self
}

/*
Same as `Coll.Add`, but panics if any inputs are redundant, as in, their primary
keys are already present in the index.
*/
func (self *Coll[Key, Val]) AddUniq(src ...Val) *Coll[Key, Val] {
	for _, src := range src {
		self.OrdMap().Add(ValidPk[Key](src), src)
	}
	return self
}

// Same as `OrdMap.Clear`.
func (self *Coll[Key, Val]) Clear() *Coll[Key, Val] {
	self.OrdMap().Clear()
	return self
}

/*
Replaces `.Slice` with the given slice and rebuilds `.Index`. Uses the slice
as-is with no reallocation. Callers must be careful to avoid modifying the
source data, which may invalidate the collection's index.
*/
func (self *Coll[Key, Val]) Reset(src ...Val) *Coll[Key, Val] {
	self.Slice = src
	self.Reindex()
	return self
}

/*
Rebuilds the inner index from the inner slice, without checking the validity of
the existing index. Can be useful for external code that directly modifies the
inner `.Slice`, for example by sorting it. This is NOT used when adding items
via `.Add`, which modifies the index incrementally rather than all-at-once.
*/
func (self *Coll[Key, Val]) Reindex() *Coll[Key, Val] {
	slice := self.Slice
	if len(slice) <= 0 {
		self.Index = nil
		return self
	}

	index := make(map[Key]int, len(slice))
	for ind, val := range slice {
		index[ValidPk[Key](val)] = ind
	}
	self.Index = index

	return self
}

/*
Swaps two elements both in `.Slice` and in `.Index`. Useful for sorting.
`.Index` may be nil, in which case it's unaffected. Slice indices must be
either equal or valid.
*/
func (self Coll[Key, _]) Swap(ind0, ind1 int) {
	if ind0 == ind1 {
		return
	}

	slice := self.Slice
	val0, val1 := slice[ind0], slice[ind1]
	slice[ind0], slice[ind1] = val1, val0

	index := self.Index
	if index != nil {
		index[ValidPk[Key](val0)], index[ValidPk[Key](val1)] = ind1, ind0
	}
}

// Implement `json.Marshaler`. Encodes the inner slice, ignoring the index.
func (self Coll[_, _]) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Slice)
}

/*
Implement `json.Unmarshaler`. Decodes the input into the inner slice and
rebuilds the index.
*/
func (self *Coll[_, _]) UnmarshalJSON(src []byte) error {
	err := json.Unmarshal(src, &self.Slice)
	self.Reindex()
	return err
}

/*
Free cast into the equivalent `*OrdMap`. Note that mutating the resulting
`OrdMap` via methods such as `OrdMap.Add` may violate guarantees of the `Coll`
type, mainly that each value is stored under the key returned by its `.Pk`
method.
*/
func (self *Coll[Key, Val]) OrdMap() *OrdMap[Key, Val] {
	return (*OrdMap[Key, Val])(self)
}

// Free cast to equivalent `LazyColl`.
func (self *Coll[Key, Val]) LazyColl() *LazyColl[Key, Val] {
	return (*LazyColl[Key, Val])(self)
}
