package lib

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

import (
	"go/ast"
	"go/parser"

	"github.com/stretchr/testify/assert"
)
import "go/token"

func TestDeclaredProperChecksumFunction(t *testing.T) {
	//Stage setting
	os.RemoveAll("db")

	x := &checksumGenerator{}
	x.Generate(map[string]tableInfo{"x": tableInfo{TableName: "y"}})

	fset := token.NewFileSet()
	my_ast, err := parser.ParseDir(fset, "db/", nil, parser.AllErrors)
	if err != nil {
		t.Error(err)
	}
	main_ast := my_ast["db"]

	f := main_ast.Files["db/checksum.go"]
	obj := f.Scope.Lookup("Checksum")

	//Generated Checksum is indeed a nullary string-type function
	checksumfn := (obj.Decl.(*ast.FuncDecl))
	assert.Empty(t, checksumfn.Type.Params.List)
	assert.Equal(t, "string", checksumfn.Type.Results.List[0].Type.(*ast.Ident).Name)

	//Cleanup
	os.RemoveAll("db")
}

//
func TestDBChecksum(t *testing.T) {

	//Let's validate that we're not subject to map random access.
	m := map[string]tableInfo{}
	for i := 0; i < 100; i++ {
		t := tableInfo{}
		t.TableName = fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(100000000)))))
		t.TableColumns = map[string]tableColumn{}
		for j := 0; j < 10; j++ {
			c := tableColumn{}
			c.ColumnName = fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(rand.Intn(100000000)))))

			switch rand.Intn(5) {
			case 0:
				c.ColumnType = "varchar"
			case 1:
				c.ColumnType = "int"
			case 2:
				c.ColumnType = "text"
			case 3:
				c.ColumnType = "datetime"
			case 4:
				c.ColumnType = "bit"
			}
			t.TableColumns[c.ColumnName] = c
			t.ColOrder = append(t.ColOrder, c)
		}
		m[t.TableName] = t
	}

	originalChecksum := codeChecksum(m)
	for i := 0; i < 1000; i++ {
		cs := codeChecksum(m)
		assert.Equal(t, originalChecksum, cs)
	}
}
