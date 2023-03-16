package gg

import "encoding/json"

/*
Syntactic shortcut for making a set from a slice, with element type inference
and capacity preallocation. Always returns non-nil, even if the input is
empty.
*/
func SetOf[A comparable](val ...A) Set[A] {
	return make(Set[A], len(val)).Add(val...)
}

/*
Syntactic shortcut for making a set from multiple slices, with element type
inference and capacity preallocation. Always returns non-nil, even if the input
is empty.
*/
func SetFrom[Slice ~[]Elem, Elem comparable](val ...Slice) Set[Elem] {
	buf := make(Set[Elem], Lens(val...))
	for _, val := range val {
		buf.Add(val...)
	}
	return buf
}

/*
Creates a set by "mapping" the elements of a given slice via the provided
function. Always returns non-nil, even if the input is empty.
*/
func SetMapped[Elem any, Val comparable](src []Elem, fun func(Elem) Val) Set[Val] {
	buf := make(Set[Val], len(src))
	if fun != nil {
		for _, val := range src {
			buf[fun(val)] = struct{}{}
		}
	}
	return buf
}

// Generic unordered set backed by a map.
type Set[A comparable] map[A]struct{}

/*
Idempotently inits the map via `make`, making it writable. The output pointer
must be non-nil.
*/
func (self *Set[A]) Init() Set[A] {
	if *self == nil {
		*self = make(Set[A])
	}
	return *self
}

// Same as `len(set)`. Nil-safe.
func (self Set[_]) Len() int { return len(self) }

// True if the set includes the given value. Nil-safe.
func (self Set[A]) Has(val A) bool { return MapHas(self, val) }

// Idempotently adds the given values to the receiver, which must be non-nil.
func (self Set[A]) Add(val ...A) Set[A] {
	for _, val := range val {
		self[val] = struct{}{}
	}
	return self
}

/*
Set union. Idempotently adds all values from the given source sets to the
receiver, which must be non-nil.
*/
func (self Set[A]) AddFrom(val ...Set[A]) Set[A] {
	for _, val := range val {
		for val := range val {
			self[val] = struct{}{}
		}
	}
	return self
}

// Deletes the given values from the receiver, which may be nil.
func (self Set[A]) Del(val ...A) Set[A] {
	for _, val := range val {
		delete(self, val)
	}
	return self
}

/*
Deletes all values present in the given source sets from the receiver, which may
be nil.
*/
func (self Set[A]) DelFrom(val ...Set[A]) Set[A] {
	for _, val := range val {
		for val := range val {
			delete(self, val)
		}
	}
	return self
}

/*
Clears and returns the receiver, which may be nil. Note that this type is
implemented as a map, and this method involves iterating the map, which is
inefficient in Go. In many cases, it's more efficient to make a new set.
*/
func (self Set[A]) Clear() Set[A] {
	for val := range self {
		delete(self, val)
	}
	return self
}

// Combination of `Set.Clear` and `Set.Add`.
func (self Set[A]) Reset(val ...A) Set[A] {
	self.Clear()
	self.Add(val...)
	return self
}

// Converts the map to a slice of its values. Order is random.
func (self Set[A]) Slice() []A { return MapKeys(self) }

/*
Returns the subset of values for which the given function returns true.
Order is random. If function is nil, output is nil.
*/
func (self Set[A]) Filter(fun func(A) bool) []A {
	var out []A
	if fun != nil {
		for val := range self {
			if fun(val) {
				out = append(out, val)
			}
		}
	}
	return out
}

// JSON-encodes as a list. Order is random.
func (self Set[A]) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Slice())
}

/*
JSON-decodes the input, which must either represent JSON "null" or a JSON list
of values compatible with the value type.
*/
func (self *Set[A]) UnmarshalJSON(src []byte) error {
	var buf []A
	err := json.Unmarshal(src, &buf)
	if err != nil {
		return err
	}

	self.Init().Reset(buf...)
	return nil
}

// Implement `fmt.GoStringer`, returning valid Go code that constructs the set.
func (self Set[A]) GoString() string {
	typ := TypeOf(self).String()

	if self == nil {
		return typ + `(nil)`
	}

	if !(len(self) > 0) {
		return typ + `{}`
	}

	var buf Buf
	buf.AppendString(typ)
	buf.AppendString(`{}.Add(`)

	var found bool
	for val := range self {
		if found {
			buf.AppendString(`, `)
		}
		found = true
		buf.AppendGoString(val)
	}

	buf.AppendString(`)`)
	return buf.String()
}
