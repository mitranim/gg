package gg

import "sort"

/*
Shortcut for making a slice out of the given values.
Workaround for the lack of type inference in type literals.
*/
func SliceOf[A any](val ...A) []A { return val }

/*
Shortcut for converting an arbitrary slice to `Slice`. Workaround for the
limitations of type inference in Go generics.
*/
func ToSlice[Src ~[]Elem, Elem any](val Src) Slice[Elem] {
	return Slice[Elem](val)
}

/*
Typedef of an arbitrary slice with various methods that duplicate global slice
functions such as `Get` or `Filter`. Useful as a shortcut for creating bound
methods that can be passed to higher-order functions. Example:

	values := []string{`one`, `two`, `three`}
	indexes := []int{0, 2}
	result := Map(indexes, ToSlice(values).Get)
	grepr.Println(result)
	// []string{`one`, `three`}
*/
type Slice[A any] []A

// Returns the underlying data pointer of the given slice.
func SliceDat[Slice ~[]Elem, Elem any](src Slice) *Elem {
	return CastUnsafe[*Elem](src)
}

// Same as global `SliceDat`.
func (self Slice[A]) Dat() *A { return SliceDat(self) }

// True if slice length is 0. The slice may or may not be nil.
func IsEmpty[Slice ~[]Elem, Elem any](val Slice) bool { return len(val) == 0 }

// Same as global `IsEmpty`.
func (self Slice[_]) IsEmpty() bool { return IsEmpty(self) }

// True if slice length is above 0.
func HasLen[Slice ~[]Elem, Elem any](val Slice) bool { return len(val) > 0 }

// Same as global `HasLen`.
func (self Slice[_]) HasLen() bool { return HasLen(self) }

// Same as `len(val)` but can be passed to higher-order functions.
func Len[Slice ~[]Elem, Elem any](val Slice) int { return len(val) }

// Same as global `Len`.
func (self Slice[_]) Len() int { return len(self) }

// Same as `len(PtrGet(val))` but can be passed to higher-order functions.
func PtrLen[Slice ~[]Elem, Elem any](val *Slice) int { return len(PtrGet(val)) }

// Same as global `PtrLen`.
func (self *Slice[_]) PtrLen() int { return PtrLen(self) }

// Same as `cap(val)` but can be passed to higher-order functions.
func Cap[Slice ~[]Elem, Elem any](val Slice) int { return cap(val) }

// Same as global `Cap`.
func (self Slice[_]) Cap() int { return cap(self) }

// Amount of unused capacity in the given slice.
func CapUnused[Slice ~[]Elem, Elem any](src Slice) int { return cap(src) - len(src) }

// Same as global `CapUnused`.
func (self Slice[_]) CapUnused() int { return CapUnused(self) }

/*
Amount of missing capacity that needs to be allocated to append the given amount
of additional elements.
*/
func CapMissing[Slice ~[]Elem, Elem any](src Slice, size int) int {
	return MaxPrim2(0, size-CapUnused(src))
}

// Same as global `CapMissing`.
func (self Slice[_]) CapMissing(size int) int { return CapMissing(self, size) }

// Counts the total length of the given slices.
func Lens[Slice ~[]Elem, Elem any](val ...Slice) int { return Sum(val, Len[Slice]) }

// Grows the length of the given slice by appending N zero values.
func GrowLen[Slice ~[]Elem, Elem any](src Slice, size int) Slice {
	return append(src, make(Slice, size)...)
}

// Same as global `GrowLen`.
func (self Slice[A]) GrowLen(size int) Slice[A] { return GrowLen(self, size) }

/*
Missing feature of the language / standard library. Ensures at least this much
unused capacity (not total capacity). If there is already enough capacity,
returns the slice as-is. Otherwise allocates a new slice, doubling the capacity
as many times as needed until it accommodates enough elements. Use this when
further growth is expected. When further growth is not expected, use
`GrowCapExact` instead. Similar to `(*bytes.Buffer).Grow` but works for all
slice types and avoids any wrapping, unwrapping, or spurious escapes to the
heap.
*/
func GrowCap[Slice ~[]Elem, Elem any](src Slice, size int) Slice {
	missing := CapMissing(src, size)
	if !(missing > 0) {
		return src
	}

	prev := MaxPrim2(0, cap(src))
	next := Or(prev, 1)
	for next < prev+size {
		next *= 2
	}

	out := make(Slice, len(src), next)
	copy(out, src)
	return out
}

// Same as global `GrowCap`.
func (self Slice[A]) GrowCap(size int) Slice[A] { return GrowCap(self, size) }

/*
Missing feature of the language / standard library. Ensures at least this much
unused capacity (not total capacity). If there is already enough capacity,
returns the slice as-is. Otherwise allocates a new slice with EXACTLY enough
additional capacity. Use this when further growth is not expected. When further
growth is expected, use `GrowCap` instead.
*/
func GrowCapExact[Slice ~[]Elem, Elem any](src Slice, size int) Slice {
	missing := CapMissing(src, size)
	if !(missing > 0) {
		return src
	}

	out := make(Slice, len(src), cap(src)+missing)
	copy(out, src)
	return out
}

