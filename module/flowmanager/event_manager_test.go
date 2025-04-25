package flowmanager

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/yrn-go/yrn/pkg/yctx"
)

func TestEventManagerTestSuite(t *testing.T) {
	suite.Run(t, new(EventManagerTestSuite))
}

type EventManagerTestSuite struct {
	suite.Suite
	pluginManagerMock    *PluginManagerMock
	statusRepositoryMock *PluginStatusRepositoryMock
	eventManager         *EventManager
	ctx                  *yctx.Context
}

func (s *EventManagerTestSuite) SetupTest() {
	s.ctx = yctx.NewContext(context.Background())
	s.pluginManagerMock = new(PluginManagerMock)
	s.statusRepositoryMock = new(PluginStatusRepositoryMock)
	s.eventManager = NewEventManager(s.pluginManagerMock, s.statusRepositoryMock)
}

func (s *EventManagerTestSuite) TestExecute_ShouldSavePluginStatus() {
	const pluginSlug = "plugin-http"
	const pluginID = "test1"

	// Plugin e Executor Mock
	executorMock := new(PluginExecutorMock)
	_ = s.eventManager.Register(FlowPlugin{
		Id:          pluginID,
		Slug:        pluginSlug,
		SchemaInput: `{"mock": true}`,
	})

	s.pluginManagerMock.
		On("GetBySlug", mock.Anything, pluginSlug).
		Return(executorMock, nil)

	executorMock.
		On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(map[string]any{"success": true}, nil)

	// Configura o mock do repositório para esperar duas chamadas de Save
	// Uma para o status inicial e outra para o status final
	s.statusRepositoryMock.
		On("Save", mock.Anything, mock.Anything).
		Return(nil).
		Twice()

	_, err := s.eventManager.Execute(s.ctx, pluginID, map[string]any{"input": "value"})

	s.NoError(err)
	s.statusRepositoryMock.AssertExpectations(s.T())
}

func (s *EventManagerTestSuite) TestExecute_ShouldHandleRepositoryError() {
	const pluginSlug = "plugin-http"
	const pluginID = "test1"

	// Plugin e Executor Mock
	executorMock := new(PluginExecutorMock)
	_ = s.eventManager.Register(FlowPlugin{
		Id:          pluginID,
		Slug:        pluginSlug,
		SchemaInput: `{"mock": true}`,
	})

	s.pluginManagerMock.
		On("GetBySlug", mock.Anything, pluginSlug).
		Return(executorMock, nil)

	executorMock.
		On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(map[string]any{"success": true}, nil)

	// Configura o mock do repositório para retornar erro
	s.statusRepositoryMock.
		On("Save", mock.Anything, mock.Anything).
		Return(fmt.Errorf("repository error")).
		Twice()

	// O erro do repositório não deve impedir a execução do plugin
	_, err := s.eventManager.Execute(s.ctx, pluginID, map[string]any{"input": "value"})

	s.NoError(err)
	s.statusRepositoryMock.AssertExpectations(s.T())
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

	// Configura o mock do repositório para todas as chamadas de Save
	s.statusRepositoryMock.
		On("Save", mock.Anything, mock.Anything).
		Return(nil)

	finalResponse, err := s.eventManager.Execute(s.ctx, "test1", map[string]any{"input": "value"})

	s.NoError(err)
	s.NotNil(finalResponse)
	s.Equal(map[string]any{"success": true}, finalResponse)
}

func (s *EventManagerTestSuite) TestExecute_ShouldCollectMetrics() {
	const pluginSlug = "plugin-http"
	const pluginID = "test1"

	// Plugin e Executor Mock
	executorMock := new(PluginExecutorMock)
	_ = s.eventManager.Register(FlowPlugin{
		Id:          pluginID,
		Slug:        pluginSlug,
		SchemaInput: `{"mock": true}`,
	})

	s.pluginManagerMock.
		On("GetBySlug", mock.Anything, pluginSlug).
		Return(executorMock, nil)

	executorMock.
		On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(map[string]any{"success": true}, nil)

	// Configura o mock do repositório
	s.statusRepositoryMock.
		On("Save", mock.Anything, mock.Anything).
		Return(nil).
		Twice()

	_, err := s.eventManager.Execute(s.ctx, pluginID, map[string]any{"input": "value"})

	s.NoError(err)

	// Verifica se as métricas foram coletadas
	metrics, exists := s.eventManager.metrics[pluginID]
	s.True(exists, "Métricas não foram coletadas para o plugin")

	// Verifica se os tempos foram registrados
	s.NotZero(metrics.StartTime)
	s.NotZero(metrics.EndTime)
	s.NotZero(metrics.ExecutionTime)

	// Verifica se o uso de memória foi registrado
	s.NotZero(metrics.MemoryBefore)
	s.NotZero(metrics.MemoryAfter)

	// Verifica se o uso de CPU foi registrado
	s.NotZero(metrics.CPUBefore)
	s.NotZero(metrics.CPUAfter)
}

