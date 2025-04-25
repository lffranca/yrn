package flowmanager

import (
	"github.com/yrn-go/yrn/pkg/yctx"
	"golang.org/x/exp/slog"
	"sync"
)

type (
	EventManagerProcessResult struct {
		Id     string
		Output any
		Error  error
	}
	EventProducerBody struct {
		Data         any            `json:"data"`
		SharedForAll map[string]any `json:"sharedForAll"`
	}
	EventManager struct {
		pluginManager        PluginManager
		plugins              map[string]FlowPlugin
		numberOfPluginsToRun int
	}
)

func NewEventManager(
	pluginManager PluginManager,
) *EventManager {
	plugins := make(map[string]FlowPlugin)

	return &EventManager{
		pluginManager,
		plugins,
		0,
	}
}

func (e *EventManager) Register(pluginExecutor FlowPlugin) (err error) {
	e.plugins[pluginExecutor.Id] = pluginExecutor
	return nil
}

func (e *EventManager) calculateNumberOfPluginsToRun(firstPluginIdToExecute string) {
	e.numberOfPluginsToRun++

	if pluginInfo, ok := e.plugins[firstPluginIdToExecute]; ok {
		for _, slugNextToBeExecuted := range pluginInfo.NextToBeExecuted {
			e.calculateNumberOfPluginsToRun(slugNextToBeExecuted)
		}
	}
}

func (e *EventManager) Execute(ctx *yctx.Context, firstPluginIdToExecute string, eventRequestData any) (finalResponse any, err error) {
	var (
		processResult        = make(chan EventManagerProcessResult)
		done                 = make(chan struct{})
		pluginEventProducer  = new(sync.Map)
		responseSharedForAll = new(sync.Map)
	)

	e.calculateNumberOfPluginsToRun(firstPluginIdToExecute)

	for slug, pluginInfo := range e.plugins {
		var pluginExecutor PluginExecutor

		pluginExecutor, err = e.pluginManager.GetBySlug(ctx, pluginInfo.Slug)
		if err != nil {
			return nil, err
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

	if ch, ok := pluginEventProducer.Load(firstPluginIdToExecute); ok {
		ch.(chan<- any) <- eventRequestData
	}

	for i := 0; i < e.numberOfPluginsToRun; i++ {
		select {
		case result := <-processResult:
			finalResponse = result.Output
			err = result.Error
		}
	}

	close(done)
	return
}

func (e *EventManager) handler(
	ctx *yctx.Context,
	pluginExecutor PluginExecutor,
	pluginInfo FlowPlugin,
	processResult chan<- EventManagerProcessResult,
	done chan struct{},
	pluginEventProducer *sync.Map,
	responseSharedForAll *sync.Map,
) chan<- any {
	var (
		eventProducer = make(chan any)
	)

	go func() {

		var parentPluginsExecuted int

		for {
			select {
			case body := <-eventProducer:
				parentPluginsExecuted++

				var (
					output any
					err    error
				)

				output, err = pluginExecutor.Do(ctx, pluginInfo.SchemaInput, body, syncMapToMap(responseSharedForAll))

				result := EventManagerProcessResult{
					Id:     pluginInfo.Id,
					Output: output,
					Error:  err,
				}

				processResult <- result

				responseSharedForAll.Store(pluginInfo.Id, output)

				for _, slugNextToBeExecuted := range pluginInfo.NextToBeExecuted {
					if ch, ok := pluginEventProducer.Load(slugNextToBeExecuted); ok {
						ch.(chan<- any) <- output
					}
				}
			case <-done:
				slog.Info("closing handler", slog.Any("plugin_id", pluginInfo.Id))
				return
			}
		}
	}()

	return eventProducer
}

func syncMapToMap(m *sync.Map) map[string]any {
	result := make(map[string]any)
	m.Range(func(key, value any) bool {
		strKey, ok := key.(string)
		if ok {
			result[strKey] = value
		}
		return true
	})
	return result
}