// Same as global `GrowCapExact`.
func (self Slice[A]) GrowCapExact(size int) Slice[A] { return GrowCapExact(self, size) }

// Zeroes each element of the given slice.
func SliceZero[A any](val []A) {
	var zero A
	for ind := range val {
		val[ind] = zero
	}
}

// Same as global `SliceZero`.
func (self Slice[_]) Zero() { SliceZero(self) }

/*
Collapses the slice's length, preserving the capacity.
Does not modify any elements.
*/
func SliceTrunc[Slice ~[]Elem, Elem any](tar *Slice) {
	if tar != nil && *tar != nil {
		*tar = (*tar)[:0]
	}
}

// Same as global `SliceTrunc`.
func (self *Slice[_]) Trunc() { SliceTrunc(self) }

/*
If the index is within bounds, returns the value at that index and true.
Otherwise returns zero value and false.
*/
func Got[A any](src []A, ind int) (A, bool) {
	if ind >= 0 && ind < len(src) {
		return src[ind], true
	}
	return Zero[A](), false
}

// Same as global `Got`.
func (self Slice[A]) Got(ind int) (A, bool) { return Got(self, ind) }

/*
If the index is within bounds, returns the value at that index.
Otherwise returns zero value.
*/
func Get[A any](src []A, ind int) A {
	val, _ := Got(src, ind)
	return val
}

// Same as global `Get`.
func (self Slice[A]) Get(ind int) A { return Get(self, ind) }

/*
If the index is within bounds, returns a pointer to the value at that index.
Otherwise returns nil.
*/
func GetPtr[A any](src []A, ind int) *A {
	if ind >= 0 && ind < len(src) {
		return &src[ind]
	}
	return nil
}

// Same as global `GetPtr`.
func (self Slice[A]) GetPtr(ind int) *A { return GetPtr(self, ind) }

/*
Sets a value at an index, same as by using the built-in square bracket syntax.
Useful as a shortcut for inline bound functions.
*/
func (self Slice[A]) Set(ind int, val A) { self[ind] = val }

/*
Returns a shallow copy of the given slice. The capacity of the resulting slice
is equal to its length.
*/
func Clone[Slice ~[]Elem, Elem any](src Slice) Slice {
	if src == nil {
		return nil
	}

	out := make(Slice, len(src))
	copy(out, src)
	return out
}

// Same as global `Clone`.
func (self Slice[A]) Clone() Slice[A] { return Clone(self) }

/*
Same as `append`, but makes a copy instead of mutating the original.
Useful when reusing one "base" slice for in multiple append calls.
*/
func CloneAppend[Slice ~[]Elem, Elem any](src Slice, val ...Elem) Slice {
	if src == nil && val == nil {
		return nil
	}

	out := make(Slice, 0, len(src)+len(val))
	out = append(out, src...)
	out = append(out, val...)
	return out
}

// Same as global `CloneAppend`.
func (self Slice[A]) CloneAppend(val ...A) Slice[A] {
	return CloneAppend(self, val...)
}

/*
Appends the given elements to the given slice. Similar to built-in `append` but
syntactically easier because the destination is not repeated.
*/
func AppendVals[Slice ~[]Elem, Elem any](tar *Slice, val ...Elem) {
	if tar != nil {
		*tar = append(*tar, val...)
	}
}

// Same as global `AppendVals`.
func (self *Slice[A]) AppendVals(val ...A) { AppendVals(self, val...) }

/*
Appends the given element to the given slice, returning the pointer to the newly
appended position in the slice.
*/
func AppendPtr[Slice ~[]Elem, Elem any](tar *Slice, val Elem) *Elem {
	*tar = append(*tar, val)
	return LastPtr(*tar)
}

// Same as global `AppendPtr`.
func (self *Slice[A]) AppendPtr(val A) *A { return AppendPtr(self, val) }

/*
Appends a zero element to the given slice, returning the pointer to the newly
appended position in the slice.
*/
func AppendPtrZero[Slice ~[]Elem, Elem any](tar *Slice) *Elem {
	return AppendPtr(tar, Zero[Elem]())
}

// Same as global `AppendPtrZero`.
func (self *Slice[A]) AppendPtrZero() *A { return AppendPtrZero(self) }

/*
Returns the first element of the given slice. If the slice is empty, returns the
zero value.
*/
func Head[Slice ~[]Elem, Elem any](val Slice) Elem { return Get(val, 0) }

// Same as global `Head`.
func (self Slice[A]) Head() A { return Head(self) }

/*
Returns a pointer to the first element of the given slice. If the slice is
empty, the pointer is nil.
*/
func HeadPtr[Slice ~[]Elem, Elem any](val Slice) *Elem { return GetPtr(val, 0) }

