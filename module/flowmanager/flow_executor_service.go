package flowmanager

import "github.com/yrn-go/yrn/pkg/yctx"

type (
	JSONSchemaValidator interface {
	}
	FlowExecutor struct {
	}
)

func (f *FlowExecutor) Execute(ctx *yctx.Context, flowId string, data []byte) (err error) {
	return
}
