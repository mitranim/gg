package gsql_test

import (
	"database/sql"
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

type Row struct {
	Name  string `db:"name"`
	Count int64  `db:"count"`
}

type StructRows struct {
	Rows []Row
	Ind  int
}

func (self StructRows) Columns() (_ []string, _ error) {
	return []string{`name`, `count`}, nil
}

func (self *StructRows) Scan(tar ...any) error {
	if len(tar) != 2 {
		return gg.Errf(`expected 2 scan destinations, got %v destinations`, len(tar))
	}
	if !self.Next() {
		return gg.Errf(`index mismatch`)
	}

	src := self.Rows[self.Ind]
	gg.ValueDerefAlloc(r.ValueOf(tar[0])).SetString(src.Name)
	gg.ValueDerefAlloc(r.ValueOf(tar[1])).SetInt(src.Count)

	self.Ind++
	return nil
}

func (self StructRows) Next() bool {
	return self.Ind >= 0 && self.Ind < len(self.Rows)
}

func (self StructRows) Close() (_ error) { return }
func (self StructRows) Err() (_ error)   { return }

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

func TestScanVal(t *testing.T) {
	defer gtest.Catch(t)

	{
		var rows ScalarRows[string]
		gtest.PanicErrIs(sql.ErrNoRows, func() {
			gsql.ScanVal[string](&rows)
		})
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		tar := gsql.ScanVal[string](&rows)
		gtest.Eq(tar, `one`)
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`, `two`}}
		gtest.PanicStr(string(gsql.ErrMultipleRows), func() {
			gsql.ScanVal[string](&rows)
		})
	}

	{
		var rows StructRows
		gtest.PanicErrIs(sql.ErrNoRows, func() {
			gsql.ScanVal[Row](&rows)
		})
	}

	{
		row := Row{Name: `one`, Count: 10}
		rows := StructRows{Rows: []Row{row}}
		gtest.Eq(gsql.ScanVal[Row](&rows), row)

		gtest.PanicErrIs(sql.ErrNoRows, func() {
			gsql.ScanVal[Row](&rows)
		})
	}

	{
		rows := StructRows{Rows: []Row{
			{Name: `one`, Count: 10},
			{Name: `two`, Count: 20},
		}}
		gtest.PanicStr(string(gsql.ErrMultipleRows), func() {
			gsql.ScanVal[Row](&rows)
		})
	}
}

func TestScanVals(t *testing.T) {
	defer gtest.Catch(t)

	{
		var rows ScalarRows[string]
		gtest.Zero(gsql.ScanVals[string](&rows))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`, `two`, `three`}}
		gtest.Equal(gsql.ScanVals[string](&rows), []string{`one`, `two`, `three`})
	}

	{
		var rows StructRows
		gtest.Zero(gsql.ScanVals[Row](&rows))
	}

	{
		row0 := Row{Name: `one`, Count: 10}
		row1 := Row{Name: `two`, Count: 20}
		rows := StructRows{Rows: []Row{row0, row1}}
		gtest.Equal(gsql.ScanVals[Row](&rows), []Row{row0, row1})
	}
}

func TestScanAnyOpt_scalar(t *testing.T) {
	defer gtest.Catch(t)

	{
		var rows ScalarRows[string]
		var tar string
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Zero(tar)
	}

	{
		var rows ScalarRows[string]
		tar := `one`
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Eq(tar, `one`)
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`, `two`, `three`}}
		var tar string
		gtest.PanicStr(`expected one row, got multiple`, func() { gsql.ScanAnyOpt(&rows, &tar) })
		gtest.Eq(tar, `one`)
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		var tar string
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Eq(tar, `one`)
	}

	{
		rows := ScalarRows[string]{Rows: []string{`two`}}
		tar := `one`
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Eq(tar, `two`)
	}
}

func TestScanAnyOpt_scalars(t *testing.T) {
	defer gtest.Catch(t)

	{
		var rows ScalarRows[string]
		var tar []string
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Zero(tar)
	}

	{
		var rows ScalarRows[string]
		tar := []string{`one`, `two`, `three`}
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Equal(tar, []string{`one`, `two`, `three`})
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`, `two`, `three`}}
		var tar []string
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Equal(tar, []string{`one`, `two`, `three`})
	}

	{
		rows := ScalarRows[string]{Rows: []string{`four`, `five`, `six`}}
		tar := []string{`one`, `two`, `three`}
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Equal(tar, []string{`one`, `two`, `three`, `four`, `five`, `six`})
	}
}

func TestScanAnyOpt_iface_to_scalar(t *testing.T) {
	defer gtest.Catch(t)

	{
		var rows ScalarRows[string]
		iface := any(``)
		gsql.ScanAnyOpt(&rows, &iface)
		gtest.Eq(iface, any(``))
	}

	{
		var rows ScalarRows[string]
		iface := any(`one`)
		gsql.ScanAnyOpt(&rows, &iface)
		gtest.Eq(iface, any(`one`))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		iface := any(``)
		gsql.ScanAnyOpt(&rows, &iface)
		gtest.Eq(iface, any(`one`))
	}

	{
		rows := ScalarRows[string]{Rows: []string{`two`}}
		iface := any(`one`)
		gsql.ScanAnyOpt(&rows, &iface)
		gtest.Eq(iface, any(`two`))
	}

	{
		var rows ScalarRows[string]
		var tar string
		iface := any(&tar)
		gsql.ScanAnyOpt(&rows, &iface)
		gtest.Eq(iface, any(&tar))
		gtest.Zero(tar)
	}

	{
		rows := ScalarRows[string]{Rows: []string{`one`}}
		var tar string
		iface := any(&tar)
		gsql.ScanAnyOpt(&rows, &iface)
		gtest.Eq(iface, any(&tar))
		gtest.Eq(tar, `one`)
		gtest.Equal(iface, any(gg.Ptr(`one`)))
	}
}

func TestScanAnyOpt_struct(t *testing.T) {
	defer gtest.Catch(t)

	{
		var rows StructRows
		var tar Row
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Zero(tar)
	}

	{
		var rows StructRows
		row := Row{Name: `one`, Count: 10}
		tar := row
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Eq(tar, row)
	}

	{
		row := Row{Name: `one`, Count: 10}
		rows := StructRows{Rows: []Row{row}}
		var tar Row
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Eq(tar, row)
	}

	{
		row := Row{Name: `one`, Count: 10}
		rows := StructRows{Rows: []Row{row}}
		tar := Row{Name: `two`, Count: 20}
		gsql.ScanAnyOpt(&rows, &tar)
		gtest.Eq(tar, row)
	}

	{
		row0 := Row{Name: `one`, Count: 10}
		row1 := Row{Name: `two`, Count: 20}
		rows := StructRows{Rows: []Row{row0, row1}}
		var tar Row
		gtest.PanicStr(string(gsql.ErrMultipleRows), func() {
			gsql.ScanAnyOpt(&rows, &tar)
		})
		gtest.Eq(tar, row0)
	}
}