// Same as global `HeadPtr`.
func (self Slice[A]) HeadPtr() *A { return HeadPtr(self) }

func PopHead[Slice ~[]Elem, Elem any](ptr *Slice) Elem {
	if ptr == nil {
		return Zero[Elem]()
	}

	head, tail := Head(*ptr), Tail(*ptr)
	*ptr = tail
	return head
}

// Same as global `PopHead`.
func (self *Slice[A]) PopHead() A { return PopHead(self) }

/*
Returns the last element of the given slice. If the slice is empty, returns the
zero value.
*/
func Last[Slice ~[]Elem, Elem any](val Slice) Elem { return Get(val, len(val)-1) }

// Same as global `Last`.
func (self Slice[A]) Last() A { return Last(self) }

/*
Returns a pointer to the last element of the given slice. If the slice is empty,
the pointer is nil.
*/
func LastPtr[Slice ~[]Elem, Elem any](val Slice) *Elem { return GetPtr(val, len(val)-1) }

// Same as global `LastPtr`.
func (self Slice[A]) LastPtr() *A { return LastPtr(self) }

/*
Returns the index of the last element of the given slice.
Same as `len(val)-1`. If slice is empty, returns -1.
*/
func LastIndex[Slice ~[]Elem, Elem any](val Slice) int {
	return len(val) - 1
}

// Same as global `LastIndex`.
func (self Slice[A]) LastIndex() int { return LastIndex(self) }

func PopLast[Slice ~[]Elem, Elem any](ptr *Slice) Elem {
	if ptr == nil {
		return Zero[Elem]()
	}

	init, last := Init(*ptr), Last(*ptr)
	*ptr = init
	return last
}

// Same as global `PopLast`.
func (self *Slice[A]) PopLast() A { return PopLast(self) }

/*
Returns the initial part of the given slice: all except the last value.
If the slice is nil, returns nil.
*/
func Init[Slice ~[]Elem, Elem any](val Slice) Slice {
	if len(val) == 0 {
		return val
	}
	return val[:len(val)-1]
}

// Same as global `Init`.
func (self Slice[A]) Init() Slice[A] { return Init(self) }

/*
Returns the tail part of the given slice: all except the first value.
If the slice is nil, returns nil.
*/
func Tail[Slice ~[]Elem, Elem any](val Slice) Slice {
	if len(val) == 0 {
		return val
	}
	return val[1:]
}

// Same as global `Tail`.
func (self Slice[A]) Tail() Slice[A] { return Tail(self) }

// Returns a subslice containing N elements from the start.
func Take[Slice ~[]Elem, Elem any](src Slice, size int) Slice {
	return src[:MaxPrim2(0, MinPrim2(size, len(src)))]
}

// Same as global `Take`.
func (self Slice[A]) Take(size int) Slice[A] { return Take(self, size) }

// Returns a subslice excluding N elements from the start.
func Drop[Slice ~[]Elem, Elem any](src Slice, size int) Slice {
	return src[MaxPrim2(0, MinPrim2(size, len(src))):]
}

// Same as global `Drop`.
func (self Slice[A]) Drop(size int) Slice[A] { return Drop(self, size) }

/*
Returns a subslice containing only elements at the start of the slice
for which the given function contiguously returned `true`.
*/
func TakeWhile[Slice ~[]Elem, Elem any](src Slice, fun func(Elem) bool) Slice {
	return Take(src, 1+FindIndex(src, func(val Elem) bool { return !fun(val) }))
}

// Same as global `TakeWhile`.
func (self Slice[A]) TakeWhile(fun func(A) bool) Slice[A] {
	return TakeWhile(self, fun)
}

/*
Returns a subslice excluding elements at the start of the slice
for which the given function contiguously returned `true`.
*/
func DropWhile[Slice ~[]Elem, Elem any](src Slice, fun func(Elem) bool) Slice {
	return Drop(src, 1+FindIndex(src, fun))
}

// Same as global `DropWhile`.
func (self Slice[A]) DropWhile(fun func(A) bool) Slice[A] {
	return DropWhile(self, fun)
}

// Calls the given function for each element of the given slice.
func Each[Slice ~[]Elem, Elem any](val Slice, fun func(Elem)) {
	if fun != nil {
		for _, val := range val {
			fun(val)
		}
	}
}

// Same as global `Each`.
func (self Slice[A]) Each(val Slice[A], fun func(A)) { Each(self, fun) }

/*
Calls the given function for each element's pointer in the given slice.
The pointer is always non-nil.
*/
func EachPtr[Slice ~[]Elem, Elem any](val Slice, fun func(*Elem)) {
	if fun != nil {
		for ind := range val {
			fun(&val[ind])
		}
	}
}

// Same as global `EachPtr`.
func (self Slice[A]) EachPtr(fun func(*A)) { EachPtr(self, fun) }

