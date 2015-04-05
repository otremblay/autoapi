package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

	definitions := msi{}

	consumes := []string{"application/json"}
	produces := consumes

	for _, t := range tables {
		props := msi{}
		def := msi{"properties": props}
		for _, c := range t.ColOrder {

			prop := msi{}
			ct := c.MappedColumnType()
			if strings.HasPrefix(ct, "int") {
				prop["type"] = "integer"
				prop["format"] = c.MappedColumnType()
			}
			if strings.HasPrefix(ct, "bool") {
				prop["type"] = "boolean"
			}
			props[c.LowercaseColumnName()] = prop
		}

		definitions[t.NormalizedTableName()] = def
		request := msi{}
		request["get"] = msi{
			"produces": produces,
			"responses": msi{
				"200": msi{
					"schema": msi{
						"items": msi{"$ref": fmt.Sprintf("#/definitions/%s", t.NormalizedTableName())},
					},
					"type": "array",
				},
			},
		}

		request["post"] = msi{
			"consumes":  consumes,
			"produces":  produces,
			"responses": msi{"405": msi{"description": "Invalid input"}},
		}

		paths[fmt.Sprintf("/%s", t.CamelCaseTableName())] = request
		if len(t.PrimaryColumns()) > 0 {
			request = msi{}
			request["get"] = msi{
				"produces":  produces,
				"responses": msi{"404": msi{"description": "Not Found"}},
				"parameters": []msi{msi{
					"name":     "id",
					"in":       "path",
					"format":   t.PrimaryColumns()[0].SwaggerFormat(),
					"type":     t.PrimaryColumns()[0].SwaggerColumnType(),
					"required": "true"},
				},
			}

			request["put"] = msi{
				"consumes": consumes,
				"produces": produces,
				"responses": msi{
					"404": msi{"description": "Not Found"},
					"405": msi{"description": "Invalid input"},
				},
				"parameters": []msi{msi{
					"name":     "id",
					"in":       "path",
					"format":   t.PrimaryColumns()[0].SwaggerFormat(),
					"type":     t.PrimaryColumns()[0].SwaggerColumnType(),
					"required": "true"},
				},
			}

			request["delete"] = msi{
				"consumes": consumes,
				"produces": produces,
				"responses": msi{
					"404": msi{"description": "Not Found"},
				},
				"parameters": []msi{msi{
					"name":     "id",
					"in":       "path",
					"format":   t.PrimaryColumns()[0].SwaggerFormat(),
					"type":     t.PrimaryColumns()[0].SwaggerColumnType(),
					"required": "true"},
				},
			}
			paths[fmt.Sprintf("/%s/{%s}", t.CamelCaseTableName(), t.PrimaryColumns()[0].LowercaseColumnName())] = request
		}
	}
	swaggermap["paths"] = paths
	swaggermap["definitions"] = definitions
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
