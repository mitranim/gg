package gg_test

import (
	"context"
	"sync"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

var DYN_CTX = gg.NewDynVar[context.Context](context.Background)
var DYN_NUM = gg.NewDynVar[int](nil)
var DYN_PTR = gg.NewDynVar[*int](func() *int { return new(int) })
var DYN_MOD = gg.NewDynVar[SomeModel](func() (_ SomeModel) { return })

func TestDynVar(t *testing.T) {
	defer gtest.Catch(t)
	t.Cleanup(gg.GlsClear)

	gid0 := gg.Gid()
	gtest.NotZero(gid0)

	var gid1 uint64
	var dyn gg.DynVar[string]

	glss0 := map[uint64]gg.GlsInternal{
		gid0: {
			gg.GlsKey(&dyn): `one`,
		},
	}

	var glss01 map[uint64]gg.GlsInternal

	gtest.Zero(dyn.Get())
	gtest.Zero(gg.Glss())

	defer dyn.Set(`one`).Use()
	gtest.Equal(gg.Glss(), glss0)
	gtest.Eq(dyn.Get(), `one`)

	step0 := make(chan struct{}, 1)
	step1 := make(chan struct{}, 1)
	step2 := make(chan struct{}, 1)
	step3 := make(chan struct{}, 1)

	go func() {
		defer gtest.Catch(t)

		gid1 = gg.Gid()
		gtest.Eq(gid1, gid0+1)
		gtest.Equal(gg.Glss(), glss0)
		gtest.Zero(dyn.Get(), `must not accidentally inherit`)

		<-step0

		defer dyn.Set(`two`).Use()

		glss01 = map[uint64]gg.GlsInternal{
			gid0: {
				gg.GlsKey(&dyn): `one`,
			},
			gid1: {
				gg.GlsKey(&dyn): `two`,
			},
		}

		gtest.Equal(gg.Glss(), glss01)
		gtest.Eq(dyn.Get(), `two`)

		gg.SendZero(step1)
		<-step2

		dyn.Clear()
		gtest.Equal(gg.Glss(), glss0)
		gtest.Zero(dyn.Get())
		gg.SendZero(step3)
	}()

	gtest.Equal(gg.Glss(), glss0)
	gtest.Eq(dyn.Get(), `one`)

	gg.SendZero(step0)
	<-step1

	gtest.Eq(gid1, gid0+1)
	gtest.MapNotEmpty(glss01)

	gtest.Equal(gg.Glss(), glss01, `must be modified by other goroutine`)
	gtest.Eq(dyn.Get(), `one`, `must be unaffected by other goroutine`)

	gg.SendZero(step2)
	<-step3

	gtest.Equal(gg.Glss(), glss0)
	gtest.Eq(dyn.Get(), `one`)

	dyn.Clear()
	gtest.Zero(gg.Glss())
	gtest.Zero(dyn.Get())
}

func TestDynVar_default(t *testing.T) {
	defer gtest.Catch(t)
	t.Cleanup(gg.GlsClear)

	def := func() int { return 123 }
	dyn := gg.NewDynVar(def)

	testDynVarDef(dyn, def, 0)
	gtest.Eq(dyn.Get(), 123)
	testDynVarDef(dyn, nil, 123)
}

func TestDynVar_default_panic_retry(t *testing.T) {
	defer gtest.Catch(t)
	t.Cleanup(gg.GlsClear)

	var count int

	def := func() int {
		count++
		if count <= 3 {
			panic(`intermittent_failure`)
		}
		return 123
	}

	dyn := gg.NewDynVar(def)
	gtest.Is(gg.DynVarDef(dyn), def)

	gtest.PanicStr(`intermittent_failure`, func() { dyn.Get() })
	gtest.PanicStr(`intermittent_failure`, func() { dyn.Get() })
	gtest.PanicStr(`intermittent_failure`, func() { dyn.Get() })

	testDynVarDef(dyn, def, 0)
	gtest.Eq(dyn.Get(), 123)
	testDynVarDef(dyn, nil, 123)
}

func testDynVarDef[A any](dyn *gg.DynVar[A], def func() A, val A) {
	gtest.Is(gg.DynVarDef(dyn), def)
	gtest.Equal(gg.DynVarVal(dyn), val)
}

func TestDynVar_GetOr(t *testing.T) {
	defer gtest.Catch(t)
	t.Cleanup(gg.GlsClear)

	var dyn gg.DynVar[int]
	gtest.Zero(dyn.GetOr(nil))
	gtest.MapEmpty(gg.Glss())

	test := func(val int) { testGetOr(&dyn, val) }

	test(10)
	test(10)
	test(10)

	dyn.Clear()

	test(20)
	test(20)
	test(20)
}

func TestDynVar_GetOr_with_default(t *testing.T) {
	defer gtest.Catch(t)
	t.Cleanup(gg.GlsClear)

	def := func() int { return 123 }
	dyn := gg.NewDynVar(def)
	testDynVarDef(dyn, def, 0)

	test := func(val int) { testGetOr(dyn, val) }

	test(10)
	test(10)
	test(10)
	testDynVarDef(dyn, def, 0)

	dyn.Clear()
	gtest.Eq(dyn.Get(), 123)

	testDynVarDef(dyn, nil, 123)

	// The function passed to `.GetOr` takes priority over the default.
	test(20)
	test(20)
	test(20)

	dyn.Clear()
	gtest.Eq(dyn.GetOr(nil), 123)
}

func testGetOr(dyn *gg.DynVar[int], val int) {
	gtest.Eq(dyn.GetOr(func() int { return val }), val)

	gtest.Equal(gg.Glss(), map[uint64]gg.GlsInternal{
		gg.Gid(): {gg.GlsKey(dyn): val},
	})
}

func BenchmarkDynVar_Get_num_empty(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)

	var dyn gg.DynVar[int]

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(dyn.Get())
	}
}