/*
Similar to `Each` but iterates two slices pairwise. If slice lengths don't
match, panics.
*/
func Each2[A, B any](one []A, two []B, fun func(A, B)) {
	validateLenMatch(len(one), len(two))

	if fun != nil {
		for ind := range one {
			fun(one[ind], two[ind])
		}
	}
}

/*
Returns the smallest value from among the inputs, which must be comparable
primitives. For non-primitives, see `Min`.
*/
func MinPrim[A LesserPrim](val ...A) A { return Fold1(val, MinPrim2[A]) }

/*
Returns the largest value from among the inputs, which must be comparable
primitives. For non-primitives, see `Max`.
*/
func MaxPrim[A LesserPrim](val ...A) A { return Fold1(val, MaxPrim2[A]) }

/*
Returns the smallest value from among the inputs. For primitive types that don't
implement `Lesser`, see `MinPrim`.
*/
func Min[A Lesser[A]](val ...A) A { return Fold1(val, Min2[A]) }

/*
Returns the largest value from among the inputs. For primitive types that don't
implement `Lesser`, see `MaxPrim`.
*/
func Max[A Lesser[A]](val ...A) A { return Fold1(val, Max2[A]) }

/*
Calls the given function for each element of the given slice and returns the
smallest value from among the results, which must be comparable primitives.
For non-primitives, see `MinBy`.
*/
func MinPrimBy[Src any, Out LesserPrim](src []Src, fun func(Src) Out) Out {
	if len(src) == 0 || fun == nil {
		return Zero[Out]()
	}

	return Fold(src[1:], fun(src[0]), func(acc Out, src Src) Out {
		return MinPrim2(acc, fun(src))
	})
}

/*
Calls the given function for each element of the given slice and returns the
smallest value from among the results. For primitive types that don't implement
`Lesser`, see `MinPrimBy`.
*/
func MinBy[Src any, Out Lesser[Out]](src []Src, fun func(Src) Out) Out {
	if len(src) == 0 || fun == nil {
		return Zero[Out]()
	}

	return Fold(src[1:], fun(src[0]), func(acc Out, src Src) Out {
		return Min2(acc, fun(src))
	})
}

/*
Calls the given function for each element of the given slice and returns the
largest value from among the results, which must be comparable primitives.
For non-primitives, see `MaxBy`.
*/
func MaxPrimBy[Src any, Out LesserPrim](src []Src, fun func(Src) Out) Out {
	if len(src) == 0 || fun == nil {
		return Zero[Out]()
	}

	return Fold(src[1:], fun(src[0]), func(acc Out, src Src) Out {
		return MaxPrim2(acc, fun(src))
	})
}

/*
Calls the given function for each element of the given slice and returns the
largest value from among the results. For primitive types that don't implement
`Lesser`, see `MaxPrimBy`.
*/
func MaxBy[Src any, Out Lesser[Out]](src []Src, fun func(Src) Out) Out {
	if len(src) == 0 || fun == nil {
		return Zero[Out]()
	}

	return Fold(src[1:], fun(src[0]), func(acc Out, src Src) Out {
		return Max2(acc, fun(src))
	})
}

// Combines all inputs via "+". If the input is empty, returns the zero value.
func Plus[A Plusable](val ...A) A { return Foldz(val, Plus2[A]) }

/*
Calls the given function on each element of the given slice and returns the sum
of all results, combined via "+".
*/
func Sum[Src any, Out Plusable](src []Src, fun func(Src) Out) Out {
	if fun == nil {
		return Zero[Out]()
	}
	return Foldz(src, func(acc Out, src Src) Out { return acc + fun(src) })
}

/*
Maps one slice to another. The resulting slice has exactly the same length as
the original. Each element is created by calling the given function on the
corresponding element of the original slice. The name refers to a well-known
functional programming abstraction which doesn't have anything in common with
the Go `map` types. Unlike many other higher-order slice functions, this one
requires a non-nil function.
*/
func Map[A, B any](src []A, fun func(A) B) []B {
	if src == nil {
		return nil
	}

	out := make([]B, 0, len(src))
	for _, val := range src {
		out = append(out, fun(val))
	}
	return out
}

/*
Similar to `Map` but instead of creating a new slice, appends to an existing
slice.
*/
func MapAppend[
	Src ~[]SrcVal,
	Tar ~[]TarVal,
	SrcVal any,
	TarVal any,
](tar *Tar, src Src, fun func(SrcVal) TarVal) {
	if tar == nil || fun == nil {
		return
	}

	buf := GrowCap(*tar, len(src))
	for _, val := range src {
		buf = append(buf, fun(val))
	}
	*tar = buf
}

/*
Similar to `Map`, but instead of creating a new slice, mutates the old one in
place by calling the given function on each element.
*/
func MapMut[Slice ~[]Elem, Elem any](src Slice, fun func(Elem) Elem) Slice {
	if fun != nil {
		for ind := range src {
			src[ind] = fun(src[ind])
		}
	}
	return src
}

// Same as global `MapMut`.
func (self Slice[A]) MapMut(fun func(A) A) Slice[A] { return MapMut(self, fun) }

