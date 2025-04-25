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

func (e *EventManager) calculateNumberOfParentPlugins(firstPluginIdToExecute string) {
	e.numberOfPluginsToRun++

	if pluginInfo, ok := e.plugins[firstPluginIdToExecute]; ok {
		for _, slugNextToBeExecuted := range pluginInfo.NextToBeExecuted {
			e.calculateNumberOfParentPlugins(slugNextToBeExecuted)
		}
	}
}

func (e *EventManager) Execute(ctx *yctx.Context, firstPluginIdToExecute string, eventRequestData any) (finalResponse any, err error) {
	var (
		processResult       = make(chan EventManagerProcessResult)
		done                = make(chan struct{})
		pluginEventProducer = new(sync.Map)
	)

	e.calculateNumberOfParentPlugins(firstPluginIdToExecute)

	slog.Info("numberOfPluginsToRun", slog.Any("total", e.numberOfPluginsToRun))

	for slug, pluginInfo := range e.plugins {
		var pluginExecutor PluginExecutor

		pluginExecutor, err = e.pluginManager.GetBySlug(ctx, pluginInfo.Slug)
		if err != nil {
			slog.Error(
				"flow plugin get by slug",
				slog.Any("plugin_info", pluginInfo),
				slog.Any("error", err),
			)

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
			),
		)
	}

	if ch, ok := pluginEventProducer.Load(firstPluginIdToExecute); ok {
		ch.(chan<- any) <- eventRequestData
	}

	for i := 0; i < e.numberOfPluginsToRun; i++ {
		select {
		case result := <-processResult:
			slog.Info("plugin response", slog.Any("result", result))
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

				output, err = pluginExecutor.Do(ctx, pluginInfo.SchemaInput, body, nil)
				processResult <- EventManagerProcessResult{
					Id:     pluginInfo.Id,
					Output: output,
					Error:  err,
				}

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
