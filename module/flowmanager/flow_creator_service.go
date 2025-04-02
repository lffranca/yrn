package flowmanager

import (
	"github.com/yrn-go/yrn/pkg/yctx"
)

type (
	FlowWriteRepository interface {
		Save(ctx *yctx.Context, flow *Flow) error
	}
)

type FlowCreator struct {
	flowWriteRepository FlowWriteRepository
}

func (f *FlowCreator) CreateFlow(ctx *yctx.Context, flow *Flow) error {
	return f.flowWriteRepository.Save(ctx, flow)
}
