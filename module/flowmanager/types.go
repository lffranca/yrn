package flowmanager

import "github.com/yrn-go/yrn/pkg/yctx"

type (
	FlowReaderRepository interface {
		GetById(ctx *yctx.Context, id string) (item *Flow, err error)
		GetAll(ctx *yctx.Context, pagination *Pagination) (items []Flow, err error)
		Count(ctx *yctx.Context) (total int, err error)
	}

	Flow struct {
		Id          string       `json:"id"`
		Name        string       `json:"name"`
		Description string       `json:"description"`
		Tenant      string       `json:"tenant"`
		Plugins     []FlowPlugin `json:"plugins"`
		Version     int          `json:"version"`
	}

	FlowPlugin struct {
		Id                          string `json:"id"`
		Slug                        string `json:"slug"`
		Name                        string `json:"name"`
		Description                 string `json:"description"`
		Version                     int    `json:"version"`
		SchemaInput                 string `json:"schema_input"`
		ContinueEvenWithError       bool   `json:"continue_even_with_error"`
		ShareResponseWithAllPlugins bool   `json:"share_response_with_all_plugins"`
	}
)
