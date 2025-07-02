package gg_test

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

type Called bool

// TODO better naming.
func (self *Called) Here() {
	gtest.False(bool(*self))
	*self = true
}

func (self *Called) Verify() {
	gtest.True(bool(*self))
	*self = false
}

// It's really tested in the "with sid" tests.
func TestSid(t *testing.T) {
	defer gtest.Catch(t)
	gtest.Zero(gg.Sid())
	gtest.Zero(gg.Sid())
	gtest.Zero(gg.Sid())
}

func TestWithSid(t *testing.T) {
	defer gtest.Catch(t)

	var called Called

	gg.WithSid(nil)

	gtest.Zero(gg.Sid())

	gg.WithSid(func() {
		sid1 := gg.Sid()
		gtest.Eq(sid1, 1)
		gtest.Eq(gg.Sid(), sid1)
		gtest.Eq(gg.Sid(), sid1)

		gg.WithSid(func() {
			sid2 := gg.Sid()
			gtest.Eq(sid2, 2)
			gtest.Eq(gg.Sid(), sid2)
			gtest.Eq(gg.Sid(), sid2)
			called.Here()
		})

		called.Verify()

		// Reuse freed.
		gg.WithSid(func() {
			sid2 := gg.Sid()
			gtest.Eq(sid2, 2)
			gtest.Eq(gg.Sid(), sid2)
			gtest.Eq(gg.Sid(), sid2)
			called.Here()
		})

		called.Verify()
		called.Here()
	})

	called.Verify()

	// Reuse freed.
	gg.WithSid(func() {
		sid1 := gg.Sid()
		gtest.Eq(sid1, 1)
		gtest.Eq(gg.Sid(), sid1)
		gtest.Eq(gg.Sid(), sid1)
		called.Here()
	})

	called.Verify()
}

func TestWithGivenSid(t *testing.T) {
	defer gtest.Catch(t)

	gg.WithGivenSid(0, nil)
	gg.WithGivenSid(1, nil)

	test := func(sid uint64) {
		var called Called

		gg.WithGivenSid(sid, func() {
			gtest.Eq(gg.Sid(), sid)

			// Enable on demand.
			// fmt.Printf("sid: %[1]v 0x%[1]x; trace:%[2]v\n\n", sid, gg.CaptureTrace(0))

			names := sidFuncNames(gg.CaptureTrace(1))

			var exp []string

			gg.Append(&exp, `gg.sid_start`)
			if sid != 0 {
				for _, char := range strconv.FormatUint(sid, 16) {
					gg.Append(&exp, `gg.sid_digit_0x`+string(char))
				}
			}
			gg.Append(&exp, `gg.sid_end`)

			gtest.Equal(names, exp, `sid =`, sid)
			called.Here()
		})

		called.Verify()
	}

	test(0x0)
	test(0x1)
	test(0x2)
	test(0x3)
	test(0x4)
	test(0x5)
	test(0x6)
	test(0x7)
	test(0x8)
	test(0x9)
	test(0xa)
	test(0xb)
	test(0xc)
	test(0xd)
	test(0xe)
	test(0xf)

	test(0x10)
	test(0x20)
	test(0x30)
	test(0x40)
	test(0x50)
	test(0x60)
	test(0x70)
	test(0x80)
	test(0x90)
	test(0xa0)
	test(0xb0)
	test(0xc0)
	test(0xd0)
	test(0xe0)
	test(0xf0)

	test(0x1010)
	test(0x2020)
	test(0x3030)
	test(0x4040)
	test(0x5050)
	test(0x6060)
	test(0x7070)
	test(0x8080)
	test(0x9090)
	test(0xa0a0)
	test(0xb0b0)
	test(0xc0c0)
	test(0xd0d0)
	test(0xe0e0)
	test(0xf0f0)

	// The 0 is in the middle because it never goes on the stack at the edges.
	test(0x1234567089abcdef)

	test(math.MaxUint64)
}

func sidFuncNames(src []uintptr) (out []string) {
	for _, addr := range gg.CaptureTrace(2) {
		name := gg.Caller(addr).Frame().Name
		if !strings.HasPrefix(name, `gg.sid_`) {
			break
		}
		out = append(out, name)
	}
	return
}

func BenchmarkSid_empty(b *testing.B) {
	defer gtest.Catch(b)

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(gg.Sid())
	}
}

func BenchmarkSid_shallow(b *testing.B) {
	defer gtest.Catch(b)

	gg.WithSid(func() {
		for ind := 0; ind < b.N; ind++ {
			gg.Nop1(gg.Sid())
		}
	})
}

func BenchmarkSid_deep(b *testing.B) {
	defer gtest.Catch(b)

	gg.WithGivenSid(math.MaxUint64, func() {
		for ind := 0; ind < b.N; ind++ {
			gg.Nop1(gg.Sid())
		}
	})
}

/*
func Benchmark_gls_shallow(b *testing.B) {
	defer gtest.Catch(b)

	gls.EnsureGoroutineId(func(uint) {
		for ind := 0; ind < b.N; ind++ {
			gg.Nop2(gls.GetGoroutineId())
		}
	})
}
*/

/*
func Test_github_com_jtolio_gls(t *testing.T) {
	defer gtest.Catch(t)

	var gro sync.WaitGroup

	for range gg.Iter(3) {
		gro.Add(1)

		go func() {
			defer gro.Done()
			gls.EnsureGoroutineId(func(id uint) {
				fmt.Printf("id: %v%v\n\n", id, gg.CaptureTrace(0))
			})
		}()
	}

	gro.Wait()
}
*/
