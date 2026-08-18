package cli

import (
	"encoding/json"
	"text/template"

	"sigs.k8s.io/yaml"
)

// Marshaller is a function that marshals input data into specified format.
type Marshaler = func(v any) ([]byte, error)

// FuncMap keeps default func map for rendering operations.
var FuncMap = template.FuncMap{
	"json": toJSON,
	"yaml": toYAML,
}

func marshal(v any, marshaler Marshaler) string {
	var (
		contents []byte
		err      error
	)

	if contents, err = marshaler(v); err != nil {
		contents = []byte(err.Error())
	}
	return string(contents)
}

func toJSON(v any) string {
	return marshal(v, json.Marshal)
}

func toYAML(v any) string {
	return marshal(v, yaml.Marshal)
}
