package lib

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func Error(msg string) error {
	return errors.New(msg)
}

func getTableInfo(db *sql.DB, dbName string) (map[string]tableInfo, error) {

	more_rows, err := db.Query("select table_name, column_name, data_type, column_key, is_nullable, extra, column_type from information_schema.columns where table_schema = ?", dbName)

	if err != nil {
		return nil, err
	}
	tables := map[string]tableInfo{}
	for more_rows.Next() {
		var tn, cn, ct, ck, nullable, extra, cdt string
		err := more_rows.Scan(&tn, &cn, &ct, &ck, &nullable, &extra, &cdt)
		if err != nil {
			return nil, err
		}
		var table tableInfo
		var ok bool
		if table, ok = tables[tn]; !ok {
			table = tableInfo{
				TableName:    tn,
				TableColumns: map[string]tableColumn{},
				ColOrder:     []tableColumn{},
				Constraints:  []string{},
			}
		}

		if strings.Contains(cdt, "unsigned") {
			ct = "u" + ct
		}
		col := tableColumn{ColumnName: cn, ColumnType: ct}

		col.Primary = ck == "PRI"
		if nullable == "NO" && extra != "auto_increment" && ct != "timestamp" {
			table.Constraints = append(table.Constraints, fmt.Sprintf(`if row.%s == %s {return libautoapi.Error("Preconditions failed, %s must be set.")}`, col.CapitalizedColumnName(), col.ColumnNullValue(), col.CapitalizedColumnName()))
		}
		table.TableColumns[cn] = col
		table.ColOrder = append(table.ColOrder, col)
		tables[tn] = table
	}
	return tables, nil
}

func Generate(db *sql.DB, dbName string) error {
	tables, err := getTableInfo(db, dbName)
	if err != nil {
		return err
	}
	err = (&dbCodeGenerator{}).Generate(tables)
	if err != nil {
		return err
	}
	err = (&httpCodeGenerator{}).Generate(tables)
	if err != nil {
		return err
	}
	err = (&mainGenerator{}).Generate(tables)
	if err != nil {
		return err
	}
	err = (&checksumGenerator{}).Generate(tables)
	return nil
}

func codeChecksum(tables map[string]tableInfo) []byte {
	checksum := md5.New()
	keys := []string{}
	for k, _ := range tables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		table := tables[k]
		checksum.Write([]byte(table.TableName))
		tckeys := []string{}
		for tck, _ := range table.TableColumns {
			tckeys = append(tckeys, tck)
		}
		sort.Strings(tckeys)
		for _, tck := range tckeys {
			col := table.TableColumns[tck]
			checksum.Write([]byte(fmt.Sprintf("%s%s%b", col.ColumnName, col.ColumnType, col.Primary)))
		}
	}
	return checksum.Sum(nil)
}

func DatabaseChecksum(db *sql.DB, dbName string) ([]byte, error) {
	ti, err := getTableInfo(db, dbName)
	if err != nil {
		return nil, err
	}
	return codeChecksum(ti), nil
}
func GetRootPath() (string, error) {
	pathsplosion, err := os.Getwd()
	if err != nil {
		return "", err
	}
	pathes := strings.Split(pathsplosion, "src/")
	if len(pathes) < 2 {
		return "", errors.New("Bad root dir, outside of a proper GOPATH?")
	}
	return pathes[1], nil
}
