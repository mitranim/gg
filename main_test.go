package gg_test

import (
	"strconv"

	"github.com/mitranim/gg"
)

func init() { gg.TraceBaseDir = gg.Cwd() }

var void struct{}

// Adapted from `reflect.ValueOf`.
func esc(val any) any {
	if trap.false {
		trap.val = val
	}
	return val
}

var trap struct {
	false bool
	val   any
}

type IntSet = gg.Set[int]

type IntMap = map[int]int

type SomeKey int64

type SomeModel struct {
	Id   SomeKey `json:"id"`
	Name string  `json:"name"`
}

func (self SomeModel) Pk() SomeKey { return self.Id }

type StructDirect struct {
	Public0 int
	Public1 string
	private *string
}

//nolint:unused
type StructIndirect struct {
	Public0 int
	Public1 *string
	private *string
}

type Outer struct {
	OuterId   int
	OuterName string
	Embed
	Inner *Inner
}

type Embed struct {
	EmbedId   int
	EmbedName string
}

type Inner struct {
	InnerId   *int
	InnerName *string
}

type SomeJsonDbMapper struct {
	SomeName  string `json:"someName" db:"some_name"`
	SomeValue string `json:"someValue" db:"some_value"`
	SomeJson  string `json:"someJson"`
	SomeDb    string `db:"some_db"`
}

type SomeColl = gg.Coll[SomeKey, SomeModel]

type SomeLazyColl = gg.LazyColl[SomeKey, SomeModel]

type IsZeroAlwaysTrue string

func (IsZeroAlwaysTrue) IsZero() bool { return true }

type IsZeroAlwaysFalse string

func (IsZeroAlwaysFalse) IsZero() bool { return false }

type FatStruct struct {
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, Id   int
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, Name string
}

type FatStructNonComparable struct {
	FatStruct
	_ []byte
}

func ComparerOf[A gg.LesserPrim](val A) Comparer[A] { return Comparer[A]{val} }

type Comparer[A gg.LesserPrim] [1]A

func (self Comparer[A]) Less(val Comparer[A]) bool { return self[0] < val[0] }

func (self Comparer[A]) Get() A { return self[0] }

func ToPair[A gg.Num](val A) (A, A) { return val - 1, val + 1 }

func True1[A any](A) bool { return true }

func False1[A any](A) bool { return false }

func Id1True[A any](val A) (A, bool) { return val, true }

func Id1False[A any](val A) (A, bool) { return val, false }

type ParserStr string

func (self *ParserStr) Parse(val string) error {
	*self = ParserStr(val)
	return nil
}

type UnmarshalerBytes []byte

func (self *UnmarshalerBytes) UnmarshalText(val []byte) error {
	*self = val
	return nil
}

// Implements `error` on the pointer type, not on the value type.
type PtrErrStr string

func (self *PtrErrStr) Error() string { return gg.PtrGet((*string)(self)) }

type StrsParser []string

func (self *StrsParser) Parse(src string) error {
	gg.Append(self, src)
	return nil
}

type IntsValue []int

func (self *IntsValue) Set(src string) error {
	return gg.ParseCatch(src, gg.AppendPtrZero(self))
}

/*
Defined as a struct to verify that the flag parser supports slices of arbitrary
types implementing `flag.Value`, even if they're not typedefs of text types,
and not something normally compatible with `gg.Parse`.
*/
type IntValue struct{ Val int }

func (self *IntValue) Set(src string) error {
	return gg.ParseCatch(src, &self.Val)
}

func intStrPair(src int) []string {
	return []string{strconv.Itoa(src - 1), strconv.Itoa(src + 1)}
}

func intPair(src int) []int {
	return []int{src - 1, src + 1}
}

type Cyclic struct {
	Id     int
	Cyclic *Cyclic
}
