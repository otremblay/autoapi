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
	if err != nil && !os.IsExist(err) {
		return err
	}
	tmpl := template.Must(template.New("httphandlers").Parse(`//WARNING.
//THIS HAS BEEN GENERATED AUTOMATICALLY BY AUTOAPI.
//IF THERE WAS A WARRANTY, MODIFYING THIS WOULD VOID IT.

package {{.Table.TableName}}

import (
"{{.DbRootPackageName}}/{{.Table.TableName}}"
dbi "{{.dbiroot}}/{{.Table.TableName}}"
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
    err := save(req)
    if err != nil {
        res.WriteHeader(500)
        fmt.Fprint(res, err)
    }
}
{{if gt $l 0}}
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
    _, get_err := {{.Table.TableName}}.GetBy{{.Table.PrimaryColumnsJoinedByAnd}}(id)
    if get_err != nil {
        fmt.Println(get_err)
        fmt.Fprintln(res, get_err)
        return
    }
    err = save(req)
    if err != nil {
        res.WriteHeader(500)
        fmt.Fprint(res, err)
    }
}
{{end}}
func save(req *http.Request) error {
    dec := json.NewDecoder(req.Body)
    row := &dbi.{{.Table.NormalizedTableName}}{}
    dec.Decode(&row)
    return {{.Table.TableName}}.Save(row)
}
{{if gt $l 0}}
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
{{end}}
`))

	path, err := GetRootPath()
	if err != nil {
		return err
	}
	for table, tinfo := range tables {

		os.Mkdir("http/"+table, 0755)
		f, err := os.Create("http/" + table + "/" + table + ".go")
		if err != nil && !os.IsExist(err) {

			return err
		}
		var b bytes.Buffer
		err = tmpl.Execute(&b, map[string]interface{}{"Table": tinfo, "DbRootPackageName": path + "/" + rootdbpath, "dbiroot": path + "/dbi"})
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
