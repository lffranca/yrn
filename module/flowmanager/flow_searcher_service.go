package flowmanager

import "github.com/yrn-go/yrn/pkg/yctx"

type (
	Pagination struct {
		Page int `json:"page"`
		Size int `json:"size"`
	}

	Page[T any] struct {
		Items      []T `json:"items"`
		TotalItems int `json:"total_items"`
		TotalPages int `json:"total_pages"`
	}

	GetAllResponseFlow struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Version     int    `json:"version"`
	}

	FlowReaderRepository interface {
		GetById(ctx *yctx.Context, id string) (item *Flow, err error)
		GetAll(ctx *yctx.Context, pagination *Pagination) (items []Flow, err error)
		Count(ctx *yctx.Context) (total int, err error)
	}

	FlowSearcher struct {
		flowReaderRepository FlowReaderRepository
	}
)

func (f *FlowSearcher) GetById(ctx *yctx.Context, id string) (*Flow, error) {
	return f.flowReaderRepository.GetById(ctx, id)
}

func (f *FlowSearcher) GetAll(ctx *yctx.Context, pagination *Pagination) (items *Page[GetAllResponseFlow], err error) {
	var (
		flows             []Flow
		total, totalPages int
	)

	flows, err = f.flowReaderRepository.GetAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	total, err = f.flowReaderRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	totalPages = total / pagination.Size

	return &Page[GetAllResponseFlow]{
		Items:      mapperFlowsToGetAllResponses(flows),
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func mapperFlowsToGetAllResponses(flows []Flow) (response []GetAllResponseFlow) {
	for _, flow := range flows {
		response = append(response, *mapperFlowToGetAllResponse(&flow))
	}

	return
}

func mapperFlowToGetAllResponse(flow *Flow) *GetAllResponseFlow {
	return &GetAllResponseFlow{
		Id:          flow.Id,
		Name:        flow.Name,
		Description: flow.Description,
		Version:     flow.Version,
	}
}
