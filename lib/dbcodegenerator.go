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
	err := os.MkdirAll("db/mysql", 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	tmpl := template.Must(template.New("class").Parse(`//WARNING.
//THIS HAS BEEN GENERATED AUTOMATICALLY BY AUTOAPI.
//IF THERE WAS A WARRANTY, MODIFYING THIS WOULD VOID IT.

package {{.Table.TableName}}

import (
"is-a-dev.com/autoapi/lib"
dbi "{{.dbipackage}}"
//"errors"
)

var DB lib.DB

{{if .Table.CacheablePrimaryColumns}}
//type {{.Table.NormalizedTableName}}Cache struct{



//}

var cache = map{{range .Table.CacheablePrimaryColumns}}[{{.MappedColumnType}}]{{end}}*dbi.{{.Table.NormalizedTableName}}{}

{{end}}


func FindWithWhere(where string, params ...interface{}) ([]*dbi.{{.Table.NormalizedTableName}}, error) {
    rows, err := DB.Query("SELECT {{.Table.QueryFieldNames}} FROM {{.Table.TableName}} " + where, params...)
    if err != nil {
        return nil,err
    }
    result := make([]*dbi.{{.Table.NormalizedTableName}},0)
    for rows.Next() {
        r := &dbi.{{.Table.NormalizedTableName}}{}
        rows.Scan(
            {{range .Table.ColOrder}}&r.{{.CapitalizedColumnName}},
            {{end}})
        {{if .Table.CacheablePrimaryColumns}}
          cache{{range .Table.CacheablePrimaryColumns}}[r.{{.CapitalizedColumnName}}]{{end}} = r
        {{end}}
        result = append(result, r)
    }
    return result, nil
}

func All() ([]*dbi.{{.Table.NormalizedTableName}}, error){
    return FindWithWhere("")
}

func GetBy{{.Table.PrimaryColumnsJoinedByAnd}}({{.Table.PrimaryColumnsParamList}}) (*dbi.{{.Table.NormalizedTableName}}, error) {
    {{if .Table.CacheablePrimaryColumns}}
      {{.Table.GenGetCache .CacheablePrimaryColumns}} 
    {{end}}
    row := &dbi.{{.Table.NormalizedTableName}}{}
    err := DB.QueryRow("SELECT {{.Table.QueryFieldNames}} FROM {{.Table.TableName}} WHERE {{.Table.PrimaryWhere}}",
    {{range .Table.PrimaryColumns}}{{.LowercaseColumnName}},
    {{end}}).Scan(
        {{range .Table.ColOrder}}&row.{{.CapitalizedColumnName}},
        {{end}})
    if err != nil {
        return nil, err
    }
    return row, nil
}

func Find({{.Table.TableName}} *dbi.{{.Table.NormalizedTableName}}) ([]*dbi.{{.Table.NormalizedTableName}}, error){
    where := []string{}
    params := []interface{}{}
{{$tn := .Table.TableName}}
{{range .Table.ColOrder}}
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


{{if .Table.PrimaryColumns }}
func DeleteBy{{.Table.PrimaryColumnsJoinedByAnd}}({{.Table.PrimaryColumnsParamList}}) (error) {
    //TODO: remove from cache.
    _, err := DB.Exec("DELETE FROM {{.Table.TableName}} WHERE {{.Table.PrimaryWhere}}",
    {{range .Table.PrimaryColumns}}{{.LowercaseColumnName}},
    {{end}})
    if err != nil {
        return err
    }
    return nil
}
{{end}}

func Save(row *dbi.{{.Table.NormalizedTableName}}) error {
    {{range .Table.Constraints}}{{.}};{{end}}
    _, err := DB.Exec("INSERT {{.Table.TableName}} VALUES({{.Table.QueryValuesSection}}) ON DUPLICATE KEY UPDATE {{.Table.UpsertDuplicate}}", 
        {{range .Table.ColOrder}}row.{{.CapitalizedColumnName}},
{{end}})
    if err != nil {return err}
        {{if .Table.CacheablePrimaryColumns}}
          cache{{range .Table.CacheablePrimaryColumns}}[row.{{.CapitalizedColumnName}}]{{end}} = row
        {{end}}
    return nil
}
`))
	path, err := GetRootPath()
	if err != nil {
		return err
	}
	for table, tinfo := range tables {
		os.Mkdir("db/mysql/"+table, 0755)
		f, err := os.Create("db/mysql/" + table + "/" + table + ".go")
		if err != nil && !os.IsExist(err) {
			return err
		}
		var b bytes.Buffer
		err = tmpl.Execute(&b, map[string]interface{}{"Table": tinfo, "dbipackage": path + "/dbi/" + table})
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
