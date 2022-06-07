package gsql_test

import (
	"testing"

	"github.com/mitranim/gg"
	s "github.com/mitranim/gg/gsql"
	gtest "github.com/mitranim/gg/gtest"
)

func TestRune_IsNull(t *testing.T) {
	defer gtest.Catch(t)

	gtest.True(s.Rune(0).IsNull())
	gtest.False(s.Rune(1).IsNull())
}

func TestRune_IsNonNull(t *testing.T) {
	defer gtest.Catch(t)

	gtest.False(s.Rune(0).IsNonNull())
	gtest.True(s.Rune(1).IsNonNull())
}

func TestRune_Clear(t *testing.T) {
	defer gtest.Catch(t)

	var tar s.Rune
	gtest.Zero(tar)

	tar = '👍'
	gtest.NotZero(tar)

	tar.Clear()
	gtest.Zero(tar)

	tar.Clear()
	gtest.Zero(tar)
}

func TestRune_String(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Str(s.Rune(0), ``)
	gtest.Str(s.Rune('👍'), `👍`)
}

func BenchmarkRune_String(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop1(s.Rune('👍'))
	}
}

func TestRune_Append(t *testing.T) {
	defer gtest.Catch(t)

	buf := gg.ToBytes(`init_`)

	gtest.Equal(s.Rune(0).Append(buf), buf)
	gtest.Equal(s.Rune(0).Append(nil), nil)

	gtest.Equal(s.Rune('👍').Append(buf), gg.ToBytes(`init_👍`))
	gtest.Equal(s.Rune('👍').Append(nil), gg.ToBytes(`👍`))
}

func TestRune_Parse(t *testing.T) {
	defer gtest.Catch(t)
	testRuneParse((*s.Rune).Parse)
}

func testRuneParse(fun func(*s.Rune, string) error) {
	gtest.ErrorStr(`unable to parse "ab" as char: too many chars`, fun(new(s.Rune), `ab`))
	gtest.ErrorStr(`unable to parse "abc" as char: too many chars`, fun(new(s.Rune), `abc`))
	gtest.ErrorStr(`unable to parse "👍👎" as char: too many chars`, fun(new(s.Rune), `👍👎`))

	var tar s.Rune

	gtest.NoError(fun(&tar, `🙂`))
	gtest.Eq(tar, '🙂')

	gtest.NoError(fun(&tar, ``))
	gtest.Zero(tar)
}

func BenchmarkRune_Parse_empty(b *testing.B) {
	var tar s.Rune

	for ind := 0; ind < b.N; ind++ {
		gg.Try(tar.Parse(``))
	}
}

func BenchmarkRune_Parse_non_empty(b *testing.B) {
	var tar s.Rune

	for ind := 0; ind < b.N; ind++ {
		gg.Try(tar.Parse(`🙂`))
	}
}

func TestRune_MarshalText(t *testing.T) {
	defer gtest.Catch(t)

	encode := func(src s.Rune) string {
		return gg.ToString(gg.Try1(src.MarshalText()))
	}

	gtest.Eq(encode(0), ``)
	gtest.Eq(encode('👍'), `👍`)
}

func BenchmarkRune_MarshalText(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop2(s.Rune('👍').MarshalText())
	}
}

func TestRune_UnmarshalText(t *testing.T) {
	defer gtest.Catch(t)
	testRuneParse(charUnmarshalText)
}

func charUnmarshalText(tar *s.Rune, src string) error {
	return tar.UnmarshalText(gg.ToBytes(src))
}

func BenchmarkRune_UnmarshalText(b *testing.B) {
	var tar s.Rune

	for ind := 0; ind < b.N; ind++ {
		gg.Try(tar.UnmarshalText(gg.ToBytes(`👍`)))
	}
}

func TestRune_MarshalJSON(t *testing.T) {
	defer gtest.Catch(t)

	encode := func(src s.Rune) string {
		return gg.ToString(gg.Try1(src.MarshalJSON()))
	}

	gtest.Eq(encode(0), `null`)
	gtest.Eq(encode('👍'), `"👍"`)
}

func BenchmarkRune_MarshalJSON_empty(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop2(s.Rune(0).MarshalJSON())
	}
}

func BenchmarkRune_MarshalJSON_non_empty(b *testing.B) {
	for ind := 0; ind < b.N; ind++ {
		gg.Nop2(s.Rune('👍').MarshalJSON())
	}
}

func TestRune_UnmarshalJSON(t *testing.T) {
	defer gtest.Catch(t)

	testRuneParse(charUnmarshalJson)

	gtest.ErrorStr(
		`cannot unmarshal number into Go value of type string`,
		new(s.Rune).UnmarshalJSON(gg.ToBytes(`123`)),
	)

	{
		tar := s.Rune('👍')
		gtest.NoError(tar.UnmarshalJSON(nil))
		gtest.Zero(tar)
	}

	{
		tar := s.Rune('👍')
		gtest.NoError(tar.UnmarshalJSON(gg.ToBytes(`null`)))
		gtest.Zero(tar)
	}
}

func charUnmarshalJson(tar *s.Rune, src string) error {
	return tar.UnmarshalJSON(gg.JsonBytes(src))
}

func BenchmarkRune_UnmarshalJSON_empty(b *testing.B) {
	var tar s.Rune

	for ind := 0; ind < b.N; ind++ {
		gg.Try(tar.UnmarshalJSON(gg.ToBytes(`null`)))
	}
}

func BenchmarkRune_UnmarshalJSON_non_empty(b *testing.B) {
	var tar s.Rune

	for ind := 0; ind < b.N; ind++ {
		gg.Try(tar.UnmarshalJSON(gg.ToBytes(`"👍"`)))
	}
}

func TestRune_Value(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Zero(gg.Try1(s.Rune(0).Value()))
	gtest.Equal(gg.Try1(s.Rune('👍').Value()), any(rune('👍')))
}

func TestRune_Scan(t *testing.T) {
	t.Run(`clear`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src any) {
			tar := s.Rune('👍')
			gtest.NoError(tar.Scan(src))
			gtest.Zero(tar)
		}

		test(nil)
		test(string(``))
		test([]byte(nil))
		test([]byte{})
		test(rune(0))
		test(s.Rune(0))
	})

	t.Run(`unclear`, func(t *testing.T) {
		defer gtest.Catch(t)

		test := func(src any, exp s.Rune) {
			var tar s.Rune
			gtest.NoError(tar.Scan(src))
			gtest.Eq(tar, exp)
		}

		test(string(`👍`), '👍')
		test([]byte(`👍`), '👍')
		test(rune('👍'), '👍')
		test(s.Rune('👍'), '👍')
	})
}

func BenchmarkRune_Scan_empty(b *testing.B) {
	var tar s.Rune
	var src []byte

	for ind := 0; ind < b.N; ind++ {
		gg.Try(tar.Scan(src))
	}
}

func BenchmarkRune_Scan_non_empty(b *testing.B) {
	var tar s.Rune
	src := string(`👍`)

	for ind := 0; ind < b.N; ind++ {
		gg.Try(tar.Scan(src))
	}
}

func Test_string_to_char_slice(t *testing.T) {
	defer gtest.Catch(t)

	gtest.Equal(
		[]rune(`👍👎🙂😄`),
		[]rune{'👍', '👎', '🙂', '😄'},
	)

	gtest.Equal(
		[]s.Rune(`👍👎🙂😄`),
		[]s.Rune{'👍', '👎', '🙂', '😄'},
	)
}
