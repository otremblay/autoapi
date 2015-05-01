package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"is-a-dev.com/autoapi/lib"

	"github.com/howeyc/gopass"
	_ "github.com/ziutek/mymysql/godrv"
)

func main() {
	dbHost := os.Args[1]
	dbName := os.Args[2]
	dbUname := os.Args[3]

	fmt.Print("Password:")
	pass := strings.TrimSpace(string(gopass.GetPasswdMasked()))
	db, err := sql.Open("mymysql", fmt.Sprintf("tcp:%s:3306*%s/%s/%s", dbHost, dbName, dbUname, pass))
	if err != nil {
		panic(err)
	}
	err = lib.Generate(db, dbName)
	if err != nil {
		panic(err)
	}
}