/*
Similar to `Map` but iterates two slices pairwise, passing each element pair to
the mapping function. If slice lengths don't match, panics.
*/
func Map2[A, B, C any](one []A, two []B, fun func(A, B) C) []C {
	validateLenMatch(len(one), len(two))

	if one == nil || two == nil {
		return nil
	}

	out := make([]C, 0, len(one))
	for ind := range one {
		out = append(out, fun(one[ind], two[ind]))
	}
	return out
}

// Similar to `Map` but excludes any zero values produced by the given function.
func MapCompact[A, B any](src []A, fun func(A) B) []B {
	if fun == nil {
		return nil
	}

	var out []B
	for _, val := range src {
		val := fun(val)
		if !IsZero(val) {
			out = append(out, val)
		}
	}
	return out
}

// Similar to `Map` but concats the slices returned by the given function.
func MapFlat[A, B any](src []A, fun func(A) []B) []B {
	if src == nil {
		return nil
	}

	var out []B
	for _, val := range src {
		out = append(out, fun(val)...)
	}
	return out
}

/*
Takes a slice and "indexes" it by using keys generated by the given function,
returning the resulting map. If the function is nil, returns nil.
*/
func Index[Slice ~[]Val, Key comparable, Val any](src Slice, fun func(Val) Key) map[Key]Val {
	if fun == nil {
		return nil
	}

	out := make(map[Key]Val, len(src))
	IndexInto(out, src, fun)
	return out
}

/*
"Indexes" the given slice by adding its values to the given map, keyed by
calling the given function for each value. If the function is nil, does
nothing.
*/
func IndexInto[Key comparable, Val any](tar map[Key]Val, src []Val, fun func(Val) Key) {
	if fun == nil {
		return
	}
	for _, val := range src {
		tar[fun(val)] = val
	}
}

/*
Takes a slice and "indexes" it by converting each element to a key-value pair,
returning the resulting map.
*/
func IndexPair[
	Slice ~[]Elem,
	Elem any,
	Key comparable,
	Val any,
](
	src Slice, fun func(Elem) (Key, Val),
) map[Key]Val {
	if fun == nil {
		return nil
	}

	out := make(map[Key]Val, len(src))
	IndexPairInto(out, src, fun)
	return out
}

/*
"Indexes" the given slice by adding key-value pairs to the given map, making
key-value pairs by calling the given function for each element. If the
function is nil, does nothing.
*/
func IndexPairInto[Elem any, Key comparable, Val any](
	tar map[Key]Val,
	src []Elem,
	fun func(Elem) (Key, Val),
) {
	if fun == nil {
		return
	}

	for _, src := range src {
		key, val := fun(src)
		tar[key] = val
	}
}

/*
Groups the elements of the given slice by using keys generated by the given
function, returning the resulting map. If the function is nil, returns nil.
*/
func Group[Slice ~[]Val, Key comparable, Val any](src Slice, fun func(Val) Key) map[Key][]Val {
	if fun == nil {
		return nil
	}

	out := map[Key][]Val{}
	GroupInto(out, src, fun)
	return out
}

/*
Groups the elements of the given slice by adding its elements to slices in the
given map, keyed by calling the given function for each element. If the
function is nil, does nothing.
*/
func GroupInto[Key comparable, Val any](tar map[Key][]Val, src []Val, fun func(Val) Key) {
	if fun == nil {
		return
	}
	for _, val := range src {
		key := fun(val)
		tar[key] = append(tar[key], val)
	}
}

/*
Somewhat similar to `Map`. Creates a slice by "mapping" source values to
outputs. Calls the given function N times, passing an index, starting with 0.
*/
func Times[A any](src int, fun func(int) A) []A {
	if !(src > 0) || fun == nil {
		return nil
	}

	buf := make([]A, src)
	for ind := range buf {
		buf[ind] = fun(ind)
	}
	return buf
}

/*
Similar to `Times` but instead of creating a new slice, appends to an existing
slice.
*/
func TimesAppend[Slice ~[]Elem, Elem any](tar *Slice, src int, fun func(int) Elem) {
	if tar == nil || fun == nil || !(src > 0) {
		return
	}

	buf := GrowCap(*tar, src)
	for ind := range Iter(src) {
		buf = append(buf, fun(ind))
	}
	*tar = buf
}

// Same as global `TimesAppend`.
func (self *Slice[A]) TimesAppend(src int, fun func(int) A) {
	TimesAppend(self, src, fun)
}

// Counts the number of elements for which the given function returns true.
func Count[A any](src []A, fun func(A) bool) int {
	var out int
	if fun != nil {
		for _, src := range src {
			if fun(src) {
				out++
			}
		}
	}
	return out
}

// Same as global `Count`.
func (self Slice[A]) Count(src []A, fun func(A) bool) int { return Count(self, fun) }

