package flowmanager

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/yrn-go/yrn/pkg/yctx"
	"log/slog"
	"testing"
)

// Mocks

type PluginExecutorMock struct {
	mock.Mock
}

func (m *PluginExecutorMock) Do(ctx *yctx.Context, schema string, data any, shared map[string]any) (any, error) {
	args := m.Called(ctx, schema, data, shared)
	return args.Get(0), args.Error(1)
}

type PluginManagerMock struct {
	mock.Mock
}

func (m *PluginManagerMock) GetBySlug(ctx *yctx.Context, slug string) (PluginExecutor, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(PluginExecutor), args.Error(1)
}

// Suite

type EventManagerTestSuite struct {
	suite.Suite
	pluginManagerMock *PluginManagerMock
	eventManager      *EventManager
	ctx               *yctx.Context
}

func TestEventManagerTestSuite(t *testing.T) {
	suite.Run(t, new(EventManagerTestSuite))
}

func (s *EventManagerTestSuite) SetupTest() {
	s.ctx = yctx.NewContext(context.Background())
	s.pluginManagerMock = new(PluginManagerMock)
	s.eventManager = NewEventManager(s.pluginManagerMock)
}

func (s *EventManagerTestSuite) TestExecute_ShouldReturnSuccess() {
	const pluginSlug = "plugin-http"

	// Plugin e Executor Mock
	executorMock := new(PluginExecutorMock)
	_ = s.eventManager.Register(FlowPlugin{
		Id:               "test1",
		Slug:             pluginSlug,
		SchemaInput:      `{"mock": true}`,
		NextToBeExecuted: []string{"test2", "test2"},
	})
	_ = s.eventManager.Register(FlowPlugin{
		Id:               "test2",
		Slug:             pluginSlug,
		SchemaInput:      `{"mock": true}`,
		NextToBeExecuted: []string{"test3", "test3", "final"},
	})
	_ = s.eventManager.Register(FlowPlugin{
		Id:               "test3",
		Slug:             pluginSlug,
		SchemaInput:      `{"mock": true}`,
		NextToBeExecuted: []string{"test4", "test4", "test4"},
	})
	_ = s.eventManager.Register(FlowPlugin{
		Id:               "test4",
		Slug:             pluginSlug,
		SchemaInput:      `{"mock": true}`,
		NextToBeExecuted: []string{"final", "final"},
	})
	_ = s.eventManager.Register(FlowPlugin{
		Id:          "final",
		Slug:        pluginSlug,
		SchemaInput: `{"mock": true}`,
	})

	s.pluginManagerMock.
		On("GetBySlug", mock.Anything, pluginSlug).
		Return(executorMock, nil)

	executorMock.
		On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(map[string]any{"success": true}, nil)

	//// Executa com timeout menor que o original
	//go func() {
	//	time.Sleep(2 * time.Second)
	//}()

	finalResponse, err := s.eventManager.Execute(s.ctx, "test1", map[string]any{"input": "value"})

	slog.Info("finalResponse", slog.Any("finalResponse", finalResponse), slog.Any("err", err))

	//s.NoError(err)
	//s.NotNil(finalResponse)
	//s.Equal(map[string]any{"success": true}, finalResponse)
}
