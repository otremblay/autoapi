package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	"text/template"

	"go/format"

	_ "github.com/ziutek/mymysql/godrv"
)

func main() {
	dbUrl := os.Args[1]
	dbName := os.Args[2]
	db, err := sql.Open("mymysql", dbUrl)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select table_name from information_schema.tables where table_schema = ?", dbName)
	if err != nil {
		panic(err)
	}
	tables := map[string]tableInfo{}
	for rows.Next() {
		var tn string
		rows.Scan(&tn)
		tables[tn] = tableInfo{
			TableName:    tn,
			TableColumns: map[string]tableColumn{},
			ColOrder:     []tableColumn{},
			Constraints:  []string{},
		}
	}

	more_rows, err := db.Query("select table_name, column_name, data_type, column_key, is_nullable, extra from information_schema.columns where table_schema = ?", dbName)

	if err != nil {
		panic(err)
	}

	for more_rows.Next() {
		var tn, cn, ct, ck, nullable, extra string
		err := more_rows.Scan(&tn, &cn, &ct, &ck, &nullable, &extra)
		if err != nil {
			panic(err)
		}
		table := tables[tn]
		col := tableColumn{ColumnName: cn, ColumnType: ct}

		col.Primary = ck == "PRI"
		fmt.Println(col, ck)
		if nullable == "NO" && extra != "auto_increment" {
			table.Constraints = append(table.Constraints, fmt.Sprintf(`if row.%s == %s {return errors.New("Preconditions failed, %s must be set.")}`, col.CapitalizedColumnName(), col.ColumnNullValue(), col.CapitalizedColumnName()))
		}
		table.TableColumns[cn] = col
		table.ColOrder = append(table.ColOrder, col)
		tables[tn] = table
	}

	err = (&generator{}).Generate(tables)
	if err != nil {
		panic(err)
	}
}

type tableInfo struct {
	TableName    string
	TableColumns map[string]tableColumn
	ColOrder     []tableColumn
	Constraints  []string
}

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

func colformat(cols []tableColumn, format string, joinstring string, str1, str2 func(tableColumn) string) string {
	result := []string{}
	for _, col := range cols {
		result = append(result, fmt.Sprintf(format, str1(col), str2(col)))
	}
	return strings.Join(result, joinstring)
}

func lcn(c tableColumn) string { return c.LowercaseColumnName() }
func mct(c tableColumn) string { return c.MappedColumnType() }

func (t tableInfo) PrimaryColumnsParamList() string {
	return colformat(t.PrimaryColumns(), "%s %s", ",", lcn, mct)
}

func (t tableInfo) UpsertDuplicate() string {
	return colformat(t.ColOrder, "%s = VALUES(%s)", ",", lcn, lcn)
}

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

type generator struct {
}

func (g *generator) Generate(tables map[string]tableInfo) error {
	err := os.Mkdir("db", 0755)
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New("class").Parse(`//WARNING.
//THIS HAS BEEN GENERATED AUTOMATICALLY BY AUTOAPI.
//IF THERE WAS A WARRANTY, MODIFYING THIS WOULD VOID IT.

package {{.TableName}}

import (
"is-a-dev.com/libautoapi"
"errors"
)

var DB libautoapi.DB

{{if .PrimaryColumns}}
//type {{.NormalizedTableName}}Cache struct{

//    rowsByKey map{{range .PrimaryColumns}}[{{.MappedColumnType}}]{{end}}*{{.NormalizedTableName}}

//}

var cache = map{{range .PrimaryColumns}}[{{.MappedColumnType}}]{{end}}*{{.NormalizedTableName}}{}

{{end}}

type {{.NormalizedTableName}} struct {
{{range .ColOrder}}{{.CapitalizedColumnName}} {{.MappedColumnType}}
{{end}}}

func New() *{{.NormalizedTableName}}{
    return &{{.NormalizedTableName}}{}
}

func All() ([]*{{.NormalizedTableName}}, error){
    rows, err := DB.Query("SELECT {{.QueryFieldNames}} FROM {{.TableName}}")
    if err != nil {
        return nil,err
    }
    result := make([]*{{.NormalizedTableName}},0)
    for rows.Next() {
        r := &{{.NormalizedTableName}}{}
        rows.Scan(
            {{range .ColOrder}}&r.{{.CapitalizedColumnName}},
            {{end}})
        {{if .PrimaryColumns}}
          cache[r.{{range .PrimaryColumns}}{{.CapitalizedColumnName}}{{end}}] = r
        {{end}}
        result = append(result, r)
    }
    return result, nil
}

func GetBy{{.PrimaryColumnsJoinedByAnd}}({{.PrimaryColumnsParamList}}) (*{{.NormalizedTableName}}, error) {
    {{if .PrimaryColumns}}
        if r, ok := cache[{{range .PrimaryColumns}}{{.LowercaseColumnName}}{{end}}]; ok { return r, nil}
    {{end}}
    row := &{{.NormalizedTableName}}{}
    err := DB.QueryRow("SELECT {{.QueryFieldNames}} FROM {{.TableName}} WHERE {{.PrimaryWhere}}",
    {{range .PrimaryColumns}}{{.LowercaseColumnName}},
    {{end}}).Scan(
        {{range .ColOrder}}&row.{{.CapitalizedColumnName}},
        {{end}})
    if err != nil {
        return nil, err
    }
    return row, nil
}

func Save(row *{{.NormalizedTableName}}) error {
    {{range .Constraints}}{{.}}{{end}}
    _, err := DB.Exec("INSERT {{.TableName}} VALUES({{.QueryValuesSection}}) ON DUPLICATE KEY UPDATE {{.UpsertDuplicate}}", 
        {{range .ColOrder}}row.{{.CapitalizedColumnName}},
{{end}})
    if err != nil {return err}
        {{if .PrimaryColumns}}
          cache[row.{{range .PrimaryColumns}}{{.CapitalizedColumnName}}{{end}}] = row
        {{end}}
    return nil
}
`))
	for table, tinfo := range tables {
		fmt.Println(tinfo.PrimaryColumns())
		os.Mkdir("db/"+table, 0755)
		f, err := os.Create("db/" + table + "/" + table + ".go")
		if err != nil {
			return err
		}
		var b bytes.Buffer
		err = tmpl.Execute(&b, tinfo)
		if err != nil {
			return err
		}
		bf, err := format.Source(b.Bytes())
		if err != nil {
			fmt.Println(b.String())
			return err
		}
		_, err = io.Copy(f, bytes.NewBuffer(bf))
		if err != nil {
			return err
		}
	}
	return nil
}