func BenchmarkDynVar_Get_num_empty_with_default(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)

	dyn := gg.NewDynVar(func() (_ int) { return })

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(dyn.Get())
	}
}

func BenchmarkDynVar_Set_num(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)

	const val = 123

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(DYN_NUM.Set(val))
	}
}

func BenchmarkDynVar_Get_num(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)
	DYN_NUM.Set(123)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(DYN_NUM.Get())
	}
}

func BenchmarkDynVar_Set_iface(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)

	val := context.Background()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(DYN_CTX.Set(val))
	}
}

func BenchmarkDynVar_Get_iface(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)
	DYN_CTX.Set(context.Background())

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(DYN_CTX.Get())
	}
}

func BenchmarkDynVar_Set_ptr(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)

	val := new(int)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(DYN_PTR.Set(val))
	}
}

func BenchmarkDynVar_Get_ptr(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)
	DYN_PTR.Set(new(int))

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(DYN_PTR.Get())
	}
}

/*
In Go 1.24.0, this is the only `DynVar.Set` benchmark that show an allocation,
presumably the copying of the value to the heap, since it doesn't fit into one
machine word in `any`.

We also allocate one GLS map per goroutine, which does not show up on this
benchmark. We could use a pool, but doesn't seem worth the bother for now.
*/
func BenchmarkDynVar_Set_mod(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)

	var val SomeModel

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(DYN_MOD.Set(val))
	}
}

func BenchmarkDynVar_Get_mod(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)
	DYN_MOD.Set(SomeModel{})

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(DYN_MOD.Get())
	}
}

/*
At the time of writing (in Go 1.24), this seems to heavily favor `sync.Map`
over mutex+map for the GLSS, even though `sync.Map` causes more allocations
and performs slightly worse in non-concurrency microbenchmarks.
*/
func BenchmarkDynVar_with_minor_concurrency(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)

	var gro sync.WaitGroup

	for ind := 0; ind < b.N; ind++ {
		for range gg.Iter(64) {
			gro.Add(1)

			go func() {
				defer gro.Done()
				defer gg.GlsClear()
				defer DYN_MOD.Set(SomeModel{}).Use()

				for range gg.Iter(32) {
					gg.Nop1(DYN_MOD.Get())
				}
			}()
		}

		gro.Wait()
	}
}

/*
func Benchmark_sync_Map_Load_ptr(b *testing.B) {
	defer gtest.Catch(b)

	var tar sync.Map
	tar.Store(gg.Gid(), new(int))

	for ind := 0; ind < b.N; ind++ {
		val, _ := tar.Load(gg.Gid())
		gg.Nop1(val.(*int))
	}
}

func Benchmark_sync_Map_Load_wide(b *testing.B) {
	defer gtest.Catch(b)

	var tar sync.Map
	tar.Store(gg.Gid(), SomeModel{})

	for ind := 0; ind < b.N; ind++ {
		val, _ := tar.Load(gg.Gid())
		gg.Nop1(val.(SomeModel))
	}
}
*/