/*
Folds the given slice by calling the given function for each element,
additionally passing the "accumulator". Returns the resulting accumulator.
*/
func Fold[Acc, Val any](src []Val, acc Acc, fun func(Acc, Val) Acc) Acc {
	if fun != nil {
		for _, val := range src {
			acc = fun(acc, val)
		}
	}
	return acc
}

/*
Short for "fold zero". Similar to `Fold` but the accumulator automatically
starts as the zero value of its type.
*/
func Foldz[Acc, Val any](src []Val, fun func(Acc, Val) Acc) Acc {
	return Fold(src, Zero[Acc](), fun)
}

/*
Similar to `Fold` but uses the first slice element as the accumulator, falling
back on zero value. The given function is invoked only for 2 or more elements.
*/
func Fold1[A any](src []A, fun func(A, A) A) A {
	if len(src) == 0 {
		return Zero[A]()
	}
	return Fold(src[1:], src[0], fun)
}

// Returns only the elements for which the given function returned `true`.
func Filter[Slice ~[]Elem, Elem any](src Slice, fun func(Elem) bool) (out Slice) {
	FilterAppend(&out, src, fun)
	return
}

// Same as global `Filter`.
func (self Slice[A]) Filter(fun func(A) bool) Slice[A] {
	return Filter(self, fun)
}

/*
Similar to `Filter` but instead of creating a new slice, appends to an existing
slice.
*/
func FilterAppend[Tar ~[]Elem, Elem any](tar *Tar, src []Elem, fun func(Elem) bool) {
	if tar == nil || fun == nil {
		return
	}

	for _, val := range src {
		if fun(val) {
			*tar = append(*tar, val)
		}
	}
}

// Same as global `FilterAppend`.
func (self *Slice[A]) FilterAppend(src []A, fun func(A) bool) {
	FilterAppend(self, src, fun)
}

/*
Inverse of `Filter`. Returns only the elements for which the given function
returned `false`.
*/
func Reject[Slice ~[]Elem, Elem any](src Slice, fun func(Elem) bool) (out Slice) {
	RejectAppend(&out, src, fun)
	return
}

// Same as global `Reject`.
func (self Slice[A]) Reject(fun func(A) bool) Slice[A] {
	return Reject(self, fun)
}

/*
Similar to `Reject` but instead of creating a new slice, appends to an existing
slice.
*/
func RejectAppend[Tar ~[]Elem, Elem any](tar *Tar, src []Elem, fun func(Elem) bool) {
	if tar == nil || fun == nil {
		return
	}

	for _, val := range src {
		if !fun(val) {
			*tar = append(*tar, val)
		}
	}
}

// Same as global `RejectAppend`.
func (self *Slice[A]) RejectAppend(src []A, fun func(A) bool) {
	RejectAppend(self, src, fun)
}

/*
Takes a slice and returns the indexes whose elements satisfy the given function.
All indexes are within the bounds of the original slice.
*/
func FilterIndex[Slice ~[]Elem, Elem any](src Slice, fun func(Elem) bool) []int {
	if fun == nil {
		return nil
	}

	var out []int
	for ind, val := range src {
		if fun(val) {
			out = append(out, ind)
		}
	}
	return out
}

// Same as global `FilterIndex`.
func (self Slice[A]) FilterIndex(fun func(A) bool) []int {
	return FilterIndex(self, fun)
}

/*
Takes a slice and returns the indexes whose elements are zero.
All indexes are within the bounds of the original slice.
*/
func ZeroIndex[Slice ~[]Elem, Elem any](src Slice) []int {
	return FilterIndex(src, IsZero[Elem])
}

// Same as global `ZeroIndex`.
func (self Slice[A]) ZeroIndex() []int { return ZeroIndex(self) }

/*
Takes a slice and returns the indexes whose elements are non-zero.
All indexes are within the bounds of the original slice.
*/
func NonZeroIndex[Slice ~[]Elem, Elem any](src Slice) []int {
	return FilterIndex(src, IsNonZero[Elem])
}

// Same as global `NonZeroIndex`.
func (self Slice[A]) NonZeroIndex() []int { return NonZeroIndex(self) }

// Returns a version of the given slice without any zero values.
func Compact[Slice ~[]Elem, Elem any](src Slice) Slice {
	return Filter(src, IsNonZero[Elem])
}

// Same as global `Compact`.
func (self Slice[A]) Compact() Slice[A] { return Compact(self) }

/*
Tests each element by calling the given function and returns the first element
for which it returns `true`. If none match, returns `-1`.
*/
func FindIndex[Slice ~[]Elem, Elem any](src Slice, fun func(Elem) bool) int {
	if fun != nil {
		for ind, val := range src {
			if fun(val) {
				return ind
			}
		}
	}
	return -1
}

// Same as global `FindIndex`.
func (self Slice[A]) FindIndex(fun func(A) bool) int {
	return FindIndex(self, fun)
}

