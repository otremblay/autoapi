package lib

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"

	"golang.org/x/tools/imports"
)

type dbCodeGenerator struct {
}

func (g *dbCodeGenerator) Generate(tables map[string]tableInfo) error {
	err := os.Mkdir("db", 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	tmpl := template.Must(template.New("class").Parse(`//WARNING.
//THIS HAS BEEN GENERATED AUTOMATICALLY BY AUTOAPI.
//IF THERE WAS A WARRANTY, MODIFYING THIS WOULD VOID IT.

package {{.TableName}}

import (
"is-a-dev.com/autoapi/lib"
//"errors"
)

var DB lib.DB

{{if .CacheablePrimaryColumns}}
//type {{.NormalizedTableName}}Cache struct{

//    rowsByKey map{{range .PrimaryColumns}}[{{.MappedColumnType}}]{{end}}*{{.NormalizedTableName}}

//}

var cache = map{{range .CacheablePrimaryColumns}}[{{.MappedColumnType}}]{{end}}*{{.NormalizedTableName}}{}

{{end}}

type {{.NormalizedTableName}} struct {
{{range .ColOrder}}{{.CapitalizedColumnName}} {{.MappedColumnType}}
{{end}}}

func New() *{{.NormalizedTableName}}{
    return &{{.NormalizedTableName}}{}
}

func FindWithWhere(where string, params ...interface{}) ([]*{{.NormalizedTableName}}, error) {
    rows, err := DB.Query("SELECT {{.QueryFieldNames}} FROM {{.TableName}} " + where, params...)
    if err != nil {
        return nil,err
    }
    result := make([]*{{.NormalizedTableName}},0)
    for rows.Next() {
        r := &{{.NormalizedTableName}}{}
        rows.Scan(
            {{range .ColOrder}}&r.{{.CapitalizedColumnName}},
            {{end}})
        {{if .CacheablePrimaryColumns}}
          cache{{range .CacheablePrimaryColumns}}[r.{{.CapitalizedColumnName}}]{{end}} = r
        {{end}}
        result = append(result, r)
    }
    return result, nil
}

func All() ([]*{{.NormalizedTableName}}, error){
    return FindWithWhere("")
}

func GetBy{{.PrimaryColumnsJoinedByAnd}}({{.PrimaryColumnsParamList}}) (*{{.NormalizedTableName}}, error) {
    {{if .CacheablePrimaryColumns}}
      {{.GenGetCache .CacheablePrimaryColumns}} 
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

func Find({{.TableName}} *{{.NormalizedTableName}}) ([]*{{.NormalizedTableName}}, error){
    where := []string{}
    params := []interface{}{}
{{$tn := .TableName}}
{{range .ColOrder}}
    if {{printf "%s%s%s" $tn "." .CapitalizedColumnName | .NullCheck}} {
        where = append(where , "{{.ColumnName}} = ?")
        params = append(params, {{$tn}}.{{.CapitalizedColumnName}})
    }
{{end}}
var resultingwhere string
if len(where)>0{
    resultingwhere = fmt.Sprintf("WHERE %s",strings.Join(where," AND "))
}
    return FindWithWhere(resultingwhere, params...)
}


{{if .PrimaryColumns }}
func DeleteBy{{.PrimaryColumnsJoinedByAnd}}({{.PrimaryColumnsParamList}}) (error) {
    //TODO: remove from cache.
    _, err := DB.Exec("DELETE FROM {{.TableName}} WHERE {{.PrimaryWhere}}",
    {{range .PrimaryColumns}}{{.LowercaseColumnName}},
    {{end}})
    if err != nil {
        return err
    }
    return nil
}
{{end}}

func Save(row *{{.NormalizedTableName}}) error {
    {{range .Constraints}}{{.}}{{end}}
    _, err := DB.Exec("INSERT {{.TableName}} VALUES({{.QueryValuesSection}}) ON DUPLICATE KEY UPDATE {{.UpsertDuplicate}}", 
        {{range .ColOrder}}row.{{.CapitalizedColumnName}},
{{end}})
    if err != nil {return err}
        {{if .CacheablePrimaryColumns}}
          cache{{range .CacheablePrimaryColumns}}[row.{{.CapitalizedColumnName}}]{{end}} = row
        {{end}}
    return nil
}
`))
	for table, tinfo := range tables {
		os.Mkdir("db/"+table, 0755)
		f, err := os.Create("db/" + table + "/" + table + ".go")
		if err != nil && !os.IsExist(err) {
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
		bf, err = imports.Process(f.Name(), bf, nil)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, bytes.NewBuffer(bf))
		if err != nil {
			return err
		}

	}
	return nil
}
