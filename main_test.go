package gg_test

import "github.com/mitranim/gg"

func init() { gg.TraceRelPath = true }

var void struct{}

type IntSet = gg.Set[int]

type SomeKey string

type SomeModel struct {
	Id   SomeKey `json:"id"`
	Name string  `json:"name"`
}

func (self SomeModel) Pk() SomeKey { return self.Id }

type SomeColl = gg.Coll[SomeKey, SomeModel]

type IsZeroAlwaysTrue string

func (IsZeroAlwaysTrue) IsZero() bool { return true }

type IsZeroAlwaysFalse string

func (IsZeroAlwaysFalse) IsZero() bool { return false }

type FatStruct struct {
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, Id   int
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, Name string
}

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
