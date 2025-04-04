package pluginhttp

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/plugincore"
	"github.com/yrn-go/yrn/pkg/yctx"
	"io"
	"net/http"
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
		requestData  *HTTPSchema
		requestBody  []byte
		resp         *http.Response
		req          *http.Request
		responseBody []byte
	)

	requestData, err = plugincore.ValidateAndGetRequestBody[HTTPSchema](
		Schema,
		schemaInputs,
		previousPluginResponse,
		responseSharedForAll,
	)
	if err != nil {
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
