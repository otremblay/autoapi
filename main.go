package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/ziutek/mymysql/godrv"
)

func main() {
	dbUrl := os.Args[1]
	dbName := os.Args[2]
	db, err := sql.Open("mymysql", dbUrl)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select table_name from information_schema.tables where table_schema = ?", dbName)
	if err != nil {
		panic(err)
	}
	tables := map[string]tableInfo{}
	for rows.Next() {
		var tn string
		rows.Scan(&tn)
		tables[tn] = tableInfo{
			TableName:    tn,
			TableColumns: map[string]tableColumn{},
			ColOrder:     []tableColumn{},
			Constraints:  []string{},
		}
	}

	more_rows, err := db.Query("select table_name, column_name, data_type, column_key, is_nullable, extra from information_schema.columns where table_schema = ?", dbName)

	if err != nil {
		panic(err)
	}

	for more_rows.Next() {
		var tn, cn, ct, ck, nullable, extra string
		err := more_rows.Scan(&tn, &cn, &ct, &ck, &nullable, &extra)
		if err != nil {
			panic(err)
		}
		table := tables[tn]
		col := tableColumn{ColumnName: cn, ColumnType: ct}

		col.Primary = ck == "PRI"
		if nullable == "NO" && extra != "auto_increment" {
			table.Constraints = append(table.Constraints, fmt.Sprintf(`if row.%s == %s {return libautoapi.Error("Preconditions failed, %s must be set.")}`, col.CapitalizedColumnName(), col.ColumnNullValue(), col.CapitalizedColumnName()))
		}
		table.TableColumns[cn] = col
		table.ColOrder = append(table.ColOrder, col)
		tables[tn] = table
	}

	err = (&dbCodeGenerator{}).Generate(tables)
	if err != nil {
		panic(err)
	}
	err = (&httpCodeGenerator{}).Generate(tables)
	if err != nil {
		panic(err)
	}
	err = (&mainGenerator{}).Generate(tables)
	if err != nil {
		panic(err)
	}
}
