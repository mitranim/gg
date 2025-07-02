package gg_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

// Semi-placeholder: needs a concurrency test.
func TestAtom(t *testing.T) {
	defer gtest.Catch(t)

	var ref gg.Atom[string]

	gtest.Zero(ref.LoadVal())
	gtest.Zero(ref.LoadPtr())

	gtest.Eq(
		gg.Tuple2(ref.LoadedVal()),
		gg.Tuple2(``, false),
	)

	val0 := ``
	ref.StorePtr(&val0)

	gtest.Eq(ref.LoadPtr(), &val0)
	gtest.Eq(ref.LoadVal(), ``)

	gtest.Eq(
		gg.Tuple2(ref.LoadedVal()),
		gg.Tuple2(``, true),
	)

	ref.StoreVal(val0)
	gtest.NotEq(ref.LoadPtr(), &val0)
	gtest.Eq(ref.LoadVal(), ``)

	gtest.Eq(
		gg.Tuple2(ref.LoadedVal()),
		gg.Tuple2(``, true),
	)

	val1 := `one`
	ref.StorePtr(&val1)

	gtest.Eq(ref.LoadPtr(), &val1)
	gtest.Eq(ref.LoadVal(), `one`)

	gtest.Eq(
		gg.Tuple2(ref.LoadedVal()),
		gg.Tuple2(`one`, true),
	)

	ref.StoreVal(val1)
	gtest.NotEq(ref.LoadPtr(), &val0)
	gtest.NotEq(ref.LoadPtr(), &val1)
	gtest.Eq(ref.LoadVal(), `one`)

	gtest.Eq(
		gg.Tuple2(ref.LoadedVal()),
		gg.Tuple2(`one`, true),
	)

	val2 := `two`

	gtest.False(ref.CompareAndSwapPtr(&val1, &val2))
	gtest.NotEq(ref.LoadPtr(), &val0)
	gtest.NotEq(ref.LoadPtr(), &val1)
	gtest.NotEq(ref.LoadPtr(), &val2)
	gtest.Eq(ref.LoadVal(), `one`)

	ref.StorePtr(&val1)
	gtest.True(ref.CompareAndSwapPtr(&val1, &val2))
	gtest.Eq(ref.LoadPtr(), &val2)
	gtest.Eq(ref.LoadVal(), `two`)

	gtest.True(ref.CompareAndSwapPtr(&val2, &val2))
	gtest.Eq(ref.LoadPtr(), &val2)
	gtest.Eq(ref.LoadVal(), `two`)

	gtest.Eq(val0, ``)
	gtest.Eq(val1, `one`)
	gtest.Eq(val2, `two`)

	val2 = `three`
	gtest.Eq(ref.LoadPtr(), &val2)
	gtest.Eq(ref.LoadVal(), `three`)
}

func TestAtom_iface(t *testing.T) {
	defer gtest.Catch(t)

	var ref gg.Atom[fmt.Stringer]

	gtest.Zero(ref.LoadPtr())
	gtest.Zero(ref.LoadVal())

	var val0 fmt.Stringer = gg.Buf(`one`)

	ref.StoreVal(val0)
	gtest.Is(ref.LoadVal(), val0)
	gtest.Eq(ref.LoadVal().String(), `one`)

	var val1 fmt.Stringer = gg.ErrStr(`two`)

	ref.StoreVal(val1)
	gtest.Is(ref.LoadVal(), val1)
	gtest.Eq(ref.LoadVal().String(), `two`)
}

func BenchmarkAtom_StorePtr(b *testing.B) {
	var ref gg.Atom[time.Time]
	val := time.Now()

	for ind := 0; ind < b.N; ind++ {
		ref.StorePtr(&val)
	}
}

func BenchmarkAtom_StoreVal(b *testing.B) {
	var ref gg.Atom[time.Time]
	val := time.Now()

	for ind := 0; ind < b.N; ind++ {
		ref.StoreVal(val)
	}
}

func BenchmarkAtom_LoadPtr(b *testing.B) {
	var ref gg.Atom[time.Time]
	ref.StoreVal(time.Now())

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(ref.LoadPtr())
	}
}

func BenchmarkAtom_LoadVal(b *testing.B) {
	var ref gg.Atom[time.Time]
	ref.StoreVal(time.Now())

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(ref.LoadVal())
	}
}
