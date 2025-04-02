package mongodb

import (
	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/yctx"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CollectionFlowName = "flow"
)

var (
	_ flowmanager.FlowReaderRepository = (*FlowRepository)(nil)
	_ flowmanager.FlowWriteRepository  = (*FlowRepository)(nil)
)

type (
	FlowRepository struct {
	}
)

func (f *FlowRepository) Save(ctx *yctx.Context, flow *flowmanager.Flow) (err error) {
	var (
		collection *mongo.Collection
	)

	collection, err = GetCollection(ctx, CollectionFlowName)
	if err != nil {
		return
	}

	_, err = collection.InsertOne(ctx.Context(), flow)
	if err != nil {
		return
	}

	return
}

func (f *FlowRepository) GetById(ctx *yctx.Context, id string) (item *flowmanager.Flow, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *FlowRepository) GetAll(ctx *yctx.Context, pagination *flowmanager.Pagination) (items []flowmanager.Flow, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *FlowRepository) Count(ctx *yctx.Context) (total int, err error) {
	//TODO implement me
	panic("implement me")
}
