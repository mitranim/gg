package gg

import r "reflect"

func cloneArray(src r.Value) {
	if src.Cap() == 0 || !IsIndirect(src.Type().Elem()) {
		return
	}

	for ind := range Iter(src.Len()) {
		ValueClone(src.Index(ind))
	}
}

func clonedArray(src r.Value) r.Value {
	if src.Cap() == 0 || !IsIndirect(src.Type().Elem()) {
		return src
	}

	out := NewElem(src.Type())
	r.Copy(out, src)
	cloneArray(out)
	return out
}

/**
Known defect: when cloning, in addition to allocating a new backing array, this
allocates a slice header, which could theoretically be avoided if we could make
just a backing array of the required size and replace the array pointer in the
slice header we already have.
*/
func cloneSlice(src r.Value) { ValueSet(src, clonedSlice(src)) }

func clonedSlice(src r.Value) r.Value {
	if src.IsNil() || src.Cap() == 0 {
		return src
	}

	out := r.MakeSlice(src.Type(), src.Len(), src.Cap())
	r.Copy(out, src)
	cloneArray(out)
	return out
}

func cloneInterface(src r.Value) { ValueSet(src, clonedInterface(src)) }

func clonedInterface(src r.Value) r.Value {
	if src.IsNil() {
		return src
	}

	elem0 := src.Elem()
	elem1 := ValueCloned(elem0)
	if elem0 == elem1 {
		return elem0
	}
	return elem1.Convert(src.Type())
}

func cloneMap(src r.Value) { ValueSet(src, clonedMap(src)) }

func clonedMap(src r.Value) r.Value {
	if src.IsNil() {
		return src
	}

	out := r.MakeMapWithSize(src.Type(), src.Len())
	iter := src.MapRange()
	for iter.Next() {
		out.SetMapIndex(ValueCloned(iter.Key()), ValueCloned(iter.Value()))
	}
	return out
}

func clonePointer(src r.Value) { ValueSet(src, clonedPointer(src)) }

func clonedPointer(src r.Value) r.Value {
	if src.IsNil() {
		return src
	}

	out := r.New(src.Type().Elem())
	out.Elem().Set(src.Elem())
	ValueClone(out.Elem())
	return out
}

func cloneStruct(src r.Value) {
	for _, field := range StructPublicFieldCache.Get(src.Type()) {
		ValueClone(src.FieldByIndex(field.Index))
	}
}

func clonedStruct(src r.Value) r.Value {
	if !IsIndirect(src.Type()) {
		return src
	}

	out := NewElem(src.Type())
	out.Set(src)
	cloneStruct(out)
	return out
}
