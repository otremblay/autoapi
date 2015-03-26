package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

type swaggerGenerator struct{}

type msi map[string]interface{}

func (sg *swaggerGenerator) Generate(tables map[string]tableInfo) error {
	dbname := os.Args[2]
	swaggermap := msi{
		"swagger": "2.0",
		"info": msi{
			"title":       fmt.Sprintf("%s Autoapi", dbname),
			"description": fmt.Sprintf("Automatic api for the %s database", dbname),
			"version":     "1.0.0", //TODO: Read from old version if available.
		},
		"paths": msi{},
	}
	paths := msi{}
	for _, t := range tables {
		request := msi{}
		request["get"] = msi{}
		request["post"] = msi{}

		paths[fmt.Sprintf("/%s", t.CamelCaseTableName())] = request
		if len(t.PrimaryColumns()) > 0 {
			request = msi{}
			request["get"] = msi{}

			request["put"] = msi{}
			request["delete"] = msi{}
			paths[fmt.Sprintf("/%s/{%s}", t.CamelCaseTableName(), t.PrimaryColumns()[0].LowercaseColumnName())] = request
		}
	}
	swaggermap["paths"] = paths

	f, err := os.Create("bin/swagger.json")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	err = enc.Encode(swaggermap)
	if err != nil {
		return err
	}
	return nil
}
