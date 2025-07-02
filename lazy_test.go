package gg_test

import (
	"sync"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

// TODO: test with concurrency.
func TestOnce(t *testing.T) {
	defer gtest.Catch(t)

	var count int

	fun := gg.Once(func() int {
		count++
		return count
	})

	gtest.Eq(fun(), 1)
	gtest.Eq(count, 1)

	gtest.Eq(fun(), 1)
	gtest.Eq(count, 1)

	gtest.Eq(fun(), 1)
	gtest.Eq(count, 1)
}

func TestOnce_panic_retry(t *testing.T) {
	defer gtest.Catch(t)

	var count int

	fun := gg.Once(func() int {
		count++
		if count <= 3 {
			panic(`intermittent_failure`)
		}
		return 123
	})

	gtest.PanicStr(`intermittent_failure`, func() { fun() })
	gtest.PanicStr(`intermittent_failure`, func() { fun() })
	gtest.PanicStr(`intermittent_failure`, func() { fun() })
	gtest.Eq(fun(), 123)
	gtest.Eq(fun(), 123)
	gtest.Eq(fun(), 123)
}

func BenchmarkOnce_make(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Once(gg.Cwd))
	}
}

func BenchmarkOnce_call(b *testing.B) {
	once := gg.Once(gg.Cwd)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(once())
	}
}

func Benchmark_sync_OnceValue_make(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(sync.OnceValue(gg.Cwd))
	}
}

func Benchmark_sync_OnceValue_call(b *testing.B) {
	once := sync.OnceValue(gg.Cwd)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(once())
	}
}

// TODO: test with concurrency.
func TestLazy(t *testing.T) {
	defer gtest.Catch(t)

	var count int

	once := gg.NewLazy(func() int {
		count++
		if count > 1 {
			panic(gg.Errf(`excessive count %v`, count))
		}
		return count
	})

	gtest.Eq(*gg.CastUnsafe[*int](once), 0)
	gtest.Eq(once.Get(), 1)
	gtest.Eq(once.Get(), once.Get())
	gtest.Eq(*gg.CastUnsafe[*int](once), 1)
}

func BenchmarkLazy(b *testing.B) {
	once := gg.NewLazy(gg.Cwd)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(once.Get())
	}
}
