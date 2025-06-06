package gsql_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gsql"
	"github.com/mitranim/gg/gtest"
)

func TestArrOf(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gsql.Arr[int](nil))
	gtest.Equal(gsql.ArrOf[int](), gsql.Arr[int](nil))
	gtest.Equal(gsql.ArrOf(10), gsql.Arr[int]{10})
	gtest.Equal(gsql.ArrOf(10, 20), gsql.Arr[int]{10, 20})
	gtest.Equal(gsql.ArrOf(10, 20, 30), gsql.Arr[int]{10, 20, 30})
}

func TestArr(t *testing.T) {
	t.Run(`String`, func(t *testing.T) {
		defer gtest.Catch(t)

		gtest.Str(gsql.Arr[int](nil), ``)
		gtest.Str(gsql.Arr[int]{}, `{}`)
		gtest.Str(gsql.Arr[int]{10}, `{10}`)
		gtest.Str(gsql.Arr[int]{10, 20}, `{10,20}`)
		gtest.Str(gsql.Arr[int]{10, 20, 30}, `{10,20,30}`)
		gtest.Str(gsql.Arr[gsql.Arr[int]]{{}, {}}, `{{},{}}`)
		gtest.Str(gsql.Arr[gsql.Arr[int]]{{10, 20}, {30, 40}}, `{{10,20},{30,40}}`)
	})

	t.Run(`Parse`, func(t *testing.T) {
		defer gtest.Catch(t)

		testParser(``, gsql.Arr[int](nil))
		testParser(`{}`, gsql.Arr[int]{})
		testParser(`{10}`, gsql.Arr[int]{10})
		testParser(`{10,20}`, gsql.Arr[int]{10, 20})
		testParser(`{10,20,30}`, gsql.Arr[int]{10, 20, 30})
		testParser(`{{},{}}`, gsql.Arr[gsql.Arr[int]]{{}, {}})
		testParser(`{{10},{20},{30,40}}`, gsql.Arr[gsql.Arr[int]]{{10}, {20}, {30, 40}})
		testParser(`{"10","20","30"}`, gsql.Arr[string]{`10`, `20`, `30`})
		testParser(`{("10","20"),("30","40")}`, gsql.Arr[string]{`("10","20")`, `("30","40")`})
		testParser(`{"(10,20)","(30,40)"}`, gsql.Arr[string]{`(10,20)`, `(30,40)`})
		testParser(`{"(\"10\",20)","(\"30\",40)"}`, gsql.Arr[string]{`("10",20)`, `("30",40)`})
		testParser(`{"(\"[\\\"one\\\"]\")","{\"(\\\"two\\\")\"}"}`, gsql.Arr[string]{`("[\"one\"]")`, `{"(\"two\")"}`})
	})
}

// TODO consider moving to `gtest`.
func testParser[
	A any,
	B interface {
		*A
		gg.Parser
	},
](src string, exp A) {
	var tar A
	gtest.NoErr(B(&tar).Parse(src))
	gtest.Equal(tar, exp)
}

func BenchmarkArr_Append(b *testing.B) {
	buf := make([]byte, 0, 4096)
	arr := gsql.ArrOf(10, 20, 30, 40, 50, 60, 70, 80)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(arr.AppendTo(buf))
	}
}

func BenchmarkArr_String(b *testing.B) {
	arr := gsql.ArrOf(10, 20, 30, 40, 50, 60, 70, 80)
	b.ResetTimer()

	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(arr.String())
	}
}
