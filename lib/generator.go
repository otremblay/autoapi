package lib

type Generator interface {
	Generate(tables map[string]tableInfo) error
}
