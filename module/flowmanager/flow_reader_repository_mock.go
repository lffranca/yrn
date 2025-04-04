package flowmanager

import (
	"github.com/stretchr/testify/mock"
	"github.com/yrn-go/yrn/pkg/yctx"
)

var (
	_ FlowReaderRepository = (*FlowReaderRepositoryMock)(nil)
)

type FlowReaderRepositoryMock struct {
	mock.Mock
}

func (m *FlowReaderRepositoryMock) GetById(ctx *yctx.Context, id string) (item *Flow, err error) {
	returns := m.MethodCalled("GetById", ctx, id)

	if index := 0; len(returns) > index {
		if itemIndex := returns.Get(index); itemIndex != nil {
			item = itemIndex.(*Flow)
		}
	}

	if index := 1; len(returns) > index {
		err = returns.Error(index)
	}

	return
}

func (m *FlowReaderRepositoryMock) GetAll(ctx *yctx.Context, pagination *Pagination) (items []Flow, err error) {
	returns := m.MethodCalled("GetAll", ctx, pagination)

	if index := 0; len(returns) > index {
		if itemIndex := returns.Get(index); itemIndex != nil {
			items = itemIndex.([]Flow)
		}
	}

	if index := 1; len(returns) > index {
		err = returns.Error(index)
	}

	return
}

func (m *FlowReaderRepositoryMock) Count(ctx *yctx.Context) (total int, err error) {
	returns := m.MethodCalled("Count", ctx)

	if index := 0; len(returns) > index {
		if itemIndex := returns.Get(index); itemIndex != nil {
			total = itemIndex.(int)
		}
	}

	if index := 1; len(returns) > index {
		err = returns.Error(index)
	}

	return
}
