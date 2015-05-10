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

type dbiCodeGenerator struct {
}

func (g *dbiCodeGenerator) Generate(tables map[string]tableInfo) error {
	err := os.Mkdir("dbi", 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	tmpl := template.Must(template.New("class").Parse(`//WARNING.
//THIS HAS BEEN GENERATED AUTOMATICALLY BY AUTOAPI.
//IF THERE WAS A WARRANTY, MODIFYING THIS WOULD VOID IT.

package {{.TableName}}

type {{.NormalizedTableName}}er interface {
 FindWithWhere(where string, params ...interface{}) ([]*{{.NormalizedTableName}}, error)
 GetBy{{.PrimaryColumnsJoinedByAnd}}({{.PrimaryColumnsParamList}}) (*{{.NormalizedTableName}}, error)
 All() ([]*{{.NormalizedTableName}}, error)
 Find({{.TableName}} *{{.NormalizedTableName}}) ([]*{{.NormalizedTableName}}, error)
{{if .PrimaryColumns }}
 DeleteBy{{.PrimaryColumnsJoinedByAnd}}({{.PrimaryColumnsParamList}}) (error)
{{end}}
 Save(row *{{.NormalizedTableName}}) error
}

type {{.NormalizedTableName}} struct {
{{range .ColOrder}}{{.CapitalizedColumnName}} {{.MappedColumnType}} ` + "`json:\"{{.ColumnName}}\"`" + `
{{end}}}

func New() *{{.NormalizedTableName}}{
    return &{{.NormalizedTableName}}{}
}

`))
	for table, tinfo := range tables {
		os.Mkdir("dbi/"+table, 0755)
		f, err := os.Create("dbi/" + table + "/" + table + ".go")
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
