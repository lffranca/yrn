package postgres

import (
	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/yctx"
)

type FlowWriteRepositoryImpl struct{}

func (f *FlowWriteRepositoryImpl) Save(ctx *yctx.Context, flow *flowmanager.Flow) error {
	//TODO implement me
	panic("implement me")
}

var _ flowmanager.FlowWriteRepository = (*FlowWriteRepositoryImpl)(nil)
