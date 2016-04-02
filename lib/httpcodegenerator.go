package lib

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

type httpCodeGenerator struct {
	DbRootPackageName string
	Verbs             string
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

{{$l := len .Table.PrimaryColumns}}
{{if .shouldGenerateGet}}
func List(res http.ResponseWriter, req *http.Request){
    enc := json.NewEncoder(res)
    
    shouldFilter := false
    filterObject := &dbi.{{.Table.NormalizedTableName}}{}
    req.ParseForm()
    if len(req.Form) > 0 {
        {{range .Table.ColOrder}}
        if _, ok := req.Form["{{.LowercaseColumnName}}"]; ok {
            form_{{.LowercaseColumnName}} := req.FormValue("{{.LowercaseColumnName}}")
            shouldFilter = true
            {{.TextRightHandConvert}}
            filterObject.{{.CapitalizedColumnName}} = parsedField 
        }
        {{end}}
    }

    if shouldFilter {
        rows, err := {{.Table.TableName}}.Find(filterObject)
        if err != nil {
            log.Println(err)
        }
        enc.Encode(rows)
        return
    }

    rows, err := {{.Table.TableName}}.All()
    if err != nil {
        log.Println(err)
    }
    enc.Encode(rows)
}


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
{{end}}

{{if .shouldGeneratePost}}
func Post(res http.ResponseWriter, req *http.Request){
    err := save(req)
    if err != nil {
        res.WriteHeader(500)
        fmt.Fprint(res, err)
    }
}
{{end}}

{{if .shouldGeneratePut}}
{{if gt $l 0}}
func Put(res http.ResponseWriter, req *http.Request){
    var err error
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
{{end}}
{{if or .shouldGeneratePut .shouldGeneratePost}}
func save(req *http.Request) error {
    dec := json.NewDecoder(req.Body)
    row := &dbi.{{.Table.NormalizedTableName}}{}
    dec.Decode(&row)
    return {{.Table.TableName}}.Save(row)
}
{{end}}
{{if .shouldGenerateDelete}}
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
		verbs := strings.Split(g.Verbs, ",")
		err = tmpl.Execute(&b, map[string]interface{}{
			"Table":                tinfo,
			"DbRootPackageName":    path + "/" + rootdbpath,
			"dbiroot":              path + "/dbi",
			"shouldGenerateGet":    stringInSlice("get", verbs),
			"shouldGeneratePost":   stringInSlice("post", verbs),
			"shouldGeneratePut":    stringInSlice("put", verbs),
			"shouldGenerateDelete": stringInSlice("delete", verbs),
		})
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
