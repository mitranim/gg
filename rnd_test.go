package gg_test

import (
	cr "crypto/rand"
	"fmt"
	"io"
	"math/big"
	mr "math/rand"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func rndSrc() io.Reader { return mr.New(mr.NewSource(0)) }

func TestRandomInt_uint64(t *testing.T) {
	defer gtest.Catch(t)

	fun := gg.RandomInt[uint64]
	src := rndSrc()

	gtest.Eq(fun(src), 13906042503472976897)
	gtest.Eq(fun(src), 14443988502964065089)
	gtest.Eq(fun(src), 17196259562495758190)
	gtest.Eq(fun(src), 8433884527138253544)
	gtest.Eq(fun(src), 16185558041432496379)
	gtest.Eq(fun(src), 15280578644808371633)
	gtest.Eq(fun(src), 16279533959527364769)
	gtest.Eq(fun(src), 5680856659213489233)
	gtest.Eq(fun(src), 18012265314398557154)
	gtest.Eq(fun(src), 12810001989293853876)
	gtest.Eq(fun(src), 8828672906944723960)
	gtest.Eq(fun(src), 11259781176380201387)
	gtest.Eq(fun(src), 6266933393232556850)
	gtest.Eq(fun(src), 8632501143404108278)
	gtest.Eq(fun(src), 6856693871787269831)
	gtest.Eq(fun(src), 6107792380522581863)
}

func TestRandomInt_int64(t *testing.T) {
	defer gtest.Catch(t)

	fun := gg.RandomInt[int64]
	src := rndSrc()

	gtest.Eq(fun(src), -4540701570236574719)
	gtest.Eq(fun(src), -4002755570745486527)
	gtest.Eq(fun(src), -1250484511213793426)
	gtest.Eq(fun(src), 8433884527138253544)
	gtest.Eq(fun(src), -2261186032277055237)
	gtest.Eq(fun(src), -3166165428901179983)
	gtest.Eq(fun(src), -2167210114182186847)
	gtest.Eq(fun(src), 5680856659213489233)
	gtest.Eq(fun(src), -434478759310994462)
	gtest.Eq(fun(src), -5636742084415697740)
	gtest.Eq(fun(src), 8828672906944723960)
	gtest.Eq(fun(src), -7186962897329350229)
	gtest.Eq(fun(src), 6266933393232556850)
	gtest.Eq(fun(src), 8632501143404108278)
	gtest.Eq(fun(src), 6856693871787269831)
	gtest.Eq(fun(src), 6107792380522581863)
}

func BenchmarkRandomInt_true_random(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.RandomInt[uint64](cr.Reader)
	}
}

func BenchmarkRandomInt_pseudo_random(b *testing.B) {
	defer gtest.Catch(b)
	src := rndSrc()
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.RandomInt[uint64](src)
	}
}

func TestRandomIntBetween_uint64(t *testing.T) {
	defer gtest.Catch(t)
	testRandomUintBetween[uint64]()
}

func TestRandomIntBetween_int64(t *testing.T) {
	defer gtest.Catch(t)
	testRandomUintBetween[int64]()
	testRandomSintBetween[int64]()
}

func testRandomUintBetween[A gg.Int]() {
	fun := gg.RandomIntBetween[A]

	gtest.PanicStr(`invalid range [0,0)`, func() {
		fun(rndSrc(), 0, 0)
	})

	gtest.PanicStr(`invalid range [1,1)`, func() {
		fun(rndSrc(), 1, 1)
	})

	gtest.PanicStr(`invalid range [2,1)`, func() {
		fun(rndSrc(), 2, 1)
	})

	testRandomIntRange[A](0, 1)
	testRandomIntRange[A](0, 2)
	testRandomIntRange[A](0, 3)
	testRandomIntRange[A](0, 16)
	testRandomIntRange[A](32, 48)

	testRandomIntRanges[A](0, 16)
	testRandomIntRanges[A](32, 48)
}

func testRandomSintBetween[A gg.Sint]() {
	gtest.PanicStr(`invalid range [1,-1)`, func() {
		gg.RandomIntBetween[A](rndSrc(), 1, -1)
	})

	gtest.PanicStr(`invalid range [2,-2)`, func() {
		gg.RandomIntBetween[A](rndSrc(), 2, -2)
	})

	testRandomIntRange[A](-10, 10)
	testRandomIntRanges[A](-10, 10)
}

func testRandomIntRange[A gg.Int](min, max A) {
	gtest.EqualSet(
		randomIntSlice(128, min, max),
		gg.Range(min, max),
		fmt.Sprintf(`expected range: [%v,%v)`, min, max),
	)
}

func testRandomIntRanges[A gg.Int](min, max int) {
	for _, max := range gg.Range(min+1, max) {
		for _, min := range gg.Range(min, max) {
			testRandomIntRange(gg.NumConv[A](min), gg.NumConv[A](max))
		}
	}
}

func randomIntSlice[A gg.Int](count int, min, max A) []A {
	src := rndSrc()
	var set gg.OrdSet[A]
	for range gg.Iter(count) {
		set.Add(gg.RandomIntBetween(src, min, max))
	}
	return set.Slice
}

func Benchmark_random_int_between_stdlib_true_random(b *testing.B) {
	defer gtest.Catch(b)
	max := new(big.Int).SetUint64(2 << 16)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Try1(cr.Int(cr.Reader, max))
	}
}

func Benchmark_random_int_between_stdlib_pseudo_random(b *testing.B) {
	defer gtest.Catch(b)
	src := rndSrc()
	max := new(big.Int).SetUint64(2 << 16)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Try1(cr.Int(src, max))
	}
}

func Benchmark_random_int_between_ours_true_random(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.RandomIntBetween(cr.Reader, 0, 2<<16)
	}
}

func Benchmark_random_int_between_ours_pseudo_random(b *testing.B) {
	defer gtest.Catch(b)
	src := rndSrc()
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.RandomIntBetween(src, 0, 2<<16)
	}
}

func TestRandomElem(t *testing.T) {
	defer gtest.Catch(t)

	gtest.PanicStr(`invalid range [0,0)`, func() {
		gg.RandomElem(rndSrc(), []int{})
	})

	testRandomElem([]int{10})
	testRandomElem([]int{10, 20})
	testRandomElem([]int{10, 20, 30})
	testRandomElem([]int{10, 20, 30, 40})
	testRandomElem([]int{10, 20, 30, 40, 50})
	testRandomElem([]int{10, 20, 30, 40, 50, 60})
}

func testRandomElem[A comparable](slice []A) {
	src := rndSrc()

	var set gg.OrdSet[A]
	for range gg.Iter(64) {
		set.Add(gg.RandomElem(src, slice))
	}

	gtest.EqualSet(set.Slice, slice)
}
