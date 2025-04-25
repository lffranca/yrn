package flowmanager

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"sync"

	"github.com/yrn-go/yrn/pkg/yctx"
	"golang.org/x/exp/slog"
)

// EventManagerProcessResult representa o resultado do processamento de um plugin
type EventManagerProcessResult struct {
	Id     string
	Output any
	Error  error
}

// EventProducerBody representa o corpo da requisição para um plugin
type EventProducerBody struct {
	Data         any            `json:"data"`
	SharedForAll map[string]any `json:"sharedForAll"`
}

// PluginMetrics contém métricas de execução de um plugin
type PluginMetrics struct {
	StartTime     time.Time
	EndTime       time.Time
	MemoryBefore  uint64
	MemoryAfter   uint64
	CPUBefore     time.Duration
	CPUAfter      time.Duration
	ExecutionTime time.Duration
}

// PluginStatus representa o status atual de um plugin
type PluginStatus struct {
	PluginID   string
	Status     string
	StartTime  time.Time
	EndTime    time.Time
	Error      error
	Metrics    PluginMetrics
	Input      any
	Output     any
	SharedData map[string]any
}

// PluginStatusRepository define a interface para o repositório de status
type PluginStatusRepository interface {
	Save(ctx *yctx.Context, status PluginStatus) error
	GetByPluginID(ctx *yctx.Context, pluginID string) (PluginStatus, error)
	GetAll(ctx *yctx.Context) ([]PluginStatus, error)
}

// EventManager gerencia a execução de plugins em um fluxo
type EventManager struct {
	pluginManager        PluginManager
	plugins              map[string]FlowPlugin
	numberOfPluginsToRun int
	metrics              map[string]PluginMetrics
	statusRepo           PluginStatusRepository
}

// NewEventManager cria uma nova instância do EventManager
func NewEventManager(
	pluginManager PluginManager,
	statusRepo PluginStatusRepository,
) *EventManager {
	plugins := make(map[string]FlowPlugin)
	metrics := make(map[string]PluginMetrics)

	return &EventManager{
		pluginManager:        pluginManager,
		plugins:              plugins,
		numberOfPluginsToRun: 0,
		metrics:              metrics,
		statusRepo:           statusRepo,
	}
}

// Register registra um novo plugin no EventManager
func (e *EventManager) Register(pluginExecutor FlowPlugin) error {
	if pluginExecutor.Id == "" {
		return errors.New("plugin ID cannot be empty")
	}
	e.plugins[pluginExecutor.Id] = pluginExecutor
	return nil
}

// calculateNumberOfPluginsToRun calcula quantos plugins serão executados
func (e *EventManager) calculateNumberOfPluginsToRun(firstPluginIdToExecute string) {
	e.numberOfPluginsToRun++

	if pluginInfo, ok := e.plugins[firstPluginIdToExecute]; ok {
		for _, slugNextToBeExecuted := range pluginInfo.NextToBeExecuted {
			e.calculateNumberOfPluginsToRun(slugNextToBeExecuted)
		}
	}
}

// getMemoryUsage retorna o uso atual de memória em bytes
func getMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// getCPUUsage retorna o tempo de CPU usado
func getCPUUsage() time.Duration {
	return time.Duration(runtime.NumGoroutine()) * time.Millisecond
}

// savePluginStatus salva o status atual do plugin
func (e *EventManager) savePluginStatus(ctx *yctx.Context, pluginID string, status string, input, output any, err error, metrics PluginMetrics, sharedData map[string]any) error {
	pluginStatus := PluginStatus{
		PluginID:   pluginID,
		Status:     status,
		StartTime:  metrics.StartTime,
		EndTime:    metrics.EndTime,
		Error:      err,
		Metrics:    metrics,
		Input:      input,
		Output:     output,
		SharedData: sharedData,
	}

	return e.statusRepo.Save(ctx, pluginStatus)
}