/*
Returns the first element for which the given function returns `true`.
If nothing is found, returns a zero value. The additional boolean indicates
whether something was actually found.
*/
func Found[Slice ~[]Elem, Elem any](src Slice, fun func(Elem) bool) (Elem, bool) {
	ind := FindIndex(src, fun)
	if ind >= 0 {
		return src[ind], true
	}
	return Zero[Elem](), false
}

// Same as global `Found`.
func (self Slice[A]) Found(fun func(A) bool) (A, bool) {
	return Found(self, fun)
}

/*
Returns the first element for which the given function returns true.
If nothing is found, returns a zero value.
*/
func Find[Slice ~[]Elem, Elem any](src Slice, fun func(Elem) bool) Elem {
	return Get(src, FindIndex(src, fun))
}

// Same as global `Find`.
func (self Slice[A]) Find(fun func(A) bool) A { return Find(self, fun) }

/*
Similar to `Found`, but instead of returning an element, returns the first
product of the given function for which the returned boolean is true. If
nothing is procured, returns zero value and false.
*/
func Procured[Out, Src any](src []Src, fun func(Src) (Out, bool)) (Out, bool) {
	if fun != nil {
		for _, src := range src {
			val, ok := fun(src)
			if ok {
				return val, true
			}
		}
	}
	return Zero[Out](), false
}

/*
Similar to `Find`, but instead of returning the first approved element,
returns the first non-zero result of the given function. If nothing is
procured, returns a zero value.
*/
func Procure[Out, Src any](src []Src, fun func(Src) Out) Out {
	if fun != nil {
		for _, src := range src {
			val := fun(src)
			if IsNonZero(val) {
				return val
			}
		}
	}
	return Zero[Out]()
}

/*
Returns a version of the given slice with the given values appended unless they
were already present in the slice. This function only appends; it doesn't
deduplicate any previously existing values in the slice, nor reorder them.
*/
func Adjoin[Slice ~[]Elem, Elem comparable](tar Slice, src ...Elem) Slice {
	RejectAppend(&tar, src, SetOf(tar...).Has)
	return tar
}

/*
Returns a version of the given slice excluding any additionally supplied
values.
*/
func Exclude[Slice ~[]Elem, Elem comparable](base Slice, sub ...Elem) Slice {
	return Reject(base, SetOf(sub...).Has)
}

/*
Returns a version of the given slice excluding any additionally supplied
values.
*/
func Subtract[Slice ~[]Elem, Elem comparable](base Slice, sub ...Slice) Slice {
	return Reject(base, SetFrom(sub...).Has)
}

// Returns intersection of two slices: elements that occur in both.
func Intersect[Slice ~[]Elem, Elem comparable](one, two Slice) Slice {
	return Filter(one, SetOf(two...).Has)
}

/*
Combines the given slices, deduplicating their elements and preserving the order
of first occurrence for each element. As a special case, if the arguments
contain exactly one non-empty slice, it's returned as-is without deduplication.
To ensure uniqueness in all cases, call `Uniq`.
*/
func Union[Slice ~[]Elem, Elem comparable](val ...Slice) Slice {
	if Count(val, HasLen[Slice]) == 1 {
		return Find(val, HasLen[Slice])
	}

	var tar Slice
	var set Set[Elem]

	for _, val := range val {
		for _, val := range val {
			if set.Has(val) {
				continue
			}
			tar = append(tar, val)
			set.Init().Add(val)
		}
	}

	return tar
}

/*
Deduplicates the elements of the given slice, preserving the order of initial
occurrence for each element. The output is always a newly allocated slice.
*/
func Uniq[Slice ~[]Elem, Elem comparable](src Slice) Slice {
	var tar Slice
	var set Set[Elem]

	for _, val := range src {
		if set.Has(val) {
			continue
		}
		tar = append(tar, val)
		set.Init().Add(val)
	}

	return tar
}

/*
True if the given slice contains the given value. Should be used ONLY for very
small inputs: no more than a few tens of elements. For larger data, consider
using indexed data structures such as sets and maps.
*/
func Has[A comparable](src []A, val A) bool {
	return Some(src, func(elem A) bool { return elem == val })
}

/*
True if the first slice has all elements from the second slice. In other words,
true if A is a superset of B: A >= B.
*/
func HasEvery[A comparable](src, exp []A) bool {
	return Every(exp, SetOf(src...).Has)
}

/*
True if the first slice has some elements from the second slice. In other words,
true if the sets A and B intersect.
*/
func HasSome[A comparable](src, exp []A) bool {
	return Some(exp, SetOf(src...).Has)
}

/*
True if the first slice has NO elements from the second slice. In other words,
true if the sets A and B don't intersect.
*/
func HasNone[A comparable](src, exp []A) bool {
	return None(exp, SetOf(src...).Has)
}

/*
True if the given function returns true for any element of the given slice.
False if the function is nil. False if the slice is empty.
*/
func Some[A any](src []A, fun func(A) bool) bool {
	if fun == nil {
		return false
	}
	for _, val := range src {
		if fun(val) {
			return true
		}
	}
	return false
}

