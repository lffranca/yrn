package flowmanager

import (
	"github.com/stretchr/testify/mock"
	"github.com/yrn-go/yrn/pkg/yctx"
)

// PluginStatusRepositoryMock implementa PluginStatusRepository
type PluginStatusRepositoryMock struct {
	mock.Mock
}

func (m *PluginStatusRepositoryMock) Save(ctx *yctx.Context, status PluginStatus) error {
	args := m.Called(ctx, status)
	return args.Error(0)
}

func (m *PluginStatusRepositoryMock) GetByPluginID(ctx *yctx.Context, pluginID string) (PluginStatus, error) {
	args := m.Called(ctx, pluginID)
	return args.Get(0).(PluginStatus), args.Error(1)
}

func (m *PluginStatusRepositoryMock) GetAll(ctx *yctx.Context) ([]PluginStatus, error) {
	args := m.Called(ctx)
	return args.Get(0).([]PluginStatus), args.Error(1)
}

// PluginExecutorMock implementa PluginExecutor
type PluginExecutorMock struct {
	mock.Mock
}

func (m *PluginExecutorMock) Do(ctx *yctx.Context, schema string, data any, shared map[string]any) (any, error) {
	args := m.Called(ctx, schema, data, shared)
	return args.Get(0), args.Error(1)
}

// PluginManagerMock implementa PluginManager
type PluginManagerMock struct {
	mock.Mock
}

func (m *PluginManagerMock) GetBySlug(ctx *yctx.Context, slug string) (PluginExecutor, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(PluginExecutor), args.Error(1)
}