func (s *EventManagerTestSuite) TestExecute_ShouldHandlePluginError() {
	const pluginSlug = "plugin-http"
	const pluginID = "test1"

	// Plugin e Executor Mock
	executorMock := new(PluginExecutorMock)
	_ = s.eventManager.Register(FlowPlugin{
		Id:          pluginID,
		Slug:        pluginSlug,
		SchemaInput: `{"mock": true}`,
	})

	s.pluginManagerMock.
		On("GetBySlug", mock.Anything, pluginSlug).
		Return(executorMock, nil)

	expectedError := fmt.Errorf("erro de execução")
	executorMock.
		On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, expectedError)

	// Configura o mock do repositório
	s.statusRepositoryMock.
		On("Save", mock.Anything, mock.Anything).
		Return(nil).
		Twice()

	_, err := s.eventManager.Execute(s.ctx, pluginID, map[string]any{"input": "value"})

	s.Error(err)
	s.Equal(expectedError, err)

	// Verifica se as métricas foram coletadas mesmo com erro
	metrics, exists := s.eventManager.metrics[pluginID]
	s.True(exists, "Métricas não foram coletadas para o plugin com erro")
	s.NotZero(metrics.ExecutionTime)
}

func (s *EventManagerTestSuite) TestExecute_ShouldHandlePanic() {
	const pluginSlug = "plugin-http"
	const pluginID = "test1"

	// Plugin e Executor Mock
	executorMock := new(PluginExecutorMock)
	_ = s.eventManager.Register(FlowPlugin{
		Id:          pluginID,
		Slug:        pluginSlug,
		SchemaInput: `{"mock": true}`,
	})

	s.pluginManagerMock.
		On("GetBySlug", mock.Anything, pluginSlug).
		Return(executorMock, nil)

	// Configura o mock para simular um panic
	executorMock.
		On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			panic("test panic")
		})

	// Configura o mock do repositório
	s.statusRepositoryMock.
		On("Save", mock.Anything, mock.Anything).
		Return(nil).
		Twice()

	// Executa o teste
	_, err := s.eventManager.Execute(s.ctx, pluginID, map[string]any{"input": "value"})

	// Verifica se o erro foi propagado corretamente
	s.Error(err)
	s.Contains(err.Error(), "panic recovered")

	// Verifica se as métricas foram coletadas
	metrics, exists := s.eventManager.metrics[pluginID]
	s.True(exists, "Métricas não foram coletadas para o plugin com panic")
	s.NotZero(metrics.ExecutionTime)
}

func (s *EventManagerTestSuite) TestExecute_ShouldHandleLongRunningPlugin() {
	const pluginSlug = "plugin-http"
	const pluginID = "test1"

	// Plugin e Executor Mock
	executorMock := new(PluginExecutorMock)
	_ = s.eventManager.Register(FlowPlugin{
		Id:          pluginID,
		Slug:        pluginSlug,
		SchemaInput: `{"mock": true}`,
	})

	s.pluginManagerMock.
		On("GetBySlug", mock.Anything, pluginSlug).
		Return(executorMock, nil)

	// Simula um plugin que demora 100ms para executar
	executorMock.
		On("Do", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			time.Sleep(100 * time.Millisecond)
		}).
		Return(map[string]any{"success": true}, nil)

	// Configura o mock do repositório
	s.statusRepositoryMock.
		On("Save", mock.Anything, mock.Anything).
		Return(nil).
		Twice()

	startTime := time.Now()
	_, err := s.eventManager.Execute(s.ctx, pluginID, map[string]any{"input": "value"})
	executionTime := time.Since(startTime)

	s.NoError(err)
	s.GreaterOrEqual(executionTime, 100*time.Millisecond, "O tempo de execução deve ser pelo menos 100ms")

	// Verifica se as métricas refletem o tempo de execução
	metrics, exists := s.eventManager.metrics[pluginID]
	s.True(exists)
	s.GreaterOrEqual(metrics.ExecutionTime, 100*time.Millisecond)
}
