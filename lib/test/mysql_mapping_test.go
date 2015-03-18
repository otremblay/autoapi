package test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"is-a-dev.com/libautoapi/db/every_type"
)

/*
                             TINYINT  -->  int8
                    UNSIGNED TINYINT  -->  uint8
                            SMALLINT  -->  int16
                   UNSIGNED SMALLINT  -->  uint16
                      MEDIUMINT, INT  -->  int32
    UNSIGNED MEDIUMINT, UNSIGNED INT  -->  uint32
                              BIGINT  -->  int64
                     UNSIGNED BIGINT  -->  uint64
                               FLOAT  -->  float32
                              DOUBLE  -->  float64
                             DECIMAL  -->  float64
                 TIMESTAMP, DATETIME  -->  time.Time
                                DATE  -->  mysql.Date
                                TIME  -->  time.Duration
                                YEAR  -->  int16
    CHAR, VARCHAR, BINARY, VARBINARY  -->  []byte
 TEXT, TINYTEXT, MEDIUMTEXT, LONGTEX  -->  []byte
BLOB, TINYBLOB, MEDIUMBLOB, LONGBLOB  -->  []byte
                                 BIT  -->  []byte
                           SET, ENUM  -->  []byte
*/
func nukeit(f reflect.StructField, _ bool) reflect.StructField { return f }
func TestMysqlMapping(t *testing.T) {
	x := every_type.EveryType{}

	tx := reflect.TypeOf(x)

	assert.Equal(t, "int8", nukeit(tx.FieldByName("ATinyint")).Type.String())
	assert.Equal(t, "uint8", nukeit(tx.FieldByName("AnUnsignedTinyint")).Type.String())
	assert.Equal(t, "int16", nukeit(tx.FieldByName("ASmallint")).Type.String())
	assert.Equal(t, "uint16", nukeit(tx.FieldByName("AnUnsignedSmallint")).Type.String())
	assert.Equal(t, "int32", nukeit(tx.FieldByName("AMediumint")).Type.String())
	assert.Equal(t, "int32", nukeit(tx.FieldByName("AnInt")).Type.String())
	assert.Equal(t, "uint32", nukeit(tx.FieldByName("AnUnsignedMediumint")).Type.String())
	assert.Equal(t, "uint32", nukeit(tx.FieldByName("AnUnsignedInt")).Type.String())
	assert.Equal(t, "int64", nukeit(tx.FieldByName("ABigint")).Type.String())
	assert.Equal(t, "uint64", nukeit(tx.FieldByName("AnUnsignedBigint")).Type.String())
	assert.Equal(t, "float32", nukeit(tx.FieldByName("AFloat")).Type.String())
	assert.Equal(t, "float64", nukeit(tx.FieldByName("ADouble")).Type.String())
	assert.Equal(t, "float64", nukeit(tx.FieldByName("ADecimal")).Type.String())
	assert.Equal(t, "time.Time", nukeit(tx.FieldByName("ATimestamp")).Type.String())
	assert.Equal(t, "time.Time", nukeit(tx.FieldByName("ADatetime")).Type.String())
	assert.Equal(t, "mysql.Date", nukeit(tx.FieldByName("ADate")).Type.String())
	assert.Equal(t, "time.Duration", nukeit(tx.FieldByName("ATime")).Type.String())
	assert.Equal(t, "int16", nukeit(tx.FieldByName("AYear")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("AChar")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("AVarchar")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("ABinary")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("AVarbinary")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("AText")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("ATinytext")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("AMediumtext")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("ALongtext")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("ABlob")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("ATinyblob")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("AMediumblob")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("ALongblob")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("ABit")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("ASet")).Type.String())
	assert.Equal(t, "[]uint8", nukeit(tx.FieldByName("AnEnum")).Type.String())
}

/*

*/
