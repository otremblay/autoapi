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

type checksumGenerator struct {
}

func (g *checksumGenerator) Generate(tables map[string]tableInfo) error {
	err := os.Mkdir("db", 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	tmpl := template.Must(template.New("dbchecksum").Parse(`//WARNING.
//THIS HAS BEEN GENERATED AUTOMATICALLY BY AUTOAPI.
//IF THERE WAS A WARRANTY, MODIFYING THIS WOULD VOID IT.

package db

func Checksum() string{
    return "{{.}}"
}

func ValidateChecksum(db *sql.DB, dbName string) error {
     b, err := libautoapi.DatabaseChecksum(db, dbName)	
     if err != nil {
         return err
     }
     if fmt.Sprintf("%x", b) != Checksum() {
        return ErrBadDatabaseChecksum
     }
     return nil
}

func MustValidateChecksum(db *sql.DB, dbName string) {
    if err := ValidateChecksum(db, dbName); err != nil{
       panic(err)
    }
}

var ErrBadDatabaseChecksum = errors.New("The code doesn't match the database's structure.")

`))

	f, err := os.Create("db/checksum.go")
	if err != nil {
		return err
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, fmt.Sprintf("%x", codeChecksum(tables)))
	if err != nil {
		return err
	}
	bf, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Println(b.String())
		fmt.Println(err)
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
	return nil
}
