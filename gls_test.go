package gg_test

import (
	"context"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

/*
Also see `dyn_var_test.go` which tests other aspects of GLS.
*/

func Test_GlsGo_GlsRun(t *testing.T) {
	defer gtest.Catch(t)
	t.Cleanup(gg.GlsClear)

	gtest.Zero(DYN_NUM.Get())
	gtest.Zero(DYN_MOD.Get())

	gid0 := gg.Gid()
	num0 := 10
	mod0 := SomeModel{Id: 20, Name: `one`}
	gls0 := map[uint64]gg.GlsInternal{
		gid0: {
			gg.GlsKey(DYN_NUM): num0,
			gg.GlsKey(DYN_MOD): mod0,
		},
	}

	defer DYN_NUM.Set(num0).Use()
	defer DYN_MOD.Set(mod0).Use()

	gtest.Eq(DYN_NUM.Get(), num0)
	gtest.Eq(DYN_MOD.Get(), mod0)
	gtest.Equal(gg.Glss(), gls0)

	// No inheritance in regular `go` (because we don't know how to).
	goWait(func() (_ struct{}) {
		gtest.Zero(DYN_NUM.Get())
		gtest.Zero(DYN_MOD.Get())
		gtest.Equal(gg.Glss(), gls0)
		return
	})

	var called Called
	conn := make(chan struct{})

	// Without overrides.
	gg.GlsGo(func() {
		defer gg.SendZero(conn)

		gtest.Eq(DYN_NUM.Get(), num0)
		gtest.Eq(DYN_MOD.Get(), mod0)

		gtest.Equal(gg.Glss(), map[uint64]gg.GlsInternal{
			gid0: {
				gg.GlsKey(DYN_NUM): num0,
				gg.GlsKey(DYN_MOD): mod0,
			},
			// Sub-goroutine gets a copy of the parent GLS.
			gg.Gid(): {
				gg.GlsKey(DYN_NUM): num0,
				gg.GlsKey(DYN_MOD): mod0,
			},
		})

		called.Here()
	})

	<-conn
	called.Verify()

	gtest.Eq(DYN_NUM.Get(), num0)
	gtest.Eq(DYN_MOD.Get(), mod0)
	gtest.Equal(gg.Glss(), gls0)

	num1 := 30
	mod1 := SomeModel{Id: 40, Name: `two`}

	subRun := func() {
		gtest.Eq(DYN_NUM.Get(), num1)
		gtest.Eq(DYN_MOD.Get(), mod1)

		gtest.Equal(gg.Glss(), map[uint64]gg.GlsInternal{
			gid0: {
				gg.GlsKey(DYN_NUM): num0,
				gg.GlsKey(DYN_MOD): mod0,
			},
			gg.Gid(): {
				gg.GlsKey(DYN_NUM): num1,
				gg.GlsKey(DYN_MOD): mod1,
			},
		})

		called.Here()
	}

	subGo := func() {
		defer gg.SendZero(conn)
		subRun()
	}

	gg.GlsGo(subGo, DYN_NUM.With(num1), DYN_MOD.With(mod1))
	<-conn
	called.Verify()

	gtest.Eq(DYN_NUM.Get(), num0)
	gtest.Eq(DYN_MOD.Get(), mod0)
	gtest.Equal(gg.Glss(), gls0)

	gg.GlsRun(subRun, DYN_NUM.With(num1), DYN_MOD.With(mod1))
	called.Verify()
}

func TestGlsSnap(t *testing.T) {
	defer gtest.Catch(t)

	// Ensure this test's GLS doesn't pollute the GLSS.
	t.Cleanup(gg.GlsClear)

	gtest.Zero(gg.Glss())
	gtest.Empty(gg.GlsSnap())

	num := 10
	mod := SomeModel{Id: 20, Name: `one`}

	DYN_NUM.Set(num)
	DYN_MOD.Set(mod)

	gtest.Equal(gg.Glss(), map[uint64]gg.GlsInternal{
		gg.Gid(): {
			gg.GlsKey(DYN_NUM): num,
			gg.GlsKey(DYN_MOD): mod,
		},
	})

	gtest.EqualSet(gg.GlsSnap(), []gg.GlsVal{
		DYN_NUM.With(num),
		DYN_MOD.With(mod),
	})
}

func TestGlsSet_Gls_Use(t *testing.T) {
	defer gtest.Catch(t)

	// Ensure this test's GLS doesn't pollute the GLSS.
	t.Cleanup(gg.GlsClear)

	dyn0 := gg.NewDynVar[int](nil)
	dyn1 := gg.NewDynVar[int](nil)
	dyn2 := gg.NewDynVar[int](nil)

	gtest.Zero(gg.Glss())
	gls0 := gg.GlsSet()
	gtest.Zero(gls0)
	gtest.Zero(gg.Glss())

	dyn0.Set(10)
	dyn1.Set(20)
	dyn2.Set(30)

	orig := map[uint64]gg.GlsInternal{
		gg.Gid(): {
			gg.GlsKey(dyn0): 10,
			gg.GlsKey(dyn1): 20,
			gg.GlsKey(dyn2): 30,
		},
	}

	gtest.Equal(gg.Glss(), orig)

	gtest.Eq(dyn0.Get(), 10)
	gtest.Eq(dyn1.Get(), 20)
	gtest.Eq(dyn2.Get(), 30)

	func() {
		defer gg.GlsSet(dyn0.With(40), dyn2.With(50)).Use()

		gtest.Equal(gg.Glss(), map[uint64]gg.GlsInternal{
			gg.Gid(): {
				gg.GlsKey(dyn0): 40,
				gg.GlsKey(dyn2): 50,
			},
		})

		gtest.Eq(dyn0.Get(), 40)
		gtest.Eq(dyn1.Get(), 0)
		gtest.Eq(dyn2.Get(), 50)
	}()

	gtest.Equal(gg.Glss(), orig)
	gls0.Use()
	gtest.Zero(gg.Glss())
}

var glsVals = []gg.GlsVal{
	DYN_CTX.With(context.Background()),
	DYN_NUM.With(123),
	DYN_PTR.With(new(int)),
	DYN_MOD.With(SomeModel{}),
}

func BenchmarkGlsCopy(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)
	gg.GlsSet(glsVals...)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.GlsCopy())
	}
}

func BenchmarkGlsSnap(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)
	gg.GlsSet(glsVals...)
	gtest.Len(gg.GlsSnap(), 4)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.GlsSnap())
	}
}

/*
For the GLSS, this benchmark and `BenchmarkGlsSet` slightly favor map+mutex
over `sync.Map`, but `BenchmarkDynVar_with_minor_concurrency` heavily favors
`sync.Map` over map+mutex.
*/
func BenchmarkGlsSet(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.GlsSet(glsVals...))
	}
}

func BenchmarkGls_Use(b *testing.B) {
	defer gtest.Catch(b)
	b.Cleanup(gg.GlsClear)

	gg.GlsSet(
		DYN_CTX.With(context.Background()),
		DYN_NUM.With(123),
		DYN_PTR.With(new(int)),
		DYN_MOD.With(SomeModel{}),
	)
	gls := gg.GlsCopy()

	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gls.Use())
	}
}
