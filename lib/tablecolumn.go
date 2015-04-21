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
	var result string
	for _, tp := range strings.Split(t.ColumnName, "_") {
		tp = strings.ToUpper(tp[0:1]) + strings.ToLower(tp[1:])
		result = result + tp
	}
	return result
}

func (t tableColumn) LowercaseColumnName() string {
	return strings.ToLower(t.ColumnName)
}

func (t tableColumn) SwaggerColumnType() string {
	ct := t.MappedColumnType()
	if strings.Contains(ct, "int") {
		return "integer"
	}
	if strings.Contains(ct, "float") {
		return "number"
	}
	if strings.Contains(t.ColumnType, "bit") {
		return "boolean"
	}

	return "string"
}

func (t tableColumn) SwaggerFormat() string {
	ct := t.MappedColumnType()
	if strings.Contains(ct, "int") {
		return ct
	}
	if strings.Contains(ct, "float") {
		return "float"
	}
	if strings.Contains(ct, "byte") {
		return "byte"
	}

	if strings.Contains(ct, "time") {
		return "date-time"
	}
	if strings.Contains(ct, "date") {
		return "date"
	}

	if strings.Contains(t.ColumnName, "password") {
		return "password"
	}

	return ""
}

func (t tableColumn) MappedColumnType() string {
	switch t.ColumnType {
	case "blob", "tinyblob", "mediumblob", "longblob",
		"binary", "varbinary", "bit", "set", "enum":
		return "[]byte"
	case "text", "tinytext", "mediumtext", "longtext", "char", "varchar":
		return "string"
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

func (t tableColumn) ColumnNullValue() string {
	switch t.ColumnType {
	case "blob", "tinyblob", "mediumblob", "longblob",
		"binary", "varbinary", "bit", "set", "enum":
		return "nil"
	case "char", "varchar", "text", "tinytext", "mediumtext", "longtext":
		return `""`
	case "tinyint", "utinyint", "smallint", "usmallint", "mediumint", "int", "umediumint", "uint", "bigint", "ubigint", "year", "float", "double", "decimal":
		return "0"
	case "time":
		return "nil"
	case "date":
		return "nil"
	case "datetime", "timestamp":
		return "nil"
	}

	return "nil"
}

func (t tableColumn) NullCheck(varname string) string {
	switch t.ColumnType {
	case "blob", "tinyblob", "mediumblob", "longblob",
		"binary", "varbinary", "bit", "set", "enum":

		return fmt.Sprintf("%s != nil", varname)
	case "text", "tinytext", "mediumtext", "longtext",
		"char", "varchar":
		return fmt.Sprintf(`%s != ""`, varname)

	case "tinyint", "utinyint", "smallint", "usmallint", "mediumint", "int", "umediumint", "uint", "bigint", "ubigint", "year", "float", "decimal", "double", "time":
		return fmt.Sprintf("%s != 0", varname)
	case "date", "datetime", "timestamp":
		return fmt.Sprintf("!%s.IsZero()", varname)
	}

	return "false"

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
