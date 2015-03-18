package lib

type jsonld struct {
	DataTypes  map[string]interface{} `json:datatypes`
	Properties map[string]interface{} `json:properties`
	Types      map[string]interface{} `json:types`
	Valid      string                 `json:valid`
}
