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
	"$id": "https://qri.io/schema/http-plugin",
	"$comment": "Schema para um plugin de chamadas HTTP",
	"title": "HTTP Plugin",
	"type": "object",
	"properties": {
		"request": {
			"type": "object",
			"description": "Configuração da requisição HTTP",
			"properties": {
				"method": {
					"type": "string",
					"enum": ["GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"],
					"description": "Método HTTP da requisição"
				},
				"url": {
					"type": "string",
					"format": "uri",
					"description": "URL para a requisição"
				},
				"headers": {
					"type": "object",
					"additionalProperties": {
						"type": "string"
					},
					"description": "Cabeçalhos da requisição"
				},
				"queryParams": {
					"type": "object",
					"additionalProperties": {
						"type": "string"
					},
					"description": "Parâmetros de consulta da URL"
				},
				"body": {
					"type": ["string", "object", "array", "null"],
					"description": "Corpo da requisição (JSON, texto ou nulo)"
				},
				"timeout": {
					"type": "integer",
					"minimum": 0,
					"description": "Tempo limite para a requisição em milissegundos"
				}
			},
			"required": ["method", "url"]
		},
		"response": {
			"type": "object",
			"description": "Configuração da resposta esperada",
			"properties": {
				"statusCode": {
					"type": "integer",
					"minimum": 100,
					"maximum": 599,
					"description": "Código de status esperado"
				},
				"headers": {
					"type": "object",
					"additionalProperties": {
						"type": "string"
					},
					"description": "Cabeçalhos da resposta"
				},
				"body": {
					"type": ["string", "object", "array", "null"],
					"description": "Corpo da resposta esperada"
				}
			}
		},
		"retry": {
			"type": "object",
			"description": "Configurações de tentativas de reexecução",
			"properties": {
				"maxAttempts": {
					"type": "integer",
					"minimum": 0,
					"description": "Número máximo de tentativas"
				},
				"delay": {
					"type": "integer",
					"minimum": 0,
					"description": "Tempo de espera entre tentativas (ms)"
				}
			},
			"required": ["maxAttempts", "delay"]
		}
	},
	"required": ["request"]
}`
