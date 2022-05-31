package gg_test

import "github.com/mitranim/gg"

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

type SomeKey string

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

type IsZeroAlwaysTrue string

func (IsZeroAlwaysTrue) IsZero() bool { return true }

type IsZeroAlwaysFalse string

func (IsZeroAlwaysFalse) IsZero() bool { return false }

type FatStruct struct {
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, Id   int
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, Name string
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

func (self *PtrErrStr) Error() string { return gg.Deref((*string)(self)) }

type StrsParser []string

func (self *StrsParser) Parse(src string) error {
	gg.AppendVals(self, src)
	return nil
}
