package test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/ziutek/mymysql/godrv"
	"is-a-dev.com/libautoapi/db"
)

func TestBadChecksum(t *testing.T) {
	dburl, dbname := "tcp:127.0.0.1:3306*autoapi/autoapi/something", "autoapi"
	fmt.Println(dburl, dbname)
	dbconn, err := sql.Open("mymysql", dburl)
	if err != nil {
		t.Error(err)
	}
	err = db.ValidateChecksum(dbconn, dbname)
	if err != nil {
		t.Error(err)
	}
	dbconn.Exec("ALTER TABLE person add column (xyzzyx varchar(255))")
	err = db.ValidateChecksum(dbconn, dbname)
	assert.Error(t, err)
	assert.Panics(t, func() { db.MustValidateChecksum(dbconn, dbname) })
	dbconn.Exec("ALTER TABLE person drop column xyzzyx")
}
