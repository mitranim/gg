package gsql_test

import (
	r "reflect"
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gsql"
	"github.com/mitranim/gg/gtest"
)

type ScalarRows[A any] struct {
	Rows []A
	Ind  int
}

func (self ScalarRows[_]) Columns() (_ []string, _ error) {
	return []string{`val`}, nil
}

func (self *ScalarRows[_]) Scan(tar ...any) error {
	if len(tar) != 1 {
		return gg.Errf(`expected 1 scan destination, got %v destinations`, len(tar))
	}
	if !self.Next() {
		return gg.Errf(`index mismatch`)
	}
	gg.ValueDerefAlloc(r.ValueOf(tar[0])).Set(r.ValueOf(self.Rows[self.Ind]))
	self.Ind++
	return nil
}

func (self ScalarRows[_]) Next() bool {
	return self.Ind >= 0 && self.Ind < len(self.Rows)
}

func (self ScalarRows[_]) Close() (_ error) { return }
func (self ScalarRows[_]) Err() (_ error)   { return }

func TestScanAny_scalar(t *testing.T) {
	defer gtest.Catch(t)

	{
		var rows ScalarRows[string]
		var tar string
		gtest.PanicStr(`no rows in result set`, func() { gsql.ScanAny(&rows, &tar) })
		gtest.Zero(tar)
	}

	{
		var rows ScalarRows[string]
		tar := `one`
		gtest.PanicStr(`no rows in result set`, func() { gsql.ScanAny(&rows, &tar) })
		gtest.Eq(tar, `one`)
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`, `two`, `three`}}
		var tar string
		gtest.PanicStr(`expected one row, got multiple`, func() { gsql.ScanAny(&rows, &tar) })
		gtest.Eq(tar, `one`)
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		var tar string
		gsql.ScanAny(&rows, &tar)
		gtest.Eq(tar, `one`)
	}

	{
		rows := ScalarRows[string]{Rows: []string{`two`}}
		tar := `one`
		gsql.ScanAny(&rows, &tar)
		gtest.Eq(tar, `two`)
	}
}

func TestScanAny_scalars(t *testing.T) {
	defer gtest.Catch(t)

	{
		var rows ScalarRows[string]
		var tar []string
		gsql.ScanAny(&rows, &tar)
		gtest.Zero(tar)
	}

	{
		var rows ScalarRows[string]
		tar := []string{`one`, `two`, `three`}
		gsql.ScanAny(&rows, &tar)
		gtest.Equal(tar, []string{`one`, `two`, `three`})
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`, `two`, `three`}}
		var tar []string
		gsql.ScanAny(&rows, &tar)
		gtest.Equal(tar, []string{`one`, `two`, `three`})
	}

	{
		rows := ScalarRows[string]{Rows: []string{`four`, `five`, `six`}}
		tar := []string{`one`, `two`, `three`}
		gsql.ScanAny(&rows, &tar)
		gtest.Equal(tar, []string{`one`, `two`, `three`, `four`, `five`, `six`})
	}
}

func TestScanAny_iface_to_scalar(t *testing.T) {
	defer gtest.Catch(t)

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		iface := any(``)

		gsql.ScanAny(&rows, &iface)
		gtest.Eq(iface, any(`one`))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`two`}}
		iface := any(`one`)

		gsql.ScanAny(&rows, &iface)
		gtest.Eq(iface, any(`two`))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		var tar string
		iface := any(tar)

		gsql.ScanAny(&rows, &iface)

		// Demonstration and reminder why other similar tests give the iface a
		// pointer to the local variable. It was impossible to modify `tar`
		// because the iface has its copy, not a pointer to it.
		gtest.Zero(tar)

		gtest.Eq(iface, any(`one`))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		var tar string
		iface := any(&tar)

		gsql.ScanAny(&rows, &iface)
		gtest.Eq(iface, any(&tar))
		gtest.Eq(tar, `one`)
		gtest.Equal(iface, any(gg.Ptr(`one`)))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`two`}}
		tar := `one`
		iface := any(&tar)

		gsql.ScanAny(&rows, &iface)
		gtest.Eq(iface, any(&tar))
		gtest.Eq(tar, `two`)
		gtest.Equal(iface, any(gg.Ptr(`two`)))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`two`}}
		tar := `one`
		ptr := &tar
		iface := any(&ptr)

		gsql.ScanAny(&rows, &iface)
		gtest.Eq(iface, any(&ptr))
		gtest.Eq(tar, `two`)
		gtest.Equal(iface, any(gg.Ptr(gg.Ptr(`two`))))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		iface := any((*string)(nil))

		gsql.ScanAny(&rows, &iface)
		gtest.Equal(iface, any(gg.Ptr(`one`)))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		iface := any((**string)(nil))

		gsql.ScanAny(&rows, &iface)
		gtest.Equal(iface, any(gg.Ptr(gg.Ptr(`one`))))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		var tar *string
		iface := any(&tar)

		gsql.ScanAny(&rows, &iface)
		gtest.Eq(iface, any(&tar))
		gtest.Equal(tar, gg.Ptr(`one`))
		gtest.Equal(iface, any(gg.Ptr(gg.Ptr(`one`))))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		var tar **string
		iface := any(&tar)

		gsql.ScanAny(&rows, &iface)
		gtest.Eq(iface, any(&tar))
		gtest.Equal(tar, gg.Ptr(gg.Ptr(`one`)))
		gtest.Equal(iface, any(gg.Ptr(gg.Ptr(gg.Ptr(`one`)))))
	}
}

func TestScanAny_iface_to_scalars(t *testing.T) {
	defer gtest.Catch(t)

	{
		rows := ScalarRows[string]{Rows: []string{`one`, `two`, `three`}}
		iface := any([]string(nil))

		gsql.ScanAny(&rows, &iface)
		gtest.Equal(iface, any([]string{`one`, `two`, `three`}))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`four`, `five`, `six`}}
		iface := any([]string{`one`, `two`, `three`})

		gsql.ScanAny(&rows, &iface)
		gtest.Equal(iface, any([]string{`one`, `two`, `three`, `four`, `five`, `six`}))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`, `two`, `three`}}
		var tar []string
		iface := any(&tar)
		gsql.ScanAny(&rows, &iface)
		gtest.Eq(iface, any(&tar))
		gtest.Equal(tar, []string{`one`, `two`, `three`})
	}

	{
		rows := ScalarRows[string]{Rows: []string{`four`, `five`, `six`}}
		tar := []string{`one`, `two`, `three`}
		iface := any(&tar)

		gsql.ScanAny(&rows, &iface)
		gtest.Eq(iface, any(&tar))
		gtest.Equal(tar, []string{`one`, `two`, `three`, `four`, `five`, `six`})
	}
}
