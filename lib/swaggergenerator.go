package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

type swaggerGenerator struct{ dbname string }

type msi map[string]interface{}

func (sg *swaggerGenerator) Generate(tables map[string]tableInfo) error {
	dbname := sg.dbname

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

		refdef := fmt.Sprintf("#/definitions/%s", t.NormalizedTableName())
		idparam := msi{}
		if len(t.PrimaryColumns()) > 0 {
			idparam = msi{
				"name":     t.PrimaryColumns()[0].LowercaseColumnName(),
				"in":       "path",
				"format":   t.PrimaryColumns()[0].SwaggerFormat(),
				"type":     t.PrimaryColumns()[0].SwaggerColumnType(),
				"required": "true"}
		}
		bodydef := msi{
			"name":     "body",
			"in":       "body",
			"required": "true",
			"schema":   msi{"$ref": refdef},
		}
		props := msi{}
		def := msi{"properties": props}
		for _, c := range t.ColOrder {

			prop := msi{}
			prop["type"] = c.SwaggerColumnType()
			prop["format"] = c.SwaggerFormat()
			props[c.LowercaseColumnName()] = prop
		}

		definitions[t.TableName] = def
		request := msi{}
		request["get"] = msi{
			"produces": produces,
			"responses": msi{
				"200": msi{
					"schema": msi{
						"$ref": refdef,
					},
					"type": "array",
				},
				"400": msi{"description": "Empty collection"},
			},
		}

		request["post"] = msi{
			"consumes":   consumes,
			"produces":   produces,
			"responses":  msi{"405": msi{"description": "Invalid input"}},
			"parameters": []msi{bodydef},
		}

		paths[fmt.Sprintf("/%s", t.CamelCaseTableName())] = request
		if len(t.PrimaryColumns()) > 0 {
			request = msi{}
			request["get"] = msi{
				"produces": produces,
				"responses": msi{
					"404": msi{"description": "Not Found"},
					"200": msi{"$ref": refdef},
				},
				"parameters": []msi{idparam},
			}

			request["put"] = msi{
				"consumes": consumes,
				"produces": produces,
				"responses": msi{
					"404": msi{"description": "Not Found"},
					"405": msi{"description": "Invalid input"},
				},
				"parameters": []msi{
					idparam,
					bodydef,
				},
			}

			request["delete"] = msi{
				"consumes": consumes,
				"produces": produces,
				"responses": msi{
					"404": msi{"description": "Not Found"},
				},
				"parameters": []msi{
					idparam,
				},
			}
			paths[fmt.Sprintf("/%s/{%s}", t.CamelCaseTableName(), t.PrimaryColumns()[0].LowercaseColumnName())] = request
		}
	}
	swaggermap["paths"] = paths
	swaggermap["definitions"] = definitions
	os.MkdirAll("bin", 0755)
	f, err := os.Create("bin/swagger.json.go")
	if err != nil {
		return err
	}

	var x bytes.Buffer
	enc := json.NewEncoder(&x)
	err = enc.Encode(swaggermap)
	if err != nil {
		return err
	}

	var tmpl, _ = template.New("swagga").Parse(swag)
	b, _ := json.Marshal(swaggermap)
	tmpl.Execute(f, string(b))

	return nil
}

var swag = `package main

import (
"net/http"
"fmt"
)
var js = ` + "`{{.}}`" + `

func swaggerresponse(res http.ResponseWriter, req *http.Request){
fmt.Fprint(res, js)
}`
