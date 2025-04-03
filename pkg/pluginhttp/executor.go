package pluginhttp

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/yctx"
	"io"
	"log/slog"
	"net/http"
	"text/template"
)

const (
	SlugHttp = "http"
)

var (
	_ flowmanager.PluginExecutor = (*Executor)(nil)
	//go:embed schema.json
	Schema []byte
)

type (
	Executor struct{}
)

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Do(ctx *yctx.Context, schemaInputs string, previousPluginResponse any, responseSharedForAll map[string]any) (output any, err error) {
	var (
		tmpl           *template.Template
		templateResult = &bytes.Buffer{}
		requestData    HTTPSchema
		requestBody    []byte
		resp           *http.Response
		req            *http.Request
		responseBody   []byte
	)

	tmpl, err = template.
		New(SlugHttp).
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

	if err = e.validate(templateResult.Bytes()); err != nil {
		return
	}

	if err = json.Unmarshal(templateResult.Bytes(), &requestData); err != nil {
		return
	}

	requestBody, err = json.Marshal(requestData.Request.Body)
	if err != nil {
		return
	}

	req, err = http.NewRequestWithContext(
		ctx.Context(),
		requestData.Request.Method,
		requestData.Request.URL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return
	}

	for key, value := range requestData.Request.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	responseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(responseBody, &output)
	return
}

func (e *Executor) validate(schemaInputs []byte) (err error) {
	var (
		schemaLoader       = gojsonschema.NewBytesLoader(Schema)
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
