package plugingdriveauth

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/plugincore"
	"github.com/yrn-go/yrn/pkg/yctx"
	"io"
	"net/http"
	"net/url"
)

const (
	SlugGDriveAuth = "gdrive-auth"
)

var (
	_ flowmanager.PluginExecutor = (*Executor)(nil)
	//go:embed schema.json
	Schema []byte
)

type (
	Executor struct{}

	AuthSchema struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		RedirectURI  string `json:"redirect_uri"`
	}
)

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Do(ctx *yctx.Context, schemaInputs string, previousPluginResponse any, responseSharedForAll map[string]any) (output any, err error) {
	var schema *AuthSchema
	schema, err = plugincore.ValidateAndGetRequestBody[AuthSchema](Schema, schemaInputs, previousPluginResponse, responseSharedForAll)
	if err != nil {
		return
	}

	form := url.Values{}
	form.Set("client_id", schema.ClientID)
	form.Set("client_secret", schema.ClientSecret)
	form.Set("code", schema.Code)
	form.Set("redirect_uri", schema.RedirectURI)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx.Context(), "POST", "https://oauth2.googleapis.com/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New("authentication failed")
		return
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &output)
	return
}
