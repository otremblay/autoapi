package lib

//Generator is the interface for generating "stuff" out of database table information.
type Generator interface {
	Generate(tables map[string]tableInfo) error
}
