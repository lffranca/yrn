package main

import (
	"encoding/json"
	"github.com/qri-io/jsonschema"
	"github.com/yrn-go/yrn/pkg/yconnector"
	"golang.org/x/exp/slog"
	"log"
)

func main() {
	var schema jsonschema.Schema

	if err := json.Unmarshal([]byte(schemaData), &schema); err != nil {
		log.Panicf("unmarshal schema: %v\n", err.Error())
	}

	appRun := yconnector.NewApp(&schema, map[string]string{})
	if err := appRun(); err != nil {
		slog.Error("server error: ", slog.Any("error", err))
	}
}

const schemaData = `{
	"$id": "https://qri.io/schema/",
	"$comment" : "sample comment",
	"title": "Person",
	"type": "object",
	"properties": {
		"firstName": {
			"type": "string"
		},
		"lastName": {
			"type": "string"
		},
		"age": {
			"description": "Age in years",
			"type": "integer",
			"minimum": 0
		},
		"friends": {
			"type" : "array",
			"items" : { "title" : "REFERENCE", "$ref" : "#" }
		}
	},
	"required": ["firstName", "lastName"]
}`
