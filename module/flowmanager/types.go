package flowmanager

import "github.com/yrn-go/yrn/pkg/yctx"

type (
	FlowReaderRepository interface {
		GetById(ctx *yctx.Context, id string) (item *Flow, err error)
		GetAll(ctx *yctx.Context, pagination *Pagination) (items []Flow, err error)
		Count(ctx *yctx.Context) (total int, err error)
	}

	Flow struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Tenant      string   `json:"tenant"`
		Plugins     []Plugin `json:"plugins"`
		Version     int      `json:"version"`
	}

	Plugin struct {
		Id                          string   `json:"id"`
		Name                        string   `json:"name"`
		Slug                        string   `json:"slug"`
		Description                 string   `json:"description"`
		Schema                      string   `json:"schema"`
		DiagramData                 string   `json:"diagram_data"`
		FlowId                      string   `json:"flow_id"`
		Tenant                      string   `json:"tenant"`
		NextSteps                   []string `json:"next_steps"`
		ContinueEvenWithError       bool     `json:"continue_even_with_error"`
		ShareResponseWithAllPlugins bool     `json:"share_response_with_all_plugins"`
	}
)
