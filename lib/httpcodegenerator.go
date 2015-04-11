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
    rows, err := {{.Table.TableName}}.All()
    if err != nil {
        log.Println(err)
    }
    enc.Encode(rows)
}

{{$l := len .Table.PrimaryColumns}}
{{if gt $l 0}}
func Get(res http.ResponseWriter, req *http.Request){
    vars := mux.Vars(req)

 
    enc := json.NewEncoder(res)

    {{if gt $l 1}}
        id_slice := strings.Split(vars["id"])
        param := id_slice["id"]
       {{.Table.FirstPrimaryColumnTypeConverter}}
    {{else}}
        param := vars["id"]
       {{.Table.FirstPrimaryColumnTypeConverter}} 
   {{end}}
 
    row, _ := {{.Table.TableName}}.GetBy{{.Table.PrimaryColumnsJoinedByAnd}}(id)
    enc.Encode(row)
}
{{end}}

func Post(res http.ResponseWriter, req *http.Request){
    save(req)
}

func Put(res http.ResponseWriter, req *http.Request){
    vars := mux.Vars(req)
    {{if gt $l 1}}
        id_slice := strings.Split(vars["id"])
        param := id_slice["id"]
       {{.Table.FirstPrimaryColumnTypeConverter}}
    {{else}}
        param := vars["id"]
       {{.Table.FirstPrimaryColumnTypeConverter}} 
   {{end}}
    row, err := {{.Table.TableName}}.GetBy{{.Table.PrimaryColumnsJoinedByAnd}}(id)
    if err != nil {
        fmt.Println(err)
        fmt.Fprintln(res, err)
        return
    }
    save(req)
}

func save(req *http.Request) error {
    dec := json.NewDecoder(req.Body)
    row := &{{.Table.TableName}}.{{.Table.NormalizedTableName}}{}
    dec.Decode(&row)
    return {{.Table.TableName}}.Save(row)
}

func Delete(res http.ResponseWriter, req *http.Request){
    vars := mux.Vars(req)
    {{if gt $l 1}}
        id_slice := strings.Split(vars["id"])
        param := id_slice["id"]
       {{.Table.FirstPrimaryColumnTypeConverter}}
    {{else}}
        param := vars["id"]
       {{.Table.FirstPrimaryColumnTypeConverter}} 
   {{end}}
    {{.Table.TableName}}.DeleteBy{{.Table.PrimaryColumnsJoinedByAnd}}(id)
}

`))

	path, err := GetRootPath()
	if err != nil {
		return err
	}
	for table, tinfo := range tables {

		os.Mkdir("http/"+table, 0755)
		f, err := os.Open("http/" + table + "/" + table + ".go")
		if err != nil && !os.IsExist(err) {

			return err
		}
		var b bytes.Buffer
		err = tmpl.Execute(&b, map[string]interface{}{"Table": tinfo, "DbRootPackageName": path + "/db"})
		if err != nil {
			fmt.Println(b.String())
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

	tmpl = tmpl
	return nil
}
