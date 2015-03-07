package main

import (
	"fmt"
	"strings"
)

type tableInfo struct {
	TableName    string
	TableColumns map[string]tableColumn
	ColOrder     []tableColumn
	Constraints  []string
}

func (t tableInfo) QueryFieldNames() string {
	return strings.Join(columnNames(t.ColOrder), ",")
}

func (t tableInfo) QueryValuesSection() string {
	return strings.Join(strings.Split(strings.Repeat("?", len(t.TableColumns)), ""), ",")
}

func (t tableInfo) NormalizedTableName() string {
	var result string = ""
	for _, tp := range strings.Split(t.TableName, "_") {
		tp = strings.ToUpper(tp[0:1]) + strings.ToLower(tp[1:])
		tp = strings.TrimSuffix(tp, "s")
		result = result + tp
	}
	return result
}

func (t tableInfo) PrimaryColumns() []tableColumn {
	result := make([]tableColumn, 0, len(t.ColOrder))
	for _, col := range t.ColOrder {
		if col.Primary {

			result = append(result, col)
		}
	}
	return result
}

func (t tableInfo) PrimaryWhere() string {
	cols := columnNames(t.PrimaryColumns())
	for i, col := range cols {
		cols[i] = fmt.Sprintf("%s = ?", col)
	}
	return strings.Join(cols, " and ")
}

func (t tableInfo) PrimaryColumnsJoinedByAnd() string {
	return strings.Join(capitalizedColumnNames(t.PrimaryColumns()), "And")
}

func (t tableInfo) PrimaryColumnsParamList() string {
	return colformat(t.PrimaryColumns(), "%s %s", ",", lcn, mct)
}

func (t tableInfo) UpsertDuplicate() string {
	return colformat(t.ColOrder, "%s = VALUES(%s)", ",", lcn, lcn)
}