// Same as global `Some`.
func (self Slice[A]) Some(fun func(A) bool) bool { return Some(self, fun) }

// Inverse of `Some`.
func None[A any](src []A, fun func(A) bool) bool { return !Some(src, fun) }

// Same as global `None`.
func (self Slice[A]) None(fun func(A) bool) bool { return None(self, fun) }

/*
Utility for comparing slices pairwise. Returns true if the slices have the same
length and the function returns true for at least one pair.
*/
func SomePair[A any](one, two []A, fun func(A, A) bool) bool {
	if len(one) != len(two) || fun == nil {
		return false
	}
	for ind := range one {
		if fun(one[ind], two[ind]) {
			return true
		}
	}
	return false
}

/*
True if the given function returns true for every element of the given slice.
False if the function is nil. True if the slice is empty.
*/
func Every[A any](src []A, fun func(A) bool) bool {
	if fun == nil {
		return false
	}
	for _, val := range src {
		if !fun(val) {
			return false
		}
	}
	return true
}

// Same as global `Every`.
func (self Slice[A]) Every(fun func(A) bool) bool { return Every(self, fun) }

/*
Utility for comparing slices pairwise. Returns true if the slices have the same
length and the function returns true for every pair.
*/
func EveryPair[A any](one, two []A, fun func(A, A) bool) bool {
	if len(one) != len(two) || fun == nil {
		return false
	}
	for ind := range one {
		if !fun(one[ind], two[ind]) {
			return false
		}
	}
	return true
}

// Concatenates the inputs. If every input is nil, output is nil.
func Concat[Slice ~[]Elem, Elem any](val ...Slice) Slice {
	if Every(val, IsZero[Slice]) {
		return nil
	}

	buf := make(Slice, 0, Lens(val...))
	for _, val := range val {
		buf = append(buf, val...)
	}
	return buf
}

// Sorts a slice of comparable primitives. For non-primitives, see `Sort`.
func SortPrim[A LesserPrim](val []A) { SortablePrim[A](val).Sort() }

/*
Sorts a slice of comparable primitives, mutating and returning that slice.
For non-primitives, see `Sort`.
*/
func SortedPrim[Slice ~[]Elem, Elem LesserPrim](val Slice) Slice {
	return Slice(SortablePrim[Elem](val).Sorted())
}

// Implements `sort.Interface`.
type SortablePrim[A LesserPrim] []A

// Implement `sort.Interface`.
func (self SortablePrim[_]) Len() int { return len(self) }

// Implement `sort.Interface`.
func (self SortablePrim[_]) Less(one, two int) bool { return self[one] < self[two] }

// Implement `sort.Interface`.
func (self SortablePrim[_]) Swap(one, two int) { self[one], self[two] = self[two], self[one] }

// Sorts the receiver, mutating it.
func (self SortablePrim[_]) Sort() { sort.Stable(NoEscUnsafe(sort.Interface(self))) }

// Sorts the receiver, mutating and returning it.
func (self SortablePrim[A]) Sorted() SortablePrim[A] {
	self.Sort()
	return self
}

// Sorts a slice of comparable primitives. For primitives, see `SortPrim`.
func Sort[A Lesser[A]](val []A) { Sortable[A](val).Sort() }

/*
Sorts a slice of comparable values, mutating and returning that slice.
For primitives, see `SortedPrim`.
*/
func Sorted[Slice ~[]Elem, Elem Lesser[Elem]](val Slice) Slice {
	return Slice(Sortable[Elem](val).Sorted())
}

// Implements `sort.Interface`.
type Sortable[A Lesser[A]] []A

// Implement `sort.Interface`.
func (self Sortable[_]) Len() int { return len(self) }

// Implement `sort.Interface`.
func (self Sortable[_]) Less(one, two int) bool { return self[one].Less(self[two]) }

// Implement `sort.Interface`.
func (self Sortable[_]) Swap(one, two int) { self[one], self[two] = self[two], self[one] }

// Sorts the receiver, mutating it.
func (self Sortable[_]) Sort() { sort.Stable(NoEscUnsafe(sort.Interface(self))) }

// Sorts the receiver, mutating and returning it.
func (self Sortable[A]) Sorted() Sortable[A] {
	self.Sort()
	return self
}

// Reverses the given slice in-place, mutating it.
func Reverse[A any](val []A) {
	ind0 := 0
	ind1 := len(val) - 1

	for ind0 < ind1 {
		val[ind0], val[ind1] = val[ind1], val[ind0]
		ind0++
		ind1--
	}
}

// Same as global `Reverse`.
func (self Slice[_]) Reverse() { Reverse(self) }

// Reverses the given slice in-place, mutating it and returning that slice.
func Reversed[Slice ~[]Elem, Elem any](val Slice) Slice {
	Reverse(val)
	return val
}

// Same as global `Reversed`.
func (self Slice[A]) Reversed() Slice[A] { return Reversed(self) }
