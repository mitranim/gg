package gg

import "encoding/json"

/*
Same as `CollOf` but for `LazyColl`. Note that while the return type is a
non-pointer for easy assignment, callers must always access `LazyColl` by
pointer to avoid redundant reindexing.
*/
func LazyCollOf[Key comparable, Val Pker[Key]](src ...Val) LazyColl[Key, Val] {
	var tar LazyColl[Key, Val]
	tar.Reset(src...)
	return tar
}

/*
Same as `CollFrom` but for `LazyColl`. Note that while the return type is a
non-pointer for easy assignment, callers must always access `LazyColl` by
pointer to avoid redundant reindexing.
*/
func LazyCollFrom[Key comparable, Val Pker[Key], Slice ~[]Val](src ...Slice) LazyColl[Key, Val] {
	var tar LazyColl[Key, Val]

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
Short for "lazy collection". Variant of `Coll` where the index is built lazily
rather than immediately. This is not the default behavior in `Coll` because it
requires various access methods such as `.Has` and `.Get` to be defined on the
pointer type rather than value type, and more importantly, it's more error
prone: the caller is responsible for making sure that the collection is always
accessed by pointer, never by value, to avoid redundant reindexing.
*/
type LazyColl[Key comparable, Val Pker[Key]] Coll[Key, Val]

// Same as `Coll.Len`.
func (self LazyColl[_, _]) Len() int { return self.coll().Len() }

// Same as `Coll.IsEmpty`.
func (self LazyColl[_, _]) IsEmpty() bool { return self.coll().IsEmpty() }

// Same as `Coll.IsNotEmpty`.
func (self LazyColl[_, _]) IsNotEmpty() bool { return self.coll().IsNotEmpty() }

// Same as `Coll.Has`. Lazily rebuilds the index if necessary.
func (self *LazyColl[Key, _]) Has(key Key) bool {
	self.ReindexOpt()
	return self.coll().Has(key)
}

// Same as `Coll.Get`. Lazily rebuilds the index if necessary.
func (self *LazyColl[Key, Val]) Get(key Key) Val {
	self.ReindexOpt()
	return self.coll().Get(key)
}

// Same as `Coll.GetReq`. Lazily rebuilds the index if necessary.
func (self *LazyColl[Key, Val]) GetReq(key Key) Val {
	self.ReindexOpt()
	return self.coll().GetReq(key)
}

// Same as `Coll.Got`. Lazily rebuilds the index if necessary.
func (self *LazyColl[Key, Val]) Got(key Key) (Val, bool) {
	self.ReindexOpt()
	return self.coll().Got(key)
}

// Same as `Coll.Ptr`. Lazily rebuilds the index if necessary.
func (self *LazyColl[Key, Val]) Ptr(key Key) *Val {
	self.ReindexOpt()
	return self.coll().Ptr(key)
}

// Same as `Coll.PtrReq`. Lazily rebuilds the index if necessary.
func (self *LazyColl[Key, Val]) PtrReq(key Key) *Val {
	self.ReindexOpt()
	return self.coll().PtrReq(key)
}

// Similar to `Coll.Add`, but does not add new entries to the index.
func (self *LazyColl[Key, Val]) Add(src ...Val) *LazyColl[Key, Val] {
	for _, val := range src {
		key := ValidPk[Key](val)
		ind, ok := self.Index[key]
		if ok {
			self.Slice[ind] = val
			continue
		}
		Append(&self.Slice, val)
	}
	return self
}

// Same as `Coll.Reset` but deletes the index instead of rebuilding it.
func (self *LazyColl[Key, Val]) Reset(src ...Val) *LazyColl[Key, Val] {
	self.Index = nil
	self.Slice = src
	return self
}

// Same as `Coll.Clear`.
func (self *LazyColl[Key, Val]) Clear() *LazyColl[Key, Val] {
	self.coll().Clear()
	return self
}

// Same as `Coll.Reindex`.
func (self *LazyColl[Key, Val]) Reindex() *LazyColl[Key, Val] {
	self.coll().Reindex()
	return self
}

/*
Rebuilds the index if the length of inner slice and index doesn't match.
This is used internally by all "read" methods on this type.
*/
func (self *LazyColl[Key, _]) ReindexOpt() {
	if len(self.Slice) != len(self.Index) {
		self.Reindex()
	}
}

// Same as `Coll.Swap` but deletes the index instead of modifying it.
func (self *LazyColl[Key, _]) Swap(ind0, ind1 int) {
	self.Index = nil
	self.coll().Swap(ind0, ind1)
}

/*
Implement `json.Marshaler`. Same as `Coll.MarshalJSON`. Encodes the inner slice,
ignoring the index.
*/
func (self LazyColl[_, _]) MarshalJSON() ([]byte, error) {
	return self.coll().MarshalJSON()
}

/*
Implement `json.Unmarshaler`. Similar to `Coll.UnmarshalJSON`, but after
decoding the input into the inner slice, deletes the index instead of
rebuilding it.
*/
func (self *LazyColl[_, _]) UnmarshalJSON(src []byte) error {
	self.Index = nil
	return json.Unmarshal(src, &self.Slice)
}

// Converts to equivalent `Coll`. Lazily rebuilds the index if necessary.
func (self *LazyColl[Key, Val]) Coll() *Coll[Key, Val] {
	self.ReindexOpt()
	return (*Coll[Key, Val])(self)
}

/*
Free cast to equivalent `Coll`. Private because careless use produces invalid
behavior. Incomplete index in `LazyColl` violates assumptions made by `Coll`.
*/
func (self *LazyColl[Key, Val]) coll() *Coll[Key, Val] {
	return (*Coll[Key, Val])(self)
}
