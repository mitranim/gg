package grepr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/grepr"
	"github.com/mitranim/gg/gtest"
)

type Struct0 struct{}

type Struct1 struct{ A int }

type Struct2 struct {
	A int
	B int
}

type Inner struct {
	InnerId   *int
	InnerName *string
}

type Embed struct {
	EmbedId   int
	EmbedName string
}

type Outer struct {
	OuterId   int
	OuterName string
	Embed
	Inner *Inner
}

type Cyclic struct {
	Id     int
	Cyclic *Cyclic
}

var testInner = Inner{
	InnerId:   gg.Ptr(30),
	InnerName: gg.Ptr(`inner`),
}

var testEmbed = Embed{EmbedId: 20}

var testOuter = Outer{
	OuterName: `outer`,
	Embed:     testEmbed,
	Inner:     &testInner,
}

func init() { gg.TraceRelPath = true }

func testRepr[A any](src A, exp string) {
	gtest.Eq(grepr.String(src), exp)
}

func TestString(t *testing.T) {
	defer gtest.Catch(t)

	testRepr(any(nil), `nil`)

	testRepr(false, `false`)
	testRepr(true, `true`)

	testRepr(-10, `-10`)
	testRepr(0, `0`)
	testRepr(10, `10`)

	testRepr(``, "``")
	testRepr(`str`, "`str`")

	testRepr([]byte(nil), `nil`)
	testRepr([]byte{}, "[]byte(``)")
	testRepr([]byte(`str`), "[]byte(`str`)")

	testRepr(gg.Buf(nil), `nil`)
	testRepr(gg.Buf{}, "gg.Buf(``)")
	testRepr(gg.Buf(`str`), "gg.Buf(`str`)")

	testRepr(Struct0{}, `grepr_test.Struct0{}`)

	testRepr(Struct1{}, `grepr_test.Struct1{}`)
	testRepr(Struct1{10}, `grepr_test.Struct1{10}`)

	testRepr(Struct2{}, `grepr_test.Struct2{}`)
	testRepr(Struct2{A: 10}, `grepr_test.Struct2{A: 10}`)
	testRepr(Struct2{B: 20}, `grepr_test.Struct2{B: 20}`)
	testRepr(Struct2{A: 10, B: 20}, `grepr_test.Struct2{
    A: 10,
    B: 20,
}`)

	testRepr(
		testOuter,
		strings.TrimSpace(`
grepr_test.Outer{
    OuterName: `+"`outer`"+`,
    Embed: grepr_test.Embed{EmbedId: 20},
    Inner: &grepr_test.Inner{
        InnerId: gg.Ptr(30),
        InnerName: gg.Ptr(`+"`inner`"+`),
    },
}
`),
	)
}

func ExampleString() {
	fmt.Println(grepr.String(testOuter))
	// Output:
	// grepr_test.Outer{
	//     OuterName: `outer`,
	//     Embed: grepr_test.Embed{EmbedId: 20},
	//     Inner: &grepr_test.Inner{
	//         InnerId: gg.Ptr(30),
	//         InnerName: gg.Ptr(`inner`),
	//     },
	// }
}

func BenchmarkString_num(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(grepr.String(10))
	}
}

func BenchmarkString_str(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(grepr.String(`str`))
	}
}

func BenchmarkString_struct_flat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(grepr.String(testEmbed))
	}
}

func BenchmarkString_struct_nested(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(grepr.String(testOuter))
	}
}

func BenchmarkString_fmt_GoStringer(b *testing.B) {
	src := gg.Set[int]{}
	gg.Nop1(fmt.GoStringer(src))
	gtest.Eq(grepr.String(src), `gg.Set[int]{}`)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gg.Nop1(grepr.String(src))
	}
}

func Test_cyclic_1(t *testing.T) {
	defer gtest.Catch(t)

	cyclic := Cyclic{Id: 10}
	cyclic.Cyclic = &cyclic

	gtest.Eq(&cyclic, cyclic.Cyclic)
	testCyclic(&cyclic)
}

func Test_cyclic_2(t *testing.T) {
	defer gtest.Catch(t)

	cyclic0 := Cyclic{Id: 10}
	cyclic1 := Cyclic{Id: 20}

	cyclic0.Cyclic = &cyclic1
	cyclic1.Cyclic = &cyclic0

	gtest.Eq(&cyclic0, cyclic1.Cyclic)
	gtest.Eq(&cyclic1, cyclic0.Cyclic)

	testCyclic(&cyclic0)
}

/*
For now, this verifies the following:

	* We eventually terminate.
	* We mark visited references.

TODO verify the exact output structure. It can be broken by unsafe hacks such as
`gg.AnyNoEscUnsafe`.
*/
func testCyclic[A any](src A) {
	gtest.TextHas(grepr.String(src), `/* visited */ (*`)
}

func BenchmarkString_cyclic(b *testing.B) {
	cyclic0 := Cyclic{Id: 10}
	cyclic1 := Cyclic{Id: 20}

	cyclic0.Cyclic = &cyclic1
	cyclic1.Cyclic = &cyclic0

	for i := 0; i < b.N; i++ {
		gg.Nop1(grepr.String(&cyclic0))
	}
}

func BenchmarkCanBackquote(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.Nop1(grepr.CanBackquote(`
one
two
three
`))
	}
}
