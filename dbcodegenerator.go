package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"
)

type dbCodeGenerator struct {
}

func (g *dbCodeGenerator) Generate(tables map[string]tableInfo) error {
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
//"errors"
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

func DeleteBy{{.PrimaryColumnsJoinedByAnd}}({{.PrimaryColumnsParamList}}) (*{{.NormalizedTableName}}, error) {
    //TODO: remove from cache.
    err := DB.Exec("DELETE FROM {{.TableName}} WHERE {{.PrimaryWhere}}",
    {{range .PrimaryColumns}}{{.LowercaseColumnName}},
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
