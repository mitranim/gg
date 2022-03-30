package gsql

import (
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
