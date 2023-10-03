//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/wo0lien/cosmoBot/internal/config"
)

func main() {
	f, err := os.Create("json_schema.json")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	schema := jsonschema.Reflect(&config.ConfigStruct{})

	b, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		panic(err)
	}

	_, err = f.Write(b)

	if err != nil {
		panic(err)
	}
}
