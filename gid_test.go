package gg_test

import (
	"math"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestGid(t *testing.T) {
	defer gtest.Catch(t)
	testGid(t, gg.Gid)
}

func TestGid_slow(t *testing.T) {
	defer gtest.Catch(t)
	testGid(t, gg.GidSlow)
}

func testGid(t *testing.T, fun func() uint64) {
	id0 := fun()
	gtest.NotZero(id0)

	gtest.Eq(fun(), id0)
	gtest.Eq(fun(), id0)
	gtest.LessPrim(0, id0)

	/**
	The +1 counters can be unstable. Most of the time this passes, but sometimes
	the GIDs turn out significantly higher than +1. We may consider changing this
	to `gtest.LessPrim` later.
	*/
	gtest.Eq(goWait(fun), id0+1)
	gtest.Eq(goWait(fun), id0+2)
	gtest.Eq(goWait(fun), id0+3)
	gtest.Eq(goWait(fun), id0+4)
}

func goWait[A any](fun func() A) A {
	conn := make(chan A)
	go sendFrom(conn, fun)
	return <-conn
}

func sendFrom[A any](conn chan A, fun func() A) { conn <- fun() }

func BenchmarkGid_empty(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Gid())
	}
}

func BenchmarkGid_shallow(b *testing.B) {
	defer gtest.Catch(b)

	gg.WithSid(func() {
		for ind := 0; ind < b.N; ind++ {
			gg.Nop1(gg.Gid())
		}
	})
}

// Note that real app code usually has deeper stacks.
func BenchmarkGid_deeper(b *testing.B) {
	defer gtest.Catch(b)

	gg.WithGivenSid(math.MaxUint64, func() {
		for ind := 0; ind < b.N; ind++ {
			gg.Nop1(gg.Gid())
		}
	})
}

func Benchmark_GidWithOverride_deeper(b *testing.B) {
	defer gtest.Catch(b)

	gg.WithGivenSid(math.MaxUint64, func() {
		for ind := 0; ind < b.N; ind++ {
			gg.Nop1(gg.GidWithOverride())
		}
	})
}

func BenchmarkGidSlow_empty(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.GidSlow())
	}
}

func BenchmarkGidSlow_shallow(b *testing.B) {
	defer gtest.Catch(b)

	gg.WithSid(func() {
		for ind := 0; ind < b.N; ind++ {
			gg.Nop1(gg.GidSlow())
		}
	})
}

// Note that real app code usually has deeper stacks.
func BenchmarkGidSlow_deeper(b *testing.B) {
	defer gtest.Catch(b)

	gg.WithGivenSid(math.MaxUint64, func() {
		for ind := 0; ind < b.N; ind++ {
			gg.Nop1(gg.GidSlow())
		}
	})
}

/*
func Benchmark_runtime_Stack(b *testing.B) {
	defer gtest.Catch(b)

	buf := make([]byte, 1024)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(runtime.Stack(buf, false))
	}
}

func Benchmark_runtime_Callers(b *testing.B) {
	defer gtest.Catch(b)

	buf := make([]uintptr, 1024)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(runtime.Callers(0, buf))
	}
}

func Benchmark_runtime_Callers_with_inline_array(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		runtimeCallersWithInlineArray()
	}
}

func runtimeCallersWithInlineArray() {
	var arr [64]uintptr
	buf := gg.NoEscUnsafe(arr[:]) // Saves an allocation.
	gg.Nop1(runtime.Callers(0, buf))
}

func Benchmark_runtime_Callers_with_pool(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		runtimeCallersWithPool()
	}
}

var pool sync.Pool

func runtimeCallersWithPool() {
	var buf []uintptr
	ptr, _ := pool.Get().(*[]uintptr)

	if ptr == nil {
		buf = make([]uintptr, 64)
		ptr = &buf
	} else {
		buf = *ptr
	}

	defer pool.Put(ptr)
	gg.Nop1(runtime.Callers(0, buf))
}

func Benchmark_runtime_Callers_with_make(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		runtimeCallersWithMake()
	}
}

func runtimeCallersWithMake() {
	gg.Nop1(runtime.Callers(0, make([]uintptr, 64)))
}

func Benchmark_runtime_Callers_frames_iterate(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		runtimeCallersFramesIterate()
	}
}

func runtimeCallersFramesIterate() {
	buf := make([]uintptr, 64)
	buf = buf[:runtime.Callers(0, buf)]
	iter := runtime.CallersFrames(buf)
	for {
		frame, more := iter.Next()
		gg.Nop1(frame.PC)
		if !more {
			break
		}
	}
}
*/
