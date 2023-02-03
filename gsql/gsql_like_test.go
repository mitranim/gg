package gsql_test

import (
	"testing"

	"github.com/mitranim/gg/gsql"
	"github.com/mitranim/gg/gtest"
)

func TestLike(t *testing.T) {
	defer gtest.Catch(t)

	test := func(src, esc string) {
		tar := gsql.Like(src)

		gtest.Eq(tar.IsNull(), src == ``)
		gtest.Eq(tar.String(), src)
		gtest.Eq(tar.Esc(), esc)
	}

	test(``, ``)
	test(` `, `% %`)
	test(`str`, `%str%`)
	test(`%`, `%\%%`)
	test(`_`, `%\_%`)
	test(`%str%`, `%\%str\%%`)
	test(`_str_`, `%\_str\_%`)
}
