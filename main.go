package main

import (
	"database/sql"
	"os"

	"is-a-dev.com/autoapi/lib"

	_ "github.com/ziutek/mymysql/godrv"
)

func main() {
	dbUrl := os.Args[1]
	dbName := os.Args[2]
	db, err := sql.Open("mymysql", dbUrl)
	if err != nil {
		panic(err)
	}
	err = lib.Generate(db, dbName)
	if err != nil {
		panic(err)
	}
}
