package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

// Placeholder, needs a concurrency test.
func TestAtom(t *testing.T) {
	defer gtest.Catch(t)

	var ref gg.Atom[string]
	gtest.Zero(ref.Load())
	gtest.Eq(
		gg.Tuple2(ref.Loaded()),
		gg.Tuple2(``, false),
	)

	ref.Store(``)
	gtest.Zero(ref.Load())
	gtest.Eq(
		gg.Tuple2(ref.Loaded()),
		gg.Tuple2(``, true),
	)

	ref.Store(`one`)
	gtest.Eq(ref.Load(), `one`)
	gtest.Eq(
		gg.Tuple2(ref.Loaded()),
		gg.Tuple2(`one`, true),
	)

	gtest.False(ref.CompareAndSwap(`three`, `two`))
	gtest.Eq(ref.Load(), `one`)

	gtest.True(ref.CompareAndSwap(`one`, `two`))
	gtest.Eq(ref.Load(), `two`)
}

func BenchmarkAtom_Store(b *testing.B) {
	var ref gg.Atom[string]

	for ind := 0; ind < b.N; ind++ {
		ref.Store(`str`)
	}
}

func BenchmarkAtom_Load(b *testing.B) {
	var ref gg.Atom[string]
	ref.Store(`str`)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(ref.Load())
	}
}

func TestChanInit(t *testing.T) {
	defer gtest.Catch(t)

	gg.ChanInit((*chan string)(nil))

	var tar chan string
	gtest.Eq(gg.ChanInit(&tar), tar)

	gtest.NotZero(tar)
	gtest.Eq(cap(tar), 0)

	prev := tar
	gtest.Eq(gg.ChanInit(&tar), prev)
	gtest.Eq(tar, prev)
}

func TestChanInitCap(t *testing.T) {
	defer gtest.Catch(t)

	gg.ChanInitCap((*chan string)(nil), 1)

	var tar chan string
	gtest.Eq(gg.ChanInitCap(&tar, 3), tar)

	gtest.NotZero(tar)
	gtest.Eq(cap(tar), 3)

	prev := tar
	gtest.Eq(gg.ChanInitCap(&tar, 5), prev)
	gtest.Eq(tar, prev)
	gtest.Eq(cap(prev), 3)
	gtest.Eq(cap(tar), 3)
}

func TestSendOpt(t *testing.T) {
	defer gtest.Catch(t)

	var tar chan string
	gg.SendOpt(tar, `val`)
	gg.SendOpt(tar, `val`)
	gg.SendOpt(tar, `val`)

	tar = make(chan string, 1)
	gg.SendOpt(tar, `val`)
	gg.SendOpt(tar, `val`)
	gg.SendOpt(tar, `val`)

	gtest.Eq(<-tar, `val`)
}

func TestSendZeroOpt(t *testing.T) {
	defer gtest.Catch(t)

	var tar chan string
	gg.SendZeroOpt(tar)
	gg.SendZeroOpt(tar)
	gg.SendZeroOpt(tar)

	tar = make(chan string, 1)
	gg.SendZeroOpt(tar)
	gg.SendZeroOpt(tar)
	gg.SendZeroOpt(tar)

	val, ok := <-tar
	gtest.Zero(val)
	gtest.True(ok)
}
