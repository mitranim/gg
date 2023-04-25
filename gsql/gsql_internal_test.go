package gsql

import (
	r "reflect"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

type Inner struct {
	InnerId   string  `db:"inner_id"`
	InnerName *string `db:"inner_name"`
}

type Outer struct {
	OuterId   int64         `db:"outer_id"`
	OuterName string        `db:"outer_name"`
	InnerZop  gg.Zop[Inner] `db:"inner_zop"`
	Inner     Inner         `db:"inner"`
}

func Test_structMetaCache(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(
		typeMetaCache.Get(gg.Type[Outer]()),
		typeMeta{
			`outer_id`:             []int{0},
			`outer_name`:           []int{1},
			`inner_zop.inner_id`:   []int{2, 0, 0},
			`inner_zop.inner_name`: []int{2, 0, 1},
			`inner.inner_id`:       []int{3, 0},
			`inner.inner_name`:     []int{3, 1},
		},
	)
}

/*
Used by `scanValsReflect`. This demonstrates doubling behavior of
`reflect.Value.Grow`. Our choice of implementation relies on this behavior. If
`reflect.Value.Grow` allocated precisely the requested amount of additional
capacity, which in our case is 1, we would have to change our strategy.

Compare our `gg.GrowCap`, which behaves like `reflect.Value.Grow`, and
`gg.GrowCapExact`, which behaves in a way that would be detrimental for
the kind of algorithm we use here.
*/
func Test_reflect_slice_grow_alloc(t *testing.T) {
	defer gtest.Catch(t)

	var tar []int

	val := r.ValueOf(&tar).Elem()

	test := func(diff, total int) {
		prevLen := len(tar)

		val.Grow(diff)
		gtest.Eq(val.Len(), len(tar))
		gtest.Eq(val.Cap(), cap(tar))
		gtest.Eq(cap(tar), total)
		gtest.Eq(len(tar), prevLen)

		val.SetLen(total)
		gtest.Eq(val.Len(), len(tar))
		gtest.Eq(len(tar), total)
	}

	test(0, 0)
	test(1, 1)
	test(1, 2)
	test(1, 4)
	test(1, 8)
	test(1, 16)
	test(1, 32)
}