// Execute inicia a execução do fluxo de plugins
func (e *EventManager) Execute(ctx *yctx.Context, firstPluginIdToExecute string, eventRequestData any) (any, error) {
	var (
		processResult        = make(chan EventManagerProcessResult, e.numberOfPluginsToRun)
		done                 = make(chan struct{})
		pluginEventProducer  = new(sync.Map)
		responseSharedForAll = new(sync.Map)
	)
	defer func() {
		// Log das métricas finais
		for pluginID, metrics := range e.metrics {
			slog.Info("plugin execution metrics",
				slog.String("plugin_id", pluginID),
				slog.Duration("execution_time", metrics.ExecutionTime),
				slog.Uint64("memory_usage_bytes", metrics.MemoryAfter-metrics.MemoryBefore),
				slog.Duration("cpu_usage", metrics.CPUAfter-metrics.CPUBefore))
		}
	}()

	e.calculateNumberOfPluginsToRun(firstPluginIdToExecute)

	// Inicializa os handlers para cada plugin
	for slug, pluginInfo := range e.plugins {
		pluginExecutor, err := e.pluginManager.GetBySlug(ctx, pluginInfo.Slug)
		if err != nil {
			slog.Error("failed to get plugin executor",
				slog.String("plugin_slug", pluginInfo.Slug),
				slog.Any("error", err))
			return nil, fmt.Errorf("failed to get plugin executor for %s: %w", pluginInfo.Slug, err)
		}

		pluginEventProducer.Store(
			slug,
			e.handler(
				ctx,
				pluginExecutor,
				pluginInfo,
				processResult,
				done,
				pluginEventProducer,
				responseSharedForAll,
			),
		)
	}

	// Inicia o fluxo com o primeiro plugin
	if ch, ok := pluginEventProducer.Load(firstPluginIdToExecute); ok {
		ch.(chan<- any) <- eventRequestData
	} else {
		return nil, fmt.Errorf("first plugin %s not found", firstPluginIdToExecute)
	}

	var (
		finalResponse any
		err           error
		completed     int
	)

	// Aguarda a conclusão de todos os plugins
	for completed < e.numberOfPluginsToRun {
		select {
		case result := <-processResult:
			completed++
			if result.Error != nil {
				slog.Error("plugin execution failed",
					slog.String("plugin_id", result.Id),
					slog.Any("error", result.Error))
				err = result.Error
			}
			finalResponse = result.Output
		}
	}

	close(done)
	return finalResponse, err
}

// handler gerencia a execução de um plugin específico
func (e *EventManager) handler(
	ctx *yctx.Context,
	pluginExecutor PluginExecutor,
	pluginInfo FlowPlugin,
	processResult chan<- EventManagerProcessResult,
	done chan struct{},
	pluginEventProducer *sync.Map,
	responseSharedForAll *sync.Map,
) chan<- any {
	eventProducer := make(chan any)

	go func() {
		var parentPluginsExecuted int

		for {
			select {
			case body := <-eventProducer:
				parentPluginsExecuted++

				// Inicia coleta de métricas
				startTime := time.Now()
				memoryBefore := getMemoryUsage()
				cpuBefore := getCPUUsage()

				var (
					output any
					err    error
				)

				// Salva status inicial
				metrics := PluginMetrics{
					StartTime:    startTime,
					MemoryBefore: memoryBefore,
					CPUBefore:    cpuBefore,
				}
				_ = e.savePluginStatus(ctx, pluginInfo.Id, "started", body, nil, nil, metrics, syncMapToMap(responseSharedForAll))

				// Executa o plugin com tratamento de panic
				func() {
					defer func() {
						if r := recover(); r != nil {
							slog.Error("panic in plugin handler",
								slog.String("plugin_id", pluginInfo.Id),
								slog.Any("recover", r))
							err = fmt.Errorf("panic recovered: %v", r)
						}
					}()

					output, err = pluginExecutor.Do(ctx, pluginInfo.SchemaInput, body, syncMapToMap(responseSharedForAll))
				}()

				// Finaliza coleta de métricas
				endTime := time.Now()
				memoryAfter := getMemoryUsage()
				cpuAfter := getCPUUsage()

				// Atualiza métricas
				metrics = PluginMetrics{
					StartTime:     startTime,
					EndTime:       endTime,
					MemoryBefore:  memoryBefore,
					MemoryAfter:   memoryAfter,
					CPUBefore:     cpuBefore,
					CPUAfter:      cpuAfter,
					ExecutionTime: endTime.Sub(startTime),
				}

				// Armazena métricas
				e.metrics[pluginInfo.Id] = metrics

				// Salva status final
				_ = e.savePluginStatus(ctx, pluginInfo.Id, "completed", body, output, err, metrics, syncMapToMap(responseSharedForAll))

				result := EventManagerProcessResult{
					Id:     pluginInfo.Id,
					Output: output,
					Error:  err,
				}

				processResult <- result

				if err == nil {
					responseSharedForAll.Store(pluginInfo.Id, output)

					for _, slugNextToBeExecuted := range pluginInfo.NextToBeExecuted {
						if ch, ok := pluginEventProducer.Load(slugNextToBeExecuted); ok {
							ch.(chan<- any) <- output
						} else {
							slog.Warn("next plugin not found",
								slog.String("current_plugin", pluginInfo.Id),
								slog.String("next_plugin", slugNextToBeExecuted))
						}
					}
				}
			case <-done:
				slog.Info("closing handler",
					slog.String("plugin_id", pluginInfo.Id),
					slog.Int("executions", parentPluginsExecuted))
				return
			}
		}
	}()

	return eventProducer
}

// syncMapToMap converte um sync.Map para map[string]any
func syncMapToMap(m *sync.Map) map[string]any {
	result := make(map[string]any)
	m.Range(func(key, value any) bool {
		if strKey, ok := key.(string); ok {
			result[strKey] = value
		}
		return true
	})
	return result
}
