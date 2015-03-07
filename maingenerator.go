package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"os"
)

type mainGenerator struct {
	rootDbPackageName       string
	rootHandlersPackageName string
}

func (g *mainGenerator) Generate(tables map[string]tableInfo) error {
	os.Mkdir("bin", 0755)
	importstmpl := template.Must(template.New("mainImports").Parse(`
package main
import({{$rootHandlersPackageName := .rootHandlersPackageName}}{{$rootdbpackagename := .rootdbpackagename}}
{{range .Tables}}"{{$rootHandlersPackageName}}/{{.TableName}}"
{{.TableName}}db "{{$rootdbpackagename}}/{{.TableName}}"
{{end}}
"net/http"
	"github.com/gorilla/mux"
"os"
"database/sql"
	_ "github.com/ziutek/mymysql/godrv"
)
`))
	importstmpl = importstmpl
	routestmpl := template.Must(template.New("mainRoutes").Parse(`
func main(){
	dbUrl := os.Args[1]
    	db, err := sql.Open("mymysql", dbUrl)
	if err != nil {
		panic(err)
	}
    {{range .Tables}}
    {{.TableName}}db.DB = db
    {{end}}
    r := mux.NewRouter()
    g := r.Methods("GET").Subrouter()
    po := r.Methods("POST").Subrouter()
    pu := r.Methods("PUT").Subrouter()
    d := r.Methods("DELETE").Subrouter()
{{range .Tables}}
g.HandleFunc("/{{.TableName}}/", {{.TableName}}.List)
g.HandleFunc("/{{.TableName}}/{id}/", {{.TableName}}.Get)
po.HandleFunc("/{{.TableName}}/", {{.TableName}}.Post)
pu.HandleFunc("/{{.TableName}}/{id}/", {{.TableName}}.Put)
d.HandleFunc("/{{.TableName}}/{id}/", {{.TableName}}.Delete)
{{end}}

http.ListenAndServe(":8080",r)
}
`))
	routestmpl = routestmpl
	var b bytes.Buffer
	err := importstmpl.Execute(&b, map[string]interface{}{"rootHandlersPackageName": "is-a-dev.com/autoapi/http", "Tables": tables, "rootdbpackagename": "is-a-dev.com/autoapi/db"})
	if err != nil {
		fmt.Println(err)
	}
	var final bytes.Buffer
	io.Copy(&final, &b)
	b = bytes.Buffer{}
	routestmpl.Execute(&b, map[string]interface{}{"Verbs": []string{"List", "Get", "Post", "Put", "Delete"}, "Tables": tables})
	io.Copy(&final, &b)
	f, err := os.Create("bin/main.go")
	if err != nil {
		return err
	}
	formatted, err := format.Source(final.Bytes())
	if err != nil {
		return err
	}
	io.Copy(f, bytes.NewBuffer(formatted))
	return nil
}
