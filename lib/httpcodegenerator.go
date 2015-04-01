package lib

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"
)

type httpCodeGenerator struct {
	DbRootPackageName string
}

func (g *httpCodeGenerator) Generate(tables map[string]tableInfo) error {
	err := os.Mkdir("http", 0755)
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("httphandlers").Parse(`//WARNING.
//THIS HAS BEEN GENERATED AUTOMATICALLY BY AUTOAPI.
//IF THERE WAS A WARRANTY, MODIFYING THIS WOULD VOID IT.

package {{.Table.TableName}}

import (
"{{.DbRootPackageName}}/{{.Table.TableName}}"

"net/http"
"fmt"
"encoding/json"
"github.com/gorilla/mux"
)



func List(res http.ResponseWriter, req *http.Request){
    enc := json.NewEncoder(res)
    rows, _ := {{.Table.TableName}}.All()
    enc.Encode(rows)
}

func Get(res http.ResponseWriter, req *http.Request){
    vars := mux.Vars(req)
    idstring := vars["id"]
 
    enc := json.NewEncoder(res)
    {{$l := len .Table.PrimaryColumns}}
    {{if gt $l 1}}
        id_slice := strings.Split(vars["id"])
    {{else}}
        param := vars["id"]
        {{.FirstPrimaryColumnTypeConverter}}
    {{end}}
    row, _ := {{.Table.TableName}}.GetBy{{.Table.PrimaryColumnsJoinedByAnd}}(id)
    enc.Encode(row)
}

func Post(res http.ResponseWriter, req *http.Request){
    dec := json.NewDecoder(req.Body)
    row := &{{.Table.TableName}}.{{.Table.NormalizedTableName}}{}
    dec.Decode(&row)
    {{.Table.TableName}}.Save(row)
}

func Put(res http.ResponseWriter, req *http.Request){
    dec := json.NewDecoder(req.Body)
    row := &{{.Table.TableName}}.{{.Table.NormalizedTableName}}{}
    dec.Decode(&row)
    {{.Table.TableName}}.Save(row)
}

func Delete(res http.ResponseWriter, req *http.Request){
    fmt.Fprintf(res, "Delete stub!")
}

`))

	path, err := GetRootPath()
	if err != nil {
		return err
	}
	for table, tinfo := range tables {

		os.Mkdir("http/"+table, 0755)
		f, err := os.Create("http/" + table + "/" + table + ".go")
		if err != nil {
			return err
		}
		var b bytes.Buffer
		err = tmpl.Execute(&b, map[string]interface{}{"Table": tinfo, "DbRootPackageName": path + "/db"})
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

	tmpl = tmpl
	return nil
}
