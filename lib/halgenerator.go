package lib

type halgenerator struct{}

/*
Generates different handlers for HTTP definition of a given collection.
Should generate an OPTIONS verb handler, that returns the available verbs
for a given endpoint.

The generated document should match the database, including required and optional
fields, default values, stuff of the like.
*/
func (h *halgenerator) Generate(tables map[string]tableInfo) error {
	return nil
}
