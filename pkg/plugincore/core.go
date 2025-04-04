package plugincore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"log/slog"
	"text/template"
)

func ValidateAndGetRequestBody[T any](schema []byte, schemaInputs string, previousPluginResponse any, responseSharedForAll map[string]any) (requestBody *T, err error) {
	var (
		tmpl           *template.Template
		templateResult = &bytes.Buffer{}
		requestData    T
	)

	tmpl, err = template.
		New("plugin_slug").
		Parse(schemaInputs)
	if err != nil {
		return
	}

	templateData := map[string]any{
		"data":         previousPluginResponse,
		"sharedForAll": responseSharedForAll,
	}

	if err = tmpl.Execute(templateResult, templateData); err != nil {
		return
	}

	if err = validate(schema, templateResult.Bytes()); err != nil {
		return
	}

	if err = json.Unmarshal(templateResult.Bytes(), &requestData); err != nil {
		return
	}

	return &requestData, nil
}

func validate(schema []byte, schemaInputs []byte) (err error) {
	var (
		schemaLoader       = gojsonschema.NewBytesLoader(schema)
		schemaInputsLoader = gojsonschema.NewBytesLoader(schemaInputs)
		result             *gojsonschema.Result
	)

	result, err = gojsonschema.Validate(schemaLoader, schemaInputsLoader)
	if err != nil {
		return
	}

	if !result.Valid() {
		var errs []any

		for _, desc := range result.Errors() {
			errs = append(errs, slog.Any(desc.Field(), desc.Description()))
		}

		slog.Error("error validate", errs...)

		return fmt.Errorf("validation failed: %v", errs)
	}

	return
}
