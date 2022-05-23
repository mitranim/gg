package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

var (
	testErr0 = error(gg.Errf(`test err 0`))
	testErr1 = error(gg.Errf(`test err 1`))
	testErr2 = error(gg.Errf(`test err 2`))
)

const (
	testErrA = gg.ErrStr(`test err A`)
	testErrB = gg.ErrStr(`test err B`)
)

func TestConc(t *testing.T) {
	defer gtest.Catch(t)

	t.Run(`no_panic`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Zero(gg.ConcCatch())

		gtest.Equal(gg.ConcCatch(nil, nil, nil), []error{nil, nil, nil})

		gtest.Equal(
			gg.ConcCatch(func() {}),
			[]error{nil},
		)

		gtest.Equal(
			gg.ConcCatch(func() {}, func() {}),
			[]error{nil, nil},
		)

		gtest.Equal(
			gg.ConcCatch(func() {}, nil, func() {}),
			[]error{nil, nil, nil},
		)

		gtest.Equal(
			gg.ConcCatch(nil, func() {}, nil, func() {}, nil),
			[]error{nil, nil, nil, nil, nil},
		)
	})

	t.Run(`only_panic`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(
			gg.ConcCatch(func() { panic(testErr0) }),
			[]error{testErr0},
		)

		gtest.Equal(
			gg.ConcCatch(
				func() { panic(testErr0) },
				func() { panic(testErr1) },
			),
			[]error{testErr0, testErr1},
		)
	})

	t.Run(`mixed`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Equal(
			gg.ConcCatch(
				func() {},
				func() { panic(testErr0) },
				func() {},
				func() { panic(testErr1) },
				func() {},
			),
			gg.Errs{nil, testErr0, nil, testErr1, nil},
		)
	})
}

func BenchmarkConcCatch_one(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		_ = gg.ConcCatch(func() { panic(testErr0) })
	}
}

func BenchmarkConcCatch_multi(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		_ = gg.ConcCatch(
			func() {},
			func() { panic(testErr0) },
			func() {},
			func() { panic(testErr1) },
			func() {},
		)
	}
}
