package lib

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"os"
	"strings"

	"golang.org/x/tools/imports"
)

type mainGenerator struct {
	rootDbPackageName       string
	rootHandlersPackageName string
	verbs                   string
}

func (g *mainGenerator) Generate(tables map[string]tableInfo) error {
	os.Mkdir("bin", 0755)
	importstmpl := template.Must(template.New("mainImports").Parse(`
package main
import({{$rootHandlersPackageName := .rootHandlersPackageName}}{{$rootdbpackagename := .rootdbpackagename}}
{{range .Tables}}"{{$rootHandlersPackageName}}/{{.TableName}}"
{{.TableName}}db "{{$rootdbpackagename}}/{{.TableName}}"
{{end}}
"{{$rootdbpackagename}}"
"{{.checksum}}"
"net/http"
	"github.com/gorilla/mux"
"os"
"database/sql"
"github.com/howeyc/gopass"
	_ "github.com/ziutek/mymysql/godrv"
)
`))
	importstmpl = importstmpl
	routestmpl := template.Must(template.New("mainRoutes").Parse(`
const (
	// Ugly way to check to see if they passed in a password
	// chance of collision with a GUID is very low
	defaultPassValue = "5e7dc6f6a1a94c39be95b88a47c2458b"
)

var (
	dbPort  string
	dbHost  string
	dbName  string
	dbUname string
	dbPass  string
	listenAddress string
)

func init() {
	flag.StringVar(&dbPort, "P", "3306", "port")
	flag.StringVar(&dbPass, "p", defaultPassValue, "password")
	flag.StringVar(&dbHost, "h", "localhost", "host")
	flag.StringVar(&dbName, "d", "", "database name")
	flag.StringVar(&dbUname, "u", "root", "username")
	flag.StringVar(&listenAddress, "l", ":8080", "listening address")
	flag.Parse()
}

func main(){
	pass := dbPass
	if pass == defaultPassValue {
		fmt.Print("Password:")
		pass = strings.TrimSpace(string(gopass.GetPasswdMasked()))
	}

	if strings.TrimSpace(dbPort) == "" {
		fmt.Println("Missing port")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if strings.TrimSpace(dbHost) == "" {
		fmt.Println("Missing host")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if strings.TrimSpace(dbName) == "" {
		fmt.Println("Missing database name")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if strings.TrimSpace(dbUname) == "" {
		fmt.Println("Missing username")
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	dbConn, err := sql.Open("mymysql", fmt.Sprintf("tcp:%s:3306*%s/%s/%s", dbHost, dbName, dbUname, pass))
	if err != nil {
		panic(err)
	}
    db.MustValidateChecksum(dbConn, dbName)
    {{range .Tables}}
    {{.TableName}}db.DB = dbConn
    {{end}}
    r := mux.NewRouter()
{{$shouldGenerateGet := .shouldGenerateGet}}
{{$shouldGeneratePost := .shouldGeneratePost}}
{{$shouldGeneratePut := .shouldGeneratePut}}
{{$shouldGenerateDelete := .shouldGenerateDelete}}
{{if $shouldGenerateGet}}
    g := r.Methods("GET").Subrouter()
{{end}}
{{if $shouldGeneratePost}}
    po := r.Methods("POST").Subrouter()
{{end}}
{{if $shouldGeneratePut}}
    pu := r.Methods("PUT").Subrouter()
{{end}}
{{if $shouldGenerateDelete}}
    d := r.Methods("DELETE").Subrouter()
{{end}}
{{range .Tables}}
{{$l := len .PrimaryColumns}}
{{if $shouldGenerateGet}}
g.HandleFunc("/{{.TableName}}", {{.TableName}}.List)
{{end}}
{{if $shouldGeneratePost}}
po.HandleFunc("/{{.TableName}}", {{.TableName}}.Post)
{{end}}
{{if gt $l 0}}
{{if $shouldGenerateGet}}
g.HandleFunc("/{{.TableName}}/{id}", {{.TableName}}.Get)
{{end}}
{{if $shouldGeneratePut}}
pu.HandleFunc("/{{.TableName}}/{id}", {{.TableName}}.Put)
{{end}}
{{if $shouldGenerateDelete}}
d.HandleFunc("/{{.TableName}}/{id}", {{.TableName}}.Delete)
{{end}}
{{end}}


{{end}}

g.HandleFunc("/swagger.json", swaggerresponse)

http.ListenAndServe(listenAddress,r)
}
`))
	routestmpl = routestmpl
	var b bytes.Buffer
	path, err := GetRootPath()
	if err != nil {
		return err
	}
	err = importstmpl.Execute(&b, map[string]interface{}{"rootHandlersPackageName": path + "/http", "Tables": tables, "rootdbpackagename": path + "/" + rootdbpath, "checksum": path + "/db"})
	if err != nil {
		fmt.Println(b.String())
		fmt.Println(err)
	}
	var final bytes.Buffer
	io.Copy(&final, &b)
	b = bytes.Buffer{}
	verbs := strings.Split(g.verbs, ",")
	routestmpl.Execute(&b, map[string]interface{}{
		"Verbs":                []string{"List", "Get", "Post", "Put", "Delete"},
		"Tables":               tables,
		"shouldGenerateGet":    stringInSlice("get", verbs),
		"shouldGeneratePost":   stringInSlice("post", verbs),
		"shouldGeneratePut":    stringInSlice("put", verbs),
		"shouldGenerateDelete": stringInSlice("delete", verbs),
	})
	io.Copy(&final, &b)
	f, err := os.Create("bin/main.go")
	if err != nil && !os.IsExist(err) {
		return err
	}

	formatted, err := format.Source(final.Bytes())
	if err != nil {
		io.Copy(f, bytes.NewBuffer(final.Bytes()))
		return err
	}
	formatted, err = imports.Process(f.Name(), formatted, nil)
	if err != nil {
		io.Copy(f, bytes.NewBuffer(formatted))
		return err
	}
	io.Copy(f, bytes.NewBuffer(formatted))
	return nil
}
