package lib

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
	case "text", "tinytext", "mediumtext", "longtext",
		"blob", "tinyblob", "mediumblob", "longblob",
		"binary", "varbinary", "bit", "set", "enum",
		"char", "varchar":
		return "[]byte"
	case "tinyint":
		return "int8"
	case "utinyint":
		return "uint8"
	case "smallint":
		return "int16"
	case "usmallint":
		return "uint16"
	case "mediumint", "int":
		return "int32"
	case "umediumint", "uint":
		return "uint32"
	case "bigint":
		return "int64"
	case "ubigint":
		return "uint64"
	case "year":
		return "int16"
	case "time":
		return "time.Duration"
	case "date":
		return "mysql.Date"
	case "datetime", "timestamp":
		return "time.Time"
	case "float":
		return "float32"
	case "decimal", "double":
		return "float64"

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
		return `nil`
	case "tinyint":
		return "0"
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
