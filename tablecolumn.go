package main

import (
	"fmt"
	"strings"
)

type tableColumn struct {
	ColumnName string
	ColumnType string
	Primary    bool
}

func (t tableColumn) CapitalizedColumnName() string {
	var result string = ""
	for _, tp := range strings.Split(t.ColumnName, "_") {
		tp = strings.ToUpper(tp[0:1]) + strings.ToLower(tp[1:])
		result = result + tp
	}
	return result
}

func (t tableColumn) LowercaseColumnName() string {
	return strings.ToLower(t.ColumnName)
}

func (tc tableColumn) MappedColumnType() string {
	switch tc.ColumnType {
	case "text", "tinytext", "mediumtext", "longtex",
		"blob", "tinyblob", "mediumblob", "longblob",
		"binary", "varbinary":
		return "[]byte"
	case "char", "varchar":
		return "string"

	}
	return "interface{}"
}

func (tc tableColumn) ColumnNullValue() string {
	switch tc.ColumnType {
	case "text", "tinytext", "mediumtext", "longtex",
		"blob", "tinyblob", "mediumblob", "longblob",
		"binary", "varbinary":
		return `nil`
	case "char", "varchar":
		return `""`

	}
	return "nil"
}

func colformat(cols []tableColumn, format string, joinstring string, str1, str2 func(tableColumn) string) string {
	result := []string{}
	for _, col := range cols {
		result = append(result, fmt.Sprintf(format, str1(col), str2(col)))
	}
	return strings.Join(result, joinstring)
}

func lcn(c tableColumn) string { return c.LowercaseColumnName() }
func mct(c tableColumn) string { return c.MappedColumnType() }

func capitalizedColumnNames(cols []tableColumn) []string {
	result := []string{}
	for _, c := range cols {
		result = append(result, c.CapitalizedColumnName())
	}
	return result
}

func columnNames(cols []tableColumn) []string {
	result := []string{}
	for _, c := range cols {
		result = append(result, c.ColumnName)
	}
	return result
}
