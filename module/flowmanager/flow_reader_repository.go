package flowmanager

import "github.com/yrn-go/yrn/pkg/yctx"

type FlowReaderRepository interface {
	GetById(ctx *yctx.Context, id string) (item *Flow, err error)
	GetAll(ctx *yctx.Context, pagination *Pagination) (items []Flow, err error)
	Count(ctx *yctx.Context) (total int, err error)
}
