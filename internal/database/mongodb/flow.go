package mongodb

import (
	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/yctx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
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
	var (
		collection *mongo.Collection
	)

	collection, err = GetCollection(ctx, CollectionFlowName)
	if err != nil {
		return
	}

	filter := bson.M{"_id": id}
	result := collection.FindOne(ctx.Context(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Err()
	}

	item = new(flowmanager.Flow)
	err = result.Decode(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (f *FlowRepository) GetAll(ctx *yctx.Context, pagination *flowmanager.Pagination) (items []flowmanager.Flow, err error) {
	var (
		collection *mongo.Collection
	)

	collection, err = GetCollection(ctx, CollectionFlowName)
	if err != nil {
		return
	}

	options := mongoOptions.Find().
		SetSkip(int64(pagination.Page * pagination.Size)).
		SetLimit(int64(pagination.Size))

	cursor, err := collection.Find(ctx.Context(), bson.M{}, options)
	if err != nil {
		return
	}
	defer cursor.Close(ctx.Context())

	err = cursor.All(ctx.Context(), &items)
	if err != nil {
		return
	}

	return items, nil
}

func (f *FlowRepository) Count(ctx *yctx.Context) (total int, err error) {
	var (
		collection *mongo.Collection
		count      int64
	)

	collection, err = GetCollection(ctx, CollectionFlowName)
	if err != nil {
		return
	}

	count, err = collection.CountDocuments(ctx.Context(), bson.M{})
	if err != nil {
		return
	}

	return int(count), nil
}
