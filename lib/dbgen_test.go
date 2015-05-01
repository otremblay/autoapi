package lib

import (
	"os"
	"testing"
)

func TestNoTablesNoFilesNoProblem(t *testing.T) {
	os.RemoveAll("db")
	dg := &dbCodeGenerator{}
	err := dg.Generate(nil)
	if err != nil {
		t.Error(err)
	}
	err = dg.Generate(map[string]tableInfo{})
	if err != nil {
		t.Error(err)
	}
	os.RemoveAll("db")
}

func TestOneTableOneColumn(t *testing.T) {
	os.RemoveAll("db")
	ti := tableInfo{}
	ti.TableName = "test"
	tc := tableColumn{}
	tc.ColumnName = "nope"
	tc.ColumnType = "varchar"
	tc.Primary = true
	ti.TableColumns = map[string]tableColumn{"nope": tc}
	ti.ColOrder = []tableColumn{tc}
	dg := &dbCodeGenerator{}
	err := dg.Generate(map[string]tableInfo{"test": ti})
	if err != nil {
		t.Error(err)
	}
	os.RemoveAll("db")
}
