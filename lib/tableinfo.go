package lib

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

func (t tableInfo) CamelCaseTableName() string {
	tn := t.NormalizedTableName()
	return strings.ToLower(string(tn[0])) + tn[1:]
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

func (t tableInfo) CacheablePrimaryColumns() []tableColumn {
	result := make([]tableColumn, 0, len(t.ColOrder))
	for _, col := range t.ColOrder {
		if col.Primary && col.MappedColumnType() != "[]byte" {
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

func (t tableInfo) GenGetCache(tc []tableColumn) string {
	if len(tc) < 1 {
		return ""
	}
	result := fmt.Sprintf("if r0, ok := cache[%s]; ok {", tc[0].LowercaseColumnName())
	if len(tc) < 2 {
		return result + "return r0,nil }"
	}

	for i, c := range tc[1:] {
		result = result + fmt.Sprintf("if r%d, ok := r%d[%s]; ok {", i+1, i, c.LowercaseColumnName())
	}
	return result + fmt.Sprintf(" return r%d,nil ", len(tc)-1) + strings.Repeat("}", len(tc))
}
