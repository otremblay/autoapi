package lib

import "fmt"

type constrainter interface {
	CheckCode() string
}

type pkconstraint struct {
	columnName string
}

func (pk *pkconstraint) CheckCode() string {
	return fmt.Sprintf(`if row.%s == %s {return lib.Error("Preconditions failed, %s must be set.")}`, pk.columnName, pk.columnName)
}

type fkconstraint struct {
	columnName string
	nullable   bool
}

type unqconstraint struct {
	columnName string
}

type chkconstraint struct {
	columnName string
	expr       string
}
